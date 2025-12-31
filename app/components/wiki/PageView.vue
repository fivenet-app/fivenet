<script lang="ts" setup>
import type { AsyncDataRequestStatus } from '#app';
import type { NavigationMenuItem } from '@nuxt/ui';
import type { ContentSurroundLink } from '@nuxt/ui/runtime/components/content/ContentSurround.vue.js';
import { emojiBlast } from 'emoji-blast';
import AccessBadges from '~/components/partials/access/AccessBadges.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import HTMLContent from '~/components/partials/content/HTMLContent.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { jsonNodeToTocLinks } from '~/utils/content';
import { getWikiWikiClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { AccessLevel } from '~~/gen/ts/resources/wiki/access';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';
import GenericImg from '../partials/elements/GenericImg.vue';
import ScrollToTop from '../partials/ScrollToTop.vue';
import List from './activity/List.vue';
import { checkPageAccess } from './helpers';
import PageSearch from './PageSearch.vue';

const props = defineProps<{
    page: Page | undefined;
    pages: PageShort[];
    navItems?: NavigationMenuItem[];
    status: AsyncDataRequestStatus;
    refresh: () => Promise<void>;
    error: Error | undefined;
}>();

const { t } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const wikiWikiClient = await getWikiWikiClient();

const confirmModal = overlay.create(ConfirmModal);

const breadcrumbs = computed(() => {
    const breadcrumbList: { label: string; icon?: string; to?: string }[] = [
        {
            label: t('common.wiki'),
            icon: 'i-mdi-home',
            to: '/wiki',
        },
    ];

    if (props.page && props.pages) {
        const addBreadcrumbs = (pages: PageShort[], currentPage: Page | PageShort) => {
            for (const page of pages) {
                if (page.id !== 0) {
                    breadcrumbList.push({
                        label: page.title || t('common.untitled'),
                        to: `/wiki/${page.job}/${page.id}/${page.slug}`,
                    });
                }

                if (page.id === currentPage.id) {
                    return true;
                }

                if (page.children && addBreadcrumbs(page.children, currentPage)) {
                    return true;
                }

                breadcrumbList.pop();
            }
            return false;
        };

        addBreadcrumbs(props.pages, props.page);
    } else if (!isRequestPending(props.status) && !props.page) {
        breadcrumbList.push({ label: t('pages.notfound.page_not_found') });
    }

    return breadcrumbList;
});

async function deletePage(id: number): Promise<void> {
    try {
        const call = wikiWikiClient.deletePage({
            id: id,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
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

const wikiService = await useWikiWiki();

const tocLinks = computedAsync(async () => props.page?.content?.content && jsonNodeToTocLinks(props.page?.content?.content));

const canAccessActivity = computed(
    () =>
        can('wiki.WikiService/ListPageActivity').value &&
        checkPageAccess(props.page?.access, props.page?.meta?.creator, AccessLevel.VIEW, props.page?.job),
);

const canAccessFiles = computed(
    () =>
        can('wiki.WikiService/ListPageActivity').value &&
        checkPageAccess(props.page?.access, props.page?.meta?.creator, AccessLevel.ACCESS, props.page?.job),
);

const accordionItems = computed(() =>
    [
        { slot: 'access' as const, label: t('common.access'), icon: 'i-mdi-lock' },
        canAccessActivity.value
            ? { slot: 'activity' as const, label: t('common.activity'), icon: 'i-mdi-comment-quote' }
            : undefined,
        canAccessFiles.value ? { slot: 'files' as const, label: t('common.file', 2), icon: 'i-mdi-file' } : undefined,
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
                  id: prev.id,
                  title: prev.title || '',
                  description: prev.description ?? '',
                  path: `/wiki/${prev.job}/${prev.id}/${prev.slug}`,
              }
            : undefined,
        next
            ? {
                  id: next.id,
                  title: next.title || '',
                  description: next.description ?? '',
                  path: `/wiki/${next.job}/${next.id}/${next.slug}`,
              }
            : undefined,
    ];
}, []);

const scrollRef = useTemplateRef('scrollRef');
</script>

<template>
    <UDashboardPanel :ui="{ body: 'py-0 sm:py-0 lg:py-6 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="`${page?.jobLabel ? page?.jobLabel + ': ' : ''}${$t('common.wiki')}`">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #default>
                    <PageSearch />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/wiki" />

                    <UButton
                        v-if="can('wiki.WikiService/UpdatePage').value"
                        color="neutral"
                        trailing-icon="i-mdi-plus"
                        @click="wikiService.createPage(page?.parentId ?? page?.id)"
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('common.page') }}
                        </span>
                    </UButton>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar class="flex lg:hidden">
                <template #default>
                    <PageSearch />
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <UPage ref="scrollRef" :ui="{ left: 'lg:top-0', root: 'lg:gap-4' }">
                <template #left>
                    <slot name="left" />
                </template>

                <UPage>
                    <UNavigationMenu class="mt-4 lg:hidden" :items="navItems" orientation="vertical" />

                    <UBreadcrumb class="pt-4 lg:pt-0" :items="breadcrumbs" />

                    <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.page')])" />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.page')])"
                        :error="error"
                        :retry="refresh"
                    />

                    <template v-else-if="!page">
                        <UPageHero
                            :title="$t('pages.notfound.page_not_found')"
                            :description="$t('pages.notfound.fun_error')"
                            :links="[
                                {
                                    label: $t('common.back'),
                                    icon: 'i-mdi-arrow-back',
                                    size: 'md',
                                    color: 'neutral',
                                    onClick: () => useRouter().back(),
                                },
                                { label: $t('common.wiki'), icon: 'i-mdi-home', size: 'md', to: '/wiki' },
                            ]"
                            :ui="{ title: 'text-3xl sm:text-4xl' }"
                        >
                            <template #headline>
                                <UBadge
                                    color="neutral"
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
                        </UPageHero>
                    </template>

                    <template v-else>
                        <UPageHeader
                            v-if="page?.meta"
                            :title="!page.meta.title ? $t('common.untitled') : page.meta.title"
                            :ui="{ root: 'py-4', wrapper: 'py-4', title: !page.meta.title ? 'italic' : '' }"
                        >
                            <template #links>
                                <UTooltip :text="$t('common.refresh')">
                                    <UButton variant="link" icon="i-mdi-refresh" @click="refresh()" />
                                </UTooltip>

                                <UTooltip
                                    v-if="
                                        can('wiki.WikiService/UpdatePage').value &&
                                        checkPageAccess(page.access, page.meta.creator, AccessLevel.EDIT, page?.job)
                                    "
                                    :text="$t('common.edit')"
                                >
                                    <UButton
                                        color="neutral"
                                        icon="i-mdi-pencil"
                                        :to="`/wiki/${page.job}/${page.id}/${page.meta.slug ?? ''}/edit`"
                                    />
                                </UTooltip>

                                <UTooltip
                                    v-if="
                                        can('wiki.WikiService/DeletePage').value &&
                                        checkPageAccess(page.access, page.meta.creator, AccessLevel.EDIT, page?.job)
                                    "
                                    :text="!page.meta.deletedAt ? $t('common.delete') : $t('common.restore')"
                                >
                                    <UButton
                                        :color="!page.meta.deletedAt ? 'error' : 'success'"
                                        :icon="!page.meta.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                                        @click="
                                            confirmModal.open({
                                                confirm: async () => page && deletePage(page.id),
                                            })
                                        "
                                    />
                                </UTooltip>
                            </template>

                            <template v-if="page.meta.updatedAt || page.meta.deletedAt" #description>
                                <div class="flex snap-x flex-row flex-wrap gap-2 overflow-x-auto pb-3 sm:pb-0">
                                    <UBadge
                                        v-if="page.meta.createdAt"
                                        class="inline-flex gap-1"
                                        color="neutral"
                                        icon="i-mdi-calendar"
                                        size="md"
                                    >
                                        {{ $t('common.created') }}
                                        <GenericTime :value="page.meta.createdAt" type="long" />
                                    </UBadge>

                                    <UBadge
                                        v-if="page.meta.updatedAt"
                                        class="inline-flex gap-1"
                                        color="neutral"
                                        icon="i-mdi-calendar-edit"
                                        size="md"
                                    >
                                        {{ $t('common.updated') }}
                                        <GenericTime :value="page.meta.updatedAt" type="long" />
                                    </UBadge>

                                    <UBadge
                                        v-if="page.meta.deletedAt"
                                        class="inline-flex gap-1"
                                        color="warning"
                                        icon="i-mdi-calendar-remove"
                                        size="md"
                                    >
                                        {{ $t('common.deleted') }}
                                        <GenericTime :value="page.meta.deletedAt" type="long" />
                                    </UBadge>

                                    <UBadge
                                        v-if="page.meta.draft"
                                        class="inline-flex gap-1"
                                        color="info"
                                        icon="i-mdi-pencil"
                                        size="md"
                                        :label="$t('common.draft')"
                                    />

                                    <UBadge
                                        v-if="page.meta.public"
                                        class="inline-flex gap-1"
                                        color="neutral"
                                        icon="i-mdi-earth"
                                        :label="$t('common.public')"
                                        size="md"
                                    />

                                    <UBadge
                                        v-if="page.meta.startpage"
                                        class="inline-flex gap-1"
                                        color="neutral"
                                        icon="i-mdi-home"
                                        size="md"
                                        :label="$t('common.startpage')"
                                    />
                                </div>

                                <p v-if="page.meta.description" class="mt-4">{{ page.meta.description }}</p>
                            </template>
                        </UPageHeader>

                        <UPageBody v-if="page.content?.content">
                            <div
                                class="mx-auto w-full max-w-(--breakpoint-xl) rounded-lg bg-neutral-100 p-4 break-words dark:bg-neutral-800"
                            >
                                <HTMLContent :value="page.content.content" />
                            </div>

                            <template v-if="surround.filter((s) => s !== undefined).length > 0">
                                <USeparator class="my-2" />

                                <UContentSurround
                                    :surround="surround as ContentSurroundLink[]"
                                    prev-icon="i-mdi-arrow-left"
                                    next-icon="i-mdi-arrow-right"
                                />
                            </template>

                            <USeparator class="my-2" />

                            <UAccordion class="print:hidden" :items="accordionItems" type="multiple" :unmount-on-hide="false">
                                <template #access>
                                    <UContainer class="mb-2">
                                        <DataNoDataBlock
                                            v-if="
                                                !page.access ||
                                                (page.access?.jobs.length === 0 && page.access?.users.length === 0)
                                            "
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

                                <template v-if="canAccessActivity" #activity>
                                    <UContainer class="mb-2">
                                        <List :page-id="page.id" />
                                    </UContainer>
                                </template>

                                <template v-if="canAccessFiles" #files>
                                    <DataNoDataBlock
                                        v-if="!page.files || page.files.length === 0"
                                        icon="i-mdi-file-search"
                                        :message="$t('common.not_found', [$t('common.file', 2)])"
                                    />
                                    <UContainer v-else class="p-2">
                                        <UPageGrid class="flex-1 sm:grid-cols-1 lg:grid-cols-1 xl:grid-cols-2">
                                            <UPageCard
                                                v-for="file in page.files"
                                                :key="file.id"
                                                :title="file.filePath"
                                                icon="i-mdi-file-document"
                                                orientation="horizontal"
                                                :ui="{ title: 'line-clamp-3! whitespace-normal!' }"
                                            >
                                                <template #default>
                                                    <div class="inline-flex items-center justify-center">
                                                        <GenericImg
                                                            v-if="file.contentType.startsWith('image/')"
                                                            :src="file.filePath"
                                                            :alt="file.filePath"
                                                            size="3xl"
                                                            class="h-full max-h-40 w-40"
                                                        />
                                                        <UIcon
                                                            v-else
                                                            class="h-20 w-20 text-3xl"
                                                            :name="
                                                                file.contentType.startsWith('video/')
                                                                    ? 'i-mdi-video'
                                                                    : 'i-mdi-file-document'
                                                            "
                                                        />
                                                    </div>
                                                </template>

                                                <template #description>
                                                    <ul>
                                                        <li>{{ file.contentType }}</li>
                                                        <li>{{ formatBytes(file.byteSize) }}</li>
                                                    </ul>
                                                </template>
                                            </UPageCard>
                                        </UPageGrid>
                                    </UContainer>
                                </template>
                            </UAccordion>
                        </UPageBody>
                    </template>

                    <template
                        v-if="(page?.meta?.toc === undefined || page?.meta?.toc === true) && tocLinks && tocLinks?.length > 0"
                        #right
                    >
                        <UContentToc
                            class="lg:col-span-2"
                            :title="$t('common.toc')"
                            :links="tocLinks"
                            :ui="{ root: 'top-0' }"
                        />
                    </template>
                </UPage>
            </UPage>

            <ScrollToTop :element="scrollRef?.$el" />
        </template>
    </UDashboardPanel>
</template>
