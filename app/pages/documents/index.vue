<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import List from '~/components/documents/List.vue';

useHead({
    title: 'pages.documents.title',
});

definePageMeta({
    title: 'pages.documents.title',
    requiresAuth: true,
    permission: 'documents.DocumentsService/ListDocuments',
});

defineProps<{
    itemsLeft: NavigationMenuItem[];
    itemsRight: NavigationMenuItem[];
}>();
</script>

<template>
    <List>
        <template v-if="itemsLeft.length > 1 || itemsRight.length > 1" #footer>
            <USeparator />

            <UDashboardToolbar>
                <div class="flex min-w-0 flex-1 flex-row flex-wrap justify-between gap-2">
                    <UNavigationMenu v-if="itemsLeft.length > 0" class="-ml-1" orientation="horizontal" :items="itemsLeft" />

                    <UNavigationMenu v-if="itemsRight.length > 0" class="-mr-1" orientation="horizontal" :items="itemsRight" />
                </div>
            </UDashboardToolbar>
        </template>
    </List>
</template>
