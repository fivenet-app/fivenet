// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: resources/users/users.proto

package users

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

// Validate checks the field values on UserShort with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserShort) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserShort with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserShortMultiError, or nil
// if none found.
func (m *UserShort) ValidateAll() error {
	return m.validate(true)
}

func (m *UserShort) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUserId() <= 0 {
		err := UserShortValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetJob()) > 20 {
		err := UserShortValidationError{
			field:  "Job",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetJobGrade() <= -1 {
		err := UserShortValidationError{
			field:  "JobGrade",
			reason: "value must be greater than -1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetFirstname()); l < 1 || l > 50 {
		err := UserShortValidationError{
			field:  "Firstname",
			reason: "value length must be between 1 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetLastname()); l < 1 || l > 50 {
		err := UserShortValidationError{
			field:  "Lastname",
			reason: "value length must be between 1 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetDateofbirth()) != 10 {
		err := UserShortValidationError{
			field:  "Dateofbirth",
			reason: "value length must be 10 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if m.Identifier != nil {

		if utf8.RuneCountInString(m.GetIdentifier()) > 64 {
			err := UserShortValidationError{
				field:  "Identifier",
				reason: "value length must be at most 64 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.JobLabel != nil {

		if utf8.RuneCountInString(m.GetJobLabel()) > 50 {
			err := UserShortValidationError{
				field:  "JobLabel",
				reason: "value length must be at most 50 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.JobGradeLabel != nil {

		if utf8.RuneCountInString(m.GetJobGradeLabel()) > 50 {
			err := UserShortValidationError{
				field:  "JobGradeLabel",
				reason: "value length must be at most 50 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.PhoneNumber != nil {

		if utf8.RuneCountInString(m.GetPhoneNumber()) > 20 {
			err := UserShortValidationError{
				field:  "PhoneNumber",
				reason: "value length must be at most 20 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Avatar != nil {

		if all {
			switch v := interface{}(m.GetAvatar()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserShortValidationError{
						field:  "Avatar",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserShortValidationError{
						field:  "Avatar",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetAvatar()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserShortValidationError{
					field:  "Avatar",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return UserShortMultiError(errors)
	}

	return nil
}

// UserShortMultiError is an error wrapping multiple validation errors returned
// by UserShort.ValidateAll() if the designated constraints aren't met.
type UserShortMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserShortMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserShortMultiError) AllErrors() []error { return m }

// UserShortValidationError is the validation error returned by
// UserShort.Validate if the designated constraints aren't met.
type UserShortValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserShortValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserShortValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserShortValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserShortValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserShortValidationError) ErrorName() string { return "UserShortValidationError" }

// Error satisfies the builtin error interface
func (e UserShortValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserShort.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserShortValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserShortValidationError{}

// Validate checks the field values on User with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *User) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on User with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in UserMultiError, or nil if none found.
func (m *User) ValidateAll() error {
	return m.validate(true)
}

func (m *User) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUserId() <= 0 {
		err := UserValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetJob()) > 20 {
		err := UserValidationError{
			field:  "Job",
			reason: "value length must be at most 20 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetJobGrade() <= -1 {
		err := UserValidationError{
			field:  "JobGrade",
			reason: "value must be greater than -1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetFirstname()); l < 1 || l > 50 {
		err := UserValidationError{
			field:  "Firstname",
			reason: "value length must be between 1 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetLastname()); l < 1 || l > 50 {
		err := UserValidationError{
			field:  "Lastname",
			reason: "value length must be between 1 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetDateofbirth()) != 10 {
		err := UserValidationError{
			field:  "Dateofbirth",
			reason: "value length must be 10 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)

	}

	if all {
		switch v := interface{}(m.GetProps()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UserValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UserValidationError{
					field:  "Props",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetProps()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserValidationError{
				field:  "Props",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetLicenses() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserValidationError{
						field:  fmt.Sprintf("Licenses[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserValidationError{
						field:  fmt.Sprintf("Licenses[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserValidationError{
					field:  fmt.Sprintf("Licenses[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Identifier != nil {

		if utf8.RuneCountInString(m.GetIdentifier()) > 64 {
			err := UserValidationError{
				field:  "Identifier",
				reason: "value length must be at most 64 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.JobLabel != nil {

		if utf8.RuneCountInString(m.GetJobLabel()) > 50 {
			err := UserValidationError{
				field:  "JobLabel",
				reason: "value length must be at most 50 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.JobGradeLabel != nil {

		if utf8.RuneCountInString(m.GetJobGradeLabel()) > 50 {
			err := UserValidationError{
				field:  "JobGradeLabel",
				reason: "value length must be at most 50 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Sex != nil {

		if l := utf8.RuneCountInString(m.GetSex()); l < 1 || l > 2 {
			err := UserValidationError{
				field:  "Sex",
				reason: "value length must be between 1 and 2 runes, inclusive",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Height != nil {
		// no validation rules for Height
	}

	if m.PhoneNumber != nil {

		if utf8.RuneCountInString(m.GetPhoneNumber()) > 20 {
			err := UserValidationError{
				field:  "PhoneNumber",
				reason: "value length must be at most 20 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Visum != nil {

		if m.GetVisum() < 0 {
			err := UserValidationError{
				field:  "Visum",
				reason: "value must be greater than or equal to 0",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Playtime != nil {

		if m.GetPlaytime() < 0 {
			err := UserValidationError{
				field:  "Playtime",
				reason: "value must be greater than or equal to 0",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.Avatar != nil {

		if all {
			switch v := interface{}(m.GetAvatar()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserValidationError{
						field:  "Avatar",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserValidationError{
						field:  "Avatar",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetAvatar()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserValidationError{
					field:  "Avatar",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Group != nil {

		if utf8.RuneCountInString(m.GetGroup()) > 50 {
			err := UserValidationError{
				field:  "Group",
				reason: "value length must be at most 50 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return UserMultiError(errors)
	}

	return nil
}

// UserMultiError is an error wrapping multiple validation errors returned by
// User.ValidateAll() if the designated constraints aren't met.
type UserMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserMultiError) AllErrors() []error { return m }

// UserValidationError is the validation error returned by User.Validate if the
// designated constraints aren't met.
type UserValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserValidationError) ErrorName() string { return "UserValidationError" }

// Error satisfies the builtin error interface
func (e UserValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUser.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserValidationError{}

// Validate checks the field values on License with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *License) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on License with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in LicenseMultiError, or nil if none found.
func (m *License) ValidateAll() error {
	return m.validate(true)
}

func (m *License) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetType()); l < 3 || l > 60 {
		err := LicenseValidationError{
			field:  "Type",
			reason: "value length must be between 3 and 60 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Label

	if len(errors) > 0 {
		return LicenseMultiError(errors)
	}

	return nil
}

// LicenseMultiError is an error wrapping multiple validation errors returned
// by License.ValidateAll() if the designated constraints aren't met.
type LicenseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LicenseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LicenseMultiError) AllErrors() []error { return m }

// LicenseValidationError is the validation error returned by License.Validate
// if the designated constraints aren't met.
type LicenseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LicenseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LicenseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LicenseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LicenseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LicenseValidationError) ErrorName() string { return "LicenseValidationError" }

// Error satisfies the builtin error interface
func (e LicenseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLicense.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LicenseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LicenseValidationError{}

// Validate checks the field values on UserLicenses with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserLicenses) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserLicenses with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserLicensesMultiError, or
// nil if none found.
func (m *UserLicenses) ValidateAll() error {
	return m.validate(true)
}

func (m *UserLicenses) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUserId() <= 0 {
		err := UserLicensesValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetLicenses() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserLicensesValidationError{
						field:  fmt.Sprintf("Licenses[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserLicensesValidationError{
						field:  fmt.Sprintf("Licenses[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserLicensesValidationError{
					field:  fmt.Sprintf("Licenses[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return UserLicensesMultiError(errors)
	}

	return nil
}

// UserLicensesMultiError is an error wrapping multiple validation errors
// returned by UserLicenses.ValidateAll() if the designated constraints aren't met.
type UserLicensesMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserLicensesMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserLicensesMultiError) AllErrors() []error { return m }

// UserLicensesValidationError is the validation error returned by
// UserLicenses.Validate if the designated constraints aren't met.
type UserLicensesValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserLicensesValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserLicensesValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserLicensesValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserLicensesValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserLicensesValidationError) ErrorName() string { return "UserLicensesValidationError" }

// Error satisfies the builtin error interface
func (e UserLicensesValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserLicenses.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserLicensesValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserLicensesValidationError{}
