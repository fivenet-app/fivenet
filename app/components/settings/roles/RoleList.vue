<script lang="ts" setup>
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import type { Job, JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Role } from '~~/gen/ts/resources/permissions/permissions';

const { t } = useI18n();

const modal = useOverlay();

const { can, activeChar } = useAuth();

const notifications = useNotificationsStore();

const completorStore = useCompletorStore();
const { getJobByName } = completorStore;

const settingsSettingsClient = await getSettingsSettingsClient();

const { data: roles, status, refresh, error } = useLazyAsyncData('settings-roles', () => getRoles());

async function getRoles(): Promise<Role[]> {
    try {
        const call = settingsSettingsClient.getRoles({});
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
        const call = settingsSettingsClient.createRole({
            job: activeChar.value!.job,
            grade: state.jobGrade.grade,
        });
        const { response } = await call;

        if (!response.role) {
            return;
        }

        notifications.add({
            title: { key: 'notifications.settings.role_created.title', parameters: {} },
            description: { key: 'notifications.settings.role_created.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        roles.value?.push(response.role!);

        await navigateTo({ name: 'settings-roles-id', params: { id: response.role.id } });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const sortedRoles = computed(() => [...(roles.value ?? [])].sort((a, b) => a.grade - b.grade));

const columns = [
    {
        accessorKey: 'rank',
        label: t('common.rank'),
    },
    {
        accessorKey: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
];

const route = useRoute('settings-roles-id');

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await createRole().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UDashboardPanelContent class="grid grid-cols-1 gap-2 lg:grid-cols-3">
        <div class="mb-2">
            <UForm
                v-if="can('settings.SettingsService/CreateRole').value"
                class="flex flex-row gap-2"
                :schema="schema"
                :state="state"
                @submit="refresh()"
            >
                <UFormField class="flex-1" name="grade" :label="$t('common.job_grade')">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.jobGrade"
                            :items="availableJobGrades"
                            by="grade"
                            searchable
                            :searchable-placeholder="$t('common.search_field')"
                        >
                            <template #item-label>
                                <template v-if="state.jobGrade">
                                    <span class="truncate">{{ state.jobGrade?.label }} ({{ state.jobGrade?.grade }})</span>
                                </template>
                            </template>

                            <template #option="{ option: jobGrade }">
                                <span class="truncate">{{ jobGrade.label }} ({{ jobGrade.grade }})</span>
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField name="submit" label="&nbsp;">
                    <UButton
                        class="flex-initial justify-end"
                        :disabled="state.jobGrade === undefined || state.jobGrade!.grade < 0 || !canSubmit"
                        :loading="!canSubmit"
                        icon="i-mdi-plus"
                        @click="
                            modal.open(ConfirmModal, {
                                title: $t('components.hints.settings_roles_list.title'),
                                description: $t('components.hints.settings_roles_list.content'),
                                icon: 'i-mdi-information-outline',
                                color: 'warning',
                                iconClass: 'text-amber-500 dark:text-amber-400',
                                confirm: onSubmitThrottle,
                            })
                        "
                    >
                        {{ $t('common.create') }}
                    </UButton>
                </UFormField>
            </UForm>

            <div>
                <SingleHint class="my-2" hint-id="settings_roles_list" />

                <DataErrorBlock
                    v-if="error"
                    :title="$t('common.unable_to_load', [$t('common.role', 2)])"
                    :error="error"
                    :retry="refresh"
                />
                <UTable
                    v-else
                    :columns="columns"
                    :data="sortedRoles"
                    :loading="isRequestPending(status)"
                    :empty-state="{
                        icon: 'i-mdi-account-group',
                        label: $t('common.not_found', [$t('common.role', 2)]),
                    }"
                >
                    <template #rank-cell="{ row: role }">
                        <div class="text-highlighted">{{ role.jobLabel }} - {{ role.jobGradeLabel }} ({{ role.grade }})</div>
                    </template>

                    <template #actions-cell="{ row: role }">
                        <UTooltip :text="$t('common.show')">
                            <UButton
                                :to="{ name: 'settings-roles-id', params: { id: role.id } }"
                                variant="link"
                                icon="i-mdi-eye"
                            />
                        </UTooltip>
                    </template>
                </UTable>

                <Pagination :status="status" :refresh="refresh" hide-buttons hide-text />

                <SingleHint class="mt-2" hint-id="settings_roles_superuser" />
            </div>
        </div>

        <div class="col-span-2 mb-2 w-full">
            <DataNoDataBlock
                v-if="!route.params.id"
                icon="i-mdi-select"
                :message="$t('common.none_selected', [$t('common.role')])"
            />
            <NuxtPage v-else @deleted="refresh()" />
        </div>
    </UDashboardPanelContent>
</template>
