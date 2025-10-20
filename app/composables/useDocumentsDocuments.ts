import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type {
    CreateDocumentResponse,
    GetDocumentResponse,
    ListDocumentPinsResponse,
    ListDocumentsRequest,
    ListDocumentsResponse,
    ToggleDocumentPinResponse,
} from '~~/gen/ts/services/documents/documents';

export async function useDocumentsDocuments() {
    const documentsDocumentsClient = await getDocumentsDocumentsClient();

    const notifications = useNotificationsStore();

    const clipboardStore = useClipboardStore();
    const { getTemplateData } = clipboardStore;

    const documents = ref<ListDocumentsResponse | undefined>();
    const pinnedDocuments = ref<ListDocumentPinsResponse | undefined>();

    const listDocuments = async (req: ListDocumentsRequest): Promise<ListDocumentsResponse> => {
        try {
            const call = documentsDocumentsClient.listDocuments(req);
            const { response } = await call;

            documents.value = response;

            return response;
        } catch (e) {
            handleGRPCError(e as RpcError);
            throw e;
        }
    };

    const getDocument = async (id: number, redirectOnError?: boolean): Promise<GetDocumentResponse> => {
        try {
            const call = documentsDocumentsClient.getDocument({
                documentId: id,
            });
            const { response } = await call;

            return response;
        } catch (e) {
            handleGRPCError(e as RpcError);

            if (redirectOnError === true) await navigateTo({ name: 'documents' });
            throw e;
        }
    };

    const createDocument = async (templateId?: number): Promise<CreateDocumentResponse> => {
        const { activeChar } = useAuth();

        const templateData = getTemplateData();
        templateData.activeChar = unref(activeChar.value!);

        try {
            const call = documentsDocumentsClient.createDocument({
                contentType: ContentType.HTML,
                templateId: templateId,
                templateData: templateData,
            });
            const { response } = await call;

            await navigateTo({
                name: 'documents-id-edit',
                params: {
                    id: response.id,
                },
            });

            return response;
        } catch (e) {
            handleGRPCError(e);
            throw e;
        }
    };

    const deleteDocument = async (id: number, restore?: boolean, reason?: string): Promise<boolean> => {
        try {
            await documentsDocumentsClient.deleteDocument({
                documentId: id,
                reason: reason,
            });

            // Navigate to document list when deletedAt timestamp is undefined
            if (!restore) {
                notifications.add({
                    title: { key: 'notifications.document_deleted.title', parameters: {} },
                    description: { key: 'notifications.document_deleted.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });

                await navigateTo({ name: 'documents' });
                return false;
            } else {
                notifications.add({
                    title: { key: 'notifications.document_restored.title', parameters: {} },
                    description: { key: 'notifications.document_restored.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });
                return true;
            }
        } catch (e) {
            handleGRPCError(e as RpcError);
            throw e;
        }
    };

    const toggleDocument = async (id: number, closed: boolean): Promise<boolean> => {
        try {
            await documentsDocumentsClient.toggleDocument({
                documentId: id,
                closed: closed,
            });

            if (!closed) {
                notifications.add({
                    title: { key: `notifications.documents.document_toggled.open.title`, parameters: {} },
                    description: { key: `notifications.documents.document_toggled.open.content`, parameters: {} },
                    type: NotificationType.SUCCESS,
                });
            } else {
                notifications.add({
                    title: { key: `notifications.documents.document_toggled.closed.title`, parameters: {} },
                    description: { key: `notifications.documents.document_toggled.closed.content`, parameters: {} },
                    type: NotificationType.SUCCESS,
                });
            }

            return closed;
        } catch (e) {
            handleGRPCError(e as RpcError);
            throw e;
        }
    };

    const changeDocumentOwner = async (id: number): Promise<void> => {
        try {
            await documentsDocumentsClient.changeDocumentOwner({
                documentId: id,
            });

            notifications.add({
                title: { key: 'notifications.documents.document_take_ownership.title', parameters: {} },
                description: { key: 'notifications.documents.document_take_ownership.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        } catch (e) {
            handleGRPCError(e as RpcError);
            throw e;
        }
    };

    const listDocumentPins = async (page: number): Promise<ListDocumentPinsResponse> => {
        const call = documentsDocumentsClient.listDocumentPins({
            pagination: {
                offset: calculateOffset(page, pinnedDocuments.value?.pagination),
            },
        });
        const { response } = await call;

        pinnedDocuments.value = response;

        return response;
    };

    const togglePin = async (documentId: number, state: boolean, personal: boolean): Promise<ToggleDocumentPinResponse> => {
        try {
            const call = documentsDocumentsClient.toggleDocumentPin({
                documentId: documentId,
                state: state,
                personal: personal,
            });
            const { response } = await call;

            return response;
        } catch (e) {
            handleGRPCError(e as RpcError);
            throw e;
        }
    };

    return {
        // State
        documents,
        pinnedDocuments,

        // Actions
        listDocuments,
        getDocument,
        createDocument,
        deleteDocument,
        toggleDocument,
        changeDocumentOwner,
        listDocumentPins,
        togglePin,
    };
}
