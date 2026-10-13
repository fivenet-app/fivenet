<script lang="ts" setup>
import UnitInfoPopover from '~/components/dispatch/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/stores/livemap';
import { type UnitStatus, StatusUnit } from '~~/gen/ts/resources/centrum/units/units';

const props = defineProps<{
    item: UnitStatus;
}>();

const { gotoCoords } = useLivemapStore();

const isMembershipChange = computed(
    () => props.item.status === StatusUnit.USER_ADDED || props.item.status === StatusUnit.USER_REMOVED,
);
const actor = computed(() => props.item.creator ?? props.item.user);
const showActor = computed(
    () =>
        actor.value !== undefined &&
        (!isMembershipChange.value ||
            props.item.user === undefined ||
            (props.item.creator !== undefined && props.item.creatorId !== props.item.userId)),
);
</script>

<template>
    <div class="flex min-w-0 flex-auto justify-between gap-2 text-xs leading-5 text-muted">
        <div class="flex min-w-0 flex-col gap-0.5">
            <span class="text-sm font-medium text-highlighted">
                {{ $t(`components.dispatch.units.feed.item.${StatusUnit[item.status]}`) }}
            </span>

            <span v-if="item.unit || item.unitId" class="inline-flex items-center gap-1">
                <UnitInfoPopover :unit-id="item.unitId" :unit="item.unit" initials-only badge show-icon />
            </span>
        </div>

        <div class="flex flex-none flex-col items-end gap-0.5 text-xs">
            <span class="text-xs/5 text-dimmed"><GenericTime :value="item.createdAt" type="compact" /></span>

            <span class="inline-flex items-center gap-1">
                <UButton
                    v-if="item.x !== undefined && item.y !== undefined"
                    variant="link"
                    size="xs"
                    icon="i-mdi-map-marker"
                    @click="gotoCoords({ x: item.x, y: item.y })"
                />
                <span v-if="isMembershipChange && item.user" class="inline-flex items-center gap-1">
                    <span class="text-dimmed">{{ $t('common.member', 1) }}:</span>
                    <CitizenInfoPopover :user="item.user" :trailing="false" text-class="text-xs" />
                </span>

                <span v-if="showActor && actor" class="inline-flex items-center gap-1">
                    <span v-if="isMembershipChange && item.user" class="text-dimmed">{{ $t('common.creator') }}:</span>
                    <CitizenInfoPopover :user="actor" :trailing="false" text-class="text-xs" />
                </span>

                <span v-if="!actor" class="text-dimmed">{{ $t('common.system') }}</span>
            </span>
        </div>
    </div>
</template>
