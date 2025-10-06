package post_logic

import (
	"strings"

	"golang.org/x/net/html"
)

type Post struct {
	Id      int
	Title   string
	Content string
}

// should get the first (or only) h1 in the html and make it the post's title
// ID should also be the index of the file (like, the 2nd file in the list)
func (p *Post) Parse(doc *html.Node) {
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
