<script lang="ts" setup>
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import ListEntry from './ListEntry.vue';

defineProps<{
    dispatches?: Dispatch[] | null;
    units: Unit[] | null;
}>();

defineEmits<{
    (e: 'goto', location: { x: number; y: number }): void;
}>();
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
                                <th
                                    scope="col"
                                    class="whitespace-nowrap py-2 pl-4 pr-3 text-left text-sm font-semibold text-gray-100 sm:pl-0"
                                >
                                    Actions
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
                                    class="whitespace-nowrap px-2 py-2 text-left text-sm font-semibold text-gray-100"
                                >
                                    Created At
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
                            <ListEntry
                                v-for="dispatch in dispatches"
                                :key="dispatch.id.toString()"
                                :dispatch="dispatch"
                                :units="units"
                                @goto="$emit('goto', $event)"
                            />
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
