<script lang="ts" setup>
import { UNIT_STATUS, Unit } from '~~/gen/ts/resources/dispatch/units';

defineProps<{
    unit: Unit;
}>();

defineEmits<{
    (e: 'goto', loc: { x: number; y: number }): void;
    (e: 'details', unit: Unit): void;
}>();
</script>

<template>
    <li class="col-span-1 flex rounded-md shadow-sm" @click="$emit('details', unit)">
        <div
            class="flex w-12 flex-shrink-0 items-center justify-center rounded-l-md text-sm font-medium text-white border-l border-t border-b"
            :style="'background-color: #' + unit.color ?? '00000'"
        >
            {{ unit.initials }}
        </div>
        <div class="flex flex-1 items-center justify-between truncate rounded-r-md border border-gray-200 bg-gray">
            <div class="flex-1 truncate px-4 py-2 text-sm">
                <span class="font-medium text-gray-100">{{ unit.name }}</span>
                <p class="text-gray-400">{{ $t('common.members', unit.users.length) }}</p>
            </div>
            <div class="flex-shrink-0 pr-5 inline-flex items-center justify-center text-white">
                {{ $t(`enums.centrum.UNIT_STATUS.${UNIT_STATUS[unit.status?.status ?? (0 as number)]}`) }}
            </div>
        </div>
    </li>
</template>
