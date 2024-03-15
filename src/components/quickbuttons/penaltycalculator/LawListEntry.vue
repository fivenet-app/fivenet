<script lang="ts" setup>
import { type SelectedPenalty } from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';
import { type Law } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    law: Law;
    count: number;
}>();

defineEmits<{
    (e: 'selected', p: SelectedPenalty): void;
}>();

const count = ref(props.count);
</script>

<template>
    <tr class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            {{ law.name }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">${{ law.fine }}</td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            {{ law.detentionTime }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-accent-200">
            {{ law.stvoPoints }}
        </td>
        <td class="break-all px-1 py-1 text-left text-sm text-accent-200">
            {{ law.description }}
        </td>
        <td class="w-20 min-w-20 px-1 py-1 text-left text-accent-200">
            <select
                v-model="count"
                name="count"
                class="mb-1 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-primary-600 sm:text-sm sm:leading-6"
                @change="$emit('selected', { law: law, count: count })"
                @focusin="focusTablet(true)"
                @focusout="focusTablet(false)"
            >
                <option v-for="(_, i) in 7" :key="i" :value="i">
                    {{ i }}
                </option>
            </select>
        </td>
    </tr>
</template>
