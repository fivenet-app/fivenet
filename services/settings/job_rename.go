package settings

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorssettings "github.com/fivenet-app/fivenet/v2026/services/settings/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

// jobReferenceColumns is deliberately explicit. Adding a job reference to the
// schema must add it here, rather than making a privileged API discover and
// execute SQL against arbitrary database metadata.
type jobReferenceColumn struct {
	relation mysql.Table
	column   mysql.ColumnString
}

var jobReferenceColumns = []jobReferenceColumn{
	{table.FivenetCalendar, table.FivenetCalendar.Job},
	{table.FivenetCalendar, table.FivenetCalendar.CreatorJob},
	{table.FivenetCalendarEntries, table.FivenetCalendarEntries.Job},
	{table.FivenetCalendarEntries, table.FivenetCalendarEntries.CreatorJob},
	{table.FivenetCalendarVisibilityCreator, table.FivenetCalendarVisibilityCreator.CreatorJob},
	{table.FivenetCentrumDispatchesHeatmaps, table.FivenetCentrumDispatchesHeatmaps.Job},
	{table.FivenetCentrumDispatchesStatus, table.FivenetCentrumDispatchesStatus.CreatorJob},
	{table.FivenetCentrumJobAccess, table.FivenetCentrumJobAccess.SourceJob},
	{table.FivenetCentrumJobAccess, table.FivenetCentrumJobAccess.Job},
	{table.FivenetCentrumMarkers, table.FivenetCentrumMarkers.Job},
	{table.FivenetCentrumSettings, table.FivenetCentrumSettings.Job},
	{table.FivenetCentrumUnits, table.FivenetCentrumUnits.Job},
	{table.FivenetCentrumUnitsStatus, table.FivenetCentrumUnitsStatus.CreatorJob},
	{table.FivenetCentrumUserLocations, table.FivenetCentrumUserLocations.Job},
	{table.FivenetDocuments, table.FivenetDocuments.CreatorJob},
	{table.FivenetDocumentsActivity, table.FivenetDocumentsActivity.CreatorJob},
	{table.FivenetDocumentsApprovalTasks, table.FivenetDocumentsApprovalTasks.Job},
	{table.FivenetDocumentsApprovalTasks, table.FivenetDocumentsApprovalTasks.CreatorJob},
	{table.FivenetDocumentsApprovals, table.FivenetDocumentsApprovals.UserJob},
	{table.FivenetDocumentsCategories, table.FivenetDocumentsCategories.Job},
	{table.FivenetDocumentsComments, table.FivenetDocumentsComments.CreatorJob},
	{table.FivenetDocumentsPins, table.FivenetDocumentsPins.Job},
	{table.FivenetDocumentsRequests, table.FivenetDocumentsRequests.CreatorJob},
	{
		table.FivenetDocumentsStampsVisibilityCreator,
		table.FivenetDocumentsStampsVisibilityCreator.CreatorJob,
	},
	{table.FivenetDocumentsStatsMetric, table.FivenetDocumentsStatsMetric.Job},
	{table.FivenetDocumentsTemplates, table.FivenetDocumentsTemplates.CreatorJob},
	{table.FivenetDocumentsVisibilityCreator, table.FivenetDocumentsVisibilityCreator.CreatorJob},
	{table.FivenetJobAssets, table.FivenetJobAssets.Job},
	{table.FivenetJobColleagueActivity, table.FivenetJobColleagueActivity.Job},
	{table.FivenetJobColleagueLabels, table.FivenetJobColleagueLabels.Job},
	{table.FivenetJobColleagueProps, table.FivenetJobColleagueProps.Job},
	{table.FivenetJobConduct, table.FivenetJobConduct.Job},
	{table.FivenetJobGroupActivity, table.FivenetJobGroupActivity.Job},
	{table.FivenetJobGroups, table.FivenetJobGroups.Job},
	{table.FivenetJobLabels, table.FivenetJobLabels.Job},
	{table.FivenetJobProps, table.FivenetJobProps.Job},
	{table.FivenetJobTimeclock, table.FivenetJobTimeclock.Job},
	{table.FivenetMailerEmails, table.FivenetMailerEmails.Job},
	{
		table.FivenetMailerEmailsVisibilityCreator,
		table.FivenetMailerEmailsVisibilityCreator.CreatorJob,
	},
	{table.FivenetMailerMessages, table.FivenetMailerMessages.CreatorJob},
	{table.FivenetMailerTemplates, table.FivenetMailerTemplates.CreatorJob},
	{table.FivenetNotifications, table.FivenetNotifications.Job},
	{table.FivenetQualifications, table.FivenetQualifications.Job},
	{table.FivenetQualifications, table.FivenetQualifications.CreatorJob},
	{table.FivenetQualificationsRequests, table.FivenetQualificationsRequests.ApproverJob},
	{table.FivenetQualificationsResults, table.FivenetQualificationsResults.CreatorJob},
	{table.FivenetStatsDailyRollup, table.FivenetStatsDailyRollup.Job},
	{table.FivenetUser, table.FivenetUser.Job},
	{table.FivenetUserJobs, table.FivenetUserJobs.Job},
	{table.FivenetUserLabelsJob, table.FivenetUserLabelsJob.Job},
	{table.FivenetVehicle, table.FivenetVehicle.Job},
	{table.FivenetVehicleActivity, table.FivenetVehicleActivity.CreatorJob},
	{table.FivenetWikiPages, table.FivenetWikiPages.Job},
	{table.FivenetWikiPagesActivity, table.FivenetWikiPagesActivity.CreatorJob},
	{table.FivenetWikiPagesVisibilityCreator, table.FivenetWikiPagesVisibilityCreator.CreatorJob},
}

