package common

import "github.com/galexrt/arpanet/proto/resources/documents"

type IJobInfo interface {
	GetJob() string
	GetJobGrade() int32
	SetJobLabel(label string)
	SetJobGradeLabel(label string)
}

type IDocumentCategory interface {
	GetCategoryId() uint64
	SetCategory(*documents.DocumentCategory)
}
