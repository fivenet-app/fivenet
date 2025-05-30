<script lang="ts" setup>
import type { Group } from '#ui/types';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const appConfig = useAppConfig();

const isOpen = ref(false);

const groups = [
    {
        key: 'pages',
        label: (q: string | undefined) => q && `${t('common.search')}: ${q}`,
        search: async (q: string) => {
            try {
                const call = $grpc.wiki.wiki.listPages({
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
] as Group[];
</script>

<template>
    <UButton
        class="w-full"
        :icon="appConfig.ui.icons.search"
        color="gray"
        :label="$t('common.search_field')"
        truncate
        aria-label="Search"
        v-bind="$attrs"
        @click="isOpen = !isOpen"
    />

    <ClientOnly>
        <UDashboardSearch
            v-model="isOpen"
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
