<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiMagnify } from '@mdi/js';

const { t } = useI18n();

const props = withDefaults(
    defineProps<{
        message?: string;
        icon?: string;
        type?: string;
        focus?: Function;
    }>(),
    {
        icon: mdiMagnify,
    },
);
function click() {
    if (props.focus) {
        props.focus();
    }
}
</script>

<template>
    <button
        type="button"
        :disabled="!focus"
        @click="click"
        class="relative block w-full p-12 text-center border-2 border-dashed rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
    >
        <SvgIcon class="w-12 h-12 mx-auto text-neutral" type="mdi" :path="icon" />
        <span class="block mt-2 text-sm font-semibold text-gray-300">
            <span v-if="message">
                {{ message }}
            </span>
            <span v-else>
                {{ $t('common.not_found', [type ?? t('common.data')]) }}
            </span>
        </span>
    </button>
</template>
