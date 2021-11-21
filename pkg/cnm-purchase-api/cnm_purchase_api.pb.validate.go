// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: ozonmp/cnm_purchase_api/v1/cnm_purchase_api.proto

package cnm_purchase_api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
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
)

// Validate checks the field values on Purchase with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Purchase) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for TotalSum

	return nil
}

// PurchaseValidationError is the validation error returned by
// Purchase.Validate if the designated constraints aren't met.
type PurchaseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PurchaseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PurchaseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PurchaseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PurchaseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PurchaseValidationError) ErrorName() string { return "PurchaseValidationError" }

// Error satisfies the builtin error interface
func (e PurchaseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPurchase.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PurchaseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PurchaseValidationError{}

// Validate checks the field values on DescribePurchaseV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *DescribePurchaseV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetPurchaseId() <= 0 {
		return DescribePurchaseV1RequestValidationError{
			field:  "PurchaseId",
			reason: "value must be greater than 0",
		}
	}

	return nil
}

// DescribePurchaseV1RequestValidationError is the validation error returned by
// DescribePurchaseV1Request.Validate if the designated constraints aren't met.
type DescribePurchaseV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribePurchaseV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribePurchaseV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribePurchaseV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribePurchaseV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribePurchaseV1RequestValidationError) ErrorName() string {
	return "DescribePurchaseV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e DescribePurchaseV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribePurchaseV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribePurchaseV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribePurchaseV1RequestValidationError{}

// Validate checks the field values on DescribePurchaseV1Response with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *DescribePurchaseV1Response) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetValue()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DescribePurchaseV1ResponseValidationError{
				field:  "Value",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// DescribePurchaseV1ResponseValidationError is the validation error returned
// by DescribePurchaseV1Response.Validate if the designated constraints aren't met.
type DescribePurchaseV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribePurchaseV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribePurchaseV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribePurchaseV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribePurchaseV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribePurchaseV1ResponseValidationError) ErrorName() string {
	return "DescribePurchaseV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DescribePurchaseV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribePurchaseV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribePurchaseV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribePurchaseV1ResponseValidationError{}

// Validate checks the field values on CreatePurchaseV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreatePurchaseV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetTotalSum() < 0 {
		return CreatePurchaseV1RequestValidationError{
			field:  "TotalSum",
			reason: "value must be greater than or equal to 0",
		}
	}

	return nil
}

// CreatePurchaseV1RequestValidationError is the validation error returned by
// CreatePurchaseV1Request.Validate if the designated constraints aren't met.
type CreatePurchaseV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePurchaseV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePurchaseV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePurchaseV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePurchaseV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePurchaseV1RequestValidationError) ErrorName() string {
	return "CreatePurchaseV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePurchaseV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePurchaseV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePurchaseV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePurchaseV1RequestValidationError{}

// Validate checks the field values on CreatePurchaseV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreatePurchaseV1Response) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for PurchaseId

	return nil
}

// CreatePurchaseV1ResponseValidationError is the validation error returned by
// CreatePurchaseV1Response.Validate if the designated constraints aren't met.
type CreatePurchaseV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePurchaseV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePurchaseV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePurchaseV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePurchaseV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePurchaseV1ResponseValidationError) ErrorName() string {
	return "CreatePurchaseV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePurchaseV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePurchaseV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePurchaseV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePurchaseV1ResponseValidationError{}

// Validate checks the field values on ListPurchasesV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListPurchasesV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetLimit() <= 0 {
		return ListPurchasesV1RequestValidationError{
			field:  "Limit",
			reason: "value must be greater than 0",
		}
	}

	if m.GetCursor() < 0 {
		return ListPurchasesV1RequestValidationError{
			field:  "Cursor",
			reason: "value must be greater than or equal to 0",
		}
	}

	return nil
}

// ListPurchasesV1RequestValidationError is the validation error returned by
// ListPurchasesV1Request.Validate if the designated constraints aren't met.
type ListPurchasesV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListPurchasesV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListPurchasesV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListPurchasesV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListPurchasesV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListPurchasesV1RequestValidationError) ErrorName() string {
	return "ListPurchasesV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListPurchasesV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListPurchasesV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListPurchasesV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListPurchasesV1RequestValidationError{}

// Validate checks the field values on ListPurchasesV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListPurchasesV1Response) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetItems() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListPurchasesV1ResponseValidationError{
					field:  fmt.Sprintf("Items[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Cursor

	return nil
}

// ListPurchasesV1ResponseValidationError is the validation error returned by
// ListPurchasesV1Response.Validate if the designated constraints aren't met.
type ListPurchasesV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListPurchasesV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListPurchasesV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListPurchasesV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListPurchasesV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListPurchasesV1ResponseValidationError) ErrorName() string {
	return "ListPurchasesV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListPurchasesV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListPurchasesV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListPurchasesV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListPurchasesV1ResponseValidationError{}

// Validate checks the field values on RemovePurchaseV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RemovePurchaseV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetPurchaseId() <= 0 {
		return RemovePurchaseV1RequestValidationError{
			field:  "PurchaseId",
			reason: "value must be greater than 0",
		}
	}

	return nil
}

// RemovePurchaseV1RequestValidationError is the validation error returned by
// RemovePurchaseV1Request.Validate if the designated constraints aren't met.
type RemovePurchaseV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemovePurchaseV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemovePurchaseV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemovePurchaseV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemovePurchaseV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemovePurchaseV1RequestValidationError) ErrorName() string {
	return "RemovePurchaseV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e RemovePurchaseV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemovePurchaseV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemovePurchaseV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemovePurchaseV1RequestValidationError{}

// Validate checks the field values on RemovePurchaseV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RemovePurchaseV1Response) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Found

	return nil
}

// RemovePurchaseV1ResponseValidationError is the validation error returned by
// RemovePurchaseV1Response.Validate if the designated constraints aren't met.
type RemovePurchaseV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemovePurchaseV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemovePurchaseV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemovePurchaseV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemovePurchaseV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemovePurchaseV1ResponseValidationError) ErrorName() string {
	return "RemovePurchaseV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e RemovePurchaseV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemovePurchaseV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemovePurchaseV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemovePurchaseV1ResponseValidationError{}