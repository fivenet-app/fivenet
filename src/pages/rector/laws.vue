<script lang="ts" setup>
import Laws from '~/components/rector/Laws.vue';

useHead({
    title: 'pages.rector.laws.title',
});
definePageMeta({
    title: 'pages.rector.laws.title',
    requiresAuth: true,
    permission: 'SuperUser',
    showQuickButtons: false,
});
</script>

<template>
    <div>
        <Laws />
    </div>
</template>
