package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var url string
	fmt.Print("Enter the URL you want to scrape: ")
	_, err := fmt.Scanln(&url)
	if err != nil {
		fmt.Printf("failed to read URL: %v\n", err)
		return
	}

	saveDirectory := "images"

	// Fetch the web page content
	doc, err := getHTMLContent(url)
	if err != nil {
		fmt.Printf("failed to fetch content from %s: %v\n", url, err)
		return
	}

	// Extract title and body text
	title, body := extractTitleAndBody(doc)
	fmt.Println("Title:", title)
	fmt.Println("Body:", body)

	// Create the save directory
	if err := os.MkdirAll(saveDirectory, os.ModePerm); err != nil {
		fmt.Printf("failed to create directory %s: %v\n", saveDirectory, err)
		return
	}

	// Scrape and download images
	if err := scrapeImages(doc, saveDirectory); err != nil {
		fmt.Printf("failed to scrape and download images: %v\n", err)
	}
}

// getHTMLContent to fetch web page content
func getHTMLContent(url string) (*goquery.Document, error) {
	response, err := http.Get(url)
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

// extractTitleAndBody to extract the title and body text from the HTML document
func extractTitleAndBody(doc *goquery.Document) (string, string) {
	title := doc.Find("title").Text()
	body := strings.Builder{}
	doc.Find("p").Each(func(index int, element *goquery.Selection) {
		body.WriteString(element.Text())
		body.WriteString(" ")
	})

	return strings.TrimSpace(title), strings.TrimSpace(body.String())
}

// scrapeImages to scrape and download images
func scrapeImages(doc *goquery.Document, saveDirectory string) error {
	doc.Find("img").Each(func(index int, element *goquery.Selection) {
		imgURL, _ := element.Attr("src")
		if imgURL != "" {
			imgResponse, err := http.Get(imgURL)
			if err != nil {
				fmt.Printf("failed to download image from %s: %v\n", imgURL, err)
				return
			}
			defer imgResponse.Body.Close()

			imgName := sanitizeFilename(imgURL)
			imgPath := path.Join(saveDirectory, imgName)
			imgFile, err := os.Create(imgPath)
			if err != nil {
				fmt.Printf("failed to create image file %s: %v\n", imgPath, err)
				return
			}
			defer imgFile.Close()

			_, err = io.Copy(imgFile, imgResponse.Body)
			if err != nil {
				fmt.Printf("failed to save image to file %s: %v\n", imgPath, err)
			}
		}
	})

	return nil
}

// sanitizeFilename to sanitize a file name
func sanitizeFilename(imgURL string) string {
	parsedURL, err := url.Parse(imgURL)
	if err != nil {
		// Handle parsing error
		return ""
	}

	filename := path.Base(parsedURL.Path)

	re := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
	filename = re.ReplaceAllString(filename, "_")

	return filename
}
