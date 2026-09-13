<script lang="ts" setup>
import type { SelectMenuItem } from '@nuxt/ui';
import { deleteSound, putSound, useSounds, type NotificationSound, type SoundKeys } from '~/composables/useSounds';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const settingsStore = useSettingsStore();
const { audio } = storeToRefs(settingsStore);
const { t } = useI18n();

const notifications = useNotificationsStore();

function notifySoundUpdate(success: boolean): void {
    const key = success ? 'success' : 'failed';
    notifications.add({
        title: { key: `components.auth.user_settings.sounds.custom_sound_${key}.title`, parameters: {} },
        description: { key: `components.auth.user_settings.sounds.custom_sound_${key}.content`, parameters: {} },
        type: success ? NotificationType.SUCCESS : NotificationType.ERROR,
    });
}

async function updateSound(key: SoundKeys, value: NotificationSound): Promise<void> {
    if (value.value === 'custom') {
        await uploadCustomSound(key);
        return;
    }

    try {
        await deleteSound(key);
        audio.value.sounds[key] = value;
        notifySoundUpdate(true);
    } catch {
        notifySoundUpdate(false);
    }
}

async function uploadCustomSound(key: SoundKeys): Promise<void> {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'audio/mpeg,audio/mp3,audio/ogg,audio/wav,audio/x-wav,audio/aac';

    input.onchange = async (event) => {
        const file = (event.target as HTMLInputElement).files?.[0];
        if (file) {
            const blob = new Blob([file], { type: file.type });
            try {
                await putSound(key, blob, file.name);
                audio.value.sounds[key] = { value: 'custom', custom: file.name };
                notifySoundUpdate(true);
            } catch {
                notifySoundUpdate(false);
            }
        }
    };

    input.click();
}

const notificationSound = useSounds('notification');
const soundDefaultItem = computed(() => ({ label: t('common.default'), value: { value: 'default' } }));

const soundsBaseItems = computed<SelectMenuItem[]>(() => [
    { label: t('components.auth.user_settings.sounds.custom_sound'), value: { value: 'custom' } },
    { label: t('common.disabled'), value: { value: 'none' } },
]);

const sounds = computed<Array<Array<{ name: SoundKeys; label: string; items: SelectMenuItem[] }>>>(() => [
    [
        {
            name: 'notification',
            label: t('components.auth.user_settings.sounds.notification'),
            items: [soundDefaultItem.value, ...soundsBaseItems.value],
        },
    ],
    [
        {
            name: 'dispatch.attention',
            label: t('components.auth.user_settings.sounds.dispatch.attention'),
            items: [soundDefaultItem.value, ...soundsBaseItems.value],
        },
        {
            name: 'dispatch.dispatchSOS',
            label: t('components.auth.user_settings.sounds.dispatch.dispatch_sos'),
            items: [soundDefaultItem.value, ...soundsBaseItems.value],
        },
        {
            name: 'dispatch.dispatchAssigned',
            label: t('components.auth.user_settings.sounds.dispatch.dispatch_assigned'),
            items: [soundDefaultItem.value, ...soundsBaseItems.value],
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
                            {{ audio.sounds[sound.name]?.custom || $t('components.auth.user_settings.sounds.custom_sound') }}
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
</template>
