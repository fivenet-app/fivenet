package qualificationsstore

import (
	"context"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CreateQualification(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
) (int64, error) {
	stmt := tQuali.
		INSERT(
			tQuali.Job,
			tQuali.Closed,
			tQuali.Draft,
			tQuali.Public,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.ContentType,
			tQuali.Content,
			tQuali.CreatorID,
			tQuali.CreatorJob,
		).
		VALUES(
			userInfo.GetJob(),
			false,
			true,
			false,
			"",
			"",
			"",
			int32(content.ContentType_CONTENT_TYPE_HTML),
			"",
			userInfo.GetUserId(),
			userInfo.GetJob(),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *Store) UpdateQualification(
	ctx context.Context,
	tx qrm.DB,
	quali *resqualifications.Qualification,
) error {
	if quali.Description != nil {
		*quali.Description = strings.TrimSuffix(quali.GetDescription(), "<br>")
	}

	stmt := tQuali.
		UPDATE(
			tQuali.Weight,
			tQuali.Closed,
			tQuali.Draft,
			tQuali.Public,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.ContentType,
			tQuali.Content,
			tQuali.DiscordSyncEnabled,
			tQuali.DiscordSettings,
			tQuali.ExamMode,
			tQuali.ExamSettings,
			tQuali.LabelSyncEnabled,
			tQuali.LabelSyncFormat,
		).
		SET(
			quali.GetWeight(),
			quali.GetClosed(),
			quali.GetDraft(),
			quali.GetPublic(),
			quali.GetAbbreviation(),
			quali.GetTitle(),
			quali.Description,
			int32(quali.GetContent().GetContentType()),
			quali.GetContent(),
			quali.GetDiscordSyncEnabled(),
			quali.GetDiscordSettings(),
			quali.GetExamMode(),
			quali.GetExamSettings(),
			quali.GetLabelSyncEnabled(),
			quali.LabelSyncFormat,
		).
		WHERE(tQuali.ID.EQ(mysql.Int64(quali.GetId()))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) DeleteQualification(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	deletedAt *timestamp.Timestamp,
) error {
	stmt := tQuali.
		UPDATE(tQuali.DeletedAt).
		SET(tQuali.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)))
	if deletedAt != nil {
		stmt = tQuali.
			UPDATE(tQuali.DeletedAt).
			SET(tQuali.DeletedAt.SET(mysql.TimestampT(deletedAt.AsTime())))
	}

	_, err := stmt.
		WHERE(tQuali.ID.EQ(mysql.Int64(qualificationId))).
		LIMIT(1).
		ExecContext(ctx, tx)
	return err
}
