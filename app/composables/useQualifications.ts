import { getQualificationsQualificationsClient } from '~~/gen/ts/clients';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { CreateQualificationResponse } from '~~/gen/ts/services/qualifications/qualifications';

export async function useQualifications() {
    const qualificationsQualificationsClient = await getQualificationsQualificationsClient();

    async function createQualification(): Promise<CreateQualificationResponse> {
        try {
            const call = qualificationsQualificationsClient.createQualification({
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
