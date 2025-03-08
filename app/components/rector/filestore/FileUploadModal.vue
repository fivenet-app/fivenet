<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import NotSupportedTabletBlock from '~/components/partials/NotSupportedTabletBlock.vue';
import { useSettingsStore } from '~/store/settings';
import type { File, FileInfo } from '~~/gen/ts/resources/filestore/file';
import type { UploadFileResponse } from '~~/gen/ts/services/rector/filestore';

const emit = defineEmits<{
    (e: 'uploaded', file: FileInfo): void;
}>();

const { $grpc } = useNuxtApp();

const appConfig = useAppConfig();

const { isOpen } = useModal();

const settingsStore = useSettingsStore();
const { nuiEnabled } = storeToRefs(settingsStore);

const schema = z.object({
    category: z.string().min(3).max(255),
    name: z.string().min(3).max(255),
    file: zodFileSingleSchema(appConfig.fileUpload.fileSizes.rector, appConfig.fileUpload.types.images),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    category: '',
    name: '',
    file: undefined,
});

const categories = ['jobassets'];

async function uploadFile(values: Schema): Promise<UploadFileResponse | undefined> {
    const file = {} as File;

    if (!values.file[0]) {
        return;
    }

    file.data = new Uint8Array(await values.file[0].arrayBuffer());

    try {
        const { response } = await $grpc.rector.rectorFilestore.uploadFile({
            prefix: values.category,
            name: values.name,
            file: file,
        });

        if (response.file) {
            emit('uploaded', response.file);
        }

        isOpen.value = false;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await uploadFile(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('common.upload') }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="category" :label="$t('common.category')" class="flex-1" required>
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.category"
                                :options="categories"
                                :searchable-placeholder="$t('common.search_field')"
                            />
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup name="name" :label="$t('common.name')" class="flex-1" required>
                        <UInput v-model="state.name" type="text" name="name" :placeholder="$t('common.name')" />
                    </UFormGroup>

                    <UFormGroup name="file" :label="$t('common.file')" class="flex-1">
                        <NotSupportedTabletBlock v-if="nuiEnabled" />
                        <template v-else>
                            <UInput
                                type="file"
                                name="file"
                                accept="image/jpeg,image/jpg,image/png"
                                :placeholder="$t('common.image')"
                                @change="state.file = $event"
                            />
                        </template>
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton block class="flex-1" color="black" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
