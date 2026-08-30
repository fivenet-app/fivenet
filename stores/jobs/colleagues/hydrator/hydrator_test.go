package colleagueshydrator

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	permsjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/permsstub"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

type testPermissions struct {
	permsstub.Permissions

	fields []string
	calls  int
}

func (p *testPermissions) AttrStringList(
	_ *userinfo.UserInfo,
	_ perms.AttrRef[perms.StringListAttr],
) (*permissionsattributes.StringList, error) {
	p.calls++
	return &permissionsattributes.StringList{Strings: p.fields}, nil
}

type testUserAwareEnricher struct {
	mstlystcdata.DummyEnricher
}

func (e *testUserAwareEnricher) EnrichJobInfoSafe(
	_ *userinfo.UserInfo,
	usrs ...common.IJobInfo,
) {
	for _, usr := range usrs {
		e.enrich(usr)
	}
}

func (e *testUserAwareEnricher) EnrichJobInfoSafeFunc(
	_ *userinfo.UserInfo,
) func(usr common.IJobInfo) {
	return func(usr common.IJobInfo) {
		e.enrich(usr)
	}
}

func (e *testUserAwareEnricher) enrich(usr common.IJobInfo) {
	usr.SetJobLabel("label:" + usr.GetJob())
	usr.SetJobGradeLabel(fmt.Sprintf("grade:%d", usr.GetJobGrade()))
}

type testColleagueStore struct {
	*jobsstore.Store

	colleagues []*jobscolleagues.Colleague
	deleted    map[int32]struct{}
}

func (s *testColleagueStore) ListColleaguesByUserIDs(
	_ context.Context,
	_ qrm.DB,
	q jobsstore.ListColleaguesByUserIDsQuery,
) ([]*jobscolleagues.Colleague, error) {
	byID := make(map[int32]*jobscolleagues.Colleague, len(s.colleagues))
	for _, colleague := range s.colleagues {
		byID[colleague.GetUserId()] = proto.Clone(colleague).(*jobscolleagues.Colleague)
	}

	out := make([]*jobscolleagues.Colleague, 0, len(q.UserIDs))
	for _, userID := range q.UserIDs {
		if _, ok := s.deleted[userID]; ok {
			continue
		}
		if colleague, ok := byID[userID]; ok {
			out = append(out, colleague)
		}
	}

	return out, nil
}

func newTestHydrator(
	t *testing.T,
	colleagues ...*jobscolleagues.Colleague,
) (*Hydrator, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return &Hydrator{
		db:       db,
		enricher: &testUserAwareEnricher{},
		store: &testColleagueStore{
			Store:      &jobsstore.Store{},
			colleagues: colleagues,
			deleted:    map[int32]struct{}{},
		},
	}, mock
}

func TestResolveFieldsRequiresExplicitRequestAndPermission(t *testing.T) {
	t.Parallel()

	permissions := &testPermissions{fields: []string{
		string(permsjobs.ColleaguesServiceGetColleagueTypesPermValueNote),
		string(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels),
	}}
	h := &Hydrator{perms: permissions}
	userInfo := &userinfo.UserInfo{}

	fields, _, err := h.resolveFields(userInfo, ResolveOpts{})
	if err != nil {
		t.Fatalf("resolveFields returned error: %v", err)
	}
	if fields.note || fields.labels {
		t.Fatalf("unexpected fields without explicit request: %+v", fields)
	}
	if permissions.calls != 0 {
		t.Fatalf(
			"expected no permission lookup without requested fields, got %d",
			permissions.calls,
		)
	}

	fields, _, err = h.resolveFields(userInfo, ResolveOpts{IncludeNote: true})
	if err != nil {
		t.Fatalf("resolveFields returned error: %v", err)
	}
	if !fields.note || fields.labels {
		t.Fatalf("unexpected fields for note request: %+v", fields)
	}

	fields, _, err = h.resolveFields(userInfo, ResolveOpts{IncludeLabels: true})
	if err != nil {
		t.Fatalf("resolveFields returned error: %v", err)
	}
	if fields.note || !fields.labels {
		t.Fatalf("unexpected fields for labels request: %+v", fields)
	}
}

func expectFallbackQuery(mock sqlmock.Sqlmock, userIDs ...int32) {
	rows := sqlmock.NewRows([]string{
		"colleague.id",
		"colleague.job",
		"colleague.job_grade",
		"colleague.firstname",
		"colleague.lastname",
		"colleague.dateofbirth",
		"colleague.phone_number",
		"colleague.profile_picture_file_id",
		"colleague.profile_picture",
		"colleague.email",
	})

	for _, userID := range userIDs {
		rows.AddRow(
			userID,
			"civilian",
			int32(2),
			fmt.Sprintf("Fallback%d", userID),
			"User",
			"1990-01-01",
			nil,
			nil,
			nil,
			nil,
		)
	}

	mock.ExpectQuery(
		`(?s)^SELECT .*FROM fivenet_user AS colleague .*colleague\.deleted_at IS NULL.*LIMIT \?;$`,
	).
		WithArgs(userIDs[0], userIDs[1], int64(len(userIDs))).
		WillReturnRows(rows)
}

