<script lang="ts" setup>
import DispatchFeed from '~/components/centrum/DispatchFeed.vue';
import DispatchesList from '~/components/centrum/DispatchesList.vue';
import Livemap from '~/components/centrum/Livemap.vue';
import UnitsList from '~/components/centrum/UnitsList.vue';
import { Dispatch } from '~~/gen/ts/resources/dispatch/dispatch';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import DispatchCreateModal from './DispatchCreateModal.vue';

const { $grpc } = useNuxtApp();

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

// TODO Handle the data storage

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

            console.log('CHANGE', resp.change);
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

const open = ref(false);
</script>

<template>
    <div class="flex-col h-full">
        <DispatchCreateModal :open="open" @close="open = false" />
        <div class="relative w-full h-full z-0 flex">
            <!-- Left column -->
            <div class="flex flex-col basis-1/3 divide-x">
                <div class="h-full">
                    <Livemap @contextmenu="open = true" />
                </div>
            </div>

            <!-- Right column -->
            <div class="flex flex-col basis-2/3 divide-y">
                <div class="basis-3/5">
                    <DispatchesList :dispatches="dispatches" :units="units" />
                </div>
                <div class="basis-1/5">
                    <DispatchFeed :units="units" />
                </div>
                <div class="basis-1/5">
                    <UnitsList :units="units" />
                </div>
            </div>
        </div>
    </div>
</template>
