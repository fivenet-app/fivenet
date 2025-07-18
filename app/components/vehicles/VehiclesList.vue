<script lang="ts" setup>
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import LicensePlate from '~/components/partials/LicensePlate.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useClipboardStore } from '~/stores/clipboard';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import type { ListVehiclesResponse } from '~~/gen/ts/services/vehicles/vehicles';
import ColleagueName from '../jobs/colleagues/ColleagueName.vue';
import VehicleInfoPopover from './VehicleInfoPopover.vue';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const props = withDefaults(
    defineProps<{
        userId?: number;
        hideOwner?: boolean;
        hideCitizenLink?: boolean;
        hideCopy?: boolean;
    }>(),
    {
        userId: undefined,
        hideOwner: false,
        hideCitizenLink: false,
        hideCopy: false,
    },
);

const clipboardStore = useClipboardStore();

const notifications = useNotificationsStore();

const { attr, attrStringList, can, isSuperuser } = useAuth();

const schema = z.object({
    licensePlate: z.string().max(32).default(''),
    model: z.string().min(6).max(32).optional(),
    userIds: z.coerce.number().array().max(5).default([]),
    wanted: z.boolean().default(false),

    sort: z.custom<TableSortable>().default({
        column: 'plate',
        direction: 'asc',
    }),
    page: pageNumberSchema,
});

const query = useSearchForm('vehicles', schema);

const hideVehicleModell = ref(false);

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`vehicles-${query.sort.column}:${query.sort.direction}-${query.page}`, () => listVehicles());

