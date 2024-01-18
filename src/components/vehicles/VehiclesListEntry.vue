<script lang="ts" setup>
import { AccountEyeIcon, ClipboardPlusIcon } from 'mdi-vue3';
import LicensePlate from '~/components/partials/LicensePlate.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { toTitleCase } from '~/utils/strings';
import { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const props = defineProps<{
    vehicle: Vehicle;
    hideOwner?: boolean;
    hideCitizenLink?: boolean;
    hideCopy?: boolean;
}>();

function addToClipboard(): void {
    clipboardStore.addVehicle(props.vehicle);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.vehicle_added.title', parameters: {} },
        content: { key: 'notifications.clipboard.vehicle_added.content', parameters: {} },
        duration: 3250,
        type: 'info',
    });
}
</script>

<template>
    <tr :key="vehicle.plate" class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <td class="max-w-[4rem] whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            <LicensePlate :plate="vehicle.plate" class="mr-2" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ vehicle.model }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ toTitleCase(vehicle.type) }}
        </td>
        <td v-if="!hideOwner" class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            <CitizenInfoPopover :user="vehicle.owner" />
        </td>
        <td
            v-if="!hideCitizenLink || !hideCopy"
            class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0"
        >
            <div class="flex flex-row justify-end">
                <button v-if="!hideCopy" class="flex-initial text-primary-500 hover:text-primary-400" @click="addToClipboard()">
                    <ClipboardPlusIcon class="ml-auto mr-2.5 h-auto w-5" />
                </button>
                <NuxtLink
                    v-if="!hideCitizenLink && can('CitizenStoreService.ListCitizens')"
                    :to="{
                        name: 'citizens-id',
                        params: { id: vehicle.owner?.userId ?? 0 },
                    }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <AccountEyeIcon class="ml-auto mr-2.5 h-auto w-5" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
