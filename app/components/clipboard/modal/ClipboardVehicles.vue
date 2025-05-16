<script lang="ts" setup>
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import type { ClipboardVehicle } from '~/stores/clipboard';
import { useClipboardStore } from '~/stores/clipboard';
import { useNotificatorStore } from '~/stores/notificator';
import type { ObjectSpecs } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = withDefaults(
    defineProps<{
        submit?: boolean;
        showSelect?: boolean;
        specs?: ObjectSpecs;
    }>(),
    {
        submit: undefined,
        showSelect: true,
        specs: undefined,
    },
);

const emit = defineEmits<{
    (e: 'statisfied', payload: boolean): void;
    (e: 'close'): void;
}>();

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const { vehicles } = storeToRefs(clipboardStore);

const selected = ref<ClipboardVehicle[]>([]);

async function select(item: ClipboardVehicle): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    } else {
        if (props.specs && props.specs.max) {
            selected.value.splice(0, selected.value.length);
        }
        selected.value.push(item);
    }

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

async function remove(item: ClipboardVehicle, notify: boolean): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    }

    clipboardStore.removeVehicle(item.plate);
    if (notify) {
        notifications.add({
            title: { key: 'notifications.clipboard.vehicle_removed.title', parameters: {} },
            description: { key: 'notifications.clipboard.vehicle_removed.content', parameters: {} },
            timeout: 3250,
            type: NotificationType.INFO,
        });
    }
}

async function removeAll(): Promise<void> {
    while (selected.value.length > 0) {
        selected.value.forEach((v) => {
            remove(v, false);
        });
    }

    if (props.specs !== undefined) {
        emit('statisfied', false);
    } else {
        emit('statisfied', true);
    }

    notifications.add({
        title: { key: 'notifications.clipboard.vehicles_removed.title', parameters: {} },
        description: { key: 'notifications.clipboard.vehicles_removed.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}

watch(props, (newVal) => {
    if (newVal.submit) {
        if (clipboardStore.activeStack) {
            clipboardStore.activeStack.vehicles.length = 0;
            selected.value.forEach((v) => clipboardStore.activeStack.vehicles.push(v));
        } else if (vehicles.value && vehicles.value[0]) {
            selected.value.unshift(vehicles.value[0]);
        }
    }
});
</script>

<template>
    <h3 class="flex items-center justify-between text-lg font-medium">
        <span>{{ $t('common.vehicle', 2) }}</span>
        <slot name="header" />
    </h3>

    <DataNoDataBlock
        v-if="vehicles?.length === 0"
        icon="i-mdi-car"
        :message="$t('components.clipboard.clipboard_modal.no_data', [$t('common.vehicle', 2)])"
        :focus="
            async () => {
                navigateTo({ name: 'vehicles' });
                $emit('close');
            }
        "
    />
    <table v-else class="min-w-full divide-y divide-gray-700">
        <thead>
            <tr>
                <th v-if="showSelect" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1" scope="col">
                    {{ $t('common.select') }}
                </th>
                <th class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-1" scope="col">
                    {{ $t('common.plate') }}
                </th>
                <th class="px-3 py-3.5 text-left text-sm font-semibold" scope="col">
                    {{ $t('common.model') }}
                </th>
                <th class="px-3 py-3.5 text-left text-sm font-semibold" scope="col">
                    {{ $t('common.owner') }}
                </th>
                <th class="relative py-3.5 pl-3 pr-4 sm:pr-0" scope="col">
                    <span class="sr-only">{{ $t('common.action', 2) }}</span>
                    <UTooltip v-if="selected.length > 0" :text="$t('common.delete')">
                        <UButton variant="link" icon="i-mdi-delete" color="error" @click="removeAll()" />
                    </UTooltip>
                </th>
            </tr>
        </thead>
        <tbody class="divide-y divide-gray-800">
            <tr v-for="item in vehicles" :key="item.plate">
                <td v-if="showSelect" class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-1">
                    <UButton
                        v-if="specs && specs.max && specs.max === 1"
                        block
                        :color="selected.includes(item) ? 'gray' : 'primary'"
                        @click="select(item)"
                    >
                        {{
                            !selected.includes(item)
                                ? $t('common.select', 1).toUpperCase()
                                : $t('common.select', 2).toUpperCase()
                        }}
                    </UButton>
                    <UCheckbox
                        v-else
                        :key="item.plate"
                        v-model="selected"
                        name="selected"
                        :checked="selected.includes(item)"
                        :value="item"
                        @click="select(item)"
                    />
                </td>
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-1">
                    {{ item.plate }}
                </td>
                <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                    {{ item.model }}
                </td>
                <td class="whitespace-nowrap px-2 py-2 text-sm sm:px-4">
                    {{ item.owner.firstname }} {{ item.owner.lastname }}
                </td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                    <UTooltip :text="$t('common.delete')">
                        <UButton variant="link" icon="i-mdi-delete" color="error" @click="remove(item, true)" />
                    </UTooltip>
                </td>
            </tr>
        </tbody>
    </table>
</template>
