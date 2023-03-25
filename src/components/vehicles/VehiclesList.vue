<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { Vehicle } from '@arpanet/gen/resources/vehicles/vehicles_pb';
import { OrderBy, PaginationRequest } from '@arpanet/gen/resources/common/database/database_pb';
import { watchDebounced } from '@vueuse/core'
import { getDMVClient } from '../../grpc/grpc';
import { FindVehiclesRequest } from '@arpanet/gen/services/dmv/vehicles_pb';
import TablePagination from '../partials/TablePagination.vue';
import VehiclesListEntry from './VehiclesListEntry.vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline';

const props = defineProps({
    userId: {
        type: Number,
        required: false,
        default: 0,
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

const search = ref<{ plate: string, model: string }>({ plate: '', model: '' });
const orderBys = ref<Array<OrderBy>>([]);
const offset = ref(0);
const totalCount = ref(0);
const listEnd = ref(0);
const vehicles = ref<Array<Vehicle>>([]);

function findVehicles(pos: number) {
    if (pos < 0) pos = 0;

    const req = new FindVehiclesRequest();
    req.setPagination((new PaginationRequest()).setOffset(pos));
    if (props.userId && props.userId > 0) {
        req.setUserId(props.userId);
    }
    req.setSearch(search.value.plate);
    req.setModel(search.value.model);
    req.setOrderbyList(orderBys.value);

    getDMVClient().
        findVehicles(req, null).
        then((resp) => {
            const pag = resp.getPagination();
            if (pag !== undefined) {
                totalCount.value = pag.getTotalCount();
                offset.value = pag.getOffset();
                listEnd.value = pag.getEnd();
            }
            vehicles.value = resp.getVehiclesList();
        });
}

function toggleOrderBy(column: string): void {
    const index = orderBys.value.findIndex((o) => {
        return o.getColumn() == column;
    });
    let orderBy: OrderBy;
    if (index > -1) {
        //@ts-ignore I just checked if it exists, so it should exist
        orderBy = orderBys.at(index);
        if (orderBy.getDesc()) {
            orderBys.value.splice(index);
        }
        else {
            orderBy.setDesc(true);
        }
    }
    else {
        orderBy = new OrderBy();
        orderBy.setColumn(column);
        orderBy.setDesc(false);
        orderBys.value.push(orderBy);
    }
    findVehicles(offset.value);
}

const searchInput = ref<HTMLInputElement | null>(null);

function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

onMounted(() => {
    if (props.userId) findVehicles(offset.value);
});

watchDebounced(search.value, () => findVehicles(offset.value), { debounce: 650, maxWait: 1500 });
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="findVehicles(0)">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">License Plate</label>
                                <div class="relative flex items-center mt-2">
                                    <input v-model="search.plate" ref="searchInput" type="text"
                                        placeholder="License plate"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="flex-1 form-control">
                                <label for="model" class="block text-sm font-medium leading-6 text-neutral">Model</label>
                                <div class="relative flex items-center mt-2">
                                    <input v-model="search.model" ref="searchInput" type="text" name="model" id="model"
                                        placeholder="Model"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <button v-if="vehicles.length == 0" type="button" @click="focusSearch()"
                            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
                            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                            <span class="block mt-2 text-sm font-semibold text-gray-300">Use the search field
                                above to search or update your query</span>
                        </button>
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th v-on:click="toggleOrderBy('plate')" scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                                            Plate
                                        </th>
                                        <th v-on:click="toggleOrderBy('model')" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">Model
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Type
                                        </th>
                                        <th v-if="!hideOwner" scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Owner
                                        </th>
                                        <th v-if="!hideOwner" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Job
                                        </th>
                                        <th v-if="!hideCitizenLink" scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Action
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <VehiclesListEntry v-for="vehicle in vehicles" :key="vehicle.getPlate()"
                                        :vehicle="vehicle" :hide-owner="hideOwner"
                                        :hide-citizen-link="hideCitizenLink" :hide-copy="hideCopy"
                                        class="transition-colors hover:bg-neutral/5" />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th v-on:click="toggleOrderBy('plate')" scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                                            Plate
                                        </th>
                                        <th v-on:click="toggleOrderBy('model')" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">Model
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Type
                                        </th>
                                        <th v-if="!hideOwner" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">Owner
                                        </th>
                                        <th v-if="!hideOwner" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Job
                                        </th>
                                        <th v-if="!hideCitizenLink" scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Action
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination :offset="offset" :entries="vehicles.length" :end="listEnd" :total="totalCount"
                                :callback="findVehicles" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
