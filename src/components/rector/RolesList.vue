<script lang="ts" setup>
import { Role } from '@fivenet/gen/resources/permissions/permissions_pb';
import { RpcError } from 'grpc-web';
import { CreateRoleRequest, GetRolesRequest } from '@fivenet/gen/services/rector/rector_pb';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline';
import RolesListEntry from './RolesListEntry.vue';
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { JobGrade } from '@fivenet/gen/resources/jobs/jobs_pb';
import { CompleteJobNamesRequest } from '@fivenet/gen/services/completor/completor_pb';
import { CheckIcon } from '@heroicons/vue/20/solid';
import { useAuthStore } from '~/store/auth';
import { watchDebounced } from '@vueuse/core';

const { $grpc } = useNuxtApp();

const store = useAuthStore();

const activeChar = computed(() => store.activeChar);

const { data: roles, pending, refresh, error } = await useLazyAsyncData('rector-roles', () => getRoles());

async function getRoles(): Promise<Array<Role>> {
    return new Promise(async (res, rej) => {
        const req = new GetRolesRequest();

        try {
            const resp = await $grpc.getRectorClient().
                getRoles(req, null);

            return res(resp.getRolesList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

let entriesJobGrades = [] as JobGrade[];
const filteredJobGrades = ref<JobGrade[]>([]);
const queryJobGrade = ref('');
const selectedJobGrade = ref<JobGrade>();

async function findJobGrades(): Promise<void> {

    const req = new CompleteJobNamesRequest();
    req.setExactMatch(true);
    req.setCurrentJob(true);

    try {
        const resp = await $grpc.getCompletorClient().
            completeJobNames(req, null);


        entriesJobGrades = resp.getJobsList()[0].getGradesList();
        filteredJobGrades.value = entriesJobGrades;
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
    }
}

async function createRole(): Promise<void> {
    return new Promise(async (res, rej) => {
        if (!selectedJobGrade.value) {
            return;
        }

        const req = new CreateRoleRequest();
        req.setJob(activeChar.value?.getJob()!);
        req.setGrade(selectedJobGrade.value.getGrade());

        try {
            const role = await $grpc.getRectorClient().
                createRole(req, null);

            if (role.hasRole()) {
                roles.value?.unshift(role.getRole()!);
            }
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watchDebounced(queryJobGrade, async () => { filteredJobGrades.value = entriesJobGrades.filter(g => g.getLabel().toLowerCase().includes(queryJobGrade.value.toLowerCase())) }, { debounce: 750, maxWait: 2000 });

onMounted(async () => {
    await findJobGrades();
});
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="flow-root mt-2">
                <div class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <form @submit.prevent="createRole()">
                            <div class="flex flex-row gap-4 mx-auto">
                                <div class="flex-1 form-control">
                                    <label for="grade" class="block text-sm font-medium leading-6 text-neutral">
                                        Job Grade
                                    </label>
                                    <div class="relative flex items-center mt-2">
                                        <Combobox as="div" v-model="selectedJobGrade" nullable>
                                            <div class="relative">
                                                <ComboboxButton as="div">
                                                    <ComboboxInput
                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        @change="queryJobGrade = $event.target.value"
                                                        :display-value="(grade: any) => grade?.getLabel()" />
                                                </ComboboxButton>

                                                <ComboboxOptions v-if="filteredJobGrades.length > 0"
                                                    class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                                                    <ComboboxOption v-for="grade in filteredJobGrades"
                                                        :key="grade.getGrade()" :value="grade" as="grade"
                                                        v-slot="{ active, selected }">
                                                        <li
                                                            :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                                                            <span :class="['block truncate', selected && 'font-semibold']">
                                                                {{ grade.getLabel() }}
                                                            </span>

                                                            <span v-if="selected"
                                                                :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                                                                <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                                            </span>
                                                        </li>
                                                    </ComboboxOption>
                                                </ComboboxOptions>
                                            </div>
                                        </Combobox>
                                    </div>
                                </div>
                                <div class="flex-1 form-control" v-can="'RectorService.CreateRole'">
                                    <button @click="createRole()"
                                        :disabled="selectedJobGrade && selectedJobGrade.getGrade() <= 0"
                                        class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500">
                                        Create
                                    </button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <DataPendingBlock v-if="pending" message="Loading roles..." />
                        <DataErrorBlock v-else-if="error" title="Unable to load roles!" :retry="refresh" />
                        <button v-else-if="roles && roles.length == 0" type="button"
                            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2">
                            <MagnifyingGlassIcon class="w-12 h-12 mx-auto text-neutral" />
                            <span class="block mt-2 text-sm font-semibold text-gray-300">
                                Use the search field above to search or update your query.
                            </span>
                        </button>
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Name
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Updated At
                                        </th>
                                        <th scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Actions
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <RolesListEntry v-for="role in roles" :role="role" />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Name
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            Updated At
                                        </th>
                                        <th scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                            Actions
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
