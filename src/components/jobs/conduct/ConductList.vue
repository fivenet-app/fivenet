<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { User } from '~~/gen/ts/resources/users/users';
import ConductCreateOrUpdateModal from '~/components/jobs/conduct/ConductCreateOrUpdateModal.vue';
import type { ListConductEntriesResponse } from '~~/gen/ts/services/jobs/conduct';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { useCompletorStore } from '~/store/completor';
import { conductTypesToBadgeColor, conductTypesToBGColor } from './helpers';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import ColleagueInfoPopover from '../colleagues/ColleagueInfoPopover.vue';
import ConductViewSlideover from './ConductViewSlideover.vue';

const props = defineProps<{
    userId?: number;
    hideUserSearch?: boolean;
}>();

const { t } = useI18n();

const completorStore = useCompletorStore();

const modal = useModal();

const slideover = useSlideover();

const schema = z.object({
    id: z.string().max(16).optional(),
    types: z.nativeEnum(ConductType).array().max(10),
    showExpired: z.boolean(),
    user: z.custom<User>().optional(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    id: undefined,
    types: [],
    showExpired: false,
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
    `jobs-conduct-${page.value}-${query.types.join(',')}-${query.showExpired}-${query.id}`,
    () => listConductEntries(),
    {
        transform: (input) => ({ ...input, entries: wrapRows(input?.entries, columns) }),
    },
);

async function listConductEntries(): Promise<ListConductEntriesResponse> {
    const entryIds = [];
    if (query.id) {
        const id = query.id.trim().replaceAll('-', '').replace(/\D/g, '');
        if (id.length > 0) {
            entryIds.push(id);
        }
    }

    const userIds = props.userId ? [props.userId] : query.user ? [query.user.userId] : [];
    try {
        const call = getGRPCJobsConductClient().listConductEntries({
            pagination: {
                offset: offset.value,
            },
            types: query.types,
            userIds: userIds,
            showExpired: query.showExpired,
            ids: entryIds,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteConductEntry(id: string): Promise<void> {
    try {
        const call = getGRPCJobsConductClient().deleteConductEntry({ id });
        await call;

        refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), { debounce: 200, maxWait: 1250 });

async function updateEntryInPlace(entry: ConductEntry): Promise<void> {
    if (data.value === null) {
        return refresh();
    }

    const idx = data.value.entries.findIndex((e) => e.id === entry.id);
    if (idx !== undefined && idx > -1) {
        data.value.entries.splice(idx, 1, entry);
    }

    refresh();
}

type CType = { status: ConductType };

const cTypes = ref<CType[]>([
    { status: ConductType.NOTE },
    { status: ConductType.NEUTRAL },
    { status: ConductType.POSITIVE },
    { status: ConductType.NEGATIVE },
    { status: ConductType.WARNING },
    { status: ConductType.SUSPENSION },
]);

const columns = [
    {
        key: 'id',
        label: t('common.id'),
    },
    {
        key: 'createdAt',
        label: t('common.created_at'),
    },
    {
        key: 'expiresAt',
        label: t('common.expires_at'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    {
        key: 'type',
        label: t('common.type'),
    },
    {
        key: 'message',
        label: t('common.message'),
    },
    {
        key: 'target',
        label: t('common.target'),
    },
    {
        key: 'creator',
        label: t('common.creator'),
        class: 'hidden lg:table-cell',
        rowClass: 'hidden lg:table-cell',
    },
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
];

const input = ref<{ input: HTMLInputElement }>();

defineShortcuts({
    '/': () => {
        input.value?.input?.focus();
    },
});
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
                <div class="flex flex-row gap-2">
                    <UFormGroup v-if="hideUserSearch !== true" name="user" :label="$t('common.search')" class="flex-1">
                        <USelectMenu
                            ref="input"
                            v-model="query.user"
                            :searchable="
                                async (query: string) => {
                                    usersLoading = true;
                                    const colleagues = await completorStore.listColleagues({
                                        search: query,
                                    });
                                    usersLoading = false;
                                    return colleagues;
                                }
                            "
                            searchable-lazy
                            :searchable-placeholder="$t('common.search_field')"
                            :search-attributes="['firstname', 'lastname']"
                            block
                            :placeholder="$t('common.colleague')"
                            trailing
                            by="userId"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                            @keydown.esc="$event.target.blur()"
                        >
                            <template #option="{ option: user }">
                                {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                            </template>
                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.creator', 2)]) }}
                            </template>

                            <template #trailing>
                                <UKbd value="/" />
                            </template>
                        </USelectMenu>
                    </UFormGroup>

                    <UFormGroup name="types" :label="$t('common.type')" class="flex-1">
                        <USelectMenu
                            v-model="query.types"
                            multiple
                            nullable
                            :options="cTypes"
                            value-attribute="status"
                            :placeholder="$t('common.na')"
                            :searchable-placeholder="$t('common.search_field')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        >
                            <template #option="{ option }">
                                <span class="truncate" :class="conductTypesToBGColor(option.status)">
                                    {{ $t(`enums.jobs.ConductType.${ConductType[option.status]}`) }}
                                </span>
                            </template>
                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>
                            <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
                        </USelectMenu>
                    </UFormGroup>

                    <UFormGroup name="id" :label="$t('common.id')" class="flex-initial">
                        <UInput v-model="query.id" type="text" name="id" :placeholder="$t('common.id')" />
                    </UFormGroup>

                    <UFormGroup
                        name="showExpired"
                        :label="$t('components.jobs.conduct.List.show_expired')"
                        class="flex-initial"
                    >
                        <UToggle v-model="query.showExpired">
                            <span class="sr-only">
                                {{ $t('components.jobs.conduct.List.show_expired') }}
                            </span>
                        </UToggle>
                    </UFormGroup>

                    <UFormGroup
                        v-if="can('JobsConductService.CreateConductEntry').value"
                        :label="$t('common.create')"
                        class="flex-initial"
                    >
                        <UButton
                            @click="
                                modal.open(ConductCreateOrUpdateModal, {
                                    onCreated: ($event) => data?.entries.unshift($event),
                                    onUpdated: ($event) => updateEntryInPlace($event),
                                })
                            "
                        >
                            {{ $t('common.create') }}
                        </UButton>
                    </UFormGroup>
                </div>
            </UForm>
        </template>
    </UDashboardToolbar>

    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.conduct_register')])" :retry="refresh" />
    <UTable
        v-else
        :loading="loading"
        :columns="columns"
        :rows="data?.entries"
        :empty-state="{ icon: 'i-mdi-list-status', label: $t('common.not_found', [$t('common.entry', 2)]) }"
        class="flex-1"
    >
        <template #createdAt-data="{ row: conduct }">
            <GenericTime :value="conduct.createdAt" />
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.expires_at') }}</dt>
                <dd class="mt-1 truncate">
                    <GenericTime v-if="conduct.expiresAt?.value" class="font-semibold" :value="conduct.expiresAt.value" />
                    <span v-else>
                        {{ $t('components.jobs.conduct.List.no_expiration') }}
                    </span>
                </dd>
            </dl>
        </template>
        <template #expiresAt-data="{ row: conduct }">
            <GenericTime v-if="conduct.expiresAt?.value" class="font-semibold" type="date" :value="conduct.expiresAt.value" />
            <span v-else>
                {{ $t('components.jobs.conduct.List.no_expiration') }}
            </span>
        </template>
        <template #type-data="{ row: conduct }">
            <UBadge :color="conductTypesToBadgeColor(conduct.type)">
                {{ $t(`enums.jobs.ConductType.${ConductType[conduct.type ?? 0]}`) }}
            </UBadge>
        </template>
        <template #message-data="{ row: conduct }">
            <p class="line-clamp-2 w-full max-w-sm whitespace-normal break-all hover:line-clamp-6">
                {{ conduct.message }}
            </p>
        </template>
        <template #target-data="{ row: conduct }">
            <ColleagueInfoPopover :user="conduct.targetUser" />
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.creator') }}</dt>
                <dd class="mt-1 truncate">
                    <ColleagueInfoPopover :user="conduct.creator.value" />
                </dd>
            </dl>
        </template>
        <template #creator-data="{ row: conduct }">
            <ColleagueInfoPopover :user="conduct.creator.value" :hide-props="true" />
        </template>
        <template #actions-data="{ row: conduct }">
            <div :key="conduct.id">
                <UButtonGroup class="inline-flex">
                    <UButton
                        variant="link"
                        icon="i-mdi-eye"
                        @click="
                            slideover.open(ConductViewSlideover, {
                                entry: conduct,
                            })
                        "
                    />

                    <UButton
                        v-if="can('JobsConductService.UpdateConductEntry').value"
                        variant="link"
                        icon="i-mdi-pencil"
                        @click="
                            modal.open(ConductCreateOrUpdateModal, {
                                entry: conduct,
                                userId: userId,
                                onCreated: ($event) => data?.entries.unshift($event),
                                onUpdated: ($event) => updateEntryInPlace($event),
                            })
                        "
                    />

                    <UButton
                        v-if="can('JobsConductService.DeleteConductEntry').value"
                        variant="link"
                        icon="i-mdi-trash-can"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteConductEntry(conduct.id),
                            })
                        "
                    />
                </UButtonGroup>
            </div>
        </template>
    </UTable>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
