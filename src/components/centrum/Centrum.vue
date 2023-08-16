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
import { Dispatch, DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatches';
import { Settings } from '~~/gen/ts/resources/dispatch/settings';
import { Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { UserShort } from '~~/gen/ts/resources/users/users';
import Feed from './Feed.vue';

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const settings = ref<Settings>();
const isDisponent = ref(false);
const ownUnit = ref<Unit>();
const feed = ref<(DispatchStatus | UnitStatus)[]>([]);
const disponents = ref<UserShort[]>([]);
const units = ref<Unit[]>([]);
const dispatches = ref<Dispatch[]>([]);

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
                isDisponent.value = resp.change.latestState.isDisponent;
            } else if (resp.change.oneofKind === 'settings') {
                settings.value = resp.change.settings;
            } else if (resp.change.oneofKind === 'disponents') {
                disponents.value = resp.change.disponents.disponents;
                // If user is not part of disponents list anymore, we need to restart the stream
                const idx = disponents.value.findIndex((d) => d.userId === activeChar.value?.userId);
                if (idx === -1) {
                    stopStream();
                    setTimeout(() => {
                        startStream();
                    }, 250);
                }
            } else if (resp.change.oneofKind === 'unitAssigned') {
                // Ignore, doesn't matter for controllers
            } else if (resp.change.oneofKind === 'unitCreated') {
                const id = resp.change.unitCreated.id;
                const idx = units.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    units.value?.unshift(resp.change.unitCreated);
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
                const idx = units.value.findIndex((u) => u.id === unitId);
                if (idx > -1) {
                    units.value[idx].status = resp.change.unitStatus;
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
                const dispatchId = resp.change.dispatchStatus.dispatchId;
                const idx = dispatches.value.findIndex((d) => d.id === dispatchId);
                if (idx > -1) {
                    dispatches.value[idx].status = resp.change.dispatchStatus;
                }
            } else {
                console.warn('Centrum: Unknown change received - Kind: ', resp.change.oneofKind, resp.change);
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

function livemapContextmenu($event: LeafletMouseEvent) {
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
                    <Livemap ref="livemapComponent" @contextmenu="livemapContextmenu($event)" />
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
