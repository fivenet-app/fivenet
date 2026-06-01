<script lang="ts" setup>
import DispatchAttributes from '~/components/dispatch/partials/DispatchAttributes.vue';
import UnitInfoPopover from '~/components/dispatch/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useLivemapStore } from '~/stores/livemap';
import { CentrumAccessLevel } from '~~/gen/ts/resources/centrum/access/access';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';
import DispatchAssignModal from '../dispatches/DispatchAssignModal.vue';
import DispatchStatusBadge from '../partials/DispatchStatusBadge.vue';
import { checkDispatchAccess } from '../helpers';

const props = defineProps<{
    dispatch: Dispatch;
}>();

const emit = defineEmits<{
    selected: [];
}>();

const overlay = useOverlay();

const dispatchAssignModal = overlay.create(DispatchAssignModal);

const { gotoCoords } = useLivemapStore();

const { selfAssign, canDo } = useCentrumStore();
</script>

<template>
    <UCard
        class="-my-[13px] -mr-[24px] -ml-[20px] flex min-w-[200px] flex-col"
        :ui="{ header: 'p-1 sm:px-2', body: 'p-1 sm:p-2 xl:mx-auto', footer: 'p-1 sm:px-2' }"
    >
        <template #header>
            <div class="grid grid-cols-1 gap-2 !text-primary md:grid-cols-2">
                <UButton
                    v-if="props.dispatch?.x !== undefined && props.dispatch?.y !== undefined"
                    variant="link"
                    icon="i-mdi-map-marker"
                    block
                    :label="$t('common.mark')"
                    @click="gotoCoords({ x: props.dispatch?.x, y: props.dispatch?.y })"
                />

                <UButton
                    variant="link"
                    icon="i-mdi-car-emergency"
                    block
                    :label="$t('common.detail', 2)"
                    @click="emit('selected')"
                />

                <UButton
                    v-if="canDo('TakeControl') && checkDispatchAccess(props.dispatch.jobs, CentrumAccessLevel.DISPATCH)"
                    class="truncate"
                    icon="i-mdi-account-multiple-plus"
                    variant="link"
                    block
                    :label="$t('common.assign')"
                    @click="
                        dispatchAssignModal.open({
                            dispatchId: props.dispatch.id,
                        })
                    "
                />

                <UButton
                    v-if="canDo('TakeDispatch') && checkDispatchAccess(props.dispatch.jobs, CentrumAccessLevel.PARTICIPATE)"
                    class="text-left"
                    icon="i-mdi-plus"
                    variant="link"
                    :label="$t('common.self_assign')"
                    @click="selfAssign(props.dispatch.id)"
                />
            </div>
        </template>

        <p class="inline-flex items-center gap-1">
            <span class="font-semibold">{{ $t('common.dispatch', 1) }}</span>
            <UButton class="font-semibold" :label="`DSP-${props.dispatch.id}`" variant="link" @click="emit('selected')" />
        </p>

        <ul role="list">
            <li class="inline-flex gap-1">
                <span class="flex-initial font-semibold">{{ $t('common.job') }}:</span>
                <span class="flex-1">
                    {{ props.dispatch.jobs?.jobs.map((j) => j.label ?? j.name).join(', ') }}
                </span>
            </li>

            <li>
                <span class="font-semibold">{{ $t('common.sent_at') }}:</span>
                {{ $d(toDate(props.dispatch.createdAt), 'short') }}
            </li>

            <li class="flex gap-1">
                <span class="flex-initial">
                    <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                </span>
                <span class="flex-1">
                    <template v-if="props.dispatch.anon">
                        {{ $t('common.anon') }}
                    </template>
                    <CitizenInfoPopover v-else-if="props.dispatch.creator" :user="props.dispatch.creator" size="sm" />
                    <template v-else>
                        {{ $t('common.unknown') }}
                    </template>
                </span>
            </li>

            <li>
                <span class="font-semibold">{{ $t('common.postal') }}:</span> {{ props.dispatch.postal ?? $t('common.na') }}
            </li>

            <li>
                <span class="font-semibold">{{ $t('common.message') }}:</span> {{ props.dispatch.message }}
            </li>

            <li class="truncate">
                <span class="font-semibold">{{ $t('common.description') }}:</span>
                {{ props.dispatch.description ?? $t('common.na') }}
            </li>

            <li class="flex gap-1">
                <span class="font-semibold">{{ $t('common.status') }}:</span>
                <DispatchStatusBadge :status="props.dispatch.status?.status" />
            </li>

            <li class="flex gap-1">
                <span class="font-semibold">{{ $t('common.attributes', 2) }}:</span>
                <DispatchAttributes :attributes="props.dispatch.attributes" size="xs" />
            </li>

            <li class="flex gap-1">
                <span class="font-semibold">{{ $t('common.unit') }}:</span>

                <span v-if="props.dispatch.units.length === 0" class="italic">{{
                    $t('enums.centrum.StatusDispatch.UNASSIGNED')
                }}</span>
                <span v-else class="grid grid-cols-2 gap-1">
                    <UnitInfoPopover
                        v-for="unit in props.dispatch.units"
                        :key="unit.unitId"
                        :unit-id="unit.unitId"
                        :unit="unit.unit"
                        :assignment="unit"
                        initials-only
                        badge
                        show-icon
                    />
                </span>
            </li>
        </ul>
    </UCard>
</template>
