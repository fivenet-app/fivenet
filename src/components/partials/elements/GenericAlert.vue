<script lang="ts" setup>
import { CloseCircleIcon, GestureTapIcon } from 'mdi-vue3';
import { type DefineComponent } from 'vue';

defineEmits<{
    (e: 'clicked'): void;
}>();

withDefaults(
    defineProps<{
        title: string;
        message?: string | undefined;
        icon?: DefineComponent;
        type?: string;
        callbackMessage?: string;
    }>(),
    {
        message: undefined,
        icon: markRaw(CloseCircleIcon),
        type: 'error',
        callbackMessage: undefined,
    },
);
</script>

<template>
    <div class="mt-6 rounded-md p-4" :class="`bg-${type}-100`">
        <div class="flex">
            <div class="flex-shrink-0">
                <component :is="icon" class="h-5 w-5" :class="`text-${type}-400`" aria-hidden="true" />
            </div>
            <div class="ml-3">
                <h3 class="text-sm font-medium" :class="`text-${type}-600`">
                    {{ title }}
                </h3>
                <div v-if="message" class="mt-2 text-sm" :class="`text-${type}-600`">
                    {{ message }}
                </div>
                <div v-if="callbackMessage" class="mt-4">
                    <div class="-mx-2 -my-1.5 flex">
                        <button
                            type="button"
                            class="flex justify-center rounded-md px-2 py-1.5 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2"
                            :class="`bg-${type}-200 text-${type}-800 hover:bg-${type}-400 focus:ring-${type}-600 focus-visible:outline-${type}-500`"
                            @click="$emit('clicked')"
                        >
                            {{ callbackMessage ?? $t('common.retry') }}
                            <GestureTapIcon class="ml-2 h-5 w-5" />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
