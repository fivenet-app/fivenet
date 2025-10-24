<script lang="ts" setup>
import type { UForm } from '#components';
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import type { Content } from '~/types/history';
import { jobAccessEntry, userAccessEntry } from '~/utils/validation';
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

const overlay = useOverlay();

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

useHead({
    title: () =>
        page.value?.meta?.title
            ? `${page.value.meta.title} - ${page.value.jobLabel} - ${t('pages.wiki.edit.title')}`
            : t('pages.wiki.edit.title'),
});

const notifications = useNotificationsStore();

const { maxAccessEntries } = useAppConfig();

const confirmModal = overlay.create(ConfirmModal);

const { ydoc, provider } = await useCollabDoc('wiki', props.pageId);

const canDo = computed(() => ({
    access: checkPageAccess(page.value?.access, page.value?.meta?.creator, AccessLevel.ACCESS, page.value?.job),
    edit: checkPageAccess(page.value?.access, page.value?.meta?.creator, AccessLevel.EDIT, page.value?.job),
    public: attr('wiki.WikiService/UpdatePage', 'Fields', 'Public').value,
}));

const schema = z.object({
    parentId: z.coerce.number(),
    meta: z.object({
        title: z.coerce.string().min(3).max(255),
        description: z.coerce.string().max(255),
        toc: z.coerce.boolean(),
        draft: z.coerce.boolean(),
        public: z.coerce.boolean(),
        startpage: z.coerce.boolean(),
    }),
    content: z.coerce.string().min(3).max(1750000),
    access: z.object({
        jobs: jobAccessEntry.array().max(maxAccessEntries).default([]),
        users: userAccessEntry.array().max(maxAccessEntries).default([]),
    }),
    files: z.custom<File>().array().max(5).default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    parentId: 0,
    meta: {
        title: '',
        description: '',
        toc: true,
        draft: true,
        public: false,
        startpage: false,
    },
    content: '',
    access: {
        jobs: [],
        users: [],
    },
    files: [],
});

