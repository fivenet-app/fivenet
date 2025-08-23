<script lang="ts" setup>
import type { CommandPaletteGroup } from '@nuxt/ui';
import { getWikiWikiClient } from '~~/gen/ts/clients';

const appConfig = useAppConfig();

const { t } = useI18n();

const wikiWikiClient = await getWikiWikiClient();

const isOpen = ref(false);

const groups = [
    {
        id: 'pages',
        label: (q: string | undefined) => q && `${t('common.search')}: ${q}`,
        search: async (q: string) => {
            try {
                const call = wikiWikiClient.listPages({
                    pagination: {
                        offset: 0,
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
                }));
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
    },
] as CommandPaletteGroup[];
</script>

<template>
    <UButton
        class="w-full"
        :icon="appConfig.ui.icons.search"
        color="neutral"
        :label="$t('common.search_field')"
        truncate
        aria-label="Search"
        v-bind="$attrs"
        @click="isOpen = !isOpen"
    />

    <ClientOnly>
        <UDashboardSearch
            v-model:open="isOpen"
            hide-color-mode
            :groups="groups"
            :empty-state="{
                icon: 'i-mdi-brain',
                label: $t('commandpalette.empty.title'),
                queryLabel: $t('commandpalette.empty.title'),
            }"
            :placeholder="`${$t('common.search_field')}`"
            :autoselect="false"
            :fuse="{ resultLimit: 6, fuseOptions: { threshold: 0.1 } }"
        />
    </ClientOnly>
</template>
