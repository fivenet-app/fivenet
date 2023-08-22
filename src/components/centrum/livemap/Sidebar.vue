<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
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
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import { default as UnitDetails } from '~/components/centrum/units/Details.vue';
import { default as UnitStatusUpdateModal } from '~/components/centrum/units/StatusUpdateModal.vue';
import { useCentrumStore } from '~/store/centrum';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { UNIT_STATUS } from '~~/gen/ts/resources/dispatch/units';
import AssignDispatchModal from '../dispatches/AssignDispatchModal.vue';
import AssignUnitModal from '../units/AssignUnitModal.vue';
import DispatchesLayer from './DispatchesLayer.vue';
import DispatchEntry from './sidebar/DispatchEntry.vue';
import JoinUnit from './sidebar/JoinUnitModal.vue';
import TakeDispatch from './sidebar/TakeDispatchModal.vue';

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const centrumStore = useCentrumStore();
const { ownDispatches, ownUnit, pendingDispatches } = storeToRefs(centrumStore);
const { startStream, stopStream } = centrumStore;

const actionsDispatch: {
    icon: DefineComponent;
    name: string;
    action?: Function;
    class?: string;
    status?: DISPATCH_STATUS;
}[] = [
    { icon: markRaw(CarBackIcon), name: 'En Route', class: 'bg-info-600', status: DISPATCH_STATUS.EN_ROUTE },
    { icon: markRaw(MarkerCheckIcon), name: 'On Scene', class: 'bg-primary-600', status: DISPATCH_STATUS.ON_SCENE },
    { icon: markRaw(HelpCircleIcon), name: 'Need Assistance', class: 'bg-warn-600', status: DISPATCH_STATUS.NEED_ASSISTANCE },
    { icon: markRaw(CheckBoldIcon), name: 'Completed', class: 'bg-success-600', status: DISPATCH_STATUS.COMPLETED },
    { icon: markRaw(ListStatusIcon), name: 'Update Status', class: 'bg-base-800' },
];

const actionsUnit: {
    icon: DefineComponent;
    name: string;
    action?: Function;
    class?: string;
    status?: UNIT_STATUS;
}[] = [
    { icon: markRaw(CarBackIcon), name: 'Unavailable', class: 'bg-error-600', status: UNIT_STATUS.UNAVAILABLE },
    { icon: markRaw(CalendarCheckIcon), name: 'Available', class: 'bg-success-600', status: UNIT_STATUS.AVAILABLE },
    { icon: markRaw(CoffeeIcon), name: 'On Break', class: 'bg-warn-600', status: UNIT_STATUS.ON_BREAK },
    { icon: markRaw(CalendarRemoveIcon), name: 'Busy', class: 'bg-info-600', status: UNIT_STATUS.BUSY },
    { icon: markRaw(ListStatusIcon), name: 'Update Status', class: 'bg-base-800' },
];

onMounted(() => {
    startStream();
});

onBeforeUnmount(() => {
    stopStream();
});

const selectUnitOpen = ref(false);

const unitStatusOpen = ref(false);
const unitStatusSelected = ref<UNIT_STATUS | undefined>();

// TODO add function to set the active dispatch (last clicked and manually selected)
const dispatchStatusSelected = ref<DISPATCH_STATUS | undefined>();

const selectedDispatch = ref<Dispatch | undefined>();
const openDispatchDetails = ref(false);
const openDispatchAssign = ref(false);
const openDispatchStatus = ref(false);
const openTakeDispatch = ref(false);

const openUnitDetails = ref(false);
const openUnitAssign = ref(false);
const openUnitStatus = ref(false);

const canStream = can('CentrumService.Stream');
</script>

