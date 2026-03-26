package errorsdocuments

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFailedQuery"},
		nil,
	)
	ErrNotFoundOrNoPerms = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrNotFoundOrNoPerms"},
		nil,
	)
	ErrTemplateNoPerms = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrTemplateNoPerms"},
		nil,
	)
	ErrPermissionDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrPermissionDenied"},
		nil,
	)
	ErrClosedDoc = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrClosedDoc.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrClosedDoc.title"},
	)
	ErrDocViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocViewDenied"},
		nil,
	)
	ErrDocUpdateDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocUpdateDenied"},
		nil,
	)
	ErrDocDeleteDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocDeleteDenied"},
		nil,
	)
	ErrDocToggleDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocToggleDenied"},
		nil,
	)
	ErrDocAccessEditDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessEditDenied"},
		nil,
	)
	ErrDocAccessViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessViewDenied"},
		nil,
	)
	ErrDocSameOwner = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocSameOwner"},
		nil,
	)
	ErrDocOwnerWrongJob = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocOwnerWrongJob"},
		nil,
	)
	ErrDocOwnerFailed = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocOwnerFailed"},
		nil,
	)
	ErrDocAccessDuplicate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessDuplicate"},
		nil,
	)
	ErrDocAccessInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessInvalid"},
		nil,
	)

	ErrTemplateFailed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrTemplateFailed"},
		nil,
	)
	ErrDocRequiredAccessTemplate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.documents.DocumentsService.ErrDocRequiredAccessTemplate.content",
		},
		&common.I18NItem{
			Key: "errors.documents.DocumentsService.ErrDocRequiredAccessTemplate.title",
		},
	)
	ErrTemplateAccessDuplicate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrTemplateAccessDuplicate"},
		nil,
	)

	ErrFeedRefsViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRefsViewDenied"},
		nil,
	)
	ErrFeedRelsViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRelsViewDenied"},
		nil,
	)
	ErrFeedRefSelf = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRefSelf"},
		nil,
	)
	ErrFeedRefAddDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRefAddDenied"},
		nil,
	)
	ErrFeedRefRemoveDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRefRemoveDenied"},
		nil,
	)
	ErrFeedRelAddDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRelAddDenied"},
		nil,
	)
	ErrFeedRelRemoveDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRelRemoveDenied"},
		nil,
	)

	ErrCommentViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrCommentViewDenied"},
		nil,
	)
	ErrCommentPostDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrCommentPostDenied"},
		nil,
	)
	ErrCommentEditDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrCommentEditDenied"},
		nil,
	)
	ErrCommentDeleteDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrCommentDeleteDenied"},
		nil,
	)

	ErrDocReqAlreadyCreated = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocReqAlreadyCreated.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocReqAlreadyCreated.title"},
	)
	ErrDocReqOwnDoc = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocReqOwnDoc"},
		nil,
	)
	ErrDocReqAlreadyCompleted = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.documents.DocumentsService.ErrDocReqAlreadyCompleted.content",
		},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocReqAlreadyCompleted.title"},
	)

	ErrApprovalTaskAlreadyHandled = common.NewI18nErr(
		codes.FailedPrecondition,
		&common.I18NItem{
			Key: "errors.documents.ApprovalService.ErrApprovalTaskAlreadyHandled.content",
		},
		&common.I18NItem{
			Key: "errors.documents.ApprovalService.ErrApprovalTaskAlreadyHandled.title",
		},
	)
	ErrApprovalSignatureRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key: "errors.documents.ApprovalService.ErrApprovalSignatureRequired.content",
		},
		&common.I18NItem{
			Key: "errors.documents.ApprovalService.ErrApprovalSignatureRequired.title",
		},
	)
	ErrApprovalDocIsDraft = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.ApprovalService.ErrApprovalDocIsDraft.content"},
		&common.I18NItem{Key: "errors.documents.ApprovalService.ErrApprovalDocIsDraft.title"},
	)

	ErrStampLimitReached = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key:        "errors.documents.DocumentsService.ErrStampLimitReached.content",
			Parameters: map[string]string{"max": "5"},
		},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrStampLimitReached.title"},
	)
)
