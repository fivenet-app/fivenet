<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { PageJobAccess, PageUserAccess } from '~~/gen/ts/resources/wiki/access';
import { AccessLevel } from '~~/gen/ts/resources/wiki/access';
import { ContentType, type Page, type PageShort } from '~~/gen/ts/resources/wiki/page';
import AccessManager from '../partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '../partials/access/helpers';
import DocEditor from '../partials/DocEditor.vue';

const props = defineProps<{
    modelValue?: Page | undefined;
    pages: PageShort[];
}>();

const emits = defineEmits<{
    (e: 'update:modelValue', value: Page | undefined): void;
    (e: 'close'): void;
}>();

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
                  content: '',
                  access: {
                      jobs: [
                          {
                              id: '0',
                              targetId: '0',
                              job: activeChar.value?.job ?? '',
                              minimumGrade: 1,
                              access: AccessLevel.VIEW,
                          },
                      ],
                      users: [],
                  },
              } as Page);
    },
    set(value) {
        emits('update:modelValue', value);
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
    content: page.value?.content ?? '',
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
    state.content = page.value.content;
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
        content: values.content,
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
                    slug: responsePage!.meta!.slug,
                },
            });
        } else {
            emits('close');
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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdatePage(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UDashboardNavbar :title="$t('common.wiki')">
            <template #right>
                <UButton
                    color="black"
                    icon="i-mdi-arrow-left"
                    @click="createPage ? navigateTo({ name: 'wiki' }) : $emit('close')"
                >
                    {{ $t('common.back') }}
                </UButton>

                <UButton type="submit" class="ml-2" trailing-icon="i-mdi-content-save">
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

        <div class="relative flex flex-1 flex-col overflow-x-auto px-8 py-2 pt-4">
            <UPage>
                <div class="flex flex-col gap-2">
                    <UFormGroup
                        v-if="!(modelValue?.meta?.createdAt && modelValue?.parentId === undefined)"
                        name="meta.parentId"
                        :label="$t('common.parent_page')"
                        class="w-full"
                    >
                        <ClientOnly>
                            <USelectMenu v-model="state.parentId" value-attribute="id" searchable-lazy :options="parentPages">
                                <template #label>
                                    <span class="truncate">
                                        {{
                                            state.parentId
                                                ? (parentPages?.find((p) => p.id === state.parentId)?.title ?? $t('common.na'))
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

                    <UFormGroup name="content" :label="$t('common.content')">
                        <ClientOnly>
                            <DocEditor v-model="state.content" />
                        </ClientOnly>
                    </UFormGroup>

                    <div class="mt-2 flex flex-col gap-2">
                        <div class="flex flex-1 gap-2">
                            <UFormGroup name="public" :label="$t('common.public')" class="flex-1">
                                <UToggle v-model="state.meta.public" :disabled="!canDo.public" />
                            </UFormGroup>

                            <UFormGroup name="closed" :label="`${$t('common.toc', 2)}?`" class="flex-1">
                                <UToggle v-model="state.meta.toc" />
                            </UFormGroup>
                        </div>

                        <UFormGroup name="access" :label="$t('common.access')">
                            <AccessManager
                                v-model:jobs="state.access.jobs"
                                v-model:users="state.access.users"
                                :target-id="page.id ?? '0'"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.wiki.AccessLevel')"
                            />
                        </UFormGroup>
                    </div>
                </div>
            </UPage>
        </div>
    </UForm>
</template>
