<script lang="ts" setup>
import type { UForm } from '#components';
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import type { Content } from '~/types/history';
import { getWikiWikiClient } from '~~/gen/ts/clients';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { File } from '~~/gen/ts/resources/file/file';
import { ObjectType } from '~~/gen/ts/resources/notifications/client_view';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { type PageJobAccess, type PageUserAccess, AccessLevel } from '~~/gen/ts/resources/wiki/access';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';
import BackButton from '../partials/BackButton.vue';
import ConfirmModal from '../partials/ConfirmModal.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
import { checkPageAccess, pageToURL } from './helpers';

const props = defineProps<{
    pageId: number;
}>();

const { t } = useI18n();

const modal = useOverlay();

const { attr, activeChar } = useAuth();

const historyStore = useHistoryStore();

const route = useRoute<'wiki-job-id-slug-edit'>();

const wikiClient = await getWikiWikiClient();

const {
    data: page,
    status,
    error,
    refresh,
} = useLazyAsyncData(`wiki-page:${route.path}`, () => getPage(parseInt(route.params.id)));

async function getPage(id: number): Promise<Page | undefined> {
    try {
        const call = wikiClient.getPage({
            id: id,
        });
        const { response } = await call;

        return response.page;
    } catch (e) {
        handleGRPCError(e as RpcError);

        await navigateTo({
            name: 'wiki-job-id-slug',
            params: { job: route.params.job, id: route.params.id, slug: [route.params.slug] },
        });
        return;
    }
}

const notifications = useNotificationsStore();

const { maxAccessEntries } = useAppConfig();

const { ydoc, provider } = await useCollabDoc('wiki', props.pageId);

const canDo = computed(() => ({
    access: checkPageAccess(page.value?.access, page.value?.meta?.creator, AccessLevel.ACCESS),
    edit: checkPageAccess(page.value?.access, page.value?.meta?.creator, AccessLevel.EDIT),
    public: attr('wiki.WikiService/UpdatePage', 'Fields', 'Public').value,
}));

const schema = z.object({
    parentId: z.coerce.number(),
    meta: z.object({
        title: z.string().min(3).max(255),
        description: z.string().max(255),
        public: z.coerce.boolean(),
        draft: z.coerce.boolean(),
        toc: z.coerce.boolean(),
    }),
    content: z.string().min(3).max(1750000),
    access: z.object({
        jobs: z.custom<PageJobAccess>().array().max(maxAccessEntries).default([]),
        users: z.custom<PageUserAccess>().array().max(maxAccessEntries).default([]),
    }),
    files: z.custom<File>().array().max(5).default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    parentId: 0,
    meta: {
        title: '',
        description: '',
        public: false,
        draft: true,
        toc: true,
    },
    content: '',
    access: {
        jobs: [],
        users: [],
    },
    files: [],
});

const { data: pages, refresh: pagesRefresh } = useLazyAsyncData(`wiki-pages-${props.pageId}-editor`, () => listPages(), {
    default: () => [],
});

