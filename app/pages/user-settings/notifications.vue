<script lang="ts" setup>
import type { SelectMenuItem } from '@nuxt/ui';
import { deleteSound, putSound } from '~/composables/useSounds';
import { notificationCategoryToIcon } from '~/components/notifications/helpers';
import { reminderTimes } from '~/types/calendar';
import { getNotificationsNotificationsClient } from '~~/gen/ts/clients';
import {
    NotificationCategory,
    NotificationKind,
    type NotificationPreference,
} from '~~/gen/ts/resources/notifications/notifications';

const { t } = useI18n();

const settingsStore = useSettingsStore();
const { audio, calendar } = storeToRefs(settingsStore);
const { activeChar } = useAuth();
const notificationsClient = await getNotificationsNotificationsClient();

const notificationCategories = [
    NotificationCategory.GENERAL,
    NotificationCategory.DOCUMENT,
    NotificationCategory.CALENDAR,
    NotificationCategory.JOBS,
    NotificationCategory.QUALIFICATIONS,
    NotificationCategory.MAILER,
    NotificationCategory.SYSTEM,
];
const notificationDeliveryScopes = [
    { category: NotificationCategory.UNSPECIFIED, labelKey: 'components.auth.user_settings.notification_delivery.all' },
    ...notificationCategories.map((category) => ({
        category,
        labelKey: `enums.notifications.NotificationCategory.${NotificationCategory[category]}`,
    })),
];

type DeliverySetting = 'inboxEnabled' | 'toastEnabled' | 'soundEnabled';
const pendingPreferenceCategories = ref<Set<NotificationCategory>>(new Set());

const { data: notificationPreferences } = useAuthedLazyAsyncData(
    'userState',
    'notification-preferences',
    async ({ signal }) => (await notificationsClient.getNotificationPreferences({}, { abort: signal })).response.preferences,
    {
        default: () => [],
        enabled: computed(() => activeChar.value !== null),
    },
);

function notificationPreference(category: NotificationCategory): NotificationPreference | undefined {
    return notificationPreferences.value.find(
        (preference) => preference.category === category && preference.kind === NotificationKind.UNSPECIFIED,
    );
}

function deliverySetting(category: NotificationCategory, setting: DeliverySetting): boolean {
    return (
        notificationPreference(category)?.[setting] ??
        notificationPreference(NotificationCategory.UNSPECIFIED)?.[setting] ??
        true
    );
}

function isPreferencePending(category: NotificationCategory): boolean {
    return pendingPreferenceCategories.value.has(category);
}

async function updateDeliverySetting(category: NotificationCategory, setting: DeliverySetting, value: boolean): Promise<void> {
    if (isPreferencePending(category)) return;

    pendingPreferenceCategories.value = new Set(pendingPreferenceCategories.value).add(category);
    const current = notificationPreference(category);
    const preference: NotificationPreference = {
        category,
        kind: NotificationKind.UNSPECIFIED,
        ...current,
        [setting]: value,
    };

    try {
        const { response } = await notificationsClient.updateNotificationPreference({ preference });
        notificationPreferences.value = response.preferences;
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        const pending = new Set(pendingPreferenceCategories.value);
        pending.delete(category);
        pendingPreferenceCategories.value = pending;
    }
}

async function resetDeliveryCategory(category: NotificationCategory): Promise<void> {
    if (isPreferencePending(category)) return;

    pendingPreferenceCategories.value = new Set(pendingPreferenceCategories.value).add(category);
    try {
        const { response } = await notificationsClient.updateNotificationPreference({
            preference: { category, kind: NotificationKind.UNSPECIFIED },
            reset: true,
        });
        notificationPreferences.value = response.preferences;
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        const pending = new Set(pendingPreferenceCategories.value);
        pending.delete(category);
        pendingPreferenceCategories.value = pending;
    }
}

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
            name: 'dispatch.attention',
            label: t('components.auth.user_settings.sounds.dispatch.attention'),
            items: [soundDefaultItem, ...soundsBaseItems.value],
        },
        {
            name: 'dispatch.dispatchSOS',
            label: t('components.auth.user_settings.sounds.dispatch.dispatch_sos'),
            items: [soundDefaultItem, ...soundsBaseItems.value],
        },
        {
            name: 'dispatch.dispatchAssigned',
            label: t('components.auth.user_settings.sounds.dispatch.dispatch_assigned'),
            items: [soundDefaultItem, ...soundsBaseItems.value],
        },
        {
            name: 'dispatch.dispatchCompleted',
            label: t('components.auth.user_settings.sounds.dispatch.dispatch_completed'),
            items: [...soundsBaseItems.value],
        },
    ],
]);
</script>

<template>
    <div class="space-y-4">
        <UPageCard
            v-if="activeChar"
            :title="$t('components.auth.user_settings.notification_delivery.title')"
            :description="$t('components.auth.user_settings.notification_delivery.description')"
        >
            <div class="space-y-2">
                <div
                    v-for="scope in notificationDeliveryScopes"
                    :key="scope.category"
                    class="grid gap-2 rounded-lg border border-default p-2 lg:grid-cols-[minmax(12rem,1fr)_auto_auto_auto_auto] lg:items-center"
                >
                    <div class="flex items-center gap-2">
                        <UIcon :name="notificationCategoryToIcon(scope.category)" class="size-5 text-muted" />
                        <span class="font-medium">
                            {{ $t(scope.labelKey) }}
                        </span>
                    </div>

                    <div class="min-w-20">
                        <UTooltip v-if="notificationPreference(scope.category)" :text="$t('common.reset')">
                            <UButton
                                color="neutral"
                                variant="ghost"
                                icon="i-mdi-restore"
                                class="!size-4 !p-0"
                                :loading="isPreferencePending(scope.category)"
                                @click="resetDeliveryCategory(scope.category)"
                            />
                        </UTooltip>
                    </div>

                    <UCheckbox
                        :model-value="deliverySetting(scope.category, 'inboxEnabled')"
                        :label="$t('components.auth.user_settings.notification_delivery.inbox')"
                        :disabled="isPreferencePending(scope.category)"
                        @update:model-value="(value) => updateDeliverySetting(scope.category, 'inboxEnabled', value === true)"
                    />
                    <UCheckbox
                        :model-value="deliverySetting(scope.category, 'toastEnabled')"
                        :label="$t('components.auth.user_settings.notification_delivery.toast')"
                        :disabled="isPreferencePending(scope.category)"
                        @update:model-value="(value) => updateDeliverySetting(scope.category, 'toastEnabled', value === true)"
                    />
                    <UCheckbox
                        :model-value="deliverySetting(scope.category, 'soundEnabled')"
                        :label="$t('components.auth.user_settings.notification_delivery.sound')"
                        :disabled="isPreferencePending(scope.category)"
                        @update:model-value="(value) => updateDeliverySetting(scope.category, 'soundEnabled', value === true)"
                    />
                </div>
            </div>
        </UPageCard>

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
