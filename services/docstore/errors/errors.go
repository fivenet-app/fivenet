package errorsdocstore

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery         = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.DocStoreService.ErrFailedQuery"}, nil)
	ErrNotFoundOrNoPerms   = common.I18nErr(codes.NotFound, &common.TranslateItem{Key: "errors.DocStoreService.ErrNotFoundOrNoPerms"}, nil)
	ErrTemplateNoPerms     = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrTemplateNoPerms"}, nil)
	ErrPermissionDenied    = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrPermissionDenied"}, nil)
	ErrClosedDoc           = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrClosedDoc"}, nil)
	ErrDocViewDenied       = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocViewDenied"}, nil)
	ErrDocUpdateDenied     = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocUpdateDenied"}, nil)
	ErrDocDeleteDenied     = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocDeleteDenied"}, nil)
	ErrDocToggleDenied     = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocToggleDenied"}, nil)
	ErrDocAccessEditDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocAccessEditDenied"}, nil)
	ErrDocAccessViewDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocAccessViewDenied"}, nil)
	ErrDocSameOwner        = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocSameOwner"}, nil)
	ErrDocOwnerWrongJob    = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocOwnerWrongJob"}, nil)
	ErrDocOwnerFailed      = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocOwnerFailed"}, nil)
	ErrDocAccessDuplicate  = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocAccessDuplicate"}, nil)

	ErrTemplateFailed            = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrTemplateFailed"}, nil)
	ErrDocRequiredAccessTemplate = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocRequiredAccessTemplate.content"}, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocRequiredAccessTemplate.title"})
	ErrTemplateAccessDuplicate   = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrTemplateAccessDuplicate"}, nil)

	ErrFeedRefsViewDenied  = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrFeedRefsViewDenied"}, nil)
	ErrFeedRelsViewDenied  = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrFeedRelsViewDenied"}, nil)
	ErrFeedRefSelf         = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrFeedRefSelf"}, nil)
	ErrFeedRefAddDenied    = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrFeedRefAddDenied"}, nil)
	ErrFeedRefRemoveDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrFeedRefRemoveDenied"}, nil)
	ErrFeedRelAddDenied    = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrFeedRelAddDenied"}, nil)
	ErrFeedRelRemoveDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrFeedRelRemoveDenied"}, nil)

	ErrCommentViewDenied   = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrCommentViewDenied"}, nil)
	ErrCommentPostDenied   = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrCommentPostDenied"}, nil)
	ErrCommentEditDenied   = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrCommentEditDenied"}, nil)
	ErrCommentDeleteDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocStoreService.ErrCommentDeleteDenied"}, nil)

	ErrDocReqAlreadyCreated   = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocReqAlreadyCreated.content"}, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocReqAlreadyCreated.title"})
	ErrDocReqOwnDoc           = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocReqOwnDoc"}, nil)
	ErrDocReqAlreadyCompleted = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocReqAlreadyCompleted.content"}, &common.TranslateItem{Key: "errors.DocStoreService.ErrDocReqAlreadyCompleted.title"})
)
