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
    <div class="w-full">
        <UButton
            :disabled="!focus"
            class="relative block w-full rounded-lg border-2 border-dashed border-base-300 p-8 text-center hover:border-base-400 focus:ring-2 focus:ring-neutral focus:ring-offset-2"
            @click="click()"
        >
            <component :is="icon" class="mx-auto size-12" />
            <span class="mt-2 block text-sm font-semibold text-gray-300">
                <span v-if="message">
                    {{ message }}
                </span>
                <span v-else>
                    {{ $t('common.not_found', [type ?? $t('common.data')]) }}
                </span>
            </span>
        </UButton>
    </div>
</template>
