<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { useLivemapStore } from '~/store/livemap';

const props = defineProps<{
    location?: Coordinate;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const livemapStore = useLivemapStore();
const { location: storeLocation } = storeToRefs(livemapStore);

interface FormData {
    message: string;
    description?: string;
    anon: boolean;
}

async function createDispatch(values: FormData): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().createDispatch({
            dispatch: {
                id: '0',
                job: '',
                message: values.message,
                description: values.description,
                anon: values.anon ?? false,
                attributes: {
                    list: [],
                },
                x: props.location ? props.location.x : storeLocation.value?.x ?? 0,
                y: props.location ? props.location.y : storeLocation.value?.y ?? 0,
                units: [],
            },
        });
        await call;

        resetForm();

        emit('close');
        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, resetForm } = useForm<FormData>({
    validationSchema: {
        message: { required: true, min: 3, max: 255 },
        description: { required: false, min: 6, max: 512 },
        anon: { required: false },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createDispatch(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <USlideover>
        <UCard
            class="flex flex-col flex-1"
            :ui="{ body: { base: 'flex-1' }, ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.centrum.create_dispatch.title') }}
                    </h3>

                    <UButton
                        color="gray"
                        variant="ghost"
                        icon="i-mdi-window-close"
                        class="-my-1"
                        @click="
                            $emit('close');
                            isOpen = false;
                        "
                    />
                </div>
            </template>

            <div>
                <div class="flex flex-1 flex-col justify-between">
                    <div class="divide-y divide-gray-200 px-2 sm:px-6">
                        <div class="mt-1">
                            <dl class="divide-y divide-neutral/10 border-b border-neutral/10">
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="message" class="block text-sm font-medium leading-6">
                                            {{ $t('common.message') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                        <VeeField
                                            type="text"
                                            name="message"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.message')"
                                            :label="$t('common.message')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="message" as="p" class="mt-2 text-sm text-error-400" />
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="description" class="block text-sm font-medium leading-6">
                                            {{ $t('common.description') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                        <VeeField
                                            type="text"
                                            name="description"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.description')"
                                            :label="$t('common.description')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="anon" class="block text-sm font-medium leading-6">
                                            {{ $t('common.anon') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                        <div class="flex h-6 items-center">
                                            <VeeField
                                                type="checkbox"
                                                name="anon"
                                                class="size-5 rounded border-gray-300 text-primary-600 focus:ring-primary-600"
                                                :placeholder="$t('common.anon')"
                                                :label="$t('common.anon')"
                                                :value="true"
                                            />
                                        </div>
                                        <VeeErrorMessage name="anon" as="p" class="mt-2 text-sm text-error-400" />
                                    </dd>
                                </div>
                            </dl>
                        </div>
                    </div>
                </div>
            </div>

            <template #footer>
                <UButton
                    class="inline-flex w-full items-center rounded-l-md px-3 py-2 text-sm font-semibold shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 sm:col-start-2"
                    :disabled="!meta.valid || !canSubmit"
                    :loading="!canSubmit"
                    @click="onSubmitThrottle"
                >
                    {{ $t('common.create') }}
                </UButton>
                <UButton
                    @click="
                        $emit('close');
                        isOpen = false;
                    "
                >
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </USlideover>
</template>
