// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/jobs/conduct.proto

package jobs

import (
	"github.com/fivenet-app/fivenet/pkg/htmlsanitizer"
)

func (m *ConductEntry) Sanitize() error {

	m.Message = htmlsanitizer.Sanitize(m.Message)

	return nil
}
