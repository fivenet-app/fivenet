<script setup lang="ts">
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue';
import { setLocale as veeValidateSetLocale } from '@vee-validate/i18n';
import { TranslateIcon } from 'mdi-vue3';
import { LocaleObject } from 'vue-i18n-routing';
import { useNotificationsStore } from '~/store/notifications';
import { useSettingsStore } from '~/store/settings';

const { locales, setLocale } = useI18n();
const settings = useSettingsStore();
const notifications = useNotificationsStore();

const languages: { name: string; iso: string }[] = [];

onMounted(async () => {
    locales.value.forEach((lang) => {
        lang = lang as LocaleObject;

        languages.push({
            name: lang.name!,
            iso: lang.iso!,
        });
    });
});

async function switchLanguage(lang: { name: string; iso: string }): Promise<void> {
    settings.setLocale(lang.iso);
    setLocale(lang.iso);
    veeValidateSetLocale(lang.iso);

    notifications.dispatchNotification({
        title: { key: 'notifications.language_switched.title', parameters: [] },
        content: { key: 'notifications.language_switched.content', parameters: [lang.name] },
        type: 'success',
    });
}
</script>

<template>
    <Menu as="div" class="relative flex-shrink-0">
        <div>
            <MenuButton
                class="flex text-sm rounded-full bg-base-800 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
            >
                <span class="sr-only">{{ $t('components.partials.sidebar_language_switcher.open_switcher') }}</span>
                <TranslateIcon
                    class="w-auto h-10 p-1 rounded-full hover:transition-colors text-base-300 bg-base-800 hover:text-base-100"
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
                class="absolute right-0 w-48 py-1 mt-2 origin-top-right rounded-md shadow-float bg-base-800 ring-1 ring-base-100 ring-opacity-5 focus:outline-none z-40"
            >
                <MenuItem v-for="item in languages" :key="item.iso" v-slot="{ active }">
                    <button
                        :class="[
                            active ? 'bg-base-800' : '',
                            'px-4 py-2 text-sm text-neutral hover:transition-colors flex flex-row w-full',
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
