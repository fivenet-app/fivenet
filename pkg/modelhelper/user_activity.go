package modelhelper

type AcitvityType string

const (
	ChangedActivityType   AcitvityType = "changed"
	CreatedActivityType   AcitvityType = "created"
	MentionedActivityType AcitvityType = "mentioned"
)
