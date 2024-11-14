<script lang="ts" setup>
import MailerList from '~/components/mailer/MailerList.vue';
import MailerSettingsModal from '~/components/mailer/MailerSettingsModal.vue';
import MailerThread from '~/components/mailer/MailerThread.vue';
import ThreadCreateOrUpdateModal from '~/components/mailer/ThreadCreateOrUpdateModal.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { mailerDB, useMailerStore } from '~/store/mailer';
import type { Email } from '~~/gen/ts/resources/mailer/email';

useHead({
    title: 'common.mail',
});
definePageMeta({
    title: 'common.mail',
    requiresAuth: true,
    permission: 'MailerService.ListThreads',
});

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const mailerStore = useMailerStore();
const { selectedThread } = storeToRefs(mailerStore);

const tabItems = [
    {
        label: t('common.all'),
    },
    {
        label: t('common.unread'),
    },
    {
        label: t('common.archive'),
    },
];
const selectedTab = ref(0);

const dropdownItems = computed(() =>
    [
        [
            {
                label: selectedThread.value?.userState?.archived ? t('common.unarchive') : t('common.archive'),
                icon: 'i-mdi-archive',
                click: () =>
                    modal.open(ConfirmModal, {
                        confirm: async () =>
                            selectedThread.value &&
                            mailerStore.setThreadState({ threadId: selectedThread.value!.id, archived: true }),
                    }),
            },
            {
                label: t('common.leave'),
                icon: 'i-mdi-door',
                click: () =>
                    modal.open(ConfirmModal, {
                        confirm: async () => selectedThread.value && mailerStore.leaveThread(selectedThread.value!.id),
                    }),
            },
        ],

        [
            can('MailerService.DeleteEmail').value
                ? {
                      label: t('common.delete'),
                      icon: 'i-mdi-trash-can-outline',
                      click: async () =>
                          modal.open(ConfirmModal, {
                              confirm: async () =>
                                  selectedThread.value && mailerStore.deleteThread({ threadId: selectedThread.value!.id }),
                          }),
                  }
                : undefined,
        ].flatMap((item) => (item !== undefined ? [item] : [])),
    ].flatMap((items) => (items.length > 0 ? [items] : [])),
);

const personalEmail: Email = {
    id: '0',
    domain: 'fivenet.app',
    label: t('common.personal_email'),
    internal: false,
};
const selectedEmail = ref<Email | undefined>(personalEmail);

const { data: emails } = useLazyAsyncData('emails', async () => {
    const emails = await mailerStore.listEmails();
    emails.unshift(personalEmail);
    return emails;
});

watch(selectedEmail, async () => loadThreads);

async function loadThreads(): Promise<void> {
    const count = await mailerDB.threads.count();

    const call = getGRPCMailerClient().listThreads({
        pagination: {
            offset: 0,
        },
        emailIds: emails.value && emails.value[0]?.id !== undefined ? [emails.value[0].id] : [],
        after: count > 0 ? undefined : toTimestamp(),
    });
    const { response } = await call;

    mailerDB.threads.bulkPut(response.threads);
}

const threads = useDexieLiveQuery(() => mailerDB.threads.toArray().then((threads) => ({ threads, loaded: true })), {
    initialValue: { threads: [], loaded: false },
});

// Filter mails based on the selected tab
const filteredThreads = computed(() => {
    if (selectedTab.value === 1) {
        return threads.value.threads.filter((thread) => !thread.userState?.archived && !!thread.userState?.lastRead);
    } else if (selectedTab.value === 2) {
        return threads.value.threads.filter((thread) => !!thread.userState?.archived);
    }

    return threads.value.threads.filter((thread) => !thread.userState?.archived);
});

const threadUserState = computed(() => selectedThread.value?.userState);

const isMailerPanelOpen = computed({
    get() {
        return !!selectedThread.value || editing.value;
    },
    set(value: boolean) {
        if (!value) {
            selectedThread.value = undefined;
        }
    },
});

// Reset selected mail if it's not in the filtered mails
watch(filteredThreads, () => {
    if (!filteredThreads.value.find((thread) => thread.id === selectedThread.value?.id)) {
        selectedThread.value = undefined;
    }
});

// Set thread as query param for persistence between reloads
const route = useRoute();
const router = useRouter();

watch(selectedThread, () => {
    if (!selectedThread.value) {
        router.replace({ query: {} });
    } else {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { thread: selectedThread.value.id }, hash: '#' });
    }
});

