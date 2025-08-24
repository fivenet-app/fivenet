<script lang="ts" setup>
import type { TableColumn } from '@nuxt/ui';
import { vMaska } from 'maska/vue';
import { z } from 'zod';
import { sexToTextColor } from '~/components/partials/citizens/helpers';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useClipboardStore } from '~/stores/clipboard';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/users';
import type { ListCitizensRequest, ListCitizensResponse } from '~~/gen/ts/services/citizens/citizens';

const { t } = useI18n();

const { attr, can } = useAuth();

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.object({
    name: z.string().max(64).optional(),
    phoneNumber: z.string().max(20).optional(),
    wanted: z.coerce.boolean().optional(),
    trafficInfractionPoints: z.coerce.number().nonnegative().optional(),
    openFines: z.coerce.number().nonnegative().optional(),
    dateofbirth: z.string().max(10).optional(),
    sorting: z
        .custom<SortByColumn>()
        .array()
        .max(3)
        .default([
            {
                id: 'name',
                desc: false,
            },
        ]),
    page: pageNumberSchema,
});

const query = useSearchForm('citizens', schema);

const { data, status, refresh, error } = useLazyAsyncData(
    `citizens-${query.sorting.column}:${query.sorting.direction}-${query.page}-${JSON.stringify(query)}`,
    () => listCitizens(),
);

