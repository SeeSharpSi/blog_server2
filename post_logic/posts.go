package post_logic

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Post struct {
	Id      int
	Title   string
	Content string
}

// Populates the Title and Content fields of a post using html 
// The title is based off of the first h1 header in the html file 
func (p *Post) Parse(html_content string) {
	r := strings.NewReader(html_content)

	doc, err := html.Parse(r) 
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err) 
	}

	h1Text, ok := getH1Text(doc)
	if !ok {
		fmt.Println("No <h1> tag found.")
	} else {
		p.Title = h1Text
		p.Content = html_content
	}
}

// getH1Text finds the first h1 node and returns its text content.
func getH1Text(doc *html.Node) (string, bool) {
	h1Node, ok := findFirstH1(doc)
	if !ok {
		return "", false
	}
	return extractText(h1Node), true
}

func findFirstH1(n *html.Node) (*html.Node, bool) {
	if n.Type == html.ElementNode && n.Data == "h1" {
		return n, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if resultNode, ok := findFirstH1(c); ok {
			return resultNode, true
		}
	}

	return nil, false
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.WriteString(extractText(c))
	}
	return b.String()
}
