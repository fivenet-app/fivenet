package qualifications

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) getQualification(
	ctx context.Context,
	qualificationId int64,
	condition mysql.BoolExpression,
	userInfo *userinfo.UserInfo,
	selectContent bool,
) (*qualifications.Qualification, error) {
	quali, err := s.store.GetQualification(
		ctx,
		qualificationId,
		condition,
		userInfo,
		selectContent,
		false,
	)
	if err != nil {
		return nil, err
	}

	return quali, nil
}

func (s *Server) checkRequirementsMetForQualification(
	ctx context.Context,
	qualificationId int64,
	userId int32,
) (bool, error) {
	return s.store.CheckRequirementsMetForQualification(ctx, qualificationId, userId)
}
