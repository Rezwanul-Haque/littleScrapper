package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// WebScraper represents a web scraping utility.
type WebScraper struct {
	URL           string
	SaveDirectory string
}

// NewWebScraper creates a new WebScraper instance.
func NewWebScraper(url, saveDirectory string) *WebScraper {
	return &WebScraper{
		URL:           url,
		SaveDirectory: saveDirectory,
	}
}

// FetchHTMLContent fetches the HTML content of the web page.
func (ws *WebScraper) FetchHTMLContent() (*goquery.Document, error) {
	response, err := http.Get(ws.URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// ExtractTitleAndBody extracts the title and body text from the HTML content.
func (ws *WebScraper) ExtractTitleAndBody(doc *goquery.Document) (string, string) {
	title := doc.Find("title").Text()
	body := strings.Builder{}
	doc.Find("p").Each(func(index int, element *goquery.Selection) {
		body.WriteString(element.Text())
		body.WriteString(" ")
	})

	return title, body.String()
}

// ScrapeImages scrapes and downloads images from the web page.
func (ws *WebScraper) ScrapeImages(doc *goquery.Document) error {
	doc.Find("img").Each(func(index int, element *goquery.Selection) {
		imgURL, _ := element.Attr("src")
		if imgURL != "" {
			imgResponse, err := http.Get(imgURL)
			if err != nil {
				fmt.Printf("Failed to download image from %s: %v\n", imgURL, err)
				return
			}
			defer imgResponse.Body.Close()

			// Extract a valid file name from the URL
			imgName := path.Base(imgURL)
			imgName = ws.sanitizeFilename(imgName)

			imgPath := path.Join(ws.SaveDirectory, imgName)
			imgFile, err := os.Create(imgPath)
			if err != nil {
				fmt.Printf("Failed to create image file %s: %v\n", imgPath, err)
				return
			}
			defer imgFile.Close()

			_, err = io.Copy(imgFile, imgResponse.Body)
			if err != nil {
				fmt.Printf("Failed to save image to file %s: %v\n", imgPath, err)
			}
		}
	})

	return nil
}

// SanitizeFilename sanitizes a file name to make it valid and clean.
func (ws *WebScraper) sanitizeFilename(filename string) string {
	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// Remove invalid characters using a regular expression
	re := regexp.MustCompile("[^a-zA-Z0-9_.-]")
	filename = re.ReplaceAllString(filename, "")

	return filename
}

func main() {
	url := "https://www.theguardian.com/politics/2018/aug/19/brexit-tory-mps-warn-of-entryism-threat-from-leave-eu-supporters"
	saveDirectory := "images"

	ws := NewWebScraper(url, saveDirectory)

	doc, err := ws.FetchHTMLContent()
	if err != nil {
		fmt.Printf("failed to fetch content from %s: %v\n", ws.URL, err)
		return
	}

	title, body := ws.ExtractTitleAndBody(doc)
	fmt.Println("Title:", title)
	fmt.Println("Body:", body)

	if err := os.MkdirAll(ws.SaveDirectory, os.ModePerm); err != nil {
		fmt.Printf("failed to create directory %s: %v\n", ws.SaveDirectory, err)
		return
	}

	if err := ws.ScrapeImages(doc); err != nil {
		fmt.Printf("failed to scrape and download images: %v\n", err)
	}
}
