<script lang="ts" setup>
import { sub } from 'date-fns';
import MessengerList from '~/components/messenger/MessengerList.vue';
import MessengerThread from '~/components/messenger/MessengerThread.vue';
import type { Thread } from '~~/gen/ts/resources/messenger/thread';

useHead({
    title: 'common.messenger',
});
definePageMeta({
    title: 'common.messenger',
    requiresAuth: true,
    permission: 'MessengerService.ListThreads',
});

const { t } = useI18n();

const tabItems = [
    {
        label: t('common.all'),
    },
    {
        label: t('common.unread'),
    },
];
const selectedTab = ref(0);

const dropdownItems = [
    [
        {
            label: t('components.messenger.mark_unread'),
            icon: 'i-mdi-check-circle-outline',
        },
        {
            label: t('components.messenger.mark_important'),
            icon: 'i-mdi-alert-circle-outline',
        },
    ],
    [
        {
            label: t('components.messenger.star_thread'),
            icon: 'i-mdi-star-circle-outline',
        },
        {
            label: t('components.messenger.mute_thread'),
            icon: 'i-mdi-pause-circle-outline',
        },
    ],
];

const threads = ref<Thread[]>([
    {
        id: '1',
        creator: {
            firstname: 'Alex Smith',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Meeting Schedule',
        closed: false,
        createdAt: toTimestamp(new Date()),
    },
    {
        id: '2',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        creator: {
            firstname: 'Jordan Brown',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        title: 'Project Update',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { minutes: 7 })),
    },
    {
        id: '3',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        creator: {
            firstname: 'Taylor Green',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        title: 'Lunch Plans',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { hours: 3 })),
    },
    {
        id: '4',
        creator: {
            firstname: 'Morgan White',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'New Proposal',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { days: 1 })),
    },
    {
        id: '5',
        creator: {
            firstname: 'Casey Gray',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Travel Itinerary',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { days: 1 })),
    },
    {
        id: '6',
        creator: {
            firstname: 'Jamie Johnson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Budget Report',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { days: 2 })),
    },
    {
        id: '7',
        creator: {
            firstname: 'Riley Davis',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Training Session',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { days: 2 })),
    },
    {
        id: '8',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        creator: {
            firstname: 'Kelly Wilson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        title: 'Happy Birthday!',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { days: 2 })),
    },
    {
        id: '9',
        creator: {
            firstname: 'Drew Moore',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Website Feedback',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { days: 5 })),
    },
    {
        id: '10',
        creator: {
            firstname: 'Jordan Taylor',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Gym Membership',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { days: 5 })),
    },
    {
        id: '11',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        creator: {
            firstname: 'Morgan Anderson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        title: 'Insurance Policy',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { days: 12 })),
    },
    {
        id: '12',
        creator: {
            firstname: 'Casey Thomas',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Book Club Meeting',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { months: 1 })),
    },
    {
        id: '13',
        creator: {
            firstname: 'Jamie Jackson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Recipe Exchange',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { months: 1 })),
    },
    {
        id: '14',
        creator: {
            firstname: 'Riley White',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Yoga Class Schedule',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { months: 1 })),
    },
    {
        id: '15',
        creator: {
            firstname: 'Kelly Harris',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Book Launch Event',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { months: 1 })),
    },
    {
        id: '16',
        creator: {
            firstname: 'Drew Martin',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Tech Conference',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { months: 1, days: 4 })),
    },
    {
        id: '17',
        creator: {
            firstname: 'Alex Thompson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Art Exhibition',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { months: 1, days: 15 })),
    },
    {
        id: '18',
        creator: {
            firstname: 'Jordan Garcia',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Networking Event',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { months: 1, days: 18 })),
    },
    {
        id: '19',
        creator: {
            firstname: 'Taylor Rodriguez',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Volunteer Opportunity',
        closed: false,
        createdAt: toTimestamp(sub(new Date(), { months: 1, days: 25 })),
    },
    {
        id: '20',
        createdAt: toTimestamp(sub(new Date(), { months: 2 })),
        userState: {
            threadId: '20',
            userId: 26061,
            muted: false,
            important: false,
            favorite: false,
        },
        title: 'Car Service Reminder',
        closed: false,
        creator: {
            firstname: 'Morgan Lopez',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        creatorJob: '',
    },
]);

// Filter mails based on the selected tab
const filteredThreads = computed(() => {
    if (selectedTab.value === 1) {
        return threads.value.filter((mail) => !!mail.userState?.lastRead);
    }

    return threads.value;
});

const selectedThread = ref<Thread | undefined>();

const isMessengerPanelOpen = computed({
    get() {
        return !!selectedThread.value;
    },
    set(value: boolean) {
        if (!value) {
            selectedThread.value = undefined;
        }
    },
});

// Reset selected mail if it's not in the filtered mails
watch(filteredThreads, () => {
    if (!filteredThreads.value.find((mail) => mail.id === selectedThread.value?.id)) {
        selectedThread.value = undefined;
    }
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel id="inbox" :width="425" :resizable="{ min: 300, max: 500 }">
            <UDashboardNavbar :title="$t('common.messenger')" :badge="filteredThreads.length">
                <template #right>
                    <UTabs
                        v-model="selectedTab"
                        :items="tabItems"
                        :ui="{ wrapper: '', list: { height: 'h-9', tab: { height: 'h-7', size: 'text-[13px]' } } }"
                    />
                </template>
            </UDashboardNavbar>

            <div class="relative flex-1 overflow-x-auto">
                <MessengerList v-model="selectedThread" :threads="filteredThreads" />
            </div>

            <UDashboardToolbar class="flex justify-between border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
                <template #right>
                    <UButtonGroup class="inline-flex">
                        <UButton color="gray" trailing-icon="i-mdi-plus">{{
                            $t('components.messenger.create_thread')
                        }}</UButton>
                    </UButtonGroup>
                </template>
            </UDashboardToolbar>
        </UDashboardPanel>

        <UDashboardPanel v-model="isMessengerPanelOpen" id="inboxmaillist" collapsible grow side="right">
            <template v-if="selectedThread">
                <UDashboardNavbar>
                    <template #toggle>
                        <UDashboardNavbarToggle icon="i-mdi-x" />

                        <UDivider orientation="vertical" class="mx-1.5 lg:hidden" />
                    </template>

                    <template #left>
                        <UTooltip :text="$t('common.archive')">
                            <UButton icon="i-mdi-archive" color="gray" variant="ghost" />
                        </UTooltip>
                    </template>

                    <template #right>
                        <UTooltip :text="$t('components.messenger.reply')">
                            <UButton icon="i-mdi-reply" color="gray" variant="ghost" />
                        </UTooltip>

                        <UTooltip :text="$t('components.messenger.forward')">
                            <UButton icon="i-mdi-forward" color="gray" variant="ghost" />
                        </UTooltip>

                        <UDivider orientation="vertical" class="mx-1.5" />

                        <UDropdown :items="dropdownItems">
                            <UButton icon="i-mdi-ellipsis-vertical" color="gray" variant="ghost" />
                        </UDropdown>
                    </template>
                </UDashboardNavbar>

                <MessengerThread :thread="selectedThread" />
            </template>
            <div v-else class="hidden flex-1 items-center justify-center lg:flex">
                <UIcon name="i-mdi-inbox" class="h-32 w-32 text-gray-400 dark:text-gray-500" />
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>
