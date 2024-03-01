// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/centrum/dispatches.proto

package centrum

import (
	"github.com/galexrt/fivenet/pkg/htmlsanitizer"
)

func (m *Dispatch) Sanitize() error {

	if m.Description != nil {
		*m.Description = htmlsanitizer.Sanitize(*m.Description)
	}

	m.Message = htmlsanitizer.Sanitize(m.Message)

	if m.Postal != nil {
		*m.Postal = htmlsanitizer.Sanitize(*m.Postal)
	}

	return nil
}

func (m *DispatchStatus) Sanitize() error {

	if m.Code != nil {
		*m.Code = htmlsanitizer.Sanitize(*m.Code)
	}

	if m.Postal != nil {
		*m.Postal = htmlsanitizer.Sanitize(*m.Postal)
	}

	if m.Reason != nil {
		*m.Reason = htmlsanitizer.Sanitize(*m.Reason)
	}

	return nil
}