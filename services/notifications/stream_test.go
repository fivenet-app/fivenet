package notifications

import (
	"context"
	"testing"
	"time"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	maileremails "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/emails"
	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/notifications"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	natsutils "github.com/fivenet-app/fivenet/v2026/pkg/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	mailerstore "github.com/fivenet-app/fivenet/v2026/stores/mailer"
	notificationsstore "github.com/fivenet-app/fivenet/v2026/stores/notifications"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type accountGroupsChangeSubscriberStub struct {
	accountGroupsChanges chan *pbuserinfo.AccountGroupsChanged
}

func (s *accountGroupsChangeSubscriberStub) SubscribeUserInfoChanges() chan *pbuserinfo.UserInfoChanged {
	return make(chan *pbuserinfo.UserInfoChanged)
}

func (s *accountGroupsChangeSubscriberStub) UnsubscribeUserInfoChanges(chan *pbuserinfo.UserInfoChanged) {
}

func (s *accountGroupsChangeSubscriberStub) SubscribeAccountGroupsChanges() chan *pbuserinfo.AccountGroupsChanged {
	return s.accountGroupsChanges
}

func (s *accountGroupsChangeSubscriberStub) UnsubscribeAccountGroupsChanges(
	chan *pbuserinfo.AccountGroupsChanged,
) {
}

type notificationsStreamServerStub struct {
	ctx  context.Context
	sent chan *pbnotifications.StreamResponse
}

type notificationStoreStub struct {
	notificationsstore.IStore
}

func (s *notificationStoreStub) CountUnread(context.Context, int32) (int64, error) {
	return 0, nil
}

type mailerStoreStub struct {
	mailerstore.IStore
}

func (s *mailerStoreStub) ListUserEmails(
	context.Context,
	qrm.DB,
	*pbuserinfo.UserInfo,
	*database.PaginationRequest,
	bool,
	bool,
) ([]*maileremails.Email, error) {
	return nil, nil
}

func (s *notificationsStreamServerStub) SetHeader(metadata.MD) error  { return nil }
func (s *notificationsStreamServerStub) SendHeader(metadata.MD) error { return nil }
func (s *notificationsStreamServerStub) SetTrailer(metadata.MD)       {}
func (s *notificationsStreamServerStub) Context() context.Context     { return s.ctx }
func (s *notificationsStreamServerStub) SendMsg(any) error            { return nil }
func (s *notificationsStreamServerStub) RecvMsg(any) error            { return nil }

func (s *notificationsStreamServerStub) Recv() (*pbnotifications.StreamRequest, error) {
	<-s.ctx.Done()
	return nil, s.ctx.Err()
}

func (s *notificationsStreamServerStub) Send(response *pbnotifications.StreamResponse) error {
	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	case s.sent <- response:
		return nil
	}
}

func receiveStreamResponse(t *testing.T, sent <-chan *pbnotifications.StreamResponse) *pbnotifications.StreamResponse {
	t.Helper()
	select {
	case response := <-sent:
		return response
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for stream response")
		return nil
	}
}

func TestStreamForwardsCanonicalAccountGroupsChanged(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		AccountId: 42,
		UserId:    7,
		Job:       "police",
		JobGrade:  4,
		Superuser: true,
		OriginalJob: &pbuserinfo.OriginalJob{
			Job:      "ambulance",
			JobGrade: 2,
		},
	}))
	defer cancel()

	natsServer := nats.NewServer(t, nats.ServerOptions{InProcess: true})
	js := natsServer.GetJS()
	_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      notifi.StreamName,
		Subjects:  []string{"notifi.>"},
		Retention: jetstream.InterestPolicy,
	})
	require.NoError(t, err)

	changes := &accountGroupsChangeSubscriberStub{
		accountGroupsChanges: make(chan *pbuserinfo.AccountGroupsChanged, 1),
	}
	server := &Server{
		js:              js,
		store:           &notificationStoreStub{},
		mailerStore:     &mailerStoreStub{},
		userinfoChanges: changes,
		metrics:         getNotificationMetrics(),
	}
	stream := &notificationsStreamServerStub{
		ctx:  ctx,
		sent: make(chan *pbnotifications.StreamResponse, 2),
	}
	done := make(chan error, 1)
	go func() { done <- server.Stream(stream) }()

	initial := receiveStreamResponse(t, stream.sent)
	assert.True(t, initial.GetNotificationState())

	changes.accountGroupsChanges <- &pbuserinfo.AccountGroupsChanged{
		AccountId:        42,
		NewGroups:        &accounts.AccountGroups{Groups: []string{"admin"}},
		CanBeSuperuser:   false,
		CanBeConfigAdmin: true,
	}
	response := receiveStreamResponse(t, stream.sent)
	change := response.GetUserEvent().GetAccountGroupsChanged()
	require.NotNil(t, change)
	assert.Equal(t, int64(42), change.GetAccountId())
	assert.Equal(t, []string{"admin"}, change.GetNewGroups().GetGroups())
	assert.False(t, change.GetCanBeSuperuser())
	assert.True(t, change.GetCanBeConfigAdmin())

	consumer, err := js.Consumer(
		ctx,
		notifi.StreamName,
		natsutils.GenerateConsumerName(42, 7, ""),
	)
	require.NoError(t, err)
	info, err := consumer.Info(ctx)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{
		"notifi.sys",
		"notifi.user.7",
		"notifi.job.ambulance",
		"notifi.job_grade.ambulance.>",
	}, info.Config.FilterSubjects)

	cancel()
	require.NoError(t, <-done)
}

