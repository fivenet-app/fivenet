package perms

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms/collections"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

var (
	muPerms   sync.Mutex
	permsList = []*Perm{}
)

type (
	Category string
	Name     string
)

type Perm struct {
	Category Category
	Name     Name
	Attrs    []Attr
	Order    int32
}

type Attr struct {
	ID            uint64
	Key           Key
	Type          permissions.AttributeTypes
	ValidValues   any
	DefaultValues any
}

func AddPermsToList(perms []*Perm) {
	muPerms.Lock()
	defer muPerms.Unlock()

	permsList = append(permsList, perms...)
}

func BuildGuard(category Category, name Name) string {
	return Guard(string(category) + "." + string(name))
}

func BuildGuardWithKey(category Category, name Name, key Key) string {
	return Guard(string(category) + "." + string(name) + "." + string(key))
}

func (p *Perms) register(ctx context.Context, defaultRolePerms []string) error {
	if err := p.cleanupRoles(ctx); err != nil {
		return fmt.Errorf("failed to cleanup roles. %w", err)
	}

	for _, perm := range permsList {
		permId, err := p.createOrUpdatePermission(ctx, perm.Category, perm.Name, perm.Order)
		if err != nil {
			return fmt.Errorf(
				"failed to create or update permission (category: %s, name: %s). %w",
				perm.Category,
				perm.Name,
				err,
			)
		}
		p.permsMap.Store(permId, &cachePerm{
			ID:        permId,
			Category:  perm.Category,
			Name:      perm.Name,
			GuardName: BuildGuard(perm.Category, perm.Name),
			Order:     &perm.Order,
		})
		p.permsGuardToIDMap.Store(BuildGuard(perm.Category, perm.Name), permId)

		for _, attr := range perm.Attrs {
			if _, err := p.registerOrUpdateAttribute(ctx, permId, attr.Key, attr.Type, attr.ValidValues); err != nil {
				return fmt.Errorf(
					"failed to register/update attribute (perm id: %d, attr: %s). %w",
					permId,
					attr.Key,
					err,
				)
			}
		}
		p.logger.Debug(
			"registered permission",
			zap.String("guard", BuildGuard(perm.Category, perm.Name)),
		)
	}

	if err := p.SetDefaultRolePerms(ctx, defaultRolePerms); err != nil {
		return fmt.Errorf("failed to set default role perms. %w", err)
	}

	return nil
}

func (p *Perms) SetDefaultRolePerms(ctx context.Context, defaultPerms []string) error {
	if len(defaultPerms) == 0 {
		return nil
	}

	role, err := p.CreateRole(ctx, DefaultRoleJob, p.startJobGrade)
	if err != nil {
		return fmt.Errorf("failed to create role. %w", err)
	}

	addPerms := []AddPerm{}
	for _, perm := range defaultPerms {
		permId, ok := p.permsGuardToIDMap.Load(perm)
		if !ok {
			p.logger.Warn("default perm not found, skipping", zap.String("guard", perm))
			continue
		}

		addPerms = append(addPerms, AddPerm{
			Id:  permId,
			Val: true,
		})
	}

	currentPerms, err := p.GetRolePermissions(ctx, role.GetId())
	if err != nil {
		return fmt.Errorf("failed to get role permissions. %w", err)
	}

	removePerms := []uint64{}
	for _, p := range currentPerms {
		if slices.ContainsFunc(addPerms, func(ap AddPerm) bool {
			return ap.Id == p.GetId()
		}) {
			// Remove perm that are already set on the role
			addPerms = slices.DeleteFunc(addPerms, func(ap AddPerm) bool {
				return ap.Id == p.GetId() && ap.Val == p.GetVal()
			})
			continue
		}

		// Perm not in the default perms? Remove it!
		removePerms = append(removePerms, p.GetId())
	}

	if len(addPerms) > 0 {
		if err := p.UpdateRolePermissions(ctx, role.GetId(), addPerms...); err != nil {
			return fmt.Errorf("failed to update role permissions. %w", err)
		}
	}

	if len(removePerms) > 0 {
		if err := p.RemovePermissionsFromRole(ctx, role.GetId(), removePerms...); err != nil {
			return fmt.Errorf("failed to remove permissions from role. %w", err)
		}
	}

	if err := p.loadRolePermissions(ctx, role.GetId()); err != nil {
		return fmt.Errorf("failed to load role permissions. %w", err)
	}

	return nil
}

