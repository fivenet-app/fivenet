<script lang="ts" setup>
import type { RoutePathSchema } from '@typed-router';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import { useSettingsStore } from '~/stores/settings';
import { reminderTimes } from '~/types/calendar';
import type { Perms } from '~~/gen/ts/perms';

const { t } = useI18n();

const { can, activeChar } = useAuth();

const settings = useSettingsStore();
const { startpage, design, streamerMode, audio, calendar } = storeToRefs(settings);

const startpages: { name: string; path: RoutePathSchema; permission?: Perms }[] = [
    { name: t('common.overview'), path: '/overview' },
    { name: t('common.mail'), path: '/mail', permission: 'mailer.MailerService/ListEmails' },
    { name: t('pages.citizens.title'), path: '/citizens', permission: 'citizens.CitizensService/ListCitizens' },
    { name: t('pages.vehicles.title'), path: '/vehicles', permission: 'vehicles.VehiclesService/ListVehicles' },
    { name: t('pages.documents.title'), path: '/documents', permission: 'documents.DocumentsService/ListDocuments' },
    { name: t('pages.jobs.overview.title'), path: '/jobs/overview', permission: 'jobs.JobsService/ListColleagues' },
    { name: t('common.calendar'), path: '/calendar' },
    {
        name: t('common.qualification', 2),
        path: '/qualifications',
        permission: 'qualifications.QualificationsService/ListQualifications',
    },
    { name: t('common.livemap'), path: '/livemap', permission: 'livemap.LivemapService/Stream' },
    { name: t('common.dispatch_center'), path: '/centrum', permission: 'centrum.CentrumService/TakeControl' },
    { name: t('common.wiki'), path: '/wiki', permission: 'wiki.WikiService/ListPages' },
];

const selectedHomepage = ref<(typeof startpages)[0]>();
watch(selectedHomepage, () => (startpage.value = selectedHomepage.value?.path ?? '/overview'));

const designDocumentsListStyle = ref(design.value.documents.listStyle === 'double');

watch(designDocumentsListStyle, async () => {
    if (designDocumentsListStyle.value) {
        design.value.documents.listStyle = 'double';
    } else {
        design.value.documents.listStyle = 'single';
    }
});

const calendarReminderTimes = [
    { label: t('components.auth.UserSettingsPanel.calendar_notifications.reminder_times.start'), value: 0 },
    ...reminderTimes.map((n) => ({ label: `${n / 60} ${t('common.time_ago.minute', n / 60)}`, value: n })),
];

const items = [
    {
        slot: 'settings',
        label: t('common.settings'),
        icon: 'i-mdi-cog',
    },
    {
        slot: 'notifications',
        label: t('common.notification', 2),
        icon: 'i-mdi-notification-settings',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});

const notificationSound = useSounds('/sounds/notification.mp3');

onBeforeMount(async () => (selectedHomepage.value = startpages.find((h) => h.path === startpage.value)));
</script>

<template>
    <UTabs v-model="selectedTab" class="w-full" :items="items" :ui="{ list: { rounded: '' } }">
        <template #settings>
            <UDashboardPanelContent>
                <UDashboardSection
                    :title="$t('common.theme')"
                    :description="$t('components.auth.UserSettingsPanel.customization')"
                >
                    <template #links>
                        <UColorModeSelect color="gray" />
                    </template>

                    <UFormGroup class="grid grid-cols-2 items-center gap-2" name="primaryColor" :label="$t('common.color')">
                        <ColorPickerTW v-model="design.ui.primary" name="primaryColor" />
                    </UFormGroup>

                    <UFormGroup
                        class="grid grid-cols-2 items-center gap-2"
                        name="grayColor"
                        :label="$t('components.auth.UserSettingsPanel.background_color')"
                    >
                        <ColorPickerTW v-model="design.ui.gray" name="grayColor" />
                    </UFormGroup>

                    <UFormGroup
                        class="grid grid-cols-2 items-center gap-2"
                        name="streamerMode"
                        :label="$t('components.auth.UserSettingsPanel.streamer_mode.title')"
                        :description="$t('components.auth.UserSettingsPanel.streamer_mode.description')"
                    >
                        <UToggle v-model="streamerMode" name="streamerMode" />
                    </UFormGroup>

                    <UFormGroup
                        class="grid grid-cols-2 items-center gap-2"
                        name="designDocumentsListStyle"
                        :label="$t('components.auth.UserSettingsPanel.documents_lists_style.title')"
                    >
                        <div class="inline-flex items-center gap-2 text-sm">
                            <span>{{ $t('components.auth.UserSettingsPanel.documents_lists_style.single') }}</span>
                            <UToggle v-model="designDocumentsListStyle" />
                            <span>{{ $t('components.auth.UserSettingsPanel.documents_lists_style.double') }}</span>
                        </div>
                    </UFormGroup>

                    <UFormGroup
                        class="grid grid-cols-2 items-center gap-2"
                        name="selectedHomepage"
                        :label="$t('components.auth.UserSettingsPanel.set_startpage.title')"
                    >
                        <ClientOnly v-if="activeChar">
                            <USelectMenu
                                v-model="selectedHomepage"
                                :options="startpages.filter((h) => h.permission === undefined || can(h.permission).value)"
                                option-attribute="name"
                                :searchable-placeholder="$t('common.search_field')"
                            />
                        </ClientOnly>
                        <p v-else class="text-sm">
                            {{ $t('components.auth.UserSettingsPanel.set_startpage.no_char_selected') }}
                        </p>
                    </UFormGroup>
                </UDashboardSection>
            </UDashboardPanelContent>
        </template>

        <template #notifications>
            <UDashboardPanelContent>
                <UDashboardSection
                    :title="$t('components.auth.UserSettingsPanel.volumes.title')"
                    :description="$t('components.auth.UserSettingsPanel.volumes.subtitle')"
                >
                    <template #links>
                        <UButton icon="i-mdi-play" @click="notificationSound.play()" />
                    </template>

                    <UFormGroup
                        class="grid grid-cols-2 items-center gap-2"
                        name="notificationsVolume"
                        :label="$t('components.auth.UserSettingsPanel.volumes.notifications_volume')"
                    >
                        <URange v-model="audio.notificationsVolume" :step="0.01" :min="0" :max="1" />
                        {{ audio.notificationsVolume <= 0 ? 0 : (audio.notificationsVolume * 100).toFixed(0) }}%
                    </UFormGroup>
                </UDashboardSection>

                <UDashboardSection
                    :title="$t('components.auth.UserSettingsPanel.calendar_notifications.title')"
                    :description="$t('components.auth.UserSettingsPanel.calendar_notifications.subtitle')"
                >
                    <UFormGroup
                        class="grid grid-cols-2 items-center gap-2"
                        name="calendarNotifications"
                        :label="$t('components.auth.UserSettingsPanel.calendar_notifications.reminder_times.name')"
                    >
                        <ClientOnly>
                            <USelectMenu
                                v-model="calendar.reminderTimes"
                                multiple
                                :options="calendarReminderTimes"
                                value-attribute="value"
                            >
                                <template #label>
                                    {{
                                        calendar.reminderTimes.length > 0
                                            ? [...calendar.reminderTimes]
                                                  .sort()
                                                  .map(
                                                      (n) =>
                                                          calendarReminderTimes.find((rt) => rt.value === n)?.label ??
                                                          $t('common.na'),
                                                  )
                                                  .join(', ')
                                            : $t('common.none')
                                    }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormGroup>
                </UDashboardSection>
            </UDashboardPanelContent>
        </template>
    </UTabs>
</template>
