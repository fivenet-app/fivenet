<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiPause, mdiRefresh } from '@mdi/js';

withDefaults(
    defineProps<{
        message: string;
        paused?: boolean;
    }>(),
    {
        paused: false,
    },
);
</script>

<template>
    <button
        type="button"
        disabled
        class="relative block w-full p-12 text-center border-2 border-dotted rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
    >
        <SvgIcon v-if="paused" class="w-12 h-12 mx-auto text-neutral" type="mdi" :path="mdiPause" />
        <SvgIcon v-else class="w-12 h-12 mx-auto text-neutral animate-spin" type="mdi" :path="mdiRefresh" />
        <span class="block mt-2 text-sm font-semibold text-gray-300">
            {{ message }}
        </span>
    </button>
</template>
