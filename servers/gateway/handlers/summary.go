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
	URL string `json:"url"`
	SecureURL string `json:"secureUrl"`
	Type string `json:"type"`
	Width int `json:"width"`
	Height int `json:"height"`
	Alt string `json:"alt"`
}

type PreviewVideo struct {
	URL string `json:"url"`
	SecureURL string `json:"secureUrl"`
	Type string `json:"type"`
	Width int `json:"width"`
	Height int `json:"height"`
}

type Metadata struct {
	Type string `json:"type"`
	URL string `json:"url"`
	Title string `json:"title"`
	SiteName string `json:"siteName"`
	Description string `json:"description"`
	Author string `json:"author"`
	Keywords []string `json:"keywords"`
	Icon PreviewImage `json:"previewImage"`
	Images []*PreviewImage `json:"images"`
	Videos []*PreviewVideo `json:"videos"`
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
	w.Header().Set("Content-Type", "application/json")
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
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
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
		parseMetaTag(token, metadata)
	}
	if token.Data == "link" {
		parseLinkTag(token, metadata)
	}
}

func parseMetaTag(token *html.Token, metadata *Metadata) {
	// extract meta tag attributes
	property, content, name := extractMetaAttributes(token)
	// check meta attribute values
	if property == "twitter:card" && metadata.Type == "" {
		metadata.Type = content
	} else if property == "og:type" {
		metadata.Type = content
	} else if property == "og:title" {
		metadata.Title = content
	} else if property == "twitter:title" && metadata.Title == "" {
		metadata.Title = content
	} else if property == "og:url" {
		metadata.URL = content
	} else if property == "og:site_name" {
		metadata.SiteName = content
	} else if property == "twitter:description" && metadata.Description == "" {
		metadata.Description = content
	} else if property == "og:description" {
		metadata.Description = content
	} else if name == "description" {
		if metadata.Description == "" {
			metadata.Description = content
		}
	} else if name == "author" {
		metadata.Author = content
	} else if name == "keywords" {
		keywords := strings.Split(content, ",")
		for i, keyword := range keywords {
			keywords[i] = strings.Trim(keyword, " ")
		}
		metadata.Keywords = keywords
	}
	if strings.HasPrefix(property, "og:image") || strings.HasPrefix(property, "twitter:image") {
		parseOGImage(property, content, metadata)
	}
	if strings.HasPrefix(property, "og:video") {
		parseOGVideo(property, content, metadata)
	}
}

func extractMetaAttributes(token *html.Token) (string, string, string) {
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
	return property, content, name
}

func parseOGImage(property string, content string, metadata *Metadata) {
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

func parseOGVideo(property string, content string, metadata *Metadata) {
	split := strings.Split(property, ":")
	if len(split) == 2 {
		newVideo := &PreviewVideo{ URL: content }
		metadata.Videos = append(metadata.Videos, newVideo)
		return
	}
	suffix := split[2]
	latestVideo := metadata.Videos[len(metadata.Videos) - 1]
	if suffix == "url" {
		latestVideo.URL = content
	}
	if suffix == "secure_url" {
		latestVideo.SecureURL = content
	}
	if suffix == "type" {
		latestVideo.Type = content
	}
	if suffix == "width" {
		newInt, err := strconv.Atoi(content)
		if err != nil {
			log.Fatal("image width must be an int")
		}
		latestVideo.Width = newInt
	}
	if suffix == "height" {
		newInt, err := strconv.Atoi(content)
		if err != nil {
			log.Fatal("image height must be an int")
		}
		latestVideo.Height = newInt
	}
}

func parseLinkTag(token *html.Token, metadata *Metadata) {
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