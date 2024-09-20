<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import type { TypedRouteFromName } from '@typed-router';
import { z } from 'zod';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import { useAuthStore } from '~/store/auth';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { SetJobsUserPropsResponse } from '~~/gen/ts/services/jobs/jobs';

useHead({
    title: 'pages.jobs.colleagues.single.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.single.title',
    requiresAuth: true,
    permission: 'JobsService.GetColleague',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-info'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return idParamRegex.test(route.params.id as string);
    },
});

const props = defineProps<{
    colleague: Colleague;
}>();

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const schema = z.object({
    note: z.string().min(0).max(512),
    reason: z.string().min(3).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    note: props.colleague?.props?.note ?? '',
    reason: '',
});
watch(props, () => {
    if (!props.colleague?.props) {
        return;
    }

    state.note = props.colleague.props.note ?? '';
});

async function setJobsUserNote(values: Schema): Promise<undefined | SetJobsUserPropsResponse> {
    if (!props.colleague) {
        return;
    }

    try {
        const call = getGRPCJobsClient().setJobsUserProps({
            reason: values.reason,
            props: {
                userId: props.colleague.userId,
                job: '',
                note: values.note,
            },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setJobsUserNote(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    editing.value = !editing.value;
}, 1000);

const canDo = computed(() => ({
    view: can('JobsService.GetColleague').value && attr('JobsService.GetColleague', 'Types', 'Note').value,
    edit:
        can('JobsService.SetJobsUserProps').value &&
        attr('JobsService.SetJobsUserProps', 'Types', 'Note').value &&
        checkIfCanAccessColleague(activeChar.value!, props.colleague, 'JobsService.SetJobsUserProps'),
}));

const editing = ref(false);

watch(editing, () => {
    if (!editing.value) {
        emits('refresh');
    }
});
</script>

<template>
    <UContainer class="w-full">
        <div class="w-full grow lg:flex lg:flex-col">
            <div class="flex-1 px-4 py-5 sm:p-0">
                <dl class="space-y-4 sm:space-y-0 xl:grid xl:grid-cols-2">
                    <div class="border-b border-gray-100 sm:flex sm:px-5 sm:py-4 dark:border-gray-800">
                        <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                            {{ $t('common.date_of_birth') }}
                        </dt>
                        <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                            {{ colleague.dateofbirth }}
                        </dd>
                    </div>

                    <div class="border-b border-gray-100 sm:flex sm:px-5 sm:py-4 dark:border-gray-800">
                        <dt class="text-sm font-medium sm:w-40 sm:shrink-0 lg:w-48">
                            {{ $t('common.phone_number') }}
                        </dt>
                        <dd class="mt-1 text-sm text-base-800 sm:col-span-2 sm:ml-6 sm:mt-0 dark:text-base-300">
                            <PhoneNumberBlock :number="colleague.phoneNumber" />
                        </dd>
                    </div>

                    <!-- Note -->
                    <div
                        v-if="colleague && canDo.view"
                        class="col-span-2 border-b border-gray-100 px-4 py-5 sm:flex dark:border-gray-800"
                    >
                        <UForm :schema="schema" :state="state" class="w-full" @submit="onSubmitThrottle">
                            <div class="flex items-center">
                                <h4 class="flex-1 text-sm font-semibold leading-6">
                                    {{ $t('common.note') }}
                                </h4>

                                <div v-if="canDo.edit" class="mb-1">
                                    <UButton
                                        v-if="!editing"
                                        icon="i-mdi-pencil"
                                        :loading="!canSubmit"
                                        @click="editing = !editing"
                                    />
                                    <template v-else>
                                        <UButtonGroup>
                                            <UButton type="submit" icon="i-mdi-content-save" :loading="!canSubmit" />
                                            <UButton
                                                color="red"
                                                icon="i-mdi-cancel"
                                                :loading="!canSubmit"
                                                @click="editing = !editing"
                                            />
                                        </UButtonGroup>
                                    </template>
                                </div>
                            </div>

                            <div class="flex flex-1 flex-col gap-2">
                                <template v-if="!editing">
                                    <div class="w-full flex-1">
                                        <p class="prose dark:prose-invert whitespace-pre-wrap text-base-800 dark:text-base-300">
                                            {{ colleague?.props?.note ?? $t('common.na') }}
                                        </p>
                                    </div>
                                </template>
                                <template v-else>
                                    <UFormGroup name="note" class="w-full">
                                        <UTextarea v-model="state.note" block :rows="6" :maxrows="10" name="note" />
                                    </UFormGroup>

                                    <UFormGroup name="reason" :label="$t('common.reason')" class="w-full" required>
                                        <UInput v-model="state.reason" type="text" name="reason" />
                                    </UFormGroup>
                                </template>
                            </div>
                        </UForm>
                    </div>
                </dl>
            </div>
        </div>
    </UContainer>
</template>