async function listPages(): Promise<PageShort[]> {
    const job = route.params.job ?? activeChar.value?.job ?? '';
    try {
        const call = wikiClient.listPages({
            pagination: {
                offset: 0,
            },
            job: job,
            rootOnly: false,
        });
        const { response } = await call;

        return response.pages;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const changed = ref(false);
const saving = ref(false);

// Track last saved string and timestamp
let lastSavedString = '';
let lastSaveTimestamp = 0;

async function saveHistory(values: Schema, name: string | undefined = undefined, type = 'wiki'): Promise<void> {
    if (saving.value) {
        return;
    }

    const now = Date.now();
    // Skip if identical to last saved or if within MIN_GAP
    if (state.content === lastSavedString || now - lastSaveTimestamp < 5000) {
        return;
    }

    saving.value = true;

    historyStore.addVersion<Content>(
        type,
        props.pageId,
        {
            content: values.content,
            files: values.files,
        },
        name,
    );

    useTimeoutFn(() => {
        saving.value = false;
    }, 1750);

    lastSavedString = state.content;
    lastSaveTimestamp = now;
}

historyStore.handleRefresh(() => saveHistory(state, 'wiki'));

watchDebounced(
    state,
    () => {
        if (changed.value) {
            saveHistory(state);
        } else {
            changed.value = true;
        }
    },
    {
        debounce: 1000,
        maxWait: 2500,
    },
);

function setFromProps(): void {
    if (!page.value) return;

    state.parentId =
        (page.value?.meta?.createdAt !== undefined && page.value?.parentId === undefined
            ? undefined
            : (page.value?.parentId ??
              (pages.value.length === 0
                  ? undefined
                  : pages.value.at(0)?.job !== undefined && pages.value.at(0)?.job === activeChar.value?.job
                    ? pages.value.at(0)?.id
                    : undefined))) ?? 0;

    state.meta.title = page.value.meta?.title ?? '';
    state.meta.description = page.value.meta?.description ?? '';
    state.content = page.value.content?.rawContent ?? '';
    state.meta.public = page.value.meta?.public ?? false;
    state.meta.toc = page.value.meta?.toc ?? true;
    state.meta.draft = page.value.meta?.draft ?? true;
    if (page.value.access) {
        state.access.jobs = page.value.access.jobs;
        state.access.users = page.value.access.users;
    }
    state.files = page.value.files;
}

const onSync = (s: boolean) => {
    if (!s) return;

    if (ydoc.getXmlFragment('content').length === 0) {
        logger.info('PageEditor - Content is empty, setting from props');
        // If the content is empty, we need to set it from the props
        setFromProps();
    }
    provider.off('sync', onSync);
};
provider.on('sync', onSync);

async function updatePage(values: Schema): Promise<void> {
    values.access.users.forEach((user) => {
        if (user.id < 0) user.id = 0;
        user.user = undefined; // Clear user object to avoid sending unnecessary data
    });
    values.access.jobs.forEach((job) => job.id < 0 && (job.id = 0));

    const req: Page = {
        id: page.value?.id ?? 0,
        job: page.value?.job ?? '',
        meta: {
            title: values.meta.title,
            description: values.meta.description,
            contentType: ContentType.HTML,
            public: values.meta.public,
            toc: values.meta.toc,
            draft: values.meta.draft,
            tags: [],
        },
        content: {
            rawContent: values.content,
        },
        parentId: values.parentId,
        access: values.access,
        files: values.files,
    };

    try {
        let responsePage: Page | undefined = undefined;
        const call = wikiClient.updatePage({
            page: req,
        });
        const { response } = await call;
        responsePage = response.page;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (responsePage) {
            page.value = responsePage;
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

type PageItem = { id: number; title: string; draft: boolean };

function pageChildrenToList(p: PageShort, prefix?: string): PageItem[] {
    const list = [];

    list.push({
        id: p.id,
        title: (prefix !== undefined ? `${prefix} > ` : '') + p.title,
        draft: p.draft,
    });
    if (p.children.length > 0) {
        p.children.filter((c) => c.id !== p.id).forEach((c) => list.push(...pageChildrenToList(c, p.title)));
    }

    return list;
}

const parentPages = computedAsync(() => {
    const pagesList = pages.value
        .flatMap((p) => pageChildrenToList(p))
        .filter((p) => !page.value?.id || p.id !== page.value?.id)
        .sort((a, b) => {
            const aDraft = a?.draft ?? false;
            const bDraft = b?.draft ?? false;

            if (aDraft !== bDraft) {
                // Drafts go last
                return aDraft ? 1 : -1;
            }
            return a.title.localeCompare(b.title);
        });

    if (page.value?.parentId && pagesList.find((p) => p.id === page.value?.parentId) === undefined) {
        pagesList.unshift({
            id: page.value.parentId,
            title: `${t('common.parent_page')} - ${t('common.id')}: ${page.value.parentId}`,
            draft: page.value.meta?.draft ?? false,
        });
    }

    return pagesList;
});

const items = [
    {
        slot: 'content' as const,
        label: t('common.content'),
        icon: 'i-mdi-pencil',
        value: 'content',
    },
    {
        slot: 'access' as const,
        label: t('common.access'),
        icon: 'i-mdi-key',
        value: 'access',
    },
];

const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'content';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

// Handle the client update event
const { sendClientView } = useClientUpdate(ObjectType.WIKI_PAGE, () =>
    notifications.add({
        title: { key: 'notifications.wiki.client_view_update.title', parameters: {} },
        description: { key: 'notifications.wiki.client_view_update.content', parameters: {} },
        duration: 7500,
        type: NotificationType.INFO,
        actions: [
            {
                label: { key: 'common.refresh', parameters: {} },
                icon: 'i-mdi-refresh',
                onClick: () => refresh(),
            },
        ],
    }),
);
sendClientView(props.pageId);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await updatePage(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

useYText(ydoc.getText('title'), toRef(state.meta, 'title'), { provider: provider });
useYText(ydoc.getText('description'), toRef(state.meta, 'description'), { provider: provider });
const detailsYdoc = ydoc.getMap('details');
useYNumber(detailsYdoc, 'parentId', toRef(state, 'parentId'), { provider: provider });
useYBoolean(detailsYdoc, 'public', toRef(state.meta, 'public'), { provider: provider });
useYBoolean(detailsYdoc, 'toc', toRef(state.meta, 'toc'), { provider: provider });
useYBoolean(detailsYdoc, 'draft', toRef(state.meta, 'draft'), { provider: provider });

// Access
useYArrayFiltered<PageJobAccess>(
    ydoc.getArray('access_jobs'),
    toRef(state.access, 'jobs'),
    { omit: ['createdAt', 'user'] },
    { provider: provider },
);
useYArrayFiltered<PageUserAccess>(
    ydoc.getArray('access_users'),
    toRef(state.access, 'users'),
    {
        omit: ['createdAt', 'user'],
    },
    { provider: provider },
);

// Files
useYArrayFiltered<File>(
    ydoc.getArray('files'),
    toRef(state, 'files'),
    {
        omit: ['createdAt', 'meta'],
    },
    { provider: provider },
);

provide('yjsDoc', ydoc);
provide('yjsProvider', provider);

const formRef = useTemplateRef<typeof UForm>('formRef');
</script>

<template>
    <UForm
        ref="formRef"
        class="flex min-h-dvh w-full max-w-full flex-1 flex-col overflow-y-auto"
        :schema="schema"
        :state="state"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="$t('common.wiki')">
            <template #right>
                <BackButton :disabled="!canSubmit" />

                <UButton v-if="page" type="submit" trailing-icon="i-mdi-content-save" :disabled="!canSubmit">
                    <span class="hidden truncate sm:block">
                        {{ $t('common.save') }}
                    </span>
                </UButton>

                <UButton
                    v-if="page?.meta?.draft"
                    type="submit"
                    color="info"
                    trailing-icon="i-mdi-publish"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    @click.prevent="
                        modal.open(ConfirmModal, {
                            title: $t('common.publish_confirm.title', { type: $t('common.document', 1) }),
                            description: $t('common.publish_confirm.description'),
                            color: 'info',
                            iconClass: 'text-info-500 dark:text-info-400',
                            icon: 'i-mdi-publish',
                            confirm: () => {
                                state.meta.draft = false;
                                formRef?.submit();
                            },
                        })
                    "
                >
                    <span class="hidden truncate sm:block">
                        {{ $t('common.publish') }}
                    </span>
                </UButton>
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent class="p-0 sm:pb-0">
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.page', 1)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.page', 1)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!page"
                icon="i-mdi-file-search"
                :message="$t('common.not_found', [$t('common.page', 1)])"
            />
            <UTabs
                v-else
                v-model="selectedTab"
                class="flex flex-1 flex-col"
                :items="items"
                :ui="{
                    wrapper: 'space-y-0 overflow-y-hidden',
                    container: 'flex flex-1 flex-col overflow-y-hidden',
                    base: 'flex flex-1 flex-col overflow-y-hidden',
                }"
            >
                <template #content>
                    <UDashboardToolbar>
                        <template #default>
                            <div class="flex w-full flex-col gap-2">
                                <UFormField
                                    v-if="!(page?.meta?.createdAt && page?.parentId === undefined)"
                                    class="w-full"
                                    name="meta.parentId"
                                    :label="$t('common.parent_page')"
                                >
                                    <div class="flex items-center gap-1">
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="state.parentId"
                                                class="flex-1"
                                                value-key="id"
                                                searchable-lazy
                                                :disabled="!canDo.edit"
                                                :items="parentPages"
                                            >
                                                <template #item-label>
                                                    <span class="truncate">
                                                        {{
                                                            state.parentId
                                                                ? (parentPages?.find((p) => p.id === state.parentId)?.title ??
                                                                  $t('common.na'))
                                                                : $t('common.none_selected', [$t('common.parent_page')])
                                                        }}
                                                    </span>
                                                </template>

                                                <template #option="{ option: opt }">
                                                    {{ opt.title }}
                                                </template>

                                                <template #option-empty="{ query: search }">
                                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                                </template>

                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.page', 2)]) }}
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>

                                        <UTooltip :text="$t('common.refresh')">
                                            <UButton variant="link" icon="i-mdi-refresh" @click="pagesRefresh()" />
                                        </UTooltip>
                                    </div>
                                </UFormField>

                                <UFormField name="meta.title" :label="$t('common.title')">
                                    <UInput v-model="state.meta.title" size="xl" :disabled="!canDo.edit" />
                                </UFormField>

                                <UFormField name="meta.description" :label="$t('common.description')">
                                    <UTextarea v-model="state.meta.description" :rows="2" :disabled="!canDo.edit" />
                                </UFormField>
                            </div>
                        </template>
                    </UDashboardToolbar>

                    <UFormField
                        class="flex flex-1 overflow-y-hidden"
                        name="content"
                        :ui="{ container: 'flex flex-1 flex-col mt-0 overflow-y-hidden', label: { wrapper: 'hidden' } }"
                        label="&nbsp;"
                    >
                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.content"
                                v-model:files="state.files"
                                class="max-w-(--breakpoint-xl) mx-auto w-full flex-1 overflow-y-hidden"
                                :disabled="!canDo.edit"
                                history-type="wiki"
                                :saving="saving"
                                enable-collab
                                :target-id="page?.id"
                                filestore-namespace="wiki"
                                :filestore-service="(opts) => wikiClient.uploadFile(opts)"
                            >
                                <template #linkModal="{ state: linkState }">
                                    <USeparator class="mt-1" :label="$t('common.or')" orientation="horizontal" />

                                    <UFormField class="w-full" name="url" :label="`${$t('common.wiki')} ${$t('common.page')}`">
                                        <ClientOnly>
                                            <USelectMenu
                                                label-key="title"
                                                searchable-lazy
                                                :items="pages"
                                                @update:model-value="($event) => (linkState.url = pageToURL($event, true))"
                                            >
                                                <template #option="{ option: opt }">
                                                    {{ opt.title }}
                                                </template>

                                                <template #option-empty="{ query: search }">
                                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                                </template>

                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.page', 2)]) }}
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>
                                    </UFormField>
                                </template>
                            </TiptapEditor>
                        </ClientOnly>
                    </UFormField>

                    <UDashboardToolbar
                        class="flex shrink-0 justify-between border-b-0 border-t border-gray-200 px-3 py-3.5 dark:border-gray-700"
                    >
                        <div class="flex flex-1 gap-2">
                            <UFormField class="flex-1" name="public" :label="$t('common.public')">
                                <USwitch v-model="state.meta.public" :disabled="!canDo.edit || !canDo.public" />
                            </UFormField>

                            <UFormField class="flex-1" name="closed" :label="`${$t('common.toc', 2)}?`">
                                <USwitch v-model="state.meta.toc" :disabled="!canDo.edit" />
                            </UFormField>
                        </div>
                    </UDashboardToolbar>
                </template>

                <template #access>
                    <div class="flex flex-1 flex-col gap-2 overflow-y-scroll px-2">
                        <UFormField name="access" :label="$t('common.access')">
                            <AccessManager
                                v-model:jobs="state.access.jobs"
                                v-model:users="state.access.users"
                                :disabled="!canDo.access"
                                :target-id="page.id ?? 0"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.wiki.AccessLevel')"
                            />
                        </UFormField>
                    </div>
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
