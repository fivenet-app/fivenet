<script lang="ts" setup>
import { EyeIcon } from 'mdi-vue3';
import Time from '~/components/partials/elements/Time.vue';
import { CONDUCT_TYPE, ConductEntry } from '~~/gen/ts/resources/jobs/conduct';
import { conductTypesToBGColor, conductTypesToRingColor, conductTypesToTextColor } from './helpers';

defineProps<{
    conduct: ConductEntry;
}>();
</script>

<template>
    <tr>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            <Time :value="conduct.createdAt" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            <Time v-if="conduct.expiresAt" :value="conduct.expiresAt" />
            <span v-else> No Expiration. </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            <div
                class="rounded-md py-1 px-2 text-sm font-medium ring-1 ring-inset"
                :class="[
                    conductTypesToBGColor(conduct.type),
                    conductTypesToRingColor(conduct.type),
                    conductTypesToTextColor(conduct.type),
                ]"
            >
                {{ $t(`enums.jobs.CONDUCT_TYPE.${CONDUCT_TYPE[conduct.type ?? (0 as number)]}`) }}
            </div>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ conduct.message }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ conduct.targetUser?.firstname }}, {{ conduct.targetUser?.lastname }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ conduct.creator?.firstname }}, {{ conduct.creator?.lastname }}
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-sm font-medium sm:pr-0">
            <div class="flex flex-row justify-end">
                <NuxtLink
                    :to="{
                        name: 'citizens-id',
                        params: { id: conduct.id.toString() ?? 0 },
                    }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <EyeIcon class="w-6 h-auto ml-auto mr-2.5" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
