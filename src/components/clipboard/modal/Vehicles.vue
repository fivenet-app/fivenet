<script lang="ts" setup>
import { CarIcon, TrashCanIcon } from 'mdi-vue3';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { ClipboardVehicle, useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { ObjectSpecs } from '~~/gen/ts/resources/documents/templates';

const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const { vehicles } = storeToRefs(clipboardStore);

const emit = defineEmits<{
    (e: 'statisfied', payload: boolean): void;
}>();

const props = withDefaults(
    defineProps<{
        submit?: boolean;
        showSelect?: boolean;
        specs?: ObjectSpecs;
    }>(),
    {
        submit: false,
        showSelect: false,
    },
);

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

    const selectedLength = BigInt(selected.value.length);
    if (props.specs) {
        if (props.specs.min && selectedLength >= props.specs.min) {
            emit('statisfied', true);
        } else if (props.specs.max && selectedLength === props.specs.max) {
            emit('statisfied', true);
        } else {
            emit('statisfied', false);
        }
    }
}

async function remove(item: ClipboardVehicle, notify: boolean): Promise<void> {
    const idx = selected.value.indexOf(item);
    if (idx !== undefined && idx > -1) {
        selected.value.splice(idx, 1);
    }

    clipboardStore.removeVehicle(item.plate);
    if (notify) {
        notifications.dispatchNotification({
            title: { key: 'notifications.clipboard.vehicle_removed.title', parameters: [] },
            content: { key: 'notifications.clipboard.vehicle_removed.content', parameters: [] },
            duration: 3500,
            type: 'info',
        });
    }
}

async function removeAll(): Promise<void> {
    while (selected.value.length > 0) {
        selected.value.forEach((v) => {
            remove(v, false);
        });
    }

    emit('statisfied', false);
    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.vehicles_removed.title', parameters: [] },
        content: { key: 'notifications.clipboard.vehicles_removed.content', parameters: [] },
        duration: 3500,
        type: 'info',
    });
}

watch(props, (newVal) => {
    if (newVal.submit) {
        if (clipboardStore.activeStack) {
            clipboardStore.activeStack.vehicles.length = 0;
            selected.value.forEach((v) => clipboardStore.activeStack.vehicles.push(v));
        } else if (vehicles.value && vehicles.value.length === 1) {
            selected.value.unshift(vehicles.value[0]);
        }
    }
});
</script>

<template>
    <h3 class="font-medium pt-2 pb-1">Vehicles</h3>
    <DataNoDataBlock
        v-if="vehicles?.length === 0"
        :icon="CarIcon"
        :message="$t('components.clipboard.clipboard_modal.no_data', [$t('common.vehicle', 2)])"
    />
    <table v-else class="min-w-full divide-y divide-gray-700">
        <thead>
            <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0" v-if="showSelect">
                    {{ $t('common.select') }}
                </th>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">
                    {{ $t('common.plate') }}
                </th>
                <th scope="col" class="py-3.5 px-3 text-left text-sm font-semibold text-white">
                    {{ $t('common.model') }}
                </th>
                <th scope="col" class="py-3.5 px-3 text-left text-sm font-semibold text-white">
                    {{ $t('common.owner') }}
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                    <span class="sr-only">{{ $t('common.action', 2) }}</span>
                    <button v-if="selected.length > 0" @click="removeAll()">
                        <TrashCanIcon class="w-6 h-6 mx-auto text-neutral" />
                    </button>
                </th>
            </tr>
        </thead>
        <tbody class="divide-y divide-gray-800">
            <tr v-for="item in vehicles" :key="item.plate">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0" v-if="showSelect">
                    <div v-if="specs && specs.max && specs.max === 1n">
                        <button
                            @click="select(item)"
                            class="inline-flex w-full justify-center rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 sm:col-start-2"
                        >
                            <span v-if="!selected.includes(item)">
                                {{ $t('common.select', 1).toUpperCase() }}
                            </span>
                            <span v-else>
                                {{ $t('common.select', 2).toUpperCase() }}
                            </span>
                        </button>
                    </div>
                    <div v-else>
                        <input
                            name="selected"
                            :key="item.plate"
                            :checked="selected.includes(item)"
                            :value="item"
                            v-model="selected"
                            type="checkbox"
                            class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-600"
                        />
                    </div>
                </td>
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
                    {{ item.plate }}
                </td>
                <td class="whitespace-nowrap py-4 px-3 text-sm text-gray-300">
                    {{ item.model }}
                </td>
                <td class="whitespace-nowrap py-4 px-3 text-sm text-gray-300">
                    {{ item.owner.firstname }}, {{ item.owner.lastname }}
                </td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                    <button @click="remove(item, true)">
                        <TrashCanIcon class="w-6 h-6 mx-auto text-neutral" />
                    </button>
                </td>
            </tr>
        </tbody>
    </table>
</template>
