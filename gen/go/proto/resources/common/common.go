package common

import "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"

// IJobInfo defines methods for accessing and mutating job and grade information.
type IJobInfo interface {
	// GetJob returns the job name or identifier.
	GetJob() string
	// SetJob sets the job name or identifier.
	SetJob(string)
	// SetJobLabel sets the label for the job.
	SetJobLabel(label string)
	// GetJobGrade returns the job grade as an integer.
	GetJobGrade() int32
	// SetJobGrade sets the job grade as an integer.
	SetJobGrade(int32)
	// SetJobGradeLabel sets the label for the job grade.
	SetJobGradeLabel(label string)
}

// IJobName defines methods for accessing and mutating job name and label.
type IJobName interface {
	// GetJob returns the job name or identifier.
	GetJob() string
	// SetJobLabel sets the label for the job.
	SetJobLabel(label string)
}

// ICategory defines methods for accessing and mutating category information.
type ICategory interface {
	// GetCategoryId returns the unique identifier for the category.
	GetCategoryId() uint64
	// SetCategory sets the category using a Category struct.
	SetCategory(*documents.Category)
}
