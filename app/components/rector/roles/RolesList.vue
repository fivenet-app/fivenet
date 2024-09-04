<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import RoleView from '~/components/rector/roles/RoleView.vue';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { Role } from '~~/gen/ts/resources/permissions/permissions';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';

const { t } = useI18n();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();
const { getJobByName } = completorStore;

const { data: roles, pending: loading, refresh, error } = useLazyAsyncData('rector-roles', () => getRoles());

async function getRoles(): Promise<Role[]> {
    try {
        const call = getGRPCRectorClient().getRoles({});
        const { response } = await call;

        return response.roles;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const job = ref<Job | undefined>();
watchOnce(roles, async () => (job.value = await getJobByName(activeChar.value!.job)));

const schema = z.object({
    jobGrade: z.custom<JobGrade>().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    jobGrade: undefined,
});

const availableJobGrades = computed(
    () => job.value?.grades.filter((g) => (roles.value?.findIndex((r) => r.grade === g.grade) ?? -1) === -1) ?? [],
);

async function createRole(): Promise<void> {
    if (state.jobGrade === undefined || state.jobGrade.grade <= 0) {
        return;
    }

    try {
        const call = getGRPCRectorClient().createRole({
            job: activeChar.value!.job,
            grade: state.jobGrade.grade,
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

const sortedRoles = computed(() => roles.value?.sort((a, b) => a.grade - b.grade));

const selectedRole = ref<Role | undefined>();

const columns = [
    {
        key: 'rank',
        label: t('common.rank'),
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
    <div class="py-2">
        <div class="px-1 sm:px-2">
            <div class="flex flex-col lg:flex-row">
                <div class="mt-2 flow-root basis-1/3">
                    <template v-if="can('RectorService.CreateRole').value">
                        <div class="sm:flex sm:items-center">
                            <div class="sm:flex-auto">
                                <UForm :schema="schema" :state="state" @submit="refresh()">
                                    <div class="flex flex-row gap-2">
                                        <UFormGroup name="grade" :label="$t('common.job_grade')" class="flex-1">
                                            <USelectMenu
                                                v-model="state.jobGrade"
                                                :options="availableJobGrades"
                                                by="grade"
                                                searchable
                                                :searchable-placeholder="$t('common.search_field')"
                                            >
                                                <template #label>
                                                    <template v-if="state.jobGrade">
                                                        <span class="truncate"
                                                            >{{ state.jobGrade?.label }} ({{ state.jobGrade?.grade }})</span
                                                        >
                                                    </template>
                                                </template>
                                                <template #option="{ option: jobGrade }">
                                                    <span class="truncate">{{ jobGrade.label }} ({{ jobGrade.grade }})</span>
                                                </template>
                                            </USelectMenu>
                                        </UFormGroup>

                                        <div class="flex flex-initial flex-col justify-end">
                                            <UButton
                                                :disabled="
                                                    state.jobGrade === undefined || state.jobGrade!.grade <= 0 || !canSubmit
                                                "
                                                :loading="!canSubmit"
                                                @click="onSubmitThrottle"
                                            >
                                                {{ $t('common.create') }}
                                            </UButton>
                                        </div>
                                    </div>
                                </UForm>
                            </div>
                        </div>
                    </template>
                    <div class="-my-2 mx-0 overflow-x-auto">
                        <div class="inline-block min-w-full px-1 py-2 align-middle">
                            <DataErrorBlock
                                v-if="error"
                                :title="$t('common.unable_to_load', [$t('common.role', 2)])"
                                :retry="refresh"
                            />
                            <UTable v-else :columns="columns" :rows="sortedRoles" :loading="loading">
                                <template #rank-data="{ row: role }">
                                    <div class="text-gray-900 dark:text-white">
                                        {{ role.jobLabel }} - {{ role.jobGradeLabel }} ({{ role.grade }})
                                    </div>
                                </template>
                                <template #actions-data="{ row: role }">
                                    <div class="text-right">
                                        <UButton variant="link" icon="i-mdi-eye" @click="selectedRole = role" />
                                    </div>
                                </template>
                            </UTable>

                            <SingleHint class="mt-2" hint-id="rector_roles_list" />

                            <SingleHint class="mt-2" hint-id="rector_roles_superuser" />
                        </div>
                    </div>
                </div>

                <div class="ml-2 flex w-full basis-2/3">
                    <template v-if="selectedRole">
                        <RoleView
                            :role-id="selectedRole.id"
                            @deleted="
                                selectedRole = undefined;
                                refresh();
                            "
                        />
                    </template>
                    <template v-else>
                        <DataNoDataBlock icon="i-mdi-select" :message="$t('common.none_selected', [$t('common.role')])" />
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>
