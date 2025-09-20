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

const { can } = useAuth();
</script>

<template>
    <UPageGrid>
        <UPageCard
            v-for="(module, index) in items.filter((i) => i.permission === undefined || can(i.permission).value)"
            :key="index"
            :to="module.to"
            :title="module.title"
            :icon="showIcon && module.icon?.startsWith('i-') ? module.icon : undefined"
            @click="$emit('selected', index)"
        >
            <template v-if="showIcon && module.icon" #leading>
                <template v-if="!module.icon.startsWith('i-')">
                    <component
                        :is="availableIcons.find((item) => item.name === module.icon)?.component ?? fallbackIcon.component"
                        v-if="module.icon"
                        class="h-10 w-10 shrink-0"
                        :class="`text-${module.color ?? 'primary'}`"
                    />
                </template>
                <template v-else>
                    <UIcon class="h-10 w-10 shrink-0" :class="`text-${module.color ?? 'primary'}`" :name="module.icon" />
                </template>
            </template>

            <template #description>
                <span class="line-clamp-2">{{ module.description }}</span>
            </template>
        </UPageCard>
    </UPageGrid>
</template>
