<script lang="ts" setup>
import { reminderTimes } from '~/types/calendar';

const { t } = useI18n();

const settings = useSettingsStore();
const { audio, calendar } = storeToRefs(settings);

const notificationSound = useSounds('/sounds/notification.mp3');

const calendarReminderTimes = [
    { label: t('components.auth.UserSettingsPanel.calendar_notifications.reminder_times.start'), value: 0 },
    ...reminderTimes.map((n) => ({ label: `${n / 60} ${t('common.time_ago.minute', n / 60)}`, value: n })),
];
</script>

<template>
    <div class="space-y-4">
        <UPageCard
            :title="$t('components.auth.UserSettingsPanel.volumes.title')"
            :description="$t('components.auth.UserSettingsPanel.volumes.subtitle')"
        >
            <template #links>
                <UButton icon="i-mdi-play" @click="notificationSound.play()" />
            </template>

            <UFormField
                class="grid grid-cols-2 items-center gap-2"
                name="notificationsVolume"
                :label="$t('components.auth.UserSettingsPanel.volumes.notifications_volume')"
            >
                <USlider v-model="audio.notificationsVolume" :step="0.01" :min="0" :max="1" />
                {{ audio.notificationsVolume <= 0 ? 0 : (audio.notificationsVolume * 100).toFixed(0) }}%
            </UFormField>
        </UPageCard>

        <UPageCard
            :title="$t('components.auth.UserSettingsPanel.calendar_notifications.title')"
            :description="$t('components.auth.UserSettingsPanel.calendar_notifications.subtitle')"
        >
            <UFormField
                class="grid grid-cols-2 items-center gap-2"
                name="calendarNotifications"
                :label="$t('components.auth.UserSettingsPanel.calendar_notifications.reminder_times.name')"
            >
                <ClientOnly>
                    <USelectMenu v-model="calendar.reminderTimes" multiple :items="calendarReminderTimes" value-key="value">
                        <template #default>
                            {{
                                calendar.reminderTimes.length > 0
                                    ? [...calendar.reminderTimes]
                                          .sort()
                                          .map(
                                              (n) =>
                                                  calendarReminderTimes.find((rt) => rt.value === n)?.label ?? $t('common.na'),
                                          )
                                          .join(', ')
                                    : $t('common.none')
                            }}
                        </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormField>
        </UPageCard>
    </div>
</template>
