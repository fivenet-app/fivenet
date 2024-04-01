<script lang="ts" setup>
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import DocumentList from '~/components/documents/DocumentList.vue';

useHead({
    title: 'pages.documents.title',
});
definePageMeta({
    title: 'pages.documents.title',
    requiresAuth: true,
    permission: 'DocStoreService.ListDocuments',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar title="Documents">
                <template #right>
                    <UInput
                        ref="input"
                        icon="i-heroicons-funnel"
                        autocomplete="off"
                        placeholder="Filter users..."
                        class="hidden lg:block"
                        @keydown.esc="$event.target.blur()"
                    >
                        <template #trailing>
                            <UKbd value="/" />
                        </template>
                    </UInput>

                    <UButton label="New user" trailing-icon="i-heroicons-plus" color="gray" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <template #left> </template>

                <template #right> </template>
            </UDashboardToolbar>

            <DocumentList />
            <ClipboardButton />
        </UDashboardPanel>
    </UDashboardPage>
</template>
