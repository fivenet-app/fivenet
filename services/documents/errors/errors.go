package errorsdocuments

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery         = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.DocumentsService.ErrFailedQuery"}, nil)
	ErrNotFoundOrNoPerms   = common.I18nErr(codes.NotFound, &common.TranslateItem{Key: "errors.DocumentsService.ErrNotFoundOrNoPerms"}, nil)
	ErrTemplateNoPerms     = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrTemplateNoPerms"}, nil)
	ErrPermissionDenied    = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrPermissionDenied"}, nil)
	ErrClosedDoc           = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrClosedDoc"}, nil)
	ErrDocViewDenied       = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocViewDenied"}, nil)
	ErrDocUpdateDenied     = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocUpdateDenied"}, nil)
	ErrDocDeleteDenied     = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocDeleteDenied"}, nil)
	ErrDocToggleDenied     = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocToggleDenied"}, nil)
	ErrDocAccessEditDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocAccessEditDenied"}, nil)
	ErrDocAccessViewDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocAccessViewDenied"}, nil)
	ErrDocSameOwner        = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocSameOwner"}, nil)
	ErrDocOwnerWrongJob    = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocOwnerWrongJob"}, nil)
	ErrDocOwnerFailed      = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocOwnerFailed"}, nil)
	ErrDocAccessDuplicate  = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocAccessDuplicate"}, nil)

	ErrTemplateFailed            = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrTemplateFailed"}, nil)
	ErrDocRequiredAccessTemplate = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocRequiredAccessTemplate.content"}, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocRequiredAccessTemplate.title"})
	ErrTemplateAccessDuplicate   = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrTemplateAccessDuplicate"}, nil)

	ErrFeedRefsViewDenied  = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrFeedRefsViewDenied"}, nil)
	ErrFeedRelsViewDenied  = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrFeedRelsViewDenied"}, nil)
	ErrFeedRefSelf         = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrFeedRefSelf"}, nil)
	ErrFeedRefAddDenied    = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrFeedRefAddDenied"}, nil)
	ErrFeedRefRemoveDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrFeedRefRemoveDenied"}, nil)
	ErrFeedRelAddDenied    = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrFeedRelAddDenied"}, nil)
	ErrFeedRelRemoveDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrFeedRelRemoveDenied"}, nil)

	ErrCommentViewDenied   = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrCommentViewDenied"}, nil)
	ErrCommentPostDenied   = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrCommentPostDenied"}, nil)
	ErrCommentEditDenied   = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrCommentEditDenied"}, nil)
	ErrCommentDeleteDenied = common.I18nErr(codes.PermissionDenied, &common.TranslateItem{Key: "errors.DocumentsService.ErrCommentDeleteDenied"}, nil)

	ErrDocReqAlreadyCreated   = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocReqAlreadyCreated.content"}, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocReqAlreadyCreated.title"})
	ErrDocReqOwnDoc           = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocReqOwnDoc"}, nil)
	ErrDocReqAlreadyCompleted = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocReqAlreadyCompleted.content"}, &common.TranslateItem{Key: "errors.DocumentsService.ErrDocReqAlreadyCompleted.title"})
)
