<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { max, min, required, mimes, size } from '@vee-validate/rules';
import { CloseIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import SquareImg from '~/components/partials/elements/SquareImg.vue';
import { useNotificatorStore } from '~/store/notificator';
import type { File } from '~~/gen/ts/resources/filestore/file';
import { User, UserProps } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    open: boolean;
    user: User;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'update:mugShot', value?: File): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

interface FormData {
    reason: string;
    mugShot?: Blob;
    reset?: boolean;
}

async function setMugShot(values: FormData): Promise<void> {
    const userProps: UserProps = {
        userId: props.user.userId,
    };
    if (!values.reset) {
        if (!values.mugShot) {
            return;
        }

        userProps.mugShot = { data: new Uint8Array(await values.mugShot.arrayBuffer()) };
    } else {
        userProps.mugShot = { data: new Uint8Array(), delete: true };
    }

    try {
        const call = $grpc.getCitizenStoreClient().setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:mugShot', response.props?.mugShot);

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('size', size);
defineRule('mimes', mimes);

const { handleSubmit, meta, setFieldValue } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
        mugShot: { required: false, mimes: ['image/jpeg', 'image/jpg', 'image/png'], size: 2000 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> => await setMugShot(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
    setFieldValue('reset', false);
}, 1000);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 bg-base-900/75 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild
                        as="template"
                        enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100"
                        leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <DialogPanel
                            class="relative h-112 w-full overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left transition-all sm:my-8 sm:max-w-2xl sm:p-6"
                        >
                            <div class="absolute right-0 top-0 block pr-4 pt-4">
                                <UButton
                                    class="rounded-md bg-neutral-50 text-gray-400 hover:text-gray-500 focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="size-5" />
                                </UButton>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                {{ $t('components.citizens.citizen_info_profile.set_mug_shot') }}
                            </DialogTitle>
                            <form @submit.prevent="onSubmitThrottle">
                                <div class="my-2 space-y-24">
                                    <div class="flex-1">
                                        <label for="reason" class="block text-sm font-medium leading-6">
                                            {{ $t('common.reason') }}
                                        </label>
                                        <VeeField
                                            type="text"
                                            name="reason"
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.reason')"
                                            :label="$t('common.reason')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="my-2 space-y-24">
                                    <div class="flex-1">
                                        <label for="mugShot" class="block text-sm font-medium leading-6">
                                            {{ $t('common.image') }}
                                        </label>
                                        <template v-if="isNUIAvailable()">
                                            <p class="text-sm">
                                                {{ $t('system.not_supported_on_tablet.title') }}
                                            </p>
                                        </template>
                                        <template v-else>
                                            <VeeField
                                                v-slot="{ handleChange, handleBlur }"
                                                type="text"
                                                name="mugShot"
                                                :placeholder="$t('common.image')"
                                                :label="$t('common.image')"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            >
                                                <UInput
                                                    type="file"
                                                    accept="image/jpeg,image/jpg,image/png"
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    @change="handleChange"
                                                    @blur="handleBlur"
                                                />
                                            </VeeField>
                                            <VeeErrorMessage name="mugShot" as="p" class="mt-2 text-sm text-error-400" />
                                        </template>
                                    </div>
                                </div>
                                <div class="flex flex-1 items-center">
                                    <SquareImg
                                        class="m-auto"
                                        :url="user?.props?.mugShot?.url"
                                        size="6xl"
                                        :alt="$t('common.mug_shot')"
                                        :no-blur="true"
                                    />
                                </div>
                                <div class="absolute bottom-0 left-0 flex w-full">
                                    <UButton
                                        class="flex-1 rounded-md bg-neutral-50 px-3.5 py-2.5 text-sm font-semibold text-gray-900 hover:bg-gray-200"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </UButton>
                                    <UButton
                                        class="flex flex-1 justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold"
                                        :disabled="!meta.valid || !canSubmit"
                                        :loading="!canSubmit"
                                        @click="
                                            setFieldValue('reset', true);
                                            onSubmitThrottle($event);
                                        "
                                    >
                                        {{ $t('common.reset') }}
                                    </UButton>
                                    <UButton
                                        class="flex flex-1 justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold"
                                        :disabled="!meta.valid || !canSubmit"
                                        :loading="!canSubmit"
                                        @click="
                                            setFieldValue('reset', false);
                                            onSubmitThrottle($event);
                                        "
                                    >
                                        {{ $t('common.save') }}
                                    </UButton>
                                </div>
                            </form>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
