<script lang="ts" setup>
import type { NavItem } from '@nuxt/content/dist/runtime/types';
import DocsHeader from '~/components/docs/DocsHeader.vue';
import PageFooter from '~/components/partials/PageFooter.vue';

const { data: navigation } = await useAsyncData<NavItem[]>('navigation', () => fetchContentNavigation(), { default: () => [] });
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
            </UContainer>
        </UMain>

        <PageFooter />
    </div>
</template>
