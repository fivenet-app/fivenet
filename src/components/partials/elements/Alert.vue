<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiCloseCircle } from '@mdi/js';

withDefaults(
    defineProps<{
        title: string;
        message?: string | undefined;
        icon?: string;
        type?: string;
    }>(),
    {
        message: undefined,
        icon: mdiCloseCircle,
        type: 'error',
    },
);
</script>

<template>
    <div :class="`p-4 mt-6 rounded-md bg-${type}-100`">
        <div class="flex">
            <div class="flex-shrink-0">
                <SvgIcon :class="`w-5 h-5 text-${type}-400`" aria-hidden="true" type="mdi" :path="icon" />
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
