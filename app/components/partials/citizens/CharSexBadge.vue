<script lang="ts" setup>
import type { BadgeColor } from '#ui/types';
import { sexToColor, sexToIcon } from './helpers';

defineProps<{
    sex: string;
}>();
</script>

<template>
    <UBadge :color="sexToColor(sex) as BadgeColor" size="lg">
        <UIcon :name="sexToIcon(sex)" />
    </UBadge>
</template>
