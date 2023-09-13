<script lang="ts" setup>
import { useCentrumStore } from '~/store/centrum';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import ListEntry from './ListEntry.vue';

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
    (e: 'details', dsp: Dispatch): void;
    (e: 'assignUnit', dsp: Dispatch): void;
    (e: 'status', dsp: Dispatch): void;
}>();

const centrumStore = useCentrumStore();
const { dispatches } = storeToRefs(centrumStore);
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-auto">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h2 class="text-base font-semibold leading-6 text-gray-100">{{ $t('common.dispatch', 2) }}</h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 -my-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <table class="min-w-full divide-y divide-base-600">
                        <thead>
                            <tr>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.action', 2) }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.id') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.created_at') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.status') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.postal') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.units') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.citizen') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.message') }}
                                </th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-base-800">
                            <ListEntry
                                v-for="[id, dispatch] in dispatches"
                                :key="id.toString()"
                                :dispatch="dispatch"
                                @goto="$emit('goto', $event)"
                                @details="$emit('details', $event)"
                                @assign-unit="$emit('assignUnit', $event)"
                                @status="$emit('status', $event)"
                            />
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
