package documentsstore

import (
	"context"
	"errors"

	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentscomment "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/comment"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CountComments(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	includeDeleted bool,
) (int32, error) {
	tComments := table.FivenetDocumentsComments.AS("comments")
	condition := tComments.DocumentID.EQ(mysql.Int64(documentID))
	if !includeDeleted {
		condition = mysql.AND(condition, tComments.DeletedAt.IS_NULL())
	}

	stmt := tComments.
		SELECT(mysql.COUNT(tComments.ID).AS("comment_count")).
		FROM(tComments).
		WHERE(condition)

	var result struct {
		CommentCount int32 `alias:"comment_count"`
	}
	if err := stmt.QueryContext(ctx, tx, &result); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return result.CommentCount, nil
}

func (s *Store) UpdateCommentsCount(ctx context.Context, tx qrm.DB, documentID int64) error {
	commentCount, err := s.CountComments(ctx, tx, documentID, false)
	if err != nil {
		return err
	}

	tMeta := table.FivenetDocumentsMeta
	updateStmt := tMeta.
		INSERT(
			tMeta.DocumentID,
			tMeta.CommentCount,
		).
		VALUES(
			documentID,
			commentCount,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tMeta.CommentCount.SET(mysql.Int32(commentCount)),
		)

	_, err = updateStmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) ListComments(
	ctx context.Context,
	documentID int64,
	userInfo *userinfo.UserInfo,
	offset int64,
	limit int64,
) ([]*documentscomment.Comment, error) {
	tComments := table.FivenetDocumentsComments.AS("comment")
	tCreator := table.FivenetUser.AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")
	tUserProps := table.FivenetUserProps.AS("user_props")

	condition := tComments.DocumentID.EQ(mysql.Int64(documentID))
	if userInfo == nil || !userInfo.GetJobAdmin() {
		condition = mysql.AND(condition, tComments.DeletedAt.IS_NULL())
	}

	stmt := tComments.
		SELECT(
			tComments.ID,
			tComments.DocumentID,
			tComments.CreatedAt,
			tComments.UpdatedAt,
			tComments.Content,
			tComments.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tUserProps.AvatarFileID.AS("creator.profile_picture_file_id"),
			tAvatar.FilePath.AS("creator.profile_picture"),
		).
		FROM(
			tComments.
				LEFT_JOIN(tCreator,
					tComments.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(condition).
		OFFSET(offset).
		ORDER_BY(tComments.CreatedAt.DESC()).
		LIMIT(limit)

	var dest []*documentscomment.Comment
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) GetComment(
	ctx context.Context,
	id int64,
	userInfo *userinfo.UserInfo,
) (*documentscomment.Comment, error) {
	tComments := table.FivenetDocumentsComments.AS("comment")
	tCreator := table.FivenetUser.AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")
	tUserProps := table.FivenetUserProps.AS("user_props")

	stmt := tComments.
		SELECT(
			tComments.ID,
			tComments.CreatedAt,
			tComments.UpdatedAt,
			tComments.DeletedAt,
			tComments.DocumentID,
			tComments.Content,
			tComments.CreatorID,
			tComments.CreatorJob,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tUserProps.AvatarFileID.AS("creator.profile_picture_file_id"),
			tAvatar.FilePath.AS("creator.profile_picture"),
		).
		FROM(
			tComments.
				LEFT_JOIN(tCreator,
					tComments.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(tComments.ID.EQ(mysql.Int64(id))).
		LIMIT(1)

	comment := &documentscomment.Comment{}
	if err := stmt.QueryContext(ctx, s.db, comment); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if comment.GetId() <= 0 {
		return nil, nil
	}

	if comment.GetCreator() != nil && userInfo != nil {
		// leave enrichment to the service
		_ = userInfo
	}

	return comment, nil
}

func (s *Store) CreateComment(
	ctx context.Context,
	tx qrm.DB,
	comment *documentscomment.Comment,
	userInfo *userinfo.UserInfo,
) (int64, error) {
	tComments := table.FivenetDocumentsComments
	stmt := tComments.
		INSERT(
			tComments.DocumentID,
			tComments.Content,
			tComments.CreatorID,
			tComments.CreatorJob,
		).
		VALUES(
			comment.GetDocumentId(),
			comment.GetContent(),
			userInfo.GetUserId(),
			userInfo.GetJob(),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   comment.GetDocumentId(),
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_ADDED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Store) UpdateComment(
	ctx context.Context,
	tx qrm.DB,
	comment *documentscomment.Comment,
	userInfo *userinfo.UserInfo,
) error {
	tComments := table.FivenetDocumentsComments
	stmt := tComments.
		UPDATE(
			tComments.Content,
		).
		SET(
			comment.GetContent(),
		).
		WHERE(mysql.AND(
			tComments.ID.EQ(mysql.Int64(comment.GetId())),
			tComments.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   comment.GetDocumentId(),
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_UPDATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteComment(
	ctx context.Context,
	tx qrm.DB,
	comment *documentscomment.Comment,
	userInfo *userinfo.UserInfo,
	deletedAt *timestamp.Timestamp,
	activityType documentsactivity.DocActivityType,
) error {
	tComments := table.FivenetDocumentsComments
	stmt := tComments.
		UPDATE(
			tComments.DeletedAt,
		).
		SET(
			tComments.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAt)),
		).
		WHERE(tComments.ID.EQ(mysql.Int64(comment.GetId()))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   comment.GetDocumentId(),
		ActivityType: activityType,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return err
	}

	return nil
}
