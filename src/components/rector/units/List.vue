<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import CreateOrUpdateUnitModal from './CreateOrUpdateUnitModal.vue';
import ListEntry from './ListEntry.vue';
import SettingsModal from './SettingsModal.vue';

const { $grpc } = useNuxtApp();

const { data: units, pending, refresh, error } = useLazyAsyncData('rector-units', () => getUnits());

async function getUnits(): Promise<Unit[]> {
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
const openSettings = ref(false);
</script>

<template>
    <CreateOrUpdateUnitModal :open="open" @close="open = false" @refresh="refresh" />
    <SettingsModal :open="openSettings" @close="openSettings = false" />

    <div class="py-2">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="flow-root mt-2">
                <div v-if="can('CentrumService.CreateOrUpdateUnit')" class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <div class="flex-initial form-control grid grid-cols-2 gap-4">
                            <button
                                v-if="can('CentrumService.CreateOrUpdateUnit')"
                                @click="open = true"
                                class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                            >
                                {{ $t('pages.rector.units.create_unit') }}
                            </button>
                            <button
                                v-if="can('CentrumService.GetSettings')"
                                @click="openSettings = true"
                                class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                            >
                                {{ $t('common.setting', 2) }}
                            </button>
                        </div>
                    </div>
                </div>
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full py-2 align-middle px-1">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.units')])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.units')])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock v-else-if="units && units.length === 0" :type="$t('common.units')" />
                        <div v-else>
                            <table class="min-w-full divide-y divide-base-600">
                                <thead>
                                    <tr>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.initials') }}
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
                                            {{ $t('common.initials') }}
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
