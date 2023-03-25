<script lang="ts" setup>
import { Vehicle } from '@arpanet/gen/resources/vehicles/vehicles_pb';
import { EyeIcon } from '@heroicons/vue/24/solid';
import { toTitleCase } from '../../utils/strings';

defineProps({
    vehicle: {
        required: true,
        type: Vehicle,
    },
    hideOwner: {
        type: Boolean,
        required: false,
        default: false,
    },
});
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
        <td class="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div v-can="'CitizenStoreService.FindUsers'">
                <router-link :to="{ name: 'Citizens: Info', params: { id: vehicle.getOwner()?.getUserId() ?? 0 } }"
                    class="text-primary-500 hover:text-primary-400"><EyeIcon class="w-6 h-auto ml-auto mr-2.5" /></router-link>
            </div>
        </td>
    </tr>
</template>
