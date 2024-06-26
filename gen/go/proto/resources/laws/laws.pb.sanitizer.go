// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/laws/laws.proto

package laws

import (
	"github.com/fivenet-app/fivenet/pkg/htmlsanitizer"
)

func (m *Law) Sanitize() error {

	if m.Description != nil {
		*m.Description = htmlsanitizer.Sanitize(*m.Description)
	}

	m.Name = htmlsanitizer.Sanitize(m.Name)

	return nil
}

func (m *LawBook) Sanitize() error {

	if m.Description != nil {
		*m.Description = htmlsanitizer.Sanitize(*m.Description)
	}

	m.Name = htmlsanitizer.Sanitize(m.Name)

	return nil
}
