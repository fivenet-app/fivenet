<script lang="ts" setup>
import { markerFallbackIcon, markerIcons } from '~/components/livemap/helpers';
import type { CardElement } from '~/utils/types';

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
            v-bind="$attrs"
            :key="index"
            :to="module.to"
            :title="module.title"
            :description="module.description"
            :icon="showIcon && module.icon?.startsWith('i-') ? module.icon : undefined"
            @click="$emit('selected', index)"
        >
            <template v-if="showIcon && !module.icon?.startsWith('i-')" #icon>
                <div v-if="showIcon">
                    <component
                        :is="markerIcons.find((item) => item.name === module.icon) ?? markerFallbackIcon"
                        v-if="module.icon"
                        class="text-primary h-10 w-10 flex-shrink-0"
                        :style="{ color: module.color }"
                    />
                </div>
            </template>

            <template #description>
                <span class="line-clamp-2">{{ module.description }}</span>
            </template>
        </UPageCard>
    </UPageGrid>
</template>
