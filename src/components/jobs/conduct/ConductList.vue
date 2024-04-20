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

const { $grpc } = useNuxtApp();

const completorStore = useCompletorStore();

const modal = useModal();

const slideover = useSlideover();

const schema = z.object({
    types: z.nativeEnum(ConductType).array().max(10),
    showExpired: z.boolean(),
    user: z.custom<User>().optional(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    types: [],
    showExpired: false,
});

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, refresh, error } = useLazyAsyncData(`jobs-conduct-${page.value}`, () => listConductEntries());

async function listConductEntries(): Promise<ListConductEntriesResponse> {
    const userIds = props.userId ? [props.userId] : query.user ? [query.user.userId] : [];
    try {
        const call = $grpc.getJobsConductClient().listConductEntries({
            pagination: {
                offset: offset.value,
            },
            types: query.types,
            userIds: userIds,
            showExpired: query.showExpired,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteConductEntry(id: string): Promise<void> {
    try {
        const call = $grpc.getJobsConductClient().deleteConductEntry({ id });
        await call;

        refresh();
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

watch(offset, async () => refresh());
watchDebounced(query, () => refresh(), { debounce: 200, maxWait: 1250 });

async function updateEntryInPlace(entry: ConductEntry): Promise<void> {
    if (data.value === null) {
        return refresh();
    }

    const idx = data.value.entries.findIndex((e) => e.id === entry.id);
    if (idx !== undefined && idx > -1) {
        data.value.entries[idx] = entry;
    }
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
        key: 'createdAt',
        label: t('common.created_at'),
    },
    {
        key: 'expiresAt',
        label: t('common.expires_at'),
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
                    <UFormGroup v-if="hideUserSearch !== true" name="user" :label="$t('common.target')" class="flex-1">
                        <UInputMenu
                            ref="input"
                            v-model="query.user"
                            nullable
                            :search="
                                async (query: string) => {
                                    usersLoading = true;
                                    const colleagues = await completorStore.listColleagues({
                                        search: query,
                                    });
                                    usersLoading = false;
                                    return colleagues;
                                }
                            "
                            :search-attributes="['firstname', 'lastname']"
                            block
                            :placeholder="$t('common.target')"
                            trailing
                            by="firstname"
                            :searchable-placeholder="$t('common.search_field')"
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
                        </UInputMenu>
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
                        v-if="can('JobsConductService.CreateConductEntry')"
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
    >
        <template #createdAt-data="{ row: conduct }">
            <GenericTime :value="conduct.createdAt" />
            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.expires_at') }}</dt>
                <dd class="mt-1 truncate">
                    <GenericTime v-if="conduct.expiresAt" class="font-semibold" :value="conduct.expiresAt" />
                    <span v-else>
                        {{ $t('components.jobs.conduct.List.no_expiration') }}
                    </span>
                </dd>
            </dl>
        </template>
        <template #expiresAt-data="{ row: conduct }">
            <GenericTime v-if="conduct.expiresAt" class="font-semibold" type="date" :value="conduct.expiresAt" />
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
                    <ColleagueInfoPopover :user="conduct.creator" />
                </dd>
            </dl>
        </template>
        <template #creator-data="{ row: conduct }">
            <ColleagueInfoPopover :user="conduct.creator" :hide-props="true" />
        </template>
        <template #actions-data="{ row: conduct }">
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
                    v-if="can('JobsConductService.UpdateConductEntry')"
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
                    v-if="can('JobsConductService.DeleteConductEntry')"
                    variant="link"
                    icon="i-mdi-trash-can"
                    @click="
                        modal.open(ConfirmModal, {
                            confirm: async () => deleteConductEntry(conduct.id),
                        })
                    "
                />
            </UButtonGroup>
        </template>
    </UTable>

    <Pagination v-model="page" :pagination="data?.pagination" />
</template>
