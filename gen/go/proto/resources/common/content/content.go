package content

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
)

// Scan implements driver.Valuer for protobuf Content.
func (x *Content) Scan(value any) error {
	var data string
	switch t := value.(type) {
	case string:
		data = t

	case []byte:
		data = string(t)

	default:
		return errors.New("invalid format for content")
	}

	if strings.HasPrefix(data, "<") {
		h, err := ParseHTML(data)
		if err != nil {
			return err
		}

		out, err := FromHTMLNode(h)
		if err != nil {
			return err
		}
		x.Content = out

		return nil
	}

	if err := protoutils.UnmarshalPartialJSON([]byte(data), x); err != nil {
		return err
	}

	// Legacy JSON content stored directly in Content.Content field
	if x.GetContent() != nil {
		x.Content = x.GetContent()

		// Make sure to make raw HTML available
		rawHtml, err := x.GetContent().ToHTML()
		if err != nil {
			return err
		}
		x.RawHtml = &rawHtml
	}

	return nil
}

// Value marshals the value into driver.Valuer.
func (x *Content) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	if s := x.GetRawHtml(); s != "" {
		h, err := ParseHTML(s)
		if err != nil {
			return nil, err
		}

		out, err := FromHTMLNode(h)
		if err != nil {
			return nil, err
		}
		x.Content = out
		x.RawHtml = nil
	}

	if x.GetTiptapJson() != nil {
		return protoutils.MarshalToJSON(&Content{
			Version:     "tiptap/v1",
			ContentType: ContentType_CONTENT_TYPE_TIPTAP_JSON,
			TiptapJson:  x.GetTiptapJson(),
		})
	} else if x.GetContent() != nil {
		// Legacy JSON content stored directly in Content.Content field
		return protoutils.MarshalToJSON(&Content{
			Version:     "legacy_json/v1",
			ContentType: ContentType_CONTENT_TYPE_HTML,
			Content:     x.GetContent(),
		})
	}

	return nil, errors.New("unsupported content type for Value()")
}

func (x *Content) Extract() *ExtractedContent {
	if x.TiptapJson != nil {
		return ExtractFromTiptap(x.GetTiptapJson())
	} else if x.Content != nil {
		return ExtractFromHTML(x.GetContent())
	}

	return &ExtractedContent{}
}

func (x *ExtractedContent) GetSummary(length int) string {
	if x == nil {
		return ""
	}

	return utils.SummaryFromText(x.Text, length)
}
