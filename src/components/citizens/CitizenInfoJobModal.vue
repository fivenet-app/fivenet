<script lang="ts" setup>
import { Job } from '@fivenet/gen/resources/jobs/jobs_pb';
import { User, UserProps } from '@fivenet/gen/resources/users/users_pb';
import { SetUserPropsRequest } from '@fivenet/gen/services/citizenstore/citizenstore_pb';
import { CompleteJobsRequest } from '@fivenet/gen/services/completor/completor_pb';
import {
    Dialog,
    DialogPanel,
    TransitionChild,
    TransitionRoot,
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions
} from '@headlessui/vue';
import { CheckIcon } from '@heroicons/vue/24/solid';
import { watchDebounced } from '@vueuse/core';
import { RpcError } from 'grpc-web';
import { useNotificationsStore } from '~/store/notifications';

const { $grpc } = useNuxtApp();
const notifications = useNotificationsStore();

const { t } = useI18n();

const props = defineProps<{
    open: boolean,
    user: User,
}>();

defineEmits<{
    (e: 'close'): void,
}>();

const queryJob = ref<string>('');
const selectedJob = ref<undefined | Job>();

const reason = ref<string>('');

const { data: jobs } = useLazyAsyncData('jobs', () => getJobs());

async function getJobs(): Promise<Array<Job>> {
    return new Promise(async (res, rej) => {
        const req = new CompleteJobsRequest();
        req.setSearch(queryJob.value);

        try {
            const resp = await $grpc.getCompletorClient().
                completeJobs(req, null);

            return res(resp.getJobsList());
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watch(jobs, () => {
    if (jobs.value) {
        selectedJob.value = jobs.value.find(j => j.getName() === props.user.getJob());
    }
});

watchDebounced(queryJob, async () => await getJobs(), { debounce: 600, maxWait: 1750 });

async function setJobProp(): Promise<void> {
    return new Promise(async (res, rej) => {
        if (!selectedJob.value || selectedJob.value.getName() === props.user.getJob()) return res();

        if (reason.value.length < 3) return res();

        const userProps = new UserProps();
        userProps.setUserId(props.user.getUserId());
        userProps.setJobName(selectedJob.value.getName());

        const req = new SetUserPropsRequest();
        req.setProps(userProps);
        req.setReason(reason.value);

        try {
            await $grpc.getCitizenStoreClient().setUserProps(req, null);

            props.user.setJob(selectedJob.value?.getName()!);
            props.user.setJobLabel(selectedJob.value?.getLabel()!);

            notifications.dispatchNotification({
                title: t('notifications.action_successfull.title'),
                content: t('notifications.action_successfull.content'),
                type: 'success'
            });

            reason.value = '';
            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <TransitionChild as="template" enter="ease-out duration-300" enter-from="opacity-0" enter-to="opacity-100"
                leave="ease-in duration-200" leave-from="opacity-100" leave-to="opacity-0">
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
            </TransitionChild>

            <div class="fixed inset-0 z-10 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild as="template" enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100" leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
                        <DialogPanel
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-850 text-neutral sm:my-8 sm:w-full sm:max-w-2xl sm:p-6 h-96">
                            <div class="my-2 space-y-24">
                                <Combobox as="div" v-model="selectedJob" nullable>
                                    <div class="relative">
                                        <ComboboxButton as="div">
                                            <ComboboxInput
                                                class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                @change="queryJob = $event.target.value"
                                                :display-value="(job: any) => job.getLabel()" autocomplete="off" />
                                        </ComboboxButton>

                                        <ComboboxOptions v-if="jobs"
                                            class="absolute z-10 w-full py-1 mt-5 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                                            <ComboboxOption v-for="job in jobs" :key="job.getName()" :value="job" as="char"
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

                                <input type="text" name="reason" id="reason"
                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    :placeholder="$t('common.reason')" v-model="reason" />
                            </div>
                            <div class="absolute bottom-0 w-full left-0 sm:flex">
                                <button type="button"
                                    class="flex-1 rounded-bd bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="$emit('close')">{{ $t('common.close', 1) }}</button>
                                <button type="button"
                                    class="flex-1 rounded-bd bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                    @click="setJobProp(); $emit('close')">{{ $t('common.save') }}</button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
