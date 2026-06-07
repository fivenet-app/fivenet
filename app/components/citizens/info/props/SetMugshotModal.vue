<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { File as FilestoreFile } from '~~/gen/ts/resources/file/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { User } from '~~/gen/ts/resources/users/user';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update:mugshot', value?: FilestoreFile): void;
}>();

const user = defineModel<User>('user', { required: true });

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

const { hasUnsavedChanges, confirmLeave } = useSnapshotChanges(state, {
    serializer: (value) =>
        JSON.stringify({
            reason: value.reason,
            reset: value.reset,
            mugshot: value.mugshot
                ? {
                      name: value.mugshot.name,
                      size: value.mugshot.size,
                      type: value.mugshot.type,
                  }
                : null,
        }),
});

const { resizeAndUpload } = useFileUploader((_) => citizensCitizensClient.uploadMugshot(_), 'documents', user.value.userId);
const { uploadImages } = useImageUpload();

async function uploadMugshot(f: File, reason: string): Promise<void> {
    try {
        const result = await uploadImages({
            files: [f],
            uploadOne: (file) => resizeAndUpload(file, reason),
            invalidTypeNotification: {
                title: {
                    key: 'components.partials.tiptap_editor.notifications.invalid_file_type_images.title',
                    parameters: {},
                },
                description: {
                    key: 'components.partials.tiptap_editor.notifications.invalid_file_type_images.content',
                    parameters: {},
                },
            },
            successNotification: {
                title: { key: 'notifications.action_successful.title', parameters: {} },
                description: { key: 'notifications.action_successful.content', parameters: {} },
            },
            onUploaded: (resp) => {
                if (user.value.props) {
                    user.value.props.mugshot = resp.file;
                } else {
                    user.value.props = { userId: user.value.userId, mugshot: resp.file };
                }
            },
        });

        if (!result.ok) return;

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
            userId: user.value.userId,
            reason: reason,
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
        if (user.value.props) {
            user.value.props.mugshot = undefined;
        } else {
            user.value.props = { userId: user.value.userId, mugshot: undefined };
        }

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as Error);
        throw e;
    }
}

const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.data.reset) {
        await deleteMugshot(user.value.props?.mugshotFileId, event.data.reason);
        return;
    }

    if (!event.data.mugshot) return;

    await uploadMugshot(event.data.mugshot, event.data.reason);
}, 1000);

const formRef = useTemplateRef('formRef');

async function closeModal(): Promise<void> {
    if (formRef.value?.loading) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal
        :title="$t('components.citizens.CitizenInfoProfile.set_mugshot')"
        :close="false"
        :dismissible="!hasUnsavedChanges && !formRef?.loading"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{ $t('components.citizens.CitizenInfoProfile.set_mugshot') }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="formRef?.loading"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField name="reason" :label="$t('common.reason')" required>
                    <UInput v-model="state.reason" class="w-full" type="text" :placeholder="$t('common.reason')" />
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

                <UButton
                    class="flex-1"
                    block
                    color="neutral"
                    :disabled="formRef?.loading"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
