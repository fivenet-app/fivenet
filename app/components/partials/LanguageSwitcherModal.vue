<script lang="ts" setup>
import type { LocaleObject } from '@nuxtjs/i18n';
import { useSettingsStore } from '~/stores/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { locale, locales } = useI18n();

const notifications = useNotificationsStore();

const settingsStore = useSettingsStore();
const { locale: userLocale } = storeToRefs(settingsStore);

const languages = ref<LocaleObject[]>([]);

onBeforeMount(async () => {
    locales.value.forEach((lang) => {
        if (typeof lang === 'string') return;

        languages.value.push({
            code: lang.code,
            language: lang.language!,
            name: lang.name!,
            icon: lang.icon ?? 'i-mdi-question',
        });
    });
});

const preventClose = ref(false);

async function switchLanguage(lang: LocaleObject): Promise<void> {
    if (locale.value === lang.code) return;

    preventClose.value = true;
    useLogger('⚙️ Settings').info('Switching language to:', lang.code);

    userLocale.value = lang.code;

    notifications.add({
        title: { key: 'notifications.language_switched.title', parameters: {} },
        description: { key: 'notifications.language_switched.content', parameters: { name: lang.name ?? lang.code } },
        type: NotificationType.SUCCESS,
        duration: 1650,
        callback: () => reloadNuxtApp({ persistState: false, force: true }),
    });
}
</script>

<template>
    <UModal :title="$t('components.language_switcher.title')" :dismissible="!preventClose">
        <template #body>
            <UPageGrid>
                <UPageCard
                    v-for="item in languages"
                    :key="item.name"
                    :title="item.name"
                    :icon="item.icon"
                    :ui="{ leadingIcon: 'size-12' }"
                    @click="switchLanguage(item)"
                />
            </UPageGrid>
        </template>

        <template #footer>
            <UButton
                class="flex-1"
                block
                color="neutral"
                :disabled="preventClose"
                :label="$t('common.close', 1)"
                @click="$emit('close', false)"
            />
        </template>
    </UModal>
</template>
