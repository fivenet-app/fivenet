<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import UnitCreateOrUpdateModal from '~/components/centrum/settings/UnitCreateOrUpdateModal.vue';
import CentrumSettingsModal from '~/components/centrum/settings/CentrumSettingsModal.vue';
import type { ListUnitsResponse } from '~~/gen/ts/services/centrum/centrum';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const modal = useModal();

const { data: units, pending: loading, refresh, error } = useLazyAsyncData('centrum-units', () => listUnits());

async function listUnits(): Promise<ListUnitsResponse> {
    try {
        const call = $grpc.getCentrumClient().listUnits({
            status: [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteUnit(id: string): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().deleteUnit({
            unitId: id,
        });
        await call;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const columns = [
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'initials',
        label: t('common.initials'),
    },
    {
        key: 'description',
        label: t('common.description'),
    },
    {
        key: 'attributes',
        label: t('common.attributes'),
    },
    {
        key: 'color',
        label: t('common.color'),
    },
    {
        key: 'homePostal',
        label: t('common.department_postal'),
    },
    {
        key: 'actions',
        label: t('common.action', 2),
    },
];
</script>

<template>
    <div>
        <UDashboardNavbar :title="$t('common.units')">
            <template #right>
                <UButton
                    v-if="can('CentrumService.Stream')"
                    icon="i-mdi-settings"
                    @click="modal.open(CentrumSettingsModal, {})"
                >
                    {{ $t('common.setting', 2) }}
                </UButton>

                <UButton
                    v-if="can('CentrumService.CreateOrUpdateUnit')"
                    trailing-icon="i-mdi-plus"
                    color="gray"
                    @click="
                        modal.open(UnitCreateOrUpdateModal, {
                            onCreated: async () => refresh(),
                            onUpdated: async () => refresh(),
                        })
                    "
                >
                    {{ $t('components.centrum.units.create_unit') }}
                </UButton>
            </template>
        </UDashboardNavbar>

        <div class="px-1 sm:px-2 lg:px-4">
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
                        <DataErrorBlock
                            v-if="error"
                            :title="$t('common.unable_to_load', [$t('common.units')])"
                            :retry="refresh"
                        />
                        <UTable
                            v-else
                            :loading="loading"
                            :columns="columns"
                            :rows="units?.units"
                            :empty-state="{ icon: 'i-mdi-car', label: $t('common.not_found', [$t('common.entry', 2)]) }"
                        >
                            <template #attributes-data="{ row: unit }">
                                <template v-if="unit.attributes !== undefined && unit.attributes?.list.length > 0">
                                    <span
                                        v-for="attribute in unit.attributes?.list"
                                        :key="attribute"
                                        class="inline-flex items-center rounded-md bg-warn-400/10 px-2 py-1 text-xs font-medium text-warn-400 ring-1 ring-inset ring-warn-400/20"
                                    >
                                        {{ $t(`components.centrum.units.attributes.${attribute}`) }}
                                    </span>
                                </template>
                                <template v-else>
                                    {{ $t('common.none', [$t('common.attributes', 2)]) }}
                                </template>
                            </template>
                            <template #color-data="{ row: unit }">
                                <ColorInput v-model="unit.color" disabled format="hex" class="h-6" />
                            </template>
                            <template #homePostal-data="{ row: unit }">
                                {{ unit.homePostal ?? $t('common.na') }}
                            </template>
                            <template #actions-data="{ row: unit }">
                                <div class="flex items-center">
                                    <UButton
                                        v-if="can('CentrumService.CreateOrUpdateUnit')"
                                        variant="link"
                                        icon="i-mdi-pencil"
                                        @click="
                                            modal.open(UnitCreateOrUpdateModal, {
                                                unit: unit,
                                                onUpdated: async () => refresh(),
                                            })
                                        "
                                    />
                                    <UButton
                                        v-if="can('CentrumService.DeleteUnit')"
                                        variant="link"
                                        icon="i-mdi-trash-can"
                                        @click="
                                            modal.open(ConfirmModal, {
                                                confirm: async () => deleteUnit(unit.id),
                                            })
                                        "
                                    />
                                </div>
                            </template>
                        </UTable>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
