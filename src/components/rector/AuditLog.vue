<script lang="ts" setup>
import { format } from 'date-fns';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { useCompletorStore } from '~/store/completor';
import { UserShort } from '~~/gen/ts/resources/users/users';
import { ViewAuditLogRequest, ViewAuditLogResponse } from '~~/gen/ts/services/rector/rector';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { AuditEntry, EventType } from '~~/gen/ts/resources/rector/audit';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';
import { useNotificatorStore } from '~/store/notificator';
import DatePicker from '~/components/partials/DatePicker.vue';
import Pagination from '../partials/Pagination.vue';

const { $grpc } = useNuxtApp();

const { d, t } = useI18n();

const completorStore = useCompletorStore();

const schema = z.object({
    from: z.date().optional(),
    to: z.date().optional(),
    method: z.string().max(64),
    service: z.string().max(64),
    search: z.string().max(128),
});

type Schema = z.output<typeof schema>;

const query = ref<Schema>({
    method: '',
    service: '',
    search: '',
});

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`rector-audit-${page.value}`, () => viewAuditLog());

async function viewAuditLog(): Promise<ViewAuditLogResponse> {
    const req: ViewAuditLogRequest = {
        pagination: {
            offset: offset.value,
        },
        userIds: [],
    };

    if (selectedCitizens.value.length > 0) {
        const users: number[] = [];
        selectedCitizens.value?.forEach((v) => users.push(v.userId));
        req.userIds = users;
    }

    if (query.value.from) {
        req.from = toTimestamp(query.value.from!);
    }
    if (query.value.to) {
        req.to = toTimestamp(query.value.to!);
    }

    if (query.value.method !== '') {
        req.method = query.value.method;
    }
    if (query.value.service !== '') {
        req.service = query.value.service;
    }

    if (query.value.search !== '') {
        req.search = query.value.search;
    }

    try {
        const call = $grpc.getRectorClient().viewAuditLog(req);
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const selectedCitizens = ref<UserShort[]>([]);

function charsGetDisplayValue(chars: UserShort[]): string {
    return chars.map((c) => `${c?.firstname} ${c?.lastname} (${c?.dateofbirth})`).join(', ');
}

watch(offset, async () => refresh());
watchDebounced(query.value, async () => refresh(), {
    debounce: 200,
    maxWait: 1250,
});
watchDebounced(selectedCitizens.value, async () => refresh(), {
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
${jsonStringify(jsonParse(logEntry.data!), 2)}
\`\`\`
`;
    } else {
        text += `Data: N/A
`;
    }

    notifications.add({
        title: { key: 'notifications.rector.audit_log.title', parameters: {} },
        description: { key: 'notifications.rector.audit_log.content', parameters: {} },
        type: 'info',
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
                        <UPopover :popper="{ placement: 'bottom-start' }">
                            <UButton
                                variant="outline"
                                color="gray"
                                block
                                icon="i-mdi-calendar-month"
                                :label="query.from ? format(query.from, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
                            />

                            <template #panel="{ close }">
                                <DatePicker v-model="query.from" @close="close" />
                            </template>
                        </UPopover>
                    </UFormGroup>

                    <UFormGroup name="to" :label="`${$t('common.time_range')} ${$t('common.to')}`" class="flex-1">
                        <UPopover :popper="{ placement: 'bottom-start' }">
                            <UButton
                                variant="outline"
                                color="gray"
                                block
                                icon="i-mdi-calendar-month"
                                :label="query.to ? format(query.to, 'dd.MM.yyyy') : 'dd.mm.yyyy'"
                            />

                            <template #panel="{ close }">
                                <DatePicker v-model="query.to" @close="close" />
                            </template>
                        </UPopover>
                    </UFormGroup>

                    <UFormGroup name="user" :label="$t('common.user')" class="flex-1">
                        <USelectMenu
                            v-model="selectedCitizens"
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
                            :search-attributes="['firstname', 'lastname']"
                            block
                            :placeholder="selectedCitizens ? charsGetDisplayValue(selectedCitizens) : $t('common.user', 2)"
                            trailing
                            by="userId"
                        >
                            <template #option="{ option: user }">
                                {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                            </template>
                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>
                            <template #empty> {{ $t('common.not_found', [$t('common.creator', 2)]) }} </template>
                        </USelectMenu>
                    </UFormGroup>

                    <UFormGroup name="service" :label="$t('common.service')" class="flex-1">
                        <UInput
                            v-model="query.service"
                            type="text"
                            name="service"
                            :placeholder="$t('common.service')"
                            block
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>

                    <UFormGroup name="method" :label="$t('common.method')" class="flex-1">
                        <UInput
                            v-model="query.method"
                            type="text"
                            name="method"
                            block
                            :placeholder="$t('common.method')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>

                    <UFormGroup name="data" :label="$t('common.data')" class="flex-1">
                        <UInput
                            v-model="query.search"
                            type="text"
                            name="data"
                            block
                            :placeholder="$t('common.search')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>
                </div>
            </UForm>
        </template>
    </UDashboardToolbar>

    <div class="inline-block w-full max-w-full px-1 py-2 align-middle">
        <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.audit_log', 2)])" :retry="refresh" />

        <UTable v-else :loading="pending" :columns="columns" :rows="data?.logs">
            <template #actions-data="{ row }">
                <UButton
                    variant="link"
                    icon="i-mdi-clipboard-plus"
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
                <span v-if="!row.data">N/A</span>
                <span v-else>
                    <VueJsonPretty
                        :data="jsonParse(row.data!)"
                        :show-icon="true"
                        :show-length="true"
                        :virtual="true"
                        :height="160"
                    />
                </span>
            </template>
        </UTable>

        <Pagination v-model="page" :pagination="data?.pagination" />
    </div>
</template>