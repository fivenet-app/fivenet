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
import { default as UnitDetails } from '~/components/centrum/units/Details.vue';
import DispatchEntry from '~/components/livemap/centrum/DispatchEntry.vue';
import { useAuthStore } from '~/store/auth';
import { useNotificationsStore } from '~/store/notifications';
import { DISPATCH_STATUS, Dispatch, DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatches';
import { Settings } from '~~/gen/ts/resources/dispatch/settings';
import { UNIT_STATUS, Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';
import JoinUnit from './centrum/JoinUnit.vue';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const settings = ref<Settings>();
const disponents = ref<UserShort[]>([]);
const ownUnit = ref<Unit>();
const units = ref<Unit[]>([]);
const dispatches = ref<Dispatch[]>([]);
const feed = ref<(DispatchStatus | UnitStatus)[]>([]);

type Action = {
    icon: DefineComponent;
    name: string;
    action?: Function;
    class?: string;
    status?: DISPATCH_STATUS | UNIT_STATUS;
};

const actionsDispatch: Action[] = [
    { icon: markRaw(CarBackIcon), name: 'En Route', class: 'bg-info-600', status: DISPATCH_STATUS.EN_ROUTE },
    { icon: markRaw(MarkerCheckIcon), name: 'On Scene', class: 'bg-primary-600', status: DISPATCH_STATUS.ON_SCENE },
    { icon: markRaw(HelpCircleIcon), name: 'Need Assistance', class: 'bg-warn-600', status: DISPATCH_STATUS.NEED_ASSISTANCE },
    { icon: markRaw(CheckBoldIcon), name: 'Completed', class: 'bg-success-600', status: DISPATCH_STATUS.COMPLETED },
    { icon: markRaw(ListStatusIcon), name: 'Update Status', class: 'bg-base-800' },
];

const actionsUnit: Action[] = [
    { icon: markRaw(CarBackIcon), name: 'Unavailable', class: 'bg-error-600', status: UNIT_STATUS.UNAVAILABLE },
    { icon: markRaw(CalendarCheckIcon), name: 'Available', class: 'bg-success-600', status: UNIT_STATUS.AVAILABLE },
    { icon: markRaw(CoffeeIcon), name: 'On Break', class: 'bg-warn-600', status: UNIT_STATUS.ON_BREAK },
    { icon: markRaw(CalendarRemoveIcon), name: 'Busy', class: 'bg-info-600', status: UNIT_STATUS.BUSY },
    { icon: markRaw(ListStatusIcon), name: 'Update Status', class: 'bg-base-800' },
];

const abort = ref<AbortController | undefined>();
const error = ref<string | null>(null);
async function startStream(): Promise<void> {
    if (abort.value !== undefined) return;

    console.debug('Centrum: Starting Data Stream');
    try {
        abort.value = new AbortController();

        const call = $grpc.getCentrumClient().stream(
            {},
            {
                abort: abort.value.signal,
            },
        );

        for await (let resp of call.responses) {
            error.value = null;

            if (resp === undefined || !resp.change) {
                continue;
            }

            console.debug('Centrum: Received change - Kind:', resp.change.oneofKind, resp.change);

            if (resp.change.oneofKind === 'latestState') {
                settings.value = resp.change.latestState.settings;
                ownUnit.value = resp.change.latestState.unit;
                units.value = resp.change.latestState.units;
                dispatches.value = resp.change.latestState.dispatches;
            } else if (resp.change.oneofKind === 'settings') {
                settings.value = resp.change.settings;
            } else if (resp.change.oneofKind === 'disponents') {
                disponents.value = resp.change.disponents.disponents;
            } else if (resp.change.oneofKind === 'unitAssigned') {
                const idx = resp.change.unitAssigned.users.findIndex((u) => u.userId !== activeChar.value?.userId);
                if (idx === -1) {
                    // User has been removed from the unit
                    ownUnit.value = undefined;

                    notifications.dispatchNotification({
                        title: { key: 'notifications.centrum.unitAssigned.removed.title', parameters: [] },
                        content: { key: 'notifications.centrum.unitAssigned.removed.content', parameters: [] },
                        type: 'success',
                    });
                } else {
                    // User has been added to unit
                    ownUnit.value = resp.change.unitAssigned;

                    notifications.dispatchNotification({
                        title: { key: 'notifications.centrum.unitAssigned.joined.title', parameters: [] },
                        content: { key: 'notifications.centrum.unitAssigned.joined.content', parameters: [] },
                        type: 'success',
                    });
                }
            } else if (resp.change.oneofKind === 'unitDeleted') {
                const id = resp.change.unitDeleted.id;
                const idx = units.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx > -1) {
                    units.value?.splice(idx, 1);
                }
            } else if (resp.change.oneofKind === 'unitUpdated') {
                const id = resp.change.unitUpdated.id;
                const idx = units.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    units.value?.unshift(resp.change.unitUpdated);
                } else {
                    units.value[idx].job = resp.change.unitUpdated.job;
                    units.value[idx].createdAt = resp.change.unitUpdated.createdAt;
                    units.value[idx].updatedAt = resp.change.unitUpdated.updatedAt;
                    units.value[idx].name = resp.change.unitUpdated.name;
                    units.value[idx].initials = resp.change.unitUpdated.initials;
                    units.value[idx].color = resp.change.unitUpdated.color;
                    units.value[idx].description = resp.change.unitUpdated.description;
                    units.value[idx].status = resp.change.unitUpdated.status;
                    units.value[idx].users = resp.change.unitUpdated.users;
                }
            } else if (resp.change.oneofKind === 'unitStatus') {
                feed.value.unshift(resp.change.unitStatus);
                const unitId = resp.change.unitStatus.unitId;
                const unit = units.value.find((u) => u.id === unitId);
                if (unit) {
                    unit.status = resp.change.unitStatus;
                }
            } else if (resp.change.oneofKind === 'dispatchCreated') {
                const id = resp.change.dispatchCreated.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchCreated);
                } else {
                    dispatches.value[idx].units = resp.change.dispatchCreated.units;
                }
            } else if (resp.change.oneofKind === 'dispatchDeleted') {
                const id = resp.change.dispatchDeleted.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx > -1) {
                    dispatches.value?.splice(idx, 1);
                }
            } else if (resp.change.oneofKind === 'dispatchUpdated') {
                const id = resp.change.dispatchUpdated.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchUpdated);
                } else {
                    dispatches.value[idx].createdAt = resp.change.dispatchUpdated.createdAt;
                    dispatches.value[idx].updatedAt = resp.change.dispatchUpdated.updatedAt;
                    dispatches.value[idx].job = resp.change.dispatchUpdated.job;
                    dispatches.value[idx].status = resp.change.dispatchUpdated.status;
                    dispatches.value[idx].message = resp.change.dispatchUpdated.message;
                    dispatches.value[idx].description = resp.change.dispatchUpdated.description;
                    dispatches.value[idx].attributes = resp.change.dispatchUpdated.attributes;
                    dispatches.value[idx].x = resp.change.dispatchUpdated.x;
                    dispatches.value[idx].y = resp.change.dispatchUpdated.y;
                    dispatches.value[idx].anon = resp.change.dispatchUpdated.anon;
                    dispatches.value[idx].userId = resp.change.dispatchUpdated.userId;
                    dispatches.value[idx].user = resp.change.dispatchUpdated.user;
                    dispatches.value[idx].units = resp.change.dispatchUpdated.units;
                }
            } else if (resp.change.oneofKind === 'dispatchStatus') {
                feed.value.unshift(resp.change.dispatchStatus);
            } else {
                console.warn('Centrum: Unknown change received - Kind: ', resp.change.oneofKind, resp.change);
            }
        }
    } catch (e) {
        const err = e as RpcError;
        error.value = err.message;
        notifications.dispatchNotification({
            content: { key: err.message, parameters: [] },
            title: { key: '', parameters: [] },
        });
        stopStream();
    }

    console.debug('Centrum: Data Stream Ended');
}

