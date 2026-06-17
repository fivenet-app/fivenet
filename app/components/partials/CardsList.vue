<script lang="ts" setup>
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
            :ui="{ title: 'w-full flex flex-row gap-2' }"
            @click="
                () => {
                    !module.to && $emit('selected', index);
                }
            "
        >
            <template #title>
                <span>{{ module.title }}</span>

                <UBadge v-if="module.deletedAt" icon="i-mdi-delete" :label="$t('common.deleted')" color="amber" />
            </template>

            <template v-if="showIcon && module.icon" #leading>
                <template v-if="!module.icon.startsWith('i-')">
                    <UIcon
                        v-if="module.icon"
                        class="h-10 w-10 shrink-0"
                        :class="`text-${module.color ?? 'primary'}`"
                        :name="convertComponentIconNameToDynamic(module.icon)"
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
