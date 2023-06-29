<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccountMultiplePlus, mdiDetails } from '@mdi/js';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';

const dispatches = [
    {
        id: 1,
        unit: 'David 1',
        status: 'EN ROUTE',
        prio: '',
        postal: '8040',
    },
    {
        id: 2,
        unit: 'Overwatch 1',
        status: 'EN ROUTE',
        prio: '',
        postal: '8040',
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
        <div class="mt-2 flow-root">
            <div class="-mx-2 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <table class="min-w-full divide-y divide-gray-300">
                        <thead>
                            <tr>
                                <th scope="col" class="relative whitespace-nowrap py-3.5 pl-3 pr-4 sm:pr-0">
                                    <span class="sr-only">Actions</span>
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-100 sm:pl-0"
                                >
                                    ID
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-100 sm:pl-0"
                                >
                                    Unit
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-3.5 text-left text-sm font-semibold text-gray-100"
                                >
                                    Status
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-3.5 text-left text-sm font-semibold text-gray-100"
                                >
                                    Prio
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-2 py-3.5 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.postal') }}
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-200">
                            <tr v-for="dispatch in dispatches" :key="dispatch.id">
                                <td
                                    class="relative whitespace-nowrap py-2 pl-2 text-right text-sm font-medium sm:pr-0 max-w-[42px] flex flex-row justify-start"
                                >
                                    <a href="#" class="text-primary-400 hover:text-primary-600" :title="$t('common.detail', 2)">
                                        <SvgIcon type="mdi" :path="mdiDetails" class="w-6 h-auto ml-auto mr-2.5" />
                                    </a>
                                    <a href="#" class="text-primary-400 hover:text-primary-600" :title="$t('common.assign')">
                                        <SvgIcon
                                            type="mdi"
                                            :path="mdiAccountMultiplePlus"
                                            class="w-6 h-auto ml-auto mr-2.5"
                                            aria-hidden="true"
                                        />
                                    </a>
                                </td>
                                <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
                                    <IDCopyBadge :id="dispatch.id.toString()" prefix="DSP" class="w-28" />
                                </td>
                                <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
                                    {{ dispatch.unit }}
                                </td>
                                <td class="whitespace-nowrap px-2 py-2 text-sm font-medium text-gray-100">
                                    {{ dispatch.status }}
                                </td>
                                <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-100">{{ dispatch.prio }}</td>
                                <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">{{ dispatch.postal }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
