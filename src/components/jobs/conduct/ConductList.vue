<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { User } from '~~/gen/ts/resources/users/users';
import ConductCreateOrUpdateModal from '~/components/jobs/conduct/ConductCreateOrUpdateModal.vue';
import type { ListConductEntriesResponse } from '~~/gen/ts/services/jobs/conduct';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { useCompletorStore } from '~/store/completor';
import { conductTypesToBGColor, conductTypesToRingColor, conductTypesToTextColor } from './helpers';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';

const props = defineProps<{
    userId?: number;
    hideUserSearch?: boolean;
}>();

const { t } = useI18n();

const { $grpc } = useNuxtApp();

type CType = { status: ConductType };

const completorStore = useCompletorStore();

const modal = useModal();

const query = ref<{ types: CType[]; showExpired?: boolean; user?: User }>({
    types: [],
    showExpired: false,
});

const usersLoading = ref(false);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending: loading, refresh, error } = useLazyAsyncData(`jobs-conduct-${page.value}`, () => listConductEntries());

async function listConductEntries(): Promise<ListConductEntriesResponse> {
    const userIds = props.userId ? [props.userId] : query.value.user ? [query.value.user.userId] : [];
    try {
        const call = $grpc.getJobsConductClient().listConductEntries({
            pagination: {
                offset: offset.value,
            },
            types: query.value.types.map((t) => t.status),
            userIds: userIds,
            showExpired: query.value.showExpired,
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
watchDebounced(query.value, () => refresh(), { debounce: 600, maxWait: 1400 });

function updateEntryInPlace(entry: ConductEntry): void {
    if (data.value === null) {
        refresh();
        return;
    }

    const idx = data.value.entries.findIndex((e) => e.id === entry.id);
    if (idx !== undefined && idx > -1) {
        data.value.entries[idx] = entry;
    }
}

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
    },
];
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <UForm :state="{}">
                        <div class="flex flex-row gap-2">
                            <UFormGroup v-if="hideUserSearch !== true" name="user" class="flex-1" :label="$t('common.target')">
                                <UInputMenu
                                    v-model="query.user"
                                    :nullable="true"
                                    :search="
                                        async (query: string) => {
                                            usersLoading = true;
                                            const colleagues = await completorStore.listColleagues({
                                                pagination: { offset: 0 },
                                                searchName: query,
                                            });
                                            usersLoading = false;
                                            return colleagues;
                                        }
                                    "
                                    :search-attributes="['firstname', 'lastname']"
                                    block
                                    :placeholder="
                                        query.user
                                            ? `${query.user?.firstname} ${query.user?.lastname} (${query.user?.dateofbirth})`
                                            : $t('common.target')
                                    "
                                    trailing
                                    by="userId"
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
                                </UInputMenu>
                            </UFormGroup>

                            <UFormGroup name="types" class="flex-1" :label="$t('common.type')">
                                <USelectMenu
                                    v-model="query.types"
                                    multiple
                                    nullable
                                    :options="cTypes"
                                    :placeholder="
                                        query.types
                                            ? query.types
                                                  .map((ct: CType) =>
                                                      $t(`enums.jobs.ConductType.${ConductType[ct.status ?? 0]}`),
                                                  )
                                                  .join(', ')
                                            : $t('common.na')
                                    "
                                >
                                    <template #option="{ option: cType }">
                                        <span :class="conductTypesToBGColor(cType.status)">
                                            {{ $t(`enums.jobs.ConductType.${ConductType[cType.status]}`) }}
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
                                class="flex-initial"
                                :label="$t('components.jobs.conduct.List.show_expired')"
                            >
                                <UToggle v-model="query.showExpired">
                                    <span class="sr-only">
                                        {{ $t('components.jobs.conduct.List.show_expired') }}
                                    </span>
                                </UToggle>
                            </UFormGroup>

                            <UFormGroup
                                v-if="can('JobsConductService.CreateConductEntry')"
                                class="flex-initial"
                                :label="$t('common.create')"
                            >
                                <UButton @click="modal.open(ConductCreateOrUpdateModal, {})">
                                    {{ $t('common.create') }}
                                </UButton>
                            </UFormGroup>
                        </div>
                    </UForm>
                </div>
            </div>

            <div class="-my-2 mx-0 overflow-x-auto">
                <div class="inline-block min-w-full px-1 py-2 align-middle">
                    <DataErrorBlock
                        v-if="error"
                        :title="$t('common.unable_to_load', [$t('common.conduct_register')])"
                        :retry="refresh"
                    />
                    <UTable
                        v-else
                        :loading="loading"
                        :columns="columns"
                        :rows="data?.entries"
                        :empty-state="{ icon: 'i-mdi-car', label: $t('common.not_found', [$t('common.entry', 2)]) }"
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
                            <GenericTime v-if="conduct.expiresAt" class="font-semibold" :value="conduct.expiresAt" />
                            <span v-else>
                                {{ $t('components.jobs.conduct.List.no_expiration') }}
                            </span>
                        </template>
                        <template #type-data="{ row: conduct }">
                            <div
                                class="rounded-md px-2 py-1 text-base font-medium ring-1 ring-inset"
                                :class="[
                                    conductTypesToBGColor(conduct.type),
                                    conductTypesToRingColor(conduct.type),
                                    conductTypesToTextColor(conduct.type),
                                ]"
                            >
                                {{ $t(`enums.jobs.ConductType.${ConductType[conduct.type ?? (0 as number)]}`) }}
                            </div>
                        </template>
                        <template #message-data="{ row: conduct }">
                            <p class="line-clamp-2 w-full max-w-sm whitespace-normal break-all hover:line-clamp-6">
                                {{ conduct.message }}
                            </p>
                        </template>
                        <template #target-data="{ row: conduct }">
                            <CitizenInfoPopover :user="conduct.targetUser" />
                            <dl class="font-normal lg:hidden">
                                <dt class="sr-only">{{ $t('common.creator') }}</dt>
                                <dd class="mt-1 truncate">
                                    <CitizenInfoPopover :user="conduct.creator" />
                                </dd>
                            </dl>
                        </template>
                        <template #creator-data="{ row: conduct }">
                            <CitizenInfoPopover :user="conduct.creator" />
                        </template>
                        <template #actions-data="{ row: conduct }">
                            <UButtonGroup class="inline-flex">
                                <UButton
                                    v-if="can('JobsConductService.UpdateConductEntry')"
                                    variant="link"
                                    icon="i-mdi-pencil"
                                    @click="
                                        modal.open(ConductCreateOrUpdateModal, {
                                            entry: conduct,
                                            userId: userId,
                                            onCreated: ($event) => data?.entries.unshift($event),
                                            onUpdate: ($event) => updateEntryInPlace($event),
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

                    <div class="flex justify-end border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
                        <UPagination
                            v-model="page"
                            :page-count="data?.pagination?.pageSize ?? 0"
                            :total="data?.pagination?.totalCount ?? 0"
                        />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
