<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { File as FilestoreFile } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:mugshot', value?: FilestoreFile): void;
}>();

const modelValue = useVModel(props, 'user', emit);

const notifications = useNotificationsStore();

const appConfig = useAppConfig();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z
    .object({
        reason: z.coerce.string().min(3).max(255),
        mugshot: z.instanceof(File).optional(),
        reset: z.coerce.boolean(),
    })
    .or(
        z.union([
            z.object({
                reason: z.coerce.string().min(3).max(255),
                mugshot: z.instanceof(File).optional(),
                reset: z.literal(false),
            }),
            z.object({
                reason: z.coerce.string().min(3).max(255),
                mugshot: z.custom<File>().optional(),
                reset: z.literal(true),
            }),
        ]),
    );

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    mugshot: undefined,
    reset: false,
});

const { resizeAndUpload } = useFileUploader((_) => citizensCitizensClient.uploadMugshot(_), 'documents', props.user.userId);

async function uploadMugshot(f: File, reason: string): Promise<void> {
    if (!f.type.startsWith('image/')) return;

    try {
        const resp = await resizeAndUpload(f, reason);

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (modelValue.value.props) {
            modelValue.value.props.mugshot = resp.file;
        } else {
            modelValue.value.props = { userId: props.user.userId, mugshot: resp.file };
        }

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as Error);
        throw e;
    }
}

async function deleteMugshot(fileId: number | undefined, reason: string): Promise<void> {
    if (fileId === undefined) return;

    try {
        await citizensCitizensClient.deleteMugshot({
            userId: props.user.userId,
            reason: reason,
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
        if (modelValue.value.props) {
            modelValue.value.props.mugshot = undefined;
        } else {
            modelValue.value.props = { userId: props.user.userId, mugshot: undefined };
        }

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as Error);
        throw e;
    }
}

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.data.reset) {
        await deleteMugshot(props.user.props?.mugshotFileId, event.data.reason);
        return;
    }

    if (!event.data.mugshot) return;

    await uploadMugshot(event.data.mugshot, event.data.reason);
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.citizens.CitizenInfoProfile.set_mugshot')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" class="w-full" />
                </UFormField>

                <UFormField name="mugshot" :label="$t('common.mugshot')">
                    <div class="flex flex-col gap-2">
                        <NotSupportedTabletBlock v-if="nuiEnabled" />
                        <UFileUpload
                            v-else
                            v-model="state.mugshot"
                            class="mx-auto max-w-md flex-1"
                            name="mugshot"
                            block
                            :disabled="formRef?.loading"
                            :accept="appConfig.fileUpload.types.images.join(',')"
                            :placeholder="$t('common.image')"
                            :label="$t('common.file_upload_label')"
                            :description="$t('common.allowed_file_types')"
                        />

                        <div class="flex w-full flex-col items-center justify-center gap-2">
                            <GenericImg
                                v-if="user.props?.mugshot"
                                size="3xl"
                                :src="`${user.props?.mugshot.filePath}?date=${new Date().getTime()}`"
                                no-blur
                            />

                            <UAlert icon="i-mdi-information-outline" :description="$t('common.image_caching')" />
                        </div>
                    </div>
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
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
                    :disabled="formRef?.loading || !user.props?.mugshotFileId"
                    :loading="formRef?.loading"
                    :label="$t('common.reset')"
                    @click="
                        state.reset = true;
                        formRef?.submit();
                    "
                />

                <UButton class="flex-1" block color="neutral" :label="$t('common.close', 1)" @click="$emit('close', false)" />
            </UFieldGroup>
        </template>
    </UModal>
</template>
