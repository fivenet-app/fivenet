// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/qualifications/qualifications.proto

package qualifications

import (
	"github.com/fivenet-app/fivenet/pkg/html/htmlsanitizer"
)

func (m *Qualification) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Abbreviation
	m.Abbreviation = htmlsanitizer.StripTags(m.Abbreviation)

	// Field: Access
	if m.Access != nil {
		if v, ok := any(m.GetAccess()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Content
	if m.Content != nil {
		if v, ok := any(m.GetContent()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := any(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Creator
	if m.Creator != nil {
		if v, ok := any(m.GetCreator()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: DeletedAt
	if m.DeletedAt != nil {
		if v, ok := any(m.GetDeletedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Description

	if m.Description != nil {
		*m.Description = htmlsanitizer.StripTags(*m.Description)
	}

	// Field: DiscordSettings
	if m.DiscordSettings != nil {
		if v, ok := any(m.GetDiscordSettings()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Exam
	if m.Exam != nil {
		if v, ok := any(m.GetExam()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: ExamSettings
	if m.ExamSettings != nil {
		if v, ok := any(m.GetExamSettings()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: LabelSyncFormat

	if m.LabelSyncFormat != nil {
		*m.LabelSyncFormat = htmlsanitizer.StripTags(*m.LabelSyncFormat)
	}

	// Field: Request
	if m.Request != nil {
		if v, ok := any(m.GetRequest()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Requirements
	for idx, item := range m.Requirements {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Result
	if m.Result != nil {
		if v, ok := any(m.GetResult()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Title
	m.Title = htmlsanitizer.Sanitize(m.Title)

	// Field: UpdatedAt
	if m.UpdatedAt != nil {
		if v, ok := any(m.GetUpdatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *QualificationDiscordSettings) Sanitize() error {
	if m == nil {
		return nil
	}

	return nil
}

func (m *QualificationExamSettings) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Time
	if m.Time != nil {
		if v, ok := any(m.GetTime()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *QualificationRequest) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: ApprovedAt
	if m.ApprovedAt != nil {
		if v, ok := any(m.GetApprovedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Approver
	if m.Approver != nil {
		if v, ok := any(m.GetApprover()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: ApproverComment

	if m.ApproverComment != nil {
		*m.ApproverComment = htmlsanitizer.StripTags(*m.ApproverComment)
	}

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := any(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: DeletedAt
	if m.DeletedAt != nil {
		if v, ok := any(m.GetDeletedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Qualification
	if m.Qualification != nil {
		if v, ok := any(m.GetQualification()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: User
	if m.User != nil {
		if v, ok := any(m.GetUser()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: UserComment

	if m.UserComment != nil {
		*m.UserComment = htmlsanitizer.StripTags(*m.UserComment)
	}

	return nil
}

func (m *QualificationRequirement) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := any(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: TargetQualification
	if m.TargetQualification != nil {
		if v, ok := any(m.GetTargetQualification()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *QualificationResult) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := any(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Creator
	if m.Creator != nil {
		if v, ok := any(m.GetCreator()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: DeletedAt
	if m.DeletedAt != nil {
		if v, ok := any(m.GetDeletedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Qualification
	if m.Qualification != nil {
		if v, ok := any(m.GetQualification()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Summary
	m.Summary = htmlsanitizer.StripTags(m.Summary)

	// Field: User
	if m.User != nil {
		if v, ok := any(m.GetUser()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *QualificationShort) Sanitize() error {
	if m == nil {
		return nil
	}

	// Field: Abbreviation
	m.Abbreviation = htmlsanitizer.StripTags(m.Abbreviation)

	// Field: CreatedAt
	if m.CreatedAt != nil {
		if v, ok := any(m.GetCreatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Creator
	if m.Creator != nil {
		if v, ok := any(m.GetCreator()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: DeletedAt
	if m.DeletedAt != nil {
		if v, ok := any(m.GetDeletedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Description

	if m.Description != nil {
		*m.Description = htmlsanitizer.StripTags(*m.Description)
	}

	// Field: ExamSettings
	if m.ExamSettings != nil {
		if v, ok := any(m.GetExamSettings()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Requirements
	for idx, item := range m.Requirements {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}

	}

	// Field: Result
	if m.Result != nil {
		if v, ok := any(m.GetResult()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	// Field: Title
	m.Title = htmlsanitizer.Sanitize(m.Title)

	// Field: UpdatedAt
	if m.UpdatedAt != nil {
		if v, ok := any(m.GetUpdatedAt()).(interface{ Sanitize() error }); ok {
			if err := v.Sanitize(); err != nil {
				return err
			}
		}
	}

	return nil
}
