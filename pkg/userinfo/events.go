package userinfo

import (
	"context"
	"fmt"
	"slices"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	pbtimestamp "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	KVBucketName = "userinfo_poll_ttl"

	// BaseSubject is the base subject for user, job, object and system notification events.
	BaseSubject = "userinfo"

	PollStreamName = "POLL_REQUESTS"
	PollSubject    = "userinfo.poll.request"

	UserInfoStreamName = "USERINFO"
	UserInfoSubject    = "userinfo.*.changes"
	UserGroupsSubject  = "userinfo.*.groups"
)

func registerStreams(ctx context.Context, js *events.JSWrapper) error {
	// Stream for incoming poll requests
	pollStreamCfg := jetstream.StreamConfig{
		Name:        PollStreamName,
		Description: "Stream for userinfo poll requests",
		Subjects:    []string{PollSubject},
		Retention:   jetstream.InterestPolicy,
	}
	if _, err := js.CreateOrUpdateStream(ctx, pollStreamCfg); err != nil {
		return fmt.Errorf("failed to create/update stream %s. %w", PollStreamName, err)
	}

	// Stream for userinfo diffs
	userinfoStreamCfg := jetstream.StreamConfig{
		Name:        UserInfoStreamName,
		Description: "Stream for userinfo changes",
		Subjects:    []string{UserInfoSubject, UserGroupsSubject},
		Retention:   jetstream.InterestPolicy,
	}
	if _, err := js.CreateOrUpdateStream(ctx, userinfoStreamCfg); err != nil {
		return fmt.Errorf("failed to create/update stream %s. %w", UserInfoStreamName, err)
	}

	return nil
}

// BuildUserInfoChangedEvent constructs the live user info change payload and enriches labels.
func BuildUserInfoChangedEvent(
	accountID int64,
	userID int32,
	changedAt *pbtimestamp.Timestamp,
	job string,
	jobGrade int32,
	enricher mstlystcdata.IEnricher,
) *pbuserinfo.UserInfoChanged {
	if changedAt == nil {
		changedAt = pbtimestamp.Now()
	}

	event := &pbuserinfo.UserInfoChanged{
		AccountId: accountID,
		UserId:    userID,
		ChangedAt: changedAt,
	}
	event.SetJob(job)
	event.SetJobGrade(jobGrade)

	if enricher != nil {
		enricher.EnrichJobInfo(event)
	}

	return event
}

// BuildUserGroupsChangedEvent constructs the live account group change payload.
func BuildUserGroupsChangedEvent(
	accountID int64,
	changedAt *pbtimestamp.Timestamp,
	groups *accounts.AccountGroups,
	canBeSuperuser bool,
) *pbuserinfo.UserGroupsChanged {
	if changedAt == nil {
		changedAt = pbtimestamp.Now()
	}

	event := &pbuserinfo.UserGroupsChanged{
		AccountId:      accountID,
		ChangedAt:      changedAt,
		CanBeSuperuser: canBeSuperuser,
	}
	if groups != nil {
		event.NewGroups = &accounts.AccountGroups{
			Groups: slices.Clone(groups.GetGroups()),
		}
	}

	return event
}
