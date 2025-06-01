import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { CreateDocumentResponse } from '~~/gen/ts/services/documents/documents';

export function useDocumentsDocuments() {
    const { $grpc } = useNuxtApp();
    const clipboardStore = useClipboardStore();
    const { getTemplateData } = clipboardStore;

    async function createDocument(templateId?: number): Promise<CreateDocumentResponse> {
        const { activeChar } = useAuth();

        const templateData = getTemplateData();
        templateData.activeChar = unref(activeChar.value!);

        try {
            const call = $grpc.documents.documents.createDocument({
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
    }

    return {
        createDocument,
    };
}
