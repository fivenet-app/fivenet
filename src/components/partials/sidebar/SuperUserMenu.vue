<script lang="ts" setup>
import {
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions,
    Switch,
    SwitchGroup,
    SwitchLabel,
} from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { CheckIcon } from 'mdi-vue3';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { Job } from '~~/gen/ts/resources/users/jobs';
import type { SetSuperUserModeRequest } from '~~/gen/ts/services/auth/auth';

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const authStore = useAuthStore();
const { activeChar, isSuperuser } = storeToRefs(authStore);
const { setAccessToken, setActiveChar, setJobProps } = authStore;

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

const superuser = ref(can('SuperUser'));

const queryJob = ref('');
const selectedJob = ref<undefined | Job>();

async function setSuperUserMode(): Promise<void> {
    try {
        const req = {
            superuser: superuser.value,
        } as SetSuperUserModeRequest;

        if (selectedJob.value) {
            req.job = selectedJob.value!.name;
        }

        const call = $grpc.getAuthClient().setSuperUserMode(req);
        const { response } = await call;

        if (superuser.value) {
            authStore.permissions.push('superuser');
        } else {
            authStore.permissions = authStore.permissions.filter((p) => p !== 'superuser');
        }

        setAccessToken(response.token, toDate(response.expires));
        setActiveChar(response.char!);
        setJobProps(response.jobProps);

        notifications.dispatchNotification({
            title: { key: 'notifications.superuser_menu.setsuperusermode.title', parameters: {} },
            content: {
                key: 'notifications.superuser_menu.setsuperusermode.title',
                parameters: { job: selectedJob.value?.label ?? activeChar.value?.jobLabel ?? 'N/A' },
            },
            type: 'info',
        });

        await navigateTo({ name: 'overview' });

        queryJob.value = '';
    } catch (e) {
        queryJob.value = '';
        superuser.value = !superuser.value;

        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const filteredJobs = computed(() => jobs.value.filter((g) => g.label.toLowerCase().includes(queryJob.value.toLowerCase())));

watch(superuser, () => setSuperUserMode());
watch(selectedJob, () => setSuperUserMode());
</script>

<template>
    <div class="flex flex-col items-center">
        <SwitchGroup as="div" class="flex items-center">
            <Switch
                v-model="superuser"
                :class="[
                    superuser ? 'bg-primary-600' : 'bg-gray-200',
                    'relative inline-flex h-4 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-600 focus:ring-offset-2',
                ]"
            >
                <span
                    aria-hidden="true"
                    :class="[
                        superuser ? 'translate-x-5' : 'translate-x-0',
                        'pointer-events-none inline-block h-3 w-3 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                    ]"
                />
            </Switch>
            <SwitchLabel as="span" class="ml-3 text-xs">
                <span class="font-medium text-gray-300">{{ $t('common.superuser') }}</span>
            </SwitchLabel>
        </SwitchGroup>

        <Combobox v-if="isSuperuser" v-model="selectedJob" as="div" nullable>
            <div class="mt-1 relative">
                <ComboboxButton as="div">
                    <ComboboxInput
                        autocomplete="off"
                        class="hidden md:block text-xs w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300"
                        :display-value="(job: any) => (job ? job?.label : '')"
                        :placeholder="`${$t('common.select')} ${$t('common.job')}`"
                        @click="listJobs"
                        @change="queryJob = $event.target.value"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </ComboboxButton>

                <ComboboxOptions
                    v-if="filteredJobs.length > 0"
                    class="absolute z-40 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                >
                    <ComboboxOption v-for="job in filteredJobs" :key="job.name" v-slot="{ active, selected }" :value="job">
                        <li
                            :class="[
                                'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                active ? 'bg-primary-500' : '',
                            ]"
                        >
                            <span :class="['block truncate', selected && 'font-semibold']">
                                {{ job.label }}
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
</template>
