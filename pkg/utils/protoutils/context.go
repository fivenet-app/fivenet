package protoutils

import (
	"context"
	"errors"
)

// IsContextCanceled reports whether err is context.Canceled or wraps it (might check proto response status in the future).
func IsContextCanceled(err error) bool {
	return errors.Is(err, context.Canceled)
}
