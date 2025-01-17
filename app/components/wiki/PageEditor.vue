<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/store/notificator';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { PageJobAccess, PageUserAccess } from '~~/gen/ts/resources/wiki/access';
import { AccessLevel } from '~~/gen/ts/resources/wiki/access';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';
import AccessManager from '../partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '../partials/access/helpers';
import TiptapEditor from '../partials/editor/TiptapEditor.vue';
import { pageToURL } from './helpers';

const props = defineProps<{
    modelValue?: Page | undefined;
    pages: PageShort[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: Page | undefined): void;
    (e: 'close'): void;
}>();

const { t } = useI18n();

const { attr, activeChar } = useAuth();

const page = computed({
    get() {
        return props.modelValue
            ? props.modelValue
            : ({
                  id: '0',
                  job: activeChar.value?.job ?? '',
                  path: '/wiki/' + (activeChar.value?.job ?? ''),
                  meta: {
                      contentType: ContentType.HTML,
                      public: false,
                      title: '',
                      description: '',
                      tags: [],
                  },
                  content: {
                      version: '',
                      rawContent: '',
                  },
                  access: {
                      jobs: [
                          {
                              id: '0',
                              targetId: '0',
                              job: activeChar.value?.job ?? '',
                              minimumGrade: 1,
                              access: AccessLevel.VIEW,
                          },
                          {
                              id: '0',
                              targetId: '0',
                              job: activeChar.value?.job ?? '',
                              minimumGrade: -1,
                              access: AccessLevel.EDIT,
                          },
                      ],
                      users: [],
                  },
              } as Page);
    },
    set(value) {
        emit('update:modelValue', value);
    },
});

const notifications = useNotificatorStore();

const { maxAccessEntries } = useAppConfig();

const canDo = computed(() => ({
    public: attr('WikiService.CreatePage', 'Fields', 'Public').value,
}));

