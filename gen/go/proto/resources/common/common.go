package common

import "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"

type IJobInfo interface {
	GetJob() string
	SetJob(string)
	SetJobLabel(label string)
	GetJobGrade() int32
	SetJobGrade(int32)
	SetJobGradeLabel(label string)
}

type IJobName interface {
	GetJob() string
	SetJobLabel(label string)
}

type ICategory interface {
	GetCategoryId() uint64
	SetCategory(*documents.Category)
}
