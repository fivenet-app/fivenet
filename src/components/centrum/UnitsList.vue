<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { UNIT_STATUS, Unit } from '~~/gen/ts/resources/dispatch/units';

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
                        <li v-for="unit in units" :key="unit.name" class="col-span-1 flex rounded-md shadow-sm">
                            <div
                                class="flex w-12 flex-shrink-0 items-center justify-center rounded-l-md text-sm font-medium text-white"
                                :style="'background-color: #' + unit.color ?? '00000'"
                            >
                                {{ unit.initials }}
                            </div>
                            <div
                                class="flex flex-1 items-center justify-between truncate rounded-r-md border-b border-r border-t border-gray-200 bg-gray"
                            >
                                <div class="flex-1 truncate px-4 py-2 text-sm">
                                    <span class="font-medium text-gray-100">{{ unit.name }}</span>
                                    <p class="text-gray-400">{{ $t('common.members', unit.users.length) }}</p>
                                </div>
                                <div class="flex-shrink-0 pr-5">
                                    <button
                                        type="button"
                                        class="inline-flex items-center justify-center text-white bg-green-700"
                                    >
                                        {{ UNIT_STATUS[unit.status as number] }}
                                    </button>
                                </div>
                            </div>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
