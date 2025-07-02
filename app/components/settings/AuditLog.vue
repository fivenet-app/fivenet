<script lang="ts" setup>
import { addDays } from 'date-fns';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
import type { JSONDataType } from 'vue-json-pretty/types/utils';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DateRangePickerPopoverClient from '~/components/partials/DateRangePickerPopover.client.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import type { AuditEntry } from '~~/gen/ts/resources/audit/audit';
import { EventType } from '~~/gen/ts/resources/audit/audit';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ViewAuditLogRequest, ViewAuditLogResponse } from '~~/gen/ts/services/settings/settings';
import { grpcMethods, grpcServices } from '~~/gen/ts/svcs';
import { eventTypeToBadgeColor } from './helpers';

const { $grpc } = useNuxtApp();

const { d, t } = useI18n();

const completorStore = useCompletorStore();

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
    sort: z.custom<TableSortable>().default({
        column: 'createdAt',
        direction: 'desc',
    }),
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

const usersLoading = ref(false);

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `settings-audit-${query.sort.column}:${query.sort.direction}-${query.page}-${query.date?.start}-${query.date?.end}-${query.methods}-${query.services}-${query.search}-${query.users.join(',')}`,
    () => viewAuditLog(),
);

async function viewAuditLog(): Promise<ViewAuditLogResponse> {
    const req: ViewAuditLogRequest = {
        pagination: {
            offset: calculateOffset(query.page, data.value?.pagination),
        },
        sort: query.sort,
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
        const call = $grpc.settings.settings.viewAuditLog(req);
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

const columns = [
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
    {
        key: 'id',
        label: t('common.id'),
    },
    {
        key: 'createdAt',
        label: t('common.created_at'),
        sortable: true,
    },
    {
        key: 'user',
        label: t('common.user', 1),
    },
    {
        key: 'service',
        label: `${t('common.service')}/${t('common.method')}`,
        sortable: true,
    },
    {
        key: 'state',
        label: t('common.state'),
        sortable: true,
    },
];

const expand = ref({
    openedRows: [],
    row: {},
});

function statesToLabel(states: { eventType: EventType }[]): string {
    return states.map((c) => t(`enums.settings.AuditLog.EventType.${EventType[c.eventType ?? 0]}`)).join(', ');
}
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm class="w-full" :schema="schema" :state="query" @submit="refresh()">
                <div class="flex flex-row flex-wrap gap-2">
                    <UFormGroup class="flex-1" name="date" :label="$t('common.time_range')">
                        <DateRangePickerPopoverClient
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
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="user" :label="$t('common.user')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="query.users"
                                multiple
                                :searchable="
                                    async (q: string) => {
                                        usersLoading = true;
                                        const users = await completorStore.completeCitizens({
                                            search: q,
                                            userIds: query.users,
                                        });
                                        usersLoading = false;
                                        return users;
                                    }
                                "
                                searchable-lazy
                                :searchable-placeholder="$t('common.search_field')"
                                :search-attributes="['firstname', 'lastname']"
                                block
                                :placeholder="$t('common.user', 2)"
                                trailing
                                value-attribute="userId"
                            >
                                <template #label="{ selected }">
                                    <span v-if="selected.length > 0" class="truncate">
                                        {{ usersToLabel(selected) }}
                                    </span>
                                </template>

                                <template #option="{ option: user }">
                                    <span class="truncate">
                                        {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                    </span>
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>

                                <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="data" :label="$t('common.data')">
                        <UInput
                            v-model="query.search"
                            type="text"
                            name="data"
                            block
                            :placeholder="$t('common.search')"
                            leading-icon="i-mdi-search"
                            :ui="{ icon: { trailing: { pointer: '' } } }"
                        >
                            <template #trailing>
                                <UButton
                                    v-show="query.search !== ''"
                                    color="gray"
                                    variant="link"
                                    icon="i-mdi-close"
                                    :padded="false"
                                    @click="query.search = ''"
                                />
                            </template>
                        </UInput>
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
                        <div class="flex flex-row flex-wrap gap-1">
                            <UFormGroup class="flex-1" name="service" :label="$t('common.service')">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.services"
                                        multiple
                                        searchable
                                        name="service"
                                        :placeholder="$t('common.service')"
                                        :options="grpcServices.map((s) => s.split('.').pop() ?? s)"
                                    >
                                        <template #option="{ option }">
                                            {{ option }}
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.service')]) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="method" :label="$t('common.method')">
                                <USelectMenu
                                    v-model="query.methods"
                                    multiple
                                    searchable
                                    name="method"
                                    :placeholder="$t('common.method')"
                                    :options="grpcMethods.filter((m) => query.services.some((s) => m.includes('.' + s + '/')))"
                                >
                                    <template #option="{ option }">
                                        {{ option.split('/').pop() }}
                                    </template>

                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.method')]) }}
                                    </template>
                                </USelectMenu>
                            </UFormGroup>

                            <UFormGroup class="flex-1" name="states" :label="$t('common.state')">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.states"
                                        multiple
                                        searchable
                                        name="states"
                                        :placeholder="$t('common.state')"
                                        :options="statesOptions"
                                        value-attribute="eventType"
                                    >
                                        <template #label="{ selected }">
                                            <span v-if="selected.length > 0">
                                                {{ statesToLabel(selected) }}
                                            </span>
                                        </template>

                                        <template #option="{ option }">
                                            {{ $t(`enums.settings.AuditLog.EventType.${EventType[option.eventType]}`) }}
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.state')]) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormGroup>
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
        v-model:expand="expand"
        class="flex-1"
        :loading="loading"
        :columns="columns"
        :rows="data?.logs"
        :empty-state="{
            icon: 'i-mdi-math-log',
            label: $t('common.not_found', [$t('common.entry', 2)]),
        }"
        sort-mode="manual"
    >
        <template #actions-data="{ row }">
            <UTooltip :text="$t('components.clipboard.clipboard_button.add')">
                <UButton variant="link" icon="i-mdi-content-copy" @click="addToClipboard(row)" />
            </UTooltip>
        </template>

        <template #createdAt-data="{ row }">
            <GenericTime :value="row.createdAt" type="long" />
        </template>

        <template #user-data="{ row }">
            <CitizenInfoPopover :user="row.user" />
        </template>

        <template #service-data="{ row }">
            <span class="dark:text-white"> {{ row.service }}/{{ row.method }} </span>
        </template>

        <template #state-data="{ row }">
            <UBadge :color="eventTypeToBadgeColor(row.state)">
                {{ $t(`enums.settings.AuditLog.EventType.${EventType[row.state]}`) }}
            </UBadge>
        </template>

        <template #expand="{ row }">
            <div class="px-2 py-1">
                <span v-if="!row.data">{{ $t('common.na') }}</span>
                <span v-else>
                    <VueJsonPretty
                        :data="JSON.parse(row.data!) as JSONDataType"
                        :show-icon="true"
                        :show-length="true"
                        :virtual="true"
                        :height="240"
                    />
                </span>
            </div>
        </template>
    </UTable>

    <Pagination v-model="query.page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
