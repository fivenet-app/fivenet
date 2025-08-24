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
    (e: 'update:mugshot', value?: FilestoreFile): void;
}>();

const modelValue = useVModel(props, 'user', emit);

const { isOpen } = useOverlay();

const notifications = useNotificationsStore();

const appConfig = useAppConfig();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z
    .object({
        reason: z.string().min(3).max(255),
        mugshot: z.custom<File>().array().min(1).max(1).default([]),
        reset: z.coerce.boolean(),
    })
    .or(
        z.union([
            z.object({
                reason: z.string().min(3).max(255),
                mugshot: z.custom<File>().array().min(1).max(1).default([]),
                reset: z.literal(false),
            }),
            z.object({
                reason: z.string().min(3).max(255),
                mugshot: z.custom<File>().array().default([]),
                reset: z.literal(true),
            }),
        ]),
    );

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    mugshot: [],
    reset: false,
});

const { resizeAndUpload } = useFileUploader((_) => citizensCitizensClient.uploadMugshot(_), 'documents', props.user.userId);

async function uploadMugshot(files: File[], reason: string): Promise<void> {
    for (const f of files) {
        if (!f.type.startsWith('image/')) continue;

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

            isOpen.value = false;
        } catch (e) {
            handleGRPCError(e as Error);
            throw e;
        }

        return;
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

        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as Error);
        throw e;
    }
}

function handleFileChanges(event: FileList) {
    state.mugshot = [...event];
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await (
        !event.data.reset
            ? uploadMugshot(event.data.mugshot, event.data.reason)
            : deleteMugshot(props.user.props?.mugshotFileId, event.data.reason)
    ).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard>
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl leading-6 font-semibold">
                            {{ $t('components.citizens.CitizenInfoProfile.set_mugshot') }}
                        </h3>

                        <UButton
                            class="-my-1"
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-window-close"
                            @click="isOpen = false"
                        />
                    </div>
                </template>

                <div>
                    <UFormField name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormField>

                    <UFormField name="mugshot" :label="$t('common.mugshot')">
                        <div class="flex flex-col gap-2">
                            <NotSupportedTabletBlock v-if="nuiEnabled" />
                            <div v-else class="flex flex-col gap-1">
                                <div class="flex flex-1 flex-row gap-1">
                                    <UInput
                                        class="flex-1"
                                        name="mugshot"
                                        type="file"
                                        :accept="appConfig.fileUpload.types.images.join(',')"
                                        block
                                        :placeholder="$t('common.image')"
                                        :disabled="!canSubmit"
                                        @change="handleFileChanges"
                                    />
                                </div>
                            </div>

                            <div class="flex w-full flex-col items-center justify-center gap-2">
                                <GenericImg
                                    v-if="user.props?.mugshot"
                                    size="3xl"
                                    :src="`${user.props?.mugshot.filePath}?date=${new Date().getTime()}`"
                                    :no-blur="true"
                                />

                                <UAlert icon="i-mdi-information-outline" :description="$t('common.image_caching')" />
                            </div>
                        </div>
                    </UFormField>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>

                        <UButton
                            class="flex-1"
                            type="submit"
                            block
                            color="error"
                            :disabled="!canSubmit || !user.props?.mugshotFileId"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.reset') }}
                        </UButton>

                        <UButton class="flex-1" block color="neutral" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
