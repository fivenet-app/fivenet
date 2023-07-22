<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiCarEmergency } from '@mdi/js';
import Details from '~/components/centrum/dispatches/Details.vue';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

defineProps<{
    dispatch: Dispatch;
}>();

const detailsOpen = ref(false);
</script>

<template>
    <li>
        <button
            type="button"
            class="text-white bg-error-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
            @click="detailsOpen = true"
        >
            <SvgIcon type="mdi" :path="mdiCarEmergency" class="h-5 w-5" aria-hidden="true" />
            <span class="mt-2 truncate">DSP-{{ dispatch.id }}</span>
            <span class="mt-2 truncate">{{ DISPATCH_STATUS[dispatch.status?.status!] }}</span>
        </button>
        <Details @close="detailsOpen = false" :dispatch="dispatch" :open="detailsOpen" />
    </li>
</template>
