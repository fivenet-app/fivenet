<script lang="ts" setup>
import Header from '~/components/docs/Header.vue';
import Footer from '~/components/docs/Footer.vue';

const { data: navigation } = await useAsyncData('navigation', () => fetchContentNavigation(), { default: () => [] });

const { data: files } = useLazyFetch<{ id: string; title: string; content: string }[]>('/docs/_content/search', {
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
        <Header />

        <UMain>
            <UContainer>
                <UPage>
                    <template #left>
                        <UAside>
                            <UNavigationTree :links="mapContentNavigation(!navigation ? [] : navigation)" />
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

        <Footer />
    </div>
</template>
