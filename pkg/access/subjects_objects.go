package access

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
)

func NewDocumentsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetDocuments,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetDocuments.ID,
			CreatedAt:  table.FivenetDocuments.CreatedAt,
			UpdatedAt:  table.FivenetDocuments.UpdatedAt,
			DeletedAt:  table.FivenetDocuments.DeletedAt,
			Public:     table.FivenetDocuments.Public,
			CreatorID:  table.FivenetDocuments.CreatorID,
			CreatorJob: table.FivenetDocuments.CreatorJob,
		},
		AccessTable: table.FivenetDocumentsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetDocumentsAccess.ID,
				TargetID: table.FivenetDocumentsAccess.TargetID,
				Access:   table.FivenetDocumentsAccess.Access,
			},
			SubjectID: table.FivenetDocumentsAccess.SubjectID,
			Effect:    table.FivenetDocumentsAccess.Effect,
		},
		CalculatedVisibilityPublicTable:       table.FivenetDocumentsVisibilityPublic,
		CalculatedVisibilityCreatorTable:      table.FivenetDocumentsVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetDocumentsVisibilitySubject,
		CalculatedVisibilityPublicTargetID:    table.FivenetDocumentsVisibilityPublic.TargetID,
		CalculatedVisibilityCreatorTargetID:   table.FivenetDocumentsVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID:  table.FivenetDocumentsVisibilityCreator.CreatorID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetDocumentsVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetDocumentsVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetDocumentsVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetDocumentsVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetDocumentsVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRulePublic},
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewDocumentTemplatesSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetDocumentsTemplates,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetDocumentsTemplates.ID,
			CreatedAt:  table.FivenetDocumentsTemplates.CreatedAt,
			UpdatedAt:  table.FivenetDocumentsTemplates.UpdatedAt,
			DeletedAt:  table.FivenetDocumentsTemplates.DeletedAt,
			CreatorJob: table.FivenetDocumentsTemplates.CreatorJob,
		},
		AccessTable: table.FivenetDocumentsTemplatesAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetDocumentsTemplatesAccess.ID,
				TargetID: table.FivenetDocumentsTemplatesAccess.TargetID,
				Access:   table.FivenetDocumentsTemplatesAccess.Access,
			},
			SubjectID: table.FivenetDocumentsTemplatesAccess.SubjectID,
			Effect:    table.FivenetDocumentsTemplatesAccess.Effect,
		},
		CalculatedVisibilitySubjectTable:     table.FivenetDocumentsTemplatesVisibilitySubject,
		CalculatedVisibilitySubjectTargetID:  table.FivenetDocumentsTemplatesVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID: table.FivenetDocumentsTemplatesVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:    table.FivenetDocumentsTemplatesVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:    table.FivenetDocumentsTemplatesVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewDocumentStampsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetDocumentsStamps,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetDocumentsStamps.ID,
			CreatedAt:  table.FivenetDocumentsStamps.CreatedAt,
			UpdatedAt:  table.FivenetDocumentsStamps.UpdatedAt,
			DeletedAt:  table.FivenetDocumentsStamps.DeletedAt,
			CreatorJob: table.FivenetDocumentsStamps.Name,
		},
		AccessTable: table.FivenetDocumentsStampsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetDocumentsStampsAccess.ID,
				TargetID: table.FivenetDocumentsStampsAccess.TargetID,
				Access:   table.FivenetDocumentsStampsAccess.Access,
			},
			SubjectID: table.FivenetDocumentsStampsAccess.SubjectID,
			Effect:    table.FivenetDocumentsStampsAccess.Effect,
		},
		CalculatedVisibilityCreatorTable:      table.FivenetDocumentsStampsVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetDocumentsStampsVisibilitySubject,
		CalculatedVisibilityCreatorTargetID:   table.FivenetDocumentsStampsVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetDocumentsStampsVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetDocumentsStampsVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetDocumentsStampsVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetDocumentsStampsVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetDocumentsStampsVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{{Kind: VisibilityRuleCreator}},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewCalendarSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetCalendar,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetCalendar.ID,
			CreatedAt:  table.FivenetCalendar.CreatedAt,
			UpdatedAt:  table.FivenetCalendar.UpdatedAt,
			DeletedAt:  table.FivenetCalendar.DeletedAt,
			Public:     table.FivenetCalendar.Public,
			CreatorID:  table.FivenetCalendar.CreatorID,
			CreatorJob: table.FivenetCalendar.CreatorJob,
		},
		AccessTable: table.FivenetCalendarAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetCalendarAccess.ID,
				TargetID: table.FivenetCalendarAccess.TargetID,
				Access:   table.FivenetCalendarAccess.Access,
			},
			SubjectID: table.FivenetCalendarAccess.SubjectID,
			Effect:    table.FivenetCalendarAccess.Effect,
		},
		CalculatedVisibilityPublicTable:       table.FivenetCalendarVisibilityPublic,
		CalculatedVisibilityCreatorTable:      table.FivenetCalendarVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetCalendarVisibilitySubject,
		CalculatedVisibilityPublicTargetID:    table.FivenetCalendarVisibilityPublic.TargetID,
		CalculatedVisibilityCreatorTargetID:   table.FivenetCalendarVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID:  table.FivenetCalendarVisibilityCreator.CreatorID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetCalendarVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetCalendarVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetCalendarVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetCalendarVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetCalendarVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRulePublic},
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewWikiPageSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetWikiPages,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetWikiPages.ID,
			CreatedAt:  table.FivenetWikiPages.CreatedAt,
			UpdatedAt:  table.FivenetWikiPages.UpdatedAt,
			DeletedAt:  table.FivenetWikiPages.DeletedAt,
			Public:     table.FivenetWikiPages.Public,
			CreatorID:  table.FivenetWikiPages.CreatorID,
			CreatorJob: table.FivenetWikiPages.Job,
		},
		AccessTable: table.FivenetWikiPagesAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetWikiPagesAccess.ID,
				TargetID: table.FivenetWikiPagesAccess.TargetID,
				Access:   table.FivenetWikiPagesAccess.Access,
			},
			SubjectID: table.FivenetWikiPagesAccess.SubjectID,
			Effect:    table.FivenetWikiPagesAccess.Effect,
		},
		CalculatedVisibilityCreatorTable:      table.FivenetWikiPagesVisibilityCreator,
		CalculatedVisibilityPublicTable:       table.FivenetWikiPagesVisibilityPublic,
		CalculatedVisibilityCreatorTargetID:   table.FivenetWikiPagesVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID:  table.FivenetWikiPagesVisibilityCreator.CreatorID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetWikiPagesVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTable:      table.FivenetWikiPagesVisibilitySubject,
		CalculatedVisibilityPublicTargetID:    table.FivenetWikiPagesVisibilityPublic.TargetID,
		CalculatedVisibilitySubjectTargetID:   table.FivenetWikiPagesVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetWikiPagesVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetWikiPagesVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetWikiPagesVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRulePublic},
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewCitizenLabelsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetUserLabelsJob,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetUserLabelsJob.ID,
			CreatedAt:  table.FivenetUserLabelsJob.CreatedAt,
			UpdatedAt:  table.FivenetUserLabelsJob.UpdatedAt,
			DeletedAt:  table.FivenetUserLabelsJob.DeletedAt,
			CreatorJob: table.FivenetUserLabelsJob.Job,
		},
		AccessTable: table.FivenetUserLabelsJobJobAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetUserLabelsJobJobAccess.ID,
				TargetID: table.FivenetUserLabelsJobJobAccess.TargetID,
				Access:   table.FivenetUserLabelsJobJobAccess.Access,
			},
			SubjectID: table.FivenetUserLabelsJobJobAccess.SubjectID,
			Effect:    table.FivenetUserLabelsJobJobAccess.Effect,
		},
		CalculatedVisibilitySubjectTable:     table.FivenetUserLabelsJobVisibilitySubject,
		CalculatedVisibilitySubjectTargetID:  table.FivenetUserLabelsJobVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID: table.FivenetUserLabelsJobVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:    table.FivenetUserLabelsJobVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:    table.FivenetUserLabelsJobVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewQualificationsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetQualifications,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetQualifications.ID,
			CreatedAt:  table.FivenetQualifications.CreatedAt,
			UpdatedAt:  table.FivenetQualifications.UpdatedAt,
			DeletedAt:  table.FivenetQualifications.DeletedAt,
			Public:     table.FivenetQualifications.Public,
			CreatorID:  table.FivenetQualifications.CreatorID,
			CreatorJob: table.FivenetQualifications.CreatorJob,
		},
		AccessTable: table.FivenetQualificationsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetQualificationsAccess.ID,
				TargetID: table.FivenetQualificationsAccess.TargetID,
				Access:   table.FivenetQualificationsAccess.Access,
			},
			SubjectID: table.FivenetQualificationsAccess.SubjectID,
			Effect:    table.FivenetQualificationsAccess.Effect,
		},
		CalculatedVisibilityPublicTable:      table.FivenetQualificationsVisibilityPublic,
		CalculatedVisibilitySubjectTable:     table.FivenetQualificationsVisibilitySubject,
		CalculatedVisibilityPublicTargetID:   table.FivenetQualificationsVisibilityPublic.TargetID,
		CalculatedVisibilitySubjectTargetID:  table.FivenetQualificationsVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID: table.FivenetQualificationsVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:    table.FivenetQualificationsVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:    table.FivenetQualificationsVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRulePublic},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewMailerEmailsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetMailerEmails,
		TargetColumns: &SubjectTargetTableColumns{
			ID:        table.FivenetMailerEmails.ID,
			CreatedAt: table.FivenetMailerEmails.CreatedAt,
			UpdatedAt: table.FivenetMailerEmails.UpdatedAt,
			DeletedAt: table.FivenetMailerEmails.DeletedAt,
			CreatorID: table.FivenetMailerEmails.UserID,
		},
		AccessTable: table.FivenetMailerEmailsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetMailerEmailsAccess.ID,
				TargetID: table.FivenetMailerEmailsAccess.TargetID,
				Access:   table.FivenetMailerEmailsAccess.Access,
			},
			SubjectID: table.FivenetMailerEmailsAccess.SubjectID,
			Effect:    table.FivenetMailerEmailsAccess.Effect,
		},
		CalculatedVisibilityCreatorTable:     table.FivenetMailerEmailsVisibilityCreator,
		CalculatedVisibilitySubjectTable:     table.FivenetMailerEmailsVisibilitySubject,
		CalculatedVisibilityCreatorTargetID:  table.FivenetMailerEmailsVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID: table.FivenetMailerEmailsVisibilityCreator.CreatorID,
		// No CreatorJob check for emails
		CalculatedVisibilitySubjectTargetID:  table.FivenetMailerEmailsVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID: table.FivenetMailerEmailsVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:    table.FivenetMailerEmailsVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:    table.FivenetMailerEmailsVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewCentrumUnitsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetCentrumUnits,
		TargetColumns: &SubjectTargetTableColumns{
			ID:        table.FivenetCentrumUnits.ID,
			CreatedAt: table.FivenetCentrumUnits.CreatedAt,
			UpdatedAt: table.FivenetCentrumUnits.UpdatedAt,
			DeletedAt: table.FivenetCentrumUnits.DeletedAt,
		},
		AccessTable: table.FivenetCentrumUnitsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetCentrumUnitsAccess.ID,
				TargetID: table.FivenetCentrumUnitsAccess.TargetID,
				Access:   table.FivenetCentrumUnitsAccess.Access,
			},
			SubjectID: table.FivenetCentrumUnitsAccess.SubjectID,
			Effect:    table.FivenetCentrumUnitsAccess.Effect,
		},
	})
}