async function listVehicles(): Promise<ListVehiclesResponse> {
    try {
        const call = $grpc.vehicles.vehicles.listVehicles({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sort,
            licensePlate: query.licensePlate,
            model: query.model,
            userIds: query.userIds,
            wanted: query.wanted,
        });
        const { response } = await call;

        if (response.vehicles.length > 0) {
            if (response.vehicles[0]?.model === undefined) {
                hideVehicleModell.value = true;
            } else {
                hideVehicleModell.value = false;
            }
        }

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const usersLoading = ref(false);

watchDebounced(query, async () => refresh(), {
    debounce: 200,
    maxWait: 1250,
});

function addToClipboard(vehicle: Vehicle): void {
    clipboardStore.addVehicle(vehicle);

    notifications.add({
        title: { key: 'notifications.clipboard.vehicle_added.title', parameters: {} },
        description: { key: 'notifications.clipboard.vehicle_added.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}

function updateVehicle(plate: string, vehicle: Vehicle): void {
    const index = data.value?.vehicles.findIndex((v) => v.plate === plate);
    if (index !== undefined && index >= 0 && data.value?.vehicles[index]) {
        data.value.vehicles[index] = vehicle;
    }
}

const columns = computed(() =>
    [
        {
            key: 'plate',
            label: t('common.plate'),
            sortable: true,
        },
        attr('vehicles.VehiclesService/ListVehicles', 'Fields', 'Wanted').value
            ? {
                  key: 'wanted',
                  label: t('common.wanted'),
                  sortable: false,
              }
            : undefined,
        {
            key: 'model',
            label: t('common.model'),
            sortable: true,
        },
        {
            key: 'type',
            label: t('common.type'),
        },
        !props.hideOwner
            ? {
                  key: 'owner',
                  label: t('common.owner'),
              }
            : undefined,
        {
            key: 'actions',
            label: t('common.action', 2),
            sortable: false,
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.input?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm class="flex w-full flex-row gap-2" :schema="schema" :state="query" @submit="refresh()">
                <UFormGroup class="flex-1" name="licensePlate" :label="$t('common.license_plate')">
                    <UInput
                        ref="input"
                        v-model="query.licensePlate"
                        type="text"
                        name="licensePlate"
                        :placeholder="$t('common.license_plate')"
                        block
                        leading-icon="i-mdi-search"
                    >
                        <template #trailing>
                            <UKbd value="/" />
                        </template>
                    </UInput>
                </UFormGroup>

                <UFormGroup v-if="!hideVehicleModell" class="flex-1" name="model" :label="$t('common.model')">
                    <UInput v-model="query.model" type="text" name="model" :placeholder="$t('common.model')" block />
                </UFormGroup>

                <UFormGroup v-if="userId === undefined" class="flex-1" name="userIds" :label="$t('common.owner')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="query.userIds"
                            name="userIds"
                            multiple
                            :searchable="
                                async (q: string): Promise<UserShort[]> => {
                                    usersLoading = true;
                                    const { response } = await $grpc.completor.completor.completeCitizens({
                                        search: q,
                                        userIds: query.userIds,
                                    });
                                    usersLoading = false;
                                    return response.users;
                                }
                            "
                            searchable-lazy
                            :search-placeholder="$t('common.search_field')"
                            :search-attributes="['firstname', 'lastname']"
                            block
                            :placeholder="$t('common.owner')"
                            trailing
                            value-attribute="userId"
                        >
                            <template #label="{ selected }">
                                <span v-if="selected.length > 0" class="truncate">
                                    {{ usersToLabel(selected) }}
                                </span>
                            </template>

                            <template #option="{ option: user }">
                                <ColleagueName class="truncate" :colleague="user" birthday />
                            </template>

                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.owner', 2)]) }} </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormGroup>

                <UFormGroup
                    v-if="attr('vehicles.VehiclesService/ListVehicles', 'Fields', 'Wanted').value"
                    class="flex flex-initial flex-col"
                    name="wanted"
                    :label="$t('common.only_wanted')"
                    :ui="{ container: 'flex-1 flex' }"
                >
                    <div class="flex flex-1 items-center">
                        <UToggle v-model="query.wanted" />
                    </div>
                </UFormGroup>
            </UForm>
        </template>
    </UDashboardToolbar>

    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [$t('common.vehicle', 2)])"
        :error="error"
        :retry="refresh"
    />

    <UTable
        v-else
        v-model:sort="query.sort"
        class="flex-1"
        :loading="loading"
        :columns="columns"
        :rows="data?.vehicles"
        :empty-state="{ icon: 'i-mdi-car', label: $t('common.not_found', [$t('common.vehicle', 2)]) }"
        sort-mode="manual"
    >
        <template #plate-data="{ row: vehicle }">
            <div class="inline-flex items-center gap-1">
                <LicensePlate :plate="vehicle.plate" class="sm:min-w-40 md:min-w-48" />
            </div>
        </template>

        <template #wanted-data="{ row: vehicle }">
            <UBadge v-if="vehicle.props?.wanted" color="error">
                {{ $t('common.wanted').toUpperCase() }}
            </UBadge>
        </template>

        <template #type-data="{ row: vehicle }">
            {{ toTitleCase(vehicle.type) }}
        </template>

        <template v-if="!hideOwner" #owner-data="{ row: vehicle }">
            <p v-if="vehicle.jobLabel" class="text-gray-900 dark:text-white">{{ vehicle.jobLabel }}</p>
            <CitizenInfoPopover v-if="vehicle.owner" :user="vehicle.owner" />
        </template>

        <template #actions-data="{ row: vehicle }">
            <div :key="vehicle.plate" class="flex flex-col justify-end md:flex-row">
                <UTooltip
                    v-if="attrStringList('vehicles.VehiclesService/ListVehicles', 'Fields').value.length > 0 || isSuperuser"
                    :text="$t('common.propertie', 2)"
                >
                    <VehicleInfoPopover :model-value="vehicle" @update:model-value="updateVehicle(vehicle.plate, $event)" />
                </UTooltip>

                <UTooltip v-if="!hideCopy" :text="$t('components.clipboard.clipboard_button.add')">
                    <UButton variant="link" icon="i-mdi-clipboard-plus" @click="addToClipboard(vehicle)" />
                </UTooltip>

                <UTooltip
                    v-if="!hideCitizenLink && vehicle.owner?.userId && can('citizens.CitizensService/ListCitizens').value"
                    :text="$t('common.show')"
                >
                    <UButton
                        variant="link"
                        icon="i-mdi-account-eye"
                        :to="{
                            name: 'citizens-id',
                            params: { id: vehicle.owner.userId },
                        }"
                    />
                </UTooltip>
            </div>
        </template>
    </UTable>

    <Pagination v-model="query.page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
