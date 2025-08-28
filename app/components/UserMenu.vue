<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui';
import SuperuserJobSelection from '~/components/partials/SuperuserJobSelection.vue';
import { useAuthStore } from '~/stores/auth';
import LanguageSwitcherModal from './partials/LanguageSwitcherModal.vue';

defineProps<{
    collapsed?: boolean;
}>();

const { isDashboardSearchModalOpen } = useDashboard();

const { can, activeChar, username, isSuperuser } = useAuth();

const { t } = useI18n();

const authStore = useAuthStore();

const overlay = useOverlay();

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
            label: t('components.auth.UserSettingsPanel.title'),
            icon: 'i-mdi-account-cog-outline',
            to: '/user-settings',
        },
        {
            label: t('common.commandpalette'),
            icon: 'i-mdi-terminal',
            shortcuts: ['CTRL', 'K'],
            onClick: () => (isDashboardSearchModalOpen.value = true),
        },
        can(['Superuser/CanBeSuperuser', 'Superuser/Superuser']).value
            ? {
                  label: `${t('common.superuser')}`,
                  icon: 'i-mdi-square-root',
                  type: 'checkbox' as const,
                  checked: isSuperuser.value,
                  onUpdateChecked(_: boolean) {
                      authStore.setSuperuserMode(!isSuperuser.value);
                  },
              }
            : undefined,
        isSuperuser.value
            ? {
                  slot: 'job' as const,
                  label: 'Select Job',
                  icon: 'i-mdi-briefcase',
                  onClick: ($event: Event) => $event.preventDefault(),
              }
            : undefined,
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
        <UChip color="error" position="top-left" :show="isSuperuser" class="w-full flex-1">
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
                        :alt="name"
                        size="xs"
                        :ui="{ rounded: 'rounded-full' }"
                    />
                </template>
            </UButton>
        </UChip>

        <template #account>
            <div class="truncate text-left">
                <p>{{ $t('components.UserDropdown.signed_in_as') }}</p>
                <p class="truncate font-medium text-highlighted">{{ username }}</p>
                <p v-if="activeChar" class="truncate font-medium text-highlighted">
                    {{ activeChar.jobLabel
                    }}<template v-if="activeChar.job !== game.unemployedJobName"> - {{ activeChar.jobGradeLabel }}</template>
                </p>
            </div>
        </template>

        <template v-if="isSuperuser" #job>
            <SuperuserJobSelection />
        </template>
    </UDropdownMenu>
</template>
