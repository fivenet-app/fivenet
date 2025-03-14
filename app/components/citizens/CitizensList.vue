<script lang="ts" setup>
import type { TableColumn } from '#ui/types';
import { vMaska } from 'maska/vue';
import { z } from 'zod';
import { sexToTextColor } from '~/components/partials/citizens/helpers';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useClipboardStore } from '~/stores/clipboard';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/users';
import type { ListCitizensRequest, ListCitizensResponse } from '~~/gen/ts/services/citizenstore/citizenstore';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const schema = z.object({
    name: z.string().max(64).optional(),
    phoneNumber: z.string().max(20).optional(),
    wanted: z.boolean().optional(),
    trafficInfractionPoints: z.coerce.number().nonnegative().optional(),
    openFines: z.coerce.number().nonnegative().optional(),
    dateofbirth: z.string().max(10).optional(),
});

type Schema = z.output<typeof schema>;

const { attr, can } = useAuth();

const query = reactive<Schema>({});

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'name',
    direction: 'asc',
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `citizens-${sort.value.column}:${sort.value.direction}-${page.value}-${JSON.stringify(query)}`,
    () => listCitizens(),
    {
        transform: (input) => ({ ...input, users: wrapRows(input?.users, columns) }),
        watch: [sort],
    },
);

async function listCitizens(): Promise<ListCitizensResponse> {
    try {
        const req: ListCitizensRequest = {
            pagination: {
                offset: offset.value,
            },
            sort: sort.value,
            search: query.name ?? '',
        };
        if (query.wanted) {
            req.wanted = query.wanted;
        }
        if (query.phoneNumber) {
            req.phoneNumber = query.phoneNumber;
        }
        if (query.trafficInfractionPoints) {
            req.trafficInfractionPoints = query.trafficInfractionPoints ?? 0;
        }
        if (query.openFines) {
            req.openFines = query.openFines ?? 0;
        }
        if (query.dateofbirth) {
            req.dateofbirth = query.dateofbirth;
        }

        const call = $grpc.citizenstore.citizenStore.listCitizens(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

function addToClipboard(user: User): void {
    clipboardStore.addUser({
        ...user,
        // @ts-expect-error wrapped table rows
        jobLabel: user.jobLabel.value,
        // @ts-expect-error wrapped table rows
        sex: user.sex.value,
        // @ts-expect-error wrapped table rows
        dateofbirth: user.dateofbirth.value,
        // @ts-expect-error wrapped table rows
        height: user.height.value,
    });

    notifications.add({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        description: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        timeout: 3250,
        type: NotificationType.INFO,
    });
}

const columns = [
    {
        key: 'name',
        label: t('common.name'),
        sortable: true,
    },
    {
        key: 'jobLabel',
        label: t('common.job'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    {
        key: 'sex',
        label: t('common.sex'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber').value
        ? { key: 'phoneNumber', label: t('common.phone_number') }
        : undefined,
    {
        key: 'dateofbirth',
        label: t('common.date_of_birth'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints').value
        ? {
              key: 'trafficInfractionPoints',
              label: t('common.traffic_infraction_points', 2),
              sortable: true,
          }
        : undefined,
    attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines').value
        ? {
              key: 'openFines',
              label: t('common.fine', 2),
              sortable: true,
          }
        : undefined,
    {
        key: 'height',
        label: t('common.height'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    can('CitizenStoreService.GetUser').value
        ? {
              key: 'actions',
              label: t('common.action', 2),
              sortable: false,
          }
        : undefined,
].filter((c) => c !== undefined) as TableColumn[];

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.input?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
            <div class="flex w-full flex-row gap-2">
                <UFormGroup class="flex-1" :label="$t('common.search')">
                    <UInput
                        ref="input"
                        v-model="query.name"
                        type="text"
                        name="name"
                        :placeholder="`${$t('common.citizen', 1)} ${$t('common.name')}`"
                        block
                        leading-icon="i-mdi-search"
                        @keydown.esc="$event.target.blur()"
                    >
                        <template #trailing>
                            <UKbd value="/" />
                        </template>
                    </UInput>
                </UFormGroup>

                <UFormGroup name="dateofbirth" :label="$t('common.date_of_birth')">
                    <UInput
                        v-model="query.dateofbirth"
                        v-maska
                        type="text"
                        name="dateofbirth"
                        data-maska="##[./]##[./]####"
                        :placeholder="`${$t('common.date_of_birth')} (DD.MM.YYYY)`"
                        block
                    />
                </UFormGroup>

                <UFormGroup
                    v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.Wanted').value"
                    name="wanted"
                    :label="$t('components.citizens.CitizensList.only_wanted')"
                    class="flex flex-initial flex-col"
                    :ui="{ container: 'flex-1 flex' }"
                >
                    <div class="flex flex-1 items-center">
                        <UToggle v-model="query.wanted">
                            <span class="sr-only">
                                {{ $t('components.citizens.CitizensList.only_wanted') }}
                            </span>
                        </UToggle>
                    </div>
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
                            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber').value"
                            name="phoneNumber"
                            :label="$t('common.phone_number')"
                            class="flex-1"
                        >
                            <UInput
                                v-model="query.phoneNumber"
                                type="tel"
                                name="phoneNumber"
                                :placeholder="$t('common.phone_number')"
                                block
                            />
                        </UFormGroup>

                        <UFormGroup
                            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'TrafficInfractionPoints').value"
                            name="trafficInfractionPoints"
                            :label="$t('common.traffic_infraction_points', 2)"
                            class="flex-1"
                        >
                            <UInput
                                v-model="query.trafficInfractionPoints"
                                type="number"
                                name="trafficInfractionPoints"
                                :smin="0"
                                :placeholder="$t('common.traffic_infraction_points')"
                                block
                            />
                        </UFormGroup>

                        <UFormGroup
                            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines').value"
                            name="openFines"
                            :label="$t('components.citizens.CitizensList.open_fine')"
                            class="flex-1"
                        >
                            <UInput
                                v-model="query.openFines"
                                type="number"
                                name="openFines"
                                :min="0"
                                :placeholder="`${$t('common.fine')}`"
                                block
                                leading-icon="i-mdi-dollar"
                            />
                        </UFormGroup>
                    </div>
                </template>
            </UAccordion>
        </UForm>
    </UDashboardToolbar>

    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
        :error="error"
        :retry="refresh"
    />
    <UTable
        v-else
        v-model:sort="sort"
        :loading="loading"
        :columns="columns"
        :rows="data?.users"
        :empty-state="{ icon: 'i-mdi-accounts', label: $t('common.not_found', [$t('common.citizen', 2)]) }"
        sort-mode="manual"
        class="flex-1"
    >
        <template #name-data="{ row: citizen }">
            <div class="inline-flex items-center gap-1 text-gray-900 dark:text-white">
                <ProfilePictureImg
                    :src="citizen.props?.mugShot?.url"
                    :name="`${citizen.firstname} ${citizen.lastname}`"
                    :alt="$t('common.mug_shot')"
                    :enable-popup="true"
                    size="sm"
                />

                <span>{{ citizen.firstname }} {{ citizen.lastname }}</span>
                <span class="lg:hidden"> ({{ citizen.dateofbirth.value }}) </span>

                <UBadge v-if="citizen.props?.wanted" color="red">
                    {{ $t('common.wanted').toUpperCase() }}
                </UBadge>
            </div>

            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.sex') }} - {{ $t('common.job') }}</dt>
                <dd class="mt-1 truncate">
                    {{ citizen.sex?.value.toUpperCase() ?? $t('common.na') }} -
                    {{ citizen.jobLabel.value ?? $t('common.na') }}
                </dd>
            </dl>
        </template>

        <template #jobLabel-data="{ row: citizen }">
            {{ citizen.jobLabel.value }}
            {{ citizen.props?.jobName || citizen.props?.jobGradeNumber ? '*' : '' }}
        </template>

        <template #sex-data="{ row: citizen }">
            <span :class="sexToTextColor(citizen.sex?.value ?? '')">
                {{ citizen.sex?.value.toUpperCase() ?? $t('common.na') }}
            </span>
        </template>

        <template #phoneNumber-data="{ row: citizen }">
            <PhoneNumberBlock :number="citizen.phoneNumber" hide-na-text />
        </template>

        <template #openFines-data="{ row: citizen }">
            <template v-if="(citizen.props?.openFines ?? 0) > 0">
                {{ $n(citizen.props?.openFines ?? 0, 'currency') }}
            </template>
        </template>

        <template #dateofbirth-data="{ row: citizen }">
            {{ citizen.dateofbirth.value }}
        </template>

        <template #height-data="{ row: citizen }"> {{ citizen.height.value ? citizen.height.value + 'cm' : '' }} </template>

        <template v-if="can('CitizenStoreService.GetUser').value" #actions-data="{ row: citizen }">
            <div :key="citizen.userId" class="flex flex-col justify-end md:flex-row">
                <UTooltip :text="$t('components.clipboard.clipboard_button.add')">
                    <UButton variant="link" icon="i-mdi-clipboard-plus" @click="addToClipboard(citizen)" />
                </UTooltip>

                <UTooltip :text="$t('common.show')">
                    <UButton
                        variant="link"
                        icon="i-mdi-eye"
                        :to="{
                            name: 'citizens-id',
                            params: { id: citizen.userId ?? 0 },
                        }"
                    />
                </UTooltip>
            </div>
        </template>
    </UTable>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
