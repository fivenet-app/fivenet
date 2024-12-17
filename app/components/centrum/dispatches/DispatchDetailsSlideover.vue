<script lang="ts" setup>
import DispatchAssignModal from '~/components/centrum/dispatches//DispatchAssignModal.vue';
import DispatchFeed from '~/components/centrum/dispatches/DispatchFeed.vue';
import DispatchStatusUpdateModal from '~/components/centrum/dispatches/DispatchStatusUpdateModal.vue';
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import DispatchAttributes from '~/components/centrum/partials/DispatchAttributes.vue';
import DispatchReferences from '~/components/centrum/partials/DispatchReferences.vue';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import { useNotificatorStore } from '~/store/notificator';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { StatusDispatch, TakeDispatchResp } from '~~/gen/ts/resources/centrum/dispatches';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    dispatchId: string;
    dispatch?: Dispatch;
}>();

const { can } = useAuth();

const modal = useModal();

const { isOpen } = useSlideover();

const { goto } = useLivemapStore();

const centrumStore = useCentrumStore();
const { dispatches, ownUnitId, timeCorrection } = storeToRefs(centrumStore);
const { canDo } = centrumStore;

const notifications = useNotificatorStore();

const dispatch = computed(() => (props.dispatch ? props.dispatch : dispatches.value.get(props.dispatchId)));

async function selfAssign(id: string): Promise<void> {
    if (ownUnitId.value === undefined) {
        notifications.add({
            title: { key: 'notifications.centrum.unitUpdated.not_in_unit.title' },
            description: { key: 'notifications.centrum.unitUpdated.not_in_unit.content' },
            type: NotificationType.ERROR,
        });

        return;
    }

    try {
        const call = getGRPCCentrumClient().takeDispatch({
            dispatchIds: [id],
            resp: TakeDispatchResp.ACCEPTED,
        });
        await call;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteDispatch(id: string): Promise<void> {
    try {
        const call = getGRPCCentrumClient().deleteDispatch({ id });
        await call;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const dispatchStatusColors = computed(() => dispatchStatusToBGColor(dispatch.value?.status?.status));

watch(dispatch, () => {
    if (dispatch.value === undefined) {
        isOpen.value = false;
    }
});
</script>

<template>
    <USlideover :ui="{ width: 'w-screen max-w-xl' }" :overlay="false">
        <UCard
            v-if="dispatch"
            :ui="{
                body: {
                    base: 'flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <div class="inline-flex items-center">
                        <IDCopyBadge :id="dispatch.id" class="mx-2" prefix="DSP" />

                        <p class="max-w-80 truncate" :title="dispatch.message">
                            {{ dispatch.message }}
                        </p>

                        <UButton
                            v-if="can('CentrumService.DeleteDispatch').value"
                            variant="link"
                            icon="i-mdi-trash-can"
                            color="red"
                            :title="$t('common.delete')"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => dispatch && deleteDispatch(dispatch.id),
                                })
                            "
                        />
                    </div>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div class="divide-y divide-gray-100 dark:divide-gray-800">
                <div>
                    <dl class="divide-neutral/10 divide-y">
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                {{ $t('common.created_at') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <GenericTime :value="dispatch.createdAt" />
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
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
                            <dt class="text-sm font-medium leading-6">
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
                                        :padded="false"
                                        @click="goto({ x: dispatch.x, y: dispatch.y })"
                                    >
                                        {{ $t('common.go_to_location') }}
                                    </UButton>
                                </div>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                {{ $t('common.description') }}
                            </dt>
                            <dd class="mt-2 text-sm sm:col-span-2 sm:mt-0">
                                <p class="max-h-14 overflow-y-scroll break-words">
                                    {{ dispatch.description ?? $t('common.na') }}
                                </p>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                {{ $t('common.attributes', 2) }}
                            </dt>
                            <dd class="mt-2 text-sm sm:col-span-2 sm:mt-0">
                                <DispatchAttributes :attributes="dispatch.attributes" />
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                {{ $t('common.reference', 2) }}
                            </dt>
                            <dd class="mt-2 text-sm sm:col-span-2 sm:mt-0">
                                <DispatchReferences :references="dispatch.references" />
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                {{ $t('common.unit') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <span v-if="dispatch.units.length === 0" class="block">
                                    {{ $t('common.units', dispatch.units.length) }}
                                </span>
                                <div v-else class="mb-1 rounded-md bg-neutral-100 dark:bg-base-900">
                                    <ul role="list" class="divide-y divide-gray-100 text-sm font-medium dark:divide-gray-800">
                                        <li
                                            v-for="unit in dispatch.units"
                                            :key="unit.unitId"
                                            class="flex items-center justify-between py-3 pl-3 pr-4"
                                        >
                                            <div class="flex flex-1 items-center">
                                                <UnitInfoPopover
                                                    :unit="unit.unit"
                                                    :assignment="unit"
                                                    class="flex items-center justify-center"
                                                >
                                                    <template #before>
                                                        <UIcon name="i-mdi-account-group" class="mr-1 size-5 shrink-0" />
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

                                <UButtonGroup class="inline-flex">
                                    <UButton
                                        v-if="canDo('TakeControl')"
                                        icon="i-mdi-account-multiple-plus"
                                        truncate
                                        @click="modal.open(DispatchAssignModal, { dispatchId: dispatchId })"
                                    >
                                        {{ $t('common.assign') }}
                                    </UButton>
                                    <UButton
                                        v-if="canDo('TakeDispatch')"
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
                            <dt class="text-sm font-medium leading-6">
                                {{ $t('common.last_update') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <GenericTime :value="dispatch.status?.createdAt" />
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
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
                                        :padded="false"
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
                            <dt class="text-sm font-medium leading-6">
                                {{ $t('common.status') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                <UButton
                                    class="rounded px-2 py-1 text-sm font-semibold"
                                    :class="dispatchStatusColors"
                                    @click="modal.open(DispatchStatusUpdateModal, { dispatchId: dispatch.id })"
                                >
                                    {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`) }}
                                    <span v-if="dispatch.status?.code">
                                        ({{ $t('common.code') }}: '{{ dispatch.status.code }}')
                                    </span>
                                </UButton>
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                {{ $t('common.code') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                {{ dispatch.status?.code ?? $t('common.na') }}
                            </dd>
                        </div>
                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                            <dt class="text-sm font-medium leading-6">
                                {{ $t('common.reason') }}
                            </dt>
                            <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                                {{ dispatch.status?.reason ?? $t('common.na') }}
                            </dd>
                        </div>
                    </dl>
                </div>

                <DispatchFeed v-if="isOpen" :dispatch-id="dispatch.id" />
            </div>

            <template #footer>
                <UButton color="black" block @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </USlideover>
</template>
