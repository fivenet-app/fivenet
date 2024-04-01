<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
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

        notifications.add({
            title: { key: 'notifications.superuser_menu.setsuperusermode.title', parameters: {} },
            description: {
                key: 'notifications.superuser_menu.setsuperusermode.content',
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
        <div as="div" class="flex items-center">
            <UToggle v-model="superuser">
                <span class="sr-only">
                    {{ $t('common.superuser') }}
                </span>
            </UToggle>
            <span class="font-medium text-gray-300">{{ $t('common.superuser') }}</span>
        </div>

        <Combobox v-if="isSuperuser" v-model="selectedJob" as="div" nullable>
            <div class="relative mt-1">
                <ComboboxButton as="div">
                    <ComboboxInput
                        autocomplete="off"
                        class="hidden w-full rounded-md border-0 bg-base-700 py-1.5 text-xs text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 md:block"
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
                    class="absolute z-40 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-xs"
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
                                <CheckIcon class="size-5" aria-hidden="true" />
                            </span>
                        </li>
                    </ComboboxOption>
                </ComboboxOptions>
            </div>
        </Combobox>
    </div>
</template>
