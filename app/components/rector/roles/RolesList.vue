<script lang="ts" setup>
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Role } from '~~/gen/ts/resources/permissions/permissions';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const modal = useModal();

const notifications = useNotificatorStore();

const { can, activeChar } = useAuth();

const completorStore = useCompletorStore();
const { getJobByName } = completorStore;

const { data: roles, pending: loading, refresh, error } = useLazyAsyncData('rector-roles', () => getRoles());

async function getRoles(): Promise<Role[]> {
    try {
        const call = $grpc.rector.rector.getRoles({});
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
    if (state.jobGrade === undefined || state.jobGrade.grade < 0) {
        return;
    }

    try {
        const call = $grpc.rector.rector.createRole({
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

        await navigateTo({ name: 'rector-roles-id', params: { id: response.role.id } });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const sortedRoles = computed(() => [...(roles.value ?? [])].sort((a, b) => a.grade - b.grade));

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

const route = useRoute('rector-roles-id');

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await createRole().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UDashboardPanelContent class="flex flex-col gap-2 lg:flex-row">
        <div class="basis-1/3">
            <UForm
                v-if="can('RectorService.CreateRole').value"
                class="flex flex-row gap-2"
                :schema="schema"
                :state="state"
                @submit="refresh()"
            >
                <UFormGroup class="flex-1" name="grade" :label="$t('common.job_grade')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.jobGrade"
                            :options="availableJobGrades"
                            by="grade"
                            searchable
                            :searchable-placeholder="$t('common.search_field')"
                        >
                            <template #label>
                                <template v-if="state.jobGrade">
                                    <span class="truncate">{{ state.jobGrade?.label }} ({{ state.jobGrade?.grade }})</span>
                                </template>
                            </template>

                            <template #option="{ option: jobGrade }">
                                <span class="truncate">{{ jobGrade.label }} ({{ jobGrade.grade }})</span>
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormGroup>

                <UFormGroup name="submit" label="&nbsp;">
                    <UButton
                        class="flex-initial justify-end"
                        :disabled="state.jobGrade === undefined || state.jobGrade!.grade < 0 || !canSubmit"
                        :loading="!canSubmit"
                        icon="i-mdi-plus"
                        @click="
                            modal.open(ConfirmModal, {
                                title: $t('components.hints.rector_roles_list.title'),
                                description: $t('components.hints.rector_roles_list.content'),
                                icon: 'i-mdi-information-outline',
                                color: 'amber',
                                iconClass: 'text-amber-500 dark:text-amber-400',
                                confirm: onSubmitThrottle,
                            })
                        "
                    >
                        {{ $t('common.create') }}
                    </UButton>
                </UFormGroup>
            </UForm>

            <div>
                <SingleHint class="my-2" hint-id="rector_roles_list" />

                <DataErrorBlock
                    v-if="error"
                    :title="$t('common.unable_to_load', [$t('common.role', 2)])"
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
                    <template #rank-data="{ row: role }">
                        <div class="text-gray-900 dark:text-white">
                            {{ role.jobLabel }} - {{ role.jobGradeLabel }} ({{ role.grade }})
                        </div>
                    </template>

                    <template #actions-data="{ row: role }">
                        <div class="text-right">
                            <UTooltip :text="$t('common.show')">
                                <UButton
                                    :to="{ name: 'rector-roles-id', params: { id: role.id } }"
                                    variant="link"
                                    icon="i-mdi-eye"
                                />
                            </UTooltip>
                        </div>
                    </template>
                </UTable>

                <Pagination :loading="loading" :refresh="refresh" hide-buttons hide-text />

                <SingleHint class="mt-2" hint-id="rector_roles_superuser" />
            </div>
        </div>

        <div class="w-full basis-2/3">
            <DataNoDataBlock
                v-if="!route.params.id"
                icon="i-mdi-select"
                :message="$t('common.none_selected', [$t('common.role')])"
            />
            <NuxtPage v-else @deleted="refresh()" />
        </div>
    </UDashboardPanelContent>
</template>
