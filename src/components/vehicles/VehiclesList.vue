<script lang="ts" setup>
import { watchDebounced } from '@vueuse/core';
import { useCompletorStore } from '~/store/completor';
import { ListVehiclesResponse } from '~~/gen/ts/services/dmv/vehicles';
import CitizenInfoPopover from '../partials/citizens/CitizenInfoPopover.vue';
import LicensePlate from '../partials/LicensePlate.vue';
import { AccountEyeIcon, ClipboardPlusIcon } from 'mdi-vue3';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import type { UserShort } from '~~/gen/ts/resources/users/users';

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

const query = ref<{ plate: string; model?: string; user_id?: number }>({
    plate: '',
    user_id: props.userId,
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * page.value : 0));

const { data, pending, refresh } = useLazyAsyncData(`vehicles-${page.value}`, () => listVehicles());

const hideModell = ref(false);

async function listVehicles(): Promise<ListVehiclesResponse> {
    try {
        const call = $grpc.getDMVClient().listVehicles({
            pagination: {
                offset: offset.value,
            },
            orderBy: [],
            userId: query.value.user_id,
            search: query.value.plate,
            model: query.value.model,
        });
        const { response } = await call;

        if (response.vehicles.length > 0) {
            if (response.vehicles[0].model === undefined) {
                hideModell.value = true;
            } else {
                hideModell.value = false;
            }
        }

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const completorStore = useCompletorStore();

const usersLoading = ref(false);
const selectedUser = ref<undefined | UserShort>();
watch(selectedUser, () => {
    if (selectedUser.value) {
        query.value.user_id = selectedUser.value.userId;
    } else {
        query.value.user_id = undefined;
    }
});

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), {
    debounce: 600,
    maxWait: 1400,
});

const clipboardStore = useClipboardStore();

const notifications = useNotificatorStore();

function addToClipboard(vehicle: Vehicle): void {
    clipboardStore.addVehicle(vehicle);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.vehicle_added.title', parameters: {} },
        content: { key: 'notifications.clipboard.vehicle_added.content', parameters: {} },
        duration: 3250,
        type: 'info',
    });
}

const columns = [
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
].filter((c) => c !== undefined);
</script>

<template>
    <div>
        <UDashboardToolbar>
            <template #default>
                <form class="w-full" @submit.prevent="refresh()">
                    <div class="flex flex-row gap-2">
                        <div class="flex-1">
                            <label for="search" class="block text-sm font-medium leading-6 text-neutral">
                                {{ $t('common.license_plate') }}
                            </label>
                            <div class="relative mt-2 flex items-center">
                                <UInput
                                    v-model="query.plate"
                                    type="text"
                                    :placeholder="$t('common.license_plate')"
                                    block
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </div>
                        </div>
                        <div v-if="!hideModell" class="flex-1">
                            <label for="model" class="block text-sm font-medium leading-6 text-neutral">
                                {{ $t('common.model') }}
                            </label>
                            <div class="relative mt-2 flex items-center">
                                <UInput
                                    v-model="query.model"
                                    type="text"
                                    name="model"
                                    :placeholder="$t('common.model')"
                                    block
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                            </div>
                        </div>

                        <div v-if="!userId" class="flex-1">
                            <label for="owner" class="block text-sm font-medium leading-6 text-neutral">
                                {{ $t('common.owner') }}
                            </label>
                            <div class="relative mt-2 items-center">
                                <UInputMenu
                                    v-model="selectedUser"
                                    :search="
                                        async (query: string) => {
                                            usersLoading = true;
                                            return await completorStore
                                                .completeCitizens({
                                                    search: query,
                                                })
                                                .finally(() => (usersLoading = false));
                                        }
                                    "
                                    :loading="usersLoading"
                                    :debounce="200"
                                    selected-icon="i-mdi-check"
                                    :search-attributes="['firstname', 'lastname']"
                                    option-attribute="'firstname','lastname'"
                                    block
                                    :placeholder="
                                        selectedUser
                                            ? `${selectedUser?.firstname} ${selectedUser?.lastname} (${selectedUser?.dateofbirth})`
                                            : $t('common.owner')
                                    "
                                    trailing
                                    by="userId"
                                >
                                    <template #option="{ option: user }">
                                        {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                    </template>
                                </UInputMenu>
                            </div>
                        </div>
                    </div>
                </form>
            </template>
        </UDashboardToolbar>

        <UTable
            :loading="pending"
            :columns="columns"
            :rows="data?.vehicles"
            :empty-state="{ icon: 'i-mdi-car', label: $t('common.not_found', [$t('common.vehicle', 2)]) }"
            :page-count="(data?.pagination?.totalCount ?? 0) / (data?.pagination?.pageSize ?? 1)"
            :total="data?.pagination?.totalCount"
        >
            <template #plate-data="{ row }">
                <LicensePlate :plate="row.plate" class="mr-2" />
            </template>
            <template #type-data="{ row }">
                {{ toTitleCase(row.type) }}
            </template>
            <template v-if="!hideOwner" #owner-data="{ row }">
                <CitizenInfoPopover :user="row.owner" />
            </template>
            <template #actions-data="{ row }">
                <div class="flex flex-row justify-end">
                    <button
                        v-if="!hideCopy"
                        class="flex-initial text-primary-500 hover:text-primary-400"
                        @click="addToClipboard(row)"
                    >
                        <ClipboardPlusIcon class="ml-auto mr-2.5 h-auto w-5" aria-hidden="true" />
                    </button>
                    <NuxtLink
                        v-if="!hideCitizenLink && can('CitizenStoreService.ListCitizens')"
                        :to="{
                            name: 'citizens-id',
                            params: { id: row.owner?.userId ?? 0 },
                        }"
                        class="flex-initial text-primary-500 hover:text-primary-400"
                    >
                        <AccountEyeIcon class="ml-auto mr-2.5 h-auto w-5" aria-hidden="true" />
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
