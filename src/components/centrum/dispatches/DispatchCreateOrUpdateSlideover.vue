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
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 max-h-[calc(100vh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
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
                            <dl class="divide-neutral/10 border-neutral/10 divide-y border-b">
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="message" class="block text-sm font-medium leading-6">
                                            {{ $t('common.message') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        <VeeField
                                            type="text"
                                            name="message"
                                            :placeholder="$t('common.message')"
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
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        <VeeField
                                            type="text"
                                            name="description"
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
                                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                        <div class="flex h-6 items-center">
                                            <VeeField
                                                type="checkbox"
                                                name="anon"
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
                <div class="min-h-[var(--header-height)]">
                    <UButton
                        block
                        color="black"
                        @click="
                            $emit('close');
                            isOpen = false;
                        "
                    >
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        block
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        {{ $t('common.create') }}
                    </UButton>
                </div>
            </template>
        </UCard>
    </USlideover>
</template>