type aclSubjectReferenceTable struct {
	relation  mysql.Table
	subjectID mysql.ColumnInteger
}

var aclSubjectReferenceTables = []aclSubjectReferenceTable{
	{table.FivenetCalendarAccess, table.FivenetCalendarAccess.SubjectID},
	{table.FivenetCalendarVisibilitySubject, table.FivenetCalendarVisibilitySubject.SubjectID},
	{table.FivenetCentrumUnitsAccess, table.FivenetCentrumUnitsAccess.SubjectID},
	{table.FivenetDocumentsAccess, table.FivenetDocumentsAccess.SubjectID},
	{table.FivenetDocumentsStampsAccess, table.FivenetDocumentsStampsAccess.SubjectID},
	{
		table.FivenetDocumentsStampsVisibilitySubject,
		table.FivenetDocumentsStampsVisibilitySubject.SubjectID,
	},
	{table.FivenetDocumentsTemplatesAccess, table.FivenetDocumentsTemplatesAccess.SubjectID},
	{
		table.FivenetDocumentsTemplatesVisibilitySubject,
		table.FivenetDocumentsTemplatesVisibilitySubject.SubjectID,
	},
	{table.FivenetDocumentsVisibilitySubject, table.FivenetDocumentsVisibilitySubject.SubjectID},
	{table.FivenetJobGroupsAccess, table.FivenetJobGroupsAccess.SubjectID},
	{table.FivenetJobGroupsVisibilitySubject, table.FivenetJobGroupsVisibilitySubject.SubjectID},
	{table.FivenetMailerEmailsAccess, table.FivenetMailerEmailsAccess.SubjectID},
	{
		table.FivenetMailerEmailsVisibilitySubject,
		table.FivenetMailerEmailsVisibilitySubject.SubjectID,
	},
	{table.FivenetQualificationsAccess, table.FivenetQualificationsAccess.SubjectID},
	{
		table.FivenetQualificationsVisibilitySubject,
		table.FivenetQualificationsVisibilitySubject.SubjectID,
	},
	{table.FivenetUserLabelsJobJobAccess, table.FivenetUserLabelsJobJobAccess.SubjectID},
	{
		table.FivenetUserLabelsJobVisibilitySubject,
		table.FivenetUserLabelsJobVisibilitySubject.SubjectID,
	},
	{table.FivenetWikiPagesAccess, table.FivenetWikiPagesAccess.SubjectID},
	{table.FivenetWikiPagesVisibilitySubject, table.FivenetWikiPagesVisibilitySubject.SubjectID},
}

