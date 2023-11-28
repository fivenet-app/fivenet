<script lang="ts" setup>
import { MagnifyIcon } from 'mdi-vue3';
import { type DefineComponent } from 'vue';

const props = withDefaults(
    defineProps<{
        message?: string;
        icon?: DefineComponent;
        type?: string;
        focus?: Function;
    }>(),
    {
        message: undefined,
        icon: markRaw(MagnifyIcon),
        type: undefined,
        focus: undefined,
    },
);

function click() {
    if (props.focus) {
        props.focus();
    }
}
</script>

<template>
    <div class="pt-4">
        <button
            type="button"
            :disabled="!focus"
            class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
            @click="click()"
        >
            <component :is="icon" class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                <span v-if="message">
                    {{ message }}
                </span>
                <span v-else>
                    {{ $t('common.not_found', [type ?? $t('common.data')]) }}
                </span>
            </span>
        </button>
    </div>
</template>
