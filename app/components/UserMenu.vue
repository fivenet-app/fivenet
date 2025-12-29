<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui';
import { useDashboard } from '@nuxt/ui/utils/dashboard';
import LanguageSwitcherModal from './partials/LanguageSwitcherModal.vue';

defineProps<{
    collapsed?: boolean;
}>();

const { activeChar, username } = useAuth();

const { t } = useI18n();

const overlay = useOverlay();

const { toggleSearch } = useDashboard({ toggleSearch: () => {} });

const languageSwitcherModal = overlay.create(LanguageSwitcherModal);

const items = computed<DropdownMenuItem[][]>(() => [
    [
        {
            slot: 'account' as const,
            label: '',
            disabled: true,
        },
    ],
    [
        {
            label: t('components.auth.AccountInfo.title'),
            icon: 'i-mdi-account-details-outline',
            to: '/auth/account-info',
        },
        {
            label: t('components.auth.user_settings_panel.title'),
            icon: 'i-mdi-account-cog-outline',
            to: '/user-settings',
        },
        {
            label: t('common.commandpalette'),
            icon: 'i-mdi-terminal',
            kbds: ['CTRL', 'K'],
            onClick: () => toggleSearch && toggleSearch(),
        },
        {
            label: t('components.language_switcher.title'),
            icon: 'i-mdi-translate',
            onClick: () => languageSwitcherModal.open(),
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
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

const { game } = useAppConfig();

const name = computed(() =>
    activeChar.value ? `${activeChar.value?.firstname} ${activeChar.value?.lastname}` : (username.value ?? t('common.na')),
);
</script>

<template>
    <UDropdownMenu
        :items="items"
        :content="{ align: 'center', collisionPadding: 12 }"
        :ui="{ content: collapsed ? 'w-48' : 'w-(--reka-dropdown-menu-trigger-width)' }"
    >
        <template #default>
            <UButton
                :label="collapsed ? undefined : name"
                color="neutral"
                variant="ghost"
                block
                :square="collapsed"
                class="data-[state=open]:bg-elevated"
                :trailing-icon="collapsed ? undefined : 'i-mdi-ellipsis-vertical'"
                :ui="{
                    trailingIcon: 'text-dimmed',
                }"
            >
                <template #leading>
                    <UAvatar
                        :src="activeChar?.profilePicture ? `/api/filestore/${activeChar.profilePicture}` : undefined"
                        :alt="getInitials(name)"
                        size="xs"
                    />
                </template>
            </UButton>
        </template>

        <template #account>
            <div class="min-w-0 truncate text-left">
                <p>{{ $t('components.UserDropdown.signed_in_as') }}</p>
                <p class="truncate font-medium text-highlighted">{{ username }}</p>
                <p v-if="activeChar" class="inline-flex items-center gap-1 font-medium text-highlighted">
                    <span>{{ activeChar.jobLabel }}</span>
                    <template v-if="activeChar.job !== game.unemployedJobName">
                        <span>-</span>
                        <span>{{ activeChar.jobGradeLabel }}</span>
                    </template>
                </p>
            </div>
        </template>
    </UDropdownMenu>
</template>
