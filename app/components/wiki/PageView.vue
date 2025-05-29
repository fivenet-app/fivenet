<script lang="ts" setup>
import { emojiBlast } from 'emoji-blast';
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import HTMLContent from '~/components/partials/content/HTMLContent.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useNotificatorStore } from '~/stores/notificator';
import { jsonNodeToTocLinks } from '~/utils/content';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { AccessLevel } from '~~/gen/ts/resources/wiki/access';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';
import ScrollToTop from '../partials/ScrollToTop.vue';
import { checkPageAccess } from './helpers';
import PageActivityList from './PageActivityList.vue';
import PageSearch from './PageSearch.vue';

const props = defineProps<{
    page: Page | undefined;
    pages: PageShort[];
    loading: boolean;
    refresh: () => Promise<void>;
    error: Error | undefined;
}>();

defineEmits<{
    (e: 'edit'): void;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const notifications = useNotificatorStore();

const breadcrumbs = computed(() => [
    {
        label: t('common.wiki'),
        icon: 'i-mdi-home',
        to: '/wiki',
    },
    ...[
        !props.page && !props.loading ? { label: t('pages.notfound.page_not_found') } : undefined,
        props.page && props.page?.id !== props.pages?.at(0)?.id ? { label: '...' } : undefined,
        props.page?.meta
            ? { label: props.page.meta.title, to: `/wiki/${props.page.job}/${props.page.id}/${props.page.meta.slug}` }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
]);

async function deletePage(id: number): Promise<void> {
    try {
        const call = $grpc.wiki.wiki.deletePage({
            id: id,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        // If the deleted page is not the "top page", navigate to it
        if (props.pages[0] && props.page?.id !== props.pages[0].id) {
            await navigateTo({
                name: 'wiki-job-id-slug',
                params: {
                    job: props.pages[0].job,
                    id: props.pages[0].id,
                    slug: [props.pages[0].slug ?? ''],
                },
            });
            return;
        }

        await navigateTo({ name: 'wiki' });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const tocLinks = computedAsync(async () => props.page?.content?.content && jsonNodeToTocLinks(props.page?.content?.content));

const accordionItems = computed(() =>
    [
        { slot: 'access', label: t('common.access'), icon: 'i-mdi-lock' },
        can('wiki.WikiService.ListPageActivity').value &&
        checkPageAccess(props.page?.access, props.page?.meta?.creator, AccessLevel.VIEW)
            ? { slot: 'activity', label: t('common.activity'), icon: 'i-mdi-comment-quote' }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

async function findSurroundingPages(
    pages: PageShort[],
    currentPage: Page | undefined,
): Promise<{ prev: PageShort | undefined; next: PageShort | undefined }> {
    if (!currentPage) return { prev: undefined, next: undefined };

    const flatPages: PageShort[] = [];
    function flattenPages(pages: PageShort[], level = 0) {
        for (const page of pages) {
            if (page.children[0] && page.children[0].id === page.id) {
                if (page.children) {
                    flattenPages(page.children, level + 1);
                }
                continue;
            }

            if (level > 0) {
                flatPages.push({ ...page, level: level });
            }

            if (page.children) {
                flattenPages(page.children, level + 1);
            }
        }
    }

    flattenPages(pages);

    const currentIndex = flatPages.findIndex((p) => p.id === currentPage.id);
    const prev = currentIndex > 0 ? flatPages[currentIndex - 1] : undefined;
    const next = currentIndex >= 0 && currentIndex < flatPages.length - 1 ? flatPages[currentIndex + 1] : undefined;

    return { prev, next };
}

const surround = computedAsync(async () => {
    const { prev, next } = await findSurroundingPages(props.pages, props.page);
    return [
        prev
            ? {
                  _id: prev.id,
                  title: prev.title || '',
                  description: prev.description ?? '',
                  _path: `/wiki/${prev.job}/${prev.id}/${prev.slug}`,
              }
            : undefined,
        next
            ? {
                  _id: next.id,
                  title: next.title || '',
                  description: next.description ?? '',
                  _path: `/wiki/${next.job}/${next.id}/${next.slug}`,
              }
            : undefined,
    ];
}, []);

const prev = computed(() => surround.value[0]);
const next = computed(() => surround.value[1]);

const scrollRef = useTemplateRef('scrollRef');
</script>

<template>
    <UDashboardNavbar :title="`${page?.jobLabel ? page?.jobLabel + ': ' : ''}${$t('common.wiki')}`">
        <template #center>
            <PageSearch />
        </template>

        <template #right>
            <PartialsBackButton fallback-to="/wiki" />

            <UButton v-if="can('wiki.WikiService.CreatePage').value" color="gray" trailing-icon="i-mdi-plus" to="/wiki/create">
                {{ $t('common.page') }}
            </UButton>
        </template>
    </UDashboardNavbar>

    <UDashboardPanelContent ref="scrollRef" class="p-0 sm:pb-0">
        <UPage class="px-8 py-2 pt-4">
            <template #left>
                <slot name="left" />
            </template>

            <UBreadcrumb class="pb-2 pt-4" :links="breadcrumbs" />

            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.page')])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.page')])"
                :error="error"
                :retry="refresh"
            />
            <template v-else-if="!page">
                <ULandingHero
                    :title="$t('pages.notfound.page_not_found')"
                    :description="$t('pages.notfound.fun_error')"
                    :links="[
                        {
                            label: $t('common.back'),
                            icon: 'i-mdi-arrow-back',
                            size: 'md',
                            color: 'gray',
                            click: () => useRouter().back(),
                        },
                        { label: $t('common.wiki'), icon: 'i-mdi-home', size: 'md', to: '/wiki' },
                    ]"
                    :ui="{ title: 'text-3xl sm:text-4xl' }"
                >
                    <template #headline>
                        <UBadge
                            color="gray"
                            variant="solid"
                            size="lg"
                            @click="
                                emojiBlast({
                                    emojis: ['😵‍💫', '🔍', '🔎', '👀'],
                                })
                            "
                            >{{ $t('pages.notfound.error') }}</UBadge
                        >
                    </template>
                </ULandingHero>
            </template>

            <template v-else>
                <UPageHeader v-if="page?.meta" :title="page.meta.title" :ui="{ wrapper: 'py-4' }">
                    <template #links>
                        <UTooltip :text="$t('common.refresh')">
                            <UButton variant="link" icon="i-mdi-refresh" @click="refresh()" />
                        </UTooltip>

                        <UTooltip
                            v-if="
                                can('wiki.WikiService.CreatePage').value &&
                                checkPageAccess(page.access, page.meta.creator, AccessLevel.EDIT)
                            "
                            :text="$t('common.edit')"
                        >
                            <UButton color="white" icon="i-mdi-pencil" @click="$emit('edit')" />
                        </UTooltip>

                        <UTooltip
                            v-if="
                                can('wiki.WikiService.DeletePage').value &&
                                checkPageAccess(page.access, page.meta.creator, AccessLevel.EDIT)
                            "
                            :text="!page.meta.deletedAt ? $t('common.delete') : $t('common.restore')"
                        >
                            <UButton
                                :color="!page.meta.deletedAt ? 'error' : 'success'"
                                :icon="!page.meta.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () => page && deletePage(page.id),
                                    })
                                "
                            />
                        </UTooltip>
                    </template>

                    <template v-if="page.meta.updatedAt || page.meta.deletedAt" #description>
                        <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-0">
                            <UBadge v-if="page.meta.createdAt" class="inline-flex gap-1" color="black" size="md">
                                <UIcon class="size-5" name="i-mdi-calendar" />
                                <span>
                                    {{ $t('common.created') }}
                                    <GenericTime :value="page.meta.createdAt" type="long" />
                                </span>
                            </UBadge>

                            <UBadge v-if="page.meta.updatedAt" class="inline-flex gap-1" color="black" size="md">
                                <UIcon class="size-5" name="i-mdi-calendar-edit" />
                                <span>
                                    {{ $t('common.updated') }}
                                    <GenericTime :value="page.meta.updatedAt" type="long" />
                                </span>
                            </UBadge>

                            <UBadge v-if="page.meta.deletedAt" class="inline-flex gap-1" color="amber" size="md">
                                <UIcon class="size-5" name="i-mdi-calendar-remove" />
                                <span>
                                    {{ $t('common.deleted') }}
                                    <GenericTime :value="page.meta.deletedAt" type="long" />
                                </span>
                            </UBadge>

                            <UBadge v-if="page.meta.public" class="inline-flex gap-1" color="black" size="md">
                                <UIcon class="size-5" name="i-mdi-earth" />
                                <span>
                                    {{ $t('common.public') }}
                                </span>
                            </UBadge>
                        </div>

                        <p v-if="page.meta.description" class="mt-4">{{ page.meta.description }}</p>
                    </template>
                </UPageHeader>

                <UPageBody v-if="page.content?.content">
                    <div class="rounded-lg bg-neutral-100 dark:bg-base-900">
                        <HTMLContent class="px-4 py-2" :value="page.content.content" />
                    </div>

                    <template v-if="surround.length > 0">
                        <UDivider class="mb-4 mt-4" />

                        <!-- UContentSurround doesn't seem to like our surround pages array -->
                        <div class="grid gap-8 sm:grid-cols-2">
                            <UContentSurroundLink v-if="prev" :link="prev" icon="i-mdi-arrow-left" />
                            <span v-else class="hidden sm:block">&nbsp;</span>
                            <UContentSurroundLink v-if="next" class="text-right" :link="next" icon="i-mdi-arrow-right" />
                        </div>
                    </template>

                    <UDivider class="mb-4 mt-4" />

                    <UAccordion class="print:hidden" multiple :items="accordionItems" :unmount="true">
                        <template #access>
                            <UContainer>
                                <DataNoDataBlock
                                    v-if="!page.access || (page.access?.jobs.length === 0 && page.access?.users.length === 0)"
                                    icon="i-mdi-file-search"
                                    :message="$t('common.not_found', [$t('common.access', 2)])"
                                />

                                <AccessBadges
                                    v-else
                                    :access-level="AccessLevel"
                                    :jobs="page?.access.jobs"
                                    :users="page?.access.users"
                                    i18n-key="enums.wiki"
                                />
                            </UContainer>
                        </template>

                        <template v-if="can('wiki.WikiService.ListPageActivity').value" #activity>
                            <UContainer>
                                <PageActivityList :page-id="page.id" />
                            </UContainer>
                        </template>
                    </UAccordion>
                </UPageBody>
            </template>

            <template v-if="page?.meta?.toc === undefined || page?.meta?.toc === true" #right>
                <PageSearch class="mb-2 !flex lg:!hidden" />

                <UContentToc :title="$t('common.toc')" :links="tocLinks" />
            </template>
        </UPage>

        <ScrollToTop :element="scrollRef?.$el" />
    </UDashboardPanelContent>
</template>
