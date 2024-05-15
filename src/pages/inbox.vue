<script lang="ts" setup>
import { sub } from 'date-fns';
import InboxList from '~/components/mailer/InboxList.vue';
import InboxMail from '~/components/mailer/InboxMail.vue';
import type { Mail } from '~~/gen/ts/resources/mailer/mail';

useHead({
    title: 'common.inbox',
});
definePageMeta({
    title: 'common.inbox',
    requiresAuth: true,
    permission: 'TODOService.TODOMethod',
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
            label: t('components.inbox.mark_unread'),
            icon: 'i-mdi-check-circle-outline',
        },
        {
            label: t('components.inbox.mark_important'),
            icon: 'i-mdi-alert-circle-outline',
        },
    ],
    [
        {
            label: t('components.inbox.star_thread'),
            icon: 'i-mdi-star-circle-outline',
        },
        {
            label: t('components.inbox.mute_thread'),
            icon: 'i-mdi-pause-circle-outline',
        },
    ],
];

const mails = ref<Mail[]>([
    {
        id: '1',
        from: {
            firstname: 'Alex Smith',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Meeting Schedule',
        body: "Hi there, just a quick reminder about our meeting scheduled for 10 AM tomorrow. We'll be discussing the new marketing strategies and I would really appreciate your input on the matter. Looking forward to a productive session.",
        createdAt: toTimestamp(new Date()),
    },
    {
        id: '2',
        unread: true,
        from: {
            firstname: 'Jordan Brown',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        subject: 'Project Update',
        body: "I wanted to provide you with the latest update on the project. We've made significant progress on the development front and I've attached a detailed report for your review. Please let me know your thoughts and any areas for improvement.",
        createdAt: toTimestamp(sub(new Date(), { minutes: 7 })),
    },
    {
        id: '3',
        unread: true,
        from: {
            firstname: 'Taylor Green',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        subject: 'Lunch Plans',
        body: 'Hey! I was wondering if you would like to grab lunch this Friday. I know a great spot downtown that serves the best Mexican cuisine. It would be a great opportunity for us to catch up and discuss the upcoming team event.',
        createdAt: toTimestamp(sub(new Date(), { hours: 3 })),
    },
    {
        id: '4',
        from: {
            firstname: 'Morgan White',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'New Proposal',
        body: "I've attached the new proposal for our next project. It outlines all the objectives, timelines, and resource allocations. I'm particularly excited about the innovative approach we're taking this time. Please have a look and let me know your thoughts.",
        createdAt: toTimestamp(sub(new Date(), { days: 1 })),
    },
    {
        id: '5',
        from: {
            firstname: 'Casey Gray',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Travel Itinerary',
        body: "Your travel itinerary for the upcoming business trip is ready. I've included all flight details, hotel reservations, and meeting schedules. Please review and let me know if there are any changes you would like to make or any additional arrangements needed.",
        createdAt: toTimestamp(sub(new Date(), { days: 1 })),
    },
    {
        id: '6',
        from: {
            firstname: 'Jamie Johnson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Budget Report',
        body: "I've completed the budget report for this quarter. It includes a detailed analysis of our expenditures and revenue, along with projections for the next quarter. I believe there are some areas where we can optimize our spending. Let's discuss this in our next finance meeting.",
        createdAt: toTimestamp(sub(new Date(), { days: 2 })),
    },
    {
        id: '7',
        from: {
            firstname: 'Riley Davis',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Training Session',
        body: "Just a reminder about the training session scheduled for next week. We'll be covering new software tools that are crucial for our workflow. It's important that everyone attends as this will greatly enhance our team's efficiency. Please confirm your availability.",
        createdAt: toTimestamp(sub(new Date(), { days: 2 })),
    },
    {
        id: '8',
        unread: true,
        from: {
            firstname: 'Kelly Wilson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        subject: 'Happy Birthday!',
        body: 'Happy Birthday! Wishing you a fantastic day filled with joy and laughter. Your dedication and hard work throughout the year have been invaluable to our team. Enjoy your day to the fullest!',
        createdAt: toTimestamp(sub(new Date(), { days: 2 })),
    },
    {
        id: '9',
        from: {
            firstname: 'Drew Moore',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Website Feedback',
        body: 'We are in the process of revamping our company website and I would greatly appreciate your feedback on the new design. Your perspective is always insightful and could help us enhance the user experience significantly. Please let me know a convenient time for you to discuss this.',
        createdAt: toTimestamp(sub(new Date(), { days: 5 })),
    },
    {
        id: '10',
        from: {
            firstname: 'Jordan Taylor',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Gym Membership',
        body: "This is a friendly reminder that your gym membership is due for renewal at the end of this month. We've added several new classes and facilities that I think you'll really enjoy. Let me know if you would like a tour of the new facilities.",
        createdAt: toTimestamp(sub(new Date(), { days: 5 })),
    },
    {
        id: '11',
        unread: true,
        from: {
            firstname: 'Morgan Anderson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        subject: 'Insurance Policy',
        body: "I'm writing to inform you that your insurance policy details have been updated. The new document outlines the changes in coverage and premium rates. It's important to review these changes to ensure they meet your needs. Please don't hesitate to contact me if you have any questions.",
        createdAt: toTimestamp(sub(new Date(), { days: 12 })),
    },
    {
        id: '12',
        from: {
            firstname: 'Casey Thomas',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Book Club Meeting',
        body: "I'm excited to remind you about our next book club meeting scheduled for next Thursday. We'll be discussing 'The Great Gatsby,' and I'm looking forward to hearing everyone's perspectives. Also, we will be choosing our next book, so bring your suggestions!",
        createdAt: toTimestamp(sub(new Date(), { months: 1 })),
    },
    {
        id: '13',
        from: {
            firstname: 'Jamie Jackson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Recipe Exchange',
        body: "Don't forget to send in your favorite recipe for our upcoming recipe exchange. It's a great opportunity to share and discover new and delicious meals. I'm particularly excited to try out new dishes and add some variety to my cooking.",
        createdAt: toTimestamp(sub(new Date(), { months: 1 })),
    },
    {
        id: '14',
        from: {
            firstname: 'Riley White',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Yoga Class Schedule',
        body: "The new schedule for yoga classes is now available. We've added some new styles and adjusted the timings to accommodate more participants. I believe these classes are a great way to relieve stress and stay healthy. Hope to see you there!",
        createdAt: toTimestamp(sub(new Date(), { months: 1 })),
    },
    {
        id: '15',
        from: {
            firstname: 'Kelly Harris',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Book Launch Event',
        body: "I'm thrilled to invite you to my book launch event next month. It's been a journey writing this book    and I'm eager to share it with you. The event will include a reading session, Q&A, and a signing opportunity. Your support would mean a lot to me.",
        createdAt: toTimestamp(sub(new Date(), { months: 1 })),
    },
    {
        id: '16',
        from: {
            firstname: 'Drew Martin',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Tech Conference',
        body: "Join us at the upcoming tech conference where we will be discussing the latest trends and innovations in technology. This is a great opportunity to network with industry leaders and learn about cutting-edge developments. Your participation would greatly contribute to our team's knowledge and growth.",
        createdAt: toTimestamp(sub(new Date(), { months: 1, days: 4 })),
    },
    {
        id: '17',
        from: {
            firstname: 'Alex Thompson',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Art Exhibition',
        body: "I wanted to invite you to check out the new art exhibition this weekend. It features some amazing contemporary artists and their latest works. It's a great opportunity to immerse yourself in the local art scene and get inspired. Let me know if you're interested in going together.",
        createdAt: toTimestamp(sub(new Date(), { months: 1, days: 15 })),
    },
    {
        id: '18',
        from: {
            firstname: 'Jordan Garcia',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Networking Event',
        body: "I'm looking forward to seeing you at the networking event next week. It's a great chance to connect with professionals from various industries and expand our professional network. There will also be guest speakers discussing key business trends. Your presence would add great value to the discussions.",
        createdAt: toTimestamp(sub(new Date(), { months: 1, days: 18 })),
    },
    {
        id: '19',
        from: {
            firstname: 'Taylor Rodriguez',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Volunteer Opportunity',
        body: "We're looking for volunteers for the upcoming community event. It's a great opportunity to give back and make a positive impact. There are various roles available, so you can choose something that aligns with your interests and skills. Let me know if you're interested and I'll provide more details.",
        createdAt: toTimestamp(sub(new Date(), { months: 1, days: 25 })),
    },
    {
        id: '20',
        from: {
            firstname: 'Morgan Lopez',
            userId: 123,
            lastname: '',
            dateofbirth: '',
            identifier: '',
            job: '',
            jobGrade: 1,
        },
        unread: false,
        subject: 'Car Service Reminder',
        body: "Just a reminder that your car is due for service next week. Regular maintenance is important to ensure your vehicle's longevity and performance. I've included the details of the service center and the recommended services in this email. Feel free to contact them directly to schedule an appointment.",
        createdAt: toTimestamp(sub(new Date(), { months: 2 })),
    },
]);

// Filter mails based on the selected tab
const filteredMails = computed(() => {
    if (selectedTab.value === 1) {
        return mails.value.filter((mail) => !!mail.unread);
    }

    return mails.value;
});

const selectedMail = ref<Mail | undefined>();

const isMailPanelOpen = computed({
    get() {
        return !!selectedMail.value;
    },
    set(value: boolean) {
        if (!value) {
            selectedMail.value = undefined;
        }
    },
});

// Reset selected mail if it's not in the filtered mails
watch(filteredMails, () => {
    if (!filteredMails.value.find((mail) => mail.id === selectedMail.value?.id)) {
        selectedMail.value = undefined;
    }
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel id="inbox" :width="400" :resizable="{ min: 300, max: 500 }">
            <UDashboardNavbar :title="$t('common.inbox')" :badge="filteredMails.length">
                <template #right>
                    <UTabs
                        v-model="selectedTab"
                        :items="tabItems"
                        :ui="{ wrapper: '', list: { height: 'h-9', tab: { height: 'h-7', size: 'text-[13px]' } } }"
                    />
                </template>
            </UDashboardNavbar>

            <InboxList v-model="selectedMail" :mails="filteredMails" />
        </UDashboardPanel>

        <UDashboardPanel v-model="isMailPanelOpen" id="inboxmaillist" collapsible grow side="right">
            <template v-if="selectedMail">
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
                        <UTooltip :text="$t('components.inbox.reply')">
                            <UButton icon="i-mdi-reply" color="gray" variant="ghost" />
                        </UTooltip>

                        <UTooltip :text="$t('components.inbox.forward')">
                            <UButton icon="i-mdi-forward" color="gray" variant="ghost" />
                        </UTooltip>

                        <UDivider orientation="vertical" class="mx-1.5" />

                        <UDropdown :items="dropdownItems">
                            <UButton icon="i-mdi-ellipsis-vertical" color="gray" variant="ghost" />
                        </UDropdown>
                    </template>
                </UDashboardNavbar>

                <InboxMail :mail="selectedMail" />
            </template>
            <div v-else class="hidden flex-1 items-center justify-center lg:flex">
                <UIcon name="i-mdi-inbox" class="h-32 w-32 text-gray-400 dark:text-gray-500" />
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>
