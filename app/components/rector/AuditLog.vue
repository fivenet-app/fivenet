<script lang="ts" setup>
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
import type { JSONDataType } from 'vue-json-pretty/types/utils';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { AuditEntry } from '~~/gen/ts/resources/rector/audit';
import { EventType } from '~~/gen/ts/resources/rector/audit';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { ViewAuditLogRequest, ViewAuditLogResponse } from '~~/gen/ts/services/rector/rector';

const { d, t } = useI18n();

const completorStore = useCompletorStore();

const schema = z.object({
    users: z.custom<UserShort>().array().max(5),
    from: z.date().optional(),
    to: z.date().optional(),
    method: z.string().max(64),
    service: z.string().max(64),
    search: z.string().max(128),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    users: [],
    method: '',
    service: '',
    search: '',
});

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `rector-audit-${page.value}-${query.from}-${query.to}-${query.method}-${query.service}-${query.search}-${query.users.map((v) => v.userId).join(':')}`,
    () => viewAuditLog(),
);

async function viewAuditLog(): Promise<ViewAuditLogResponse> {
    const req: ViewAuditLogRequest = {
        pagination: {
            offset: offset.value,
        },
        userIds: [],
    };

    req.userIds = query.users.map((v) => v.userId);

    if (query.from) {
        req.from = toTimestamp(query.from!);
    }
    if (query.to) {
        req.to = toTimestamp(query.to!);
    }

    if (query.method !== '') {
        req.method = query.method;
    }
    if (query.service !== '') {
        req.service = query.service;
    }

    if (query.search !== '') {
        req.search = query.search;
    }

    try {
        const call = getGRPCRectorClient().viewAuditLog(req);
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
    },
    {
        key: 'user',
        label: t('common.user', 1),
    },
    {
        key: 'service',
        label: `${t('common.service')}/${t('common.method')}`,
    },
    {
        key: 'state',
        label: t('common.state'),
    },
    {
        key: 'data',
        label: t('common.data'),
    },
];
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
                <div class="flex flex-row flex-wrap gap-2">
                    <UFormGroup name="from" :label="`${$t('common.time_range')} ${$t('common.from')}`" class="flex-1">
                        <DatePickerPopoverClient
                            v-model="query.from"
                            :popover="{ popper: { placement: 'bottom-start' } }"
                            :date-picker="{ clearable: true }"
                        />
                    </UFormGroup>

                    <UFormGroup name="to" :label="`${$t('common.time_range')} ${$t('common.to')}`" class="flex-1">
                        <DatePickerPopoverClient
                            v-model="query.to"
                            :popover="{ popper: { placement: 'bottom-start' } }"
                            :date-picker="{ clearable: true }"
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
                        <UInput v-model="query.service" type="text" name="service" :placeholder="$t('common.service')" block />
                    </UFormGroup>

                    <UFormGroup name="method" :label="$t('common.method')" class="flex-1">
                        <UInput v-model="query.method" type="text" name="method" block :placeholder="$t('common.method')" />
                    </UFormGroup>

                    <UFormGroup name="data" :label="$t('common.data')" class="flex-1">
                        <UInput v-model="query.search" type="text" name="data" block :placeholder="$t('common.search')" />
                    </UFormGroup>
                </div>
            </UForm>
        </template>
    </UDashboardToolbar>

    <div class="relative overflow-x-auto">
        <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.audit_log', 2)])" :retry="refresh" />

        <UTable v-else :loading="loading" :columns="columns" :rows="data?.logs">
            <template #actions-data="{ row }">
                <UButton
                    variant="link"
                    icon="i-mdi-content-copy"
                    :title="$t('components.clipboard.clipboard_button.add')"
                    @click="addToClipboard(row)"
                />
            </template>
            <template #createdAt-data="{ row }">
                <GenericTime :value="row.createdAt" type="long" />
            </template>
            <template #user-data="{ row }">
                <CitizenInfoPopover :user="row.user" />
            </template>
            <template #service-data="{ row }"> {{ row.service }}/{{ row.method }} </template>
            <template #state-data="{ row }">
                {{ EventType[row.state] }}
            </template>
            <template #data-data="{ row }">
                <span v-if="!row.data">{{ $t('common.na') }}</span>
                <span v-else>
                    <VueJsonPretty
                        :data="JSON.parse(row.data!) as JSONDataType"
                        :show-icon="true"
                        :show-length="true"
                        :virtual="true"
                        :height="160"
                    />
                </span>
            </template>
        </UTable>

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
    </div>
</template>
