<script lang="ts" setup>
import type { TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import LicensePlate from '~/components/partials/LicensePlate.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useClipboardStore } from '~/stores/clipboard';
import { getVehiclesVehiclesClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import type { ListVehiclesResponse } from '~~/gen/ts/services/vehicles/vehicles';
import ColleagueName from '../jobs/colleagues/ColleagueName.vue';
import SelectMenu from '../partials/SelectMenu.vue';
import VehicleInfoPopover from './VehicleInfoPopover.vue';

const { t } = useI18n();

const completorStore = useCompletorStore();

const vehiclesVehiclesClient = await getVehiclesVehiclesClient();

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

    sorting: z
        .object({
            columns: z.custom<SortByColumn>().array().max(3).default([]),
        })
        .default({
            columns: [
                {
                    id: 'plate',
                    desc: false,
                },
            ],
        }),
    page: pageNumberSchema,
});

const query = useSearchForm('vehicles', schema);

const hideVehicleModell = ref(false);

const { data, status, refresh, error } = useLazyAsyncData(
    () => `vehicles-${JSON.stringify(query.sorting)}-${query.page}`,
    () => listVehicles(),
);

async function listVehicles(): Promise<ListVehiclesResponse> {
    try {
        const call = vehiclesVehiclesClient.listVehicles({
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sorting,
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

watchDebounced(query, async () => refresh(), {
    debounce: 200,
    maxWait: 1250,
});

function addToClipboard(vehicle: Vehicle): void {
    clipboardStore.addVehicle(vehicle);

    notifications.add({
        title: { key: 'notifications.clipboard.vehicle_added.title', parameters: {} },
        description: { key: 'notifications.clipboard.vehicle_added.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

function updateVehicle(plate: string, vehicle: Vehicle): void {
    const index = data.value?.vehicles.findIndex((v) => v.plate === plate);
    if (index !== undefined && index >= 0 && data.value?.vehicles[index]) {
        data.value.vehicles[index] = vehicle;
    }
}

const UBadge = resolveComponent('UBadge');
const UButton = resolveComponent('UButton');
const appConfig = useAppConfig();

const columns = computed(() =>
    (
        [
            {
                accessorKey: 'plate',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.plate'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) =>
                    h('div', { class: 'inline-flex items-center gap-1' }, [
                        h(LicensePlate, { plate: row.original.plate, class: 'sm:min-w-40 md:min-w-48' }),
                    ]),
            },
            attr('vehicles.VehiclesService/ListVehicles', 'Fields', 'Wanted').value
                ? {
                      accessorKey: 'wanted',
                      header: t('common.wanted'),
                      cell: ({ row }) =>
                          row.original.props?.wanted
                              ? h(UBadge, { color: 'error' }, () => $t('common.wanted').toUpperCase())
                              : undefined,
                  }
                : undefined,
            {
                accessorKey: 'model',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.model'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
            },
            {
                accessorKey: 'type',
                header: t('common.type'),
                cell: ({ row }) => toTitleCase(row.original.type),
            },
            !props.hideOwner
                ? {
                      accessorKey: 'owner',
                      header: t('common.owner'),
                  }
                : undefined,
            {
                id: 'actions',
            },
        ] as TableColumn<Vehicle>[]
    ).flatMap((item) => (item !== undefined ? [item] : [])),
);

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.inputRef?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm class="my-2 flex w-full flex-row gap-2" :schema="schema" :state="query" @submit="refresh()">
                <UFormField class="flex-1" name="licensePlate" :label="$t('common.license_plate')">
                    <UInput
                        ref="input"
                        v-model="query.licensePlate"
                        type="text"
                        name="licensePlate"
                        :placeholder="$t('common.license_plate')"
                        block
                        leading-icon="i-mdi-search"
                        class="w-full"
                    >
                        <template #trailing>
                            <UKbd value="/" />
                        </template>
                    </UInput>
                </UFormField>

                <UFormField v-if="!hideVehicleModell" class="flex-1" name="model" :label="$t('common.model')">
                    <UInput v-model="query.model" type="text" name="model" :placeholder="$t('common.model')" class="w-full" />
                </UFormField>

                <UFormField v-if="userId === undefined" class="flex-1" name="userIds" :label="$t('common.owner')">
                    <SelectMenu
                        v-model="query.userIds"
                        name="userIds"
                        multiple
                        :searchable="
                            async (q: string): Promise<UserShort[]> =>
                                await completorStore.completeCitizens({
                                    search: q,
                                    userIds: query.userIds,
                                })
                        "
                        searchable-key="completor-citizens"
                        :filter-fields="['firstname', 'lastname']"
                        class="w-full"
                        :placeholder="$t('common.owner')"
                        trailing
                        value-key="userId"
                    >
                        <template #item-label="{ item }">
                            <ColleagueName class="truncate" :colleague="item" birthday />
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.owner', 2)]) }} </template>
                    </SelectMenu>
                </UFormField>

                <UFormField
                    v-if="attr('vehicles.VehiclesService/ListVehicles', 'Fields', 'Wanted').value"
                    class="flex flex-initial flex-col"
                    name="wanted"
                    :label="$t('common.only_wanted')"
                    :ui="{ container: 'flex-1 flex' }"
                >
                    <div class="flex flex-1 items-center">
                        <USwitch v-model="query.wanted" />
                    </div>
                </UFormField>
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
        v-model:sorting="query.sorting.columns"
        class="flex-1"
        :loading="isRequestPending(status)"
        :columns="columns"
        :data="data?.vehicles"
        :pagination-options="{ manualPagination: true }"
        :sorting-options="{ manualSorting: true }"
        :empty="$t('common.not_found', [$t('common.vehicle', 2)])"
        sticky
    >
        <template v-if="!hideOwner" #owner-cell="{ row: vehicle }">
            <p v-if="vehicle.original.jobLabel" class="text-highlighted">{{ vehicle.original.jobLabel }}</p>
            <CitizenInfoPopover v-if="vehicle.original.owner" :user="vehicle.original.owner" />
        </template>

        <template #actions-cell="{ row: vehicle }">
            <div class="flex flex-col justify-end md:flex-row">
                <UTooltip
                    v-if="attrStringList('vehicles.VehiclesService/ListVehicles', 'Fields').value.length > 0 || isSuperuser"
                    :text="$t('common.propertie', 2)"
                >
                    <VehicleInfoPopover
                        :model-value="vehicle.original"
                        @update:model-value="updateVehicle(vehicle.original.plate, $event)"
                    />
                </UTooltip>

                <UTooltip v-if="!hideCopy" :text="$t('components.clipboard.clipboard_button.add')">
                    <UButton variant="link" icon="i-mdi-clipboard-plus" @click="addToClipboard(vehicle.original)" />
                </UTooltip>

                <UTooltip
                    v-if="
                        !hideCitizenLink && vehicle.original.owner?.userId && can('citizens.CitizensService/ListCitizens').value
                    "
                    :text="$t('common.show')"
                >
                    <UButton
                        variant="link"
                        icon="i-mdi-account-eye"
                        :to="{
                            name: 'citizens-id',
                            params: { id: vehicle.original.owner.userId },
                        }"
                    />
                </UTooltip>
            </div>
        </template>
    </UTable>

    <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
</template>
