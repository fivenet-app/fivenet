<script lang="ts" setup>
import type { NavItem } from '@nuxt/content/dist/runtime/types';
import DocsHeader from '~/components/docs/DocsHeader.vue';
import PageFooter from '~/components/partials/PageFooter.vue';

const { data: navigation } = await useAsyncData<NavItem[]>('navigation', () => fetchContentNavigation(), { default: () => [] });

const { data: files } = useLazyFetch<{ id: string; title: string; content: string }[]>('/help/_content/search', {
    default: () => [],
    server: false,
});

const searchFiles = computed(() =>
    files.value?.map((file) => ({
        _id: file.id,
        title: file.title,
        body: null,
    })),
);
</script>

<template>
    <div>
        <DocsHeader />

        <UMain>
            <UContainer>
                <UPage>
                    <template #left>
                        <UAside>
                            <UNavigationTree :links="mapContentNavigation(!navigation ? [] : navigation[0].children ?? [])" />
                        </UAside>
                    </template>

                    <slot />
                </UPage>

                <ClientOnly>
                    <LazyUContentSearch
                        :files="searchFiles"
                        :navigation="navigation"
                        :empty-state="{
                            icon: 'i-mdi-globe-model',
                            label: $t('commandpalette.empty.title'),
                            queryLabel: $t('commandpalette.empty.title'),
                        }"
                        :placeholder="`${$t('common.search')}...`"
                    />
                </ClientOnly>
            </UContainer>
        </UMain>

        <PageFooter />
    </div>
</template>
