// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/internet/domain.proto

package internet

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on TLD with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *TLD) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TLD with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in TLDMultiError, or nil if none found.
func (m *TLD) ValidateAll() error {
	return m.validate(true)
}

func (m *TLD) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TLDValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TLDValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TLDValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if l := utf8.RuneCountInString(m.GetName()); l < 2 || l > 24 {
		err := TLDValidationError{
			field:  "Name",
			reason: "value length must be between 2 and 24 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Internal

	if m.UpdatedAt != nil {

		if all {
			switch v := interface{}(m.GetUpdatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, TLDValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, TLDValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TLDValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.DeletedAt != nil {

		if all {
			switch v := interface{}(m.GetDeletedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, TLDValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, TLDValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetDeletedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TLDValidationError{
					field:  "DeletedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.CreatorId != nil {
		// no validation rules for CreatorId
	}

	if len(errors) > 0 {
		return TLDMultiError(errors)
	}

	return nil
}

// TLDMultiError is an error wrapping multiple validation errors returned by
// TLD.ValidateAll() if the designated constraints aren't met.
type TLDMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TLDMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TLDMultiError) AllErrors() []error { return m }

// TLDValidationError is the validation error returned by TLD.Validate if the
// designated constraints aren't met.
type TLDValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TLDValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TLDValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TLDValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TLDValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TLDValidationError) ErrorName() string { return "TLDValidationError" }

// Error satisfies the builtin error interface
func (e TLDValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTLD.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TLDValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TLDValidationError{}

// Validate checks the field values on Domain with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Domain) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Domain with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in DomainMultiError, or nil if none found.
func (m *Domain) ValidateAll() error {
	return m.validate(true)
}

func (m *Domain) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DomainValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DomainValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DomainValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for TldId

	// no validation rules for Active

	if utf8.RuneCountInString(m.GetName()) > 128 {
		err := DomainValidationError{
			field:  "Name",
			reason: "value length must be at most 128 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.UpdatedAt != nil {

		if all {
			switch v := interface{}(m.GetUpdatedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DomainValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DomainValidationError{
						field:  "UpdatedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DomainValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.DeletedAt != nil {

		if all {
			switch v := interface{}(m.GetDeletedAt()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DomainValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DomainValidationError{
						field:  "DeletedAt",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetDeletedAt()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DomainValidationError{
					field:  "DeletedAt",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Tld != nil {

		if all {
			switch v := interface{}(m.GetTld()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, DomainValidationError{
						field:  "Tld",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, DomainValidationError{
						field:  "Tld",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetTld()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return DomainValidationError{
					field:  "Tld",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.TransferCode != nil {

		if utf8.RuneCountInString(m.GetTransferCode()) != 10 {
			err := DomainValidationError{
				field:  "TransferCode",
				reason: "value length must be 10 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)

		}

		if !_Domain_TransferCode_Pattern.MatchString(m.GetTransferCode()) {
			err := DomainValidationError{
				field:  "TransferCode",
				reason: "value does not match regex pattern \"^[0-9A-Z]{6}$\"",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.ApproverJob != nil {
		// no validation rules for ApproverJob
	}

	if m.ApproverId != nil {
		// no validation rules for ApproverId
	}

	if m.CreatorJob != nil {
		// no validation rules for CreatorJob
	}

	if m.CreatorId != nil {
		// no validation rules for CreatorId
	}

	if len(errors) > 0 {
		return DomainMultiError(errors)
	}

	return nil
}

// DomainMultiError is an error wrapping multiple validation errors returned by
// Domain.ValidateAll() if the designated constraints aren't met.
type DomainMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DomainMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DomainMultiError) AllErrors() []error { return m }

// DomainValidationError is the validation error returned by Domain.Validate if
// the designated constraints aren't met.
type DomainValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DomainValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DomainValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DomainValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DomainValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DomainValidationError) ErrorName() string { return "DomainValidationError" }

// Error satisfies the builtin error interface
func (e DomainValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDomain.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DomainValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DomainValidationError{}

var _Domain_TransferCode_Pattern = regexp.MustCompile("^[0-9A-Z]{6}$")
