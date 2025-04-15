<script lang="ts" setup>
defineProps<{
    title: string;
}>();
</script>

<template>
    <UPage>
        <UHeader :title="title" :items="[]" />

        <UMain>
            <slot />
        </UMain>

        <UFooter>
            <template #top>
                <UFooterColumns :items="[]" />
            </template>
        </UFooter>
    </UPage>
</template>
