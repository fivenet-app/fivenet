<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { LeafletMouseEvent } from 'leaflet';
import { HelpCircleIcon } from 'mdi-vue3';
import Livemap from '~/components/centrum/Livemap.vue';
import CreateOrUpdateModal from '~/components/centrum/dispatches/CreateOrUpdateModal.vue';
import { default as DispatchesList } from '~/components/centrum/dispatches/List.vue';
import { default as UnitsList } from '~/components/centrum/units/List.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useAuthStore } from '~/store/auth';
import { DISPATCH_STATUS, Dispatch, DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatches';
import { Settings } from '~~/gen/ts/resources/dispatch/settings';
import { Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';
import Feed from './Feed.vue';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const settings = ref<Settings>();
const isDisponent = ref(false);
const feed = ref<(DispatchStatus | UnitStatus)[]>([]);
const disponents = ref<UserShort[]>([]);
const units = ref<Unit[]>([]);
const dispatches = ref<Dispatch[]>([]);

function addOrUpdateUnit(unit: Unit): void {
    const idx = units.value?.findIndex((d) => d.id === unit.id) ?? -1;
    if (idx === -1) {
        units.value?.unshift(unit);
    } else {
        units.value[idx].job = unit.job;
        units.value[idx].createdAt = unit.createdAt;
        units.value[idx].updatedAt = unit.updatedAt;
        units.value[idx].name = unit.name;
        units.value[idx].initials = unit.initials;
        units.value[idx].color = unit.color;
        units.value[idx].description = unit.description;
        units.value[idx].status = unit.status;
        units.value[idx].users = unit.users;
    }
}

function addOrUpdateDispatch(dispatch: Dispatch): void {
    const idx = dispatches.value?.findIndex((d) => d.id === dispatch.id) ?? -1;
    if (idx === -1) {
        dispatches.value?.unshift(dispatch);
    } else {
        dispatches.value[idx].createdAt = dispatch.createdAt;
        dispatches.value[idx].updatedAt = dispatch.updatedAt;
        dispatches.value[idx].job = dispatch.job;
        dispatches.value[idx].status = dispatch.status;
        dispatches.value[idx].message = dispatch.message;
        dispatches.value[idx].description = dispatch.description;
        dispatches.value[idx].attributes = dispatch.attributes;
        dispatches.value[idx].x = dispatch.x;
        dispatches.value[idx].y = dispatch.y;
        dispatches.value[idx].anon = dispatch.anon;
        dispatches.value[idx].userId = dispatch.userId;
        dispatches.value[idx].user = dispatch.user;
        if (dispatch.units.length == 0) {
            dispatches.value[idx].units.length = 0;
        } else {
            dispatches.value[idx].units = dispatch.units;
        }
    }
}

function removeUnit(unit: Unit): void {
    const idx = units.value?.findIndex((d) => d.id === unit.id) ?? -1;
    if (idx > -1) {
        units.value?.splice(idx, 1);
    }
}

function removeDispatchFromList(id: bigint): void {
    const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
    if (idx > -1) {
        dispatches.value?.splice(idx, 1);
    }
}

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
                units.value = resp.change.latestState.units;
                dispatches.value = resp.change.latestState.dispatches;
                isDisponent.value = resp.change.latestState.isDisponent;
            } else if (resp.change.oneofKind === 'settings') {
                settings.value = resp.change.settings;
            } else if (resp.change.oneofKind === 'disponents') {
                disponents.value = resp.change.disponents.disponents;
                // If user is not part of disponents list anymore, we need to reset the lists and restart the stream
                const idx = disponents.value.findIndex((d) => d.userId === activeChar.value?.userId);
                if (idx === -1) {
                    stopStream();
                    setTimeout(() => {
                        startStream();
                    }, 250);
                }
            } else if (resp.change.oneofKind === 'unitAssigned') {
                // Ignore, doesn't matter for controllers
            } else if (resp.change.oneofKind === 'unitDeleted') {
                removeUnit(resp.change.unitDeleted);
            } else if (resp.change.oneofKind === 'unitUpdated') {
                addOrUpdateUnit(resp.change.unitUpdated);
            } else if (resp.change.oneofKind === 'unitStatus') {
                const id = resp.change.unitStatus.id;
                let idx = dispatches.value.findIndex((d) => d.id === id);
                if (idx === -1) {
                    units.value?.unshift(resp.change.unitStatus);
                } else {
                    units.value[idx] = resp.change.unitStatus;
                }

                if (resp.change.unitStatus.status) {
                    feed.value.unshift(resp.change.unitStatus.status);
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
                removeDispatchFromList(resp.change.dispatchDeleted.id);
            } else if (resp.change.oneofKind === 'dispatchUpdated') {
                addOrUpdateDispatch(resp.change.dispatchUpdated);
            } else if (resp.change.oneofKind === 'dispatchStatus') {
                const id = resp.change.dispatchStatus.id;
                let idx = dispatches.value.findIndex((d) => d.id === id);
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchStatus);
                } else {
                    dispatches.value[idx] = resp.change.dispatchStatus;
                }

                if (resp.change.dispatchStatus.status) {
                    feed.value.unshift(resp.change.dispatchStatus.status);
                }

                if (
                    resp.change.dispatchStatus.status?.status === DISPATCH_STATUS.UNIT_UNASSIGNED ||
                    resp.change.dispatchStatus.status?.status === DISPATCH_STATUS.ARCHIVED
                ) {
                    removeDispatchFromList(id);
                }
            } else {
                console.warn('Centrum: Unknown change received - Kind: ', resp.change.oneofKind, resp.change);
            }

            if (resp.restart !== undefined && resp.restart) {
                stopStream();
                setTimeout(() => {
                    startStream();
                }, 250);
            }
        }
    } catch (e) {
        const err = e as RpcError;
        error.value = err.message;
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

const createOrUpdateModal = ref<InstanceType<typeof CreateOrUpdateModal>>();
const livemapComponent = ref<InstanceType<typeof Livemap>>();

const open = ref(false);
const location = ref<{ x: number; y: number }>({ x: 0, y: 0 });

function livemapCreateDispatch($event: LeafletMouseEvent) {
    goto({ x: $event.latlng.lat, y: $event.latlng.lng });

    open.value = true;
}

function goto(e: { x: number; y: number }) {
    if (createOrUpdateModal.value) {
        createOrUpdateModal.value.location = { x: e.x, y: e.y };
    }
    if (livemapComponent.value) {
        livemapComponent.value.location = { x: e.x, y: e.y };
    }
}

async function takeControl(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().takeControl({
                signon: true,
            });
            await call;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <div class="flex-col h-full relative">
        <div
            v-if="error || abort === undefined"
            class="absolute inset-0 flex justify-center items-center z-20"
            style="background-color: rgba(62, 60, 62, 0.5)"
        >
            <DataPendingBlock v-if="!error" :message="$t('components.livemap.starting_datastream')" />
            <DataErrorBlock v-else="error" :title="$t('components.livemap.failed_datastream')" :retry="startStream" />
        </div>
        <div v-else-if="!isDisponent">
            <div class="absolute inset-0 flex justify-center items-center z-20" style="background-color: rgba(62, 60, 62, 0.5)">
                <button
                    @click="takeControl()"
                    type="button"
                    class="relative block w-full p-12 text-center border-2 border-dotted rounded-lg border-base-300 hover:border-base-400 focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2"
                >
                    <HelpCircleIcon class="w-12 h-12 mx-auto text-neutral" />
                    <span class="block mt-2 text-sm font-semibold text-gray-300">
                        No one is in the dispatch center. Want to take control? Click here.
                    </span>
                </button>
            </div>
        </div>

        <CreateOrUpdateModal ref="createOrUpdateModal" :open="open" @close="open = false" :location="location" />

        <div class="relative w-full h-full z-0 flex">
            <!-- Left column -->
            <div class="flex flex-col basis-1/3 divide-x">
                <div class="h-full">
                    <Livemap ref="livemapComponent" @create-dispatch="livemapCreateDispatch($event)" />
                </div>
            </div>

            <!-- Right column -->
            <div class="flex flex-col basis-2/3 divide-y">
                <div class="basis-3/5 max-h-[60%]">
                    <DispatchesList :dispatches="dispatches" :units="units" @goto="goto($event)" />
                </div>
                <div class="basis-1/5 max-h-[20%]">
                    <Feed :items="feed" />
                </div>
                <div class="basis-1/5 max-h-[20%]">
                    <UnitsList :units="units" @goto="goto($event)" />
                </div>
            </div>
        </div>
    </div>
</template>
