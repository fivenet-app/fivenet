<script lang="ts" setup>
defineProps<{
    closed: boolean | undefined;
}>();
</script>

<template>
    <UBadge v-if="closed" color="red" class="inline-flex gap-1" size="md">
        <UIcon name="i-mdi-lock" color="red" class="size-5" />
        <span>
            {{ $t('common.close', 2) }}
        </span>
    </UBadge>
    <UBadge v-else color="green" class="inline-flex gap-1" size="md">
        <UIcon name="i-mdi-lock-open-variant" color="green" class="size-5" />
        <span>
            {{ $t('common.open', 2) }}
        </span>
    </UBadge>
</template>
