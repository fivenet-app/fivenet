<script lang="ts" setup>
import { CloseCircleIcon } from 'mdi-vue3';
import { type DefineComponent } from 'vue';

withDefaults(
    defineProps<{
        title: string;
        message?: string | undefined;
        icon?: DefineComponent;
        type?: string;
    }>(),
    {
        message: undefined,
        icon: markRaw(CloseCircleIcon),
        type: 'error',
    },
);
</script>

<template>
    <div :class="`p-4 mt-6 rounded-md bg-${type}-100`">
        <div class="flex">
            <div class="flex-shrink-0">
                <component :is="icon" :class="`w-5 h-5 text-${type}-400`" aria-hidden="true" />
            </div>
            <div class="ml-3">
                <h3 :class="`text-sm font-medium text-${type}-600`">
                    {{ title }}
                </h3>
                <div v-if="message" :class="`mt-2 text-sm text-${type}-600`">
                    {{ message }}
                </div>
            </div>
        </div>
    </div>
</template>
