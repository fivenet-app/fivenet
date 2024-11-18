<script lang="ts" setup>
import EmailCreateForm from '~/components/mailer/EmailCreateForm.vue';
import EmailSettingsModal from '~/components/mailer/EmailSettingsModal.vue';
import MailerList from '~/components/mailer/MailerList.vue';
import MailerThread from '~/components/mailer/MailerThread.vue';
import ThreadCreateOrUpdateModal from '~/components/mailer/ThreadCreateOrUpdateModal.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { mailerDB, useMailerStore } from '~/store/mailer';

useHead({
    title: 'common.mail',
});
definePageMeta({
    title: 'common.mail',
    requiresAuth: true,
    permission: 'MailerService.ListEmails',
});

const { t } = useI18n();

const { can, isSuperuser } = useAuth();

const modal = useModal();

const mailerStore = useMailerStore();
const { emails, selectedEmail, selectedThread } = storeToRefs(mailerStore);

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

watch(selectedEmail, async () => {
    Promise.all([loadThreads()]);
});

async function loadThreads(): Promise<void> {
    if (!selectedEmail.value?.id) {
        return;
    }

    const count = await mailerDB.threads.count();
    await mailerStore.listThreads({
        pagination: {
            offset: 0,
        },
        emailIds: [selectedEmail.value.id],
        after: count > 0 ? undefined : toTimestamp(),
    });
}

const threads = useDexieLiveQuery(
    () =>
        mailerDB.threads
            .orderBy('id')
            .reverse()
            .toArray()
            .then((threads) => ({ threads: threads, loaded: true })),
    {
        initialValue: { threads: [], loaded: false },
    },
);

// Filter mails based on the selected tab
const filteredThreads = computed(() => {
    if (selectedTab.value === 1) {
        return threads.value.threads.filter((thread) => !thread.state?.archived && !!thread.state?.unread);
    } else if (selectedTab.value === 2) {
        return threads.value.threads.filter((thread) => !!thread.state?.archived);
    }

    return threads.value.threads.filter((thread) => !thread.state?.archived);
});

const threadState = computed(() => selectedThread.value?.state);

