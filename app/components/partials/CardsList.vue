<script lang="ts" setup>
import type { ContextMenuItem } from '@nuxt/ui';
import DeletedAtBadge from '~/components/partials/DeletedAtBadge.vue';
import type { CardElement } from '~/utils/types';

withDefaults(
    defineProps<{
        items: CardElement[];
        showIcon?: boolean;
        getContextMenuItems?: (item: CardElement, idx: number) => ContextMenuItem[][];
    }>(),
    {
        showIcon: true,
        getContextMenuItems: undefined,
    },
);

defineEmits<{
    (e: 'selected', idx: number): void;
}>();

const { can } = useAuth();
</script>

<template>
    <UPageGrid>
        <template
            v-for="(module, index) in items.filter((i) => i.permission === undefined || can(i.permission).value)"
            :key="module.to ?? index"
        >
            <UContextMenu
                :items="getContextMenuItems?.(module, index) ?? []"
                :disabled="(getContextMenuItems?.(module, index) ?? []).length === 0"
            >
                <UPageCard
                    :to="module.to"
                    :title="module.label"
                    :icon="showIcon && module.icon?.startsWith('i-') ? module.icon : undefined"
                    :ui="{ title: 'w-full flex flex-row gap-2' }"
                    @click="
                        () => {
                            !module.to && $emit('selected', index);
                        }
                    "
                >
                    <template #title>
                        <span>{{ module.label }}</span>

                        <DeletedAtBadge v-if="module.deletedAt" hide-date icon="i-mdi-delete" :deleted-at="module.deletedAt" />
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
                            <UIcon
                                class="h-10 w-10 shrink-0"
                                :class="`text-${module.color ?? 'primary'}`"
                                :name="module.icon"
                            />
                        </template>
                    </template>

                    <template #description>
                        <span class="line-clamp-2">{{ module.description }}</span>
                    </template>
                </UPageCard>
            </UContextMenu>
        </template>
    </UPageGrid>
</template>
