<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useAuthStore } from '~/stores/auth';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificationsStore();

const appConfig = useAppConfig();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const schema = z
    .object({
        avatar: z.custom<File>().array().min(1).max(1),
        reset: z.boolean(),
    })
    .or(
        z.union([
            z.object({
                avatar: z.custom<File>().array().min(1).max(1),
                reset: z.literal(false),
            }),
            z.object({ avatar: z.custom<File>().array(), reset: z.literal(true) }),
        ]),
    );

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    avatar: [],
    reset: false,
});

const { resizeAndUpload } = useFileUploader((_) => $grpc.citizens.citizens.uploadAvatar(_), 'documents', 0);

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

            activeChar.value!.avatar = resp.file?.filePath;

            isOpen.value = false;
        } catch (e) {
            handleGRPCError(e as Error);
            throw e;
        }

        return;
    }
}

async function deleteAvatar(): Promise<void> {
    try {
        await $grpc.citizens.citizens.deleteAvatar({});

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        activeChar.value!.avatar = undefined;

        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as Error);
        throw e;
    }
}

function handleFileChanges(event: FileList) {
    state.avatar = [...event];
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await (!event.data.reset ? uploadAvatar(event.data.avatar) : deleteAvatar()).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.jobs.self_service.set_profile_picture') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="avatar" :label="$t('common.avatar')">
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
                                    v-if="activeChar?.avatar"
                                    size="3xl"
                                    :src="`${activeChar.avatar}?date=${new Date().getTime()}`"
                                    :no-blur="true"
                                />

                                <UAlert icon="i-mdi-information-outline" :description="$t('common.image_caching')" />
                            </div>
                        </div>
                    </UFormGroup>
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
                            :disabled="!canSubmit || !activeChar?.avatar"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.reset') }}
                        </UButton>

                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
