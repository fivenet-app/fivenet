<script lang="ts" setup>
import { SelectedPenalty } from '~/utils/penalty';
import { Law } from '~~/gen/ts/resources/laws/laws';

defineProps<{
    law: Law;
}>();

defineEmits<{
    (e: 'selected', p: SelectedPenalty): void;
}>();

const count = ref(0);
</script>

<template>
    <tr>
        <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ law.name }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">${{ law.fine }}</td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ law.detentionTime }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ law.stvoPoints }}
        </td>
        <td class="px-2 py-2 text-sm text-base-200">
            {{ law.description }}
        </td>
        <td class="w-20 min-w-20 px-2 py-2 text-sm text-base-200">
            <select
                name="count"
                v-model="count"
                @change="$emit('selected', { law, count })"
                class="mb-1 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-primary-600 sm:text-sm sm:leading-6"
            >
                <option v-for="(_, i) in 7">
                    {{ i }}
                </option>
            </select>
        </td>
    </tr>
</template>
