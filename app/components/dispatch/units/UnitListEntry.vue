<script lang="ts" setup>
import { defaultUnitIcon, unitStatusToBGColor } from '~/components/dispatch/helpers';
import UnitDetailsSlideover from '~/components/dispatch/units/UnitDetailsSlideover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { rgbBlack } from '~/utils/color';
import { type Unit, StatusUnit } from '~~/gen/ts/resources/centrum/units/units';

const props = defineProps<{
    unit: Unit;
    showMembers?: boolean;
}>();

const overlay = useOverlay();

const unitDetailsSlideover = overlay.create(UnitDetailsSlideover, {
    props: {
        unit: props.unit,
    },
});

const unitColorHex = computed(() => hexToRgb(props.unit.color, rgbBlack)!);
const isBright = computed(() => isColorBright(unitColorHex.value));
</script>

<template>
    <li class="col-span-1 overflow-hidden rounded-md shadow-xs" @click="unitDetailsSlideover.open()">
        <div class="flex">
            <div
                class="flex w-12 shrink-0 items-center justify-center border-y border-l border-default text-sm font-medium"
                :class="isBright ? 'text-black' : 'text-white'"
                :style="'background-color: ' + unit.color"
            >
                {{ unit.initials }}
            </div>

            <div class="flex flex-1 items-center justify-between truncate border border-default">
                <div class="flex-1 px-1 py-1.5 text-sm">
                    <div class="inline-flex items-center gap-1 font-medium">
                        <UIcon
                            v-if="unit.icon && unit.icon !== defaultUnitIcon"
                            class="!size-4"
                            :name="convertComponentIconNameToDynamic(unit.icon)"
                            :style="{ color: unit.color ?? 'currentColor' }"
                        />

                        <span>{{ unit.name }}</span>
                    </div>

                    <p class="text-sm text-muted">
                        {{ $t('common.member', unit.users.length) }}
                    </p>
                </div>
            </div>

            <div
                class="flex w-20 shrink-0 items-center justify-center border-y border-r border-default text-center text-sm font-medium"
                :class="unitStatusToBGColor(unit.status?.status)"
            >
                {{ $t(`enums.centrum.StatusUnit.${StatusUnit[unit.status?.status ?? 0]}`) }}
            </div>
        </div>

        <div
            v-if="showMembers && unit.users.length > 0"
            class="border-x border-b border-default bg-elevated px-2 py-1.5"
            @click.stop
        >
            <ul class="max-h-20 divide-y divide-default overflow-y-auto text-sm" role="list">
                <li v-for="user in unit.users" :key="user.userId" class="py-0.5 first:pt-0 last:pb-0">
                    <CitizenInfoPopover :user="user.user" show-avatar-in-name />
                </li>
            </ul>
        </div>
    </li>
</template>
