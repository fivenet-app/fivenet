<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import {
    CalendarCheckIcon,
    CalendarRemoveIcon,
    CarBackIcon,
    CarEmergencyIcon,
    CheckBoldIcon,
    ChevronDownIcon,
    CoffeeIcon,
    HelpCircleIcon,
    HoopHouseIcon,
    InformationOutlineIcon,
    ListStatusIcon,
    MarkerCheckIcon,
} from 'mdi-vue3';
import { DefineComponent } from 'vue';
import { default as DispatchDetails } from '~/components/centrum/dispatches/Details.vue';
import { default as DispatchStatusUpdateModal } from '~/components/centrum/dispatches/StatusUpdateModal.vue';
import { dispatchStatusToBGColor, unitStatusToBGColor } from '~/components/centrum/helpers';
import { default as UnitDetails } from '~/components/centrum/units/Details.vue';
import { default as UnitStatusUpdateModal } from '~/components/centrum/units/StatusUpdateModal.vue';
import { useCentrumStore } from '~/store/centrum';
import { useNotificationsStore } from '~/store/notifications';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { CENTRUM_MODE } from '~~/gen/ts/resources/dispatch/settings';
import { UNIT_STATUS } from '~~/gen/ts/resources/dispatch/units';
import DispatchesLayer from './DispatchesLayer.vue';
import JoinUnit from './JoinUnitModal.vue';
import TakeDispatchModal from './TakeDispatchModal.vue';

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const centrumStore = useCentrumStore();
const { settings, ownDispatches, ownUnit, pendingDispatches } = storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

const notifications = useNotificationsStore();

const actionsUnit: {
    icon: DefineComponent;
    name: string;
    action?: Function;
    class?: string;
    status?: UNIT_STATUS;
}[] = [
    { icon: markRaw(CarBackIcon), name: 'Unavailable', status: UNIT_STATUS.UNAVAILABLE },
    { icon: markRaw(CalendarCheckIcon), name: 'Available', status: UNIT_STATUS.AVAILABLE },
    { icon: markRaw(CoffeeIcon), name: 'On Break', status: UNIT_STATUS.ON_BREAK },
    { icon: markRaw(CalendarRemoveIcon), name: 'Busy', status: UNIT_STATUS.BUSY },
    { icon: markRaw(ListStatusIcon), name: 'components.centrum.update_unit_status.title', class: 'bg-base-800' },
];

const actionsDispatch: {
    icon: DefineComponent;
    name: string;
    action?: Function;
    class?: string;
    status?: DISPATCH_STATUS;
}[] = [
    { icon: markRaw(CarBackIcon), name: 'En Route', status: DISPATCH_STATUS.EN_ROUTE },
    { icon: markRaw(MarkerCheckIcon), name: 'On Scene', status: DISPATCH_STATUS.ON_SCENE },
    { icon: markRaw(HelpCircleIcon), name: 'Need Assistance', status: DISPATCH_STATUS.NEED_ASSISTANCE },
    { icon: markRaw(CheckBoldIcon), name: 'Completed', status: DISPATCH_STATUS.COMPLETED },
    { icon: markRaw(ListStatusIcon), name: 'components.centrum.update_dispatch_status.title', class: 'bg-base-800' },
];

onBeforeMount(async () => setTimeout(async () => startStream(), 250));

onBeforeUnmount(async () => {
    stopStream();
    centrumStore.$reset();
});

const selectUnitOpen = ref(false);

const selectedDispatch = ref<Dispatch | undefined>();
const openDispatchDetails = ref(false);
const openDispatchStatus = ref(false);
const openTakeDispatch = ref(false);

const openUnitDetails = ref(false);
const openUnitStatus = ref(false);

const canStream = can('CentrumService.Stream');

async function updateDispatchStatus(dispatchId: bigint, status: DISPATCH_STATUS): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().updateDispatchStatus({
                dispatchId: dispatchId,
                status: status,
            });
            await call;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function updateDspStatus(dispatchId?: bigint, status?: DISPATCH_STATUS): Promise<void> {
    if (!dispatchId) {
        notifications.dispatchNotification({
            title: { key: 'notifications.centrum.sidebar.no_dispatch_selected.title', parameters: [] },
            content: { key: 'notifications.centrum.sidebar.no_dispatch_selected.content', parameters: [] },
            type: 'error',
        });
        return;
    }

    if (status === undefined) {
        openDispatchStatus.value = true;
        return;
    }

    updateDispatchStatus(dispatchId, status);
    notifications.dispatchNotification({
        title: { key: 'notifications.centrum.sidebar.dispatch_status_updated.title', parameters: [] },
        content: { key: 'notifications.centrum.sidebar.dispatch_status_updated.content', parameters: [] },
        type: 'success',
    });
}

