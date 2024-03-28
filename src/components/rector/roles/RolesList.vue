<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { watchOnce } from '@vueuse/core';
import { CheckIcon, SelectIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { Role } from '~~/gen/ts/resources/permissions/permissions';
import { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import RoleView from '~/components/rector/roles/RoleView.vue';
import RolesListEntry from '~/components/rector/roles/RolesListEntry.vue';
import GenericTable from '~/components/partials/elements/GenericTable.vue';

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const completorStore = useCompletorStore();
const { getJobByName } = completorStore;

const { data: roles, pending, refresh, error } = useLazyAsyncData('rector-roles', () => getRoles());

async function getRoles(): Promise<Role[]> {
    try {
        const call = $grpc.getRectorClient().getRoles({});
        const { response } = await call;

        return response.roles;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const job = ref<Job | undefined>();
watchOnce(roles, async () => (job.value = await getJobByName(activeChar.value!.job)));

const selectedJobGrade = ref<JobGrade | null>(null);
const queryJobGradeRaw = ref('');
const queryJobGrade = computed(() => queryJobGradeRaw.value.toLowerCase());
const availableJobGrades = computed(
    () => job.value?.grades.filter((g) => (roles.value?.findIndex((r) => r.grade === g.grade) ?? -1) === -1) ?? [],
);

async function createRole(): Promise<void> {
    if (selectedJobGrade.value === null || selectedJobGrade.value.grade <= 0) {
        return;
    }

    try {
        const call = $grpc.getRectorClient().createRole({
            job: activeChar.value!.job,
            grade: selectedJobGrade.value.grade,
        });
        const { response } = await call;

        if (!response.role) {
            return;
        }

        notifications.dispatchNotification({
            title: { key: 'notifications.rector.role_created.title', parameters: {} },
            content: { key: 'notifications.rector.role_created.content', parameters: {} },
            type: 'success',
        });

        roles.value?.push(response.role!);

        selectedRole.value = response.role;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const sortedRoles = computed(() => roles.value?.sort((a, b) => a.grade - b.grade));

const selectedRole = ref<Role | undefined>();
</script>

<template>
    <div class="py-2">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="flex flex-col lg:flex-row">
                <div class="mt-2 flow-root basis-1/3">
                    <template v-if="can('RectorService.CreateRole')">
                        <div class="sm:flex sm:items-center">
                            <div class="sm:flex-auto">
                                <form @submit.prevent="createRole()">
                                    <div class="mx-auto flex flex-row gap-4">
                                        <div class="form-control flex-1">
                                            <label for="grade" class="block text-sm font-medium leading-6 text-neutral">
                                                {{ $t('common.job_grade') }}
                                            </label>
                                            <Combobox
                                                v-model="selectedJobGrade"
                                                as="div"
                                                class="relative mt-2 flex w-full items-center"
                                                nullable
                                            >
                                                <div class="relative w-full">
                                                    <ComboboxButton as="div" class="w-full">
                                                        <ComboboxInput
                                                            autocomplete="off"
                                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                            :display-value="
                                                                (grade: any) =>
                                                                    grade ? `${grade?.label} (${grade?.grade})` : ''
                                                            "
                                                            @change="queryJobGradeRaw = $event.target.value"
                                                            @focusin="focusTablet(true)"
                                                            @focusout="focusTablet(false)"
                                                        />
                                                    </ComboboxButton>

                                                    <ComboboxOptions
                                                        class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                                    >
                                                        <ComboboxOption
                                                            v-for="grade in availableJobGrades.filter((g) =>
                                                                g.label.toLowerCase().includes(queryJobGrade),
                                                            )"
                                                            v-slot="{ active, selected }"
                                                            :key="grade.grade"
                                                            :value="grade"
                                                        >
                                                            <li
                                                                :class="[
                                                                    'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                    active ? 'bg-primary-500' : '',
                                                                ]"
                                                            >
                                                                <span :class="['block truncate', selected && 'font-semibold']">
                                                                    {{ grade.label }} ({{ grade.grade }})
                                                                </span>

                                                                <span
                                                                    v-if="selected"
                                                                    :class="[
                                                                        active ? 'text-neutral' : 'text-primary-500',
                                                                        'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                                    ]"
                                                                >
                                                                    <CheckIcon class="h-5 w-5" aria-hidden="true" />
                                                                </span>
                                                            </li>
                                                        </ComboboxOption>
                                                    </ComboboxOptions>
                                                </div>
                                            </Combobox>
                                        </div>
                                        <div class="form-control flex flex-initial flex-col justify-end">
                                            <button
                                                type="submit"
                                                class="inline-flex rounded-md px-3 py-2 text-sm font-semibold text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                                :disabled="selectedJobGrade === null || selectedJobGrade!.grade <= 0"
                                                :class="[
                                                    selectedJobGrade === null || selectedJobGrade!.grade <= 0
                                                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                                ]"
                                            >
                                                {{ $t('common.create') }}
                                            </button>
                                        </div>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </template>
                    <div class="-my-2 mx-0 overflow-x-auto">
                        <div class="inline-block min-w-full px-1 py-2 align-middle">
                            <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.role', 2)])" />
                            <DataErrorBlock
                                v-else-if="error"
                                :title="$t('common.unable_to_load', [$t('common.role', 2)])"
                                :retry="refresh"
                            />
                            <DataNoDataBlock v-else-if="roles && roles.length === 0" :type="$t('common.role', 2)" />
                            <GenericTable v-else>
                                <template #thead>
                                    <tr>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold text-neutral sm:pr-0"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </template>
                                <template #tbody>
                                    <RolesListEntry
                                        v-for="role in sortedRoles"
                                        :key="role.id"
                                        :class="selectedRole?.id === role.id ? 'bg-base-800' : ''"
                                        :role="role"
                                        @selected="selectedRole = role"
                                    />
                                </template>
                            </GenericTable>

                            <SingleHint hint-id="rector_roles_list" />
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
                        <DataNoDataBlock :icon="SelectIcon" :message="$t('common.none_selected', [$t('common.role')])" />
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>
