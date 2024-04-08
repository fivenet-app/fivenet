<script lang="ts" setup>
import AttrView from '~/components/rector/attrs/AttrView.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { Role } from '~~/gen/ts/resources/permissions/permissions';
import { Job } from '~~/gen/ts/resources/users/jobs';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const notifications = useNotificatorStore();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const { data: roles, pending: loading, refresh, error } = useLazyAsyncData('rector-roles', () => getRoles());

async function getRoles(): Promise<Role[]> {
    try {
        const call = $grpc.getRectorClient().getRoles({
            lowestRank: true,
        });
        const { response } = await call;

        return response.roles;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const selectedJob = ref<Job | undefined>(undefined);

const availableJobs = computed(
    () => jobs.value?.filter((j) => (roles.value?.findIndex((r) => r.job === j.name) ?? -1) === -1) ?? [],
);

async function createRole(): Promise<void> {
    if (selectedJob.value === undefined || selectedJob.value?.name === undefined) {
        return;
    }

    try {
        const call = $grpc.getRectorClient().createRole({
            job: selectedJob.value?.name,
            grade: 1,
        });
        const { response } = await call;

        if (!response.role) {
            return;
        }

        notifications.add({
            title: { key: 'notifications.rector.role_created.title', parameters: {} },
            description: { key: 'notifications.rector.role_created.content', parameters: {} },
            type: 'success',
        });

        roles.value?.push(response.role!);

        selectedRole.value = response.role;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

onBeforeMount(async () => await listJobs());

const selectedRole = ref<Role | undefined>();

const sortedRoles = computed(() => roles.value?.sort((a, b) => (a.jobLabel ?? '').localeCompare(b.jobLabel ?? '')));

const columns = [
    {
        key: 'job',
        label: t('common.job'),
    },
    {
        key: 'actions',
        label: t('common.action', 2),
    },
];
</script>

<template>
    <div class="py-2">
        <div class="px-1 sm:px-2">
            <div class="flex flex-col lg:flex-row">
                <div class="mt-2 flow-root basis-1/3">
                    <div v-if="can('RectorService.CreateRole')" class="sm:flex sm:items-center">
                        <div class="sm:flex-auto">
                            <UForm :schema="undefined" :state="{}" @submit="refresh()">
                                <div class="flex flex-row gap-2">
                                    <UFormGroup class="flex-1" name="grade" :label="$t('common.job')">
                                        <USelectMenu v-model="selectedJob" :options="availableJobs" by="label">
                                            <template #label>
                                                <template v-if="selectedJob">
                                                    <span class="truncate"
                                                        >{{ selectedJob?.label }} ({{ selectedJob.name }})</span
                                                    >
                                                </template>
                                            </template>
                                            <template #option="{ option: job }">
                                                <span class="truncate">{{ job.label }} ({{ job.name }})</span>
                                            </template>
                                        </USelectMenu>
                                    </UFormGroup>

                                    <div class="flex flex-initial flex-col justify-end">
                                        <UButton :disabled="selectedJob === undefined" @click="createRole()">
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
                                :retry="refresh"
                            />
                            <template v-else>
                                <UTable :columns="columns" :rows="sortedRoles" :loading="loading">
                                    <template #job-data="{ row: role }">
                                        <span>{{ role.jobLabel }} ({{ role.job }})</span>
                                    </template>
                                    <template #actions-data="{ row: role }">
                                        <div class="text-right">
                                            <UButton
                                                class="place-self-end"
                                                variant="link"
                                                icon="i-mdi-eye"
                                                @click="selectedRole = role"
                                            />
                                        </div>
                                    </template>
                                </UTable>
                            </template>
                        </div>
                    </div>
                </div>
                <div class="ml-2 flex w-full basis-2/3">
                    <template v-if="selectedRole">
                        <AttrView :role-id="selectedRole.id" @deleted="refresh()" />
                    </template>
                    <template v-else>
                        <DataNoDataBlock icon="i-mdi-select" :message="$t('common.none_selected', [$t('common.job', 2)])" />
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>
