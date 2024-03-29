<script lang="ts" setup>
import { GenderFemaleIcon, GenderMaleIcon } from 'mdi-vue3';

defineProps<{
    sex: string;
}>();
</script>

<template>
    <span
        v-if="sex === 'f'"
        class="inline-flex items-center rounded-full bg-pink-200 px-2 py-0.5 text-sm font-medium text-pink-600"
    >
        <GenderFemaleIcon class="size-5" aria-hidden="true" />
    </span>

    <span
        v-else-if="sex === 'm'"
        class="inline-flex items-center rounded-full bg-blue-200 px-2 py-0.5 text-sm font-medium text-blue-600"
    >
        <GenderMaleIcon class="size-5" aria-hidden="true" />
    </span>
</template>