func TestListByUserIDFallsBackForMissingColleagueJob(t *testing.T) {
	t.Parallel()

	h, mock := newTestHydrator(
		t,
		&jobscolleagues.Colleague{
			UserId:      1,
			Job:         "police",
			JobGrade:    4,
			Firstname:   "Alice",
			Lastname:    "Active",
			Dateofbirth: "1990-02-03",
			PhoneNumber: func() *string { v := "555-1000"; return &v }(),
			Props: &jobscolleagues.ColleagueProps{
				UserId: 1,
				Job:    "police",
			},
		},
	)
	expectFallbackQuery(mock, 1, 2)

	got, err := h.ListByUserID(
		t.Context(),
		nil,
		nil,
		[]int32{1, 2},
		ResolveOpts{
			Scope: JobScope{Mode: JobScopeExplicit, Job: "police"},
		},
	)
	if err != nil {
		t.Fatalf("ListByUserID returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 colleagues, got %d", len(got))
	}

	first := got[0]
	if first.GetUserId() != 1 {
		t.Fatalf("expected first colleague user_id 1, got %d", first.GetUserId())
	}
	if first.GetProps() == nil {
		t.Fatalf("expected first colleague props to be retained")
	}
	if first.GetJob() != "police" {
		t.Fatalf("expected first colleague job police, got %q", first.GetJob())
	}
	if first.GetJobLabel() != "label:police" {
		t.Fatalf("expected first colleague job label to be enriched, got %q", first.GetJobLabel())
	}
	if first.GetJobGradeLabel() != "grade:4" {
		t.Fatalf(
			"expected first colleague job grade label to be enriched, got %q",
			first.GetJobGradeLabel(),
		)
	}

	second := got[1]
	if second.GetUserId() != 2 {
		t.Fatalf("expected second colleague user_id 2, got %d", second.GetUserId())
	}
	if second.GetProps() != nil {
		t.Fatalf("expected fallback colleague props to stay nil")
	}
	if second.GetJob() != "civilian" {
		t.Fatalf("expected fallback colleague job civilian, got %q", second.GetJob())
	}
	if second.GetJobLabel() != "label:civilian" {
		t.Fatalf(
			"expected fallback colleague job label to be enriched, got %q",
			second.GetJobLabel(),
		)
	}
	if second.GetJobGradeLabel() != "grade:2" {
		t.Fatalf(
			"expected fallback colleague job grade label to be enriched, got %q",
			second.GetJobGradeLabel(),
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations not met: %v", err)
	}
}

func TestHydrateTargetsUsesFallbackForMissingColleagueJob(t *testing.T) {
	t.Parallel()

	h, mock := newTestHydrator(
		t,
		&jobscolleagues.Colleague{
			UserId:      1,
			Job:         "police",
			JobGrade:    4,
			Firstname:   "Alice",
			Lastname:    "Active",
			Dateofbirth: "1990-02-03",
			PhoneNumber: func() *string { v := "555-1000"; return &v }(),
			Props: &jobscolleagues.ColleagueProps{
				UserId: 1,
				Job:    "police",
			},
		},
	)
	expectFallbackQuery(mock, 1, 2)

	targets := []Target{
		{
			UserID: 1,
			Set: func(colleague *jobscolleagues.Colleague) {
				if colleague == nil {
					t.Fatal("expected colleague for user 1")
				}
			},
		},
		{
			UserID: 2,
			Set: func(colleague *jobscolleagues.Colleague) {
				if colleague == nil {
					t.Fatal("expected fallback colleague for user 2")
				}
				if colleague.GetProps() != nil {
					t.Fatal("expected fallback colleague props to stay nil")
				}
				if colleague.GetJob() != "civilian" {
					t.Fatalf("expected fallback colleague job civilian, got %q", colleague.GetJob())
				}
			},
		},
	}

	if err := h.HydrateTargets(
		t.Context(),
		nil,
		nil,
		targets,
		ResolveOpts{
			Scope: JobScope{Mode: JobScopeExplicit, Job: "police"},
		},
	); err != nil {
		t.Fatalf("HydrateTargets returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations not met: %v", err)
	}
}

func TestListByUserIDSkipsDeletedColleagueRows(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	h := &Hydrator{
		db:       db,
		enricher: &testUserAwareEnricher{},
		store: &testColleagueStore{
			Store:      &jobsstore.Store{},
			colleagues: []*jobscolleagues.Colleague{},
			deleted: map[int32]struct{}{
				7: {},
			},
		},
	}

	mock.ExpectQuery(
		`(?s)^SELECT .*FROM fivenet_user AS colleague .*colleague\.deleted_at IS NULL.*LIMIT \?;$`,
	).
		WithArgs(int32(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"colleague.id",
			"colleague.job",
			"colleague.job_grade",
			"colleague.firstname",
			"colleague.lastname",
			"colleague.dateofbirth",
			"colleague.phone_number",
			"colleague.profile_picture_file_id",
			"colleague.profile_picture",
			"colleague.email",
		}))

	got, err := h.ListByUserID(
		t.Context(),
		nil,
		nil,
		[]int32{7},
		ResolveOpts{
			Scope: JobScope{Mode: JobScopeExplicit, Job: "police"},
		},
	)
	if err != nil {
		t.Fatalf("ListByUserID returned error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected deleted colleague to be skipped, got %d results", len(got))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations not met: %v", err)
	}
}
