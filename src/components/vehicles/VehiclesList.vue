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
import VehiclesListEntry from '~/components/vehicles/VehiclesListEntry.vue';

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

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
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
                        <div class="mx-auto flex flex-row gap-4">
                            <div class="form-control flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.license_plate') }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        ref="searchInput"
                                        v-model="query.plate"
                                        type="text"
                                        :placeholder="$t('common.license_plate')"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div class="form-control flex-1">
                                <label for="model" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.model') }}
                                </label>
                                <div class="relative mt-2 flex items-center">
                                    <input
                                        v-model="query.model"
                                        type="text"
                                        name="model"
                                        :placeholder="$t('common.model')"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                            <div v-if="!userId" class="form-control flex-1">
                                <label for="owner" class="block text-sm font-medium leading-6 text-neutral">
                                    {{ $t('common.owner') }}
                                </label>
                                <div class="relative mt-2 items-center">
                                    <Combobox v-model="selectedChar" as="div" nullable>
                                        <div class="relative">
                                            <ComboboxButton as="div">
                                                <ComboboxInput
                                                    autocomplete="off"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    :display-value="
                                                        (char: any) =>
                                                            char
                                                                ? `${char?.firstname} ${char?.lastname} (${char?.dateofbirth})`
                                                                : ''
                                                    "
                                                    :placeholder="$t('common.owner')"
                                                    @change="queryChar = $event.target.value"
                                                    @focusin="focusTablet(true)"
                                                    @focusout="focusTablet(false)"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                v-if="chars !== null && chars.length > 0"
                                                class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="char in chars"
                                                    :key="char?.userId"
                                                    v-slot="{ active, selected }"
                                                    :value="char"
                                                    as="char"
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
                                                            ({{ char?.dateofbirth }})
                                                        </span>

                                                        <span
                                                            v-if="selected"
                                                            :class="[
                                                                active ? 'text-neutral' : 'text-primary-500',
                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                            ]"
                                                        >
                                                            <CheckIcon class="h-5 w-5" aria-hidden="true" />
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
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
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
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1"
                                        >
                                            {{ $t('common.plate') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.model') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.type') }}
                                        </th>
                                        <th
                                            v-if="!hideOwner"
                                            scope="col"
                                            class="px-2 py-3.5 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.owner') }}
                                        </th>
                                        <th
                                            v-if="!hideCitizenLink && !hideCopy"
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold text-neutral sm:pr-0"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <VehiclesListEntry
                                        v-for="vehicle in data?.vehicles"
                                        :key="vehicle.plate"
                                        :vehicle="vehicle"
                                        :hide-owner="hideOwner"
                                        :hide-citizen-link="hideCitizenLink"
                                        :hide-copy="hideCopy"
                                    />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1"
                                        >
                                            {{ $t('common.plate') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.model') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.type') }}
                                        </th>
                                        <th
                                            v-if="!hideOwner"
                                            scope="col"
                                            class="px-2 py-3.5 text-left text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.owner') }}
                                        </th>
                                        <th
                                            v-if="!hideCitizenLink && !hideCopy"
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold text-neutral sm:pr-0"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                            </table>

                            <TablePagination
                                :pagination="data?.pagination"
                                :refresh="refresh"
                                @offset-change="offset = $event"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
