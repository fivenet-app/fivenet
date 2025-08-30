<script lang="ts" setup>
import { UBadge, UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { addDays } from 'date-fns';
import { h } from 'vue';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
import type { JSONDataType } from 'vue-json-pretty/types/utils';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DateRangePickerClient from '~/components/partials/DateRangePicker.client.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import { type AuditEntry, EventType } from '~~/gen/ts/resources/audit/audit';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ViewAuditLogRequest, ViewAuditLogResponse } from '~~/gen/ts/services/settings/settings';
import { grpcMethods, grpcServices } from '~~/gen/ts/svcs';
import CitizenInfoPopover from '../partials/citizens/CitizenInfoPopover.vue';
import SelectMenu from '../partials/SelectMenu.vue';
import { eventTypeToBadgeColor } from './helpers';

const { d, t } = useI18n();

const completorStore = useCompletorStore();

const settingsSettingsClient = await getSettingsSettingsClient();

const schema = z.object({
    users: z.coerce.number().array().max(5).default([]),
    date: z
        .object({
            start: z.coerce.date(),
            end: z.coerce.date(),
        })
        .optional(),
    services: z.string().max(64).array().max(10).default([]),
    methods: z.string().max(64).array().max(10).default([]),
    states: z.nativeEnum(EventType).array().max(10).default([]),
    search: z.string().max(64).default(''),
    sorting: z
        .object({
            columns: z
                .custom<SortByColumn>()
                .array()
                .max(3)
                .default([
                    {
                        id: 'createdAt',
                        desc: true,
                    },
                ]),
        })
        .default({ columns: [{ id: 'createdAt', desc: true }] }),
    page: pageNumberSchema,
});

const query = useSearchForm('settings_auditlog', schema);

const eventTypes = Object.keys(EventType)
    .map((eventType) => EventType[eventType as keyof typeof EventType])
    .filter((eventType) => {
        if (typeof eventType === 'string') {
            return false;
        } else if (typeof eventType === 'number' && eventType === 0) {
            return false;
        }
        return true;
    });
const statesOptions = eventTypes.map((eventType) => ({ eventType: eventType }));

const { data, status, refresh, error } = useLazyAsyncData(
    () =>
        `settings-audit-${JSON.stringify(query.sorting)}-${query.page}-${query.date?.start}-${query.date?.end}-${query.methods}-${query.services}-${query.search}-${query.users.join(',')}`,
    () => viewAuditLog(),
);

async function viewAuditLog(): Promise<ViewAuditLogResponse> {
    const req: ViewAuditLogRequest = {
        pagination: {
            offset: calculateOffset(query.page, data.value?.pagination),
        },
        sort: query.sorting,
        userIds: query.users,
        services: query.services,
        // Make sure to remove the service from the beginning
        methods: query.methods.map((m) => m.split('/').pop() ?? m),
        states: query.states,
    };

    if (query.date) {
        req.from = toTimestamp(query.date.start);
        req.to = toTimestamp(query.date.end);
    }

    if (query.search !== '') {
        req.search = query.search;
    }

    try {
        const call = settingsSettingsClient.viewAuditLog(req);
        const { response } = await call;

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

const notifications = useNotificationsStore();

async function addToClipboard(logEntry: AuditEntry): Promise<void> {
    const user = logEntry.user;
    let text = `**Audit Log Entry ${logEntry.id} - ${d(toDate(logEntry.createdAt)!, 'short')}**

`;
    if (user) {
        text += `User: ${user?.firstname} ${user?.lastname} (${user?.userId}; ${user?.identifier})
`;
    }
    text += `Action: \`${logEntry.service}/${logEntry.method}\`
Event: \`${EventType[logEntry.state]}\`
`;
    if (logEntry.data) {
        text += `Data:
\`\`\`json
${JSON.stringify(JSON.parse(logEntry.data!), undefined, 2)}
\`\`\`
`;
    } else {
        text += `Data: ${t('common.na')}
`;
    }

    notifications.add({
        title: { key: 'notifications.settings.audit_log.title', parameters: {} },
        description: { key: 'notifications.settings.audit_log.content', parameters: {} },
        type: NotificationType.INFO,
    });

    return copyToClipboardWrapper(text);
}

const appConfig = useAppConfig();

const columns = computed(
    () =>
        [
            {
                id: 'expand',
                cell: ({ row }) =>
                    h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        icon: 'i-lucide-chevron-down',
                        square: true,
                        'aria-label': 'Expand',
                        ui: {
                            leadingIcon: ['transition-transform', row.getIsExpanded() ? 'duration-200 rotate-180' : ''],
                        },
                        onClick: () => row.toggleExpanded(),
                    }),
            },
            {
                id: 'actions',
                cell: ({ row }) =>
                    h(UTooltip, { text: t('components.clipboard.clipboard_button.add') }, () =>
                        h(UButton, {
                            variant: 'link',
                            icon: 'i-mdi-content-copy',
                            onClick: () => addToClipboard(row.original),
                        }),
                    ),
            },
            {
                accessorKey: 'id',
                header: t('common.id'),
                cell: ({ row }) => row.original.id,
            },
            {
                accessorKey: 'createdAt',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.created_at'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt, type: 'long' }),
            },
            {
                accessorKey: 'user',
                header: t('common.user'),
                cell: ({ row }) => h(CitizenInfoPopover, { user: row.original.user }),
            },
            {
                accessorKey: 'service',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.service'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                cell: ({ row }) => row.original.service,
            },
            {
                accessorKey: 'state',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.state'),
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
                    h(UBadge, {
                        color: eventTypeToBadgeColor(row.original.state),
                        label: t(`enums.settings.AuditLog.EventType.${EventType[row.original.state]}`),
                    }),
            },
        ] as TableColumn<AuditEntry>[],
);

function statesToLabel(states: { eventType: EventType }[]): string {
    return states.map((c) => t(`enums.settings.AuditLog.EventType.${EventType[c.eventType ?? 0]}`)).join(', ');
}
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm class="w-full" :schema="schema" :state="query" @submit="refresh()">
                <div class="flex flex-row flex-wrap gap-2">
                    <UFormField class="flex-1" name="date" :label="$t('common.time_range')">
                        <DateRangePickerClient
                            v-model="query.date"
                            class="flex-1"
                            mode="date"
                            :popover="{ class: 'flex-1' }"
                            :date-picker="{
                                mode: 'dateTime',
                                disabledDates: [{ start: addDays(new Date(), 1), end: null }],
                                is24Hr: true,
                                clearable: true,
                            }"
                        />
                    </UFormField>

                    <UFormField class="flex-1" name="user" :label="$t('common.user')">
                        <SelectMenu
                            v-model="query.users"
                            multiple
                            :searchable="
                                async (q: string) =>
                                    await completorStore.completeCitizens({
                                        search: q,
                                        userIds: query.users,
                                    })
                            "
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['firstname', 'lastname']"
                            block
                            :placeholder="$t('common.user', 2)"
                            trailing
                            value-key="userId"
                        >
                            <template #item-label="{ item }">
                                {{ userToLabel(item) }}
                            </template>

                            <template #item="{ item }">
                                {{ userToLabel(item) }}
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                        </SelectMenu>
                    </UFormField>

                    <UFormField class="flex-1" name="data" :label="$t('common.data')">
                        <UInput
                            v-model="query.search"
                            type="text"
                            name="data"
                            block
                            :placeholder="$t('common.search')"
                            leading-icon="i-mdi-search"
                            :ui="{ trailing: 'pe-1' }"
                        >
                            <template #trailing>
                                <UButton
                                    v-show="query.search !== ''"
                                    color="neutral"
                                    variant="link"
                                    icon="i-mdi-close"
                                    :aria-label="query.search ? 'Hide search' : 'Show search'"
                                    :aria-pressed="query.search"
                                    aria-controls="search"
                                    @click="query.search = ''"
                                />
                            </template>
                        </UInput>
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
                        <div class="flex flex-row flex-wrap gap-1">
                            <UFormField class="flex-1" name="service" :label="$t('common.service')">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.services"
                                        multiple
                                        name="service"
                                        :placeholder="$t('common.service')"
                                        :items="grpcServices.map((s) => s.split('.').pop() ?? s)"
                                    >
                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.service')]) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>

                            <UFormField class="flex-1" name="method" :label="$t('common.method')">
                                <USelectMenu
                                    v-model="query.methods"
                                    multiple
                                    name="method"
                                    :placeholder="$t('common.method')"
                                    :items="grpcMethods.filter((m) => query.services.some((s) => m.includes('.' + s + '/')))"
                                >
                                    <template #item="{ item }">
                                        {{ item.split('/').pop() }}
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.method')]) }}
                                    </template>
                                </USelectMenu>
                            </UFormField>

                            <UFormField class="flex-1" name="states" :label="$t('common.state')">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.states"
                                        multiple
                                        name="states"
                                        :placeholder="$t('common.state')"
                                        :items="statesOptions"
                                        value-key="eventType"
                                    >
                                        <template #item-label="{ item }">
                                            {{ statesToLabel([item]) }}
                                        </template>

                                        <template #item="{ item }">
                                            {{ $t(`enums.settings.AuditLog.EventType.${EventType[item.eventType]}`) }}
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.state')]) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>
                        </div>
                    </template>
                </UAccordion>
            </UForm>
        </template>
    </UDashboardToolbar>

    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [$t('common.audit_log', 2)])"
        :error="error"
        :retry="refresh"
    />

    <UTable
        v-else
        class="flex-1"
        :loading="isRequestPending(status)"
        :columns="columns"
        :data="data?.logs"
        :pagination-options="{ manualPagination: true }"
        :sorting-options="{ manualSorting: true }"
        :empty="$t('common.not_found', [$t('common.entry', 2)])"
        :ui="{ tr: 'data-[expanded=true]:bg-elevated/50' }"
        sticky
    >
        <template #expanded="{ row }">
            <div class="px-2 py-1">
                <span v-if="!row.original.data">{{ $t('common.na') }}</span>
                <template v-else>
                    <VueJsonPretty
                        :data="JSON.parse(row.original.data!) as JSONDataType"
                        :show-icon="true"
                        :show-length="true"
                        :virtual="true"
                        :height="240"
                    />
                </template>
            </div>
        </template>
    </UTable>

    <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
</template>
