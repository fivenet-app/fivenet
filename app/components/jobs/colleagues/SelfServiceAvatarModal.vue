<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificationsStore();

const appConfig = useAppConfig();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.union([
    z.object({
        profilePicture: z.instanceof(File).optional(),
        reset: z.literal(false),
    }),
    z.object({ profilePicture: z.custom<File>().optional(), reset: z.literal(true) }),
]);

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    profilePicture: undefined,
    reset: false,
});

const { resizeAndUpload } = useFileUploader((_) => citizensCitizensClient.uploadAvatar(_), 'documents', 0);

async function uploadAvatar(f: File): Promise<void> {
    if (!f.type.startsWith('image/')) return;

    try {
        const resp = await resizeAndUpload(f);

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        activeChar.value!.profilePicture = resp.file?.filePath;

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as Error);
        throw e;
    }
}

async function deleteAvatar(): Promise<void> {
    try {
        await citizensCitizensClient.deleteAvatar({});

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        activeChar.value!.profilePicture = undefined;

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as Error);
        throw e;
    }
}

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.data.reset) {
        await deleteAvatar();
        return;
    }

    if (!event.data.profilePicture) return;

    await uploadAvatar(event.data.profilePicture);
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.jobs.self_service.set_profile_picture')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField name="profilePicture" :label="$t('common.profile_picture')">
                    <div class="flex flex-col gap-2">
                        <NotSupportedTabletBlock v-if="nuiEnabled" />
                        <div v-else class="flex flex-col gap-1">
                            <div class="flex flex-1 flex-row gap-1">
                                <UFileUpload
                                    v-model="state.profilePicture"
                                    class="flex-1"
                                    name="mugshot"
                                    :accept="appConfig.fileUpload.types.images.join(',')"
                                    block
                                    :disabled="formRef?.loading"
                                    :placeholder="$t('common.image')"
                                    :label="$t('common.file_upload_label')"
                                    :description="$t('common.allowed_file_types')"
                                />
                            </div>
                        </div>

                        <div class="flex w-full flex-col items-center justify-center gap-2">
                            <GenericImg
                                v-if="activeChar?.profilePicture"
                                size="3xl"
                                :src="`${activeChar.profilePicture}?date=${new Date().getTime()}`"
                                :no-blur="true"
                            />

                            <UAlert icon="i-mdi-information-outline" :description="$t('common.image_caching')" />
                        </div>
                    </div>
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    block
                    :disabled="formRef?.loading"
                    :loading="formRef?.loading"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />

                <UButton
                    class="flex-1"
                    block
                    color="error"
                    :disabled="formRef?.loading || !activeChar?.profilePicture"
                    :loading="formRef?.loading"
                    :label="$t('common.reset')"
                    @click="
                        state.reset = true;
                        formRef?.submit();
                    "
                />

                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UButtonGroup>
        </template>
    </UModal>
</template>
