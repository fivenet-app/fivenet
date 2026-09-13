<script lang="ts" setup>
import NotificationPreferences from '~/components/user-settings/NotificationPreferences.vue';
import NotificationSounds from '~/components/user-settings/NotificationSounds.vue';
import { reminderTimes } from '~/types/calendar';

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { calendar } = storeToRefs(settingsStore);
const calendarReminderTimes = computed(() => [
    ...reminderTimes.map((n) => ({
        label:
            n === 0
                ? t('components.auth.user_settings.calendar_notifications.reminder_times.start')
                : `${n / 60} ${t('common.time_ago.minute', n / 60)}`,
        value: n,
    })),
]);
</script>

<template>
    <div class="space-y-4">
        <NotificationPreferences />

        <UPageCard
            :title="$t('components.auth.user_settings.calendar_notifications.title')"
            :description="$t('components.auth.user_settings.calendar_notifications.description')"
        >
            <UFormField
                class="grid grid-cols-2 items-center gap-2"
                name="calendarNotifications"
                :label="$t('components.auth.user_settings.calendar_notifications.reminder_times.name')"
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
                                    : $t('common.none_selected')
                            }}
                        </template>
                    </USelectMenu>
                </ClientOnly>
            </UFormField>
        </UPageCard>

        <NotificationSounds />
    </div>
</template>
