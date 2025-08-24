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

const sorting = useVModel(props, 'modelValue', emit, {
    deep: true,
});

const { custom } = useAppConfig();

function toggleDirection(): void {
    if (sorting.value.direction === 'asc') {
        sorting.value = {
            id: sorting.value.id,
            desc: true,
        };
    } else {
        sorting.value = {
            id: sorting.value.id,
            desc: false,
        };
    }
}

function changeColumn(col: string): void {
    sorting.value = {
        id: col,
        direction: sorting.value.direction,
    };
}
</script>

<template>
    <div class="flex gap-2">
        <ClientOnly v-if="fields.length > 1">
            <USelectMenu
                class="w-full"
                :model-value="sorting.id"
                :placeholder="$t('common.na')"
                value-key="value"
                :items="fields"
                @update:model-value="changeColumn($event)"
            >
                <template #item-label>
                    {{ fields.find((f) => f.value === sorting.id)?.label ?? $t('common.na') }}
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
                :icon="sorting.direction === 'asc' ? custom.icons.sortAscIcon : custom.icons.sortDescIcon"
                color="neutral"
                variant="ghost"
                @click="toggleDirection"
            />
        </UTooltip>
    </div>
</template>
