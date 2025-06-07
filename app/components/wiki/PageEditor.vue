<script lang="ts" setup>
import type { UForm } from '#components';
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useNotificatorStore } from '~/stores/notificator';
import type { Content } from '~/types/history';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { File } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { PageJobAccess, PageUserAccess } from '~~/gen/ts/resources/wiki/access';
import { AccessLevel } from '~~/gen/ts/resources/wiki/access';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';
import BackButton from '../partials/BackButton.vue';
import ConfirmModal from '../partials/ConfirmModal.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
import { pageToURL } from './helpers';

const props = defineProps<{
    pageId: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const modal = useModal();

const { attr, activeChar } = useAuth();

const historyStore = useHistoryStore();

const route = useRoute<'wiki-job-id-slug-edit'>();

const {
    data: page,
    pending: loading,
    error,
    refresh,
} = useLazyAsyncData(`wiki-page:${route.path}`, () => getPage(parseInt(route.params.id)));

async function getPage(id: number): Promise<Page | undefined> {
    try {
        const call = $grpc.wiki.wiki.getPage({
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

const notifications = useNotificatorStore();

const { maxAccessEntries } = useAppConfig();

const { ydoc, provider } = useCollabDoc('wiki', props.pageId);

watchOnce(page, () => provider.connect());

const canDo = computed(() => ({
    public: attr('wiki.WikiService.UpdatePage', 'Fields', 'Public').value,
}));

const schema = z.object({
    parentId: z.number(),
    meta: z.object({
        title: z.string().min(3).max(255),
        description: z.string().max(255),
        public: z.boolean(),
        draft: z.boolean(),
        toc: z.boolean(),
    }),
    content: z.string().min(3).max(1750000),
    access: z.object({
        jobs: z.custom<PageJobAccess>().array().max(maxAccessEntries),
        users: z.custom<PageUserAccess>().array().max(maxAccessEntries),
    }),
    files: z.custom<File>().array().max(5),
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
        const call = $grpc.wiki.wiki.listPages({
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

provider.once('loadContent', () => setFromProps());

async function updatePage(values: Schema): Promise<void> {
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
            version: '',
            rawContent: values.content,
        },
        parentId: values.parentId,
        access: values.access,
        files: values.files,
    };

    try {
        let responsePage: Page | undefined = undefined;
        const call = $grpc.wiki.wiki.updatePage({
            page: req,
        });
        const { response } = await call;
        responsePage = response.page;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
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
        slot: 'content',
        label: t('common.content'),
        icon: 'i-mdi-pencil',
    },
    {
        slot: 'access',
        label: t('common.access', 1),
        icon: 'i-mdi-key',
    },
];

const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});

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
        class="min-h-dscreen flex w-full max-w-full flex-1 flex-col overflow-y-auto"
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
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.page', 1)])" />
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
                    list: { rounded: '' },
                }"
            >
                <template #content>
                    <UDashboardToolbar>
                        <template #default>
                            <div class="flex w-full flex-col gap-2">
                                <UFormGroup
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
                                                value-attribute="id"
                                                searchable-lazy
                                                :options="parentPages"
                                            >
                                                <template #label>
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
                                </UFormGroup>

                                <UFormGroup name="meta.title" :label="$t('common.title')">
                                    <UInput v-model="state.meta.title" size="xl" />
                                </UFormGroup>

                                <UFormGroup name="meta.description" :label="$t('common.description')">
                                    <UTextarea v-model="state.meta.description" :rows="2" />
                                </UFormGroup>
                            </div>
                        </template>
                    </UDashboardToolbar>

                    <UFormGroup
                        class="flex flex-1 overflow-y-hidden"
                        name="content"
                        :ui="{ container: 'flex flex-1 flex-col mt-0 overflow-y-hidden', label: { wrapper: 'hidden' } }"
                        label="&nbsp;"
                    >
                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.content"
                                v-model:files="state.files"
                                class="mx-auto w-full max-w-screen-xl flex-1 overflow-y-hidden"
                                history-type="wiki"
                                :target-id="page?.id"
                                filestore-namespace="wiki"
                                :filestore-service="(opts) => $grpc.wiki.wiki.uploadFile(opts)"
                            >
                                <template #linkModal="{ state: linkState }">
                                    <UDivider class="mt-1" :label="$t('common.or')" orientation="horizontal" />

                                    <UFormGroup class="w-full" name="url" :label="`${$t('common.wiki')} ${$t('common.page')}`">
                                        <ClientOnly>
                                            <USelectMenu
                                                label-attribute="title"
                                                searchable-lazy
                                                :options="pages"
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
                                    </UFormGroup>
                                </template>
                            </TiptapEditor>
                        </ClientOnly>
                    </UFormGroup>

                    <UDashboardToolbar
                        class="flex shrink-0 justify-between border-b-0 border-t border-gray-200 px-3 py-3.5 dark:border-gray-700"
                    >
                        <div class="flex flex-1 gap-2">
                            <UFormGroup class="flex-1" name="public" :label="$t('common.public')">
                                <UToggle v-model="state.meta.public" :disabled="!canDo.public" />
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="closed" :label="`${$t('common.toc', 2)}?`">
                                <UToggle v-model="state.meta.toc" />
                            </UFormGroup>
                        </div>
                    </UDashboardToolbar>
                </template>

                <template #access>
                    <div class="flex flex-1 flex-col gap-2 overflow-y-scroll px-2">
                        <UFormGroup name="access" :label="$t('common.access')">
                            <AccessManager
                                v-model:jobs="state.access.jobs"
                                v-model:users="state.access.users"
                                :target-id="page.id ?? 0"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.wiki.AccessLevel')"
                            />
                        </UFormGroup>
                    </div>
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
