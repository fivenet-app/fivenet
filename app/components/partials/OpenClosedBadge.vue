<script lang="ts" setup>
defineProps<{
    closed: boolean | undefined;
}>();
</script>

<template>
    <UBadge v-if="closed" v-bind="$attrs" color="error" class="inline-flex gap-1">
        <UIcon name="i-mdi-lock" color="error" class="size-5" />
        <span>
            {{ $t('common.close', 2) }}
        </span>
    </UBadge>
    <UBadge v-else v-bind="$attrs" color="green" class="inline-flex gap-1">
        <UIcon name="i-mdi-lock-open-variant" color="green" class="size-5" />
        <span>
            {{ $t('common.open', 2) }}
        </span>
    </UBadge>
</template>
