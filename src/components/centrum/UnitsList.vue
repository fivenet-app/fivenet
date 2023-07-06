<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import UnitsListEntry from './UnitsListEntry.vue';

const { $grpc } = useNuxtApp();

const { data: units, pending, refresh, error } = useLazyAsyncData(`centrum-units`, () => listUnits());

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
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-scroll">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h1 class="text-base font-semibold leading-6 text-gray-100">Active Units</h1>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <ul role="list" class="mt-3 grid grid-cols-1 gap-5 sm:grid-cols-2 sm:gap-2 lg:grid-cols-3">
                        <UnitsListEntry
                            v-for="unit in units?.sort((a, b) => {
                                if (a.status < b.status) return -1;
                                if (a.status > b.status) return 1;
                                return 0;
                            })"
                            :unit="unit"
                            :key="unit.id.toString()"
                        />
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
