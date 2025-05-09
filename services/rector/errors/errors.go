package errorsrector

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery       = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.RectorService.ErrFailedQuery"}, nil)
	ErrInvalidRequest    = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.RectorService.ErrInvalidRequest"}, nil)
	ErrNoPermission      = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.RectorService.ErrNoPermission"}, nil)
	ErrRoleAlreadyExists = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.RectorService.ErrRoleAlreadyExists"}, nil)
	ErrOwnRoleDeletion   = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.RectorService.ErrOwnRoleDeletion"}, nil)
	ErrInvalidAttrs      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.RectorService.ErrInvalidAttrs"}, nil)
	ErrInvalidPerms      = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.RectorService.ErrInvalidPerms"}, nil)
)
