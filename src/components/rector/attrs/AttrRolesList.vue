<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { CheckIcon } from 'mdi-vue3';
import AttrView from '~/components/rector/attrs/AttrView.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { Role } from '~~/gen/ts/resources/permissions/permissions';
import { Job } from '~~/gen/ts/resources/users/jobs';
import GenericTable from '~/components/partials/elements/GenericTable.vue';

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const { data: roles, pending, refresh, error } = useLazyAsyncData('rector-roles', () => getRoles());

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

const selectedJob = ref<Job | null>(null);
const queryJobRaw = ref('');
const queryJob = computed(() => queryJobRaw.value.toLowerCase());
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

const selectedRole = ref<Role | undefined>();

const sortedRoles = computed(() => roles.value?.sort((a, b) => (a.jobLabel ?? '').localeCompare(b.jobLabel ?? '')));

onBeforeMount(async () => await listJobs());
</script>

<template>
    <div class="py-2">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="flex flex-col lg:flex-row">
                <div class="mt-2 flow-root basis-1/3">
                    <div v-if="can('RectorService.CreateRole')" class="sm:flex sm:items-center">
                        <div class="sm:flex-auto">
                            <UForm :state="{}">
                                <div class="mx-auto flex flex-row gap-4">
                                    <div class="flex-1">
                                        <label for="job" class="block text-sm font-medium leading-6">
                                            {{ $t('common.job') }}
                                        </label>
                                        <Combobox
                                            v-model="selectedJob"
                                            as="div"
                                            class="relative mt-2 flex w-full items-center"
                                            nullable
                                        >
                                            <div class="relative w-full">
                                                <ComboboxButton as="div" class="w-full">
                                                    <ComboboxInput
                                                        autocomplete="off"
                                                        class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :display-value="
                                                            (job: any) => (job ? `${job?.label} (${job?.name})` : '')
                                                        "
                                                        @change="queryJobRaw = $event.target.value"
                                                        @focusin="focusTablet(true)"
                                                        @focusout="focusTablet(false)"
                                                    />
                                                </ComboboxButton>

                                                <ComboboxOptions
                                                    class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                                >
                                                    <ComboboxOption
                                                        v-for="job in availableJobs.filter((g) =>
                                                            g.label.toLowerCase().includes(queryJob),
                                                        )"
                                                        :key="job.name"
                                                        v-slot="{ active, selected }"
                                                        :value="job"
                                                    >
                                                        <li
                                                            :class="[
                                                                'relative cursor-default select-none py-2 pl-8 pr-4',
                                                                active ? 'bg-primary-500' : '',
                                                            ]"
                                                        >
                                                            <span :class="['block truncate', selected && 'font-semibold']">
                                                                {{ job.label }} ({{ job.name }})
                                                            </span>

                                                            <span
                                                                v-if="selected"
                                                                :class="[
                                                                    active ? 'text-neutral' : 'text-primary-500',
                                                                    'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                                ]"
                                                            >
                                                                <CheckIcon class="size-5" />
                                                            </span>
                                                        </li>
                                                    </ComboboxOption>
                                                </ComboboxOptions>
                                            </div>
                                        </Combobox>
                                    </div>
                                    <div class="flex flex-initial flex-col justify-end">
                                        <UButton :disabled="selectedJob === null" @click="createRole()">
                                            {{ $t('common.create') }}
                                        </UButton>
                                    </div>
                                </div>
                            </UForm>
                        </div>
                    </div>
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
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold sm:pr-0"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </template>
                                <template #tbody>
                                    <tr v-for="role in sortedRoles" :key="role.id" class="even:bg-base-800">
                                        <td
                                            class="whitespace-nowrap py-2 pl-2 pr-3 text-sm font-medium"
                                            :title="`ID: ${role.id}`"
                                        >
                                            {{ role.jobLabel }} ({{ role.job }})
                                        </td>
                                        <td class="whitespace-nowrap py-2 pl-3 pr-2 text-right text-sm font-medium">
                                            <div class="flex flex-row justify-end">
                                                <UButton variant="link" icon="i-mdi-eye" @click="selectedRole = role" />
                                            </div>
                                        </td>
                                    </tr>
                                </template>
                            </GenericTable>
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
