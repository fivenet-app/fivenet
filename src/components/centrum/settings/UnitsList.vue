<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { Unit } from '~~/gen/ts/resources/centrum/units';
import UnitCreateOrUpdateModal from '~/components/centrum/settings/UnitCreateOrUpdateModal.vue';
import UnitsListEntry from '~/components/centrum/settings/UnitsListEntry.vue';
import CentrumSettingsModal from '~/components/centrum/settings/CentrumSettingsModal.vue';
import GenericTable from '~/components/partials/elements/GenericTable.vue';

const { $grpc } = useNuxtApp();

const { data: units, pending, refresh, error } = useLazyAsyncData('centrum-units', () => getUnits());

async function getUnits(): Promise<Unit[]> {
    try {
        const call = $grpc.getCentrumClient().listUnits({
            status: [],
        });
        const { response } = await call;

        return response.units;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const open = ref(false);
const openSettings = ref(false);
</script>

<template>
    <UnitCreateOrUpdateModal v-if="open" :open="open" @close="open = false" @created="refresh()" />
    <CentrumSettingsModal :open="openSettings" @close="openSettings = false" />

    <div class="py-2">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="mt-2 flow-root">
                <div v-if="can('CentrumService.CreateOrUpdateUnit')" class="sm:flex sm:items-center">
                    <div class="sm:flex-auto">
                        <div class="grid flex-initial grid-cols-2 gap-4">
                            <UButton
                                v-if="can('CentrumService.CreateOrUpdateUnit')"
                                class="inline-flex rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                @click="open = true"
                            >
                                {{ $t('components.centrum.units.create_unit') }}
                            </UButton>
                            <UButton
                                v-if="can('CentrumService.Stream')"
                                class="inline-flex rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                @click="openSettings = true"
                            >
                                {{ $t('common.setting', 2) }}
                            </UButton>
                        </div>
                    </div>
                </div>
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.units')])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.units')])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock v-else-if="units && units.length === 0" :type="$t('common.units')" />
                        <template v-else>
                            <GenericTable>
                                <template #thead>
                                    <tr>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.name') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.initials') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.description') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.attributes', 2) }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.color') }}
                                        </th>
                                        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                                            {{ $t('common.department_postal') }}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold text-neutral sm:pr-0"
                                        >
                                            {{ $t('common.action', 2) }}
                                        </th>
                                    </tr>
                                </template>
                                <template #tbody>
                                    <UnitsListEntry v-for="unit in units" :key="unit.id" :unit="unit" @updated="refresh()" />
                                </template>
                            </GenericTable>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
