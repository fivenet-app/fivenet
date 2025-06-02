<script lang="ts" setup>
import { sexToColor, sexToIcon } from './helpers';

defineProps<{
    sex: string;
}>();
</script>

<template>
    <UBadge :color="sexToColor(sex)" size="lg">
        <UIcon class="h-3.5" :name="sexToIcon(sex)" />
    </UBadge>
</template>
