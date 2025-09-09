<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import { getSettingsSettingsClient } from '~~/gen/ts/clients';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Role } from '~~/gen/ts/resources/permissions/permissions';

const { t } = useI18n();

const { can } = useAuth();

const notifications = useNotificationsStore();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const settingsSettingsClient = await getSettingsSettingsClient();

const { data: roles, status, refresh, error } = useLazyAsyncData('settings-roles', () => getRoles());

async function getRoles(): Promise<Role[]> {
    try {
        const call = settingsSettingsClient.getRoles({
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
        const call = settingsSettingsClient.createRole({
            job: state.job?.name,
            grade: 1,
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

        await navigateTo({ name: 'settings-limiter-job', params: { job: response.role.job } });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

onBeforeMount(async () => await listJobs());

const sortedRoles = computed(() => [...(roles.value ?? [])].sort((a, b) => (a.jobLabel ?? '').localeCompare(b.jobLabel ?? '')));

const appConfig = useAppConfig();

const columns = computed(
    () =>
        [
            {
                accessorKey: 'job',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.job'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                meta: {
                    class: {
                        td: 'text-highlighted',
                    },
                },
                cell: ({ row }) => `${row.original.jobLabel} (${row.original.job})`,
            },
            {
                id: 'actions',
                cell: ({ row }) =>
                    h(UTooltip, { text: t('common.show') }, () =>
                        h(UButton, {
                            to: { name: 'settings-limiter-job', params: { job: row.original.job } },
                            variant: 'link',
                            icon: 'i-mdi-eye',
                        }),
                    ),
            },
        ] as TableColumn<Role>[],
);

const route = useRoute('settings-limiter-job');

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async () => {
    canSubmit.value = false;
    await createRole().finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('pages.settings.limiter.title')">
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
                        class="flex flex-row gap-2"
                        :schema="schema"
                        :state="state"
                        @submit="refresh()"
                    >
                        <UFormField class="flex-1" name="grade" :label="$t('common.job')">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.job"
                                    class="w-full"
                                    :items="availableJobs"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :filter-fields="['label', 'name']"
                                >
                                    <template v-if="state.job" #default>
                                        {{ state.job?.label }} ({{ state.job.name }})
                                    </template>

                                    <template #item="{ item }"> {{ item.label }} ({{ item.name }}) </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormField>

                        <UFormField name="submit" label="&nbsp;">
                            <UButton
                                :disabled="state.job === undefined || !canSubmit"
                                :loading="!canSubmit"
                                icon="i-mdi-plus"
                                @click="onSubmitThrottle"
                            >
                                {{ $t('common.create') }}
                            </UButton>
                        </UFormField>
                    </UForm>

                    <div>
                        <DataErrorBlock
                            v-if="error"
                            :title="$t('common.unable_to_load', [$t('common.job', 2)])"
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
                    </div>
                </div>

                <div class="col-span-2 mb-2 w-full">
                    <DataNoDataBlock
                        v-if="!route.params.job"
                        icon="i-mdi-select"
                        :message="$t('common.none_selected', [$t('common.job')], 2)"
                    />
                    <NuxtPage v-else @deleted="refresh()" />
                </div>
            </div>
        </template>
    </UDashboardPanel>
</template>
