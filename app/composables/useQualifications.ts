import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { CreateQualificationResponse } from '~~/gen/ts/services/qualifications/qualifications';

export function useQualifications() {
    const { $grpc } = useNuxtApp();

    async function createQualification(): Promise<CreateQualificationResponse> {
        try {
            const call = $grpc.qualifications.qualifications.createQualification({
                contentType: ContentType.HTML,
            });
            const { response } = await call;

            await navigateTo({
                name: 'qualifications-id-edit',
                params: { id: response.qualificationId },
            });

            return response;
        } catch (e) {
            handleGRPCError(e as RpcError);
            throw e;
        }
    }

    return {
        createQualification,
    };
}