async function updateUnitStatus(id: bigint, status: UNIT_STATUS): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().updateUnitStatus({
                unitId: id,
                status: status,
            });
            await call;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function updateUtStatus(id: bigint, status?: UNIT_STATUS): Promise<void> {
    if (status === undefined) {
        openUnitStatus.value = true;
        return;
    }

    updateUnitStatus(id, status);
    notifications.dispatchNotification({
        title: { key: 'notifications.centrum.sidebar.dispatch_status_updated.title', parameters: [] },
        content: { key: 'notifications.centrum.sidebar.dispatch_status_updated.content', parameters: [] },
        type: 'success',
    });
}
</script>

<template>
    <Livemap>
        <template v-slot:default v-if="canStream">
            <DispatchesLayer :show-all-dispatches="settings.mode === CENTRUM_MODE.SIMPLIFIED" @goto="$emit('goto', $event)" />
        </template>
        <template v-slot:afterMap v-if="canStream">
            <div class="lg:inset-y-0 lg:flex lg:w-50 lg:flex-col">
                <!-- Dispatch -->
                <TakeDispatchModal
                    v-if="ownUnit"
                    :open="openTakeDispatch"
                    @close="openTakeDispatch = false"
                    @goto="$emit('goto', $event)"
                />

                <template v-if="selectedDispatch">
                    <DispatchDetails
                        :dispatch="selectedDispatch"
                        :open="openDispatchDetails"
                        @close="openDispatchDetails = false"
                        @goto="$emit('goto', $event)"
                    />

                    <DispatchStatusUpdateModal
                        :dispatch="selectedDispatch"
                        :open="openDispatchStatus"
                        @close="openDispatchStatus = false"
                    />
                </template>

                <div class="h-full flex grow gap-y-5 overflow-y-auto bg-base-600 px-4 py-0.5">
                    <nav class="flex flex-1 flex-col">
                        <ul role="list" class="flex flex-1 flex-col gap-y-2 divide-y divide-base-400">
                            <li>
                                <ul role="list" class="-mx-2 mt-2 space-y-1">
                                    <li>
                                        <template v-if="ownUnit">
                                            <button
                                                @click="openUnitDetails = true"
                                                type="button"
                                                class="text-white hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                                :class="unitStatusToBGColor(ownUnit.status?.status)"
                                            >
                                                <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                                <span class="mt-2 truncate">{{ ownUnit.initials }}: {{ ownUnit.name }}</span>
                                                <span class="mt-2 truncate">
                                                    {{
                                                        $t(
                                                            `enums.centrum.UNIT_STATUS.${
                                                                UNIT_STATUS[ownUnit.status?.status ?? (0 as number)]
                                                            }`,
                                                        )
                                                    }}
                                                </span>
                                            </button>

                                            <UnitDetails
                                                :unit="ownUnit"
                                                :ownUnit="ownUnit"
                                                :open="openUnitDetails"
                                                @close="openUnitDetails = false"
                                            />
                                        </template>
                                        <button
                                            @click="selectUnitOpen = true"
                                            type="button"
                                            class="text-white bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                        >
                                            <template v-if="!ownUnit" class="flex w-full flex-col items-center">
                                                <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                                <span class="mt-2 truncate">{{ $t('common.no_own_unit') }}</span>
                                            </template>
                                            <template v-else class="truncate">{{ $t('common.leave_unit') }}</template>
                                        </button>

                                        <JoinUnit
                                            :open="selectUnitOpen"
                                            @close="selectUnitOpen = false"
                                            @joined="ownUnit = $event"
                                        />
                                    </li>
                                </ul>
                            </li>
                            <template v-if="ownUnit">
                                <li>
                                    <ul role="list" class="-mx-2 space-y-1">
                                        <li>
                                            <Disclosure as="div" v-slot="{ open }">
                                                <DisclosureButton
                                                    class="flex w-full items-start justify-between text-left text-white"
                                                >
                                                    <span class="text-base-100 leading-7">
                                                        <div class="text-xs font-semibold leading-6 text-base-100">
                                                            {{ $t('common.unit') }}
                                                        </div>
                                                    </span>
                                                    <span class="ml-6 flex h-7 items-center">
                                                        <ChevronDownIcon
                                                            :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                                            aria-hidden="true"
                                                        />
                                                    </span>
                                                </DisclosureButton>
                                                <DisclosurePanel>
                                                    <div class="flex flex-row gap-2">
                                                        <div class="w-full grid grid-cols-2 gap-0.5">
                                                            <UnitStatusUpdateModal
                                                                :unit="ownUnit"
                                                                :open="openUnitStatus"
                                                                @close="openUnitStatus = false"
                                                            />

                                                            <button
                                                                v-for="(item, idx) in actionsUnit"
                                                                :key="item.name"
                                                                type="button"
                                                                class="text-white bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                                                :class="[
                                                                    idx >= actionsUnit.length - 1 ? 'col-span-2' : '',
                                                                    item.status ? unitStatusToBGColor(item.status) : item.class,
                                                                    item.class,
                                                                ]"
                                                                @click="updateUtStatus(ownUnit.id, item.status)"
                                                            >
                                                                <component
                                                                    :is="item.icon ?? HoopHouseIcon"
                                                                    class="text-base-100 group-hover:text-white h-5 w-5 shrink-0"
                                                                    aria-hidden="true"
                                                                />
                                                                <span class="mt-1">
                                                                    {{
                                                                        item.status
                                                                            ? $t(
                                                                                  `enums.centrum.UNIT_STATUS.${
                                                                                      UNIT_STATUS[item.status ?? (0 as number)]
                                                                                  }`,
                                                                              )
                                                                            : $t(item.name)
                                                                    }}
                                                                </span>
                                                            </button>
                                                        </div>
                                                    </div>
                                                </DisclosurePanel>
                                            </Disclosure>
                                        </li>
                                    </ul>
                                </li>
                                <li>
                                    <ul role="list" class="-mx-2 space-y-1">
                                        <div class="text-xs font-semibold leading-6 text-base-100">
                                            {{ $t('common.dispatch') }} {{ $t('common.status') }}
                                        </div>
                                        <li>
                                            <div class="grid grid-cols-2 gap-0.5">
                                                <button
                                                    v-for="(item, idx) in actionsDispatch"
                                                    :key="item.name"
                                                    type="button"
                                                    class="text-white bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                                    :class="[
                                                        idx >= actionsDispatch.length - 1 ? 'col-span-2' : '',
                                                        item.status ? dispatchStatusToBGColor(item.status) : item.class,
                                                        item.class,
                                                    ]"
                                                    @click="updateDspStatus(selectedDispatch?.id, item.status)"
                                                >
                                                    <component
                                                        :is="item.icon ?? HoopHouseIcon"
                                                        class="text-base-100 group-hover:text-white h-5 w-5 shrink-0"
                                                        aria-hidden="true"
                                                    />
                                                    <span class="mt-1">
                                                        {{
                                                            item.status
                                                                ? $t(
                                                                      `enums.centrum.DISPATCH_STATUS.${
                                                                          DISPATCH_STATUS[item.status ?? (0 as number)]
                                                                      }`,
                                                                  )
                                                                : $t(item.name)
                                                        }}
                                                    </span>
                                                </button>
                                            </div>
                                        </li>
                                    </ul>
                                </li>
                                <li>
                                    <div class="text-xs font-semibold leading-6 text-base-100">
                                        {{ $t('common.your_dispatches') }}
                                    </div>
                                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                                        <li v-if="ownDispatches.length === 0">
                                            <button
                                                type="button"
                                                class="text-white bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                            >
                                                <CarEmergencyIcon class="h-5 w-5" aria-hidden="true" />
                                                <span class="mt-2 truncate">{{ $t('common.no_assigned_dispatches') }}</span>
                                            </button>
                                        </li>
                                        <template v-else>
                                            <li class="flex flex-row items-center" v-for="dispatch in ownDispatches">
                                                <div class="mr-1.5">
                                                    <input
                                                        name="active"
                                                        type="radio"
                                                        class="h-4 w-4 border-gray-300 text-primary-600 focus:ring-primary-600"
                                                        v-model="selectedDispatch"
                                                        v-bind:value="dispatch"
                                                    />
                                                </div>
                                                <button
                                                    type="button"
                                                    class="text-white bg-error-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                                    @click="
                                                        selectedDispatch = dispatch;
                                                        openDispatchDetails = true;
                                                    "
                                                >
                                                    <span class="font-bold truncate inline-flex items-center">
                                                        <CarEmergencyIcon class="h-4 w-4 mr-0.5" /> DSP-{{ dispatch.id }}</span
                                                    >
                                                    <span class="mt-2 truncate">
                                                        {{
                                                            $t(
                                                                `enums.centrum.DISPATCH_STATUS.${
                                                                    DISPATCH_STATUS[dispatch.status?.status ?? (0 as number)]
                                                                }`,
                                                            )
                                                        }}
                                                    </span>
                                                </button>
                                            </li>
                                        </template>
                                    </ul>
                                </li>
                            </template>
                        </ul>
                    </nav>
                </div>

                <!-- "Take Dispatches" Button -->
                <span v-if="ownUnit" class="fixed inline-flex z-90 bottom-2 right-1/2">
                    <span class="flex absolute h-3 w-3 top-0 right-0 -mt-1 -mr-1" v-if="pendingDispatches.length > 0">
                        <span
                            class="animate-ping absolute inline-flex h-full w-full rounded-full bg-error-400 opacity-75"
                        ></span>
                        <span class="relative inline-flex rounded-full h-3 w-3 bg-error-500"></span>
                    </span>
                    <button
                        type="button"
                        @click="openTakeDispatch = true"
                        class="flex items-center justify-center w-12 h-12 rounded-full bg-primary-500 shadow-float text-neutral hover:bg-primary-400"
                    >
                        <CarEmergencyIcon class="w-10 h-auto" />
                    </button>
                </span>
            </div>
        </template>
    </Livemap>
</template>