func TestApplyUserInfoChanged(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Job:      "police",
		JobGrade: 1,
	}

	applyUserInfoChanged(current, &pbuserinfo.UserInfoChanged{
		NewJob:      new("ems"),
		NewJobGrade: new(int32(3)),
	})

	assert.Equal(t, "ems", current.GetJob())
	assert.Equal(t, int32(3), current.GetJobGrade())
}

func TestApplyUserInfoChangedIgnoresNil(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Job:      "police",
		JobGrade: 1,
	}

	applyUserInfoChanged(current, nil)

	assert.Equal(t, "police", current.GetJob())
	assert.Equal(t, int32(1), current.GetJobGrade())
}

func TestApplyAccountGroupsChanged(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Groups: &accounts.AccountGroups{Groups: []string{"old"}},
	}

	applyAccountGroupsChanged(current, &pbuserinfo.AccountGroupsChanged{
		NewGroups:      &accounts.AccountGroups{Groups: []string{"supporter", "donator"}},
		CanBeSuperuser: true,
	})

	assert.Equal(t, []string{"supporter", "donator"}, current.GetGroups().GetGroups())
	assert.True(t, current.GetCanBeSuperuser())
	assert.False(t, current.GetCanBeConfigAdmin())
}

func TestApplyAccountGroupsChangedUpdatesConfigAdminWithoutChangingJob(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Job:              "police",
		JobGrade:         3,
		CanBeSuperuser:   false,
		CanBeConfigAdmin: false,
	}
	applyAccountGroupsChanged(current, &pbuserinfo.AccountGroupsChanged{
		CanBeSuperuser:   false,
		CanBeConfigAdmin: true,
	})

	assert.Equal(t, "police", current.GetJob())
	assert.Equal(t, int32(3), current.GetJobGrade())
	assert.False(t, current.GetCanBeSuperuser())
	assert.True(t, current.GetCanBeConfigAdmin())
}

func TestApplyAccountGroupsChangedClearsNil(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Groups: &accounts.AccountGroups{Groups: []string{"old"}},
	}

	applyAccountGroupsChanged(current, &pbuserinfo.AccountGroupsChanged{})

	assert.Nil(t, current.GetGroups())
	assert.False(t, current.GetCanBeSuperuser())
	assert.False(t, current.GetCanBeConfigAdmin())
	assert.False(t, current.GetSuperuser())
}

func TestApplyAccountGroupsChangedRestoresOriginalJobWhenRevokingSuperuser(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Job:       "ems-super",
		JobGrade:  7,
		Superuser: true,
		OriginalJob: &pbuserinfo.OriginalJob{
			Job:      "ems",
			JobGrade: 2,
		},
	}

	applyAccountGroupsChanged(current, &pbuserinfo.AccountGroupsChanged{
		CanBeSuperuser: false,
	})

	assert.False(t, current.GetCanBeSuperuser())
	assert.False(t, current.GetCanBeConfigAdmin())
	assert.False(t, current.GetSuperuser())
	assert.Equal(t, "ems", current.GetJob())
	assert.Equal(t, int32(2), current.GetJobGrade())
}

func TestBuildSubjectsAccountOnly(t *testing.T) {
	t.Parallel()

	s := &Server{}

	baseSubjects, additionalSubjects, err := s.buildSubjects(t.Context(), &pbuserinfo.UserInfo{
		AccountId: 42,
	})
	require.NoError(t, err)

	assert.Equal(t, []string{
		"notifi.sys",
	}, baseSubjects)
	assert.Empty(t, additionalSubjects)
}

func TestValidPreferenceScope(t *testing.T) {
	t.Parallel()

	assert.True(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{}))
	assert.True(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
	}))
	assert.True(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_JOBS,
		Kind:     resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_LEADERSHIP_ADDED,
	}))
	assert.True(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_JOBS,
		Kind:     resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_MEMBER_ADDED,
	}))
	assert.True(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Kind:     resourcesnotifications.NotificationKind_NOTIFICATION_KIND_DOCUMENT_APPROVAL_ASSIGNED,
	}))
	assert.True(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Kind:     resourcesnotifications.NotificationKind_NOTIFICATION_KIND_DOCUMENT_REQUEST_DECIDED,
	}))
	assert.False(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Kind:     resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_LEADERSHIP_ADDED,
	}))
	assert.False(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory(99),
	}))
	assert.False(t, validPreferenceScope(nil))
}

func TestHasDeliveryPreferenceOverride(t *testing.T) {
	t.Parallel()

	assert.False(t, hasDeliveryPreferenceOverride(nil))
	assert.False(t, hasDeliveryPreferenceOverride(&resourcesnotifications.NotificationPreference{}))
	assert.True(
		t,
		hasDeliveryPreferenceOverride(
			&resourcesnotifications.NotificationPreference{InboxEnabled: new(true)},
		),
	)
}
