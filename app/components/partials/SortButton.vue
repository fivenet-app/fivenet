<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        modelValue: TableSortable;
        fields: { label: string; value: string }[];
    }>(),
    {},
);

const emit = defineEmits<{
    (e: 'update:modelValue', v: TableSortable): void;
}>();

const sort = useVModel(props, 'modelValue', emit, {
    deep: true,
});

const { ui } = useAppConfig();

function toggleDirection(): void {
    if (sort.value.direction === 'asc') {
        sort.value = {
            column: sort.value.column,
            direction: 'desc',
        };
    } else {
        sort.value = {
            column: sort.value.column,
            direction: 'asc',
        };
    }
}

function changeColumn(col: string): void {
    sort.value = {
        column: col,
        direction: sort.value.direction,
    };
}
</script>

<template>
    <div class="flex gap-2">
        <ClientOnly v-if="fields.length > 1">
            <USelectMenu
                class="w-full"
                :model-value="sort.column"
                :placeholder="$t('common.na')"
                value-attribute="value"
                :options="fields"
                @update:model-value="changeColumn($event)"
            >
                <template #label>
                    {{ fields.find((f) => f.value === sort.column)?.label ?? $t('common.na') }}
                </template>

                <template #option="{ option: field }">
                    {{ field.label }}
                </template>

                <template #empty> {{ $t('common.not_found', [$t('common.field', 2)]) }} </template>
            </USelectMenu>
        </ClientOnly>

        <UTooltip :text="$t('common.sort_direction')">
            <UButton
                square
                trailing
                :icon="sort.direction === 'asc' ? ui.table.default.sortAscIcon : ui.table.default.sortDescIcon"
                color="gray"
                variant="ghost"
                @click="toggleDirection"
            />
        </UTooltip>
    </div>
</template>
