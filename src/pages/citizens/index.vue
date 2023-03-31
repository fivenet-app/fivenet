<script lang="ts" setup>
import CitizensList from '../../components/citizens/CitizensList.vue';
import ClipboardButton from '../../components/clipboard/ClipboardButton.vue';

useHead({
    title: 'Citizens',
});
definePageMeta({
    title: 'Citizens',
    requiresAuth: true,
    permission: 'CitizenStoreService.FindUsers',
});
</script>

<template>
    <div class="w-full">
        <CitizensList />
        <ClipboardButton />
    </div>
</template>
