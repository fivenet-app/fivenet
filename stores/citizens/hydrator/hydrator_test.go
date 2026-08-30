package citizenshydrator

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/permsstub"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
)

type testPermissions struct {
	permsstub.Permissions

	fields []string
}

func (p *testPermissions) AttrStringList(
	_ *userinfo.UserInfo,
	_ perms.AttrRef[perms.StringListAttr],
) (*permissionsattributes.StringList, error) {
	return &permissionsattributes.StringList{Strings: p.fields}, nil
}

type testUserAwareEnricher struct {
	mstlystcdata.DummyEnricher

	calls int
}

func (e *testUserAwareEnricher) EnrichJobInfoSafe(
	userInfo *userinfo.UserInfo,
	usrs ...common.IJobInfo,
) {
	for _, usr := range usrs {
		e.calls++
		usr.SetJob("unemployed")
		usr.SetJobGrade(0)
		usr.SetJobLabel("N/A")
		usr.SetJobGradeLabel("N/A")
	}
}

func (e *testUserAwareEnricher) EnrichJobInfoSafeFunc(
	userInfo *userinfo.UserInfo,
) func(usr common.IJobInfo) {
	return func(usr common.IJobInfo) {
		e.calls++
		usr.SetJob("unemployed")
		usr.SetJobGrade(0)
		usr.SetJobLabel("N/A")
		usr.SetJobGradeLabel("N/A")
	}
}

func TestEnrichShortsSafeRedactsJobInfo(t *testing.T) {
	t.Parallel()

	enricher := &testUserAwareEnricher{}
	h := &Hydrator{
		enricher: enricher,
	}

	user := &usershort.UserShort{
		UserId:   42,
		Job:      "police",
		JobGrade: 4,
	}

	h.enrichShortsSafe(nil, user, nil)

	if enricher.calls != 1 {
		t.Fatalf("expected 1 enrichment call, got %d", enricher.calls)
	}
	if got := user.GetJob(); got != "unemployed" {
		t.Fatalf("expected job to be redacted, got %q", got)
	}
	if got := user.GetJobGrade(); got != 0 {
		t.Fatalf("expected job grade to be redacted, got %d", got)
	}
	if got := user.GetJobLabel(); got != "N/A" {
		t.Fatalf("expected job label to be redacted, got %q", got)
	}
	if got := user.GetJobGradeLabel(); got != "N/A" {
		t.Fatalf("expected job grade label to be redacted, got %q", got)
	}
}

func TestResolveFieldsUsesCitizenPermissions(t *testing.T) {
	t.Parallel()

	h := &Hydrator{perms: &testPermissions{fields: []string{
		string(permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber),
		string(permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsBloodType),
		string(permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsMugshot),
	}}}

	fields, err := h.resolveFields(&userinfo.UserInfo{})
	if err != nil {
		t.Fatalf("resolveFields returned error: %v", err)
	}
	if !fields.phoneNumber || !fields.bloodType || !fields.mugshot {
		t.Fatalf("expected permitted fields to be resolved: %#v", fields)
	}
	if fields.email || fields.openFines || fields.licenses || fields.labels {
		t.Fatalf("unexpected fields resolved: %#v", fields)
	}
}

func TestListBasicByUserIDOnlySelectsPermittedPhoneNumber(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	h := &Hydrator{
		db: db,
		perms: &testPermissions{
			fields: []string{
				string(permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber),
			},
		},
		enricher: mstlystcdata.NewDummyUserAwareEnricher(),
	}
	mock.ExpectQuery(`(?s)^SELECT .*user_short\.phone_number.*FROM fivenet_user AS user_short .*user_short\.deleted_at IS NULL.*LIMIT \?;$`).
		WithArgs(int32(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_short.id",
			"user_short.firstname",
			"user_short.lastname",
			"user_short.job",
			"user_short.job_grade",
			"user_short.phone_number",
		}).AddRow(42, "Ada", "Lovelace", "police", 4, "555-0100"))

	users, err := h.ListBasicByUserID(t.Context(), nil, &userinfo.UserInfo{}, []int32{42})
	if err != nil {
		t.Fatalf("ListBasicByUserID returned error: %v", err)
	}
	if len(users) != 1 || users[0].GetPhoneNumber() != "555-0100" {
		t.Fatalf("expected permission-gated phone number, got %#v", users)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations were not met: %v", err)
	}
}
