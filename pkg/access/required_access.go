package access

import (
	"fmt"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
)

const defaultAccessEntryLimit = 20

type AccessEntryLimitError struct {
	Kind   string
	Max    int
	Actual int
}

func (e *AccessEntryLimitError) Error() string {
	return fmt.Sprintf("%s access entries exceed max %d (got %d)", e.Kind, e.Max, e.Actual)
}

func CloneAccess(in *resourcesaccess.Access) *resourcesaccess.Access {
	if in == nil {
		return &resourcesaccess.Access{}
	}

	return &resourcesaccess.Access{
		Jobs:           cloneJobAccessEntries(in.GetJobs()),
		Users:          cloneUserAccessEntries(in.GetUsers()),
		Qualifications: cloneQualificationAccessEntries(in.GetQualifications()),
	}
}

func NormalizeAccess(
	current *resourcesaccess.Access,
	required *resourcesaccess.Access,
	fallback *resourcesaccess.Access,
	maxEntries int,
) (*resourcesaccess.Access, error) {
	out, err := ApplyRequiredAccessOverlay(current, required, maxEntries)
	if err != nil {
		return nil, err
	}

	return EnsureMinimumAccess(out, fallback, maxEntries)
}

func ApplyRequiredAccessOverlay(
	current *resourcesaccess.Access,
	required *resourcesaccess.Access,
	maxEntries int,
) (*resourcesaccess.Access, error) {
	out := CloneAccess(current)
	if required == nil || required.IsEmpty() {
		return out, validateAccessLimit(out, maxEntries)
	}

	out.Jobs = mergeRequiredJobAccessEntries(out.GetJobs(), required.GetJobs())
	out.Users = mergeRequiredUserAccessEntries(out.GetUsers(), required.GetUsers())
	out.Qualifications = mergeRequiredQualificationAccessEntries(
		out.GetQualifications(),
		required.GetQualifications(),
	)

	if err := validateAccessLimit(out, maxEntries); err != nil {
		return nil, err
	}

	return out, nil
}

func EnsureMinimumAccess(
	current *resourcesaccess.Access,
	fallback *resourcesaccess.Access,
	maxEntries int,
) (*resourcesaccess.Access, error) {
	out := CloneAccess(current)
	if out.IsEmpty() && fallback != nil {
		out = CloneAccess(fallback)
	}

	if err := validateAccessLimit(out, maxEntries); err != nil {
		return nil, err
	}

	return out, nil
}

type accessEntryOps[T any, K comparable] struct {
	keyFn         func(T) K
	accessFn      func(T) int32
	setAccessFn   func(T, int32)
	requiredFn    func(T) bool
	setRequiredFn func(T, bool)
	cloneFn       func(T) T
}

func mergeRequiredEntries[T any, K comparable](
	current []T,
	required []T,
	ops accessEntryOps[T, K],
) []T {
	out := cloneEntries(current, ops.cloneFn)
	if len(required) == 0 {
		return out
	}

	indexByKey := make(map[K]int, len(out))
	for i, entry := range out {
		indexByKey[ops.keyFn(entry)] = i
	}

	for _, entry := range required {
		if !ops.requiredFn(entry) {
			continue
		}

		key := ops.keyFn(entry)
		if idx, ok := indexByKey[key]; ok {
			if ops.accessFn(out[idx]) < ops.accessFn(entry) {
				ops.setAccessFn(out[idx], ops.accessFn(entry))
			}
			ops.setRequiredFn(out[idx], true)
			continue
		}

		cloned := ops.cloneFn(entry)
		ops.setRequiredFn(cloned, true)
		out = append(out, cloned)
		indexByKey[key] = len(out) - 1
	}

	return out
}

func cloneEntries[T any](entries []T, cloneFn func(T) T) []T {
	if len(entries) == 0 {
		return []T{}
	}

	out := make([]T, 0, len(entries))
	for _, entry := range entries {
		out = append(out, cloneFn(entry))
	}

	return out
}

