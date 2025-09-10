<script lang="ts" setup>
import type { TableColumn } from '@nuxt/ui';
import { type ClipboardVehicle, useClipboardStore } from '~/stores/clipboard';
import type { ObjectSpecs } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        submit?: boolean;
        showSelect?: boolean;
        specs?: ObjectSpecs;
        hideHeader?: boolean;
    }>(),
    {
        submit: undefined,
        showSelect: true,
        specs: undefined,
        hideHeader: false,
    },
);

const emit = defineEmits<{
    (e: 'statisfied', payload: boolean): void;
    (e: 'close'): void;
}>();

const { t } = useI18n();

const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const { vehicles } = storeToRefs(clipboardStore);

const selected = ref<string[]>([]);

async function select(): Promise<void> {
    const selectedLength = selected.value.length;
    if (props.specs) {
        if (props.specs.min !== undefined && selectedLength >= props.specs.min) {
            emit('statisfied', true);
        } else if (props.specs.max !== undefined && selectedLength === props.specs.max) {
            emit('statisfied', true);
        } else {
            emit('statisfied', false);
        }
    } else {
        emit('statisfied', true);
    }
}

watch(selected, () => select());

async function remove(item: string, notify: boolean): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    }

    clipboardStore.removeVehicle(item);
    if (notify) {
        notifications.add({
            title: { key: 'notifications.clipboard.vehicle_removed.title', parameters: {} },
            description: { key: 'notifications.clipboard.vehicle_removed.content', parameters: {} },
            duration: 3250,
            type: NotificationType.INFO,
        });
    }
}

async function removeAll(): Promise<void> {
    // Make a shallow copy to avoid mutation issues
    const toRemove = [...selected.value];
    toRemove.forEach((v) => {
        remove(v, false);
    });
    selected.value = [];

    if (props.specs !== undefined) {
        emit('statisfied', false);
    } else {
        emit('statisfied', true);
    }

    notifications.add({
        title: { key: 'notifications.clipboard.vehicles_removed.title', parameters: {} },
        description: { key: 'notifications.clipboard.vehicles_removed.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const columns = computed(() =>
    (
        [
            props.showSelect
                ? {
                      id: 'actions',
                  }
                : undefined,
            {
                accesssorKey: 'plate',
                header: t('common.plate'),
                cell: ({ row }) => h('span', { class: 'text-highlighted' }, row.original.plate),
            },
            {
                accesssorKey: 'model',
                header: t('common.model'),
                cell: ({ row }) => h('span', {}, row.original.model),
            },
            {
                accesssorKey: 'owner',
                header: t('common.owner'),
                cell: ({ row }) => h('span', {}, `${row.original.owner.firstname} ${row.original.owner.lastname}`),
            },
            {
                id: 'delete',
            },
        ] as TableColumn<ClipboardVehicle>[]
    ).filter((c) => c !== undefined),
);

watch(props, (newVal) => {
    if (newVal.submit) {
        if (clipboardStore.activeStack) {
            clipboardStore.activeStack.vehicles.length = 0;
            selected.value.forEach((v) =>
                clipboardStore.activeStack.vehicles.push(clipboardStore.vehicles.find((veh) => veh.plate === v)!),
            );
        } else if (vehicles.value && vehicles.value[0]) {
            selected.value.unshift(vehicles.value[0].plate);
        }
    }
});
</script>

<template>
    <div>
        <h3 v-if="!hideHeader" class="flex items-center justify-between text-lg font-medium">
            <span>{{ $t('common.vehicle', 2) }}</span>
            <slot name="header" />
        </h3>

        <UTable :columns="columns" :data="vehicles" :empty="$t('common.not_found', [$t('common.vehicle', 2)])">
            <template #actions-cell="{ row }">
                <URadioGroup
                    v-if="specs && specs.max && specs.max === 1"
                    :model-value="selected[0]"
                    name="selected"
                    :items="[row.original.plate]"
                    value-key="plate"
                    :ui="{ label: 'hidden' }"
                    @update:model-value="(v) => (selected = [v])"
                />
                <UCheckboxGroup
                    v-else
                    :key="row.original.plate"
                    v-model="selected"
                    name="selected"
                    :items="[row.original.plate]"
                    value-key="plate"
                    :ui="{ label: 'hidden' }"
                />
            </template>

            <template v-if="selected.length > 0" #actions-header>
                <UTooltip :text="$t('common.delete')">
                    <UButton variant="link" icon="i-mdi-delete" color="error" size="xs" @click="removeAll()" />
                </UTooltip>
            </template>

            <template #delete-cell="{ row }">
                <UTooltip :text="$t('common.delete')">
                    <UButton variant="link" icon="i-mdi-delete" color="error" @click="remove(row.original.plate, true)" />
                </UTooltip>
            </template>
        </UTable>
    </div>
</template>
