<script lang="ts" setup>
import type { CardElement } from '~/utils/types';
import { availableIcons, fallbackIcon } from './icons';

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

defineOptions({
    inheritAttrs: false,
});

const { can } = useAuth();
</script>

<template>
    <UPageGrid>
        <UPageCard
            v-for="(card, index) in items.filter((i) => i.permission === undefined || can(i.permission).value)"
            v-bind="$attrs"
            :key="index"
            :to="card.to"
            :title="card.title"
            :description="card.description"
            :icon="showIcon && card.icon?.startsWith('i-') ? card.icon : undefined"
            @click="$emit('selected', index)"
        >
            <template v-if="showIcon && card.icon" #leading>
                <template v-if="!card.icon.startsWith('i-')">
                    <component
                        :is="availableIcons.find((item) => item.name === card.icon) ?? fallbackIcon"
                        v-if="card.icon"
                        class="text-primary h-10 w-10 shrink-0"
                        :class="card.color && `text-${card.color}-500 dark:text-${card.color}-400`"
                    />
                </template>
                <template v-else>
                    <UIcon
                        :name="card.icon"
                        class="text-primary h-10 w-10 shrink-0"
                        :class="`text-${card.color}-500 dark:text-${card.color}-400`"
                    />
                </template>
            </template>

            <template #description>
                <span class="line-clamp-2">{{ card.description }}</span>
            </template>
        </UPageCard>
    </UPageGrid>
</template>
