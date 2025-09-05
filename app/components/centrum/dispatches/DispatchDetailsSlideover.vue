<script lang="ts" setup>
import DispatchAssignModal from '~/components/centrum/dispatches//DispatchAssignModal.vue';
import DispatchFeed from '~/components/centrum/dispatches/DispatchFeed.vue';
import DispatchStatusUpdateModal from '~/components/centrum/dispatches/DispatchStatusUpdateModal.vue';
import { checkDispatchAccess, dispatchStatusToBGColor, dispatchStatusToIcon } from '~/components/centrum/helpers';
import DispatchAttributes from '~/components/centrum/partials/DispatchAttributes.vue';
import DispatchReferences from '~/components/centrum/partials/DispatchReferences.vue';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import { CentrumAccessLevel } from '~~/gen/ts/resources/centrum/access';
import { type Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

const props = defineProps<{
    dispatchId: number;
    dispatch?: Dispatch;
}>();

const emit = defineEmits<{
    close: [boolean];
}>();

const { can } = useAuth();

const overlay = useOverlay();

const { goto } = useLivemapStore();

const centrumStore = useCentrumStore();
const { dispatches, timeCorrection } = storeToRefs(centrumStore);
const { canDo, selfAssign } = centrumStore;

const centrumCentrumClient = await getCentrumCentrumClient();

const dispatch = computed(() => (props.dispatch ? props.dispatch : dispatches.value.get(props.dispatchId)!));

async function deleteDispatch(id: number): Promise<void> {
    try {
        const call = centrumCentrumClient.deleteDispatch({ id });
        await call;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const dispatchStatusColors = computed(() => dispatchStatusToBGColor(dispatch.value?.status?.status));

watch(dispatch, () => {
    if (dispatch.value === undefined) {
        emit('close', false);
    }
});

const canAccessDispatch = computed(() => ({
    participate: checkDispatchAccess(dispatch.value?.jobs, CentrumAccessLevel.PARTICIPATE),
    dispatch: checkDispatchAccess(dispatch.value?.jobs, CentrumAccessLevel.DISPATCH),
}));

const confirmModal = overlay.create(ConfirmModal);
const dispatchAssignModal = overlay.create(DispatchAssignModal);
const dispatchStatusUpdateModal = overlay.create(DispatchStatusUpdateModal);
</script>

<template>
    <USlideover :overlay="false">
        <template #title>
            <div class="inline-flex items-center">
                <IDCopyBadge :id="dispatch.id" class="mx-2" prefix="DSP" />

                <p class="max-w-80 flex-1 truncate">
                    {{ dispatch.message }}
                </p>
            </div>
        </template>

        <template #body>
            <div class="divide-y divide-default">
                <div>
                    <dl class="divide-neutral/10 divide-y">
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.job') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <span>
                                    {{ dispatch.jobs?.jobs.map((j) => j.label ?? j.name).join(', ') }}
                                </span>
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.sent_at') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <GenericTime :value="dispatch.createdAt" />
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.sent_by') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <span v-if="dispatch.anon">
                                    {{ $t('common.anon') }}
                                </span>
                                <CitizenInfoPopover v-else-if="dispatch.creator" :user="dispatch.creator" />
                                <span v-else>
                                    {{ $t('common.unknown') }}
                                </span>
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.location') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <div class="sm:inline-flex sm:flex-row sm:gap-2">
                                    <span class="block">
                                        {{ $t('common.postal') }}:
                                        {{ dispatch.postal ?? $t('common.na') }}
                                    </span>
                                    <UButton
                                        size="xs"
                                        variant="link"
                                        icon="i-mdi-map-marker"
                                        @click="goto({ x: dispatch.x, y: dispatch.y })"
                                    >
                                        {{ $t('common.go_to_location') }}
                                    </UButton>
                                </div>
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.description') }}
                            </dt>
                            <dd class="mt-2 text-sm sm:col-span-2 sm:mt-0">
                                <p class="max-h-14 overflow-y-scroll break-words">
                                    {{ dispatch.description ?? $t('common.na') }}
                                </p>
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.attributes', 2) }}
                            </dt>
                            <dd class="mt-2 text-sm sm:col-span-2 sm:mt-0">
                                <DispatchAttributes :attributes="dispatch.attributes" />
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.reference', 2) }}
                            </dt>
                            <dd class="mt-2 text-sm sm:col-span-2 sm:mt-0">
                                <DispatchReferences :references="dispatch.references" />
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.unit', 2) }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <span v-if="dispatch.units.length === 0" class="block">
                                    {{ $t('common.units', dispatch.units.length) }}
                                </span>
                                <div v-else class="mb-1 rounded-md bg-neutral-100 dark:bg-neutral-900">
                                    <ul class="divide-y divide-gray-100 text-sm font-medium dark:divide-gray-800" role="list">
                                        <li
                                            v-for="unit in dispatch.units"
                                            :key="unit.unitId"
                                            class="flex items-center justify-between py-3 pr-4 pl-3"
                                        >
                                            <div class="flex flex-1 items-center">
                                                <UnitInfoPopover
                                                    class="flex items-center justify-center"
                                                    :unit="unit.unit"
                                                    :assignment="unit"
                                                >
                                                    <template #before>
                                                        <UIcon class="mr-1 size-5 shrink-0" name="i-mdi-account-group" />
                                                    </template>
                                                </UnitInfoPopover>
                                                <span
                                                    v-if="unit.expiresAt"
                                                    class="ml-2 inline-flex flex-1 items-center truncate"
                                                >
                                                    -
                                                    {{
                                                        useLocaleTimeAgo(toDate(unit.expiresAt, timeCorrection), {
                                                            showSecond: true,
                                                            updateInterval: 1_000,
                                                        }).value
                                                    }}
                                                </span>
                                            </div>
                                        </li>
                                    </ul>
                                </div>

                                <UButtonGroup class="inline-flex w-full">
                                    <UButton
                                        v-if="canDo('TakeControl') && canAccessDispatch.dispatch"
                                        icon="i-mdi-account-multiple-plus"
                                        truncate
                                        @click="dispatchAssignModal.open({ dispatchId: dispatchId })"
                                    >
                                        {{ $t('common.assign') }}
                                    </UButton>
                                    <UButton
                                        v-if="canDo('TakeDispatch') && canAccessDispatch.participate"
                                        icon="i-mdi-plus"
                                        truncate
                                        @click="selfAssign(dispatch.id)"
                                    >
                                        {{ $t('common.self_assign') }}
                                    </UButton>
                                </UButtonGroup>
                            </dd>
                        </div>
                    </dl>
                </div>

                <div>
                    <dl class="divide-neutral/10 divide-y">
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.last_update') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <GenericTime :value="dispatch.status?.createdAt" />
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.location') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <div class="sm:inline-flex sm:flex-row sm:gap-2">
                                    <span class="block">
                                        {{ $t('common.postal') }}:
                                        {{ dispatch.status?.postal ?? $t('common.na') }}
                                    </span>
                                    <UButton
                                        v-if="dispatch.status?.x !== undefined && dispatch.status?.y !== undefined"
                                        size="xs"
                                        variant="link"
                                        icon="i-mdi-map-marker"
                                        @click="
                                            goto({
                                                x: dispatch.status?.x,
                                                y: dispatch.status?.y,
                                            })
                                        "
                                    >
                                        {{ $t('common.go_to_location') }}
                                    </UButton>
                                    <span v-else>{{ $t('common.no_location') }}</span>
                                </div>
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.status') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UButton
                                    class="rounded-sm px-2 py-1 text-sm font-semibold"
                                    :class="dispatchStatusColors"
                                    :icon="dispatchStatusToIcon(dispatch.status?.status)"
                                    :disabled="!canAccessDispatch.participate"
                                    @click="dispatchStatusUpdateModal.open({ dispatchId: dispatch.id })"
                                >
                                    {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`) }}
                                    <span v-if="dispatch.status?.code">
                                        ({{ $t('common.code') }}: '{{ dispatch.status.code }}')
                                    </span>
                                </UButton>
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.code') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                {{ dispatch.status?.code ?? $t('common.na') }}
                            </dd>
                        </div>

                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm leading-6 font-medium">
                                {{ $t('common.reason') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                {{ dispatch.status?.reason ?? $t('common.na') }}
                            </dd>
                        </div>
                    </dl>
                </div>

                <div>
                    <DispatchFeed :dispatch-id="dispatch.id" />
                </div>
            </div>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton color="neutral" class="flex-1" block @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>

                <UTooltip
                    v-if="can('centrum.CentrumService/DeleteDispatch').value && canAccessDispatch.dispatch"
                    :text="$t('common.delete')"
                >
                    <UButton
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            confirmModal.open({
                                confirm: async () => dispatch && deleteDispatch(dispatch.id),
                            })
                        "
                    />
                </UTooltip>
            </UButtonGroup>
        </template>
    </USlideover>
</template>
