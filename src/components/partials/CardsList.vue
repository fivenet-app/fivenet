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
    <div class="sm:px-2">
        <UPageGrid>
            <UPageCard
                v-for="(module, index) in items.filter((i) => i.permission === undefined || can(i.permission))"
                :key="index"
                :href="module.to"
                :title="module.title"
                :description="module.description"
                :to="module.to"
                :icon="!showIcon ? undefined : module.icon"
                @click="$emit('selected', index)"
            >
                <template #description>
                    <span class="line-clamp-2">{{ module.description }}</span>
                </template>
            </UPageCard>
        </UPageGrid>
    </div>
</template>
