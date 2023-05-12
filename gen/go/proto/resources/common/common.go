package common

import "github.com/galexrt/fivenet/gen/go/proto/resources/documents"

const (
	SuperuserCategoryPerm   = "Superuser"
	SuperuserAnyAccessName  = "AnyAccess"
	SuperuserAnyAccess      = SuperuserCategoryPerm + SuperuserAnyAccessName
	SuperuserAnyAccessGuard = "superuser-anyaccess"
)

type IJobInfo interface {
	GetJob() string
	GetJobGrade() int32
	SetJobLabel(label string)
	SetJobGradeLabel(label string)
}

type IJobName interface {
	GetJob() string
	SetJobLabel(label string)
}

type IDocumentCategory interface {
	GetCategoryId() uint64
	SetCategory(*documents.DocumentCategory)
}
