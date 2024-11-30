<script lang="ts" setup>
defineProps<{
    title: string;
}>();
</script>

<template>
    <UPage>
        <UHeader :title="title" :links="[]" />

        <UMain>
            <slot />
        </UMain>

        <UFooter>
            <template #top>
                <UFooterColumns :links="[]" />
            </template>
        </UFooter>
    </UPage>
</template>