<template>
    <Livemap>
        <template v-slot:default v-if="canStream">
            <DispatchesLayer
                @select="
                    selectedDispatch = $event;
                    openDispatchDetails = true;
                "
            />
        </template>
        <template v-slot:afterMap v-if="canStream">
            <div class="lg:inset-y-0 lg:flex lg:w-50 lg:flex-col">
                <!-- Own Unit Modals -->
                <template v-if="ownUnit">
                    <UnitStatusUpdateModal
                        :open="unitStatusOpen"
                        @close="unitStatusOpen = false"
                        :unit="ownUnit"
                        :status="unitStatusSelected"
                    />
                </template>

                <!-- Dispatch -->
                <TakeDispatch :open="openTakeDispatch" @close="openTakeDispatch = false" @goto="$emit('goto', $event)" />

                <template v-if="selectedDispatch">
                    <DispatchDetails
                        :dispatch="selectedDispatch"
                        :open="openDispatchDetails"
                        @close="openDispatchDetails = false"
                        @goto="$emit('goto', $event)"
                        @assign-unit="
                            selectedDispatch = $event;
                            openDispatchAssign = true;
                        "
                        @status="
                            selectedDispatch = $event;
                            openDispatchStatus = true;
                        "
                    />
                    <AssignDispatchModal
                        :open="openDispatchAssign"
                        :dispatch="selectedDispatch"
                        @close="openDispatchAssign = false"
                    />
                    <DispatchStatusUpdateModal
                        :open="openDispatchStatus"
                        :dispatch="selectedDispatch"
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
                                                class="text-white bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
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
                                            <AssignUnitModal
                                                :open="openUnitAssign"
                                                :unit="ownUnit"
                                                @close="openUnitAssign = false"
                                            />
                                            <UnitStatusUpdateModal
                                                :open="openUnitStatus"
                                                :unit="ownUnit"
                                                @close="openUnitStatus = false"
                                            />
                                        </template>
                                        <button
                                            @click="selectUnitOpen = true"
                                            type="button"
                                            class="text-white bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                        >
                                            <span v-if="!ownUnit" class="flex w-full flex-col items-center">
                                                <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                                <span class="mt-2 truncate">{{ $t('common.no_own_unit') }}</span>
                                            </span>
                                            <span v-else class="truncate">{{ $t('common.leave_unit') }}</span>
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
                                                    <span class="text-base-200 leading-7">
                                                        <div class="text-xs font-semibold leading-6 text-base-200">
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
                                                            <button
                                                                v-for="(item, idx) in actionsUnit"
                                                                :key="item.name"
                                                                type="button"
                                                                class="text-white bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                                                :class="[
                                                                    idx >= actionsUnit.length - 1 ? 'col-span-2' : '',
                                                                    item.class,
                                                                ]"
                                                                @click="
                                                                    unitStatusSelected = item.status;
                                                                    unitStatusOpen = true;
                                                                "
                                                            >
                                                                <component
                                                                    :is="item.icon ?? HoopHouseIcon"
                                                                    class="text-base-200 group-hover:text-white h-5 w-5 shrink-0"
                                                                    aria-hidden="true"
                                                                />
                                                                <span class="mt-1">{{ item.name }}</span>
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
                                        <div class="text-xs font-semibold leading-6 text-base-200">
                                            {{ $t('common.dispatch', 2) }}
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
                                                        item.class,
                                                        dispatchStatusToBGColor(item.status),
                                                    ]"
                                                    @click="
                                                        dispatchStatusSelected = item.status;
                                                        openUnitStatus = true;
                                                    "
                                                >
                                                    <component
                                                        :is="item.icon ?? HoopHouseIcon"
                                                        class="text-base-200 group-hover:text-white h-5 w-5 shrink-0"
                                                        aria-hidden="true"
                                                    />
                                                    <span class="mt-1">{{ item.name }}</span>
                                                </button>
                                            </div>
                                        </li>
                                    </ul>
                                </li>
                                <li>
                                    <div class="text-xs font-semibold leading-6 text-base-200">
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
                                            <DispatchEntry
                                                v-for="dispatch in ownDispatches"
                                                :dispatch="dispatch"
                                                @goto="$emit('goto', $event)"
                                                @details="
                                                    selectedDispatch = $event;
                                                    openDispatchDetails = true;
                                                "
                                            />
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
