<script lang="ts" setup>
defineProps<{
    closed: boolean | undefined;
}>();
</script>

<template>
    <UBadge v-if="closed" class="inline-flex gap-1" color="error" v-bind="$attrs">
        <UIcon class="size-5" name="i-mdi-lock" color="error" />
        <span>
            {{ $t('common.close', 2) }}
        </span>
    </UBadge>
    <UBadge v-else class="inline-flex gap-1" color="green" v-bind="$attrs">
        <UIcon class="size-5" name="i-mdi-lock-open-variant" color="green" />
        <span>
            {{ $t('common.open', 2) }}
        </span>
    </UBadge>
</template>
