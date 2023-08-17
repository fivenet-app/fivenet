<script lang="ts" setup>
import { AccountIcon } from 'mdi-vue3';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

const props = defineProps<{
    dispatch: Dispatch;
}>();

const expiresAt = props.dispatch.units.find((u) => u.expiresAt !== undefined)?.expiresAt;
</script>

<template>
    <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
        <dt class="text-sm font-medium leading-6 text-white">
            DSP-{{ dispatch.id.toString() }}
            <br />
            <span v-if="expiresAt">
                Expires in {{ useLocaleTimeAgo(toDate(expiresAt), { showSecond: true, updateInterval: 1000 }).value }}
            </span>
        </dt>
        <dd class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0">
            <ul role="list" class="border divide-y rounded-md divide-base-200 border-base-200">
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                    {{ $t('common.message') }}: {{ dispatch.message }}
                </li>
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">{{ dispatch.x }}, {{ dispatch.y }}</li>
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                    <div class="flex items-center flex-1">
                        <AccountIcon class="flex-shrink-0 w-5 h-5 text-base-400" aria-hidden="true" />
                        {{ $t('common.members', 2) }}:
                        <span v-if="dispatch.units.length === 0">No members </span>
                        <span v-else class="flex-1 ml-2 truncate">
                            <span v-for="unit in dispatch.units">
                                {{ unit.unit?.name }}
                                ({{ unit.unit?.initials }})
                            </span>
                        </span>
                    </div>
                </li>
            </ul>
        </dd>
    </div>
</template>
