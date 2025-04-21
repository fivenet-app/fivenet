<script setup lang="ts">
import SuperUserJobSelection from '~/components/partials/SuperUserJobSelection.vue';
import { useAuthStore } from '~/stores/auth';
import LanguageSwitcherModal from './partials/LanguageSwitcherModal.vue';

defineProps<{
    collapsed?: boolean;
}>();

const { can, activeChar, username, isSuperuser } = useAuth();

const { t } = useI18n();

const authStore = useAuthStore();

const overlay = useOverlay();

const langSwitcherModal = overlay.create(LanguageSwitcherModal, {});

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
            label: t('components.auth.AccountInfo.title'),
            icon: 'i-mdi-account-cog-outline',
            to: '/auth/account-info',
        },
        {
            label: t('common.setting', 2),
            icon: 'i-mdi-cog-outline',
            to: '/settings',
        },
        {
            label: t('common.commandpalette'),
            icon: 'i-mdi-terminal',
            shortcuts: ['meta', 'K'],
            onClick: () => console.log('isDashboardSearchModalOpen'),
        },
        can(['CanBeSuper', 'SuperUser']).value
            ? {
                  label: `${t('common.superuser')}: ${isSuperuser.value ? t('common.enabled') : t('common.disabled')}`,
                  icon: 'i-mdi-square-root',
                  onClick: () => authStore.setSuperUserMode(!isSuperuser.value),
              }
            : undefined,
        isSuperuser.value
            ? {
                  slot: 'job',
                  label: 'Select Job',
                  icon: 'i-mdi-briefcase',
                  onClick: ($event: Event) => $event.preventDefault(),
              }
            : undefined,
        {
            label: t('components.language_switcher.title'),
            icon: 'i-mdi-translate',
            onClick: () => langSwitcherModal.open(),
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
        <UChip
            class="w-full"
            color="error"
            :text="$t('common.superuser')"
            position="top-left"
            :show="isSuperuser"
            :ui="{ base: 'top-0 left-1/2' }"
        >
            <UButton
                v-bind="{
                    label: collapsed ? undefined : name,
                }"
                color="neutral"
                variant="ghost"
                block
                :square="collapsed"
                class="data-[state=open]:bg-(--ui-bg-elevated)"
                :ui="{
                    trailingIcon: 'text-(--ui-text-dimmed)',
                }"
            >
                <template #leading>
                    <UAvatar :src="activeChar?.avatar?.url" :alt="$t('common.avatar')" :text="getInitials(name)" size="2xs" />
                </template>
            </UButton>
        </UChip>

        <template #account>
            <div class="truncate text-left">
                <p>{{ $t('components.UserDropdown.signed_in_as') }}</p>
                <p class="truncate font-medium text-neutral-900 dark:text-white">{{ username }}</p>
                <p v-if="activeChar" class="truncate font-medium text-neutral-900 dark:text-white">
                    {{ activeChar.jobLabel
                    }}<template v-if="activeChar.job !== game.unemployedJobName"> - {{ activeChar.jobGradeLabel }}</template>
                </p>
            </div>
        </template>

        <template v-if="isSuperuser" #job>
            <SuperUserJobSelection />
        </template>
    </UDropdownMenu>
</template>
