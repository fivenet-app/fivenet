<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
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

const overlay = useOverlay();

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

const columns = computed(
    () =>
        [
            {
                accessorKey: 'rank',
                header: t('common.rank'),
                meta: {
                    class: {
                        td: 'text-highlighted',
                    },
                },
                cell: ({ row }) => `${row.original.jobLabel} - ${row.original.jobGradeLabel} (${row.original.grade})`,
            },
            {
                id: 'actions',
                cell: ({ row }) =>
                    h(UTooltip, { text: $t('common.show') }, [
                        h(UButton, {
                            to: { name: 'settings-roles-id', params: { id: row.original.id } },
                            variant: 'link',
                            icon: 'i-mdi-eye',
                        }),
                    ]),
            },
        ] as TableColumn<Role>[],
);

const route = useRoute('settings-roles-id');

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await createRole().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');

const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('pages.settings.roles.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/settings" />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <div class="grid grid-cols-1 gap-2 lg:grid-cols-3">
                <div class="mb-2">
                    <UForm
                        v-if="can('settings.SettingsService/CreateRole').value"
                        ref="formRef"
                        class="flex flex-row gap-2"
                        :schema="schema"
                        :state="state"
                        @submit="onSubmitThrottle"
                    >
                        <UFormField class="flex-1" name="grade" :label="$t('common.job_grade')">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.jobGrade"
                                    class="w-full"
                                    :items="availableJobGrades"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :disabled="availableJobGrades.length === 0"
                                >
                                    <template v-if="state.jobGrade" #default>
                                        <span class="truncate">{{ state.jobGrade?.label }} ({{ state.jobGrade?.grade }})</span>
                                    </template>

                                    <template #item="{ item }">
                                        <span class="truncate">{{ item.label }} ({{ item.grade }})</span>
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
                                    confirmModal.open({
                                        title: $t('components.hints.settings_roles_list.title'),
                                        description: $t('components.hints.settings_roles_list.content'),
                                        icon: 'i-mdi-information-outline',
                                        color: 'warning',
                                        iconClass: 'text-amber-500 dark:text-amber-400',
                                        confirm: async () => await formRef?.submit(),
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
                            :pagination-options="{ manualPagination: true }"
                            :sorting-options="{ manualSorting: true }"
                            :empty="$t('common.not_found', [$t('common.role', 2)])"
                            sticky
                        />

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
            </div>
        </template>
    </UDashboardPanel>
</template>
