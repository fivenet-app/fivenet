<script lang="ts" setup>
import Feed from './Feed.vue';
import { default as DispatchesList } from '~/components/centrum/dispatches/List.vue';
import Livemap from '~/components/centrum/Livemap.vue';
import { default as UnitsList } from '~/components/centrum/units/List.vue';
import { Dispatch, DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatch';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { Unit, UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import CreateOrUpdateModal from '~/components/centrum/dispatches/CreateOrUpdateModal.vue';
import { LeafletMouseEvent } from 'leaflet';

const { $grpc } = useNuxtApp();

const feed = ref<(DispatchStatus | UnitStatus)[]>([]);

const { data: dispatches } = useLazyAsyncData(`centrum-dispatches`, () => listDispatches());

async function listDispatches(): Promise<Array<Dispatch>> {
    return new Promise(async (res, rej) => {
        try {
            const req = {
                status: [],
            };

            const call = $grpc.getCentrumClient().listDispatches(req);
            const { response } = await call;

            return res(response.dispatches);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { data: units } = useLazyAsyncData(`centrum-units`, () => listUnits());

async function listUnits(): Promise<Array<Unit>> {
    return new Promise(async (res, rej) => {
        try {
            const req = {
                status: [],
            };

            const call = $grpc.getCentrumClient().listUnits(req);
            const { response } = await call;

            return res(response.units);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
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
            if (!dispatches.value) {
                continue;
            }

            console.debug('Centrum: Received change - Kind:', resp.change.oneofKind);

            if (resp.change.oneofKind === 'dispatchUpdate') {
                const id = resp.change.dispatchUpdate.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchUpdate);
                } else {
                    dispatches.value![idx] = resp.change.dispatchUpdate;
                }
            } else if (resp.change.oneofKind === 'dispatchStatus') {
                feed.value.unshift(resp.change.dispatchStatus);
            } else if (resp.change.oneofKind === 'dispatchUnassigned') {
                const id = resp.change.dispatchUnassigned.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchUnassigned);
                } else {
                    dispatches.value![idx].units = resp.change.dispatchUnassigned.units;
                }
            } else if (resp.change.oneofKind === 'dispatchAssigned') {
                const id = resp.change.dispatchAssigned.id;
                const idx = dispatches.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    dispatches.value?.unshift(resp.change.dispatchAssigned);
                } else {
                    dispatches.value![idx].units = resp.change.dispatchAssigned.units;
                }
            } else if (resp.change.oneofKind === 'unitUpdate') {
                const id = resp.change.unitUpdate.id;
                const idx = units.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx === -1) {
                    units.value?.unshift(resp.change.unitUpdate);
                } else {
                    units.value![idx] = resp.change.unitUpdate;
                }
            } else if (resp.change.oneofKind === 'unitStatus') {
                feed.value.unshift(resp.change.unitStatus);
            } else if (resp.change.oneofKind === 'unitAssigned') {
                // TODO
            } else if (resp.change.oneofKind === 'unitDeleted') {
                const id = resp.change.unitDeleted;
                const idx = units.value?.findIndex((d) => d.id === id) ?? -1;
                if (idx > -1) {
                    units.value?.splice(idx, 1);
                }
            } else {
                console.log('Centrum: Unknown change received - Kind: ', resp.change.oneofKind, resp.change);
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
</script>

<template>
    <div class="flex-col h-full">
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
                <div class="basis-3/5">
                    <DispatchesList :dispatches="dispatches" :units="units" @goto="goto($event)" />
                </div>
                <div class="basis-1/5">
                    <Feed :items="feed" />
                </div>
                <div class="basis-1/5">
                    <UnitsList :units="units" @goto="goto($event)" />
                </div>
            </div>
        </div>
    </div>
</template>
