package documents

import (
	"testing"

	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	"github.com/stretchr/testify/assert"
)

func TestDocumentRequestNotificationAction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		requestType documentsactivity.DocActivityType
		expected    string
	}{
		"access": {
			requestType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_ACCESS,
			expected:    "access",
		},
		"closure": {
			requestType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_CLOSURE,
			expected:    "closure",
		},
		"opening": {
			requestType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OPENING,
			expected:    "opening",
		},
		"update": {
			requestType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_UPDATE,
			expected:    "update",
		},
		"owner change": {
			requestType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE,
			expected:    "owner_change",
		},
		"deletion": {
			requestType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_DELETION,
			expected:    "deletion",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, documentRequestNotificationAction(test.requestType))
		})
	}

	assert.Equal(
		t,
		"created",
		documentRequestNotificationAction(documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_UNSPECIFIED),
	)
}
