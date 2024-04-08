<script lang="ts" setup>
withDefaults(
    defineProps<{
        doubleHeader?: boolean;
    }>(),
    {
        doubleHeader: false,
    },
);
</script>

<template>
    <table class="min-w-full divide-y divide-base-600">
        <thead>
            <slot name="thead" />
        </thead>
        <tbody class="divide-y divide-base-800">
            <slot name="tbody" />
        </tbody>
        <thead v-if="doubleHeader">
            <slot name="thead" />
        </thead>
    </table>
</template>
