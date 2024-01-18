<script setup lang="ts">
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue';
import { TranslateIcon } from 'mdi-vue3';
import { type LocaleObject } from 'vue-i18n-routing';
import { useNotificatorStore } from '~/store/notificator';

const { locale, locales } = useI18n();
const notifications = useNotificatorStore();

type Language = { name: string; iso: string };
const languages = ref<Language[]>([]);

onMounted(async () => {
    locales.value.forEach((lang) => {
        if (typeof lang === 'string') {
            return;
        }

        lang = lang as LocaleObject;
        languages.value.push({
            name: lang.name!,
            iso: lang.iso!,
        });
    });
});

async function switchLanguage(lang: Language): Promise<void> {
    console.debug('Switching language to:', lang);

    locale.value = lang.iso;

    notifications.dispatchNotification({
        title: { key: 'notifications.language_switched.title', parameters: {} },
        content: { key: 'notifications.language_switched.content', parameters: { name: lang.name } },
        type: 'success',
    });
}
</script>

<template>
    <Menu as="div" class="relative flex-shrink-0">
        <div>
            <MenuButton
                class="flex rounded-full bg-base-800 text-sm ring-2 ring-base-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
            >
                <span class="sr-only">{{ $t('components.partials.sidebar_language_switcher.open_switcher') }}</span>
                <TranslateIcon
                    class="h-10 w-auto rounded-full bg-base-800 p-1 text-base-300 hover:text-base-100 hover:transition-colors"
                />
            </MenuButton>
        </div>
        <transition
            enter-active-class="transition duration-100 ease-out"
            enter-from-class="transform scale-95 opacity-0"
            enter-to-class="transform scale-100 opacity-100"
            leave-active-class="transition duration-75 ease-in"
            leave-from-class="transform scale-100 opacity-100"
            leave-to-class="transform scale-95 opacity-0"
        >
            <MenuItems
                class="absolute right-0 z-40 mt-2 w-48 origin-top-right rounded-md bg-base-800 py-1 shadow-float ring-1 ring-base-100 ring-opacity-5 focus:outline-none"
            >
                <MenuItem v-for="item in languages" :key="item.iso" v-slot="{ active }">
                    <button
                        :class="[
                            active ? 'bg-primary-500' : '',
                            'flex w-full flex-row px-4 py-2 text-sm text-neutral hover:transition-colors',
                        ]"
                        @click="switchLanguage(item)"
                    >
                        {{ item.name }}
                    </button>
                </MenuItem>
            </MenuItems>
        </transition>
    </Menu>
</template>
