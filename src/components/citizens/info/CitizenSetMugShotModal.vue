<script lang="ts" setup>
import { max, min, required, mimes, size } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import SquareImg from '~/components/partials/elements/SquareImg.vue';
import { useNotificatorStore } from '~/store/notificator';
import type { File } from '~~/gen/ts/resources/filestore/file';
import { User, UserProps } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'update:mugShot', value?: File): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

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

        isOpen.value = false;
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
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.citizens.citizen_info_profile.set_mug_shot') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UForm :state="{}" @submit.prevent="onSubmitThrottle">
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
                            :url="user?.props?.mugShot?.url"
                            class="m-auto"
                            size="3xl"
                            :alt="$t('common.mug_shot')"
                            :no-blur="true"
                        />
                    </div>
                </UForm>
            </div>

            <template #footer>
                <div class="flex">
                    <UButton @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
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
            </template>
        </UCard>
    </UModal>
</template>
