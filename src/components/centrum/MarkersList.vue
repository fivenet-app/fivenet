<script lang="ts" setup>
import { useLivemapStore } from '~/store/livemap';
import MarkersListEntry from '~/components/centrum/MarkersListEntry.vue';
import GenericTable from '../partials/GenericTable.vue';

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const livemapStore = useLivemapStore();
const { markersMarkers } = storeToRefs(livemapStore);
</script>

<template>
    <div class="h-full overflow-y-auto px-4 sm:px-6 lg:px-8">
        <div class="sm:flex sm:items-center">
            <div class="inline-flex items-center sm:flex-auto">
                <h2 class="inline-flex flex-1 items-center text-base font-semibold leading-6 text-gray-100">
                    {{ $t('common.marker', 2) }}
                </h2>
                <h2 class="text-base font-semibold text-gray-100">
                    {{ $t('common.count') }}:
                    {{ [...markersMarkers.values()].length }}
                </h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <GenericTable>
                        <template #thead>
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
                                    {{ $t('common.created_at') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.expires_at') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.name') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.type') }}
                                </th>
                                <th
                                    scope="col"
                                    class="whitespace-nowrap px-1 py-1 text-left text-sm font-semibold text-gray-100"
                                >
                                    {{ $t('common.description') }}
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
                                    {{ $t('common.job') }}
                                </th>
                            </tr>
                        </template>
                        <template #tbody>
                            <MarkersListEntry
                                v-for="[key, marker] in markersMarkers"
                                :key="key"
                                :marker="marker"
                                @goto="$emit('goto', $event)"
                            />
                        </template>
                    </GenericTable>
                </div>
            </div>
        </div>
    </div>
</template>