async function stopStream(): Promise<void> {
    console.debug('Centrum: Stopping Data Stream');
    abort.value?.abort();
    abort.value = undefined;
}

onMounted(() => {
    startStream();
});

onBeforeUnmount(() => {
    stopStream();
});

const unitOpen = ref(false);
const selectUnitOpen = ref(false);
</script>

<template>
    <!-- Sidebar component, swap this element with another sidebar if you like -->
    <div class="h-full flex grow gap-y-5 overflow-y-auto bg-base-600 px-4 py-0.5">
        <nav class="flex flex-1 flex-col">
            <ul role="list" class="flex flex-1 flex-col gap-y-2 divide-y divide-base-400">
                <li>
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li>
                            <div v-if="ownUnit">
                                <button
                                    v-if="ownUnit"
                                    @click="unitOpen = true"
                                    type="button"
                                    class="text-white bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                >
                                    <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                    <span class="mt-2 truncate">{{ ownUnit.initials }}: {{ ownUnit.name }}</span>
                                </button>
                                <UnitDetails :unit="ownUnit" :ownUnit="ownUnit" :open="unitOpen" @close="unitOpen = false" />
                            </div>
                            <button
                                @click="selectUnitOpen = true"
                                type="button"
                                class="text-white bg-info-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                            >
                                <span v-if="!ownUnit" class="flex w-full flex-col items-center">
                                    <InformationOutlineIcon class="h-5 w-5" aria-hidden="true" />
                                    <span class="mt-2 truncate">Not in any Unit.</span>
                                </span>
                                <span v-else class="truncate">Leave Unit</span>
                            </button>

                            <JoinUnit
                                :open="selectUnitOpen"
                                @close="selectUnitOpen = false"
                                @joined="ownUnit = $event"
                                :own-unit="ownUnit"
                                :units="units"
                            />
                        </li>
                    </ul>
                </li>
                <li v-if="ownUnit">
                    <ul role="list" class="-mx-2 space-y-1">
                        <li>
                            <Disclosure as="div" v-slot="{ open }">
                                <DisclosureButton class="flex w-full items-start justify-between text-left text-white">
                                    <span class="text-base-200 leading-7">
                                        <div class="text-xs font-semibold leading-6 text-base-200">Unit</div>
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
                                                :class="[idx >= actionsUnit.length - 1 ? 'col-span-2' : '', item.class]"
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
                <li v-if="ownUnit">
                    <ul role="list" class="-mx-2 space-y-1">
                        <div class="text-xs font-semibold leading-6 text-base-200">Dispatches</div>
                        <li>
                            <div class="grid grid-cols-2 gap-0.5">
                                <button
                                    v-for="(item, idx) in actionsDispatch"
                                    :key="item.name"
                                    type="button"
                                    class="text-white bg-primary hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                                    :class="[idx >= actionsDispatch.length - 1 ? 'col-span-2' : '', item.class]"
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
                <li v-if="ownUnit">
                    <div class="text-xs font-semibold leading-6 text-base-200">Your Dispatches</div>
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                        <li v-if="!dispatches || dispatches.length === 0">
                            <button
                                type="button"
                                class="text-white bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
                            >
                                <CarEmergencyIcon class="h-5 w-5" aria-hidden="true" />
                                <span class="mt-2 truncate">No assigned Dispatches.</span>
                            </button>
                        </li>
                        <DispatchEntry v-else v-for="dispatch in dispatches" :dispatch="dispatch" />
                    </ul>
                </li>
            </ul>
        </nav>
    </div>
</template>
