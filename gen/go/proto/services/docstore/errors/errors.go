package errorsdocstore

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedQuery         = status.Error(codes.Internal, "errors.DocStoreService.ErrFailedQuery")
	ErrNotFoundOrNoPerms   = status.Error(codes.NotFound, "errors.DocStoreService.ErrNotFoundOrNoPerms")
	ErrTemplateNoPerms     = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrTemplateNoPerms")
	ErrPermissionDenied    = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrPermissionDenied")
	ErrClosedDoc           = status.Error(codes.Canceled, "errors.DocStoreService.ErrClosedDoc")
	ErrDocViewDenied       = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrDocViewDenied")
	ErrDocUpdateDenied     = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrDocUpdateDenied")
	ErrDocDeleteDenied     = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrDocDeleteDenied")
	ErrDocToggleDenied     = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrDocToggleDenied")
	ErrDocAccessEditDenied = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrDocAccessEditDenied")
	ErrDocAccessViewDenied = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrDocAccessViewDenied")

	ErrTemplateFailed = status.Error(codes.InvalidArgument, "errors.DocStoreService.ErrTemplateFailed")

	ErrFeedRefsViewDenied  = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrFeedRefsViewDenied")
	ErrFeedRelsViewDenied  = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrFeedRelsViewDenied")
	ErrFeedRefSelf         = status.Error(codes.InvalidArgument, "errors.DocStoreService.ErrFeedRefSelf")
	ErrFeedRefAddDenied    = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrFeedRefAddDenied")
	ErrFeedRefRemoveDenied = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrFeedRefRemoveDenied")
	ErrFeedRelAddDenied    = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrFeedRelAddDenied")
	ErrFeedRelRemoveDenied = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrFeedRelRemoveDenied")

	ErrCommentViewDenied   = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrCommentViewDenied")
	ErrCommentPostDenied   = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrCommentPostDenied")
	ErrCommentEditDenied   = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrCommentEditDenied")
	ErrCommentDeleteDenied = status.Error(codes.PermissionDenied, "errors.DocStoreService.ErrCommentDeleteDenied")

	ErrDocReqAlreadyCreated = status.Error(codes.InvalidArgument, "errors.DocStoreService.ErrDocReqAlreadyCreated.title;errors.DocStoreService.ErrDocReqAlreadyCreated.content")
)
