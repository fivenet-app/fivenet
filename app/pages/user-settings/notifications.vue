<script lang="ts" setup>
import type { SelectMenuItem } from '@nuxt/ui';
import { deleteSound, putSound } from '~/composables/useSounds';
import { reminderTimes } from '~/types/calendar';

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { audio, calendar } = storeToRefs(settingsStore);

async function updateSound(key: SoundKeys, value: NotificationSound): Promise<void> {
    if (value.value === 'custom') {
        // If "switching" to custom, open the file dialog
        await uploadCustomSound(key);
        return;
    } else {
        await deleteSound(key);
    }

    audio.value.sounds[key] = value;
}

async function uploadCustomSound(key: SoundKeys): Promise<void> {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'audio/mpeg,audio/mp3,audio/ogg,audio/wav,audio/x-wav,audio/aac';

    input.onchange = async (event) => {
        const file = (event.target as HTMLInputElement).files?.[0];
        if (file) {
            const blob = new Blob([file], { type: file.type });
            await putSound(key, blob, file.name);

            audio.value.sounds[key] = { value: 'custom', custom: file.name };
        }
    };

    input.click();
}

const notificationSound = useSounds('notification');

const calendarReminderTimes = computed(() => [
    ...reminderTimes.map((n) => ({
        label:
            n === 0
                ? t('components.auth.user_settings.calendar_notifications.reminder_times.start')
                : `${n / 60} ${t('common.time_ago.minute', n / 60)}`,
        value: n,
    })),
]);

const soundDefaultItem = { label: t('common.default'), value: { value: 'default' } };

const soundsBaseItems = computed<SelectMenuItem[]>(() => [
    { label: t('components.auth.user_settings.sounds.custom_sound'), value: { value: 'custom' } },
    { label: t('common.disabled'), value: { value: 'none' } },
]);

const sounds = computed<
    Array<
        Array<{
            name: SoundKeys;
            label: string;
            items: SelectMenuItem[];
        }>
    >
>(() => [
    [
        {
            name: 'notification',
            label: t('components.auth.user_settings.sounds.notification'),
            items: [soundDefaultItem, ...soundsBaseItems.value],
        },
    ],
    [
        {
            name: 'centrum.attention',
            label: t('components.auth.user_settings.sounds.centrum.attention'),
            items: [soundDefaultItem, ...soundsBaseItems.value],
        },
        {
            name: 'centrum.dispatchSOS',
            label: t('components.auth.user_settings.sounds.centrum.dispatch_sos'),
            items: [soundDefaultItem, ...soundsBaseItems.value],
        },
        {
            name: 'centrum.dispatchAssigned',
            label: t('components.auth.user_settings.sounds.centrum.dispatch_assigned'),
            items: [soundDefaultItem, ...soundsBaseItems.value],
        },
        {
            name: 'centrum.dispatchCompleted',
            label: t('components.auth.user_settings.sounds.centrum.dispatch_completed'),
            items: [...soundsBaseItems.value],
        },
    ],
]);
</script>

<template>
    <div class="space-y-4">
        <UPageCard
            :description="$t('components.auth.user_settings.volumes.subtitle')"
            :ui="{ body: 'w-full', wrapper: 'w-full', title: 'flex w-full flex-row' }"
        >
            <template #title>
                <span class="flex-1">
                    {{ $t('components.auth.user_settings.volumes.title') }}
                </span>

                <UButton icon="i-mdi-play" @click="notificationSound.play()" />
            </template>

            <UFormField
                class="grid grid-cols-2 items-center gap-2"
                name="notificationsVolume"
                :label="$t('components.auth.user_settings.volumes.notifications_volume')"
            >
                <USlider v-model="audio.notificationsVolume" :step="0.01" :min="0" :max="1" />
                <span> {{ audio.notificationsVolume <= 0 ? 0 : (audio.notificationsVolume * 100).toFixed(0) }}% </span>
            </UFormField>
        </UPageCard>

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

        <UPageCard
            :title="$t('components.auth.user_settings.sounds.title')"
            :description="$t('components.auth.user_settings.sounds.description')"
        >
            <template v-for="(category, idx) in sounds" :key="idx">
                <UFormField
                    v-for="sound in category"
                    :key="sound.label"
                    class="grid grid-cols-2 items-center gap-2"
                    :name="sound.name"
                    :label="sound.label"
                    :ui="{ container: 'flex flex-col gap-2' }"
                >
                    <UFieldGroup>
                        <USelectMenu
                            class="w-full"
                            :model-value="audio.sounds[sound.name]"
                            value-key="value"
                            :items="sound.items"
                            :ui="{ base: 'line-clamp-2' }"
                            @update:model-value="(value) => updateSound(sound.name, value)"
                        >
                            <template v-if="audio.sounds[sound.name]?.value === 'custom'" #default>
                                {{
                                    audio.sounds[sound.name]?.custom || $t('components.auth.user_settings.sounds.custom_sound')
                                }}
                            </template>
                        </USelectMenu>

                        <UButton
                            :disabled="audio.sounds[sound.name]?.value === 'none'"
                            icon="i-mdi-play"
                            variant="outline"
                            @click="() => useSounds(sound.name).play()"
                        />
                    </UFieldGroup>

                    <UButton
                        v-if="audio.sounds[sound.name]?.value === 'custom'"
                        block
                        icon="i-mdi-upload"
                        variant="subtle"
                        :label="$t('components.auth.user_settings.sounds.select_custom_sound')"
                        @click="() => uploadCustomSound(sound.name)"
                    />
                </UFormField>

                <USeparator v-if="idx + 1 < sounds.length" class="my-2" />
            </template>
        </UPageCard>
    </div>
</template>
