<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useConfirmDialog } from '@vueuse/core';
import { AccountIcon, CloseIcon, MapMarkerIcon, PencilIcon, PlusIcon, TrashCanIcon } from 'mdi-vue3';
import AssignDispatchModal from '~/components/centrum/dispatches//AssignDispatchModal.vue';
import DispatchFeed from '~/components/centrum/dispatches/DispatchFeed.vue';
import DispatchStatusUpdateModal from '~/components/centrum/dispatches/DispatchStatusUpdateModal.vue';
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import Time from '~/components/partials/elements/Time.vue';
import { useCentrumStore } from '~/store/centrum';
import { useNotificatorStore } from '~/store/notificator';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { Settings } from '~~/gen/ts/resources/dispatch/settings';
import { TakeDispatchResp } from '~~/gen/ts/services/centrum/centrum';

const props = defineProps<{
    open: boolean;
    dispatch: Dispatch;
    settings?: Settings;
}>();

defineEmits<{
    (e: 'close'): void;
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const centrumStore = useCentrumStore();
const { ownUnitId } = storeToRefs(centrumStore);
const { canDo } = centrumStore;

const notificationsStore = useNotificatorStore();

async function selfAssign(id: bigint): Promise<void> {
    if (ownUnitId.value === undefined) {
        notificationsStore.dispatchNotification({
            title: { key: 'notifications.centrum.unitUpdated.not_in_unit.title' },
            content: { key: 'notifications.centrum.unitUpdated.not_in_unit.content' },
            type: 'error',
        });

        return;
    }

    try {
        const call = $grpc.getCentrumClient().takeDispatch({
            dispatchIds: [id],
            resp: TakeDispatchResp.ACCEPTED,
        });
        await call;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteDispatch(id: bigint): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().deleteDispatch({ id });
        await call;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => await deleteDispatch(id));

const dispatchStatusColors = computed(() => dispatchStatusToBGColor(props.dispatch.status?.status ?? 0));

const openAssign = ref(false);
const openStatus = ref(false);
</script>

<template>
    <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(dispatch.id)" />

    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
            <div class="fixed inset-0" />

            <div class="fixed inset-0 overflow-hidden">
                <div class="absolute inset-0 overflow-hidden">
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-2xl pl-10 sm:pl-16">
                        <TransitionChild
                            as="template"
                            enter="transform transition ease-in-out duration-150 sm:duration-300"
                            enter-from="translate-x-full"
                            enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-150 sm:duration-300"
                            leave-from="translate-x-0"
                            leave-to="translate-x-full"
                        >
                            <DialogPanel class="pointer-events-auto w-screen max-w-2xl">
                                <form class="flex h-full flex-col divide-y divide-gray-200 bg-gray-900 shadow-xl">
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle
                                                    class="inline-flex items-center text-base font-semibold leading-6 text-neutral"
                                                >
                                                    {{ $t('common.dispatch') }}:
                                                    <IDCopyBadge :id="dispatch.id" class="ml-2 mr-2" prefix="DSP" />
                                                    <p class="truncate max-w-[20rem]" :title="dispatch.message">
                                                        {{ dispatch.message }}
                                                    </p>
                                                    <button
                                                        v-if="can('CentrumService.DeleteDispatch')"
                                                        type="button"
                                                        class="flex-initial text-neutral hover:text-gray-300 ml-1"
                                                        :title="$t('common.delete')"
                                                        @click="reveal()"
                                                    >
                                                        <TrashCanIcon class="h-6 w-6" aria-hidden="true" />
                                                    </button>
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-neutral"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                                        <CloseIcon class="h-6 w-6" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                            <div class="mt-1">
                                                <p class="text-sm text-primary-300">
                                                    {{ $t('common.description') }}:
                                                    {{ dispatch.description ?? $t('common.na') }}
                                                </p>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="divide-y divide-neutral/10">
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.created_at') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <Time :value="dispatch.createdAt" />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.sent_by') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <span v-if="dispatch.anon">
                                                                    {{ $t('common.anon') }}
                                                                </span>
                                                                <CitizenInfoPopover
                                                                    v-else-if="dispatch.creator"
                                                                    :user="dispatch.creator"
                                                                />
                                                                <span v-else>
                                                                    {{ $t('common.unknown') }}
                                                                </span>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.location') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <span class="block">
                                                                    {{ $t('common.postal') }}:
                                                                    {{ dispatch.postal ?? $t('common.na') }}
                                                                </span>
                                                                <button
                                                                    v-if="dispatch.x && dispatch.y"
                                                                    type="button"
                                                                    class="inline-flex items-center text-primary-400 hover:text-primary-600"
                                                                    @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
                                                                >
                                                                    <MapMarkerIcon class="w-5 h-5 mr-1" aria-hidden="true" />
                                                                    {{ $t('common.go_to_location') }}
                                                                </button>
                                                                <span v-else>{{ $t('common.no_location') }}</span>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.attributes', 2) }}
                                                            </dt>
                                                            <dd class="mt-2 text-sm text-gray-400 sm:col-span-2 sm:mt-0">
                                                                <template
                                                                    v-if="
                                                                        dispatch.attributes !== undefined &&
                                                                        dispatch.attributes?.list.length > 0
                                                                    "
                                                                >
                                                                    <span
                                                                        v-for="attribute in dispatch.attributes?.list"
                                                                        :key="attribute"
                                                                        class="inline-flex items-center rounded-md bg-error-400/10 px-2 py-1 text-xs font-medium text-error-400 ring-1 ring-inset ring-error-400/20"
                                                                    >
                                                                        {{
                                                                            $t(
                                                                                `components.centrum.update_dispatch.attributes.${attribute}`,
                                                                            )
                                                                        }}
                                                                    </span>
                                                                </template>
                                                                <span v-else>
                                                                    {{
                                                                        $t('common.none_selected', [$t('common.attributes', 2)])
                                                                    }}
                                                                </span>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.units') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <span class="block">
                                                                    {{ $t('common.unit', dispatch.units.length) }}
                                                                </span>
                                                                <ul
                                                                    role="list"
                                                                    class="border divide-y rounded-md divide-base-200 border-base-200"
                                                                >
                                                                    <li
                                                                        v-for="unit in dispatch.units"
                                                                        :key="unit.unitId.toString()"
                                                                        class="flex items-center justify-between py-3 pl-3 pr-4 text-sm"
                                                                    >
                                                                        <div class="flex items-center flex-1">
                                                                            <UnitInfoPopover
                                                                                :unit="unit.unit"
                                                                                :assignment="unit"
                                                                                class="flex items-center justify-center"
                                                                            >
                                                                                <template #before>
                                                                                    <AccountIcon
                                                                                        class="flex-shrink-0 w-5 h-5 text-base-400 mr-1"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </template>
                                                                            </UnitInfoPopover>
                                                                            <span
                                                                                v-if="unit.expiresAt"
                                                                                class="flex-1 ml-2 truncate inline-flex items-center"
                                                                            >
                                                                                -
                                                                                {{
                                                                                    useLocaleTimeAgo(toDate(unit.expiresAt), {
                                                                                        showSecond: true,
                                                                                        updateInterval: 1000,
                                                                                    }).value
                                                                                }}
                                                                            </span>
                                                                        </div>
                                                                    </li>
                                                                </ul>

                                                                <AssignDispatchModal
                                                                    v-if="openAssign"
                                                                    :open="openAssign"
                                                                    :dispatch="dispatch"
                                                                    @close="openAssign = false"
                                                                />

                                                                <span class="mt-2 isolate inline-flex rounded-md shadow-sm">
                                                                    <button
                                                                        v-if="canDo('TakeControl')"
                                                                        type="button"
                                                                        class="flex flex-row items-center rounded bg-neutral/10 px-2 py-1 text-xs font-semibold text-neutral shadow-sm hover:bg-neutral/20"
                                                                        @click="openAssign = true"
                                                                    >
                                                                        <PencilIcon class="h-6 w-6" />
                                                                        <span class="truncate ml-0.5">
                                                                            {{ $t('common.assign') }}
                                                                        </span>
                                                                    </button>
                                                                    <button
                                                                        v-if="canDo('TakeDispatch')"
                                                                        type="button"
                                                                        class="flex flex-row items-center ml-2 rounded bg-neutral/10 px-2 py-1 text-xs font-semibold text-neutral shadow-sm hover:bg-neutral/20"
                                                                        @click="selfAssign(dispatch.id)"
                                                                    >
                                                                        <PlusIcon class="h-6 w-6" />
                                                                        <span class="truncate ml-0.5">{{
                                                                            $t('common.self_assign')
                                                                        }}</span>
                                                                    </button>
                                                                </span>
                                                            </dd>
                                                        </div>
                                                    </dl>
                                                </div>
                                                <div class="mt-1">
                                                    <dl class="divide-y divide-neutral/10">
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.last_update') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <Time :value="dispatch.status?.createdAt" />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.location') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <span class="block">
                                                                    {{ $t('common.postal') }}:
                                                                    {{ dispatch.status?.postal ?? $t('common.na') }}
                                                                </span>
                                                                <button
                                                                    v-if="dispatch.status?.x && dispatch.status?.y"
                                                                    type="button"
                                                                    class="inline-flex items-center text-primary-400 hover:text-primary-600"
                                                                    @click="
                                                                        $emit('goto', {
                                                                            x: dispatch.status?.x,
                                                                            y: dispatch.status?.y,
                                                                        })
                                                                    "
                                                                >
                                                                    <MapMarkerIcon class="w-5 h-5 mr-1" aria-hidden="true" />
                                                                    {{ $t('common.go_to_location') }}
                                                                </button>
                                                                <span v-else>{{ $t('common.no_location') }}</span>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.status') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <DispatchStatusUpdateModal
                                                                    v-if="openStatus"
                                                                    :open="openStatus"
                                                                    :dispatch-id="dispatch.id"
                                                                    @close="openStatus = false"
                                                                />

                                                                <button
                                                                    type="button"
                                                                    class="rounded px-2 py-1 text-xs font-semibold text-neutral shadow-sm hover:bg-neutral/20"
                                                                    :class="dispatchStatusColors"
                                                                    @click="openStatus = true"
                                                                >
                                                                    {{
                                                                        $t(
                                                                            `enums.centrum.StatusDispatch.${
                                                                                StatusDispatch[dispatch.status?.status ?? 0]
                                                                            }`,
                                                                        )
                                                                    }}
                                                                    <span v-if="dispatch.status?.code">
                                                                        ({{ $t('common.code') }}: '{{ dispatch.status.code }}')
                                                                    </span>
                                                                </button>
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.code') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                {{ dispatch.status?.code ?? $t('common.na') }}
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                {{ $t('common.reason') }}
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                {{ dispatch.status?.reason ?? $t('common.na') }}
                                                            </dd>
                                                        </div>
                                                    </dl>
                                                </div>

                                                <DispatchFeed :dispatch-id="dispatch.id" />
                                            </div>
                                        </div>
                                    </div>
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <button
                                            type="button"
                                            class="w-full rounded-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
                                            @click="$emit('close')"
                                        >
                                            {{ $t('common.close') }}
                                        </button>
                                    </div>
                                </form>
                            </DialogPanel>
                        </TransitionChild>
                    </div>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
