<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { useConfirmDialog } from '@vueuse/core';
import { PencilIcon, TrashCanIcon } from 'mdi-vue3';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import CreateOrUpdateUnitModal from './CreateOrUpdateUnitModal.vue';

const props = defineProps<{
    unit: Unit;
}>();

const { $grpc } = useNuxtApp();

async function deleteUnit(id: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().deleteUnit({
                unitId: id,
            });
            await call;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteUnit(id));

const open = ref(false);
</script>

<template>
    <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(unit.id)" />

    <CreateOrUpdateUnitModal v-if="can('CentrumService.CreateOrUpdateUnit')" :unit="unit" :open="open" @close="open = false" />

    <tr :key="unit.id?.toString()">
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
            <input type="color" disabled :value="`#${unit.color ?? 'ffffff'}`" class="h-6" />
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
            <div class="flex flex-row justify-end">
                <button
                    v-if="can('CentrumService.CreateOrUpdateUnit')"
                    @click="open = true"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <PencilIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                </button>
                <button
                    v-if="can('CentrumService.DeleteUnit')"
                    @click="reveal(props.unit.id)"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <TrashCanIcon class="h-6 w-6 text-primary-500" aria-hidden="true" />
                </button>
            </div>
        </td>
    </tr>
</template>
