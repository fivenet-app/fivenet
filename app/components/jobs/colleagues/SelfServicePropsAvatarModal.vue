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

const schema = z
    .object({
        profilePicture: z.custom<File>().array().min(1).max(1).default([]),
        reset: z.coerce.boolean(),
    })
    .or(
        z.union([
            z.object({
                profilePicture: z.custom<File>().array().min(1).max(1).default([]),
                reset: z.literal(false),
            }),
            z.object({ profilePicture: z.custom<File>().array(), reset: z.literal(true) }),
        ]),
    );

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    profilePicture: [],
    reset: false,
});

const { resizeAndUpload } = useFileUploader((_) => citizensCitizensClient.uploadAvatar(_), 'documents', 0);

async function uploadAvatar(files: File[]): Promise<void> {
    for (const f of files) {
        if (!f.type.startsWith('image/')) continue;

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

        return;
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

function handleFileChanges(event: FileList) {
    state.profilePicture = [...event];
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await (!event.data.reset ? uploadAvatar(event.data.profilePicture) : deleteAvatar()).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
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
                                    class="flex-1"
                                    name="mugshot"
                                    :accept="appConfig.fileUpload.types.images.join(',')"
                                    block
                                    :placeholder="$t('common.image')"
                                    :disabled="!canSubmit"
                                    @update:model-value="($event) => handleFileChanges($event)"
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
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />

                <UButton
                    class="flex-1"
                    block
                    color="error"
                    :disabled="!canSubmit || !activeChar?.profilePicture"
                    :loading="!canSubmit"
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
