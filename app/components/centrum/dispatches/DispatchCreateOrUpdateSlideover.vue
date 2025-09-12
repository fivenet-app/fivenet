<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { useLivemapStore } from '~/stores/livemap';
import type { Coordinate } from '~/types/livemap';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';

const props = defineProps<{
    location?: Coordinate;
    dispatch?: Dispatch;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { activeChar } = useAuth();

const livemapStore = useLivemapStore();
const { location: storeLocation } = storeToRefs(livemapStore);

const centrumCentrumClient = await getCentrumCentrumClient();

const { data: dispatchTargetJobs } = useLazyAsyncData('centrum-dispatches-target-jobs', async () => {
    try {
        const call = centrumCentrumClient.listDispatchTargetJobs({});
        const { response } = await call;

        return response.jobs ?? [];
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
});

const schema = z.object({
    message: z.string().min(3).max(255),
    description: z.union([z.string().min(3).max(512), z.string().length(0).optional()]),
    anon: z.coerce.boolean(),
    jobs: z.object({
        jobs: z.array(z.string().min(1).max(32)).min(1).max(5).default([]),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    message: '',
    description: '',
    anon: false,
    jobs: {
        jobs: [],
    },
});

async function createDispatch(values: Schema): Promise<void> {
    try {
        const call = centrumCentrumClient.createDispatch({
            dispatch: {
                id: 0,
                job: '',
                jobs: {
                    jobs: [],
                },
                message: values.message,
                description: values.description,
                anon: values.anon,
                attributes: {
                    list: [],
                },
                x: props.location?.x ?? storeLocation.value?.x ?? 0,
                y: props.location?.y ?? storeLocation.value?.y ?? 0,
                units: [],
            },
        });
        await call;

        emit('close');
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(dispatchTargetJobs, (jobs) => {
    if (!jobs || jobs?.length <= 0) {
        state.jobs.jobs = [];
        return;
    }

    state.jobs.jobs = [jobs[0]?.name ?? activeChar.value!.job];
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createDispatch(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <USlideover :title="$t('components.centrum.create_dispatch.title')" :overlay="false">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <dl class="divide-y divide-default">
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="message">
                                {{ $t('common.message') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField class="flex-1" name="message" required>
                                <UInput
                                    v-model="state.message"
                                    type="text"
                                    name="message"
                                    :placeholder="$t('common.message')"
                                    class="w-full"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="description">
                                {{ $t('common.description') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="description" class="flex-1">
                                <UTextarea
                                    v-model="state.description"
                                    type="text"
                                    name="description"
                                    :placeholder="$t('common.description')"
                                    class="w-full"
                                    :rows="3"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="anon">
                                {{ $t('common.anon') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="anon">
                                <UCheckbox v-model="state.anon" name="anon" :placeholder="$t('common.anon')" />
                            </UFormField>
                        </dd>
                    </div>

                    <div
                        v-if="dispatchTargetJobs && dispatchTargetJobs.length > 0"
                        class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                    >
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="jobs.jobs">
                                {{ $t('common.job') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="jobs.jobs" required>
                                <USelectMenu
                                    v-model="state.jobs.jobs"
                                    name="jobs.jobs"
                                    multiple
                                    :placeholder="$t('common.job')"
                                    :filter-fields="['name', 'label']"
                                    value-key="name"
                                    label-key="label"
                                    :items="dispatchTargetJobs"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :disabled="dispatchTargetJobs.length <= 1"
                                    class="w-full"
                                >
                                    <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                                </USelectMenu>
                            </UFormField>
                        </dd>
                    </div>
                </dl>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.create')"
                    @click="() => formRef?.submit()"
                />

                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close')" />
            </UButtonGroup>
        </template>
    </USlideover>
</template>