const schema = z.object({
    parentId: z.string().optional(),
    meta: z.object({
        title: z.string().min(3).max(255),
        description: z.string().max(255),
        public: z.boolean(),
        toc: z.boolean(),
    }),
    content: z.string().min(3).max(1750000),
    access: z.object({
        jobs: z.custom<PageJobAccess>().array().max(maxAccessEntries),
        users: z.custom<PageUserAccess>().array().max(maxAccessEntries),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    parentId: undefined,
    meta: {
        title: page.value?.meta?.title ?? '',
        description: page.value?.meta?.description ?? '',
        public: page.value?.meta?.public ?? false,
        toc: page.value?.meta?.toc ?? true,
    },
    content: page.value?.content?.rawContent ?? '',
    access: {
        jobs: [],
        users: [],
    },
});

const createPage = computed(() => page.value.id === '0');

function setFromProps(): void {
    state.parentId =
        page.value?.meta?.createdAt !== undefined && page.value?.parentId === undefined
            ? undefined
            : (page.value?.parentId ??
              (props.pages.length === 0
                  ? undefined
                  : props.pages.at(0)?.job !== undefined && props.pages.at(0)?.job === activeChar.value?.job
                    ? props.pages.at(0)?.id
                    : undefined));

    state.meta.title = page.value.meta?.title ?? '';
    state.meta.description = page.value.meta?.description ?? '';
    state.meta.public = page.value.meta?.public ?? false;
    state.meta.toc = page.value.meta?.toc ?? true;
    state.content = page.value.content?.rawContent ?? '';
    if (page.value.access) {
        state.access = page.value.access;
    }
}

setFromProps();

async function createOrUpdatePage(values: Schema): Promise<void> {
    const req: Page = {
        id: props.modelValue?.id ?? '0',
        job: props.modelValue?.job ?? '',
        meta: {
            title: values.meta.title,
            description: values.meta.description,
            contentType: ContentType.HTML,
            public: values.meta.public,
            tags: [],
        },
        content: {
            version: '',
            rawContent: values.content,
        },
        parentId: values.parentId,
        access: values.access,
    };

    try {
        let responsePage: Page | undefined = undefined;
        if (createPage.value) {
            const call = getGRPCWikiClient().createPage({
                page: req,
            });
            const { response } = await call;
            responsePage = response.page;
        } else {
            const call = getGRPCWikiClient().updatePage({
                page: req,
            });
            const { response } = await call;
            responsePage = response.page;
        }

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (responsePage) {
            page.value = responsePage;
        }

        if (createPage.value) {
            navigateTo({
                name: 'wiki-job-id-slug',
                params: {
                    job: responsePage!.job,
                    id: responsePage!.id,
                    slug: [responsePage!.meta!.slug ?? ''],
                },
            });
        } else {
            emit('close');
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

type PageItem = { id: string; title: string };

function pageChildrenToList(p: PageShort, prefix?: string): PageItem[] {
    const list = [];

    list.push({
        id: p.id,
        title: (prefix !== undefined ? `${prefix} > ` : '') + p.title,
    });
    if (p.children.length > 0) {
        p.children.filter((c) => c.id !== p.id).forEach((c) => list.push(...pageChildrenToList(c, p.title)));
    }

    return list;
}

const parentPages = computedAsync(() =>
    props.pages
        .filter((p) => !props.modelValue?.id || p.id === props.modelValue?.id)
        .flatMap((p) => pageChildrenToList(p))
        .sort((a, b) => a.title.localeCompare(b.title)),
);

const items = [
    {
        slot: 'edit',
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
const route = useRoute();

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
    await createOrUpdatePage(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm
        :schema="schema"
        :state="state"
        class="flex min-h-screen w-full max-w-full flex-1 flex-col overflow-y-auto"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="$t('common.wiki')">
            <template #right>
                <UButton
                    color="black"
                    icon="i-mdi-arrow-left"
                    :disabled="!canSubmit"
                    @click="createPage ? navigateTo({ name: 'wiki' }) : $emit('close')"
                >
                    {{ $t('common.back') }}
                </UButton>

                <UButton type="submit" class="ml-2" trailing-icon="i-mdi-content-save" :disabled="!canSubmit">
                    <span class="hidden truncate sm:block">
                        <template v-if="!page.id">
                            {{ $t('common.create') }}
                        </template>
                        <template v-else>
                            {{ $t('common.save') }}
                        </template>
                    </span>
                </UButton>
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent class="p-0">
            <UTabs
                v-model="selectedTab"
                :items="items"
                class="flex flex-1 flex-col"
                :ui="{
                    wrapper: 'space-y-0 overflow-y-hidden',
                    container: 'flex flex-1 flex-col overflow-y-hidden',
                    base: 'flex flex-1 flex-col overflow-y-hidden',
                    list: { rounded: '' },
                }"
            >
                <template #edit>
                    <UDashboardToolbar>
                        <template #default>
                            <div class="flex w-full flex-col gap-2">
                                <UFormGroup
                                    v-if="!(modelValue?.meta?.createdAt && modelValue?.parentId === undefined)"
                                    name="meta.parentId"
                                    :label="$t('common.parent_page')"
                                    class="w-full"
                                >
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="state.parentId"
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

                                            <template #empty> {{ $t('common.not_found', [$t('common.page', 2)]) }} </template>
                                        </USelectMenu>
                                    </ClientOnly>
                                </UFormGroup>

                                <UFormGroup name="meta.title" :label="$t('common.title')">
                                    <UInput v-model="state.meta.title" size="xl" />
                                </UFormGroup>

                                <UFormGroup name="meta.description" :label="$t('common.description')">
                                    <UTextarea v-model="state.meta.description" />
                                </UFormGroup>
                            </div>
                        </template>
                    </UDashboardToolbar>

                    <UFormGroup
                        name="content"
                        class="flex flex-1 overflow-y-hidden"
                        :ui="{ container: 'flex flex-1 mt-0 overflow-y-hidden', label: { wrapper: 'hidden' } }"
                        label="&nbsp;"
                    >
                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.content"
                                class="mx-auto w-full max-w-screen-xl flex-1 overflow-y-hidden"
                                rounded="rounded-none"
                            >
                                <template #linkModal="{ state: linkState }">
                                    <UDivider :label="$t('common.or')" orientation="horizontal" class="mt-1" />

                                    <UFormGroup name="url" :label="`${$t('common.wiki')} ${$t('common.page')}`" class="w-full">
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
                            <UFormGroup name="public" :label="$t('common.public')" class="flex-1">
                                <UToggle v-model="state.meta.public" :disabled="!canDo.public" />
                            </UFormGroup>

                            <UFormGroup name="closed" :label="`${$t('common.toc', 2)}?`" class="flex-1">
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
                                :target-id="page.id ?? '0'"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.wiki.AccessLevel')"
                            />
                        </UFormGroup>
                    </div>
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