func (p *Perms) createOrUpdatePermission(
	ctx context.Context,
	category Category,
	name Name,
	order int32,
) (uint64, error) {
	perm, err := p.loadPermissionFromDatabaseByCategoryName(ctx, category, name)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, fmt.Errorf("failed to load permission from database by guard. %w", err)
		}
	}

	if perm != nil {
		if Category(perm.GetCategory()) != category || Name(perm.GetName()) != name ||
			(perm.Order == nil || perm.GetOrder() != order) {
			if err := p.UpdatePermission(ctx, perm.GetId(), category, name); err != nil {
				return perm.GetId(), fmt.Errorf("failed to update permission. %w", err)
			}
		}

		return perm.GetId(), nil
	}

	return p.CreatePermission(ctx, category, name)
}

func (p *Perms) registerOrUpdateAttribute(
	ctx context.Context,
	permId uint64,
	key Key,
	aType permissions.AttributeTypes,
	validValues any,
) (uint64, error) {
	attr, err := p.getAttributeFromDatabase(ctx, permId, key)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, fmt.Errorf("failed to get attribute from database. %w", err)
		}
	}

	attrValidValues := &permissions.AttributeValues{}
	if attr != nil && attr.ValidValues != nil {
		if err := p.convertRawValue(attrValidValues, *attr.ValidValues, aType); err != nil {
			return 0, fmt.Errorf("failed to convert raw value. %w", err)
		}
	}

	var validValsOut string
	// If the valid values is a nil or a string, don't do anything extra just set to an empty string
	if validValues != nil {
		switch v := validValues.(type) {
		case string:
			if v != "" {
				validValsOut = v
			}

		default:
			if aType == permissions.StringListAttributeType {
				marshaled, err := json.Marshal(v)
				if err != nil {
					return 0, fmt.Errorf("failed to marshal valid values to string. %w", err)
				}
				validValsOut = "{\"stringList\":{\"strings\":" + string(marshaled) + "}}"
			}
		}
	}
	if validValsOut == "" {
		validValsOut = "{}"
	}

	validVals := &permissions.AttributeValues{}
	if err := p.convertRawValue(validVals, validValsOut, aType); err != nil {
		return 0, fmt.Errorf("failed to convert raw value for valid values. %w", err)
	}

	if attr != nil && attr.ID > 0 {
		if err := p.addOrUpdateAttributeInMap(permId, attr.ID, key, aType, validVals); err != nil {
			return 0, fmt.Errorf("failed to add or update attribute in map. %w", err)
		}

		if attr.Type != string(aType) || (attr.ValidValues == nil || validVals != attrValidValues) {
			if err := p.UpdateAttribute(ctx, attr.ID, permId, key, aType, validVals); err != nil {
				return attr.ID, fmt.Errorf("failed to update attribute. %w", err)
			}
			return attr.ID, nil
		}

		return attr.ID, nil
	}

	attrId, err := p.CreateAttribute(ctx, permId, key, aType, validVals)
	if err != nil {
		return 0, fmt.Errorf("failed to create attribute. %w", err)
	}

	if err := p.addOrUpdateAttributeInMap(permId, attrId, key, aType, validVals); err != nil {
		return 0, fmt.Errorf("failed to add or update attribute in map after creation. %w", err)
	}

	return attrId, nil
}

