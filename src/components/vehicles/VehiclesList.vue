<script lang="ts" setup>
import { ref, onMounted, watch } from 'vue';
import { Vehicle } from '@arpanet/gen/resources/vehicles/vehicles_pb';
import { OrderBy, PaginationRequest, PaginationResponse } from '@arpanet/gen/resources/common/database/database_pb';
import { watchDebounced } from '@vueuse/core'
import { getCompletorClient, getDMVClient, handleRPCError } from '../../grpc/grpc';
import { FindVehiclesRequest } from '@arpanet/gen/services/dmv/vehicles_pb';
import TablePagination from '../partials/TablePagination.vue';
import VehiclesListEntry from './VehiclesListEntry.vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline';
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { UserShort } from '@arpanet/gen/resources/users/users_pb';
import { CompleteCharNamesRequest } from '@arpanet/gen/services/completor/completor_pb';
import {
    CheckIcon,
} from '@heroicons/vue/20/solid';

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

const search = ref<{ plate: string, model: string, user_id: number }>({ plate: '', model: '', user_id: 0 });
const orderBys = ref<Array<OrderBy>>([]);
const pagination = ref<PaginationResponse>();
const vehicles = ref<Array<Vehicle>>([]);

const entriesChars = ref<UserShort[]>([]);
const queryChar = ref('');
const selectedChar = ref<undefined | UserShort>(undefined);

function findVehicles(pos: number) {
    if (pos < 0) pos = 0;

    const req = new FindVehiclesRequest();
    req.setPagination((new PaginationRequest()).setOffset(pos));
    if (props.userId && props.userId > 0) {
        req.setUserId(props.userId);
    } else {
        req.setUserId(search.value.user_id);
    }
    req.setSearch(search.value.plate);
    req.setModel(search.value.model);
    req.setOrderbyList(orderBys.value);

    getDMVClient().
        findVehicles(req, null).
        then((resp) => {
            const pag = resp.getPagination();
            pagination.value = resp.getPagination();
            vehicles.value = resp.getVehiclesList();
        }).
        catch((err) => handleRPCError(err));
}

async function findChars(): Promise<void> {
    if (queryChar.value === "") {
        return;
    }

    const req = new CompleteCharNamesRequest();
    req.setSearch(queryChar.value);

    const resp = await getCompletorClient().
        completeCharNames(req, null);

    entriesChars.value = resp.getUsersList();
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
    findVehicles(pagination.value?.getOffset()!);
}

const searchInput = ref<HTMLInputElement | null>(null);

function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

onMounted(() => {
    if (props.userId) findVehicles(pagination.value?.getOffset()!);
});

watchDebounced(queryChar, async () => await findChars(), { debounce: 600, maxWait: 1250 });
watch(selectedChar, () => {
    if (selectedChar && selectedChar.value?.getUserId()) {
        search.value.user_id = selectedChar.value?.getUserId();
    } else {
        search.value.user_id = 0;
    }
});
watchDebounced(search.value, () => findVehicles(pagination.value?.getOffset()!), { debounce: 650, maxWait: 1500 });
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="findVehicles(0)">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">License
                                    Plate</label>
                                <div class="relative flex items-center mt-2">
                                    <input v-model="search.plate" ref="searchInput" type="text" placeholder="License plate"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="flex-1 form-control">
                                <label for="model" class="block text-sm font-medium leading-6 text-neutral">Model</label>
                                <div class="relative flex items-center mt-2">
                                    <input v-model="search.model" type="text" name="model" id="model" placeholder="Model"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="flex-1 form-control" v-if="!props.userId">
                                <label for="owner" class="block text-sm font-medium leading-6 text-neutral">Owner</label>
                                <div class="relative items-center mt-2">
                                    <Combobox as="div" v-model="selectedChar" nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    @change="queryChar = $event.target.value"
                                                    :display-value="(char: any) => char ? `${char?.getFirstname()} ${char?.getLastname()}` : ''"
                                                    placeholder="Owner" />
                                            </ComboboxButton>

                                            <ComboboxOptions v-if="entriesChars.length > 0"
                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                                                <ComboboxOption v-for="char in entriesChars" :key="char?.getIdentifier()"
                                                    :value="char" as="char" v-slot="{ active, selected }">
                                                    <li
                                                        :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                                                        <span :class="['block truncate', selected && 'font-semibold']">
                                                            {{ char?.getFirstname() }} {{ char?.getLastname() }}
                                                        </span>

                                                        <span v-if="selected"
                                                            :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                                                            <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                                        </span>
                                                    </li>
                                                </ComboboxOption>
                                            </ComboboxOptions>
                                        </div>
                                    </Combobox>
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
                            <span class="block mt-2 text-sm font-semibold text-gray-300">
                                Use the search field above to search or update your query.
                            </span>
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
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Model
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Type
                                        </th>
                                        <th v-if="!hideOwner" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Owner
                                        </th>
                                        <th v-if="!hideOwner" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Job
                                        </th>
                                        <th v-if="!hideCitizenLink" scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Actions
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <VehiclesListEntry v-for="vehicle in vehicles" :key="vehicle.getPlate()"
                                        :vehicle="vehicle" :hide-owner="hideOwner" :hide-citizen-link="hideCitizenLink"
                                        :hide-copy="hideCopy" class="transition-colors hover:bg-neutral/5" />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th v-on:click="toggleOrderBy('plate')" scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                                            Plate
                                        </th>
                                        <th v-on:click="toggleOrderBy('model')" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Model
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Type
                                        </th>
                                        <th v-if="!hideOwner" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Owner
                                        </th>
                                        <th v-if="!hideOwner" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Job
                                        </th>
                                        <th v-if="!hideCitizenLink && !hideCopy" scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Actions
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination :pagination="pagination!" :callback="findVehicles" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
