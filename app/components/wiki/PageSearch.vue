<script lang="ts" setup>
import type { Group } from '#ui/types';

const { t } = useI18n();

const appConfig = useAppConfig();

const isOpen = ref(false);

const groups = [
    {
        key: 'pages',
        label: (q: string | undefined) => q && `${t('common.search')}: ${q}`,
        search: async (q: string) => {
            try {
                const call = getGRPCWikiClient().listPages({
                    pagination: {
                        offset: 0,
                    },
                    rootOnly: false,
                    search: q.trim(),
                });
                const { response } = await call;

                return response.pages.flatMap((page) => ({
                    id: page.id,
                    label: page.title,
                    suffix: `${page.description} ${page.jobLabel}`,
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
        v-bind="$attrs"
        :icon="appConfig.ui.icons.search"
        color="gray"
        class="w-full min-w-40 sm:min-w-64"
        :label="$t('common.search_field')"
        truncate
        aria-label="Search"
        @click="isOpen = !isOpen"
    />

    <ClientOnly>
        <UModal v-model="isOpen" :ui="{ width: 'w-full sm:max-w-5xl' }">
            <UCommandPalette
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
        </UModal>
    </ClientOnly>
</template>
