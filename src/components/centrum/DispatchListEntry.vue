<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccountMultiplePlus, mdiDetails } from '@mdi/js';
import { DISPATCH_STATUS, Dispatch } from '../../../gen/ts/resources/dispatch/dispatch';
import DispatchDetails from './DispatchDetails.vue';

defineProps<{
    dispatch: Dispatch;
}>();

const open = ref(false);
</script>

<template>
    <DispatchDetails @close="open = false" :dispatch="dispatch" :open="open" />
    <td
        class="relative whitespace-nowrap py-2 pl-2 text-right text-sm font-medium sm:pr-0 max-w-[42px] flex flex-row justify-start"
    >
        <button class="text-primary-400 hover:text-primary-600" :title="$t('common.detail', 2)" @click="open = true">
            <SvgIcon type="mdi" :path="mdiDetails" class="w-6 h-auto ml-auto mr-2.5" />
        </button>
        <button class="text-primary-400 hover:text-primary-600" :title="$t('common.assign')">
            <SvgIcon type="mdi" :path="mdiAccountMultiplePlus" class="w-6 h-auto ml-auto mr-2.5" aria-hidden="true" />
        </button>
    </td>
    <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
        {{ dispatch.id }}
    </td>
    <td class="whitespace-nowrap px-2 py-2 text-sm font-medium text-gray-100">
        {{ DISPATCH_STATUS[dispatch.status as number] }}
    </td>
    <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
        {{ dispatch.units }}
    </td>
    <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">{{ dispatch.marker }}</td>
    <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-100">{{ dispatch.message }}</td>
</template>
