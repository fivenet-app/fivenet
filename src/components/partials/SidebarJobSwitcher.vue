<script lang="ts" setup>
import { useAuthStore } from '../../store/auth';
import { computed, ref, onMounted } from 'vue';
import { Job } from '@arpanet/gen/resources/jobs/jobs_pb';
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { CheckIcon } from '@heroicons/vue/20/solid';
import { CompleteJobNamesRequest } from '@arpanet/gen/services/completor/completor_pb';
import { watchDebounced } from '@vueuse/shared';
import { SetJobRequest } from '@arpanet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';

const { $grpc } = useNuxtApp();
const store = useAuthStore();

const activeChar = computed(() => store.$state.activeChar);

let entriesJobs = [] as Job[];
const queryJob = ref('');
const selectedJob = ref<Job>();

async function findJobs(): Promise<void> {
    const req = new CompleteJobNamesRequest();
    req.setSearch(queryJob.value);

    try {
        const resp = await $grpc.getCompletorClient().
            completeJobNames(req, null)

        entriesJobs = resp.getJobsList();
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}

onMounted(async () => {
    findJobs();
});

async function setJob(): Promise<void> {
    const req = new SetJobRequest();
    req.setCharId(activeChar.value?.getUserId()!);
    req.setJob(selectedJob.value?.getName()!);
    const grades = selectedJob.value?.getGradesList()!;
    req.setJobGrade(grades[grades.length - 1].getGrade());

    try {
        const resp = await $grpc.getAuthClient().
            setJob(req, null);

        await store.updateAccessToken(resp.getToken());
        await store.updateActiveChar(resp.getChar()!);
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}

watchDebounced(selectedJob, () => setJob());
</script>

<template>
    <Combobox as="div" v-model="selectedJob">
        <div class="relative">
            <ComboboxButton as="div">
                <ComboboxInput
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    @change="queryJob = $event.target.value" :display-value="(job: any) => job?.getLabel()" />
            </ComboboxButton>

            <ComboboxOptions v-if="entriesJobs.length > 0"
                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                <ComboboxOption v-for="job in entriesJobs" :key="job.getName()" :value="job" as="job"
                    v-slot="{ active, selected }">
                    <li
                        :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                        <span :class="['block truncate', selected && 'font-semibold']">
                            {{ job.getLabel() }}
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
</template>
