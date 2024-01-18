<script lang="ts" setup>
import { CancelIcon, CloseCircleIcon, RefreshIcon } from 'mdi-vue3';

defineProps<{
    title?: string;
    message?: string;
    retry?: () => Promise<any>;
    retryMessage?: string;
}>();

const disabled = ref(true);

const timeout = ref<NodeJS.Timeout | undefined>();
onMounted(
    () =>
        (timeout.value = setTimeout(() => {
            timeout.value = undefined;
            disabled.value = false;
        }, 1250)),
);
onBeforeUnmount(() => {
    if (timeout.value) {
        clearTimeout(timeout.value);
        timeout.value = undefined;
    }
});
</script>

<template>
    <div class="mx-auto max-w-md rounded-md bg-error-100 p-4">
        <div class="flex">
            <div class="flex-shrink-0">
                <CloseCircleIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
            </div>
            <div class="ml-3">
                <h3 class="text-sm font-medium text-error-800">
                    {{ title ?? $t('components.partials.data_error_block.default_title') }}
                </h3>
                <div class="mt-2 text-sm text-error-700">
                    <p>
                        {{ message ?? $t('components.partials.data_error_block.default_message') }}
                    </p>
                </div>
                <div v-if="retry" class="mt-4">
                    <div class="-mx-2 -my-1.5 flex">
                        <button
                            type="button"
                            class="flex justify-center rounded-md px-2 py-1.5 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-error-600 focus:ring-offset-2 focus-visible:outline-error-500"
                            :disabled="disabled"
                            :class="[
                                disabled
                                    ? 'disabled bg-base-500 text-gray-200 hover:bg-base-400 focus-visible:outline-base-500'
                                    : 'bg-error-200 text-error-800 hover:bg-error-300 hover:bg-primary-400',
                            ]"
                            @click="retry()"
                        >
                            {{ retryMessage ?? $t('common.retry') }}
                            <RefreshIcon v-if="timeout === undefined" class="ml-2 h-5 w-5" />
                            <CancelIcon v-else class="ml-2 h-5 w-5" />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
