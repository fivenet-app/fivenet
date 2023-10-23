<script lang="ts" setup>
import { useCentrumStore } from '~/store/centrum';
import ListEntry from './ListEntry.vue';

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const centrumStore = useCentrumStore();
const { getSortedDispatches } = storeToRefs(centrumStore);
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-auto">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto inline-flex items-center">
                <h2 class="flex-1 text-base font-semibold leading-6 text-gray-100">
                    {{ $t('common.dispatch', 2) }}
                </h2>
                <h2 class="text-base font-semibold text-gray-100">
                    {{ $t('components.centrum.livemap.total_dispatches') }}: {{ getSortedDispatches.length }}
                </h2>
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
                                v-for="dispatch in getSortedDispatches"
                                :key="dispatch.id.toString()"
                                :dispatch="dispatch"
                                @goto="$emit('goto', $event)"
                            />
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
