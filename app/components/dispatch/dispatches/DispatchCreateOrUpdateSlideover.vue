<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { useLivemapStore } from '~/stores/livemap';
import { getCentrumDispatchesClient } from '~~/gen/ts/clients';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';
import type { Coords } from '~~/gen/ts/resources/livemap/coords';

const props = defineProps<{
    location?: Coords;
    dispatch?: Dispatch;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { activeChar } = useAuth();

const livemapStore = useLivemapStore();
const { location: storeLocation, showLocationMarker, markerCoordPickerActive } = storeToRefs(livemapStore);

const centrumDispatchesClient = await getCentrumDispatchesClient();

const { data: dispatchTargetJobs } = useAuthedLazyAsyncData(
    'userState',
    'centrum-dispatches-target-jobs',
    async ({ signal }) => {
        try {
            const call = centrumDispatchesClient.listDispatchTargetJobs({}, { abort: signal });
            const { response } = await call;

            return response.jobs ?? [];
        } catch (e) {
            handleGRPCError(e as RpcError);
            throw e;
        }
    },
);

const sortedDispatchTargetJobs = computed(() => {
    const activeJob = activeChar.value?.job;
    return [...(dispatchTargetJobs.value ?? [])].sort((a, b) => Number(b.name === activeJob) - Number(a.name === activeJob));
});

const schema = z.object({
    message: z.coerce.string().min(3).max(255),
    description: z.union([z.string().min(3).max(512), z.string().length(0).optional()]),
    anon: z.coerce.boolean(),
    x: z.coerce.number(),
    y: z.coerce.number(),
    jobs: z.object({
        jobs: z.coerce.string().min(1).max(32).array().min(0).max(5).default([]),
    }),
});

type Schema = z.output<typeof schema>;

const initialX = props.location?.x ?? storeLocation.value?.x ?? 0;
const initialY = props.location?.y ?? storeLocation.value?.y ?? 0;

const state = reactive<Schema>({
    message: '',
    description: '',
    anon: false,
    x: initialX,
    y: initialY,
    jobs: {
        jobs: [],
    },
});

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(state, {
    serializer: (value) =>
        JSON.stringify({
            message: value.message,
            description: value.description,
            anon: value.anon,
            x: value.x,
            y: value.y,
            jobs: [...value.jobs.jobs].sort(),
        }),
});

const isPickingCoordinates = computed<boolean>(() => markerCoordPickerActive.value === true);

function stopCoordinatePicking(): void {
    markerCoordPickerActive.value = false;
    showLocationMarker.value = false;
}

watch(
    () => storeLocation.value,
    (location) => {
        if (!isPickingCoordinates.value || !location) return;

        state.x = location.x;
        state.y = location.y;
    },
);

function toggleCoordinatePicker(): void {
    if (markerCoordPickerActive.value) {
        stopCoordinatePicking();
        return;
    }

    markerCoordPickerActive.value = true;
    storeLocation.value = { x: state.x, y: state.y };
    showLocationMarker.value = true;
}

async function createDispatch(values: Schema): Promise<void> {
    try {
        stopCoordinatePicking();

        const call = centrumDispatchesClient.createDispatch({
            dispatch: {
                id: 0,
                job: '',
                jobs: {
                    jobs: values.jobs.jobs.map((name) => ({ name })),
                },
                message: values.message,
                description: values.description,
                anon: values.anon,
                attributes: {
                    list: [],
                },
                x: values.x,
                y: values.y,
                units: [],
            },
        });
        await call;

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

watch(sortedDispatchTargetJobs, (jobs) => {
    if (!jobs || jobs?.length <= 0) {
        state.jobs.jobs = [];
        syncSnapshot();
        return;
    }

    state.jobs.jobs = [jobs[0]!.name];
    syncSnapshot();
});

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    await createDispatch(event.data);
});

const formRef = useTemplateRef('formRef');

async function closeSlideover(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}

const initialDismissGuard = ref<boolean>(true);
const { start: startInitialDismissGuardTimeout, stop: stopInitialDismissGuardTimeout } = useTimeoutFn(
    () => {
        initialDismissGuard.value = false;
    },
    250,
    { immediate: false },
);

onMounted(() => {
    startInitialDismissGuardTimeout();
});

onBeforeUnmount(() => {
    stopInitialDismissGuardTimeout();
    stopCoordinatePicking();
});
</script>

<template>
    <USlideover
        :title="$t('components.dispatch.create_dispatch.title')"
        :close="false"
        :modal="false"
        :dismissible="!initialDismissGuard && !isPickingCoordinates && !hasUnsavedChanges && canSubmit"
        :overlay="false"
        :ui="{ content: isPickingCoordinates ? 'max-w-sm' : 'max-w-xl' }"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.dispatch.create_dispatch.title') }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeSlideover"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="submit">
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
                                    class="w-full"
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
                            <UFormField class="flex-1" name="description">
                                <UTextarea
                                    v-model="state.description"
                                    class="w-full"
                                    type="text"
                                    name="description"
                                    :placeholder="$t('common.description')"
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
                                    class="w-full"
                                    name="jobs.jobs"
                                    multiple
                                    :placeholder="$t('common.job')"
                                    :filter-fields="['name', 'label']"
                                    value-key="name"
                                    label-key="label"
                                    :items="sortedDispatchTargetJobs"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    :disabled="dispatchTargetJobs.length <= 1"
                                >
                                    <template #empty> {{ $t('common.not_found', [$t('common.job', 2)]) }} </template>
                                </USelectMenu>
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="x">
                                {{ $t('common.longitude') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="x">
                                <UInputNumber
                                    v-model="state.x"
                                    class="w-full"
                                    name="x"
                                    :step="0.00001"
                                    :placeholder="$t('common.longitude')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="y">
                                {{ $t('common.latitude') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="y">
                                <UInputNumber
                                    v-model="state.y"
                                    class="w-full"
                                    name="y"
                                    :step="0.00001"
                                    :placeholder="$t('common.latitude')"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium" />
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UButton
                                :variant="isPickingCoordinates ? 'solid' : 'soft'"
                                :color="isPickingCoordinates ? 'warning' : 'neutral'"
                                icon="i-mdi-crosshairs-gps"
                                :label="isPickingCoordinates ? $t('common.apply') : $t('common.select')"
                                type="button"
                                @click="toggleCoordinatePicker"
                            />
                        </dd>
                    </div>
                </dl>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('common.create')"
                    @click="() => formRef?.submit()"
                />

                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="closeSlideover" />
            </UFieldGroup>
        </template>
    </USlideover>
</template>
