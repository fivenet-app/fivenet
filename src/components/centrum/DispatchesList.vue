<script lang="ts" setup>
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatch';
import DispatchListEntry from './DispatchListEntry.vue';

const dispatches: Dispatch[] = [
    {
        id: 1n,
        job: 'ambulance',
        message: 'Notruf',
        status: DISPATCH_STATUS.EN_ROUTE,
        units: [],
        attributes: {},
    },
];
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-scroll">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h1 class="text-base font-semibold leading-6 text-gray-100">Dispatches</h1>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <table class="min-w-full divide-y divide-gray-300">
                        <thead>
                            <tr>
                                <th scope="col" class="relative whitespace-nowrap py-2 pl-3 pr-4 sm:pr-0">
                                    <span class="sr-only">Actions</span>
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap py-2 pl-4 pr-3 text-left text-sm font-semibold text-gray-100 sm:pl-0"
                                >
                                    ID
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-2 text-left text-sm font-semibold text-gray-100"
                                >
                                    Status
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap py-2 pl-4 pr-3 text-left text-sm font-semibold text-gray-100 sm:pl-0"
                                >
                                    Units
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-2 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.location') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-2 text-left text-sm font-semibold text-gray-100"
                                >
                                    Message
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-200">
                            <tr v-for="dispatch in dispatches" :key="dispatch.id.toString()">
                                <DispatchListEntry :dispatch="dispatch" />
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
