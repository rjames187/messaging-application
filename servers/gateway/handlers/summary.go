package handlers

import (
	"net/http"

	"golang.org/x/net/html"
)

type PreviewImage struct {
	HREF string
	Sizes string
	Type string
}

type Metadata struct {
	Type string
	URL string
	Title string
	SiteName string
	Description string
	Author string
	Keywords []string
	Icon PreviewImage
	Images []PreviewImage
}

func SummaryHandler(w http.ResponseWriter, r *http.Request) {

}

func fetchHTML(url string) {
	
}

func extractSummary(resp *http.Response) {
	tokenizer := html.NewTokenizer(resp.Body)
	for {

	}
}