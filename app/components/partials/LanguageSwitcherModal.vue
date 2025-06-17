<script lang="ts" setup>
import type { LocaleObject } from '@nuxtjs/i18n';
import { useSettingsStore } from '~/stores/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const { locale, locales } = useI18n();

const { isOpen } = useModal();

const notifications = useNotificationsStore();

const settings = useSettingsStore();
const { locale: userLocale } = storeToRefs(settings);

const languages = ref<LocaleObject[]>([]);

onBeforeMount(async () => {
    locales.value.forEach((lang) => {
        if (typeof lang === 'string') {
            return;
        }

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
    if (locale.value === lang.code) {
        return;
    }

    preventClose.value = true;
    useLogger('⚙️ Settings').info('Switching language to:', lang.code);

    userLocale.value = lang.code;

    notifications.add({
        title: { key: 'notifications.language_switched.title', parameters: {} },
        description: { key: 'notifications.language_switched.content', parameters: { name: lang.name ?? lang.code } },
        type: NotificationType.SUCCESS,
        timeout: 1650,
        callback: () => reloadNuxtApp({ persistState: false, force: true }),
    });
}
</script>

<template>
    <UModal :prevent-close="preventClose">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.language_switcher.title') }}
                    </h3>

                    <UButton
                        class="-my-1"
                        :disabled="preventClose"
                        color="gray"
                        variant="ghost"
                        icon="i-mdi-window-close"
                        @click="isOpen = false"
                    />
                </div>
            </template>

            <UPageGrid>
                <UPageCard
                    v-for="item in languages"
                    :key="item.name"
                    :title="item.name"
                    :icon="item.icon"
                    @click="switchLanguage(item)"
                />
            </UPageGrid>

            <template #footer>
                <UButton class="flex-1" block color="black" :disabled="preventClose" @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
