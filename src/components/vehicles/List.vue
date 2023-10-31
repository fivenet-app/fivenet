<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { watchDebounced } from '@vueuse/core';
import { useRouteHash } from '@vueuse/router';
import { CarSearchIcon, CheckIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import { useCompletorStore } from '~/store/completor';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { ListVehiclesResponse } from '~~/gen/ts/services/dmv/vehicles';
import ListEntry from '~/components/vehicles/ListEntry.vue';

const { $grpc } = useNuxtApp();

const props = withDefaults(
    defineProps<{
        userId?: number;
        hideOwner?: boolean;
        hideCitizenLink?: boolean;
        hideCopy?: boolean;
    }>(),
    {
        userId: 0,
        hideOwner: false,
        hideCitizenLink: false,
        hideCopy: false,
    },
);

const query = ref<{ plate: string; model?: string; user_id?: number }>({
    plate: '',
});
const offset = ref(0n);

const hash = useRouteHash();
if (!props.hideOwner) {
    if (hash.value !== undefined && hash.value !== null) {
        query.value = unmarshalHashToObject(hash.value as string);
    }
}

const { data, pending, refresh, error } = useLazyAsyncData(`vehicles-${offset.value}`, () => {
    if (!props.hideOwner) {
        hash.value = marshalObjectToHash(query.value);
    }
    return listVehicles();
});

async function listVehicles(): Promise<ListVehiclesResponse> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDMVClient().listVehicles({
                pagination: {
                    offset: offset.value,
                },
                orderBy: [],
                userId: props.userId && props.userId > 0 ? props.userId : query.value.user_id,
                search: query.value.plate,
                model: query.value.model,
            });
            const { response } = await call;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            throw e;
        }
    });
}

const searchInput = ref<HTMLInputElement | null>(null);
function focusSearch(): void {
    if (searchInput.value) {
        searchInput.value.focus();
    }
}

const completorStore = useCompletorStore();

const queryChar = ref('');
const selectedChar = ref<undefined | UserShort>(undefined);

const { data: chars, refresh: charsRefresh } = useLazyAsyncData(
    `chars-${queryChar.value}`,
    () =>
        completorStore.completeCitizens({
            search: queryChar.value,
        }),
    {
        immediate: false,
    },
);

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), {
    debounce: 600,
    maxWait: 1400,
});
watchDebounced(queryChar, async () => charsRefresh(), {
    debounce: 600,
    maxWait: 1250,
});
watch(selectedChar, () => {
    if (selectedChar && selectedChar.value?.userId) {
        query.value.user_id = selectedChar.value?.userId;
    } else {
        query.value.user_id = 0;
    }
});
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="flex flex-row gap-4 mx-auto">
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.license_plate') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input
                                        v-model="query.plate"
                                        ref="searchInput"
                                        type="text"
                                        :placeholder="$t('common.license_plate')"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="flex-1 form-control">
                                <label for="model" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.model') }}
                                </label>
                                <div class="relative flex items-center mt-2">
                                    <input
                                        v-model="query.model"
                                        type="text"
                                        name="model"
                                        :placeholder="$t('common.model')"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="flex-1 form-control" v-if="!userId">
                                <label for="owner" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.owner') }}
                                </label>
                                <div class="relative items-center mt-2">
                                    <Combobox as="div" v-model="selectedChar" nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    autocomplete="off"
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    @change="queryChar = $event.target.value"
                                                    :display-value="
                                                        (char: any) => (char ? `${char?.firstname} ${char?.lastname}` : '')
                                                    "
                                                    :placeholder="$t('common.owner')"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                v-if="chars !== null && chars.length > 0"
                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="char in chars"
                                                    :key="char?.userId"
                                                    :value="char"
                                                    as="char"
                                                    v-slot="{ active, selected }"
                                                >
                                                    <li
                                                        :class="[
                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                            active ? 'bg-primary-500' : '',
                                                        ]"
                                                    >
                                                        <span :class="['block truncate', selected && 'font-semibold']">
                                                            {{ char?.firstname }}
                                                            {{ char?.lastname }}
                                                        </span>

                                                        <span
                                                            v-if="selected"
                                                            :class="[
                                                                active ? 'text-neutral' : 'text-primary-500',
                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                            ]"
                                                        >
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
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full py-2 align-middle px-1">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.vehicle', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.vehicle', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data?.vehicles.length === 0"
                            :icon="CarSearchIcon"
                            :focus="focusSearch"
                            :type="$t('common.vehicle', 2)"
                        />
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0"
                                        >
                                            {{ $t('common.plate') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.model') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.type') }}
                                        </th>
                                        <th
                                            v-if="!hideOwner"
                                            scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.owner') }}
                                        </th>
                                        <th
                                            v-if="!hideCitizenLink && !hideCopy"
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <ListEntry
                                        v-for="vehicle in data?.vehicles"
                                        :key="vehicle.plate"
                                        :vehicle="vehicle"
                                        :hide-owner="hideOwner"
                                        :hide-citizen-link="hideCitizenLink"
                                        :hide-copy="hideCopy"
                                        class="transition-colors hover:bg-neutral/5"
                                    />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0"
                                        >
                                            {{ $t('common.plate') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.model') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.type') }}
                                        </th>
                                        <th
                                            v-if="!hideOwner"
                                            scope="col"
                                            class="py-3.5 px-2 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.owner') }}
                                        </th>
                                        <th
                                            v-if="!hideCitizenLink && !hideCopy"
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination
                                :pagination="data?.pagination"
                                @offset-change="offset = $event"
                                :refresh="refresh"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
