package documents

import (
	"testing"

	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
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
					Email:     new("canonical@example.test"),
					OpenFines: new(int64(42)),
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

func TestRenderTemplateDefaultsMissingUserProps(t *testing.T) {
	t.Parallel()

	server := &Server{}
	tmpl := &documentstemplates.Template{
		Content: "{{ (index .Users 0).Props.BloodType }}",
	}
	data := &resolvedTemplateData{Users: []*users.User{{UserId: 42}}}

	_, _, content, err := server.renderTemplate(tmpl, data)

	require.NoError(t, err)
	require.Empty(t, content)
	require.NotNil(t, data.Users[0].GetProps())
	require.NotNil(t, data.Users[0].GetProps().BloodType)
}

func TestRenderTemplateFindsDynamicProperty(t *testing.T) {
	t.Parallel()

	server := &Server{}
	value := &properties.PropertyValue{}
	value.SetStringValue("blue")
	tmpl := &documentstemplates.Template{
		Content: `<span data-template-var="property (index .Users 0).Properties &#34;user.favorite_color&#34;">{{ (property (index .Users 0).Properties "user.favorite_color").GetValue.GetStringValue }}</span>`,
	}
	data := &resolvedTemplateData{Users: []*users.User{{
		Properties: &properties.Properties{Properties: []*properties.Property{{
			Key:   "user.favorite_color",
			Value: value,
		}}},
	}}}

	_, _, content, err := server.renderTemplate(tmpl, data)

	require.NoError(t, err)
	require.Equal(t, "blue", content)
}

func TestRenderTemplateFindsDynamicPropertyOnUserProps(t *testing.T) {
	t.Parallel()

	server := &Server{}
	value := &properties.PropertyValue{}
	value.SetStringValue("blue")
	tmpl := &documentstemplates.Template{
		Content: `<span data-template-var="property (index .Users 0).Props.Properties &#34;user.favorite_color&#34;">{{ (property (index .Users 0).Props.Properties "user.favorite_color").GetValue.GetStringValue }}</span>`,
	}
	data := &resolvedTemplateData{Users: []*users.User{{
		Properties: &properties.Properties{Properties: []*properties.Property{{
			Key:   "user.favorite_color",
			Value: value,
		}}},
	}}}

	_, _, content, err := server.renderTemplate(tmpl, data)

	require.NoError(t, err)
	require.Equal(t, "blue", content)
}

func TestRenderTemplateSupportsTemplateVariablePresets(t *testing.T) {
	t.Parallel()

	server := &Server{}
	propertyValue := &properties.PropertyValue{}
	propertyValue.SetBoolValue(true)
	tmpl := &documentstemplates.Template{
		ContentTitle: `{{ now | date "02.01.2006 15:04" }}|{{ now | date "02.01.2006" }}|{{ now | date "15:04" }}|{{ .ActiveChar.GetFirstname }}|{{ .ActiveChar.GetLastname }}|{{ .ActiveChar.GetDateofbirth }}|{{ .ActiveChar.GetPhoneNumber }}|{{ .ActiveChar.GetProps.GetNamePrefix }}|{{ .ActiveChar.GetProps.GetNameSuffix }}|{{ .ActiveChar.GetEmail }}|{{ (first .Users).GetFirstname }}|{{ (first .Users).GetLastname }}|{{ (first .Users).GetDateofbirth }}|{{ (first .Users).GetSex }}|{{ (first .Users).GetHeight }}|{{ (first .Users).GetProps.GetWanted }}|{{ (first .Users).GetPhoneNumber }}|{{ (property (first .Users).GetProperties "user.wanted").GetValue.GetBoolValue }}|{{ (propertyValue (first .Users).GetProperties "user.wanted").GetBoolValue }}`,
	}
	data := &resolvedTemplateData{
		ActiveChar: &jobscolleagues.Colleague{
			Firstname:   "Active",
			Lastname:    "Colleague",
			Dateofbirth: "2000-01-02",
			PhoneNumber: new("555-0100"),
			Email:       new("active@example.test"),
			Props: &jobscolleagues.ColleagueProps{
				NamePrefix: new("Dr."),
				NameSuffix: new("Jr."),
			},
		},
		Users: []*users.User{{
			Firstname:   "Citizen",
			Lastname:    "Example",
			Dateofbirth: "2001-02-03",
			Sex:         new("female"),
			Height:      new(float32(1.75)),
			PhoneNumber: new("555-0101"),
			Props:       &usersprops.UserProps{Wanted: new(true)},
			Properties: &properties.Properties{Properties: []*properties.Property{{
				Key:   "user.wanted",
				Value: propertyValue,
			}}},
		}},
	}

	title, _, _, err := server.renderTemplate(tmpl, data)

	require.NoError(t, err)
	require.Contains(
		t,
		title,
		"Active|Colleague|2000-01-02|555-0100|Dr.|Jr.|active@example.test|Citizen|Example|2001-02-03|female|1.75|true|555-0101|true|true",
	)
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
		PhoneNumber:          new("555-0100"),
		ProfilePicture:       new("/avatars/42.png"),
		ProfilePictureFileId: new(int64(42)),
		Props: &jobscolleagues.ColleagueProps{
			NamePrefix: new("Dr."),
			NameSuffix: new("Jr."),
		},
	}}

	title, state, content, err := server.renderTemplate(tmpl, data)

	require.NoError(t, err)
	require.Equal(t, "Dr. Active Colleague Jr.", title)
	require.Equal(t, "555-0100", state)
	require.Equal(t, "/avatars/42.png (42)", content)
}

