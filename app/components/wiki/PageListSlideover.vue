<script lang="ts" setup>
const emits = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

// TODO
</script>

<template>
    <USlideover :title="$t('common.page', 2)" :ui="{ content: 'max-w-3xl' }"> </USlideover>
</template>
