<script lang="ts" setup>
import { type LocaleObject } from 'vue-i18n-routing';
import { useNotificatorStore } from '~/store/notificator';

const { locale, locales } = useI18n();

const notifications = useNotificatorStore();

const modal = useModal();

const languages = ref<LocaleObject[]>([]);

onMounted(async () => {
    locales.value.forEach((lang) => {
        if (typeof lang === 'string') {
            return;
        }

        lang = lang as LocaleObject;
        languages.value.push({
            code: lang.code,
            name: lang.name!,
            iso: lang.iso!,
            icon: lang.icon ?? 'i-mdi-question',
        });
    });
});

async function switchLanguage(lang: LocaleObject): Promise<void> {
    console.debug('Switching language to:', lang.name);

    locale.value = lang.iso!;

    notifications.add({
        title: { key: 'notifications.language_switched.title', parameters: {} },
        description: { key: 'notifications.language_switched.content', parameters: { name: lang.name! } },
        type: 'success',
        timeout: 2000,
        callback: () => reloadNuxtApp({ persistState: false, force: true }),
    });
}
</script>

<template>
    <UModal>
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.language_switcher.title') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="modal.close()" />
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

            <div class="mt-5 gap-2 sm:mt-4 sm:flex">
                <UButton block @click="modal.close()">
                    {{ $t('common.close', 1) }}
                </UButton>
            </div>
        </UCard>
    </UModal>
</template>