const isMailerPanelOpen = computed({
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

onBeforeMount(async () => {
    await mailerStore.listEmails();

    if (!route.query.thread) {
        return;
    }

    selectedThread.value = await mailerStore.getThread(route.query.thread as string);
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel id="mailerthreadlist" :width="450" :resizable="{ min: 325, max: 550 }">
            <UDashboardNavbar :title="$t('common.mail')" :badge="filteredThreads.length">
                <template #right>
                    <UButton
                        v-if="can('MailerService.CreateThread').value && selectedEmail"
                        color="gray"
                        trailing-icon="i-mdi-plus"
                        @click="modal.open(ThreadCreateOrUpdateModal, {})"
                    >
                        {{ $t('components.mailer.create_thread') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar
                v-if="selectedEmail"
                :ui="{ wrapper: 'p-0 gap-x-0', container: 'gap-x-0 justify-stretch items-stretch h-full flex flex-1 flex-col' }"
            >
                <div class="inline-flex gap-1 bg-gray-100 p-1 dark:bg-gray-800">
                    <ClientOnly>
                        <USelectMenu
                            v-model="selectedEmail"
                            :options="emails"
                            :placeholder="$t('common.mail')"
                            searchable
                            :searchable-placeholder="$t('common.search_field')"
                            :search-attributes="['label', 'email']"
                            trailing
                            by="id"
                            class="flex-1"
                        >
                            <template #label>
                                <span class="truncate">
                                    {{
                                        (selectedEmail?.label && selectedEmail?.label !== ''
                                            ? selectedEmail?.label + ' (' + selectedEmail.email + ')'
                                            : undefined) ??
                                        (selectedEmail?.userId
                                            ? $t('common.personal_email') + ' (' + selectedEmail.email + ')'
                                            : undefined) ??
                                        selectedEmail?.email ??
                                        $t('common.none')
                                    }}
                                </span>
                            </template>

                            <template #option="{ option }">
                                <span class="truncate">
                                    {{
                                        (option?.label && option?.label !== ''
                                            ? option?.label + ' (' + option.email + ')'
                                            : undefined) ??
                                        (option?.userId ? $t('common.personal_email') : undefined) ??
                                        option?.email ??
                                        $t('common.none')
                                    }}
                                </span>
                            </template>

                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.mail', 2)]) }} </template>
                        </USelectMenu>
                    </ClientOnly>
                </div>

                <UTabs
                    v-model="selectedTab"
                    :items="tabItems"
                    :ui="{ wrapper: 'w-full h-full space-y-0', list: { rounded: '' } }"
                />
            </UDashboardToolbar>

            <template v-if="selectedEmail">
                <div class="relative flex-1 overflow-x-auto">
                    <MailerList v-model="selectedThread" :threads="filteredThreads" :loaded="threads.loaded" />
                </div>

                <UDashboardToolbar class="flex justify-between border-t border-gray-200 px-3 py-3.5 dark:border-gray-700">
                    <template #left>
                        <UButton
                            color="gray"
                            trailing-icon="i-mdi-cog"
                            :label="$t('common.settings')"
                            @click="() => modal.open(EmailSettingsModal, {})"
                        />
                    </template>
                </UDashboardToolbar>
            </template>
            <div v-else class="flex flex-1 flex-col items-center">
                <div class="flex flex-1 flex-col items-center justify-center gap-2 text-gray-400 dark:text-gray-500">
                    <UIcon name="i-mdi-email-multiple" class="h-32 w-32" />
                    <EmailCreateForm personal-email />
                </div>
            </div>
        </UDashboardPanel>

        <UDashboardPanel v-if="selectedEmail" id="mailerthreadview" v-model="isMailerPanelOpen" collapsible grow side="right">
            <template v-if="selectedThread">
                <UDashboardNavbar>
                    <template #toggle>
                        <UDashboardNavbarToggle icon="i-mdi-close" />

                        <UDivider orientation="vertical" class="mx-1.5 lg:hidden" />
                    </template>

                    <template #left>
                        <UTooltip :text="$t('components.mailer.mark_unread')">
                            <UButton
                                :icon="!threadState?.unread ? 'i-mdi-check-circle-outline' : 'i-mdi-check-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.state = await mailerStore.setThreadState({
                                            threadId: selectedThread!.id,
                                            unread: !threadState?.unread,
                                        }))
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('components.mailer.mark_important')">
                            <UButton
                                :icon="!threadState?.important ? 'i-mdi-alert-circle-outline' : 'i-mdi-alert-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.state = await mailerStore.setThreadState({
                                            threadId: selectedThread!.id,
                                            important: !threadState?.important,
                                        }))
                                "
                            />
                        </UTooltip>
                    </template>

                    <template #right>
                        <UTooltip :text="$t('components.mailer.star_thread')">
                            <UButton
                                :icon="!threadState?.favorite ? 'i-mdi-star-circle-outline' : 'i-mdi-star-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.state = await mailerStore.setThreadState({
                                            threadId: selectedThread!.id,
                                            favorite: !threadState?.favorite,
                                        }))
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('components.mailer.mute_thread')">
                            <UButton
                                :icon="!threadState?.muted ? 'i-mdi-pause-circle-outline' : 'i-mdi-pause-circle'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    async () =>
                                        (selectedThread!.state = await mailerStore.setThreadState({
                                            threadId: selectedThread!.id,
                                            muted: !threadState?.muted,
                                        }))
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="threadState?.archived ? $t('common.unarchive') : $t('common.archive')">
                            <UButton
                                :icon="threadState?.archived ? 'i-mdi-archive' : 'i-mdi-archive-outline'"
                                color="gray"
                                variant="ghost"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () =>
                                            (selectedThread!.state = await mailerStore.setThreadState({
                                                threadId: selectedThread!.id,
                                                archived: !threadState?.archived,
                                            })),
                                    })
                                "
                            />
                        </UTooltip>

                        <UTooltip v-if="isSuperuser" :text="$t('common.delete')">
                            <UButton
                                icon="i-mdi-trash-can-outline"
                                color="gray"
                                variant="ghost"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () =>
                                            selectedEmail?.id &&
                                            selectedThread &&
                                            mailerStore.deleteThread({
                                                emailId: selectedEmail.id,
                                                threadId: selectedThread.id,
                                            }),
                                    })
                                "
                            />
                        </UTooltip>
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
