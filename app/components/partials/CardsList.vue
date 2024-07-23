<script lang="ts" setup>
import { type CardElement } from '~/utils/types';

withDefaults(
    defineProps<{
        items: CardElement[];
        showIcon?: boolean;
    }>(),
    {
        showIcon: true,
    },
);

defineEmits<{
    (e: 'selected', idx: number): void;
}>();
</script>

<template>
    <UPageGrid>
        <UPageCard
            v-for="(module, index) in items.filter((i) => i.permission === undefined || can(i.permission).value)"
            :key="index"
            :to="module.to"
            :title="module.title"
            :description="module.description"
            :icon="!showIcon ? undefined : module.icon"
            @click="$emit('selected', index)"
        >
            <template #description>
                <span class="line-clamp-2">{{ module.description }}</span>
            </template>
        </UPageCard>
    </UPageGrid>
</template>
