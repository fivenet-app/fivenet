<script lang="ts" setup>
import { addDays, subDays } from 'date-fns';
import { Timeline, type TimelineGroup, type TimelineItem, type TimelineMarker } from 'vue-timeline-chart';
import 'vue-timeline-chart/style.css';
import type { TimeclockEntry } from '~~/gen/ts/resources/jobs/timeclock';

const props = defineProps<{
    data: TimeclockEntry[];
    from: Date;
    to: Date;
}>();

const groups = computed<TimelineGroup[]>(() =>
    props.data
        .filter((d, idx, self) => self.findIndex((o) => o.userId === d.userId) === idx)
        .map((d) => ({
            id: d.userId.toString(),
            label: `${d.user?.firstname} ${d.user?.lastname}`,
        })),
);

const items = computed<TimelineItem[]>(() =>
    props.data.map((d) =>
        d.startTime && d.endTime
            ? {
                  type: 'range',
                  group: d.userId.toString(),
                  start: toDate(d.startTime).getTime(),
                  end: toDate(d.endTime).getTime(),
              }
            : d.startTime
              ? {
                    type: 'point',
                    group: d.userId.toString(),
                    start: toDate(d.startTime).getTime(),
                }
              : {
                    type: 'point',
                    group: d.userId.toString(),
                    start: toDate(d.date).getTime(),
                },
    ),
);

const markers = computed<TimelineMarker[]>(() =>
    [
        mouseHoverPosition.value
            ? ({
                  type: 'marker',
                  id: 'mousehover',
                  start: mouseHoverPosition.value,
              } as TimelineMarker)
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const mouseHoverPosition = ref<number | null>(null);

function onMousemoveTimeline({ time }: { time: number }) {
    mouseHoverPosition.value = time;
}
function onMouseleaveTimeline() {
    mouseHoverPosition.value = null;
}

const minViewportDuration = 60000 * 60 * 4;
</script>

<template>
    <div class="flex flex-1 flex-col divide-y divide-gray-300 dark:divide-gray-700">
        <div class="flex min-w-full gap-1 px-1">
            <p class="line-clamp-1 flex-1 text-left hover:line-clamp-none">
                <span class="font-semibold">{{ $t('common.date') }}:</span>

                {{
                    mouseHoverPosition
                        ? $d(new Date(mouseHoverPosition), 'long')
                        : $t('components.jobs.timeclock.timeline.hover')
                }}
            </p>

            <slot name="caption" />
        </div>

        <Timeline
            class="timeclockTimeline flex-1"
            :groups="groups"
            :items="items"
            :markers="markers"
            :viewport-min="subDays(from, 1).getTime()"
            :viewport-max="addDays(to, 1).getTime()"
            :min-viewport-duration="minViewportDuration"
            @mousemove-timeline="onMousemoveTimeline"
            @mouseleave-timeline="onMouseleaveTimeline"
        />
    </div>
</template>

<style scoped>
.timeclockTimeline :deep(.timeline) {
    height: 100%;

    .groups {
        height: 100%;
    }
}
</style>
