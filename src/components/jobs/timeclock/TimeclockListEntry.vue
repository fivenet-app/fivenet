<script lang="ts" setup>
import { CalendarIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { TimeclockEntry } from '~~/gen/ts/resources/jobs/timeclock';

const props = defineProps<{
    entry: TimeclockEntry;
    first?: Date;
    showDate: boolean;
}>();

// Soooo math with a computer is pretty cool, right? Right guys? Prevent ".349999999" seconds from appearing
const spentTime = parseFloat(((Math.round(props.entry.spentTime * 100) / 100) * 60 * 60).toPrecision(2));
</script>

<template>
    <tr class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <td
            v-if="showDate"
            class="inline-flex items-center whitespace-nowrap py-1 pl-4 pr-3 text-base font-medium text-neutral sm:pl-1"
        >
            <template v-if="first">
                <CalendarIcon class="h-5 w-5 pr-2" />
                {{ $d(first, 'date') }}
            </template>
        </td>
        <td v-if="!showDate" class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            <CitizenInfoPopover :user="entry.user" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ entry.spentTime > 0 ? fromSecondsToFormattedDuration(spentTime, { seconds: false }) : '' }}
            <template v-if="entry.startTime !== undefined">
                <span
                    class="inline-flex items-center rounded-md bg-success-500/10 px-2 py-1 text-xs font-medium text-success-400 ring-1 ring-inset ring-success-500/20"
                >
                    {{ $t('common.active') }}
                </span>
            </template>
        </td>
    </tr>
</template>
