<script lang="ts" setup>
import { mimes, required, size } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import type { SetProfilePictureRequest } from '~~/gen/ts/services/citizenstore/citizenstore';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificatorStore();

const modal = useModal();

interface FormData {
    avatar?: Blob;
    reset?: boolean;
}

async function setProfilePicture(values: FormData): Promise<void> {
    const req = {} as SetProfilePictureRequest;
    if (!values.reset) {
        if (!values.avatar) {
            return;
        }

        req.avatar = { data: new Uint8Array(await values.avatar.arrayBuffer()) };
    } else {
        req.avatar = { data: new Uint8Array(), delete: true };
    }

    try {
        const call = $grpc.getCitizenStoreClient().setProfilePicture(req);
        const { response } = await call;

        if (activeChar.value) {
            activeChar.value.avatar = response.avatar;
        }

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        modal.close();
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('size', size);
defineRule('mimes', mimes);

const { handleSubmit, meta, setFieldValue } = useForm<FormData>({
    validationSchema: {
        avatar: { required: false, mimes: ['image/jpeg', 'image/jpg', 'image/png'], size: 2000 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await setProfilePicture(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
    setFieldValue('reset', false);
}, 1000);

const nuiAvailable = ref(isNUIAvailable());
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.jobs.self_service.set_profile_picture') }}
                    </h3>

                    <UButton
                        color="gray"
                        variant="ghost"
                        icon="i-heroicons-x-mark-20-solid"
                        class="-my-1"
                        @click="modal.close()"
                    />
                </div>
            </template>

            <div>
                <form @submit.prevent="onSubmitThrottle">
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="avatar" class="block text-sm font-medium leading-6">
                                {{ $t('common.avatar') }}
                            </label>
                            <template v-if="nuiAvailable">
                                <p class="text-sm">
                                    {{ $t('system.not_supported_on_tablet.title') }}
                                </p>
                            </template>
                            <template v-else>
                                <VeeField
                                    v-slot="{ handleChange, handleBlur }"
                                    name="avatar"
                                    :placeholder="$t('common.image')"
                                    :label="$t('common.image')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                >
                                    <UInput
                                        type="file"
                                        accept="image/jpeg,image/jpg,image/png"
                                        @change="handleChange"
                                        @blur="handleBlur"
                                    />
                                </VeeField>
                                <VeeErrorMessage name="avatar" as="p" class="mt-2 text-sm text-error-400" />
                            </template>
                        </div>
                    </div>
                    <div class="flex flex-1 items-center">
                        <ProfilePictureImg
                            class="m-auto"
                            :url="activeChar?.avatar?.url"
                            :name="`${activeChar?.firstname} ${activeChar?.lastname}`"
                            size="huge"
                            :no-blur="true"
                        />
                    </div>
                </form>
            </div>

            <div class="mt-5 gap-2 sm:mt-4 sm:flex">
                <UButton
                    class="flex-1 rounded-md bg-neutral-50 px-3.5 py-2.5 text-sm font-semibold text-gray-900 hover:bg-gray-200"
                    @click="modal.close()"
                >
                    {{ $t('common.close', 1) }}
                </UButton>
                <UButton
                    class="flex flex-1 justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold"
                    :disabled="nuiAvailable || !meta.valid || !canSubmit || !activeChar?.avatar"
                    :class="[
                        nuiAvailable || !meta.valid || !canSubmit || !activeChar?.avatar
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-error-500 hover:bg-error-400 focus-visible:outline-error-500',
                    ]"
                    @click="
                        setFieldValue('reset', true);
                        onSubmitThrottle($event);
                    "
                >
                    <template v-if="!canSubmit">
                        <LoadingIcon class="mr-2 size-5 animate-spin" />
                    </template>
                    {{ $t('common.reset') }}
                </UButton>
                <UButton
                    type="submit"
                    class="flex flex-1 justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold"
                    :disabled="nuiAvailable || !meta.valid || !canSubmit"
                    :class="[
                        nuiAvailable || !meta.valid || !canSubmit
                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                            : 'bg-primary-500 hover:bg-primary-400',
                    ]"
                >
                    <template v-if="!canSubmit">
                        <LoadingIcon class="mr-2 size-5 animate-spin" />
                    </template>
                    {{ $t('common.save') }}
                </UButton>
            </div>
        </UCard>
    </UModal>
</template>
