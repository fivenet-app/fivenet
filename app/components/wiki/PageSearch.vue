<script lang="ts" setup>
import type { CommandPaletteGroup, CommandPaletteItem } from '@nuxt/ui';
import { computed } from 'vue';
import { pageToURL } from './helpers';

const appConfig = useAppConfig();

const { t } = useI18n();

const isOpen = ref<boolean>(false);

const { listPages: listWikiPages } = await useWikiWiki();

const searchTerm = ref('');
const searchTermDebounced = debouncedRef(searchTerm, 200);

const pagesKey = computed(() => `wiki-pages-search-${searchTermDebounced.value}`);

const { data: pages, status } = useLazyAsyncData(pagesKey, () => listPages(searchTerm.value));

async function listPages(q: string): Promise<CommandPaletteItem[]> {
    if (q.length < 3) return [];

    const response = await listWikiPages({
        pagination: {
            offset: 0,
            pageSize: 6,
        },
        rootOnly: false,
        search: q.trim().substring(0, 64),
    });

    return response.pages.flatMap((page) => ({
        id: page.id,
        label: page.title,
        suffix: `${page.description} ${page.jobLabel}`,
        to: pageToURL(page),
        onSelect: () => {
            // Close the search modal when selecting an item
            isOpen.value = false;
        },
    }));
}

const groups = computed<CommandPaletteGroup[]>(() => [
    {
        id: 'pages',
        label: t('common.search'),
        items: pages.value || [],
    },
]);
</script>

<template>
    <UButton
        class="w-full"
        :icon="appConfig.ui.icons.search"
        color="neutral"
        variant="outline"
        :label="$t('common.search_field')"
        truncate
        aria-label="Search"
        v-bind="$attrs"
        @click="isOpen = !isOpen"
    />

    <UModal v-model:open="isOpen">
        <template #content>
            <ClientOnly>
                <UCommandPalette
                    v-model:search-term="searchTerm"
                    :loading="isRequestPending(status)"
                    :color-mode="false"
                    :groups="groups"
                    :placeholder="`${$t('common.search_field')}`"
                    :fuse="{ resultLimit: 6, fuseOptions: { threshold: 0.1 } }"
                >
                    <template #empty>
                        {{ $t('commandpalette.empty.title') }}
                    </template>
                </UCommandPalette>
            </ClientOnly>
        </template>
    </UModal>
</template>
