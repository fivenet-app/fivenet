<script lang="ts" setup>
import DispatchStatusInfoPopover from '~/components/dispatch/dispatches/DispatchStatusInfoPopover.vue';
import UnitInfoPopover from '~/components/dispatch/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/stores/livemap';
import { type DispatchStatus, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';

withDefaults(
    defineProps<{
        item: DispatchStatus;
        showId?: boolean;
    }>(),
    {
        showId: false,
    },
);

const { gotoCoords } = useLivemapStore();
</script>

<template>
    <div class="flex min-w-0 flex-auto justify-between gap-2 text-xs leading-5 text-gray-200">
        <div class="flex min-w-0 flex-col gap-0.5">
            <span class="inline-flex items-center gap-1">
                <span class="text-sm font-medium text-highlighted">{{
                    $t(`components.dispatch.dispatches.feed.item.${StatusDispatch[item.status]}`)
                }}</span>

                <DispatchStatusInfoPopover v-if="showId" :status="item" />
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
                <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
            </span>
        </div>
    </div>
</template>
