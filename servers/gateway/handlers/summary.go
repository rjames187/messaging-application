package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type PreviewImage struct {
	URL string 
	SecureURL string
	Type string
	Width int
	Height int
	Alt string
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
	Images []*PreviewImage
}

func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	metadata, err := fetchHTML(url)
	if err != nil {
		if strings.HasPrefix(err.Error(), "error fetching html:") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
	}
	enc := json.NewEncoder(w)
	enc.Encode(metadata)
}

func fetchHTML(url string) (*Metadata, error) {
	resp, err := http.Get(url)
	if err != nil {
		message := fmt.Sprintf("error fetching html: %v", err)
		return nil, errors.New(message)
	}
	defer resp.Body.Close()
	metadata := extractSummary(resp)
	return metadata, nil
}

func extractSummary(resp *http.Response) *Metadata {
	tokenizer := html.NewTokenizer(resp.Body)
	metadata := &Metadata{}
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			log.Fatalf("error tokenizing HTML: %v", err)
		}

		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if token.Data == "head" {
				break
			}
			if token.Data == "title" {
				tokenType := tokenizer.Next()
				if tokenType == html.TextToken {
					token := tokenizer.Token()
					if metadata.Title == "" {
						metadata.Title = token.Data
					}
				}
				continue
			}
			parseToken(&token, metadata)
		}
	}
	return metadata
}

func parseToken(token *html.Token, metadata *Metadata) { 
	if token.Data == "meta" {
		// extract meta tag attributes
		var property string
		var content string
		var name string
		for _, attr := range token.Attr {
			if attr.Key == "property" {
				property = attr.Val
			}
			if attr.Key == "content" {
				content = attr.Val
			}
			if attr.Key == "name" {
				name = attr.Val
			}
		}
		// check meta attribute values
		if property == "og:type" {
			metadata.Type = content
		}
		if property == "og:url" {
			metadata.URL = content
		}
		if property == "og:site_name" {
			metadata.SiteName = content
		}
		if property == "og:description" {
			metadata.Description = content
		}
		if name == "description" {
			if metadata.Description == "" {
				metadata.Description = content
			}
		}
		if name == "author" {
			metadata.Author = content
		}
		if name == "keywords" {
			keywords := strings.Split(content, ",")
			for i, keyword := range keywords {
				keywords[i] = strings.Trim(keyword, " ")
			}
			metadata.Keywords = keywords
		}
		if strings.HasPrefix(property, "og:image") {
			split := strings.Split(property, ":")
			if len(split) == 2 {
				newImage := &PreviewImage{ URL: content }
				metadata.Images = append(metadata.Images, newImage)
				return
			}
			suffix := split[2]
			latestImage := metadata.Images[len(metadata.Images) - 1]
			if suffix == "url" {
				latestImage.URL = content
			}
			if suffix == "secure_url" {
				latestImage.SecureURL = content
			}
			if suffix == "type" {
				latestImage.Type = content
			}
			if suffix == "width" {
				newInt, err := strconv.Atoi(content)
				if err != nil {
					log.Fatal("image width must be an int")
				}
				latestImage.Width = newInt
			}
			if suffix == "height" {
				newInt, err := strconv.Atoi(content)
				if err != nil {
					log.Fatal("image height must be an int")
				}
				latestImage.Height = newInt
			}
			if suffix == "alt" {
				latestImage.Alt = content
			}
		}
	}
	if token.Data == "link" {
		// extract link tag attributes
		var rel string
		var href string
		var sizes string
		var tipe string
		for _, attr := range token.Attr {
			if attr.Key == "rel" {
				rel = attr.Val
			}
			if attr.Key == "href" {
				href = attr.Val
			}
			if attr.Key == "sizes" {
				sizes = attr.Val
			}
			if attr.Key == "type" {
				tipe = attr.Val
			}
		}
		if rel == "icon" {
			sizesSlice := strings.Split(sizes, "x")
			height, err := strconv.Atoi(sizesSlice[0])
			if err != nil {
				log.Fatal("can't convert icon height to int")
			}
			width, err := strconv.Atoi(sizesSlice[1])
			if err != nil {
				log.Fatal("can't convert icon width to int")
			}
			newImage := PreviewImage{
				URL: href,
				Type: tipe,
				Height: height,
				Width: width,
			}
			metadata.Icon = newImage
		}
	}
}