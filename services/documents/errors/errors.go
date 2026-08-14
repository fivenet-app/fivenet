package errorsdocuments

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFailedQuery.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFailedQuery.title"},
	)
	ErrNotFoundOrNoPerms = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrNotFoundOrNoPerms.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrNotFoundOrNoPerms.title"},
	)
	ErrTemplateNoPerms = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrTemplateNoPerms.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrTemplateNoPerms.title"},
	)
	ErrPermissionDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrPermissionDenied.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrPermissionDenied.title"},
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
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessEditDenied.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessEditDenied.title"},
	)
	ErrDocAccessViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessViewDenied.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessViewDenied.title"},
	)
	ErrDocSameOwner = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocSameOwner.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocSameOwner.title"},
	)
	ErrDocOwnerWrongJob = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocOwnerWrongJob.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocOwnerWrongJob.title"},
	)
	ErrDocOwnerFailed = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocOwnerFailed.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocOwnerFailed.title"},
	)
	ErrDocAccessDuplicate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessDuplicate.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessDuplicate.title"},
	)
	ErrDocAccessInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessInvalid.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocAccessInvalid.title"},
	)

	ErrTemplateFailed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrTemplateFailed.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrTemplateFailed.title"},
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
		&common.I18NItem{
			Key: "errors.documents.DocumentsService.ErrTemplateAccessDuplicate.content",
		},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrTemplateAccessDuplicate.title"},
	)

	ErrFeedRefsViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRefsViewDenied.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRefsViewDenied.title"},
	)
	ErrFeedRelsViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRelsViewDenied.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRelsViewDenied.title"},
	)
	ErrFeedRefSelf = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRefSelf.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrFeedRefSelf.title"},
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
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrCommentPostDenied.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrCommentPostDenied.title"},
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
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocReqOwnDoc.content"},
		&common.I18NItem{Key: "errors.documents.DocumentsService.ErrDocReqOwnDoc.title"},
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
			Key:        "errors.documents.StampsService.ErrStampLimitReached.content",
			Parameters: map[string]string{"max": "5"},
		},
		&common.I18NItem{Key: "errors.documents.StampsService.ErrStampLimitReached.title"},
	)
	ErrNoStatsCategories = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.documents.StampsService.ErrStampLimitReached.content"},
		&common.I18NItem{Key: "errors.documents.StampsService.ErrStampLimitReached.title"},
	)
)
