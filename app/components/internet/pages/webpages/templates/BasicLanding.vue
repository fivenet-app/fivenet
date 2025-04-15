<script lang="ts" setup>
defineProps<{
    title: string;
}>();
</script>

<template>
    <UPage>
        <UHeader :title="title">
            <UNavigationMenu :items="items" />
        </UHeader>

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
