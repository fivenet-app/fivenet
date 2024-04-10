<script setup lang="ts">
import type { NavItem } from '@nuxt/content/dist/runtime/types';

const navigation = inject<Ref<NavItem[]>>('navigation');
</script>

<template>
    <div>
        <DocsHeader />

        <UContainer>
            <UPage>
                <template #left>
                    <UAside>
                        <UNavigationTree :links="mapContentNavigation(navigation)" />
                    </UAside>
                </template>

                <slot />
            </UPage>
        </UContainer>
    </div>
</template>
