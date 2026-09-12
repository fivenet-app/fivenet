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
	// BaseSubject is the base subject for user, job, object and system notification events.
	BaseSubject = "userinfo"

	UserInfoStreamName = "USERINFO"
	UserInfoSubject    = "userinfo.*.changes"
)

func registerUserInfoStream(ctx context.Context, js *events.JSWrapper) error {
	userinfoStreamCfg := jetstream.StreamConfig{
		Name:        UserInfoStreamName,
		Description: "Stream for userinfo changes",
		Subjects:    []string{UserInfoSubject},
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

// BuildAccountGroupsChangedEvent constructs the live account group change payload.
func BuildAccountGroupsChangedEvent(
	accountID int64,
	changedAt *pbtimestamp.Timestamp,
	groups *accounts.AccountGroups,
	canBeSuperuser bool,
	canBeConfigAdmin bool,
) *pbuserinfo.AccountGroupsChanged {
	if changedAt == nil {
		changedAt = pbtimestamp.Now()
	}

	event := &pbuserinfo.AccountGroupsChanged{
		AccountId:        accountID,
		ChangedAt:        changedAt,
		CanBeSuperuser:   canBeSuperuser,
		CanBeConfigAdmin: canBeConfigAdmin,
	}
	if groups != nil {
		event.NewGroups = &accounts.AccountGroups{
			Groups: slices.Clone(groups.GetGroups()),
		}
	}

	return event
}