onMounted(async () => {
    loadThreads();

    if (!route.query.thread) {
        return;
    }

    selectedThread.value = await mailerStore.getThread(route.query.thread as string);
});

const editing = ref(false);
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel id="mailerthreadlist" :width="450" :resizable="{ min: 325, max: 550 }">
            <UDashboardNavbar :title="$t('common.mail')" :badge="filteredThreads.length">
                <template #right>
                    <UButton
                        v-if="can('MailerService.CreateThread').value"
                        color="gray"
                        trailing-icon="i-mdi-plus"
                        @click="modal.open(ThreadCreateOrUpdateModal, {})"
                    >
                        {{ $t('components.mailer.create_thread') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar
                :ui="{ wrapper: 'p-0 gap-x-0', container: 'gap-x-0 justify-stretch items-stretch h-full flex flex-1 flex-col' }"
            >
                <div class="bg-gray-100 p-1 dark:bg-gray-800">
                    <USelectMenu
                        v-model="selectedEmail"
                        :options="emails"
                        :placeholder="$t('common.mail')"
                        :searchable-placeholder="$t('common.search_field')"
                        :search-attributes="['label']"
                        trailing
                        by="id"
                    >
                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.mail', 2)]) }} </template>
                    </USelectMenu>
                </div>

                <UTabs
                    v-model="selectedTab"
                    :items="tabItems"
                    :ui="{ wrapper: 'w-full h-full space-y-0', list: { rounded: '' } }"
                />
            </UDashboardToolbar>

            <div class="relative flex-1 overflow-x-auto">
                <MailerList v-model="selectedThread" :threads="filteredThreads" :loaded="threads.loaded" />
            </div>

            <UDashboardToolbar class="flex justify-between border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
                <template #left>
                    <UButton color="gray" trailing-icon="i-mdi-cog" @click="() => modal.open(MailerSettingsModal, {})" />
                </template>
            </UDashboardToolbar>
        </UDashboardPanel>

        <UDashboardPanel id="mailerthreadview" v-model="isMailerPanelOpen" collapsible grow side="right">
            <template v-if="selectedThread">
                <UDashboardNavbar>
                    <template #toggle>
                        <UDashboardNavbarToggle icon="i-mdi-close" />

                        <UDivider orientation="vertical" class="mx-1.5 lg:hidden" />
                    </template>

                    <template #left>
                        <UTooltip :text="$t('components.mailer.mark_unread')">
                            <UButton
                                :icon="!threadUserState?.unread ? 'i-mdi-check-circle-outline' : 'i-mdi-check-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.userState = await mailerStore.setThreadState({
                                            threadId: selectedThread!.id,
                                            unread: !threadUserState?.unread,
                                        }))
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('components.mailer.mark_important')">
                            <UButton
                                :icon="!threadUserState?.important ? 'i-mdi-alert-circle-outline' : 'i-mdi-alert-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.userState = await mailerStore.setThreadState({
                                            threadId: selectedThread!.id,
                                            important: !threadUserState?.important,
                                        }))
                                "
                            />
                        </UTooltip>
                    </template>

                    <template #right>
                        <UTooltip :text="$t('components.mailer.star_thread')">
                            <UButton
                                :icon="!threadUserState?.favorite ? 'i-mdi-star-circle-outline' : 'i-mdi-star-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.userState = await mailerStore.setThreadState({
                                            threadId: selectedThread!.id,
                                            favorite: !threadUserState?.favorite,
                                        }))
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('components.mailer.mute_thread')">
                            <UButton
                                :icon="!threadUserState?.muted ? 'i-mdi-pause-circle-outline' : 'i-mdi-pause-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.userState = await mailerStore.setThreadState({
                                            threadId: selectedThread!.id,
                                            muted: !threadUserState?.muted,
                                        }))
                                "
                            />
                        </UTooltip>

                        <UDivider orientation="vertical" class="mx-1.5" />

                        <UDropdown :items="dropdownItems">
                            <UButton icon="i-mdi-ellipsis-vertical" color="gray" variant="ghost" />
                        </UDropdown>
                    </template>
                </UDashboardNavbar>

                <MailerThread :thread-id="selectedThread.id" />
            </template>

            <div
                v-else
                class="hidden flex-1 flex-col items-center justify-center gap-2 text-gray-400 lg:flex dark:text-gray-500"
            >
                <UIcon name="i-mdi-email-multiple" class="h-32 w-32" />
                <p>{{ $t('common.none_selected', [$t('common.mail')]) }}</p>
            </div>
        </UDashboardPanel>
    </UDashboardPage>
</template>
