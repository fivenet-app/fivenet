package errorsdocuments

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"google.golang.org/grpc/codes"
)

var (
	ErrFailedQuery = common.NewI18nErr(
		codes.Internal,
		&common.I18NItem{Key: "errors.DocumentsService.ErrFailedQuery"},
		nil,
	)
	ErrNotFoundOrNoPerms = common.NewI18nErr(
		codes.NotFound,
		&common.I18NItem{Key: "errors.DocumentsService.ErrNotFoundOrNoPerms"},
		nil,
	)
	ErrTemplateNoPerms = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrTemplateNoPerms"},
		nil,
	)
	ErrPermissionDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrPermissionDenied"},
		nil,
	)
	ErrClosedDoc = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrClosedDoc"},
		nil,
	)
	ErrDocViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocViewDenied"},
		nil,
	)
	ErrDocUpdateDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocUpdateDenied"},
		nil,
	)
	ErrDocDeleteDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocDeleteDenied"},
		nil,
	)
	ErrDocToggleDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocToggleDenied"},
		nil,
	)
	ErrDocAccessEditDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocAccessEditDenied"},
		nil,
	)
	ErrDocAccessViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocAccessViewDenied"},
		nil,
	)
	ErrDocSameOwner = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocSameOwner"},
		nil,
	)
	ErrDocOwnerWrongJob = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocOwnerWrongJob"},
		nil,
	)
	ErrDocOwnerFailed = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocOwnerFailed"},
		nil,
	)
	ErrDocAccessDuplicate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocAccessDuplicate"},
		nil,
	)
	ErrDocAccessInvalid = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocAccessInvalid"},
		nil,
	)

	ErrTemplateFailed = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrTemplateFailed"},
		nil,
	)
	ErrDocRequiredAccessTemplate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocRequiredAccessTemplate.content"},
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocRequiredAccessTemplate.title"},
	)
	ErrTemplateAccessDuplicate = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrTemplateAccessDuplicate"},
		nil,
	)

	ErrFeedRefsViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrFeedRefsViewDenied"},
		nil,
	)
	ErrFeedRelsViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrFeedRelsViewDenied"},
		nil,
	)
	ErrFeedRefSelf = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrFeedRefSelf"},
		nil,
	)
	ErrFeedRefAddDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrFeedRefAddDenied"},
		nil,
	)
	ErrFeedRefRemoveDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrFeedRefRemoveDenied"},
		nil,
	)
	ErrFeedRelAddDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrFeedRelAddDenied"},
		nil,
	)
	ErrFeedRelRemoveDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrFeedRelRemoveDenied"},
		nil,
	)

	ErrCommentViewDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrCommentViewDenied"},
		nil,
	)
	ErrCommentPostDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrCommentPostDenied"},
		nil,
	)
	ErrCommentEditDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrCommentEditDenied"},
		nil,
	)
	ErrCommentDeleteDenied = common.NewI18nErr(
		codes.PermissionDenied,
		&common.I18NItem{Key: "errors.DocumentsService.ErrCommentDeleteDenied"},
		nil,
	)

	ErrDocReqAlreadyCreated = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocReqAlreadyCreated.content"},
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocReqAlreadyCreated.title"},
	)
	ErrDocReqOwnDoc = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocReqOwnDoc"},
		nil,
	)
	ErrDocReqAlreadyCompleted = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocReqAlreadyCompleted.content"},
		&common.I18NItem{Key: "errors.DocumentsService.ErrDocReqAlreadyCompleted.title"},
	)

	ErrApprovalTaskAlreadyHandled = common.NewI18nErr(
		codes.FailedPrecondition,
		&common.I18NItem{Key: "errors.DocumentsService.ErrApprovalTaskAlreadyHandled.content"},
		&common.I18NItem{Key: "errors.DocumentsService.ErrApprovalTaskAlreadyHandled.title"},
	)
	ErrApprovalSignatureRequired = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrApprovalSignatureRequired.content"},
		&common.I18NItem{Key: "errors.DocumentsService.ErrApprovalSignatureRequired.title"},
	)
	ErrApprovalDocIsDraft = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrApprovalDocIsDraft.content"},
		&common.I18NItem{Key: "errors.DocumentsService.ErrApprovalDocIsDraft.title"},
	)
	ErrApprovalCreatorCannotDecide = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{Key: "errors.DocumentsService.ErrApprovalCreatorCannotDecide.content"},
		&common.I18NItem{Key: "errors.DocumentsService.ErrApprovalCreatorCannotDecide.title"},
	)

	ErrStampLimitReached = common.NewI18nErr(
		codes.InvalidArgument,
		&common.I18NItem{
			Key:        "errors.DocumentsService.ErrStampLimitReached.content",
			Parameters: map[string]string{"max": "5"},
		},
		&common.I18NItem{Key: "errors.DocumentsService.ErrStampLimitReached.title"},
	)
)
