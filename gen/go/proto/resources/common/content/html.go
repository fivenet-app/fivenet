package content

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

func (n *RichTextHtmlNode) populateFrom(htmlNode *html.Node) error {
	if htmlNode.Parent == nil {
		n.Type = NodeType_NODE_TYPE_DOC
	} else {
		n.Type = NodeType_NODE_TYPE_ELEMENT
	}

	switch htmlNode.Type {
	case html.ElementNode:
		n.Tag = htmlNode.Data

	case html.DocumentNode:

	default:
		return errors.New("given node needs to be an element or document")
	}

	if len(htmlNode.Attr) > 0 {
		n.Attrs = make(map[string]string)
		var a html.Attribute
		for _, a = range htmlNode.Attr {
			key := strings.ToLower(a.Key)
			switch key {
			case "id":
				// Skip empty id attribute
				if strings.TrimSpace(a.Val) == "" {
					continue
				}

				n.Id = &a.Val

			case "style":
				// Skip empty style attribute
				val := strings.TrimSpace(a.Val)
				if val == "" {
					continue
				}

				// Clean style options - maybe sort them in the future?
				if !strings.HasSuffix(val, ";") {
					val += ";"
				}

				n.Attrs[key] = val

			case "class":
				// Skip empty class attribute
				if strings.TrimSpace(a.Val) == "" {
					continue
				}

				fallthrough

			default:
				// Don't skip empty valued attributes as they might have a meaning on the frontend-side
				n.Attrs[key] = a.Val
			}
		}
	}

	e := htmlNode.FirstChild
	for e != nil {
		switch e.Type {
		case html.TextNode:
			n.Type = NodeType_NODE_TYPE_TEXT

			data := strings.TrimSpace(e.Data)
			// If text data is not empty after timming spaces
			if len(data) > 0 {
				n.Content = append(n.Content, &RichTextHtmlNode{
					Type: NodeType_NODE_TYPE_TEXT,
					Text: &e.Data,
				})
			}

		case html.ElementNode:
			if n.Content == nil {
				n.Content = make([]*RichTextHtmlNode, 0)
			}

			jsonElemNode := &RichTextHtmlNode{}
			if err := jsonElemNode.populateFrom(e); err != nil {
				return err
			}

			n.Content = append(n.Content, jsonElemNode)
		}

		e = e.NextSibling
	}

	if len(n.GetTag()) == 2 && utils.IsHeaderTag(n.GetTag()) {
		// Either empty id or "broken" id tag
		if n.Id == nil || n.GetId() == "" || utils.IsHeaderTag(n.GetId()) {
			if n.Text != nil && n.GetText() != "" {
				id := utils.SlugNoDots(fmt.Sprintf("%s-%s", n.GetTag(), n.GetText()))
				n.Id = &id
			} else if len(n.GetContent()) > 0 {
				id := utils.SlugNoDots(fmt.Sprintf("%s-%s", n.GetTag(), walkContentForText(n.GetContent())))
				n.Id = &id
			}
		}
	}

	return nil
}

func walkContentForText(ns []*RichTextHtmlNode) string {
	text := ""
	for i := range ns {
		element := ns[i]
		if element.Text == nil || element.GetText() == "" {
			text += walkContentForText(element.GetContent())
		} else {
			text += element.GetText()
			break
		}
	}

	return text
}