func mergeRequiredJobAccessEntries(
	current []*resourcesaccess.JobAccess,
	required []*resourcesaccess.JobAccess,
) []*resourcesaccess.JobAccess {
	return mergeRequiredEntries(
		current,
		required,
		accessEntryOps[*resourcesaccess.JobAccess, string]{
			keyFn: func(entry *resourcesaccess.JobAccess) string {
				return subjectJobAccessKey(entry.GetJob(), entry.GetMinimumGrade())
			},
			accessFn: func(entry *resourcesaccess.JobAccess) int32 {
				return entry.GetAccess()
			},
			setAccessFn: func(entry *resourcesaccess.JobAccess, access int32) {
				entry.SetAccess(access)
			},
			requiredFn: func(entry *resourcesaccess.JobAccess) bool {
				return entry.GetRequired()
			},
			setRequiredFn: func(entry *resourcesaccess.JobAccess, required bool) {
				entry.SetRequired(required)
			},
			cloneFn: cloneJobAccessEntry,
		},
	)
}

func mergeRequiredUserAccessEntries(
	current []*resourcesaccess.UserAccess,
	required []*resourcesaccess.UserAccess,
) []*resourcesaccess.UserAccess {
	return mergeRequiredEntries(
		current,
		required,
		accessEntryOps[*resourcesaccess.UserAccess, int32]{
			keyFn: func(entry *resourcesaccess.UserAccess) int32 {
				return entry.GetUserId()
			},
			accessFn: func(entry *resourcesaccess.UserAccess) int32 {
				return entry.GetAccess()
			},
			setAccessFn: func(entry *resourcesaccess.UserAccess, access int32) {
				entry.SetAccess(access)
			},
			requiredFn: func(entry *resourcesaccess.UserAccess) bool {
				return entry.GetRequired()
			},
			setRequiredFn: func(entry *resourcesaccess.UserAccess, required bool) {
				entry.SetRequired(required)
			},
			cloneFn: cloneUserAccessEntry,
		},
	)
}

func mergeRequiredQualificationAccessEntries(
	current []*resourcesaccess.QualificationAccess,
	required []*resourcesaccess.QualificationAccess,
) []*resourcesaccess.QualificationAccess {
	return mergeRequiredEntries(
		current,
		required,
		accessEntryOps[*resourcesaccess.QualificationAccess, int64]{
			keyFn: func(entry *resourcesaccess.QualificationAccess) int64 {
				return entry.GetQualificationId()
			},
			accessFn: func(entry *resourcesaccess.QualificationAccess) int32 {
				return entry.GetAccess()
			},
			setAccessFn: func(entry *resourcesaccess.QualificationAccess, access int32) {
				entry.SetAccess(access)
			},
			requiredFn: func(entry *resourcesaccess.QualificationAccess) bool {
				return entry.GetRequired()
			},
			setRequiredFn: func(entry *resourcesaccess.QualificationAccess, required bool) {
				entry.SetRequired(required)
			},
			cloneFn: cloneQualificationAccessEntry,
		},
	)
}

func cloneJobAccessEntry(in *resourcesaccess.JobAccess) *resourcesaccess.JobAccess {
	if in == nil {
		return nil
	}

	out := *in
	return &out
}

func cloneUserAccessEntry(in *resourcesaccess.UserAccess) *resourcesaccess.UserAccess {
	if in == nil {
		return nil
	}

	out := *in
	return &out
}

func cloneQualificationAccessEntry(in *resourcesaccess.QualificationAccess) *resourcesaccess.QualificationAccess {
	if in == nil {
		return nil
	}

	out := *in
	return &out
}

func cloneJobAccessEntries(entries []*resourcesaccess.JobAccess) []*resourcesaccess.JobAccess {
	return cloneEntries(entries, cloneJobAccessEntry)
}

func cloneUserAccessEntries(entries []*resourcesaccess.UserAccess) []*resourcesaccess.UserAccess {
	return cloneEntries(entries, cloneUserAccessEntry)
}

func cloneQualificationAccessEntries(
	entries []*resourcesaccess.QualificationAccess,
) []*resourcesaccess.QualificationAccess {
	return cloneEntries(entries, cloneQualificationAccessEntry)
}

func validateAccessLimit(access *resourcesaccess.Access, maxEntries int) error {
	if access == nil || maxEntries <= 0 {
		return nil
	}

	if len(access.GetJobs()) > maxEntries {
		return &AccessEntryLimitError{Kind: "jobs", Max: maxEntries, Actual: len(access.GetJobs())}
	}
	if len(access.GetUsers()) > maxEntries {
		return &AccessEntryLimitError{Kind: "users", Max: maxEntries, Actual: len(access.GetUsers())}
	}
	if len(access.GetQualifications()) > maxEntries {
		return &AccessEntryLimitError{
			Kind:   "qualifications",
			Max:    maxEntries,
			Actual: len(access.GetQualifications()),
		}
	}

	return nil
}