func (s *Server) RenameJob(
	ctx context.Context,
	req *pbsettings.RenameJobRequest,
) (*pbsettings.RenameJobResponse, error) {
	if req.GetSourceJob() == req.GetTargetJob() {
		return &pbsettings.RenameJobResponse{}, nil
	}
	// Job keys are stored under case-insensitive MySQL collations in typical
	// installations. Treat a case-only change as invalid rather than allowing
	// it to collapse the source and target into the same database key.
	if strings.EqualFold(req.GetSourceJob(), req.GetTargetJob()) {
		return nil, errorssettings.ErrInvalidRequest
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	defer tx.Rollback() //nolint:errcheck // Commit below makes this a no-op.

	if err := renameJob(ctx, tx, req.GetSourceJob(), req.GetTargetJob()); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorssettings.ErrInvalidRequest
		}
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	// The mutation is durable at this point. Cache refreshes are best effort: a
	// failed refresh must not make callers retry an already-completed rename.
	if err := s.perms.ReloadJob(ctx, req.GetSourceJob()); err != nil {
		s.logger.Warn("failed to reload source job permissions after rename", zap.Error(err))
	}
	if err := s.perms.ReloadJob(ctx, req.GetTargetJob()); err != nil {
		s.logger.Warn("failed to reload target job permissions after rename", zap.Error(err))
	}
	// Job-list attributes can belong to roles of other jobs. Reload every
	// remaining job so their in-memory attribute views observe the rewrite.
	roles, err := s.perms.GetRoles(ctx, false)
	if err != nil {
		s.logger.Warn("failed to list roles after job rename", zap.Error(err))
	} else {
		jobs := map[string]struct{}{}
		for _, role := range roles {
			jobs[role.GetJob()] = struct{}{}
		}
		for job := range jobs {
			if err := s.perms.ReloadJob(ctx, job); err != nil {
				s.logger.Warn("failed to reload job permissions after rename", zap.Error(err))
			}
		}
	}
	if s.jobs != nil {
		if err := s.jobs.Refresh(ctx); err != nil {
			s.logger.Warn("failed to refresh jobs after rename", zap.Error(err))
		}
	}

	return &pbsettings.RenameJobResponse{}, nil
}

