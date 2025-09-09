<script lang="ts" setup>
import { sexToColor, sexToIcon } from './helpers';

defineProps<{
    sex: string;
}>();
</script>

<template>
    <UBadge :color="sexToColor(sex)" size="lg" :icon="sexToIcon(sex)" />
</template>
