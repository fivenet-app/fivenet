<script lang="ts" setup>
import { Penalty, SelectedPenalty } from '~/utils/penalty';

defineProps<{
    penalty: Penalty,
}>();

defineEmits<{
    (e: 'selected', p: SelectedPenalty): void,
}>();

const count = ref(0);
</script>

<template>
    <tr>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ penalty.name }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            ${{ penalty.fine }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ penalty.detentionTime }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ penalty.stvoPoints }}
        </td>
        <td class="px-2 py-2 text-sm text-base-200">
            {{ penalty.description }}
        </td>
        <td class="w-20 min-w-20 px-2 py-2 text-sm text-base-200">
            <select id="count" name="count" v-model="count" @change="$emit('selected', { penalty, count })"
                class="mt-2 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6">
                <option v-for="(_, i) in 7">{{ i }}</option>
            </select>
        </td>
    </tr>
</template>
