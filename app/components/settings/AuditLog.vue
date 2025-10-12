<script lang="ts" setup>
import { UBadge, UButton, UTooltip } from '#components';
import { CalendarDate } from '@internationalized/date';
import type { TableColumn } from '@nuxt/ui';
import { addDays } from 'date-fns';
import { h } from 'vue';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
import type { JSONDataType } from 'vue-json-pretty/types/utils';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import InputDateRangePopover from '~/components/partials/InputDateRangePopover.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import { type AuditEntry, EventAction, EventResult } from '~~/gen/ts/resources/audit/audit';
import type { SortByColumn } from '~~/gen/ts/resources/common/database/database';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ViewAuditLogRequest, ViewAuditLogResponse } from '~~/gen/ts/services/settings/settings';
import { grpcMethods, grpcServices } from '~~/gen/ts/svcs';
import CitizenInfoPopover from '../partials/citizens/CitizenInfoPopover.vue';
import SelectMenu from '../partials/SelectMenu.vue';
import { eventActionToBadgeColor, eventResultToBadgeColor } from './helpers';

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
    services: z.coerce.string().max(64).array().max(10).default([]),
    methods: z.coerce.string().max(64).array().max(10).default([]),
    actions: z.enum(EventAction).array().max(6).default([]),
    results: z.enum(EventResult).array().max(6).default([]),
    search: z.coerce.string().max(64).default(''),
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

const eventActions = Object.keys(EventAction)
    .map((e) => EventAction[e as keyof typeof EventAction])
    .filter((e) => {
        if (typeof e === 'string') {
            return false;
        } else if (typeof e === 'number' && e === 0) {
            return false;
        }
        return true;
    });
const actionOptions = eventActions.map((e) => ({
    label: t(`enums.settings.AuditLog.EventAction.${EventAction[e]}`),
    value: e,
}));

const eventResults = Object.keys(EventResult)
    .map((e) => EventResult[e as keyof typeof EventResult])
    .filter((e) => {
        if (typeof e === 'string') {
            return false;
        } else if (typeof e === 'number' && e === 0) {
            return false;
        }
        return true;
    });