func (p *Perms) cleanupRoles(ctx context.Context) error {
	if err := p.updateDefaultRole(ctx); err != nil {
		return err
	}

	tJobs := tables.Jobs().AS("job")
	tJobsGrades := tables.JobsGrades().AS("job_grade")

	stmt := tJobs.
		SELECT(
			tJobs.Name,
			tJobs.Label,
			tJobsGrades.JobName.AS("job_name"),
			tJobsGrades.Grade,
			tJobsGrades.Label,
		).
		FROM(
			tJobs.
				INNER_JOIN(tJobsGrades,
					tJobsGrades.JobName.EQ(tJobs.Name),
				),
		)

	var dest []*jobs.Job
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query jobs and job grades. %w", err)
		}
	}

	existingRoles, err := p.GetRoles(ctx, false)
	if err != nil {
		return fmt.Errorf("failed to get roles. %w", err)
	}

	// Remove default role from the list of roles
	existingRoles = slices.DeleteFunc(existingRoles, func(role *permissions.Role) bool {
		return role.GetJob() == DefaultRoleJob && role.GetGrade() == p.startJobGrade
	})

	// Iterate over current job and job grades to find any non-existent roles in our database
	for i := range dest {
		for _, grade := range dest[i].GetGrades() {
			role, err := p.GetRoleByJobAndGrade(ctx, dest[i].GetName(), grade.GetGrade())
			if err != nil {
				return fmt.Errorf("failed to get role by job and grade. %w", err)
			}
			if role == nil {
				continue
			}

			index := slices.IndexFunc(existingRoles, func(r *permissions.Role) bool {
				return role.GetId() == r.GetId()
			})
			if index > -1 {
				existingRoles = slices.Delete(existingRoles, index, index+1)
			}
		}
	}

	// Skip removal of roles if the cleanup is disabled
	if !p.cleanupRolesForMissingJobs {
		p.logger.Debug(
			"skipping cleanup of roles for missing jobs",
			zap.Any("roles", existingRoles),
		)
		return nil
	}

	// Remove all roles that shouldn't exist anymore
	for i := range existingRoles {
		if err := p.DeleteRole(ctx, existingRoles[i].GetId()); err != nil {
			return fmt.Errorf("failed to delete role. %w", err)
		}
	}

	return nil
}

func (p *Perms) updateDefaultRole(ctx context.Context) error {
	// Make sure the start job grade is set correctly
	stmt := tRoles.
		UPDATE(
			tRoles.Grade,
		).
		SET(
			tRoles.Grade.SET(jet.Int32(p.startJobGrade)),
		).
		WHERE(tRoles.Job.EQ(jet.String(DefaultRoleJob))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("failed to update default role grade. %w", err)
	}

	return nil
}

func (p *Perms) getActiveJobs(ctx context.Context) ([]string, error) {
	stmt := tRoles.
		SELECT(
			tRoles.Job,
		).
		FROM(tRoles).
		WHERE(
			tRoles.Job.NOT_EQ(jet.String(DefaultRoleJob)),
		).
		GROUP_BY(tRoles.Job)

	var dest []string
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query active jobs. %w", err)
		}
	}

	return dest, nil
}

func (p *Perms) ApplyJobPermissions(ctx context.Context, job string) error {
	jobs := []string{}
	if job != "" {
		jobs = append(jobs, job)
	} else {
		var err error
		jobs, err = p.getActiveJobs(ctx)
		if err != nil {
			return fmt.Errorf("failed to get active jobs. %w", err)
		}
	}

	for _, job := range jobs {
		if err := p.applyJobPermissions(ctx, job); err != nil {
			return fmt.Errorf("failed to apply job permissions for job %s. %w", job, err)
		}
	}

	if err := p.publishMessage(ctx, JobLimitsUpdatedSubject, &permissions.RoleIDEvent{
		Job: job,
	}); err != nil {
		return fmt.Errorf("failed to publish job perm update message for job %s. %w", job, err)
	}

	return nil
}

func (p *Perms) applyJobPermissions(ctx context.Context, job string) error {
	roles, err := p.GetJobRoles(ctx, job)
	if err != nil {
		return fmt.Errorf("failed to get job roles for job %s. %w", job, err)
	}

	if len(roles) == 0 {
		return nil
	}

	jps, err := p.GetJobPermissions(ctx, job)
	if err != nil {
		return fmt.Errorf("failed to get job permissions for job %s. %w", job, err)
	}

	for _, role := range roles {
		ps, err := p.GetRolePermissions(ctx, role.GetId())
		if err != nil {
			return fmt.Errorf(
				"failed to get role permissions for role ID %d. %w",
				role.GetId(),
				err,
			)
		}

		if len(ps) == 0 {
			continue
		}

		if len(jps) == 0 {
			p.logger.Debug(
				"removing all job permissions from role due to job perms change",
				zap.String("job", job),
			)
			pIds := []uint64{}
			for _, p := range ps {
				pIds = append(pIds, p.GetId())
			}

			if err := p.RemovePermissionsFromRole(ctx, role.GetId(), pIds...); err != nil {
				return fmt.Errorf(
					"failed to remove permissions from role ID %d. %w",
					role.GetId(),
					err,
				)
			}
			continue
		}

		toRemove := []uint64{}
		for _, p := range ps {
			if !slices.ContainsFunc(jps, func(in *permissions.Permission) bool {
				return in.GetId() == p.GetId() && in.GetVal()
			}) {
				toRemove = append(toRemove, p.GetId())
			}
		}

		if len(toRemove) > 0 {
			p.logger.Debug(
				"removing permissions from role due to job perms change",
				zap.String("job", job),
				zap.Int("perms_length", len(toRemove)),
			)
			if err := p.RemovePermissionsFromRole(ctx, role.GetId(), toRemove...); err != nil {
				return fmt.Errorf(
					"failed to remove permissions from role ID %d. %w",
					role.GetId(),
					err,
				)
			}
		}
	}

	if err := p.applyJobPermissionsToAttrs(ctx, roles, jps); err != nil {
		return fmt.Errorf("failed to apply job permissions to attributes for job %s. %w", job, err)
	}

	return nil
}

