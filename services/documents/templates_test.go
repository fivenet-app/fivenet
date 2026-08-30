package documents

import (
	"testing"

	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRenderTemplateUsesResolvedUserProps(t *testing.T) {
	t.Parallel()

	server := &Server{}
	tmpl := &documentstemplates.Template{
		ContentTitle: "{{ (index .Users 0).Props.Email }}",
		State:        "{{ (index .Users 0).Props.OpenFines }}",
		Content:      "{{ (index .Users 0).Firstname }}",
	}
	data := &resolvedTemplateData{
		Users: []*users.User{
			{
				Firstname: "Canonical",
				Props: &usersprops.UserProps{
					Email:     stringPtr("canonical@example.test"),
					OpenFines: int64Ptr(42),
				},
			},
		},
	}

	title, state, content, err := server.renderTemplate(tmpl, data)

	require.NoError(t, err)
	require.Equal(t, "canonical@example.test", title)
	require.Equal(t, "42", state)
	require.Equal(t, "Canonical", content)
}

func TestRenderTemplateUsesColleagueActiveChar(t *testing.T) {
	t.Parallel()

	server := &Server{}
	tmpl := &documentstemplates.Template{
		ContentTitle: "{{ .ActiveChar.Props.NamePrefix }} {{ .ActiveChar.Firstname }} {{ .ActiveChar.Lastname }} {{ .ActiveChar.Props.NameSuffix }}",
		State:        "{{ .ActiveChar.PhoneNumber }}",
		Content:      "{{ .ActiveChar.ProfilePicture }} ({{ .ActiveChar.ProfilePictureFileId }})",
	}
	data := &resolvedTemplateData{ActiveChar: &jobscolleagues.Colleague{
		Firstname:            "Active",
		Lastname:             "Colleague",
		PhoneNumber:          stringPtr("555-0100"),
		ProfilePicture:       stringPtr("/avatars/42.png"),
		ProfilePictureFileId: int64Ptr(42),
		Props: &jobscolleagues.ColleagueProps{
			NamePrefix: stringPtr("Dr."),
			NameSuffix: stringPtr("Jr."),
		},
	}}

	title, state, content, err := server.renderTemplate(tmpl, data)

	require.NoError(t, err)
	require.Equal(t, "Dr. Active Colleague Jr.", title)
	require.Equal(t, "555-0100", state)
	require.Equal(t, "/avatars/42.png (42)", content)
}

func TestRenderTemplateStripsTemplateVarSpan(t *testing.T) {
	t.Parallel()

	server := &Server{}
	tmpl := &documentstemplates.Template{
		Content: `<p>Hello <span data-template-var="(index .Users 0).Firstname" class="template-var">{{ (index .Users 0).Firstname }}</span></p>`,
	}
	data := &resolvedTemplateData{
		Users: []*users.User{{Firstname: "Smith"}},
	}

	_, _, content, err := server.renderTemplate(tmpl, data)

	require.NoError(t, err)
	require.Equal(t, "<p>Hello Smith</p>", content)
}

func TestStripTemplateVarSpansPreservesTrimMarkers(t *testing.T) {
	t.Parallel()

	content := `<p><span class="template-var" data-template-var=".Firstname">{{- .Firstname -}}</span></p>`

	stripped, err := stripTemplateActionSpans(content)
	require.NoError(t, err)
	require.Equal(t, `<p>{{- .Firstname -}}</p>`, stripped)
}

func TestStripTemplateActionSpansPreservesOtherHTMLAttributes(t *testing.T) {
	t.Parallel()

	content := `<p data-custom="keep"><span data-template-block="if .Active" data-left-trim="true">{{- if .Active -}}</span><span data-template-block-end="end" data-right-trim="true">{{ end }}</span><span data-keep="yes">content</span></p>`

	stripped, err := stripTemplateActionSpans(content)
	require.NoError(t, err)
	require.Equal(t, `<p data-custom="keep">{{- if .Active -}}{{ end }}<span data-keep="yes">content</span></p>`, stripped)
}

func TestStripTemplateActionSpansUnwrapsNestedActions(t *testing.T) {
	t.Parallel()

	content := `<p><span data-template-var=".Outer">before <span data-template-block="if .Active">{{ if .Active }}</span> after</span></p>`

	stripped, err := stripTemplateActionSpans(content)
	require.NoError(t, err)
	require.Equal(t, `<p>before {{ if .Active }} after</p>`, stripped)
}

func TestValidateTemplateRequirementsAfterResolution(t *testing.T) {
	t.Parallel()

	required := true
	tmpl := &documentstemplates.Template{
		Schema: &documentstemplates.TemplateSchema{
			Requirements: &documentstemplates.TemplateRequirements{
				Users: &documentstemplates.ObjectSpecs{Required: &required, Min: int32Ptr(2)},
			},
		},
	}

	err := validateTemplateRequirements(tmpl, &resolvedTemplateData{Users: []*users.User{{}}})

	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func stringPtr(v string) *string { return &v }

func int32Ptr(v int32) *int32 { return &v }

func int64Ptr(v int64) *int64 { return &v }
