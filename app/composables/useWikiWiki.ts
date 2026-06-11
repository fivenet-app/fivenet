import { getWikiWikiClient } from '~~/gen/ts/clients';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { Page } from '~~/gen/ts/resources/wiki/page';
import type {
    CreatePageResponse,
    DeletePageResponse,
    GetPageResponse,
    ListPageActivityRequest,
    ListPageActivityResponse,
    ListPagesRequest,
    ListPagesResponse,
    MovePageResponse,
    UpdatePageResponse,
} from '~~/gen/ts/services/wiki/wiki';

export async function useWikiWiki() {
    const wikiWikiClient = await getWikiWikiClient();

    async function runWikiCall<T>(call: any): Promise<T> {
        try {
            const { response } = await call;
            return response as T;
        } catch (e) {
            handleGRPCError(e);
            throw e;
        }
    }

    async function createPage(parentId?: number): Promise<CreatePageResponse> {
        const response = await runWikiCall<CreatePageResponse>(
            wikiWikiClient.createPage({
                contentType: ContentType.HTML,
                parentId: parentId,
            }),
        );

        await navigateTo({
            name: 'wiki-job-id-slug-edit',
            params: {
                job: response.job,
                id: response.id,
                slug: '',
            },
        });

        return response;
    }

    async function movePage(payload: { pageId: number; beforeId?: number; afterId?: number }): Promise<MovePageResponse> {
        return runWikiCall<MovePageResponse>(
            wikiWikiClient.movePage({
                pageId: payload.pageId,
                beforeId: payload.beforeId,
                afterId: payload.afterId,
            }),
        );
    }

    async function listPages(request: ListPagesRequest): Promise<ListPagesResponse> {
        return runWikiCall<ListPagesResponse>(wikiWikiClient.listPages(request));
    }

    async function getPage(id: number): Promise<Page | undefined> {
        const response = await runWikiCall<GetPageResponse>(
            wikiWikiClient.getPage({
                id: id,
            }),
        );

        return response.page;
    }

    async function updatePage(page: Page): Promise<Page | undefined> {
        const response = await runWikiCall<UpdatePageResponse>(
            wikiWikiClient.updatePage({
                page: page,
            }),
        );

        return response.page;
    }

    async function deletePage(id: number): Promise<void> {
        await runWikiCall<DeletePageResponse>(
            wikiWikiClient.deletePage({
                id: id,
            }),
        );
    }

    async function listPageActivity(request: ListPageActivityRequest): Promise<ListPageActivityResponse> {
        return runWikiCall<ListPageActivityResponse>(wikiWikiClient.listPageActivity(request));
    }

    return {
        deletePage,
        getPage,
        createPage,
        listPageActivity,
        listPages,
        movePage,
        updatePage,
    };
}
