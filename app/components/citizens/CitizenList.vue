<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { vMaska } from 'maska/vue';
import { h } from 'vue';
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
import CitizenLabelModal from './CitizenLabelModal.vue';

const { t } = useI18n();

const { attr, can } = useAuth();

const overlay = useOverlay();

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.object({
    name: z.string().max(64).optional(),
    phoneNumber: z.string().max(20).optional(),
    wanted: z.coerce.boolean().optional(),
    trafficInfractionPoints: z.coerce.number().nonnegative().optional(),
    openFines: z.coerce.number().nonnegative().optional(),
    dateofbirth: z.string().max(10).optional(),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'name',
                        desc: false,
                    },
                ]),
        })
        .default({ columns: [{ id: 'name', desc: false }] }),
    page: pageNumberSchema,
});

const query = useSearchForm('citizens', schema);

const { data, status, refresh, error } = useLazyAsyncData(
    `citizens-${JSON.stringify(query.sorting)}-${query.page}-${JSON.stringify(query)}`,
    () => listCitizens(),
);

async function listCitizens(): Promise<ListCitizensResponse> {
    try {
        const req: ListCitizensRequest = {
            pagination: {
                offset: calculateOffset(query.page, data.value?.pagination),
            },
            sort: query.sorting,
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
    clipboardStore.addUser(user);

    notifications.add({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        description: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const appConfig = useAppConfig();

const columns = computed(() =>
    (
        [
            {
                accessorKey: 'name',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.name'),
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
                accessorKey: 'jobLabel',
                header: t('common.job'),
                meta: {
                    class: {
                        td: 'hidden lg:table-cell',
                        th: 'hidden lg:table-cell',
                    },
                },
                cell: ({ row }) =>
                    `${row.original.jobLabel}${row.original.props?.jobName || row.original.props?.jobGradeNumber ? '*' : ''}`,
            },
            {
                accessorKey: 'sex',
                header: t('common.sex'),
                meta: {
                    class: {
                        td: 'hidden lg:table-cell',
                        th: 'hidden lg:table-cell',
                    },
                },
                cell: ({ row }) =>
                    h(
                        'span',
                        { class: sexToTextColor(row.original.sex ?? '') },
                        row.original.sex?.toUpperCase() ?? t('common.na'),
                    ),
            },
            attr('citizens.CitizensService/ListCitizens', 'Fields', 'PhoneNumber').value
                ? {
                      accessorKey: 'phoneNumber',
                      header: t('common.phone_number'),
                      cell: ({ row }) => h(PhoneNumberBlock, { number: row.original.phoneNumber, hideNaText: true }),
                  }
                : undefined,
            {
                accessorKey: 'dateofbirth',
                header: t('common.date_of_birth'),
                meta: {
                    class: {
                        td: 'hidden lg:table-cell',
                        th: 'hidden lg:table-cell',
                    },
                },
                cell: ({ row }) => row.original.dateofbirth,
            },
            attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints').value
                ? {
                      accessorKey: 'trafficInfractionPoints',
                      header: ({ column }) => {
                          const isSorted = column.getIsSorted();

                          return h(UButton, {
                              color: 'neutral',
                              variant: 'ghost',
                              label: t('common.traffic_infraction_points', 2),
                              icon: isSorted
                                  ? isSorted === 'asc'
                                      ? appConfig.custom.icons.sortAsc
                                      : appConfig.custom.icons.sortDesc
                                  : appConfig.custom.icons.sort,
                              class: '-mx-2.5',
                              onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                          });
                      },
                      sortable: true,
                      cell: ({ row }) => row.original.props?.trafficInfractionPoints,
                  }
                : undefined,
            attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.OpenFines').value
                ? {
                      accessorKey: 'openFines',
                      header: ({ column }) => {
                          const isSorted = column.getIsSorted();

                          return h(UButton, {
                              color: 'neutral',
                              variant: 'ghost',
                              label: t('common.fine', 2),
                              icon: isSorted
                                  ? isSorted === 'asc'
                                      ? appConfig.custom.icons.sortAsc
                                      : appConfig.custom.icons.sortDesc
                                  : appConfig.custom.icons.sort,
                              class: '-mx-2.5',
                              onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                          });
                      },
                      sortable: true,
                      cell: ({ row }) =>
                          row.original.props?.openFines !== undefined && row.original.props?.openFines > 0
                              ? $n(row.original.props?.openFines, 'currency')
                              : '',
                  }
                : undefined,
            {
                accessorKey: 'height',
                header: t('common.height'),
                meta: {
                    class: {
                        td: 'hidden lg:table-cell',
                        th: 'hidden lg:table-cell',
                    },
                },
                cell: ({ row }) => (row.original.height ? `${row.original.height}cm` : ''),
            },
            can('citizens.CitizensService/GetUser').value
                ? {
                      id: 'actions',
                      cell: ({ row }) =>
                          h('div', {}, [
                              h(UTooltip, { text: t('components.clipboard.clipboard_button.add') }, [
                                  h(UButton, {
                                      variant: 'link',
                                      icon: 'i-mdi-clipboard-plus',
                                      onClick: () => addToClipboard(row.original),
                                  }),
                              ]),
                              h(UTooltip, { text: t('common.show') }, [
                                  h(UButton, {
                                      variant: 'link',
                                      icon: 'i-mdi-eye',
                                      to: {
                                          name: 'citizens-id',
                                          params: { id: row.original.userId ?? 0 },
                                      },
                                  }),
                              ]),
                          ]),
                  }
                : undefined,
        ] as TableColumn<User>[]
    ).filter((c) => c !== undefined),
);

const citizenLabelModal = overlay.create(CitizenLabelModal);

const input = useTemplateRef('input');

defineShortcuts({
    '/': () => input.value?.inputRef?.focus(),
});
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.citizens.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <UButton
                        v-if="can('citizens.CitizensService/ManageLabels').value"
                        :label="$t('common.label', 2)"
                        icon="i-mdi-tag"
                        @click="citizenLabelModal.open({})"
                    />
                </template>
            </UDashboardNavbar>

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
                                class="w-full"
                                leading-icon="i-mdi-search"
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
                                class="w-full"
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
                        class="my-2"
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
                                        class="w-full"
                                    />
                                </UFormField>

                                <UFormField
                                    v-if="
                                        attr('citizens.CitizensService/ListCitizens', 'Fields', 'TrafficInfractionPoints').value
                                    "
                                    class="flex-1"
                                    name="trafficInfractionPoints"
                                    :label="$t('common.traffic_infraction_points', 2)"
                                >
                                    <UInputNumber
                                        v-model="query.trafficInfractionPoints"
                                        name="trafficInfractionPoints"
                                        :min="0"
                                        :placeholder="$t('common.traffic_infraction_points')"
                                        class="w-full"
                                    />
                                </UFormField>

                                <UFormField
                                    v-if="attr('citizens.CitizensService/ListCitizens', 'Fields', 'UserProps.OpenFines').value"
                                    class="flex-1"
                                    name="openFines"
                                    :label="$t('components.citizens.CitizenList.open_fine')"
                                >
                                    <UInputNumber
                                        v-model="query.openFines"
                                        name="openFines"
                                        :min="0"
                                        :step="1000"
                                        :placeholder="`${$t('common.fine')}`"
                                        class="w-full"
                                        :format-options="{
                                            style: 'currency',
                                            currency: 'USD',
                                            currencyDisplay: 'code',
                                            currencySign: 'accounting',
                                        }"
                                    />
                                </UFormField>
                            </div>
                        </template>
                    </UAccordion>
                </UForm>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.citizen', 2)])"
                :error="error"
                :retry="refresh"
            />

            <UTable
                v-else
                v-model:sorting="query.sorting.columns"
                class="flex-1"
                :loading="isRequestPending(status)"
                :columns="columns"
                :data="data?.users"
                :empty="$t('common.not_found', [$t('common.citizen', 2)])"
                :sorting-options="{ manualSorting: true }"
                :pagination-options="{ manualPagination: true }"
                sticky
            >
                <template #name-cell="{ row }">
                    <div class="inline-flex items-center gap-1 text-highlighted">
                        <ProfilePictureImg
                            :src="row.original.props?.mugshot?.filePath"
                            :name="`${row.original.firstname} ${row.original.lastname}`"
                            :alt="$t('common.mugshot')"
                            :enable-popup="true"
                            size="sm"
                        />

                        <span>{{ row.original.firstname }} {{ row.original.lastname }}</span>

                        <UBadge v-if="row.original.props?.wanted" color="error">
                            {{ $t('common.wanted').toUpperCase() }}
                        </UBadge>
                    </div>
                </template>
            </UTable>

            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" />
        </template>
    </UDashboardPanel>
</template>
