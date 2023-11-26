<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useConfirmDialog } from '@vueuse/core';
import { PencilIcon, TrashCanIcon } from 'mdi-vue3';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { Unit } from '~~/gen/ts/resources/centrum/units';
import CreateOrUpdateUnitModal from '~/components/centrum/settings/CreateOrUpdateUnitModal.vue';

const props = defineProps<{
    unit: Unit;
}>();

const emits = defineEmits<{
    (e: 'updated'): void;
}>();

const { $grpc } = useNuxtApp();

async function deleteUnit(id: string): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().deleteUnit({
            unitId: id,
        });
        await call;

        emits('updated');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const color = ref(`#${props.unit.color ?? 'ffffff'}`);

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteUnit(id));

const open = ref(false);
</script>

<template>
    <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(unit.id)" />
    <CreateOrUpdateUnitModal
        v-if="can('CentrumService.CreateOrUpdateUnit')"
        :unit="unit"
        :open="open"
        @close="open = false"
        @updated="$emit('updated')"
    />

    <tr>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ unit.name }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ unit.initials }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ unit.description }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
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
                {{ $t('common.none_selected', [$t('common.attributes', 2)]) }}
            </template>
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            <ColorInput v-model="color" disabled format="hex" class="h-6" />
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div class="flex flex-row justify-end">
                <button
                    v-if="can('CentrumService.CreateOrUpdateUnit')"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                    @click="open = true"
                >
                    <PencilIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                </button>
                <button
                    v-if="can('CentrumService.DeleteUnit')"
                    type="button"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                    @click="reveal(unit.id)"
                >
                    <TrashCanIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                </button>
            </div>
        </td>
    </tr>
</template>