func (n *RichTextHtmlNode) populateTo(htmlNode *html.Node) {
	if n.GetTag() != "" {
		htmlNode.Data = n.GetTag()
		htmlNode.Type = html.ElementNode
	} else {
		htmlNode.Type = html.DocumentNode
	}

	if n.Id != nil && n.GetId() != "" {
		// Make sure that headers have id
		if len(n.GetTag()) == 2 && utils.IsHeaderTag(n.GetTag()) {
			// Either empty id or "broken" id tag
			if n.GetId() == "" || utils.IsHeaderTag(n.GetId()) {
				if n.Text != nil && n.GetText() != "" {
					id := utils.SlugNoDots(fmt.Sprintf("%s-%s", n.GetTag(), n.GetText()))
					n.Id = &id
				} else if len(n.GetContent()) > 0 {
					id := utils.SlugNoDots(fmt.Sprintf("%s-%s", n.GetTag(), walkContentForText(n.GetContent())))
					n.Id = &id
				}
			}
		}

		htmlNode.Attr = append(htmlNode.Attr, html.Attribute{
			Key: "id",
			Val: n.GetId(),
		})
	}

	keys := make([]string, 0, len(n.GetAttrs()))
	for k := range n.GetAttrs() {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		htmlNode.Attr = append(htmlNode.Attr, html.Attribute{
			Key: k,
			Val: n.GetAttrs()[k],
		})
	}

	if n.Text != nil && n.GetText() != "" {
		htmlNode.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: n.GetText(),
		})
	}

	for _, e := range n.GetContent() {
		htmlElem := &html.Node{}
		e.populateTo(htmlElem)
		htmlNode.AppendChild(htmlElem)
	}
}

func ParseHTML(in string) (*html.Node, error) {
	d, err := html.Parse(strings.NewReader(in))
	if err != nil {
		return nil, err
	}

	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			body.Data = "body"
			body.Parent = nil
			body.PrevSibling = nil
			body.NextSibling = nil
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(d.LastChild)

	if body != nil {
		return body, nil
	}

	return d, nil
}

func FromHTML(in string) (*RichTextHtmlNode, error) {
	h, err := ParseHTML(in)
	if err != nil {
		return nil, err
	}

	out, err := FromHTMLNode(h)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// FromHTMLNode converts html.Node to RichTextHtmlNode.
func FromHTMLNode(node *html.Node) (*RichTextHtmlNode, error) {
	jNode := &RichTextHtmlNode{}
	if err := jNode.populateFrom(node); err != nil {
		return nil, err
	}

	return jNode, nil
}

// ToHTMLNode converts RichTextHtmlNode to html.Node.
func (n *RichTextHtmlNode) ToHTMLNode() (*html.Node, error) {
	node := &html.Node{}

	n.populateTo(node)

	return node, nil
}

// ToHTML HTML potentially not pretty.
func (n *RichTextHtmlNode) ToHTML() (string, error) {
	h, err := n.ToHTMLNode()
	if err != nil {
		return "", err
	}

	htmlBuff := &bytes.Buffer{}
	if err := html.Render(htmlBuff, h); err != nil {
		return "", err
	}

	out := htmlBuff.String()
	out = strings.ReplaceAll(out, "<body>", "")
	out = strings.ReplaceAll(out, "</body>", "")

	return out, nil
}

func PrettyHTML(in string) (string, error) {
	in = strings.ReplaceAll(in, "<body>\n", "")
	in = strings.ReplaceAll(in, "\n</body>", "")

	return gohtml.Format(in), nil
}

func ExtractFromHTML(doc *RichTextHtmlNode) *ExtractedContent {
	if doc == nil {
		return &ExtractedContent{}
	}

	h, err := doc.ToHTMLNode()
	if err != nil {
		return &ExtractedContent{}
	}

	textBuff := &strings.Builder{}
	wordCount := uint32(0)
	var crawler func(*html.Node)
	firstHeading := ""
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && utils.IsHeaderTag(node.Data) && firstHeading == "" {
			// Extract heading text
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.TextNode {
					firstHeading = strings.TrimSpace(child.Data)
					break
				}
			}
		}

		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			if text != "" {
				textBuff.WriteString(text + " ")
				wordCount += uint32(len(strings.Fields(text)))
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(h)

	text := strings.TrimSpace(textBuff.String())

	return &ExtractedContent{
		Text:         text,
		WordCount:    wordCount,
		FirstHeading: firstHeading,
	}
}
