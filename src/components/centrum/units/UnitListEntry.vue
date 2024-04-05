<script lang="ts" setup>
import { unitStatusToBGColor } from '~/components/centrum/helpers';
import { type RGB } from '~/utils/colour';
import { StatusUnit, Unit } from '~~/gen/ts/resources/centrum/units';
import UnitDetailsSlideover from '~/components/centrum/units/UnitDetailsSlideover.vue';

const props = defineProps<{
    unit: Unit;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const unitColorHex = computed(() => hexToRgb(props.unit.color) ?? ({ r: 0, g: 0, b: 0 } as RGB));
const isBright = computed(() => isColourBright(unitColorHex.value));

const open = ref(false);
</script>

<template>
    <li class="col-span-1 flex rounded-md shadow-sm" @click="open = true">
        <UnitDetailsSlideover :open="open" :unit="unit" @close="open = false" @goto="$emit('goto', $event)" />

        <div
            class="flex w-12 shrink-0 items-center justify-center rounded-l-md border-y border-l text-sm font-medium"
            :class="isBright ? 'text-black' : 'text-neutral'"
            :style="'background-color: ' + unit.color"
        >
            {{ unit.initials }}
        </div>
        <div class="flex flex-1 items-center justify-between truncate border border-gray-200">
            <div class="flex-1 px-1 py-2 text-sm">
                <span class="font-medium">{{ unit.name }}</span>
                <p :class="unit.users.length === 0 ? 'text-gray-400' : 'text-gray-300'">
                    {{ $t('common.member', unit.users.length) }}
                </p>
            </div>
        </div>
        <div
            class="flex w-20 shrink-0 items-center justify-center rounded-r-md border-y border-r text-center text-sm font-medium"
            :class="unitStatusToBGColor(unit.status?.status)"
        >
            {{ $t(`enums.centrum.StatusUnit.${StatusUnit[unit.status?.status ?? 0]}`) }}
        </div>
    </li>
</template>
