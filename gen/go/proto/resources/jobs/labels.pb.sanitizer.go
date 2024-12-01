// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/jobs/labels.proto

package jobs

import (
	"github.com/fivenet-app/fivenet/pkg/html/htmlsanitizer"
)

func (m *Label) Sanitize() error {

	m.Color = htmlsanitizer.StripTags(m.Color)

	return nil
}