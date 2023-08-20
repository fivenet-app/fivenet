<script lang="ts" setup>
import { AccountEyeIcon, ClipboardPlusIcon } from 'mdi-vue3';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { toTitleCase } from '~/utils/strings';
import { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';

const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const props = defineProps<{
    vehicle: Vehicle;
    hideOwner?: boolean;
    hideCitizenLink?: boolean;
    hideCopy?: boolean;
}>();

function addToClipboard(): void {
    clipboardStore.addVehicle(props.vehicle);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.vehicle_added.title', parameters: [] },
        content: { key: 'notifications.clipboard.vehicle_added.content', parameters: [] },
        duration: 3500,
        type: 'info',
    });
}
</script>

<template>
    <tr :key="vehicle.plate">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ vehicle.plate }}
        </td>
        <td class="whitespace-nowrap px-1 py-1text-sm text-base-200">
            {{ vehicle.model }}
        </td>
        <td class="whitespace-nowrap px-1 py-1text-sm text-base-200">
            {{ toTitleCase(vehicle.type) }}
        </td>
        <td v-if="!hideOwner" class="whitespace-nowrap px-1 py-1text-sm text-base-200">
            {{ vehicle.owner?.firstname }}, {{ vehicle.owner?.lastname }}
        </td>
        <td
            v-if="!hideCitizenLink && !hideCopy"
            class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0"
        >
            <div class="flex flex-row justify-end">
                <button v-if="!hideCopy" class="flex-initial text-primary-500 hover:text-primary-400" @click="addToClipboard()">
                    <ClipboardPlusIcon class="w-6 h-auto ml-auto mr-2.5" />
                </button>
                <NuxtLink
                    v-if="!hideCitizenLink && can('CitizenStoreService.ListCitizens')"
                    :to="{
                        name: 'citizens-id',
                        params: { id: vehicle.owner?.userId ?? 0 },
                    }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <AccountEyeIcon class="w-6 h-auto ml-auto mr-2.5" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
