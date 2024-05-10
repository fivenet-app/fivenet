<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import type { TypedRouteFromName } from '@typed-router';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { SetJobsUserPropsResponse } from '~~/gen/ts/services/jobs/jobs';
import { checkIfCanAccessColleague } from '~/components/jobs/colleagues/helpers';
import { useAuthStore } from '~/store/auth';

useHead({
    title: 'pages.citizens.id.title',
});
definePageMeta({
    title: 'pages.citizens.id.title',
    requiresAuth: true,
    permission: 'JobsService.GetColleague',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-info'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const props = defineProps<{
    colleague: Colleague;
}>();

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const schema = z.object({
    note: z.string().min(0).max(512),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    note: props.colleague?.props?.note ?? '',
});

async function setJobsUserNote(values: Schema): Promise<void | SetJobsUserPropsResponse> {
    if (!props.colleague) {
        return;
    }

    try {
        const call = $grpc.getJobsClient().setJobsUserProps({
            reason: '',
            props: {
                userId: props.colleague.userId,
                job: '',
                note: values.note,
            },
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setJobsUserNote(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
    editing.value = !editing.value;
}, 1000);

const canEdit =
    can('JobsService.SetJobsUserProps') &&
    attr('JobsService.SetJobsUserProps', 'Types', 'Note') &&
    checkIfCanAccessColleague(activeChar.value!, props.colleague, 'JobsService.SetJobsUserProps');

const editing = ref(false);

watch(editing, () => {
    if (!editing.value) {
        emits('refresh');
    }
});
</script>

<template>
    <UContainer class="w-full">
        <!-- Note -->
        <UForm v-if="colleague" :schema="schema" :state="state" class="w-full flex-col" @submit="onSubmitThrottle">
            <div class="flex items-center">
                <h4 v-if="canEdit" class="flex-1 text-base font-semibold leading-6">{{ $t('common.note') }}:</h4>

                <template v-if="canEdit">
                    <UButton
                        v-if="!editing"
                        variant="link"
                        icon="i-mdi-pencil"
                        :loading="!canSubmit"
                        @click="editing = !editing"
                    />
                    <div v-else class="flex flex-row gap-1">
                        <UButton variant="link" icon="i-mdi-content-save" :loading="!canSubmit" @click="onSubmitThrottle" />
                        <UButton variant="link" icon="i-mdi-cancel" :loading="!canSubmit" @click="editing = !editing" />
                    </div>
                </template>
            </div>

            <div class="flex flex-1">
                <template v-if="!editing">
                    <div class="w-full flex-1">
                        <p class="prose prose-invert">
                            {{ colleague?.props?.note }}
                        </p>
                    </div>
                </template>
                <template v-else>
                    <UFormGroup name="note" class="w-full">
                        <UTextarea
                            v-model="state.note"
                            block
                            :rows="6"
                            :maxrows="10"
                            name="note"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>
                </template>
            </div>
        </UForm>
    </UContainer>
</template>
