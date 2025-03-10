package documents

func (x *Template) GetJob() string {
	return x.CreatorJob
}

func (x *Template) SetJobLabel(label string) {
	x.CreatorJobLabel = &label
}

func (x *TemplateShort) GetJob() string {
	return x.CreatorJob
}

func (x *TemplateShort) SetJobLabel(label string) {
	x.CreatorJobLabel = &label
}

// pkg/access compatibility

func (x *TemplateJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *TemplateJobAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *TemplateUserAccess) GetAccess() AccessLevel {
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

func (x *TemplateUserAccess) GetId() uint64 {
	return 0
}

func (x *TemplateUserAccess) GetTargetId() uint64 {
	return 0
}

func (x *TemplateUserAccess) SetAccess(access AccessLevel) {}

func (x *TemplateUserAccess) GetUserId() int32 {
	return 0
}

func (x *TemplateUserAccess) SetUserId(userId int32) {}
