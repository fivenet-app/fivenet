<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import GenericImg from '~/components/partials/elements/GenericImg.vue';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useNotificatorStore } from '~/stores/notificator';
import { useSettingsStore } from '~/stores/settings';
import type { File as FilestoreFile } from '~~/gen/ts/resources/filestore/file';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserProps } from '~~/gen/ts/resources/users/props';
import type { User } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    user: User;
}>();

const emit = defineEmits<{
    (e: 'update:mugShot', value?: FilestoreFile): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const appConfig = useAppConfig();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const mugShotSchema = zodFileSingleSchema(appConfig.fileUpload.fileSizes.images, appConfig.fileUpload.types.images);
const schema = z
    .object({
        reason: z.string().min(3).max(255),
        mugShot: mugShotSchema,
        reset: z.boolean(),
    })
    .or(
        z.union([
            z.object({ reason: z.string().min(3).max(255), mugShot: z.custom<FileList>(), reset: z.boolean() }),
            z.object({ reason: z.string().min(3).max(255), mugShot: z.undefined(), reset: z.boolean() }),
        ]),
    );

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    mugShot: undefined,
    reset: false,
});

async function setMugShot(values: Schema): Promise<void> {
    const userProps: UserProps = {
        userId: props.user.userId,
    };
    if (!values.reset) {
        if (!values.mugShot || !values.mugShot[0]) {
            return;
        }

        userProps.mugShot = { data: new Uint8Array(await values.mugShot[0].arrayBuffer()) };
    } else {
        userProps.mugShot = { data: new Uint8Array(), delete: true };

        state.reset = false;
    }

    try {
        const call = $grpc.citizens.citizens.setUserProps({
            props: userProps,
            reason: values.reason,
        });
        const { response } = await call;

        emit('update:mugShot', response.props?.mugShot);

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
    await setMugShot(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.citizens.CitizenInfoProfile.set_mug_shot') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                    </UFormGroup>

                    <UFormGroup name="mugShot" :label="$t('common.image')">
                        <NotSupportedTabletBlock v-if="nuiEnabled" />
                        <template v-else>
                            <UInput
                                type="file"
                                name="mugShot"
                                :accept="appConfig.fileUpload.types.images.join(',')"
                                :placeholder="$t('common.image')"
                                @change="state.mugShot = $event"
                            />
                        </template>
                    </UFormGroup>

                    <div class="mt-4 flex flex-1 items-center">
                        <GenericImg
                            class="m-auto"
                            :src="user?.props?.mugShot?.url"
                            size="3xl"
                            :alt="$t('common.mug_shot')"
                            :no-blur="true"
                        />
                    </div>
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
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                            @click="state.reset = true"
                        >
                            {{ $t('common.reset') }}
                        </UButton>

                        <UButton class="flex-1" block color="black" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
