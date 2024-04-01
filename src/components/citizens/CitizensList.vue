<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { watchDebounced } from '@vueuse/core';
import { vMaska } from 'maska';
import { ChevronDownIcon, ClipboardPlusIcon, EyeIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { attr } from '~/composables/can';
import GenericInput from '~/composables/partials/forms/GenericInput.vue';
import { ListCitizensRequest, ListCitizensResponse } from '~~/gen/ts/services/citizenstore/citizenstore';
import type { User } from '~~/gen/ts/resources/users/users';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import PhoneNumberBlock from '../partials/citizens/PhoneNumberBlock.vue';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const query = ref<{
    name?: string;
    phoneNumber?: string;
    wanted?: boolean;
    trafficInfractionPoints?: number;
    fines?: number;
    dateofbirth?: string;
}>({});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`citizens-${page.value}-${jsonStringify(query.value)}`, () =>
    listCitizens(),
);

async function listCitizens(): Promise<ListCitizensResponse> {
    try {
        const req: ListCitizensRequest = {
            pagination: {
                offset: offset.value,
            },
            searchName: query.value.name ?? '',
        };
        if (query.value.wanted) {
            req.wanted = query.value.wanted;
        }
        if (query.value.phoneNumber) {
            req.phoneNumber = query.value.phoneNumber;
        }
        if (query.value.trafficInfractionPoints) {
            req.trafficInfractionPoints = query.value.trafficInfractionPoints ?? 0;
        }
        if (query.value.fines) {
            req.openFines = BigInt(query.value.fines?.toString() ?? '0');
        }
        if (query.value.dateofbirth) {
            req.dateofbirth = query.value.dateofbirth;
        }

        const call = $grpc.getCitizenStoreClient().listCitizens(req);
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

watch(offset, async () => refresh());
watchDebounced(query.value, () => refresh(), { debounce: 600, maxWait: 1400 });

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

function addToClipboard(user: User): void {
    clipboardStore.addUser(user);

    notifications.add({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        description: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        timeout: 3250,
        type: 'info',
    });
}

const columns = [
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'jobLabel',
        label: t('common.job'),
    },
    {
        key: 'sex',
        label: t('common.sex'),
    },
    attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')
        ? { key: 'phoneNumber', label: t('common.phone_number') }
        : undefined,
    {
        key: 'dateofbirth',
        label: t('common.date_of_birth'),
    },
    attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints')
        ? {
              key: 'trafficInfractionPoints',
              label: t('common.traffic_infraction_points'),
          }
        : undefined,
    attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines')
        ? {
              key: 'openFines',
              label: t('common.fine', 2),
          }
        : undefined,
    {
        key: 'height',
        label: t('common.height'),
    },
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
].filter((c) => c !== undefined);
</script>

<template>
    <div>
        <UDashboardToolbar>
            <template #default>
                <form class="w-full" @submit.prevent="refresh()">
                    <div class="w-full flex flex-row gap-2">
                        <div>
                            <label for="searchName" class="block text-sm font-medium leading-6 text-neutral">
                                {{ $t('common.search') }}
                                {{ $t('common.citizen', 1) }}
                            </label>
                            <div class="relative mt-2">
                                <GenericInput
                                    ref="searchInput"
                                    v-model="query.name"
                                    type="text"
                                    name="searchName"
                                    :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                                    block
                                />
                            </div>
                        </div>
                        <div>
                            <label for="dateofbirth" class="block text-sm font-medium leading-6 text-neutral">
                                {{ $t('common.search') }}
                                {{ $t('common.date_of_birth') }}
                            </label>
                            <div class="relative mt-2">
                                <UInput
                                    v-model="query.dateofbirth"
                                    v-maska
                                    type="text"
                                    name="dateofbirth"
                                    data-maska="##.##.####"
                                    :placeholder="`${$t('common.date_of_birth')} (DD.MM.YYYY)`"
                                    block
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </div>
                        </div>
                        <div v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.Wanted')" class="flex-initial">
                            <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                {{ $t('components.citizens.citizens_list.only_wanted') }}
                            </label>
                            <div class="relative mt-3 flex items-center">
                                <UToggle v-model="query.wanted">
                                    <span class="sr-only">
                                        {{ $t('components.citizens.citizens_list.only_wanted') }}
                                    </span>
                                </UToggle>
                            </div>
                        </div>
                    </div>
                    <Disclosure v-slot="{ open }" as="div" class="pt-2">
                        <DisclosureButton class="flex w-full items-start justify-between text-left text-sm text-neutral">
                            <span class="leading-7 text-accent-200">{{ $t('common.advanced_search') }}</span>
                            <span class="ml-6 flex h-7 items-center">
                                <ChevronDownIcon
                                    :class="[open ? 'upsidedown' : '', 'size-5 transition-transform']"
                                    aria-hidden="true"
                                />
                            </span>
                        </DisclosureButton>
                        <DisclosurePanel class="mt-2 pr-4">
                            <div class="flex flex-row gap-2">
                                <div v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')">
                                    <label for="searchPhone" class="block text-sm font-medium leading-6 text-neutral">
                                        {{ $t('common.search') }}
                                        {{ $t('common.phone_number') }}
                                    </label>
                                    <div class="relative mt-2">
                                        <UInput
                                            v-model="query.phoneNumber"
                                            type="tel"
                                            name="searchPhone"
                                            :placeholder="$t('common.phone_number')"
                                            block
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </div>
                                </div>
                                <div>
                                    <label
                                        for="trafficInfractionPoints"
                                        class="block text-sm font-medium leading-6 text-neutral"
                                    >
                                        {{ $t('common.search') }}
                                        {{ $t('common.traffic_infraction_points', 2) }}
                                    </label>
                                    <div class="relative mt-2">
                                        <UInput
                                            v-model="query.trafficInfractionPoints"
                                            type="number"
                                            name="trafficInfractionPoints"
                                            :placeholder="$t('common.traffic_infraction_points')"
                                            block
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </div>
                                </div>
                                <div v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines')">
                                    <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                        {{ $t('components.citizens.citizens_list.open_fine') }}
                                    </label>
                                    <div class="relative mt-2">
                                        <UInput
                                            v-model="query.fines"
                                            type="number"
                                            name="fine"
                                            :placeholder="`${$t('common.fine')}`"
                                            block
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                    </div>
                                </div>
                            </div>
                        </DisclosurePanel>
                    </Disclosure>
                </form>
            </template>
        </UDashboardToolbar>

        <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.citizen', 2)])" :retry="refresh" />
        <UTable
            :loading="pending"
            :columns="columns"
            :rows="data?.users"
            :empty-state="{ icon: 'i-mdi-car', label: $t('common.not_found', [$t('common.citizen', 2)]) }"
            :page-count="(data?.pagination?.totalCount ?? 0) / (data?.pagination?.pageSize ?? 1)"
            :total="data?.pagination?.totalCount"
        >
            <template #name-data="{ row }">
                <span>{{ row.firstname }} {{ row.lastname }}</span>
                <span class="lg:hidden"> ({{ row.dateofbirth }}) </span>

                <span
                    v-if="row.props?.wanted"
                    class="ml-1 rounded-md bg-error-100 px-2 py-0.5 text-sm font-medium text-error-700"
                >
                    {{ $t('common.wanted').toUpperCase() }}
                </span>

                <dl class="font-normal lg:hidden">
                    <dt class="sr-only">{{ $t('common.sex') }} - {{ $t('common.job') }}</dt>
                    <dd class="mt-1 truncate text-accent-200">{{ row.sex!.toUpperCase() }} - {{ row.jobLabel }}</dd>
                </dl>
            </template>
            <template #jobLabel-data="{ row }">
                {{ row.jobLabel }}
            </template>
            <template #sex-data="{ row }">
                {{ row.sex!.toUpperCase() }}
            </template>
            <template #phoneNumber-data="{ row }">
                <PhoneNumberBlock :number="row.phoneNumber" />
            </template>
            <template #openFines-data="{ row }">
                <template v-if="(row.props?.openFines ?? 0n) > 0n">
                    {{ $n(parseInt((row?.props?.openFines ?? 0n).toString()), 'currency') }}
                </template>
            </template>
            <template #height-data="{ row }"> {{ row.height }}cm </template>
            <template #actions-data="{ row }">
                <div v-if="can('CitizenStoreService.GetUser')" class="flex flex-col justify-end gap-1 md:flex-row">
                    <span class="flex-initial text-primary-500 hover:text-primary-400" @click="addToClipboard(row)">
                        <ClipboardPlusIcon class="ml-auto mr-2.5 h-auto w-5" aria-hidden="true" />
                    </span>

                    <NuxtLink
                        :to="{
                            name: 'citizens-id',
                            params: { id: row.userId ?? 0 },
                        }"
                        class="flex-initial text-primary-500 hover:text-primary-400"
                    >
                        <EyeIcon class="ml-auto mr-2.5 h-auto w-5" aria-hidden="true" />
                    </NuxtLink>
                </div>
            </template>
        </UTable>

        <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
            <UPagination
                v-model="page"
                :page-count="parseInt(data?.pagination?.pageSize.toString() ?? '0')"
                :total="parseInt(data?.pagination?.totalCount.toString() ?? '0')"
            />
        </div>
    </div>
</template>
