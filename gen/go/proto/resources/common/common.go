package common

import "github.com/galexrt/fivenet/gen/go/proto/resources/documents"

const (
	SuperuserAnyAccess      = "Superuser.AnyAccess"
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
