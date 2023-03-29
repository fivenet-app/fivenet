<script lang="ts" setup>
import { Vehicle } from '@arpanet/gen/resources/vehicles/vehicles_pb';
import { ClipboardDocumentIcon, EyeIcon } from '@heroicons/vue/24/solid';
import { useStore } from '../../store/store';
import { toTitleCase } from '../../utils/strings';

const store = useStore();

const props = defineProps({
    vehicle: {
        required: true,
        type: Vehicle,
    },
    hideOwner: {
        type: Boolean,
        required: false,
        default: false,
    },
    hideCitizenLink: {
        type: Boolean,
        required: false,
        default: false,
    },
    hideCopy: {
        type: Boolean,
        required: false,
        default: false,
    },
});

function addToClipboard() {
    store.dispatch('clipboard/addVehicle', props.vehicle);
}
</script>

<template>
    <tr :key="vehicle.getPlate()">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ vehicle.getPlate() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ vehicle.getModel() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ toTitleCase(vehicle.getType()) }}
        </td>
        <td v-if="!hideOwner" class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ vehicle.getOwner()?.getFirstname() }}, {{ vehicle.getOwner()?.getLastname() }}
        </td>
        <td v-if="!hideOwner" class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ vehicle.getOwner()?.getJobLabel() }}
        </td>
        <td v-if="!hideCitizenLink && !hideCopy"
            class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div class="flex flex-row justify-end">
                <router-link v-if="!hideCitizenLink" v-can="'CitizenStoreService.FindUsers'"
                    :to="{ name: 'Citizens: Info', params: { id: vehicle.getOwner()?.getUserId() ?? 0 } }"
                    class="flex-initial text-primary-500 hover:text-primary-400">
                    <EyeIcon class="w-6 h-auto ml-auto mr-2.5" />
                </router-link>
                <button v-if="!hideCopy" class="flex-initial text-primary-500 hover:text-primary-400"
                    @click="addToClipboard()">
                    <ClipboardDocumentIcon class="w-6 h-auto ml-auto mr-2.5" />
                </button>
            </div>
        </td>
    </tr>
</template>
