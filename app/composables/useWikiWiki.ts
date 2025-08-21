import { getWikiWikiClient } from '~~/gen/ts/clients';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { CreatePageResponse } from '~~/gen/ts/services/wiki/wiki';

export async function useWikiWiki() {
    const wikiWikiClient = await getWikiWikiClient();

    async function createPage(parentId?: number): Promise<CreatePageResponse> {
        try {
            const call = wikiWikiClient.createPage({
                contentType: ContentType.HTML,
                parentId: parentId,
            });
            const { response } = await call;

            await navigateTo({
                name: 'wiki-job-id-slug-edit',
                params: {
                    job: response.job,
                    id: response.id,
                    slug: '',
                },
            });

            return response;
        } catch (e) {
            handleGRPCError(e);
            throw e;
        }
    }

    return {
        createPage,
    };
}
