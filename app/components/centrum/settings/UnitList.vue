<script lang="ts" setup>
import UnitAttributes from '~/components/centrum/partials/UnitAttributes.vue';
import UnitCreateOrUpdateModal from '~/components/centrum/settings/UnitCreateOrUpdateModal.vue';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { availableIcons, fallbackIcon } from '~/components/partials/icons';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ListUnitsResponse } from '~~/gen/ts/services/centrum/centrum';

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const notifications = useNotificationsStore();

const { data: units, pending: loading, refresh, error } = useLazyAsyncData('centrum-units', () => listUnits());

async function listUnits(): Promise<ListUnitsResponse> {
    try {
        const call = $grpc.centrum.centrum.listUnits({
            status: [],
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteUnit(id: number): Promise<void> {
    try {
        const call = $grpc.centrum.centrum.deleteUnit({
            unitId: id,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const columns = [
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
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
        key: 'color',
        label: t('common.color'),
    },
    {
        key: 'icon',
        label: t('common.icon'),
    },
    {
        key: 'attributes',
        label: t('common.attributes'),
    },
    {
        key: 'homePostal',
        label: t('common.department_postal'),
    },
];
</script>

<template>
    <UDashboardNavbar :title="$t('common.unit', 2)">
        <template #right>
            <PartialsBackButton fallback-to="/centrum" />

            <UTooltip v-if="can('centrum.CentrumService/Stream').value" :text="$t('common.setting', 2)">
                <UButton icon="i-mdi-settings" to="/centrum/settings">
                    <span class="hidden truncate sm:block">
                        {{ $t('common.setting', 2) }}
                    </span>
                </UButton>
            </UTooltip>

            <UButton
                v-if="can('centrum.CentrumService/CreateOrUpdateUnit').value"
                trailing-icon="i-mdi-plus"
                color="gray"
                @click="
                    modal.open(UnitCreateOrUpdateModal, {
                        onCreated: async () => refresh(),
                        onUpdated: async () => refresh(),
                    })
                "
            >
                <span class="hidden truncate sm:block">
                    {{ $t('components.centrum.units.create_unit') }}
                </span>
            </UButton>
        </template>
    </UDashboardNavbar>

    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.unit', 2)])" :error="error" :retry="refresh" />
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

        <template #icon-data="{ row: unit }">
            <component
                :is="availableIcons.find((item) => item.name === unit.icon)?.component ?? fallbackIcon.component"
                class="size-5"
                :fill="unit.color ?? 'currentColor'"
            />
        </template>

        <template #homePostal-data="{ row: unit }">
            {{ unit.homePostal ?? $t('common.na') }}
        </template>

        <template #actions-data="{ row: unit }">
            <div :key="unit.id" class="flex items-center">
                <UTooltip v-if="can('centrum.CentrumService/CreateOrUpdateUnit').value" :text="$t('common.update')">
                    <UButton
                        variant="link"
                        icon="i-mdi-pencil"
                        @click="
                            modal.open(UnitCreateOrUpdateModal, {
                                unit: unit,
                                onUpdated: async () => refresh(),
                            })
                        "
                    />
                </UTooltip>

                <UTooltip v-if="can('centrum.CentrumService/DeleteUnit').value" :text="$t('common.delete')">
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteUnit(unit.id),
                            })
                        "
                    />
                </UTooltip>
            </div>
        </template>
    </UTable>
</template>
