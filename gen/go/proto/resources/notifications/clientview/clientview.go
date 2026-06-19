package notificationsclientview

type VisibilityMode int

const (
	VisibilityUnsupported VisibilityMode = iota
	VisibilityTargetAccess
	VisibilityJobScoped
)

type TypeSpec struct {
	NatsKey           string
	AccessRegistryKey string
	Visibility        VisibilityMode
}

var objectTypeSpecs = map[ObjectType]TypeSpec{
	ObjectType_OBJECT_TYPE_CITIZEN: {
		NatsKey:           "citizen",
		AccessRegistryKey: "citizen",
		Visibility:        VisibilityTargetAccess,
	},
	ObjectType_OBJECT_TYPE_DOCUMENT: {
		NatsKey:           "document",
		AccessRegistryKey: "documents",
		Visibility:        VisibilityTargetAccess,
	},
	ObjectType_OBJECT_TYPE_WIKI_PAGE: {
		NatsKey:           "wiki_page",
		AccessRegistryKey: "wiki_page",
		Visibility:        VisibilityTargetAccess,
	},
	ObjectType_OBJECT_TYPE_JOBS_COLLEAGUE: {
		NatsKey:    "jobs_colleague",
		Visibility: VisibilityJobScoped,
	},
	ObjectType_OBJECT_TYPE_JOBS_CONDUCT: {
		NatsKey:    "jobs_conduct",
		Visibility: VisibilityJobScoped,
	},
}

func (x ObjectType) Spec() (TypeSpec, bool) {
	spec, ok := objectTypeSpecs[x]
	return spec, ok
}
