<script lang="ts" setup>
import type { Sort } from '~~/gen/ts/resources/common/database/database';

const props = withDefaults(
    defineProps<{
        modelValue: Sort;
        fields: { label: string; value: string }[];
    }>(),
    {},
);

const emit = defineEmits<{
    (e: 'update:modelValue', v: Sort): void;
}>();

const sorting = useVModel(props, 'modelValue', emit, {
    deep: true,
});

const { custom } = useAppConfig();

function toggleDirection(): void {
    if (sorting.value.columns.length === 0) {
        sorting.value = {
            columns: [
                {
                    id: sorting.value.columns.at(0)?.id || props.fields[0]?.value || '',
                    desc: true,
                },
            ],
        };
    } else {
        sorting.value = {
            columns: [
                {
                    id: sorting.value.columns.at(0)?.id || props.fields[0]?.value || '',
                    desc: false,
                },
            ],
        };
    }
}

function changeColumn(col: string): void {
    sorting.value = {
        columns: [
            {
                id: col,
                desc: sorting.value.columns.at(0)?.desc || false,
            },
        ],
    };
}
</script>

<template>
    <div class="flex gap-2">
        <ClientOnly v-if="fields.length > 1">
            <USelectMenu
                class="w-full"
                :placeholder="$t('common.na')"
                value-key="value"
                :items="fields"
                @update:model-value="($event) => changeColumn($event)"
            >
                <template #empty> {{ $t('common.not_found', [$t('common.field', 2)]) }} </template>
            </USelectMenu>
        </ClientOnly>

        <UTooltip :text="$t('common.sort_direction')">
            <UButton
                square
                trailing
                :icon="sorting.columns[0]?.desc ? custom.icons.sortDesc : custom.icons.sortAsc"
                color="neutral"
                variant="ghost"
                @click="() => toggleDirection()"
            />
        </UTooltip>
    </div>
</template>
