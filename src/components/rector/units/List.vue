<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import CreateOrUpdateUnitModal from './CreateOrUpdateUnitModal.vue';
import ListEntry from './ListEntry.vue';

const { $grpc } = useNuxtApp();

const { data: units, pending, refresh, error } = useLazyAsyncData('rector-units', () => getUnits());

async function getUnits(): Promise<Array<Unit>> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().listUnits({
                status: [],
            });
            const { response } = await call;

            return res(response.units);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const open = ref(false);
</script>

<template>
    <CreateOrUpdateUnitModal :open="open" @close="open = false" @refresh="refresh" />
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="flow-root mt-2">
                <div v-can="'RectorService.CreateRole'" class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <div class="flex-initial form-control" v-can="'DocStoreService.CreateTemplate'">
                            <button
                                @click="open = true"
                                class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                            >
                                {{ $t('pages.rector.units.create_unit') }}
                            </button>
                        </div>
                    </div>
                </div>
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.unit', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.unit', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock v-else-if="units && units.length === 0" :type="$t('common.unit', 2)" />
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.initials', 2) }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.description') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.color') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody class="divide-y divide-base-800">
                                    <ListEntry v-for="unit in units" :unit="unit" />
                                </tbody>
                                <thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.initials', 2) }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.description') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.color') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </thead>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
