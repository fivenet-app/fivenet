<script lang="ts" setup>
import { max, min, required, mimes, size } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import type { File, FileInfo } from '~~/gen/ts/resources/filestore/file';
import type { UploadFileResponse } from '~~/gen/ts/services/rector/filestore';

const props = defineProps<{
    uploadedFile?: (file: FileInfo) => void;
}>();

const { $grpc } = useNuxtApp();

interface FormData {
    prefix: string;
    name: string;
    file: Blob;
}

const { isOpen } = useModal();

async function uploadFile(values: FormData): Promise<UploadFileResponse> {
    const file = {} as File;

    file.data = new Uint8Array(await values.file.arrayBuffer());

    try {
        const { response } = await $grpc.getRectorFilestoreClient().uploadFile({
            prefix: values.prefix,
            name: values.name,
            file,
        });

        if (response.file && props.uploadedFile) {
            props.uploadedFile(response.file);
        }

        isOpen.value = false;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const prefixes = ['jobassets'];

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('mimes', mimes);
defineRule('size', size);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        prefix: { required: true },
        name: { required: true, min: 3, max: 128 },
        file: { required: true, mimes: ['image/jpeg', 'image/jpg', 'image/png'], size: 2000 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<any> => await uploadFile(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
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
                <form @submit.prevent="onSubmitThrottle">
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="prefix" class="block text-sm font-medium leading-6">
                                {{ $t('common.category') }}
                            </label>
                            <VeeField
                                v-slot="{ field }"
                                name="prefix"
                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                :placeholder="$t('common.category')"
                                :label="$t('common.category')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            >
                                <select
                                    v-bind="field"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                >
                                    <option
                                        v-for="prefix in prefixes"
                                        :key="prefix"
                                        :selected="field.value === prefix"
                                        :value="prefix"
                                    >
                                        {{ prefix }}
                                    </option>
                                </select>
                            </VeeField>
                            <VeeErrorMessage name="prefix" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </div>
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="name" class="block text-sm font-medium leading-6">
                                {{ $t('common.name') }}
                            </label>
                            <VeeField
                                type="text"
                                name="name"
                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                :placeholder="$t('common.name')"
                                :label="$t('common.name')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </div>
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="file" class="block text-sm font-medium leading-6">
                                {{ $t('common.image') }}
                            </label>
                            <template v-if="isNUIAvailable()">
                                <p class="text-sm">
                                    {{ $t('system.not_supported_on_tablet.title') }}
                                </p>
                            </template>
                            <template v-else>
                                <VeeField
                                    v-slot="{ handleChange, handleBlur }"
                                    name="file"
                                    :placeholder="$t('common.image')"
                                    :label="$t('common.image')"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                >
                                    <input
                                        type="file"
                                        accept="image/jpeg,image/jpg,image/png"
                                        @change="handleChange"
                                        @blur="handleBlur"
                                    />
                                </VeeField>
                                <VeeErrorMessage name="file" as="p" class="mt-2 text-sm text-error-400" />
                            </template>
                        </div>
                    </div>
                    <div class="mt-5 gap-2 sm:mt-4 sm:flex">
                        <UButton @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                        <UButton type="submit" :disabled="!meta.valid || !canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </div>
                </form>
            </div>
        </UCard>
    </UModal>
</template>
