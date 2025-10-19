<script lang="ts" setup>
import type { CommandPaletteGroup, CommandPaletteItem } from '@nuxt/ui';
import { getWikiWikiClient } from '~~/gen/ts/clients';

const appConfig = useAppConfig();

const { t } = useI18n();

const isOpen = ref(false);

const wikiWikiClient = await getWikiWikiClient();

const searchTerm = ref('');
const searchTermDebounced = debouncedRef(searchTerm, 200);

const { data: pages, status } = useLazyAsyncData(
    `mailer-threads-search-${searchTermDebounced.value}`,
    () => listPages(searchTerm.value),
    {
        watch: [searchTermDebounced],
    },
);

async function listPages(q: string): Promise<CommandPaletteItem[]> {
    if (q.length < 3) return [];

    try {
        const call = wikiWikiClient.listPages({
            pagination: {
                offset: 0,
                pageSize: 6,
            },
            rootOnly: false,
            search: q.trim().substring(0, 64),
        });
        const { response } = await call;

        return response.pages.flatMap((page) => ({
            id: page.id,
            label: page.title,
            suffix: `${page.description} ${page.jobLabel}`,
            to: `/wiki/${page.job}/${page.id}/${page.slug ?? ''}`,
            onSelect: () => {
                // Close the search modal when selecting an item
                isOpen.value = false;
            },
        }));
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const groups: CommandPaletteGroup[] = [
    {
        id: 'pages',
        label: t('common.search'),
        items: pages.value || [],
    },
];
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
