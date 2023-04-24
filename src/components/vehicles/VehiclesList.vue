<script lang="ts" setup>
import { ref, watch } from 'vue';
import { Vehicle } from '@fivenet/gen/resources/vehicles/vehicles_pb';
import { OrderBy, PaginationRequest, PaginationResponse } from '@fivenet/gen/resources/common/database/database_pb';
import { watchDebounced } from '@vueuse/core'
import { FindVehiclesRequest } from '@fivenet/gen/services/dmv/vehicles_pb';
import TablePagination from '~/components/partials/TablePagination.vue';
import VehiclesListEntry from './VehiclesListEntry.vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline';
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { UserShort } from '@fivenet/gen/resources/users/users_pb';
import { CompleteCharNamesRequest } from '@fivenet/gen/services/completor/completor_pb';
import { CheckIcon } from '@heroicons/vue/20/solid';
import { RpcError } from 'grpc-web';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';

const { $grpc } = useNuxtApp();

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

const entriesChars = ref<UserShort[]>([]);
const queryChar = ref('');
const selectedChar = ref<undefined | UserShort>(undefined);

async function findChars(): Promise<void> {
    if (queryChar.value === '') {
        return;
    }

    const req = new CompleteCharNamesRequest();
    req.setSearch(queryChar.value);

    const resp = await $grpc.getCompletorClient().
        completeCharNames(req, null);

    entriesChars.value = resp.getUsersList();
}

const search = ref<{ plate: string, model: string, user_id: number }>({ plate: '', model: '', user_id: 0 });
const orderBys = ref<Array<OrderBy>>([]);
const pagination = ref<PaginationResponse>();
const offset = ref(0);

const { data: vehicles, pending, refresh, error } = await useLazyAsyncData(`vehicles-${offset.value}`, () => findVehicles());

async function findVehicles(): Promise<Array<Vehicle>> {
    return new Promise(async (res, rej) => {
        const req = new FindVehiclesRequest();
        req.setPagination((new PaginationRequest()).setOffset(offset.value));
        if (props.userId && props.userId > 0) {
            req.setUserId(props.userId);
        } else {
            req.setUserId(search.value.user_id);
        }
        req.setSearch(search.value.plate);
        req.setModel(search.value.model);
        req.setOrderbyList(orderBys.value);

        try {
            const resp = await $grpc.getDMVClient().
                findVehicles(req, null);

            pagination.value = resp.getPagination();
            return res(resp.getVehiclesList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function toggleOrderBy(column: string): Promise<void> {
    const index = orderBys.value.findIndex((o) => {
        return o.getColumn() == column;
    });
    let orderBy: OrderBy;
    if (index > -1) {
        //@ts-ignore I just checked if it exists, so it should exist
        orderBy = orderBys.value.at(index);
        if (orderBy.getDesc()) {
            orderBys.value.splice(index);
        }
        else {
            orderBy.setDesc(true);
        }
    } else {
        orderBy = new OrderBy();
        orderBy.setColumn(column);
        orderBy.setDesc(false);
        orderBys.value.push(orderBy);
    }

    return refresh();
}

const searchInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

watch(offset, async () => refresh());
watchDebounced(search.value, async () => refresh(), { debounce: 600, maxWait: 1400 });
watchDebounced(queryChar, async () => await findChars(), { debounce: 600, maxWait: 1250 });
watch(selectedChar, () => {
    if (selectedChar && selectedChar.value?.getUserId()) {
        search.value.user_id = selectedChar.value?.getUserId();
    } else {
        search.value.user_id = 0;
    }
});
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="findVehicles()">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.license_plate') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input v-model="search.plate" ref="searchInput" type="text" :placeholder="$t('common.license_plate')"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="flex-1 form-control">
                                <label for="model" class="block text-sm font-medium leading-6 text-neutral">{{ $t('common.model') }}</label>
                                <div class="relative flex items-center mt-2">
                                    <input v-model="search.model" type="text" name="model" id="model" :placeholder="$t('common.model')"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="flex-1 form-control" v-if="!userId">
                                <label for="owner" class="block text-sm font-medium leading-6 text-neutral">{{ $t('common.owner') }}</label>
                                <div class="relative items-center mt-2">
                                    <Combobox as="div" v-model="selectedChar" nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    @change="queryChar = $event.target.value"
                                                    :display-value="(char: any) => char ? `${char?.getFirstname()} ${char?.getLastname()}` : ''"
                                                    :placeholder="$t('common.owner')" />
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
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.vehicle', 2)])" />
                        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.vehicle', 2)])" :retry="refresh" />
                        <button v-else-if="vehicles && vehicles.length == 0" type="button" @click="focusSearch()"
                            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
                            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                            <span class="block mt-2 text-sm font-semibold text-gray-300">
                                {{ $t('common.not_found', [$t('common.vehicle', 2)]) }}
                            </span>
                        </button>
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th v-on:click="toggleOrderBy('plate')" scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                                            {{ $t('common.plate') }}
                                        </th>
                                        <th v-on:click="toggleOrderBy('model')" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.model') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.type') }}
                                        </th>
                                        <th v-if="!hideOwner" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.owner') }}
                                        </th>
                                        <th v-if="!hideCitizenLink && !hideCopy" scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <VehiclesListEntry v-for="vehicle in vehicles" :key="vehicle.getPlate()"
                                        :vehicle="vehicle" :hide-citizen-link="hideCitizenLink"
                                        :hide-copy="hideCopy" class="transition-colors hover:bg-neutral/5" />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th v-on:click="toggleOrderBy('plate')" scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                                            {{ $t('common.plate') }}
                                        </th>
                                        <th v-on:click="toggleOrderBy('model')" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.model') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.type') }}
                                        </th>
                                        <th v-if="!hideOwner" scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.owner') }}
                                        </th>
                                        <th v-if="!hideCitizenLink && !hideCopy" scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination :pagination="pagination" @offset-change="offset = $event" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
