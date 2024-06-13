<script lang="ts" setup>
import { z } from 'zod';
import { ListVehiclesResponse } from '~~/gen/ts/services/dmv/vehicles';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import LicensePlate from '~/components/partials/LicensePlate.vue';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

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

const schema = z.object({
    licensePlate: z.string().max(32),
    model: z.string().min(6).max(32).optional(),
    userId: z.number().optional(),
});

type Schema = z.output<typeof schema>;

const query = ref<Schema>({
    licensePlate: '',
    userId: props.userId,
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const hideVehicleModell = ref(false);

const { data, pending: loading, refresh, error } = useLazyAsyncData(`vehicles-${page.value}`, () => listVehicles());

async function listVehicles(): Promise<ListVehiclesResponse> {
    try {
        const call = getGRPCDMVClient().listVehicles({
            pagination: {
                offset: offset.value,
            },
            orderBy: [],
            licensePlate: query.value.licensePlate,
            model: query.value.model,
            userId: query.value.userId,
        });
        const { response } = await call;

        if (response.vehicles.length > 0) {
            if (response.vehicles[0].model === undefined) {
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
const selectedUser = ref<undefined | UserShort>();
watch(selectedUser, () => (query.value.userId = selectedUser.value?.userId));

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), {
    debounce: 200,
    maxWait: 1250,
});

const clipboardStore = useClipboardStore();

const notifications = useNotificatorStore();

function addToClipboard(vehicle: Vehicle): void {
    clipboardStore.addVehicle(vehicle);

    notifications.add({
        title: { key: 'notifications.clipboard.vehicle_added.title', parameters: {} },
        description: { key: 'notifications.clipboard.vehicle_added.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}

const columns = computed(() =>
    [
        {
            key: 'plate',
            label: t('common.plate'),
        },
        {
            key: 'model',
            label: t('common.model'),
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

const input = ref<{ input: HTMLInputElement }>();

defineShortcuts({
    '/': () => {
        input.value?.input?.focus();
    },
});
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm :schema="schema" :state="query" class="flex w-full flex-row gap-2" @submit="refresh()">
                <UFormGroup name="licensePlate" :label="$t('common.license_plate')" class="flex-1">
                    <UInput
                        ref="input"
                        v-model="query.licensePlate"
                        type="text"
                        name="licensePlate"
                        :placeholder="$t('common.license_plate')"
                        block
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    >
                        <template #trailing>
                            <UKbd value="/" />
                        </template>
                    </UInput>
                </UFormGroup>

                <UFormGroup v-if="!hideVehicleModell" name="model" :label="$t('common.model')" class="flex-1">
                    <UInput
                        v-model="query.model"
                        type="text"
                        name="model"
                        :placeholder="$t('common.model')"
                        block
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </UFormGroup>

                <UFormGroup v-if="userId === undefined" name="selectedUser" :label="$t('common.owner')" class="flex-1">
                    <UInputMenu
                        v-model="selectedUser"
                        name="selectedUser"
                        nullable
                        :search="
                            async (query: string): Promise<UserShort[]> => {
                                usersLoading = true;
                                const { response } = await getGRPCCompletorClient().completeCitizens({
                                    search: query,
                                });
                                usersLoading = false;
                                return response.users;
                            }
                        "
                        search-lazy
                        :search-placeholder="$t('common.search_field')"
                        :search-attributes="['firstname', 'lastname']"
                        block
                        :placeholder="$t('common.owner')"
                        trailing
                        by="userId"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    >
                        <template #option="{ option: user }">
                            {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                        </template>
                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>
                        <template #empty> {{ $t('common.not_found', [$t('common.owner', 2)]) }} </template>
                    </UInputMenu>
                </UFormGroup>
            </UForm>
        </template>
    </UDashboardToolbar>

    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.vehicle', 2)])" :retry="refresh" />

    <UTable
        v-else
        :loading="loading"
        :columns="columns"
        :rows="data?.vehicles"
        :empty-state="{ icon: 'i-mdi-car', label: $t('common.not_found', [$t('common.vehicle', 2)]) }"
        class="flex-1"
    >
        <template #plate-data="{ row: vehicle }">
            <LicensePlate :plate="vehicle.plate" class="mr-2" />
        </template>
        <template #type-data="{ row: vehicle }">
            {{ toTitleCase(vehicle.type) }}
        </template>
        <template v-if="!hideOwner" #owner-data="{ row: vehicle }">
            <CitizenInfoPopover :user="vehicle.owner" />
        </template>
        <template #actions-data="{ row: vehicle }">
            <div :key="vehicle.plate" class="flex flex-col justify-end md:flex-row">
                <UButton v-if="!hideCopy" variant="link" icon="i-mdi-clipboard-plus" @click="addToClipboard(vehicle)" />

                <UButton
                    v-if="!hideCitizenLink && can('CitizenStoreService.ListCitizens')"
                    variant="link"
                    icon="i-mdi-account-eye"
                    :to="{
                        name: 'citizens-id',
                        params: { id: vehicle.owner?.userId ?? 0 },
                    }"
                />
            </div>
        </template>
    </UTable>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
