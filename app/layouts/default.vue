<script lang="ts" setup>
import type { NavigationMenuItem } from '@nuxt/ui';
import ClipboardModal from '~/components/clipboard/modal/ClipboardModal.vue';
import Banners from '~/components/partials/Banners.vue';
import SuperuserJobToggle from '~/components/partials/SuperuserJobToggle.vue';
import MathCalculatorDrawer from '~/components/quickbuttons/mathcalculator/MathCalculatorDrawer.vue';
import NotepadDrawer from '~/components/quickbuttons/notepad/NotepadDrawer.vue';
import PenaltyCalculatorDrawer from '~/components/quickbuttons/penaltycalculator/PenaltyCalculatorDrawer.vue';
import TopLogoDropdown from '~/components/TopLogoDropdown.vue';
import UserMenu from '~/components/UserMenu.vue';
import { useMailerStore } from '~/stores/mailer';

const { t } = useI18n();

const { can, activeChar, jobProps, isSuperuser, canBeSuperuser } = useAuth();

const { isDashboardSidebarSlideoverOpen, isHelpSlideoverOpen } = useDashboard();

const overlay = useOverlay();

const { website } = useAppConfig();

const mailerStore = useMailerStore();
const { unreadCount } = storeToRefs(mailerStore);

const { navigationItems } = useAppFeatures();
const items = computed(() =>
    navigationItems.value.map((item) =>
        item.to === '/mail'
            ? {
                  ...item,
                  badge: unreadCount.value > 0 ? (unreadCount.value <= 9 ? unreadCount.value.toString() : '9+') : undefined,
              }
            : item,
    ),
);

const footerLinks = computed(() =>
    [
        website.statsPage
            ? {
                  label: t('pages.stats.title'),
                  icon: 'i-mdi-analytics',
                  to: '/stats',
              }
            : undefined,
        {
            label: t('common.help'),
            icon: 'i-mdi-question-mark-circle-outline',
            tooltip: {
                kbds: ['?'],
            },
            onClick: () => (isHelpSlideoverOpen.value = true),
        },
        {
            label: t('common.about'),
            icon: 'i-mdi-about-circle-outline',
            to: '/about',
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const clipboardModal = overlay.create(ClipboardModal);

const clipboardLink = computed(() =>
    [
        activeChar.value &&
        can([
            'documents.DocumentsService/UpdateDocument',
            'citizens.CitizensService/GetUser',
            'vehicles.VehiclesService/ListVehicles',
        ]).value
            ? {
                  label: t('common.clipboard'),
                  icon: 'i-mdi-clipboard-list-outline',
                  tooltip: {
                      text: t('common.clipboard'),
                      kbds: ['Q', 'C'],
                  },
                  kbds: ['Q', 'C'],
                  onClick: () => clipboardModal.open(),
              }
            : undefined,
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

const penaltyCalculatorDrawer = overlay.create(PenaltyCalculatorDrawer);
const mathCalculatorDrawer = overlay.create(MathCalculatorDrawer);
const notepadDrawer = overlay.create(NotepadDrawer);

const quickAccessButtons = computed<NavigationMenuItem[]>(() =>
    [
        jobProps.value?.quickButtons?.penaltyCalculator || isSuperuser.value
            ? {
                  label: t('components.penaltycalculator.title'),
                  icon: 'i-mdi-gavel',
                  tooltip: {
                      text: t('components.penaltycalculator.title'),
                      kbds: ['Q', 'P'],
                  },
                  kbds: ['Q', 'P'],
                  onClick: () => {
                      isDashboardSidebarSlideoverOpen.value = false;
                      penaltyCalculatorDrawer.open();
                  },
              }
            : undefined,
        {
            label: t('components.mathcalculator.title'),
            icon: 'i-mdi-calculator',
            tooltip: {
                text: t('components.mathcalculator.title'),
                kbds: ['Q', 'M'],
            },
            kbds: ['Q', 'M'],
            onClick: () => {
                isDashboardSidebarSlideoverOpen.value = false;
                mathCalculatorDrawer.open();
            },
        },
        {
            label: t('components.notepad.title'),
            icon: 'i-mdi-notebook',
            tooltip: {
                text: t('components.notepad.title'),
                kbds: ['Q', 'N'],
            },
            kbds: ['Q', 'N'],
            onClick: () => {
                isDashboardSidebarSlideoverOpen.value = false;
                notepadDrawer.open();
            },
        },
    ].flatMap((item) => (item !== undefined ? [item] : [])),
);

defineShortcuts(extractShortcutsFromNavItems(items.value, '-'));
defineShortcuts(extractShortcutsFromNavItems(clipboardLink.value, '-'));
defineShortcuts(extractShortcuts(quickAccessButtons.value, '-'));
</script>

<template>
    <UDashboardGroup unit="rem" storage="local">
        <UDashboardSidebar
            id="default"
            v-model:open="isDashboardSidebarSlideoverOpen"
            class="bg-elevated/25"
            :default-size="16.5"
            :min-size="13.5"
            :max-size="23"
            collapsible
            resizable
            :ui="{ footer: 'lg:border-t lg:border-default' }"
        >
            <template #header="{ collapsed }">
                <div class="flex w-full min-w-0 items-center gap-1">
                    <TopLogoDropdown class="min-w-0 flex-1" :collapsed="collapsed" />
                    <NotificationsNotificationPopover v-if="!collapsed" class="shrink-0" />
                </div>
            </template>

            <template #default="{ collapsed }">
                <UDashboardSearchButton :collapsed="collapsed" :label="$t('common.search_field')" />

                <div v-if="collapsed" class="flex w-full justify-center">
                    <NotificationsNotificationPopover />
                </div>

                <UNavigationMenu orientation="vertical" tooltip popover :items="items" :collapsed="collapsed" />

                <template v-if="clipboardLink.length > 0">
                    <USeparator />

                    <UNavigationMenu orientation="vertical" tooltip popover :items="clipboardLink" :collapsed="collapsed" />
                </template>

                <template v-if="quickAccessButtons">
                    <USeparator />

                    <UNavigationMenu
                        orientation="vertical"
                        tooltip
                        popover
                        :items="quickAccessButtons"
                        :collapsed="collapsed"
                    />
                </template>

                <div class="flex-1" />

                <template v-if="canBeSuperuser || isSuperuser">
                    <SuperuserJobToggle :collapsed="collapsed" />

                    <USeparator />
                </template>

                <UNavigationMenu orientation="vertical" tooltip popover :items="footerLinks" :collapsed="collapsed" />
            </template>

            <template #footer="{ collapsed }">
                <UserMenu :collapsed="collapsed" />
            </template>
        </UDashboardSidebar>

        <slot />

        <Banners />

        <ClientOnly>
            <LazyPartialsCommandSearch v-if="activeChar" :children="items" />

            <LazyPartialsWebSocketStatusOverlay />

            <LazyPartialsEventsLayer />
        </ClientOnly>
    </UDashboardGroup>
</template>
