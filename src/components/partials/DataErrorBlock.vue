<script lang="ts" setup>
import { XCircleIcon } from '@heroicons/vue/20/solid';

defineProps({
    title: {
        type: String,
        required: false,
    },
    message: {
        type: String,
        required: false,
    },
    retry: {
        type: Function,
        required: false,
    },
    retryMessage: {
        type: String,
        required: false,
    },

});
</script>

<template>
    <div class="p-4 rounded-md bg-red-50 max-w-md mx-auto">
        <div class="flex">
            <div class="flex-shrink-0">
                <XCircleIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
            </div>
            <div class="ml-3">
                <h3 class="text-sm font-medium text-error-800">{{ title ?? $t('components.partials.data_error_block.default_title') }}</h3>
                <div class="mt-2 text-sm text-error-700">
                    <p>
                        {{ message ?? $t('components.partials.data_error_block.default_message') }}
                    </p>
                </div>
                <div v-if="retry" class="mt-4">
                    <div class="-mx-2 -my-1.5 flex">
                        <button type="button"
                            class="rounded-md bg-red-50 px-2 py-1.5 text-sm font-medium text-error-800 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-red-600 focus:ring-offset-2 focus:ring-offset-red-50"
                            @click="retry!()">
                            {{ retryMessage ?? $t('common.retry') }}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
