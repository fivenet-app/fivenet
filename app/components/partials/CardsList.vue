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
            v-for="(module, index) in items.filter((i) => i.permission === undefined || can(i.permission).value)"
            :key="index"
            :to="module.to"
            :title="module.title"
            :description="module.description"
            :icon="showIcon && module.icon?.startsWith('i-') ? module.icon : undefined"
            v-bind="$attrs"
            @click="$emit('selected', index)"
        >
            <template v-if="showIcon && module.icon" #icon>
                <template v-if="!module.icon.startsWith('i-')">
                    <component
                        :is="availableIcons.find((item) => item.name === module.icon)?.component ?? fallbackIcon"
                        v-if="module.icon"
                        class="text-primary h-10 w-10 shrink-0"
                        :class="module.color && `text-${module.color}-500 dark:text-${module.color}-400`"
                    />
                </template>
                <template v-else>
                    <UIcon
                        class="text-primary h-10 w-10 shrink-0"
                        :class="`text-${module.color}-500 dark:text-${module.color}-400`"
                        :name="module.icon"
                    />
                </template>
            </template>

            <template #description>
                <span class="line-clamp-2">{{ module.description }}</span>
            </template>
        </UPageCard>
    </UPageGrid>
</template>
