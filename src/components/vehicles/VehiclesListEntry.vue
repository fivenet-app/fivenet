<script setup lang="ts">
import { Vehicle } from '@arpanet/gen/resources/vehicles/vehicles_pb';

defineProps({
    'vehicle': {
        required: true,
        type: Vehicle,
    },
    'hideOwner': {
        type: Boolean,
        required: false,
        default: false,
    },
});
</script>

<template>
    <tr :key="vehicle.getPlate()">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
            {{ vehicle.getPlate() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ vehicle.getModel() }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ vehicle.getType().toUpperCase() }}
        </td>
        <td v-if="!hideOwner" class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ vehicle.getOwner()?.getFirstname() }}, {{ vehicle.getOwner()?.getLastname() }}
        </td>
        <td v-if="!hideOwner" class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            {{ vehicle.getOwner()?.getJobLabel() }}
        </td>
        <td class="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div v-can="'CitizenStoreService.FindUsers'">
                <router-link :to="{ name: 'Citizens Info', params: { id: vehicle.getOwner()?.getUserId().toString() } }"
                    class="text-indigo-400 hover:text-indigo-300">VIEW OWNER</router-link>
            </div>
        </td>
    </tr>
</template>
