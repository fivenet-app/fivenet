<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import type { SetProfilePictureRequest } from '~~/gen/ts/services/citizenstore/citizenstore';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificatorStore();

const { isOpen } = useModal();

const MAX_FILE_SIZE = 2 * 1024 * 1024;
const ACCEPTED_IMAGE_TYPES = ['image/jpeg', 'image/jpg', 'image/png'];

const schema = z.object({
    reason: z.string().min(3).max(255),
    avatar: z.custom<FileList>().superRefine((files, ctx) => {
        if (!files || files.length === 0) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: 'File must be provided',
            });
            return false;
        }

        if (!ACCEPTED_IMAGE_TYPES.includes(files[0].type)) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: 'File must be a valid image type',
            });
            return false;
        }

        if (files[0].size > MAX_FILE_SIZE) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: 'File must be less than 5MB',
            });
            return false;
        }

        return true;
    }),
    reset: z.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    reason: '',
    avatar: undefined,
    reset: false,
});

async function setProfilePicture(values: Schema): Promise<void> {
    const req = {} as SetProfilePictureRequest;
    if (!values.reset) {
        if (!values.avatar) {
            return;
        }

        req.avatar = { data: new Uint8Array(await values.avatar[0].arrayBuffer()) };
    } else {
        req.avatar = { data: new Uint8Array(), delete: true };

        state.reset = false;
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

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setProfilePicture(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const nuiAvailable = ref(isNUIAvailable());
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :state="state" @submit="onSubmitThrottle">
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
                            @click="isOpen = false"
                        />
                    </div>
                </template>

                <div>
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
                            size="3xl"
                            :no-blur="true"
                        />
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton
                            type="submit"
                            block
                            color="red"
                            class="flex-1"
                            :disabled="nuiAvailable || !canSubmit || !activeChar?.avatar"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.reset') }}
                        </UButton>

                        <UButton
                            type="submit"
                            block
                            class="flex-1"
                            :disabled="nuiAvailable || !canSubmit"
                            :loading="!canSubmit"
                        >
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
