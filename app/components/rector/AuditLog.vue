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
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { AuditEntry } from '~~/gen/ts/resources/rector/audit';
import { EventType } from '~~/gen/ts/resources/rector/audit';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { ViewAuditLogRequest, ViewAuditLogResponse } from '~~/gen/ts/services/rector/rector';
import { grpcMethods, grpcServices } from '~~/gen/ts/svcs';
import { eventTypeToBadgeColor } from './helpers';

const { $grpc } = useNuxtApp();

const { d, t } = useI18n();

const completorStore = useCompletorStore();

const schema = z.object({
    users: z.custom<UserShort>().array().max(5),
    date: z
        .object({
            start: z.date(),
            end: z.date(),
        })
        .optional(),
    services: z.string().max(64).array().max(10),
    methods: z.string().max(64).array().max(10),
    search: z.string().max(64),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    users: [],
    services: [],
    methods: [],
    search: '',
});

const usersLoading = ref(false);

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const sort = useRouteQueryObject<TableSortable>('sort', {
    column: 'createdAt',
    direction: 'desc',
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `rector-audit-${sort.value.column}:${sort.value.direction}-${page.value}-${query.date?.start}-${query.date?.end}-${query.methods}-${query.services}-${query.search}-${query.users.map((v) => v.userId).join(':')}`,
    () => viewAuditLog(),
    {
        watch: [sort],
    },
);

async function viewAuditLog(): Promise<ViewAuditLogResponse> {
    const req: ViewAuditLogRequest = {
        pagination: {
            offset: offset.value,
        },
        sort: sort.value,
        userIds: [],
        services: query.services,
        // Make sure to remove the service from the beginning
        methods: query.methods.map((m) => m.split('/').pop() ?? m),
    };

    req.userIds = query.users.map((v) => v.userId);

    if (query.date) {
        req.from = toTimestamp(query.date.start);
        req.to = toTimestamp(query.date.end);
    }

    if (query.search !== '') {
        req.search = query.search;
    }

    try {
        const call = $grpc.rector.rector.viewAuditLog(req);
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), {
    debounce: 200,
    maxWait: 1250,
});

const notifications = useNotificatorStore();

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
        title: { key: 'notifications.rector.audit_log.title', parameters: {} },
        description: { key: 'notifications.rector.audit_log.content', parameters: {} },
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
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
                <div class="flex flex-row flex-wrap gap-2">
                    <UFormGroup name="date" :label="$t('common.time_range')" class="flex-1">
                        <DateRangePickerPopoverClient
                            v-model="query.date"
                            mode="date"
                            class="flex-1"
                            :popover="{ class: 'flex-1' }"
                            :date-picker="{
                                mode: 'dateTime',
                                disabledDates: [{ start: addDays(new Date(), 1), end: null }],
                                is24Hr: true,
                                clearable: true,
                            }"
                        />
                    </UFormGroup>

                    <UFormGroup name="user" :label="$t('common.user')" class="flex-1">
                        <ClientOnly>
                            <USelectMenu
                                v-model="query.users"
                                multiple
                                :searchable="
                                    async (query: string) => {
                                        usersLoading = true;
                                        const users = await completorStore.completeCitizens({
                                            search: query,
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
                                by="userId"
                            >
                                <template #label>
                                    <span v-if="query.users.length" class="truncate">
                                        {{ usersToLabel(query.users) }}
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

                    <UFormGroup name="service" :label="$t('common.service')" class="flex-1">
                        <ClientOnly>
                            <USelectMenu
                                v-model="query.services"
                                multiple
                                searchable
                                name="service"
                                :placeholder="$t('common.service')"
                                :options="grpcServices"
                            >
                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.service')]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup name="method" :label="$t('common.method')" class="flex-1">
                        <USelectMenu
                            v-model="query.methods"
                            multiple
                            searchable
                            name="method"
                            :placeholder="$t('common.method')"
                            :options="grpcMethods.filter((m) => query.services.some((s) => m.startsWith(s + '/')))"
                        >
                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.method')]) }}
                            </template>
                        </USelectMenu>
                    </UFormGroup>

                    <UFormGroup name="data" :label="$t('common.data')" class="flex-1">
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
        :loading="loading"
        :columns="columns"
        :rows="data?.logs"
        :empty-state="{
            icon: 'i-mdi-math-log',
            label: $t('common.not_found', [$t('common.entry', 2)]),
        }"
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
                {{ $t(`enums.rector.AuditLog.EventType.${EventType[row.state]}`) }}
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

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
