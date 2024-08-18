<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useLivemapStore } from '~/store/livemap';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';

const props = defineProps<{
    location?: Coordinate;
    dispatch?: Dispatch;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { isOpen } = useSlideover();

const livemapStore = useLivemapStore();
const { location: storeLocation } = storeToRefs(livemapStore);

const schema = z.object({
    message: z.string().min(3).max(255),
    description: z.union([z.string().min(3).max(512), z.string().length(0).optional()]),
    anon: z.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    message: '',
    description: '',
    anon: false,
});

async function createDispatch(values: Schema): Promise<void> {
    try {
        const call = getGRPCCentrumClient().createDispatch({
            dispatch: {
                id: '0',
                job: '',
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
        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createDispatch(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <USlideover :ui="{ width: 'w-screen max-w-xl' }" :overlay="false">
        <UForm :schema="schema" :state="state" class="flex flex-1" @submit="onSubmitThrottle">
            <UCard
                class="flex flex-1 flex-col"
                :ui="{
                    body: {
                        base: 'flex-1 h-full max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
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
                    <dl class="divide-neutral/10 divide-y">
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="message" class="block text-sm font-medium leading-6">
                                    {{ $t('common.message') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="message">
                                    <UInput
                                        v-model="state.message"
                                        type="text"
                                        name="message"
                                        :placeholder="$t('common.message')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="description" class="block text-sm font-medium leading-6">
                                    {{ $t('common.description') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="description">
                                    <UInput
                                        v-model="state.description"
                                        type="text"
                                        name="description"
                                        :placeholder="$t('common.description')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </UFormGroup>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                <label for="anon" class="block text-sm font-medium leading-6">
                                    {{ $t('common.anon') }}
                                </label>
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UFormGroup name="anon">
                                    <UCheckbox v-model="state.anon" name="anon" :placeholder="$t('common.anon')" />
                                </UFormGroup>
                            </dd>
                        </div>
                    </dl>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.create') }}
                        </UButton>

                        <UButton
                            block
                            class="flex-1"
                            color="black"
                            @click="
                                $emit('close');
                                isOpen = false;
                            "
                        >
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </USlideover>
</template>
