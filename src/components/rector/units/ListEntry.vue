<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiPencil, mdiTrashCan } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import CreateOrUpdateUnitModal from './CreateOrUpdateUnitModal.vue';

const props = defineProps<{
    unit: Unit;
}>();

const { $grpc } = useNuxtApp();

async function deleteUnit(): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getUnitClient().deleteUnit({
                unitId: props.unit.id,
            });
            await call;

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const open = ref(false);
</script>

<template>
    <CreateOrUpdateUnitModal :unit="unit" :open="open" @close="open = false" />
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
                <button @click="open = true" class="flex-initial text-primary-500 hover:text-primary-400">
                    <SvgIcon class="h-6 w-6 text-primary-500" aria-hidden="true" type="mdi" :path="mdiPencil" />
                </button>
                <button
                    v-can="'UnitService.DeleteUnit'"
                    @click="deleteUnit"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <SvgIcon class="h-6 w-6 text-primary-500" aria-hidden="true" type="mdi" :path="mdiTrashCan" />
                </button>
            </div>
        </td>
    </tr>
</template>
