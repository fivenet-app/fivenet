<script lang="ts" setup>
import List from '~/components/citizens/List.vue';
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';

useHead({
    title: 'pages.citizens.title',
});
definePageMeta({
    title: 'pages.citizens.title',
    requiresAuth: true,
    permission: 'CitizenStoreService.ListCitizens',
});
</script>

<template>
    <div class="w-full">
        <List />
        <ClipboardButton />
    </div>
</template>
