<script lang="ts" setup>
import UnitAttributes from '~/components/centrum/partials/UnitAttributes.vue';
import UnitCreateOrUpdateModal from '~/components/centrum/settings/UnitCreateOrUpdateModal.vue';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import type { ListUnitsResponse } from '~~/gen/ts/services/centrum/centrum';

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const { data: units, pending: loading, refresh, error } = useLazyAsyncData('centrum-units', () => listUnits());

async function listUnits(): Promise<ListUnitsResponse> {
    try {
        const call = getGRPCCentrumClient().listUnits({
            status: [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteUnit(id: string): Promise<void> {
    try {
        const call = getGRPCCentrumClient().deleteUnit({
            unitId: id,
        });
        await call;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const columns = [
    {
        key: 'name',
        label: t('common.name'),
        sortable: true,
    },
    {
        key: 'initials',
        label: t('common.initials'),
        sortable: true,
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
        sortable: false,
    },
];
</script>

<template>
    <UDashboardNavbar :title="$t('common.unit')">
        <template #right>
            <UButton color="black" icon="i-mdi-arrow-back" to="/centrum">
                {{ $t('common.back') }}
            </UButton>

            <UButton v-if="can('CentrumService.Stream').value" icon="i-mdi-settings" to="/centrum/settings">
                {{ $t('common.setting', 2) }}
            </UButton>

            <UButton
                v-if="can('CentrumService.CreateOrUpdateUnit').value"
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

    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.unit', 2)])" :retry="refresh" />
    <UTable
        v-else
        :loading="loading"
        :columns="columns"
        :rows="units?.units"
        :empty-state="{ icon: 'i-mdi-car', label: $t('common.not_found', [$t('common.unit', 2)]) }"
    >
        <template #name-data="{ row: unit }">
            <div class="text-gray-900 dark:text-white">
                {{ unit.name }}
            </div>
        </template>
        <template #attributes-data="{ row: unit }">
            <UnitAttributes :attributes="unit.attributes" />
        </template>
        <template #color-data="{ row: unit }">
            <ColorPickerClient v-model="unit.color" disabled hide-icon />
        </template>
        <template #homePostal-data="{ row: unit }">
            {{ unit.homePostal ?? $t('common.na') }}
        </template>
        <template #actions-data="{ row: unit }">
            <div .key="unit.id" class="flex items-center">
                <UButton
                    v-if="can('CentrumService.CreateOrUpdateUnit').value"
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
                    v-if="can('CentrumService.DeleteUnit').value"
                    variant="link"
                    icon="i-mdi-trash-can"
                    color="red"
                    @click="
                        modal.open(ConfirmModal, {
                            confirm: async () => deleteUnit(unit.id),
                        })
                    "
                />
            </div>
        </template>
    </UTable>
</template>
