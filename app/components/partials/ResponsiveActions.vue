<script lang="ts" setup>
import type { DropdownMenuItem } from '@nuxt/ui';
import type { ResponsiveActionEntry, ResponsiveActionItem, ResponsiveActionSeparator } from './ResponsiveActions.types';

const props = withDefaults(
    defineProps<{
        items: ResponsiveActionEntry[];
        label: string;
        icon?: string;
    }>(),
    {
        icon: 'i-mdi-menu',
    },
);

type NestedDropdownMenuItem = DropdownMenuItem & {
    children?: NestedDropdownMenuItem[];
};

function isSeparator(item: ResponsiveActionEntry): item is ResponsiveActionSeparator {
    return item.kind === 'separator';
}

function toDropdownItem(item: ResponsiveActionItem): NestedDropdownMenuItem {
    const mapped: NestedDropdownMenuItem = {
        label: item.label,
        icon: item.icon,
        color: item.color === 'neutral' ? undefined : item.color,
        disabled: item.disabled,
        kbds: item.kbds,
        to: item.to,
    };

    if (item.children?.length) {
        mapped.children = item.children.map((child) => toDropdownItem(child));
    } else if (item.onClick) {
        mapped.onClick = () => item.onClick?.();
    }

    return mapped;
}

function getMobileItems(): DropdownMenuItem[][] {
    return props.items.length > 0
        ? [props.items.map((item) => (isSeparator(item) ? ({ type: 'separator' } as DropdownMenuItem) : toDropdownItem(item)))]
        : [];
}

const hasItems = computed(() => props.items.some((item) => !isSeparator(item)));

function getDesktopItemKey(item: ResponsiveActionEntry, index: number): string {
    return isSeparator(item) ? `separator-${index}` : `${item.label}-${index}`;
}
</script>

<template>
    <div v-if="hasItems" class="w-full">
        <div class="mx-auto hidden w-full max-w-(--breakpoint-xl) flex-1 flex-row gap-2 overflow-x-auto md:flex">
            <template v-for="(item, index) in props.items" :key="getDesktopItemKey(item, index)">
                <div v-if="!isSeparator(item)" class="min-w-0 flex-1" :class="item.class">
                    <UTooltip :text="item.tooltip ?? item.label" :kbds="item.kbds">
                        <UButton
                            v-if="!item.children?.length"
                            block
                            :label="item.label"
                            :icon="item.icon"
                            :color="item.color ?? 'neutral'"
                            :variant="item.variant ?? 'ghost'"
                            :disabled="item.disabled"
                            :to="item.to"
                            @click="item.onClick"
                        />

                        <UDropdownMenu
                            v-else
                            class="flex w-full"
                            :items="item.children.map((child) => toDropdownItem(child))"
                            :content="{ align: 'center' }"
                            :ui="{ content: 'w-48' }"
                        >
                            <UButton
                                class="group"
                                block
                                :label="item.label"
                                :icon="item.icon"
                                :color="item.color ?? 'neutral'"
                                :variant="item.variant ?? 'ghost'"
                                :disabled="item.disabled"
                                trailing-icon="i-mdi-chevron-down"
                                :ui="{
                                    trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                                    label: 'flex-1',
                                }"
                            />
                        </UDropdownMenu>
                    </UTooltip>
                </div>
            </template>
        </div>

        <div class="flex md:hidden">
            <UDropdownMenu :items="getMobileItems()" :content="{ align: 'end' }" :ui="{ content: 'w-64' }">
                <template #default>
                    <UButton
                        class="ml-auto data-[state=open]:bg-elevated"
                        color="neutral"
                        variant="ghost"
                        :label="props.label"
                        :icon="props.icon"
                        trailing-icon="i-mdi-chevron-down"
                        :ui="{ trailingIcon: 'text-dimmed' }"
                    />
                </template>
            </UDropdownMenu>
        </div>
    </div>
</template>
