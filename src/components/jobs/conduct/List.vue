<script lang="ts" setup>
import Time from '~/components/partials/elements/Time.vue';
import { CONDUCT_TYPE, ConductEntry } from '~~/gen/ts/resources/jobs/conduct';
import { conductTypesToBGColor, conductTypesToRingColor, conductTypesToTextColor } from './helpers';

const statuses = {
    Paid: 'text-green-700 bg-green-50 ring-green-600/20',
    Withdraw: 'text-gray-600 bg-gray-50 ring-gray-500/10',
    Overdue: 'text-red-700 bg-red-50 ring-red-600/10',
};
const entries: ConductEntry[] = [
    {
        id: 0n,
        job: 'ambulance',
        description: '',
        targetUserId: 26061,
        targetUser: {
            userId: 26061,
            identifier: '',
            job: 'ambulance',
            dateofbirth: '01.08.1997',
            firstname: 'Prof. Dr. Philipp',
            lastname: 'Scott',
            jobGrade: 20,
        },
        creatorId: 26061,
        creator: {
            userId: 26061,
            identifier: '',
            job: 'ambulance',
            dateofbirth: '01.08.1997',
            firstname: 'Prof. Dr. Philipp',
            lastname: 'Scott',
            jobGrade: 20,
        },
        type: CONDUCT_TYPE.NEGATIVE,
    },
];
</script>

<template>
    <div>
        <div class="mt-6 overflow-hidden border-t border-gray-100">
            <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div class="mx-auto max-w-2xl lg:mx-0 lg:max-w-none">
                    <table class="w-full text-left">
                        <thead class="sr-only">
                            <tr>
                                <th>Type</th>
                                <th class="hidden sm:table-cell">Client</th>
                                <th>More details</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="conduct in entries" :key="conduct.id.toString()">
                                <td class="relative py-5 pr-6">
                                    <div
                                        class="rounded-md py-1 px-2 text-xs font-medium ring-1 ring-inset"
                                        :class="[
                                            conductTypesToBGColor(conduct.type),
                                            conductTypesToTextColor(conduct.type),
                                            conductTypesToRingColor(conduct.type),
                                        ]"
                                    >
                                        {{ CONDUCT_TYPE[conduct.type] }}
                                    </div>
                                </td>
                                <td class="relative py-5 pr-6">
                                    <div class="flex gap-x-6">
                                        <div class="flex-auto">
                                            <div class="flex items-start gap-x-3">
                                                <div class="text-sm font-medium leading-6 text-gray-900">
                                                    {{ conduct.description }}
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="absolute bottom-0 right-full h-px w-screen bg-gray-100" />
                                    <div class="absolute bottom-0 left-0 h-px w-screen bg-gray-100" />
                                </td>
                                <td class="hidden py-5 pr-6 sm:table-cell">
                                    <div class="text-xs leading-5 text-gray-500">{{ conduct.description }}</div>
                                </td>
                                <td class="py-5 text-right">
                                    <div class="flex justify-end text-xs leading-5 text-gray-500">
                                        <Time :value="conduct.createdAt" />
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
