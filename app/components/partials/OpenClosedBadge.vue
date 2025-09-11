<script lang="ts" setup>
defineProps<{
    closed: boolean | undefined;
}>();
</script>

<template>
    <UBadge v-if="closed" color="error" icon="i-mdi-lock" :label="$t('common.close', 2)" v-bind="$attrs" />
    <UBadge v-else color="green" icon="i-mdi-lock-open-variant" :label="$t('common.open', 2)" v-bind="$attrs" />
</template>
