<script lang="ts" setup>
import { ClipboardDocumentIcon, EyeIcon } from '@heroicons/vue/24/solid';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { toTitleCase } from '~/utils/strings';
import { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';

const store = useClipboardStore();
const notifications = useNotificationsStore();

const { t } = useI18n();

const props = defineProps<{
    vehicle: Vehicle;
    hideOwner?: boolean;
    hideCitizenLink?: boolean;
    hideCopy?: boolean;
}>();

function addToClipboard(): void {
    store.addVehicle(props.vehicle);

    notifications.dispatchNotification({
        title: t('notifications.clipboard.vehicle_added.title'),
        content: t('notifications.clipboard.vehicle_added.content'),
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
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ vehicle.model }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ toTitleCase(vehicle.type) }}
        </td>
        <td v-if="!hideOwner" class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ vehicle.owner?.firstname }}, {{ vehicle.owner?.lastname }}
        </td>
        <td
            v-if="!hideCitizenLink && !hideCopy"
            class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0"
        >
            <div class="flex flex-row justify-end">
                <button v-if="!hideCopy" class="flex-initial text-primary-500 hover:text-primary-400" @click="addToClipboard()">
                    <ClipboardDocumentIcon class="w-6 h-auto ml-auto mr-2.5" />
                </button>
                <NuxtLink
                    v-if="!hideCitizenLink"
                    v-can="'CitizenStoreService.ListCitizens'"
                    :to="{
                        name: 'citizens-id',
                        params: { id: vehicle.owner?.userId ?? 0 },
                    }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <EyeIcon class="w-6 h-auto ml-auto mr-2.5" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
