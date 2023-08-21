<script lang="ts" setup>
import { AccountIcon } from 'mdi-vue3';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

const props = defineProps<{
    dispatch: Dispatch;
}>();

defineEmits<{
    (e: 'select', id: bigint): void;
    (e: 'goto', loc: Coordinate): void;
}>();

const expiresAt = props.dispatch.units.find((u) => u.expiresAt !== undefined)?.expiresAt;
</script>

<template>
    <div class="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
        <dt class="text-sm font-medium leading-6 text-white">
            <div class="flex h-6 items-center">
                <input
                    type="checkbox"
                    name="selected"
                    checked
                    class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-600 h-6 w-6"
                    @change="$emit('select')"
                />
                <IDCopyBadge class="ml-2" prefix="DSP" :id="dispatch.id" />
            </div>
            <div v-if="expiresAt" class="mt-1 text-white text-sm">
                Expires in {{ useLocaleTimeAgo(toDate(expiresAt), { showSecond: true, updateInterval: 1000 }).value }}
            </div>
        </dt>
        <dd class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0">
            <ul role="list" class="border divide-y rounded-md divide-base-200 border-base-200">
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                    {{ $t('common.message') }}: {{ dispatch.message }}
                </li>
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                    <button
                        v-if="dispatch.x && dispatch.y"
                        type="button"
                        class="text-primary-400 hover:text-primary-600"
                        @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
                    >
                        {{ $t('common.go_to_location') }}
                    </button>
                    <span v-else>{{ $t('common.no_location') }}</span>
                </li>
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                    <div class="flex items-center flex-1">
                        <AccountIcon class="flex-shrink-0 w-5 h-5 text-base-400" aria-hidden="true" />
                        {{ $t('common.member', 2) }}:
                        <span v-if="dispatch.units.length === 0">{{ $t('common.member', 0) }}</span>
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
