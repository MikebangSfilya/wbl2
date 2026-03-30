package parser

import (
	"bytes"
	"net/url"

	"golang.org/x/net/html"
)

func ExtractLinks(body []byte, baseURL *url.URL) ([]string, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var links []string
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a", "link":
				links = appendIfValid(links, getAttr(n, "href"), baseURL)
			case "img", "script":
				links = appendIfValid(links, getAttr(n, "src"), baseURL)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, nil
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func appendIfValid(links []string, val string, baseURL *url.URL) []string {
	if val == "" {
		return links
	}
	parsed, err := url.Parse(val)
	if err != nil {
		return links
	}
	resolved := baseURL.ResolveReference(parsed)
	return append(links, resolved.String())
}
