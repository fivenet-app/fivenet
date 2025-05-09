// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/mailer/settings.proto

package mailer

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/html/htmlsanitizer"
)

func (m *EmailSettings) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: BlockedEmails
	for idx, item := range m.BlockedEmails {
		_, _ = idx, item

		m.BlockedEmails[idx] = htmlsanitizer.StripTags(m.BlockedEmails[idx])

	}

	// Field: Signature

	if m.Signature != nil {
		*m.Signature = htmlsanitizer.Sanitize(*m.Signature)
	}

	return nil
}
