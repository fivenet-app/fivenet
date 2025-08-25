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
</script>

<template>
    <USlideover :overlay="false">
        <UForm class="flex flex-1" :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard
                class="flex flex-1 flex-col"
                :ui="{
                    body: {
                        base: 'flex-1 h-full max-h-[calc(100dvh-(2*var(--ui-header-height)))] overflow-y-auto',
                        padding: 'px-1 py-2 sm:p-2',
                    },
                }"
            >
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl leading-6 font-semibold">
                            {{ $t('components.centrum.create_dispatch.title') }}
                        </h3>

                        <UButton
                            class="-my-1"
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-window-close"
                            @click="$emit('close')"
                        />
                    </div>
                </template>

                <div>
                    <dl class="divide-neutral/10 divide-y">
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                <label class="block text-sm leading-6 font-medium" for="message">
                                    {{ $t('common.message') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormField name="message">
                                    <UInput
                                        v-model="state.message"
                                        type="text"
                                        name="message"
                                        :placeholder="$t('common.message')"
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
                                <UFormField name="description">
                                    <UInput
                                        v-model="state.description"
                                        type="text"
                                        name="description"
                                        :placeholder="$t('common.description')"
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
                                <UFormField name="jobs.jobs">
                                    <USelectMenu
                                        v-model="state.jobs.jobs"
                                        name="jobs.jobs"
                                        multiple
                                        :placeholder="$t('common.job')"
                                        option-attribute="label"
                                        :search-attributes="['name', 'label']"
                                        value-key="name"
                                        :items="dispatchTargetJobs"
                                        searchable
                                        searchable-lazy
                                        :searchable-placeholder="$t('common.search_field')"
                                        :disabled="dispatchTargetJobs.length <= 1"
                                    >
                                        <template #item-label="{ item }">
                                            <span class="truncate">{{
                                                item.length > 0
                                                    ? item.map((j: { label: string }) => j.label).join(', ')
                                                    : '&nbsp;'
                                            }}</span>
                                        </template>

                                        <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                                    </USelectMenu>
                                </UFormField>
                            </dd>
                        </div>
                    </dl>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.create') }}
                        </UButton>

                        <UButton class="flex-1" block color="neutral" @click="$emit('close')">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </USlideover>
</template>
