<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import AttrView from '~/components/rector/attrs/AttrView.vue';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Role } from '~~/gen/ts/resources/permissions/permissions';
import type { Job } from '~~/gen/ts/resources/users/jobs';

const { t } = useI18n();

const { can } = useAuth();

const notifications = useNotificatorStore();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const { data: roles, pending: loading, refresh, error } = useLazyAsyncData('rector-roles', () => getRoles());

async function getRoles(): Promise<Role[]> {
    try {
        const call = getGRPCRectorClient().getRoles({
            lowestRank: true,
        });
        const { response } = await call;

        return response.roles;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const schema = z.object({
    job: z.custom<Job>().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    job: undefined,
});

const availableJobs = computed(
    () => jobs.value?.filter((j) => (roles.value?.findIndex((r) => r.job === j.name) ?? -1) === -1) ?? [],
);

async function createRole(): Promise<void> {
    if (state.job === undefined || state.job?.name === undefined) {
        return;
    }

    try {
        const call = getGRPCRectorClient().createRole({
            job: state.job?.name,
            grade: 1,
        });
        const { response } = await call;

        if (!response.role) {
            return;
        }

        notifications.add({
            title: { key: 'notifications.rector.role_created.title', parameters: {} },
            description: { key: 'notifications.rector.role_created.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        roles.value?.push(response.role!);

        selectedRole.value = response.role;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

onBeforeMount(async () => await listJobs());

const selectedRole = ref<Role | undefined>();

const sortedRoles = computed(() => [...(roles.value ?? [])].sort((a, b) => (a.jobLabel ?? '').localeCompare(b.jobLabel ?? '')));

const columns = [
    {
        key: 'job',
        label: t('common.job'),
    },
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
];

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await createRole().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <div class="relative overflow-x-auto">
        <div class="px-1 sm:px-2">
            <div class="flex flex-col lg:flex-row">
                <div class="mt-2 flow-root basis-1/3">
                    <div v-if="can('RectorService.CreateRole').value" class="sm:flex sm:items-center">
                        <div class="sm:flex-auto">
                            <UForm :schema="schema" :state="state" @submit="refresh()">
                                <div class="flex flex-row gap-2">
                                    <UFormGroup class="flex-1" name="grade" :label="$t('common.job')">
                                        <ClientOnly>
                                            <USelectMenu
                                                v-model="state.job"
                                                :options="availableJobs"
                                                by="name"
                                                searchable
                                                :searchable-placeholder="$t('common.search_field')"
                                            >
                                                <template #label>
                                                    <template v-if="state.job">
                                                        <span class="truncate"
                                                            >{{ state.job?.label }} ({{ state.job.name }})</span
                                                        >
                                                    </template>
                                                </template>

                                                <template #option="{ option: job }">
                                                    <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                                </template>
                                            </USelectMenu>
                                        </ClientOnly>
                                    </UFormGroup>

                                    <div class="flex flex-initial flex-col justify-end">
                                        <UButton
                                            :disabled="state.job === undefined || !canSubmit"
                                            :loading="!canSubmit"
                                            icon="i-mdi-plus"
                                            @click="onSubmitThrottle"
                                        >
                                            {{ $t('common.create') }}
                                        </UButton>
                                    </div>
                                </div>
                            </UForm>
                        </div>
                    </div>

                    <div class="-my-2 mx-0 overflow-x-auto">
                        <div class="inline-block min-w-full px-1 py-2 align-middle">
                            <DataErrorBlock
                                v-if="error"
                                :title="$t('common.unable_to_load', [$t('common.job', 2)])"
                                :error="error"
                                :retry="refresh"
                            />
                            <UTable
                                v-else
                                :columns="columns"
                                :rows="sortedRoles"
                                :loading="loading"
                                :empty-state="{
                                    icon: 'i-mdi-account-group',
                                    label: $t('common.not_found', [$t('common.role', 2)]),
                                }"
                            >
                                <template #job-data="{ row: role }">
                                    <div class="text-gray-900 dark:text-white">{{ role.jobLabel }} ({{ role.job }})</div>
                                </template>

                                <template #actions-data="{ row: role }">
                                    <div class="text-right">
                                        <UTooltip :text="$t('common.show')">
                                            <UButton
                                                class="place-self-end"
                                                variant="link"
                                                icon="i-mdi-eye"
                                                @click="selectedRole = role"
                                            />
                                        </UTooltip>
                                    </div>
                                </template>
                            </UTable>
                        </div>
                    </div>
                </div>
                <div class="mt-0 mt-4 w-full basis-2/3 lg:ml-2">
                    <AttrView
                        v-if="selectedRole"
                        :role-id="selectedRole.id"
                        @deleted="
                            selectedRole = undefined;
                            refresh();
                        "
                    />
                    <template v-else>
                        <DataNoDataBlock icon="i-mdi-select" :message="$t('common.none_selected', [$t('common.job', 2)])" />
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>
