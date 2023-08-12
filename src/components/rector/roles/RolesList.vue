<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced, watchOnce } from '@vueuse/core';
import { CheckIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import { JobGrade } from '~~/gen/ts/resources/jobs/jobs';
import { Role } from '~~/gen/ts/resources/permissions/permissions';
import RolesListEntry from './RolesListEntry.vue';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const notifications = useNotificationsStore();

const { activeChar } = storeToRefs(authStore);

const { data: roles, pending, refresh, error } = useLazyAsyncData('rector-roles', () => getRoles());

async function getRoles(): Promise<Role[]> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getRectorClient().getRoles({});
            const { response } = await call;

            return res(response.roles);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

let entriesJobGrades = [] as JobGrade[];
const filteredJobGrades = ref<JobGrade[]>([]);
const queryJobGrade = ref('');
const selectedJobGrade = ref<JobGrade>();

async function findJobGrades(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCompletorClient().completeJobs({
                currentJob: true,
                exactMatch: true,
            });
            const { response } = await call;

            entriesJobGrades = response.jobs[0].grades;
            filteredJobGrades.value = entriesJobGrades.filter(
                (g) => (roles.value?.findIndex((r) => r.grade === g.grade) ?? -1) === -1,
            );

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function createRole(): Promise<void> {
    return new Promise(async (res, rej) => {
        if (!selectedJobGrade.value) {
            return res();
        }

        try {
            const call = $grpc.getRectorClient().createRole({
                job: activeChar.value?.job!,
                grade: selectedJobGrade.value.grade,
            });
            const { response } = await call;

            if (!response.role) {
                return res();
            }

            roles.value?.unshift(response.role!);

            notifications.dispatchNotification({
                title: { key: 'notifications.rector.role_created.title', parameters: [] },
                content: { key: 'notifications.rector.role_created.content', parameters: [] },
                type: 'success',
            });

            await navigateTo({
                name: 'rector-roles-id',
                params: {
                    id: response.role?.id?.toString(),
                },
            });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watchDebounced(
    queryJobGrade,
    async () => {
        filteredJobGrades.value = entriesJobGrades.filter((g) =>
            g.label.toLowerCase().includes(queryJobGrade.value.toLowerCase()),
        );
    },
    { debounce: 600, maxWait: 1750 },
);

watchOnce(roles, async () => await findJobGrades());
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="flow-root mt-2">
                <div v-if="can('RectorService.CreateRole')" class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <form @submit.prevent="createRole()">
                            <div class="flex flex-row gap-4 mx-auto">
                                <div class="flex-1 form-control">
                                    <label for="grade" class="block text-sm font-medium leading-6 text-neutral">
                                        {{ $t('common.job_grade') }}
                                    </label>
                                    <Combobox
                                        as="div"
                                        v-model="selectedJobGrade"
                                        class="relative flex items-center mt-2 w-full"
                                        nullable
                                    >
                                        <div class="relative w-full">
                                            <ComboboxButton as="div" class="w-full">
                                                <ComboboxInput
                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    @change="queryJobGrade = $event.target.value"
                                                    :display-value="(grade: any) => grade?.label"
                                                />
                                            </ComboboxButton>

                                            <ComboboxOptions
                                                v-if="filteredJobGrades.length > 0"
                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                            >
                                                <ComboboxOption
                                                    v-for="grade in filteredJobGrades"
                                                    :key="grade.grade"
                                                    :value="grade"
                                                    as="grade"
                                                    v-slot="{ active, selected }"
                                                >
                                                    <li
                                                        :class="[
                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                            active ? 'bg-primary-500' : '',
                                                        ]"
                                                    >
                                                        <span :class="['block truncate', selected && 'font-semibold']">
                                                            {{ grade.label }}
                                                        </span>

                                                        <span
                                                            v-if="selected"
                                                            :class="[
                                                                active ? 'text-neutral' : 'text-primary-500',
                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                            ]"
                                                        >
                                                            <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                                        </span>
                                                    </li>
                                                </ComboboxOption>
                                            </ComboboxOptions>
                                        </div>
                                    </Combobox>
                                </div>
                                <div class="flex-initial form-control flex flex-col justify-end">
                                    <button
                                        type="submit"
                                        :disabled="selectedJobGrade && selectedJobGrade.grade <= 0"
                                        class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                    >
                                        Create
                                    </button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.role', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.role', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock v-else-if="roles && roles.length === 0" :type="$t('common.role', 2)" />
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <RolesListEntry v-for="role in roles" :role="role" />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