const { data: pages, refresh: pagesRefresh } = useLazyAsyncData(`wiki-pages-id:${props.pageId}-editor`, () => listPages(), {
    default: () => [] as PageShort[],
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
    if (saving.value) return;

    const now = Date.now();
    // Skip if identical to last saved or if within MIN_GAP
    if (state.content === lastSavedString || now - lastSaveTimestamp < 5000) return;

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

    state.parentId = page.value?.parentId ?? 0;
    state.meta.title = page.value.meta?.title ?? '';
    state.meta.description = page.value.meta?.description ?? '';
    state.content = page.value.content?.rawContent ?? '';
    state.meta.toc = page.value.meta?.toc ?? true;
    state.meta.draft = page.value.meta?.draft ?? true;
    state.meta.public = page.value.meta?.public ?? false;
    state.meta.startpage = page.value.meta?.startpage ?? false;
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
            startpage: values.meta.startpage,
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

    pagesList.unshift({
        id: 0,
        title: t('common.none_selected', [t('common.parent_page')]),
        draft: false,
    });

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

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.wiki')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <BackButton :disabled="!canSubmit" />

                    <UButton
                        v-if="page"
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canSubmit"
                        :label="$t('common.save')"
                        :ui="{ label: 'hidden truncate sm:block' }"
                        @click="formRef?.submit()"
                    />

                    <UButton
                        v-if="page?.meta?.draft"
                        color="info"
                        trailing-icon="i-mdi-publish"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        :label="$t('common.publish')"
                        :ui="{ label: 'hidden truncate sm:block' }"
                        @click="
                            confirmModal.open({
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
                    />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <UForm
                ref="formRef"
                :schema="schema"
                :state="state"
                class="flex min-h-full w-full max-w-full flex-1 flex-col overflow-y-auto"
                @submit="onSubmitThrottle"
            >
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
                    variant="link"
                    :unmount-on-hide="false"
                    :ui="{ content: 'h-full' }"
                >
                    <template #content>
                        <UDashboardPanel :ui="{ root: 'h-full min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
                            <template #header>
                                <UDashboardToolbar>
                                    <template #default>
                                        <div class="mx-auto my-2 flex w-full max-w-(--breakpoint-xl) flex-col gap-2">
                                            <UFormField class="flex-1" name="meta.parentId" :label="$t('common.parent_page')">
                                                <div class="flex items-center gap-1">
                                                    <ClientOnly>
                                                        <USelectMenu
                                                            v-model="state.parentId"
                                                            class="flex-1"
                                                            value-key="id"
                                                            label-key="title"
                                                            :disabled="!canDo.edit"
                                                            :items="parentPages"
                                                        >
                                                            <template #default>
                                                                {{
                                                                    state.parentId
                                                                        ? (parentPages?.find((p) => p.id === state.parentId)
                                                                              ?.title ?? $t('common.na'))
                                                                        : $t('common.none_selected', [$t('common.parent_page')])
                                                                }}
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
                                                <UInput
                                                    v-model="state.meta.title"
                                                    size="xl"
                                                    class="w-full"
                                                    :disabled="!canDo.edit"
                                                />
                                            </UFormField>

                                            <UFormField name="meta.description" :label="$t('common.description')">
                                                <UTextarea
                                                    v-model="state.meta.description"
                                                    class="w-full"
                                                    :rows="2"
                                                    :disabled="!canDo.edit"
                                                />
                                            </UFormField>

                                            <div class="flex flex-1 gap-2">
                                                <UFormField
                                                    class="flex-1 md:grid md:grid-cols-2 md:items-center"
                                                    name="meta.public"
                                                    :label="$t('common.public')"
                                                >
                                                    <USwitch
                                                        v-model="state.meta.public"
                                                        :disabled="!canDo.edit || !canDo.public"
                                                    />
                                                </UFormField>

                                                <UFormField
                                                    v-if="canDo.public"
                                                    class="flex-1 md:grid md:grid-cols-2 md:items-center"
                                                    name="meta.startpage"
                                                    :label="`${$t('common.startpage')}?`"
                                                >
                                                    <USwitch
                                                        v-model="state.meta.startpage"
                                                        :disabled="!canDo.edit || !canDo.public"
                                                    />
                                                </UFormField>

                                                <UFormField
                                                    class="flex-1 md:grid md:grid-cols-2 md:items-center"
                                                    name="meta.toc"
                                                    :label="`${$t('common.toc', 2)}?`"
                                                >
                                                    <USwitch v-model="state.meta.toc" :disabled="!canDo.edit" />
                                                </UFormField>
                                            </div>
                                        </div>
                                    </template>
                                </UDashboardToolbar>
                            </template>

                            <template #body>
                                <ClientOnly>
                                    <TiptapEditor
                                        v-model="state.content"
                                        v-model:files="state.files"
                                        name="content"
                                        class="mx-auto my-2 h-full w-full max-w-(--breakpoint-xl) flex-1 overflow-y-hidden"
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

                                            <UFormField
                                                class="w-full"
                                                name="url"
                                                :label="`${$t('common.wiki')} ${$t('common.page')}`"
                                            >
                                                <ClientOnly>
                                                    <USelectMenu
                                                        label-key="title"
                                                        :items="pages"
                                                        @update:model-value="
                                                            ($event) => (linkState.url = pageToURL($event, true))
                                                        "
                                                    >
                                                        <template #empty>
                                                            {{ $t('common.not_found', [$t('common.page', 2)]) }}
                                                        </template>
                                                    </USelectMenu>
                                                </ClientOnly>
                                            </UFormField>
                                        </template>
                                    </TiptapEditor>
                                </ClientOnly>
                            </template>
                        </UDashboardPanel>
                    </template>

                    <template #access>
                        <UDashboardPanel :ui="{ root: 'min-h-0' }">
                            <template #body>
                                <div class="mx-auto w-full max-w-(--breakpoint-xl)">
                                    <UPageCard :title="$t('common.access')">
                                        <UFormField name="access">
                                            <AccessManager
                                                v-model:jobs="state.access.jobs"
                                                v-model:users="state.access.users"
                                                :disabled="!canDo.access"
                                                :target-id="page.id ?? 0"
                                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.wiki.AccessLevel')"
                                                name="access"
                                            />
                                        </UFormField>
                                    </UPageCard>
                                </div>
                            </template>
                        </UDashboardPanel>
                    </template>
                </UTabs>
            </UForm>
        </template>
    </UDashboardPanel>
</template>