func TestRenderTemplateUsesNilSafeProtoGetters(t *testing.T) {
	t.Parallel()

	server := &Server{}
	tmpl := &documentstemplates.Template{
		Content: "{{ (index .Vehicles 0).GetPlate }} {{ (index .Vehicles 0).GetOwner.GetLastname }}",
	}
	data := &resolvedTemplateData{
		Vehicles: []*resourcesvehicles.Vehicle{{Plate: "ABC 123"}},
	}

	_, _, content, err := server.renderTemplate(tmpl, data)

	require.NoError(t, err)
	require.Equal(t, "ABC 123 ", content)
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

func TestStripTemplateVarSpansPreservesQuotedArguments(t *testing.T) {
	t.Parallel()

	content := `<p><span class="template-var" data-template-var="now | date &#34;02.01.2006 15:04&#34;">{{ now | date "02.01.2006 15:04" }}</span></p>`

	stripped, err := stripTemplateActionSpans(content)
	require.NoError(t, err)
	require.Equal(t, `<p>{{ now | date "02.01.2006 15:04" }}</p>`, stripped)
}

func TestStripTemplateActionSpansKeepsAttributeEntitiesEscaped(t *testing.T) {
	t.Parallel()

	content := `<p><span data-template-var="now | date &#34;02.01.2006 15:04&#34;">{{ now | date "02.01.2006 15:04" }}</span><span data-keep="now | date &#34;02.01.2006 15:04&#34;">keep</span></p>`

	stripped, err := stripTemplateActionSpans(content)
	require.NoError(t, err)
	require.Equal(
		t,
		`<p>{{ now | date "02.01.2006 15:04" }}<span data-keep="now | date &#34;02.01.2006 15:04&#34;">keep</span></p>`,
		stripped,
	)
}

func TestStripTemplateActionSpansPreservesOtherHTMLAttributes(t *testing.T) {
	t.Parallel()

	content := `<p data-custom="keep"><span data-template-block="if .Active" data-left-trim="true">{{- if .Active -}}</span><span data-template-block-end="end" data-right-trim="true">{{ end }}</span><span data-keep="yes">content</span></p>`

	stripped, err := stripTemplateActionSpans(content)
	require.NoError(t, err)
	require.Equal(
		t,
		`<p data-custom="keep">{{- if .Active -}}{{ end }}<span data-keep="yes">content</span></p>`,
		stripped,
	)
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

	tmpl := &documentstemplates.Template{
		Schema: &documentstemplates.TemplateSchema{
			Requirements: &documentstemplates.TemplateRequirements{
				Users: &documentstemplates.ObjectSpecs{Required: true, Min: new(int32(2))},
			},
		},
	}

	err := validateTemplateRequirements(tmpl, &resolvedTemplateData{Users: []*users.User{{}}})

	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
