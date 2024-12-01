package errorsjobs

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery        = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.JobsService.ErrFailedQuery"}, nil)
	ErrPropsAbsenceDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.JobsService.ErrPropsAbsenceDenied"}, nil)
	ErrPropsNoteDenied    = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.JobsService.ErrPropsNoteDenied"}, nil)
	ErrPropsLabelsDenied  = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.JobsService.ErrPropsLabelsDenied"}, nil)
	ErrReasonRequired     = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.JobsService.ErrReasonRequired"}, nil)
	ErrNotFoundOrNoPerms  = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.JobsService.ErrNotFoundOrNoPerms"}, nil)
)
