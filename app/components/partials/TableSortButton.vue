<script lang="ts" setup>
import type { DropdownMenuItem } from '@nuxt/ui';

type SortableColumn = {
    clearSorting: () => void;
    getIsSorted: () => false | 'asc' | 'desc';
    toggleSorting: (desc?: boolean, isMulti?: boolean) => void;
};

const props = withDefaults(
    defineProps<{
        column: SortableColumn;
        label: string;
        clearable?: boolean;
        mode?: 'menu' | 'split' | 'toggle';
    }>(),
    {
        clearable: true,
        mode: 'split',
    },
);

const { custom } = useAppConfig();
const { t } = useI18n();

function getSortDirection(): false | 'asc' | 'desc' {
    return props.column.getIsSorted();
}

function getSortIcon(): string {
    const sortDirection = getSortDirection();

    if (sortDirection === 'asc') {
        return custom.icons.sortAsc;
    }

    if (sortDirection === 'desc') {
        return custom.icons.sortDesc;
    }

    return custom.icons.sort;
}

function getNextSortLabel(): string {
    return getSortDirection() === 'asc' ? 'descending' : 'ascending';
}

function toggleSort(): void {
    props.column.toggleSorting(getSortDirection() === 'asc');
}

function selectAscending(): void {
    if (getSortDirection() === 'asc') {
        if (props.clearable) {
            props.column.clearSorting();
        }
        return;
    }

    props.column.toggleSorting(false);
}

function selectDescending(): void {
    if (getSortDirection() === 'desc') {
        if (props.clearable) {
            props.column.clearSorting();
        }
        return;
    }

    props.column.toggleSorting(true);
}

function getItems(): DropdownMenuItem[] {
    const sortDirection = getSortDirection();

    const items: DropdownMenuItem[] = [
        {
            label: t('common.asc'),
            type: 'checkbox',
            icon: custom.icons.sortAsc,
            checked: sortDirection === 'asc',
            onSelect: () => selectAscending(),
        },
        {
            label: t('common.desc'),
            type: 'checkbox',
            icon: custom.icons.sortDesc,
            checked: sortDirection === 'desc',
            onSelect: () => selectDescending(),
        },
    ];

    if (sortDirection && props.clearable) {
        items.push({
            label: t('common.clear'),
            icon: 'i-mdi-close',
            onSelect: () => props.column.clearSorting(),
        });
    }

    return items;
}
</script>

<template>
    <div class="-mx-2.5 inline-flex items-center">
        <UButton
            v-if="mode === 'toggle'"
            color="neutral"
            variant="ghost"
            :label="label"
            :icon="getSortIcon()"
            :aria-label="`Sort by ${getNextSortLabel()}`"
            @click="toggleSort"
        />

        <UDropdownMenu v-else-if="mode === 'menu'" :items="getItems()" :content="{ align: 'start' }">
            <UButton
                class="data-[state=open]:bg-elevated"
                color="neutral"
                variant="ghost"
                :label="label"
                :icon="getSortIcon()"
                :aria-label="`Sort by ${getNextSortLabel()}`"
            />
        </UDropdownMenu>

        <template v-else>
            <UButton
                class="rounded-r-none"
                color="neutral"
                variant="ghost"
                :label="label"
                :icon="getSortIcon()"
                :aria-label="`Sort by ${getNextSortLabel()}`"
                @click="toggleSort"
            />

            <UDropdownMenu :items="getItems()" :content="{ align: 'start' }">
                <UButton
                    class="-ml-px rounded-l-none data-[state=open]:bg-elevated"
                    square
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-chevron-down"
                    :aria-label="`${label} sort options`"
                />
            </UDropdownMenu>
        </template>
    </div>
</template>
