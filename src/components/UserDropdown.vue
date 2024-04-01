<script setup lang="ts">
import { useAuthStore } from '~/store/auth';

const { isHelpSlideoverOpen } = useDashboard();
const { isDashboardSearchModalOpen } = useUIState();
const { metaSymbol } = useShortcuts();

const { t } = useI18n();

const authStore = useAuthStore();
const { activeChar, username } = storeToRefs(authStore);

const items = computed(() => [
    [
        {
            slot: 'account',
            label: '',
            disabled: true,
        },
    ],
    [
        {
            label: t('common.setting', 2),
            icon: 'i-heroicons-cog-8-tooth',
            to: '/settings',
        },
        {
            label: t('common.command_palette'),
            icon: 'i-heroicons-command-line',
            shortcuts: [metaSymbol.value, 'K'],
            click: () => {
                isDashboardSearchModalOpen.value = true;
            },
        },
        {
            label: 'Help',
            icon: 'i-heroicons-question-mark-circle',
            shortcuts: ['?'],
            click: () => (isHelpSlideoverOpen.value = true),
        },
    ],
    [
        {
            label: 'Documentation',
            icon: 'i-heroicons-book-open',
            to: 'https://ui.nuxt.com/pro/getting-started',
            target: '_blank',
        },
        {
            label: 'GitHub repository',
            icon: 'i-simple-icons-github',
            to: 'https://github.com/nuxt-ui-pro/dashboard',
            target: '_blank',
        },
    ],
    [
        {
            label: t('components.partials.sidebar.change_character'),
            icon: 'i-mdi-account-switch',
            to: '/auth/character-selector',
        },
        {
            label: t('common.sign_out'),
            icon: 'i-mdi-logout',
            to: '/auth/logout',
        },
    ],
]);
</script>

<template>
    <UDropdown
        mode="hover"
        :items="items"
        :ui="{ width: 'w-full', item: { disabled: 'cursor-text select-text' } }"
        :popper="{ strategy: 'absolute', placement: 'top' }"
        class="w-full"
    >
        <template #default="{ open }">
            <UButton
                color="gray"
                variant="ghost"
                class="w-full"
                :label="activeChar ? `${activeChar?.firstname} ${activeChar?.lastname}` : username"
                :class="[open && 'bg-gray-50 dark:bg-gray-800']"
            >
                <template #leading>
                    <UAvatar :src="activeChar?.avatar?.url" size="2xs" />
                </template>

                <template #trailing>
                    <UIcon name="i-heroicons-ellipsis-vertical" class="ml-auto size-5" />
                </template>
            </UButton>
        </template>

        <template #account>
            <div class="text-left">
                <p>Signed in as</p>
                <p class="truncate font-medium text-gray-900 dark:text-white">{{ username }}</p>
            </div>
        </template>
    </UDropdown>
</template>