func (p *Perms) applyJobPermissionsToAttrs(
	ctx context.Context,
	roles collections.Roles,
	jps []*permissions.Permission,
) error {
	if len(roles) == 0 {
		return nil
	}

	for _, role := range roles {
		attrs, err := p.GetRoleAttributes(ctx, role.GetJob(), role.GetGrade())
		if err != nil {
			return fmt.Errorf(
				"failed to get role attributes for job %s and grade %d. %w",
				role.GetJob(),
				role.GetGrade(),
				err,
			)
		}

		if len(attrs) == 0 {
			continue
		}

		if len(jps) == 0 {
			p.logger.Debug(
				"removing all attributes from role due to job perms change",
				zap.String("job", role.GetJob()),
			)
			if err := p.RemoveAttributesFromRole(ctx, role.GetId(), attrs...); err != nil {
				return fmt.Errorf(
					"failed to remove all attributes from role %d. %w",
					role.GetId(),
					err,
				)
			}
			continue
		}

		toRemove := []*permissions.RoleAttribute{}
		toUpdate := []*permissions.RoleAttribute{}
		for _, attr := range attrs {
			maxValues, _ := p.GetJobAttributeValue(ctx, role.GetJob(), attr.GetAttrId())

			if slices.ContainsFunc(jps, func(in *permissions.Permission) bool {
				return in.GetId() == attr.GetPermissionId()
			}) {
				if _, changed := attr.GetValue().Check(permissions.AttributeTypes(attr.GetType()), attr.GetValidValues(), maxValues); changed {
					p.logger.Debug(
						"attribute changed on role due to job perms change",
						zap.String("job", role.GetJob()),
						zap.Uint64("attr_id", attr.GetAttrId()),
						zap.Any("attr_value", attr.GetValue()),
						zap.Any("attr_valid_value", attr.GetValidValues()),
						zap.Any("attr_max_values", maxValues),
					)
					toUpdate = append(toUpdate, attr)
				} else {
					p.logger.Debug("attribute not changed on role due to job perms change", zap.String("job", role.GetJob()), zap.Uint64("attr_id", attr.GetAttrId()), zap.Any("attr_value", attr.GetValue()), zap.Any("attr_valid_value", attr.GetValidValues()), zap.Any("attr_max_values", maxValues))
				}
			} else {
				toRemove = append(toRemove, attr)
			}
		}

		if len(toRemove) > 0 {
			p.logger.Debug(
				"removing attribute from role due to job perms change",
				zap.String("job", role.GetJob()),
				zap.Int("perms_length", len(toRemove)),
			)
			if err := p.RemoveAttributesFromRole(ctx, role.GetId(), toRemove...); err != nil {
				return fmt.Errorf("failed to remove attributes from role %d. %w", role.GetId(), err)
			}
		}

		if len(toUpdate) > 0 {
			p.logger.Debug(
				"updating attribute on role due to job perms change",
				zap.String("job", role.GetJob()),
				zap.Int("perms_length", len(toUpdate)),
			)
			if err := p.UpdateRoleAttributes(ctx, role.GetJob(), role.GetId(), toUpdate...); err != nil {
				return fmt.Errorf("failed to update attributes for role %d. %w", role.GetId(), err)
			}
		}
	}

	return nil
}
