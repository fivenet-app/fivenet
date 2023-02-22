package model

type Job struct {
	ID uint `gorm:"primaryKey"`

	Name  string `gorm:"type:varchar(50);uniqueIndex:job_name"`
	Label string `gorm:"type:varchar(50);uniqueIndex:job_name"`

	Grades []JobGrade
}

type JobGrade struct {
	ID uint `gorm:"primaryKey"`

	JobName uint `gorm:"column:job_name;uniqueIndex:idx_job_grade"`
	Grade   int  `gorm:"uniqueIndex:idx_job_grade"`
	Label   string
}

func (JobGrade) TableName() string {
	return "job_grades"
}
