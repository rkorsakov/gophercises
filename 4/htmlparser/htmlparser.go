package htmlparser

import (
	"golang.org/x/net/html"
	"io"
	"log"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func (l Link) String() string {
	return "Link{\n" +
		"  Href: " + "\"" + l.Href + "\",\n" +
		"  Text: " + "\"" + l.Text + "\",\n" +
		"}"
}

func ParseHTML(r io.Reader) []Link {
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	var links []Link
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var link Link
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link.Href = attr.Val
					break
				}
			}
			var textBuilder strings.Builder
			var collectText func(*html.Node)
			collectText = func(node *html.Node) {
				if node.Type == html.TextNode {
					text := strings.TrimSpace(node.Data)
					if text != "" {
						if textBuilder.Len() > 0 {
							textBuilder.WriteString(" ")
						}
						textBuilder.WriteString(text)
					}
				}
				for child := node.FirstChild; child != nil; child = child.NextSibling {
					if child.Type != html.ElementNode || child.Data != "a" {
						collectText(child)
					}
				}
			}
			collectText(n)
			link.Text = strings.Join(strings.Fields(textBuilder.String()), " ")
			links = append(links, link)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}