const resultOptions = eventResults.map((e) => ({
    label: t(`enums.settings.AuditLog.EventResult.${EventResult[e]}`),
    value: e,
}));

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
        actions: query.actions,
        results: query.results,
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
    text += `Service/Method: \`${logEntry.service}/${logEntry.method}\`
Action: \`${EventAction[logEntry.action]}\`
Result: \`${EventResult[logEntry.result]}\`
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
                    h(UTooltip, { text: t('common.expand_collapse') }, () =>
                        h(UButton, {
                            color: 'neutral',
                            variant: 'ghost',
                            icon: 'i-mdi-chevron-down',
                            square: true,
                            'aria-label': 'Expand',
                            ui: {
                                leadingIcon: ['transition-transform', row.getIsExpanded() ? 'duration-200 rotate-180' : ''],
                            },
                            onClick: () => row.toggleExpanded(),
                        }),
                    ),
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
                        label: `${t('common.service')} / ${t('common.method')}`,
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                meta: {
                    class: {
                        td: 'text-default',
                    },
                },
                cell: ({ row }) => `${row.original.service}/${row.original.method}`,
            },
            {
                accessorKey: 'action',
                header: t('common.action'),
                cell: ({ row }) =>
                    h(UBadge, {
                        color: eventActionToBadgeColor(row.original.action),
                        label: t(`enums.settings.AuditLog.EventAction.${EventAction[row.original.action]}`),
                    }),
            },
            {
                accessorKey: 'result',
                header: t('common.result'),
                cell: ({ row }) =>
                    h(UBadge, {
                        color: eventResultToBadgeColor(row.original.result),
                        label: t(`enums.settings.AuditLog.EventResult.${EventResult[row.original.result]}`),
                    }),
            },
        ] as TableColumn<AuditEntry>[],
);

function actionsToLabel(actions: EventAction[]): string {
    return actions.map((c) => t(`enums.settings.AuditLog.EventAction.${EventAction[c ?? 0]}`)).join(', ');
}

function resultsToLabel(results: EventResult[]): string {
    return results.map((c) => t(`enums.settings.AuditLog.EventResult.${EventResult[c ?? 0]}`)).join(', ');
}

const dataToggled = ref(false);

const today = new Date();
const tomorrow = addDays(today, 1);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.audit_log')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/settings" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <template #default>
                    <UForm class="my-2 flex w-full flex-1 flex-col gap-2" :schema="schema" :state="query" @submit="refresh()">
                        <div class="flex flex-1 flex-row gap-2">
                            <UFormField class="flex-1" name="date" :label="$t('common.time_range')">
                                <InputDateRangePopover
                                    v-model="query.date"
                                    class="w-full"
                                    :max-value="
                                        new CalendarDate(tomorrow.getFullYear(), tomorrow.getMonth() + 1, tomorrow.getDate())
                                    "
                                    time
                                    clearable
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
                                    searchable-key="completor-citizens"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :filter-fields="['firstname', 'lastname']"
                                    block
                                    :placeholder="$t('common.user', 2)"
                                    trailing
                                    value-key="userId"
                                    class="w-full"
                                >
                                    <template #item-label="{ item }">
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
                                    class="w-full"
                                    :ui="{ trailing: 'pe-1' }"
                                >
                                    <template #trailing>
                                        <UButton
                                            v-if="query.search !== ''"
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

                        <UCollapsible>
                            <UButton
                                class="group"
                                color="neutral"
                                variant="ghost"
                                trailing-icon="i-mdi-chevron-down"
                                :label="$t('common.advanced_search')"
                                :ui="{
                                    trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                                }"
                                block
                            />

                            <template #content>
                                <div class="flex flex-row flex-wrap gap-1">
                                    <UFormField class="flex-1" name="service" :label="$t('common.service')">
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="query.services"
                                                multiple
                                                name="service"
                                                :placeholder="$t('common.service')"
                                                :items="grpcServices.map((s) => s.split('.').pop() ?? s)"
                                                class="w-full"
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
                                            :items="
                                                grpcMethods.filter((m) => query.services.some((s) => m.includes('.' + s + '/')))
                                            "
                                            class="w-full"
                                        >
                                            <template #item-label="{ item }">
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
                                                v-model="query.actions"
                                                multiple
                                                name="states"
                                                :placeholder="$t('common.state')"
                                                :items="actionOptions"
                                                label-key="label"
                                                value-key="value"
                                                class="w-full"
                                            >
                                                <template v-if="query.actions.length" #default>
                                                    {{ actionsToLabel(query.actions) }}
                                                </template>

                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.state')]) }}
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>
                                    </UFormField>

                                    <UFormField class="flex-1" name="results" :label="$t('common.result')">
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="query.results"
                                                multiple
                                                name="results"
                                                :placeholder="$t('common.result')"
                                                :items="resultOptions"
                                                label-key="label"
                                                value-key="value"
                                                class="w-full"
                                            >
                                                <template v-if="query.actions.length" #default>
                                                    {{ resultsToLabel(query.results) }}
                                                </template>

                                                <template #empty>
                                                    {{ $t('common.not_found', [$t('common.state')]) }}
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>
                                    </UFormField>
                                </div>
                            </template>
                        </UCollapsible>
                    </UForm>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.audit_log', 2)])"
                :error="error"
                :retry="refresh"
            />

            <UTable
                v-else
                v-model:sorting="query.sorting.columns"
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
                    <UCard :ui="{ header: 'p-2 sm:px-2', body: 'p-1 sm:p-1', footer: '' }">
                        <template #header>
                            <div class="flex flex-row items-center gap-2">
                                <div class="flex-1 font-semibold text-highlighted">
                                    {{ $t('common.data') }}
                                </div>

                                <div>
                                    <UTooltip :text="$t('common.expand_collapse')">
                                        <UButton
                                            icon="i-mdi-chevron-double-down"
                                            variant="link"
                                            size="sm"
                                            class="place-self-end"
                                            :class="dataToggled ? 'rotate-180' : ''"
                                            :ui="{
                                                leadingIcon: 'transition-transform duration-200',
                                            }"
                                            :data-state="dataToggled ? 'open' : 'closed'"
                                            @click="dataToggled = !dataToggled"
                                        />
                                    </UTooltip>
                                </div>
                            </div>
                        </template>

                        <span v-if="!row.original.data">{{ $t('common.na') }}</span>
                        <template v-else>
                            <VueJsonPretty
                                :data="JSON.parse(row.original.data!) as JSONDataType"
                                show-icon
                                show-length
                                virtual
                                :height="dataToggled ? 240 : 800"
                                show-line-number
                            />
                        </template>

                        <template v-if="row.original.meta" #footer>
                            <div class="flex flex-row items-center justify-between gap-2">
                                <div>
                                    <span class="font-semibold">{{ $t('common.duration') }}</span
                                    >: {{ row.original.meta.meta['duration_ms'] ?? $t('common.na') }} ms
                                </div>

                                <div>
                                    <span class="font-semibold">{{ $t('common.code') }}</span
                                    >: {{ row.original.meta.meta['code'] ?? $t('common.na') }}
                                </div>
                            </div>
                        </template>
                    </UCard>
                </template>
            </UTable>
        </template>

        <template #footer>
            <Pagination v-model="query.page" :pagination="data?.pagination" :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>
