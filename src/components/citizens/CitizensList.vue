<script lang="ts" setup>
import { vMaska } from 'maska';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { attr } from '~/composables/can';
import GenericInput from '~/composables/partials/forms/GenericInput.vue';
import { ListCitizensRequest, ListCitizensResponse } from '~~/gen/ts/services/citizenstore/citizenstore';
import type { User } from '~~/gen/ts/resources/users/users';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import PhoneNumberBlock from '../partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '../partials/citizens/ProfilePictureImg.vue';

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
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

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
            req.openFines = BigInt(query.value.fines?.toString());
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
                <UForm class="w-full" :state="{}" @submit="refresh()">
                    <div class="flex w-full flex-row gap-2">
                        <UFormGroup class="flex-1" :label="`${$t('common.citizen', 1)} ${$t('common.name')}`">
                            <GenericInput
                                v-model="query.name"
                                type="text"
                                name="searchName"
                                block
                                :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                            />
                        </UFormGroup>

                        <UFormGroup :label="$t('common.date_of_birth')">
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
                        </UFormGroup>

                        <UFormGroup
                            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.Wanted')"
                            :label="$t('components.citizens.citizens_list.only_wanted')"
                        >
                            <UToggle v-model="query.wanted">
                                <span class="sr-only">
                                    {{ $t('components.citizens.citizens_list.only_wanted') }}
                                </span>
                            </UToggle>
                        </UFormGroup>
                    </div>

                    <UAccordion
                        class="mt-2"
                        color="white"
                        variant="soft"
                        size="sm"
                        :items="[{ label: $t('common.advanced_search'), slot: 'search' }]"
                    >
                        <template #search>
                            <div class="flex flex-row gap-2">
                                <UFormGroup
                                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')"
                                    class="flex-1"
                                    :label="$t('common.phone_number')"
                                >
                                    <UInput
                                        v-model="query.phoneNumber"
                                        type="tel"
                                        name="searchPhone"
                                        :placeholder="$t('common.phone_number')"
                                        block
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'TrafficInfractionPoints')"
                                    class="flex-1"
                                    :label="$t('common.traffic_infraction_points', 2)"
                                >
                                    <UInput
                                        v-model="query.trafficInfractionPoints"
                                        type="number"
                                        name="trafficInfractionPoints"
                                        :placeholder="$t('common.traffic_infraction_points')"
                                        block
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </UFormGroup>

                                <UFormGroup
                                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines')"
                                    class="flex-1"
                                    :label="$t('components.citizens.citizens_list.open_fine')"
                                >
                                    <UInput
                                        v-model="query.fines"
                                        type="number"
                                        name="fine"
                                        :placeholder="`${$t('common.fine')}`"
                                        block
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </UFormGroup>
                            </div>
                        </template>
                    </UAccordion>
                </UForm>
            </template>
        </UDashboardToolbar>

        <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.citizen', 2)])" :retry="refresh" />
        <UTable
            :loading="pending"
            :columns="columns"
            :rows="data?.users"
            :empty-state="{ icon: 'i-mdi-accounts', label: $t('common.not_found', [$t('common.citizen', 2)]) }"
        >
            <template #name-data="{ row: citizen }">
                <div class="inline-flex items-center">
                    <ProfilePictureImg
                        :url="citizen.props?.mugShot?.url"
                        :name="`${citizen.firstname} ${citizen.lastname}`"
                        :alt="$t('common.mug_shot')"
                        :enable-popup="true"
                        size="sm"
                        class="mr-2"
                    />

                    <span>{{ citizen.firstname }} {{ citizen.lastname }}</span>
                    <span class="lg:hidden"> ({{ citizen.dateofbirth }}) </span>
                </div>

                <span
                    v-if="citizen.props?.wanted"
                    class="ml-1 rounded-md bg-error-100 px-2 py-0.5 text-sm font-medium text-error-700"
                >
                    {{ $t('common.wanted').toUpperCase() }}
                </span>

                <dl class="font-normal lg:hidden">
                    <dt class="sr-only">{{ $t('common.sex') }} - {{ $t('common.job') }}</dt>
                    <dd class="mt-1 truncate">{{ citizen.sex!.toUpperCase() }} - {{ citizen.jobLabel }}</dd>
                </dl>
            </template>
            <template #jobLabel-data="{ row: citizen }">
                {{ citizen.jobLabel }}
            </template>
            <template #sex-data="{ row: citizen }">
                {{ citizen.sex!.toUpperCase() }}
            </template>
            <template #phoneNumber-data="{ row: citizen }">
                <PhoneNumberBlock :number="citizen.phoneNumber" />
            </template>
            <template #openFines-data="{ row: citizen }">
                <template v-if="(citizen.props?.openFines ?? 0n) > 0n">
                    {{ $n(parseInt((citizen?.props?.openFines ?? 0n).toString()), 'currency') }}
                </template>
            </template>
            <template #height-data="{ row: citizen }"> {{ citizen.height }}cm </template>
            <template #actions-data="{ row: citizen }">
                <div v-if="can('CitizenStoreService.GetUser')" class="flex flex-col justify-end gap-1 md:flex-row">
                    <UButton variant="link" icon="i-mdi-clipboard-plus" @click="addToClipboard(citizen)" />

                    <UButton
                        variant="link"
                        icon="i-mdi-eye"
                        :to="{
                            name: 'citizens-id',
                            params: { id: citizen.userId ?? 0 },
                        }"
                    />
                </div>
            </template>
        </UTable>

        <div class="flex justify-end border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
            <UPagination
                v-model="page"
                :page-count="data?.pagination?.pageSize ?? 0"
                :total="data?.pagination?.totalCount ?? 0"
            />
        </div>
    </div>
</template>