async function listCitizens(): Promise<ListCitizensResponse> {
    try {
        const req: ListCitizensRequest = {
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: { columns: query.sorting },
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

        const call = citizensCitizensClient.listCitizens(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

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
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const columns = [
    {
        accessorKey: 'name',
        label: t('common.name'),
        sortable: true,
    },
    {
        accessorKey: 'jobLabel',
        label: t('common.job'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    {
        accessorKey: 'sex',
        label: t('common.sex'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    attr('citizens.CitizensService/ListCitizens', 'Fields', 'PhoneNumber').value
        ? { accessorKey: 'phoneNumber', label: t('common.phone_number') }
        : undefined,
    {
        accessorKey: 'dateofbirth',
        label: t('common.date_of_birth'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints').value
        ? {
              accessorKey: 'trafficInfractionPoints',
              label: t('common.traffic_infraction_points', 2),
              sortable: true,
          }
        : undefined,
    attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.OpenFines').value
        ? {
              accessorKey: 'openFines',
              label: t('common.fine', 2),
              sortable: true,
          }
        : undefined,
    {
        accessorKey: 'height',
        label: t('common.height'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    can('citizens.CitizensService/GetUser').value
        ? {
              accessorKey: 'actions',
              label: t('common.action', 2),
              sortable: false,
          }
        : undefined,
].filter((c) => c !== undefined) as TableColumn[];

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.inputRef?.focus(),
});
</script>

<template>
    <UDashboardToolbar>
        <UForm class="w-full" :schema="schema" :state="query" @submit="refresh()">
            <div class="flex w-full flex-row gap-2">
                <UFormField class="flex-1" :label="$t('common.search')" name="name">
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
                </UFormField>

                <UFormField name="dateofbirth" :label="$t('common.date_of_birth')">
                    <UInput
                        v-model="query.dateofbirth"
                        v-maska
                        type="text"
                        name="dateofbirth"
                        :placeholder="`${$t('common.date_of_birth')} (DD.MM.YYYY)`"
                        block
                        data-maska="##.##.####"
                    />
                </UFormField>

                <UFormField
                    v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.Wanted').value"
                    class="flex flex-initial flex-col"
                    name="wanted"
                    :label="$t('common.only_wanted')"
                    :ui="{ container: 'flex-1 flex' }"
                >
                    <div class="flex flex-1 items-center">
                        <USwitch v-model="query.wanted" />
                    </div>
                </UFormField>
            </div>

            <UAccordion
                class="mt-2"
                color="neutral"
                variant="soft"
                size="sm"
                :items="[{ label: $t('common.advanced_search'), slot: 'search' as const }]"
            >
                <template #search>
                    <div class="flex flex-row gap-2">
                        <UFormField
                            v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'PhoneNumber').value"
                            class="flex-1"
                            name="phoneNumber"
                            :label="$t('common.phone_number')"
                        >
                            <UInput
                                v-model="query.phoneNumber"
                                type="tel"
                                name="phoneNumber"
                                :placeholder="$t('common.phone_number')"
                                block
                            />
                        </UFormField>

                        <UFormField
                            v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'TrafficInfractionPoints').value"
                            class="flex-1"
                            name="trafficInfractionPoints"
                            :label="$t('common.traffic_infraction_points', 2)"
                        >
                            <UInput
                                v-model="query.trafficInfractionPoints"
                                type="number"
                                name="trafficInfractionPoints"
                                :min="0"
                                :placeholder="$t('common.traffic_infraction_points')"
                                block
                            />
                        </UFormField>

                        <UFormField
                            v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.OpenFines').value"
                            class="flex-1"
                            name="openFines"
                            :label="$t('components.citizens.CitizenList.open_fine')"
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
                        </UFormField>
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
        v-model:sorting="query.sorting"
        class="flex-1"
        :loading="isRequestPending(status)"
        :columns="columns"
        :data="data?.users"
        :empty-state="{ icon: 'i-mdi-accounts', label: $t('common.not_found', [$t('common.citizen', 2)]) }"
        sort-mode="manual"
    >
        <template #name-cell="{ row: citizen }">
            <div class="inline-flex items-center gap-1 text-highlighted">
                <ProfilePictureImg
                    :src="citizen.original.props?.mugshot?.filePath"
                    :name="`${citizen.original.firstname} ${citizen.original.lastname}`"
                    :alt="$t('common.mugshot')"
                    :enable-popup="true"
                    size="sm"
                />

                <span>{{ citizen.original.firstname }} {{ citizen.original.lastname }}</span>
                <span class="lg:hidden"> ({{ citizen.original.dateofbirth.value }}) </span>

                <UBadge v-if="citizen.original.props?.wanted" color="error">
                    {{ $t('common.wanted').toUpperCase() }}
                </UBadge>
            </div>

            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.sex') }} - {{ $t('common.job') }}</dt>
                <dd class="mt-1 truncate">
                    {{ citizen.original.sex?.value.toUpperCase() ?? $t('common.na') }} -
                    {{ citizen.original.jobLabel.value ?? $t('common.na') }}
                </dd>
            </dl>
        </template>

        <template #jobLabel-cell="{ row: citizen }">
            {{ citizen.original.jobLabel.value }}
            {{ citizen.original.props?.jobName || citizen.original.props?.jobGradeNumber ? '*' : '' }}
        </template>

        <template #sex-cell="{ row: citizen }">
            <span :class="sexToTextColor(citizen.original.sex?.value ?? '')">
                {{ citizen.original.sex?.value.toUpperCase() ?? $t('common.na') }}
            </span>
        </template>

        <template #phoneNumber-cell="{ row: citizen }">
            <PhoneNumberBlock :number="citizen.original.phoneNumber" hide-na-text />
        </template>

        <template #openFines-cell="{ row: citizen }">
            <template v-if="(citizen.original.props?.openFines ?? 0) > 0">
                {{ $n(citizen.original.props?.openFines ?? 0, 'currency') }}
            </template>
        </template>

        <template #dateofbirth-cell="{ row: citizen }">
            {{ citizen.original.dateofbirth.value }}
        </template>

        <template #height-cell="{ row: citizen }">
            {{ citizen.original.height.value ? citizen.original.height.value + 'cm' : '' }}
        </template>

        <template v-if="can('citizens.CitizensService/GetUser').value" #actions-cell="{ row: citizen }">
            <div :key="citizen.original.userId" class="flex flex-col justify-end md:flex-row">
                <UTooltip :text="$t('components.clipboard.clipboard_button.add')">
                    <UButton variant="link" icon="i-mdi-clipboard-plus" @click="addToClipboard(citizen.original)" />
                </UTooltip>

                <UTooltip :text="$t('common.show')">
                    <UButton
                        variant="link"
                        icon="i-mdi-eye"
                        :to="{
                            name: 'citizens-id',
                            params: { id: citizen.original.userId ?? 0 },
                        }"
                    />
                </UTooltip>
            </div>
        </template>
    </UTable>

    <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" />
</template>
