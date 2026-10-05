<script setup lang="ts">
import type { CalendarDate, Time } from '@internationalized/date';
import InputTimePicker from './InputTimePicker.vue';

const props = defineProps<{
    modelValue: Date | undefined;
    permanent: boolean;
}>();

const emits = defineEmits<{
    (e: 'update:modelValue', value: Date | undefined): void;
    (e: 'update:permanent', value: boolean): void;
}>();

const { d, t } = useI18n();

const quickPicks = [
    { key: '30_minutes', seconds: 30 * 60 },
    { key: '1_hour', seconds: 60 * 60 },
    { key: '2_hours', seconds: 2 * 60 * 60 },
    { key: '1_day', seconds: 24 * 60 * 60 },
];

const calendarValue = computed<CalendarDate | undefined>({
    get: () => (props.modelValue ? dateToCalendarDate(props.modelValue) : undefined),
    set: (value) => {
        if (!value) return;

        const date = calendarDateToDate(value);
        const time = props.modelValue ? dateToTime(props.modelValue) : undefined;
        date.setHours(time?.hour ?? 0, time?.minute ?? 0, time?.second ?? 0, time?.millisecond ?? 0);
        emits('update:modelValue', date);
        emits('update:permanent', false);
    },
});

const timeValue = computed<Time | undefined>({
    get: () => (props.modelValue ? dateToTime(props.modelValue) : undefined),
    set: (value) => {
        if (!value) return;

        const date = props.modelValue ? new Date(props.modelValue) : new Date();
        date.setHours(value.hour, value.minute, value.second, value.millisecond);
        emits('update:modelValue', date);
        emits('update:permanent', false);
    },
});

function selectQuickPick(seconds: number): void {
    emits('update:modelValue', new Date(Date.now() + seconds * 1000));
    emits('update:permanent', false);
}

function selectPermanent(): void {
    emits('update:modelValue', undefined);
    emits('update:permanent', true);
}
</script>

<template>
    <UPopover :content="{ align: 'start' }" :modal="true">
        <UButton class="group w-full justify-between" color="neutral" variant="subtle">
            <span class="truncate">
                {{ permanent ? $t('common.permanent') : modelValue ? d(modelValue, 'short') : $t('common.pick_date') }}
            </span>
            <template #trailing>
                <UIcon
                    class="size-5 shrink-0 text-dimmed transition-transform duration-200 group-data-[state=open]:rotate-180"
                    name="i-mdi-calendar-clock"
                />
            </template>
        </UButton>

        <template #content>
            <div class="flex items-stretch divide-default sm:divide-x">
                <div class="hidden flex-col justify-center sm:flex">
                    <UButton
                        v-for="quickPick in quickPicks"
                        :key="quickPick.key"
                        class="rounded-none px-4"
                        :label="t(`common.wanted_duration.${quickPick.key}`)"
                        color="neutral"
                        variant="ghost"
                        truncate
                        @click="selectQuickPick(quickPick.seconds)"
                    />
                    <UButton
                        class="rounded-none px-4"
                        :label="$t('common.permanent')"
                        color="neutral"
                        variant="ghost"
                        truncate
                        @click="selectPermanent"
                    />
                </div>

                <div class="space-y-2 p-2">
                    <UCalendar v-model="calendarValue" />
                    <InputTimePicker v-model="timeValue" class="w-full" :hour-cycle="24" />
                </div>
            </div>
        </template>
    </UPopover>
</template>
