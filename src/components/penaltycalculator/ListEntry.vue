<script lang="ts" setup>
import { SelectedPenalty } from '~/components/penaltycalculator/Calculator.vue';
import { Law } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    law: Law;
    count: bigint;
}>();

defineEmits<{
    (e: 'selected', p: SelectedPenalty): void;
}>();

const count = ref(props.count);
</script>

<template>
    <tr>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ law.name }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">${{ law.fine }}</td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ law.detentionTime }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ law.stvoPoints }}
        </td>
        <td class="break-all text-sm px-1 py-1 text-left text-base-200">
            {{ law.description }}
        </td>
        <td class="w-20 min-w-20 px-1 py-1 text-left text-base-200">
            <select
                name="count"
                v-model="count"
                @change="$emit('selected', { law: law, count: BigInt(count) })"
                class="mb-1 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-primary-600 sm:text-sm sm:leading-6"
            >
                <option v-for="(_, i) in 7">
                    {{ i }}
                </option>
            </select>
        </td>
    </tr>
</template>
