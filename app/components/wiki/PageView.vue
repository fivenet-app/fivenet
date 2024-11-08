<script lang="ts" setup>
import type { TocLink } from '@nuxt/content';
import { emojiBlast } from 'emoji-blast';
import { useNotificatorStore } from '~/store/notificator';
import slug from '~/utils/slugify';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { AccessLevel } from '~~/gen/ts/resources/wiki/access';
import type { Page, PageShort } from '~~/gen/ts/resources/wiki/page';
import ConfirmModal from '../partials/ConfirmModal.vue';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
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

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const notifications = useNotificatorStore();

const breadcrumbs = computed(() => [
    {
        label: t('common.wiki'),
        icon: 'i-mdi-brain',
        to: '/wiki',
    },
    ...[
        !props.page ? { label: t('pages.notfound.page_not_found') } : undefined,
        props.page && props.page?.id !== props.pages?.at(0)?.id ? { label: '...' } : undefined,
        props.page?.meta
            ? { label: props.page.meta.title, to: `/wiki/${props.page.job}/${props.page.id}/${props.page.meta.slug}` }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
]);

async function deletePage(id: string): Promise<void> {
    try {
        const call = getGRPCWikiClient().deletePage({
            id: id,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo({ name: 'wiki' });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function walk(nodes: ChildNode[]) {
    const headers: TocLink[] = [];
    nodes.forEach((n, idx) => {
        const node = n as HTMLElement;

        const sub = Array.from(node.childNodes);
        if (sub.length) {
            walk(sub);
        }

        if (/h[1-6]/i.test(node.tagName)) {
            if (node.id === '') {
                node.id = slug(node.textContent?.substring(0, 64) ?? `${node.tagName}-${idx}`);
            }
            headers.push({
                id: node.id,
                depth: parseInt(node.tagName.replace('H', '')),
                text: node.innerText,
            });
        }
    });
    return headers;
}

const route = useRoute();
watch(
    () => route.hash,
    () => handleHashChange(route.hash),
);

const handleHashChange = (hash: string) => {
    if (hash) {
        const targetElement = document.querySelector(hash); // Find element with id
        if (targetElement) {
            targetElement.scrollIntoView({
                behavior: 'smooth',
            }); // Smooth scroll to the element
        }
    }
};

const contentRef = useTemplateRef('contentRef');
const tocLinks = computedAsync(async () => walk(Array.from(Array.from(contentRef.value?.childNodes ?? []))));
watchOnce(tocLinks, async () => {
    await nextTick();
    handleHashChange(route.hash);
});

const accordionItems = computed(() =>
    [
        { slot: 'access', label: t('common.access'), icon: 'i-mdi-lock' },
        can('WikiService.ListPageActivity').value
            ? { slot: 'activity', label: t('common.activity'), icon: 'i-mdi-comment-quote' }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);
</script>

<template>
    <UDashboardNavbar :title="`${page?.jobLabel ? page?.jobLabel + ': ' : ''}${$t('common.wiki')}`">
        <template #center>
            <PageSearch />
        </template>

        <template #right>
            <UButton v-if="can('WikiService.CreatePage')" color="gray" trailing-icon="i-mdi-plus" to="/wiki/create">
                {{ $t('common.page') }}
            </UButton>
        </template>
    </UDashboardNavbar>

    <div class="flex flex-1 flex-col px-8 py-2 pt-4">
        <UPage>
            <template #left>
                <slot name="left" />
            </template>

            <UBreadcrumb class="pb-2 pt-4" :links="breadcrumbs" />

            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.page')])" />
            <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.page')])" :retry="refresh" />
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
                <UPageHeader
                    v-if="page?.meta"
                    :title="page.meta.title"
                    :description="page.meta.description"
                    :ui="{ wrapper: 'py-4' }"
                >
                    <template #links>
                        <UTooltip :text="$t('common.refresh')">
                            <UButton variant="link" icon="i-mdi-refresh" @click="refresh()" />
                        </UTooltip>

                        <UTooltip v-if="can('WikiService.CreatePage').value" :text="$t('common.edit')">
                            <UButton color="white" icon="i-mdi-pencil" @click="$emit('edit')" />
                        </UTooltip>

                        <UTooltip v-if="can('WikiService.DeletePage').value" :text="$t('common.delete')">
                            <UButton
                                color="red"
                                icon="i-mdi-trash-can"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () => page && deletePage(page.id),
                                    })
                                "
                            />
                        </UTooltip>
                    </template>
                </UPageHeader>

                <UPageBody prose class="pb-8">
                    <!-- eslint-disable vue/no-v-html -->
                    <div ref="contentRef" class="prose dark:prose-invert" v-html="page.content" />
                </UPageBody>

                <UDivider class="mb-4" />

                <UAccordion multiple :items="accordionItems" :unmount="true" class="print:hidden">
                    <template #access>
                        <UContainer>
                            <DataNoDataBlock
                                v-if="!page.access || (page.access?.jobs.length === 0 && page.access?.users.length === 0)"
                                icon="i-mdi-file-search"
                                :message="$t('common.not_found', [$t('common.access', 2)])"
                            />

                            <div v-else class="flex flex-col gap-2">
                                <div class="flex flex-row flex-wrap gap-1">
                                    <UBadge
                                        v-for="entry in page.access?.jobs"
                                        :key="entry.id"
                                        color="black"
                                        class="inline-flex gap-1"
                                        size="md"
                                    >
                                        <span class="size-2 rounded-full bg-info-500" />
                                        <span>
                                            {{ entry.jobLabel
                                            }}<span
                                                v-if="entry.minimumGrade > 0"
                                                :title="`${entry.jobLabel} - ${$t('common.rank')} ${entry.minimumGrade}`"
                                            >
                                                ({{ entry.jobGradeLabel }})</span
                                            >
                                            -
                                            {{ $t(`enums.docstore.AccessLevel.${AccessLevel[entry.access]}`) }}
                                        </span>
                                    </UBadge>
                                </div>

                                <div class="flex flex-row flex-wrap gap-1">
                                    <UBadge
                                        v-for="entry in page.access?.users"
                                        :key="entry.id"
                                        color="black"
                                        class="inline-flex gap-1"
                                        size="md"
                                    >
                                        <span class="size-2 rounded-full bg-amber-500" />
                                        <span :title="`${$t('common.id')} ${entry.userId}`">
                                            {{ entry.user?.firstname }}
                                            {{ entry.user?.lastname }} -
                                            {{ $t(`enums.docstore.AccessLevel.${AccessLevel[entry.access]}`) }}
                                        </span>
                                    </UBadge>
                                </div>
                            </div>
                        </UContainer>
                    </template>

                    <template v-if="can('WikiService.ListPageActivity').value" #activity>
                        <UContainer>
                            <PageActivityList :page-id="page.id" />
                        </UContainer>
                    </template>
                </UAccordion>
            </template>

            <template v-if="page?.meta?.toc === undefined || page?.meta?.toc === true" #right>
                <PageSearch class="!flex lg:!hidden" />

                <UContentToc :title="$t('common.toc')" :links="tocLinks" />
            </template>
        </UPage>
    </div>
</template>
