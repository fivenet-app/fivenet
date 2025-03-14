<script setup lang="ts">
import SuperUserJobSelection from '~/components/partials/SuperUserJobSelection.vue';
import { useAuthStore } from '~/stores/auth';
import LanguageSwitcherModal from './partials/LanguageSwitcherModal.vue';

const { isDashboardSearchModalOpen } = useUIState();
const { metaSymbol } = useShortcuts();

const { can, activeChar, username, isSuperuser } = useAuth();

const { t } = useI18n();

const authStore = useAuthStore();

const modal = useModal();

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
            shortcuts: [metaSymbol.value, 'K'],
            click: () => (isDashboardSearchModalOpen.value = true),
        },
        can(['CanBeSuper', 'SuperUser']).value
            ? {
                  label: `${t('common.superuser')}: ${isSuperuser.value ? t('common.enabled') : t('common.disabled')}`,
                  icon: 'i-mdi-square-root',
                  click: () => authStore.setSuperUserMode(!isSuperuser.value),
              }
            : undefined,
        isSuperuser.value
            ? {
                  slot: 'job',
                  label: 'Select Job',
                  icon: 'i-mdi-briefcase',
                  click: ($event: Event) => $event.preventDefault(),
              }
            : undefined,
        {
            label: t('components.language_switcher.title'),
            icon: 'i-mdi-translate',
            click: () => modal.open(LanguageSwitcherModal, {}),
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

const open = ref(false);
</script>

<template>
    <UDropdown
        v-model:open="open"
        :items="items"
        :ui="{ width: 'w-full', item: { disabled: 'cursor-text select-text' } }"
        :popper="{ strategy: 'absolute', placement: 'top' }"
        class="w-full"
        mode="hover"
    >
        <UChip
            class="w-full"
            color="red"
            :text="$t('common.superuser')"
            position="top-left"
            :show="isSuperuser"
            :ui="{ base: 'top-0 left-1/2' }"
        >
            <UButton
                color="gray"
                variant="ghost"
                class="w-full"
                :label="name"
                :class="[open && 'bg-gray-50 dark:bg-gray-800']"
                @click="open = !open"
                @touchstart.passive="open = !open"
            >
                <template #leading>
                    <UAvatar :src="activeChar?.avatar?.url" :alt="$t('common.avatar')" :text="getInitials(name)" size="2xs" />
                </template>

                <template #trailing>
                    <UIcon name="i-mdi-ellipsis-vertical" class="ml-auto size-5" />
                </template>
            </UButton>
        </UChip>

        <template #account>
            <div class="truncate text-left">
                <p>{{ $t('components.UserDropdown.signed_in_as') }}</p>
                <p class="truncate font-medium text-gray-900 dark:text-white">{{ username }}</p>
                <p v-if="activeChar" class="truncate font-medium text-gray-900 dark:text-white">
                    {{ activeChar.jobLabel
                    }}<template v-if="activeChar.job !== game.unemployedJobName"> - {{ activeChar.jobGradeLabel }}</template>
                </p>
            </div>
        </template>

        <template v-if="isSuperuser" #job>
            <SuperUserJobSelection />
        </template>
    </UDropdown>
</template>