func renameJob(ctx context.Context, tx *sql.Tx, source, target string) error {
	tJobs := table.FivenetJobs
	tJobGrades := table.FivenetJobsGrades

	job := struct {
		Label sql.NullString `alias:"label"`
	}{}
	if err := tJobs.SELECT(tJobs.Label).
		FROM(tJobs).
		WHERE(tJobs.Name.EQ(mysql.String(source))).
		LIMIT(1).
		FOR(mysql.UPDATE()).
		QueryContext(ctx, tx, &job); err != nil {
		return err
	}
	if _, err := tJobs.INSERT(tJobs.Name, tJobs.Label).
		VALUES(target, job.Label).
		ON_DUPLICATE_KEY_UPDATE(tJobs.Name.SET(tJobs.Name)).
		ExecContext(ctx, tx); err != nil {
		return err
	}
	// Target grades are authoritative if the upstream already created them.
	if _, err := tJobGrades.INSERT(tJobGrades.JobName, tJobGrades.Grade, tJobGrades.Label).
		QUERY(tJobGrades.SELECT(mysql.String(target), tJobGrades.Grade, tJobGrades.Label).
			FROM(tJobGrades).
			WHERE(tJobGrades.JobName.EQ(mysql.String(source)))).
		ON_DUPLICATE_KEY_UPDATE(tJobGrades.JobName.SET(tJobGrades.JobName)).
		ExecContext(ctx, tx); err != nil {
		return err
	}
	if _, err := tJobGrades.DELETE().
		WHERE(tJobGrades.JobName.EQ(mysql.String(source))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	if err := mergeJobLimits(ctx, tx, source, target); err != nil {
		return err
	}
	if err := mergeRoles(ctx, tx, source, target); err != nil {
		return err
	}
	if err := rewriteRoleJobAttributes(ctx, tx, source, target); err != nil {
		return err
	}
	if err := mergeACLSubjects(ctx, tx, source, target); err != nil {
		return err
	}
	for _, ref := range jobReferenceColumns {
		// Some job-scoped tables have a unique key containing the job name. A
		// target job may already have the equivalent row, in which case target
		// data wins and the remaining source duplicate is removed below.
		_, err := ref.relation.UPDATE(ref.column).
			SET(mysql.String(target)).
			WHERE(ref.column.EQ(mysql.String(source))).
			ExecContext(ctx, tx)
		if err != nil && !dbutils.IsDuplicateError(err) {
			return fmt.Errorf("update %s.%s: %w", ref.relation.TableName(), ref.column.Name(), err)
		}
		if err != nil {
			// Jet does not expose MySQL's UPDATE IGNORE. Keep this narrowly
			// scoped fallback to preserve target-wins merge semantics.
			if _, err := tx.ExecContext(
				ctx,
				fmt.Sprintf(
					"UPDATE IGNORE `%s` SET `%s` = ? WHERE `%s` = ?",
					ref.relation.TableName(),
					ref.column.Name(),
					ref.column.Name(),
				),
				target,
				source,
			); err != nil {
				return fmt.Errorf(
					"merge duplicate %s.%s: %w",
					ref.relation.TableName(),
					ref.column.Name(),
					err,
				)
			}
		}
		if _, err := ref.relation.DELETE().
			WHERE(ref.column.EQ(mysql.String(source))).
			ExecContext(ctx, tx); err != nil {
			return fmt.Errorf(
				"remove duplicate %s.%s: %w",
				ref.relation.TableName(),
				ref.column.Name(),
				err,
			)
		}
	}
	_, err := tJobs.DELETE().WHERE(tJobs.Name.EQ(mysql.String(source))).ExecContext(ctx, tx)
	return err
}

func mergeJobLimits(ctx context.Context, tx *sql.Tx, source, target string) error {
	tJobPerms := table.FivenetRbacJobPermissions
	if _, err := tJobPerms.INSERT(tJobPerms.Job, tJobPerms.PermissionID, tJobPerms.Val).
		QUERY(tJobPerms.SELECT(mysql.String(target), tJobPerms.PermissionID, tJobPerms.Val).
			FROM(tJobPerms).
			WHERE(tJobPerms.Job.EQ(mysql.String(source)))).
		ON_DUPLICATE_KEY_UPDATE(tJobPerms.Job.SET(tJobPerms.Job)).
		ExecContext(ctx, tx); err != nil {
		return err
	}
	if _, err := tJobPerms.DELETE().
		WHERE(tJobPerms.Job.EQ(mysql.String(source))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	tJobAttrs := table.FivenetRbacJobAttrs
	if _, err := tJobAttrs.INSERT(tJobAttrs.Job, tJobAttrs.AttrID, tJobAttrs.MaxValues).
		QUERY(tJobAttrs.SELECT(mysql.String(target), tJobAttrs.AttrID, tJobAttrs.MaxValues).
			FROM(tJobAttrs).
			WHERE(tJobAttrs.Job.EQ(mysql.String(source)))).
		ON_DUPLICATE_KEY_UPDATE(tJobAttrs.Job.SET(tJobAttrs.Job)).
		ExecContext(ctx, tx); err != nil {
		return err
	}
	if _, err := tJobAttrs.DELETE().
		WHERE(tJobAttrs.Job.EQ(mysql.String(source))).
		ExecContext(ctx, tx); err != nil {
		return err
	}
	return nil
}

func mergeRoles(ctx context.Context, tx *sql.Tx, source, target string) error {
	tRoles := table.FivenetRbacRoles
	tRolePerms := table.FivenetRbacRolesPermissions
	tRoleAttrs := table.FivenetRbacRolesAttrs

	type role struct {
		ID    int64 `alias:"id"`
		Grade int32 `alias:"grade"`
	}
	var sourceRoles []role
	if err := tRoles.SELECT(tRoles.ID, tRoles.Grade).
		FROM(tRoles).
		WHERE(tRoles.Job.EQ(mysql.String(source))).
		FOR(mysql.UPDATE()).
		QueryContext(ctx, tx, &sourceRoles); err != nil {
		return err
	}
	for _, sourceRole := range sourceRoles {
		sourceID, grade := sourceRole.ID, sourceRole.Grade
		targetRole := struct {
			ID int64 `alias:"id"`
		}{}
		err := tRoles.SELECT(tRoles.ID).
			FROM(tRoles).
			WHERE(mysql.AND(
				tRoles.Job.EQ(mysql.String(target)),
				tRoles.Grade.EQ(mysql.Int32(grade)),
			)).
			LIMIT(1).
			FOR(mysql.UPDATE()).
			QueryContext(ctx, tx, &targetRole)
		if errors.Is(err, qrm.ErrNoRows) {
			if _, err = tRoles.UPDATE(tRoles.Job).
				SET(mysql.String(target)).
				WHERE(tRoles.ID.EQ(mysql.Int64(sourceID))).
				ExecContext(ctx, tx); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		targetID := targetRole.ID
		if _, err := tRolePerms.INSERT(tRolePerms.RoleID, tRolePerms.PermissionID, tRolePerms.Val).
			QUERY(tRolePerms.SELECT(mysql.Int64(targetID), tRolePerms.PermissionID, tRolePerms.Val).
				FROM(tRolePerms).
				WHERE(tRolePerms.RoleID.EQ(mysql.Int64(sourceID)))).
			ON_DUPLICATE_KEY_UPDATE(tRolePerms.RoleID.SET(tRolePerms.RoleID)).
			ExecContext(ctx, tx); err != nil {
			return err
		}
		if _, err := tRoleAttrs.INSERT(tRoleAttrs.RoleID, tRoleAttrs.AttrID, tRoleAttrs.Value).
			QUERY(tRoleAttrs.SELECT(mysql.Int64(targetID), tRoleAttrs.AttrID, tRoleAttrs.Value).
				FROM(tRoleAttrs).
				WHERE(tRoleAttrs.RoleID.EQ(mysql.Int64(sourceID)))).
			ON_DUPLICATE_KEY_UPDATE(tRoleAttrs.RoleID.SET(tRoleAttrs.RoleID)).
			ExecContext(ctx, tx); err != nil {
			return err
		}
		if _, err := tRoles.DELETE().
			WHERE(tRoles.ID.EQ(mysql.Int64(sourceID))).
			ExecContext(ctx, tx); err != nil {
			return err
		}
	}
	return nil
}

func rewriteRoleJobAttributes(ctx context.Context, tx *sql.Tx, source, target string) error {
	tRoleAttrs := table.FivenetRbacRolesAttrs
	tAttrs := table.FivenetRbacAttrs

	type roleAttribute struct {
		RoleID int64  `alias:"role_id"`
		AttrID int64  `alias:"attr_id"`
		Type   string `alias:"type"`
		Value  string `alias:"value"`
	}
	var attrs []roleAttribute
	if err := tRoleAttrs.SELECT(tRoleAttrs.RoleID, tRoleAttrs.AttrID, tAttrs.Type, tRoleAttrs.Value).
		FROM(tRoleAttrs.INNER_JOIN(tAttrs, tAttrs.ID.EQ(tRoleAttrs.AttrID))).
		WHERE(tAttrs.Type.IN(mysql.String("JobList"), mysql.String("JobGradeList"))).
		FOR(mysql.UPDATE()).
		QueryContext(ctx, tx, &attrs); err != nil {
		return err
	}
	for _, attr := range attrs {
		roleID, attrID, typ, value := attr.RoleID, attr.AttrID, attr.Type, attr.Value
		v := &permissionsattributes.AttributeValues{}
		if err := protojson.Unmarshal([]byte(value), v); err != nil {
			return err
		}
		changed := false
		switch typ {
		case "JobList":
			list := v.GetJobList()
			if list == nil {
				continue
			}
			for i, job := range list.Strings {
				if job == source {
					list.Strings[i] = target
					changed = true
				}
			}
			list.Strings = slices.Compact(slices.Sorted(slices.Values(list.Strings)))
		case "JobGradeList":
			list := v.GetJobGradeList()
			if list == nil {
				continue
			}
			if list.Jobs == nil {
				list.Jobs = map[string]int32{}
			}
			if list.Grades == nil {
				list.Grades = map[string]*permissionsattributes.JobGrades{}
			}
			if grade, ok := list.Jobs[source]; ok {
				if grade > list.Jobs[target] {
					list.Jobs[target] = grade
				}
				delete(list.Jobs, source)
				changed = true
			}
			if grades, ok := list.Grades[source]; ok {
				dest := list.Grades[target]
				if dest == nil {
					dest = &permissionsattributes.JobGrades{}
					list.Grades[target] = dest
				}
				dest.Grades = slices.Compact(
					slices.Sorted(slices.Values(append(dest.Grades, grades.Grades...))),
				)
				delete(list.Grades, source)
				changed = true
			}
		}
		if !changed {
			continue
		}
		out, err := protoutils.MarshalToJSON(v)
		if err != nil {
			return err
		}
		if _, err := tRoleAttrs.UPDATE(tRoleAttrs.Value).
			SET(mysql.String(string(out))).
			WHERE(mysql.AND(
				tRoleAttrs.RoleID.EQ(mysql.Int64(roleID)),
				tRoleAttrs.AttrID.EQ(mysql.Int64(attrID)),
			)).
			ExecContext(ctx, tx); err != nil {
			return err
		}
	}
	return nil
}

func mergeACLSubjects(ctx context.Context, tx *sql.Tx, source, target string) error {
	tScopes := table.FivenetACLSubjectJobGradeScopes
	tSubjects := table.FivenetACLSubjects

	type scope struct {
		ID    int64 `alias:"subject_id"`
		Grade int32 `alias:"minimum_grade"`
	}
	var sourceScopes []scope
	if err := tScopes.SELECT(tScopes.SubjectID, tScopes.MinimumGrade).
		FROM(tScopes).
		WHERE(tScopes.Job.EQ(mysql.String(source))).
		FOR(mysql.UPDATE()).
		QueryContext(ctx, tx, &sourceScopes); err != nil {
		return err
	}
	for _, scope := range sourceScopes {
		sourceID, grade := scope.ID, scope.Grade
		targetScope := struct {
			ID int64 `alias:"subject_id"`
		}{}
		err := tScopes.SELECT(tScopes.SubjectID).
			FROM(tScopes).
			WHERE(mysql.AND(
				tScopes.Job.EQ(mysql.String(target)),
				tScopes.MinimumGrade.EQ(mysql.Int32(grade)),
			)).
			LIMIT(1).
			FOR(mysql.UPDATE()).
			QueryContext(ctx, tx, &targetScope)
		if errors.Is(err, qrm.ErrNoRows) {
			if _, err = tScopes.UPDATE(tScopes.Job).
				SET(mysql.String(target)).
				WHERE(tScopes.SubjectID.EQ(mysql.Int64(sourceID))).
				ExecContext(ctx, tx); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		targetID := targetScope.ID
		for _, relation := range aclSubjectReferenceTables {
			_, err := relation.relation.UPDATE(relation.subjectID).
				SET(mysql.Int64(targetID)).
				WHERE(relation.subjectID.EQ(mysql.Int64(sourceID))).
				ExecContext(ctx, tx)
			if err != nil && !dbutils.IsDuplicateError(err) {
				return fmt.Errorf(
					"update ACL subject reference %s: %w",
					relation.relation.TableName(),
					err,
				)
			}
			if err != nil {
				// See the matching job-reference fallback above.
				if _, err := tx.ExecContext(
					ctx,
					fmt.Sprintf(
						"UPDATE IGNORE `%s` SET subject_id = ? WHERE subject_id = ?",
						relation.relation.TableName(),
					),
					targetID,
					sourceID,
				); err != nil {
					return fmt.Errorf(
						"merge duplicate ACL subject reference %s: %w",
						relation.relation.TableName(),
						err,
					)
				}
			}
			if _, err := relation.relation.DELETE().
				WHERE(relation.subjectID.EQ(mysql.Int64(sourceID))).
				ExecContext(ctx, tx); err != nil {
				return fmt.Errorf(
					"remove duplicate ACL subject reference %s: %w",
					relation.relation.TableName(),
					err,
				)
			}
		}
		if _, err := tSubjects.DELETE().
			WHERE(tSubjects.ID.EQ(mysql.Int64(sourceID))).
			ExecContext(ctx, tx); err != nil {
			return err
		}
	}
	return nil
}
