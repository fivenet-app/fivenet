package content

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/pkg/html/htmlsanitizer"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
	"google.golang.org/protobuf/encoding/protojson"
)

type NodeType string

const (
	Version_v0 = "v0"

	DocNodeType NodeType = "doc"

	ElementNodeType NodeType = "element"
	TextNodeType    NodeType = "text"
	CommentNodeType NodeType = "comment"
)

// Scan implements driver.Valuer for protobuf Content.
func (x *Content) Scan(value any) error {
	switch t := value.(type) {
	case string:
		if strings.HasPrefix(t, "{") {
			if err := protojson.Unmarshal([]byte(t), x); err != nil {
				return err
			}

			return x.Populate()
		}

		h, err := ParseHTML(t)
		if err != nil {
			return err
		}

		out, err := FromHTMLNode(h)
		if err != nil {
			return err
		}

		v := Version_v0
		x.Version = &v
		x.Content = out
		x.RawContent = &t

	case []byte:
		if strings.HasPrefix(string(t), "{") {
			if err := protojson.Unmarshal(t, x); err != nil {
				return err
			}

			return x.Populate()
		}

		hRaw := string(t)
		h, err := ParseHTML(hRaw)
		if err != nil {
			return err
		}

		out, err := FromHTMLNode(h)
		if err != nil {
			return err
		}

		v := Version_v0
		x.Version = &v
		x.Content = out
		x.RawContent = &hRaw

	default:
		return fmt.Errorf("invalid format for content")
	}

	return nil
}

// Value marshals the value into driver.Valuer.
func (x *Content) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	// If the raw content isn't nil, need to "encode" it to `JSONNode` for the `Content` field
	if x.RawContent != nil && *x.RawContent != "" {
		h, err := ParseHTML(*x.RawContent)
		if err != nil {
			return nil, err
		}

		x.Content, err = FromHTMLNode(h)
		if err != nil {
			return nil, err
		}
	}

	return protoutils.Marshal(&Content{
		Version: x.Version,
		Content: x.Content,
	})
}

func (x *Content) Populate() error {
	out, err := x.Content.ToHTMLP()
	if err != nil {
		return err
	}

	x.RawContent = &out

	return nil
}

func (n *JSONNode) populateFrom(htmlNode *html.Node) error {
	if htmlNode.Parent == nil {
		n.Type = string(DocNodeType)
	} else {
		n.Type = string(ElementNodeType)
	}

	switch htmlNode.Type {
	case html.ElementNode:
		n.Tag = htmlNode.Data
		break

	case html.DocumentNode:
		break

	default:
		return fmt.Errorf("given node needs to be an element or document")
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

				n.Id = a.Val

			case "style":
				// Skip empty style attribute
				if strings.TrimSpace(a.Val) == "" {
					continue
				}

				// Clean style options - maybe sort them in the future?
				val := strings.Replace(a.Val, " ", "", -1)
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
			n.Type = string(TextNodeType)

			data := strings.TrimSpace(e.Data)
			// If text data is not empty after timming spaces
			if len(data) > 0 {
				n.Content = append(n.Content, &JSONNode{
					Type: string(TextNodeType),
					Text: e.Data,
				})
			}

		case html.ElementNode:
			if n.Content == nil {
				n.Content = make([]*JSONNode, 0)
			}

			jsonElemNode := &JSONNode{}
			if err := jsonElemNode.populateFrom(e); err != nil {
				return err
			}

			n.Content = append(n.Content, jsonElemNode)
		}

		e = e.NextSibling
	}

	if strings.HasPrefix(n.Tag, "h") {
		if n.Id == "" {
			if n.Text != "" {
				n.Id = utils.SlugNoDots(fmt.Sprintf("%s-%s", n.Tag, n.Text))
			} else if len(n.Content) > 0 {
				n.Id = utils.SlugNoDots(fmt.Sprintf("%s-%s", n.Tag, n.Content[0].Text))
			}
		}
	}

	return nil
}

func (n *JSONNode) populateTo(htmlNode *html.Node) {
	if n.Tag != "" {
		htmlNode.Data = n.Tag
		htmlNode.Type = html.ElementNode
	} else {
		htmlNode.Type = html.DocumentNode
	}

	if n.Id != "" {
		htmlNode.Attr = append(htmlNode.Attr, html.Attribute{
			Key: "id",
			Val: n.Id,
		})
	}

	keys := make([]string, 0, len(n.Attrs))
	for k := range n.Attrs {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		htmlNode.Attr = append(htmlNode.Attr, html.Attribute{
			Key: k,
			Val: n.Attrs[k],
		})
	}

	if n.Text != "" {
		htmlNode.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: n.Text,
		})
	}

	for _, e := range n.Content {
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

// FromHTMLNode
func FromHTMLNode(node *html.Node) (*JSONNode, error) {
	jNode := &JSONNode{}
	if err := jNode.populateFrom(node); err != nil {
		return nil, err
	}

	return jNode, nil
}

// ToHTMLNode
func (n *JSONNode) ToHTMLNode() (*html.Node, error) {
	node := &html.Node{}

	n.populateTo(node)

	return node, nil
}

// ToHTML HTML potentially not pretty
func (n *JSONNode) ToHTML() (string, error) {
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

// ToHTMLP Pretty HTML
func (n *JSONNode) ToHTMLP() (string, error) {
	h, err := n.ToHTML()
	if err != nil {
		return "", err
	}

	h = strings.ReplaceAll(h, "<body>\n", "")
	h = strings.ReplaceAll(h, "\n</body>", "")

	return gohtml.Format(h), nil
}

func PrettyHTML(in string) (string, error) {
	in = strings.ReplaceAll(in, "<body>\n", "")
	in = strings.ReplaceAll(in, "\n</body>", "")

	return gohtml.Format(in), nil
}

func (c *Content) GetSummary(length int) string {
	if c.RawContent == nil {
		return ""
	}

	return utils.StringFirstN(htmlsanitizer.StripTags(*c.RawContent), length)
}
