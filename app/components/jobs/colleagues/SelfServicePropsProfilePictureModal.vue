<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { useNotificatorStore } from '~/stores/notificator';
import { useSettingsStore } from '~/stores/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { SetProfilePictureRequest } from '~~/gen/ts/services/citizenstore/citizenstore';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificatorStore();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const appConfig = useAppConfig();

const schema = z.object({
    avatar: zodFileSingleSchema(appConfig.fileUpload.fileSizes.images, appConfig.fileUpload.types.images),
    reset: z.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    avatar: undefined,
    reset: false,
});

async function setProfilePicture(values: Schema): Promise<void> {
    const req = {} as SetProfilePictureRequest;
    if (!values.reset) {
        if (!values.avatar || !values.avatar[0]) {
            return;
        }

        req.avatar = { data: new Uint8Array(await values.avatar[0].arrayBuffer()) };
    } else {
        req.avatar = { data: new Uint8Array(), delete: true };

        state.reset = false;
    }

    try {
        const call = $grpc.citizenstore.citizenStore.setProfilePicture(req);
        const { response } = await call;

        if (activeChar.value) {
            activeChar.value.avatar = response.avatar;
        }

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await setProfilePicture(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :state="state" :schema="schema" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.jobs.self_service.set_profile_picture') }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <div>
                        <UFormGroup name="avatar" class="flex-1" :label="$t('common.avatar')">
                            <NotSupportedTabletBlock v-if="nuiEnabled" />
                            <UInput
                                v-else
                                type="file"
                                name="avatar"
                                :placeholder="$t('common.image')"
                                :accept="appConfig.fileUpload.types.images.join(',')"
                                @change="state.avatar = $event"
                            />
                        </UFormGroup>
                    </div>

                    <div class="flex flex-1 items-center">
                        <ProfilePictureImg
                            class="m-auto"
                            :src="activeChar?.avatar?.url"
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
                            color="error"
                            class="flex-1"
                            :disabled="nuiEnabled || !canSubmit || !activeChar?.avatar"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.reset') }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="nuiEnabled || !canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
