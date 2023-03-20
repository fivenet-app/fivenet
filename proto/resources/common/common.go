package common

type IJobInfo interface {
	GetJob() string
	GetJobGrade() int32
	SetJobLabel(label string)
	SetJobGradeLabel(label string)
}
