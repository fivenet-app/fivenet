package documents

func (x *DocumentShort) SetCategory(cat *DocumentCategory) {
	x.Category = cat
}

func (x *Document) SetCategory(cat *DocumentCategory) {
	x.Category = cat
}

func (x *DocumentJobAccess) GetJobGrade() int32 {
	return x.MinimumGrade
}

func (x *DocumentJobAccess) SetJobLabel(label string) {
	x.JobLabel = label
}

func (x *DocumentJobAccess) SetJobGradeLabel(label string) {
	x.JobGradeLabel = label
}
