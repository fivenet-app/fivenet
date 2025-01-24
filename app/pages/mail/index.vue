<script lang="ts" setup>
import EmailSettingsModal from '~/components/mailer/EmailSettingsModal.vue';
import { canAccess } from '~/components/mailer/helpers';
import MailerThread from '~/components/mailer/MailerThread.vue';
import MessagesSearch from '~/components/mailer/MessagesSearch.vue';
import TemplatesModal from '~/components/mailer/TemplatesModal.vue';
import ThreadCreateOrUpdateModal from '~/components/mailer/ThreadCreateOrUpdateModal.vue';
import ThreadList from '~/components/mailer/ThreadList.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DashboardPanel from '~/components/partials/DashboardPanel.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useMailerStore } from '~/store/mailer';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';
import type { ListThreadsResponse } from '~~/gen/ts/services/mailer/mailer';

useHead({
    title: 'common.mail',
});
definePageMeta({
    title: 'common.mail',
    requiresAuth: true,
    permission: 'MailerService.ListEmails',
});

const { t } = useI18n();

const { isSuperuser } = useAuth();

const modal = useModal();

const mailerStore = useMailerStore();
const { draft, emails, selectedEmail, selectedThread, threads, unreadThreadIds } = storeToRefs(mailerStore);

const tabItems = [
    {
        label: t('common.all'),
        slot: 'all',
    },
    {
        label: t('common.unread'),
        slot: 'unread',
    },
    {
        label: t('common.archive'),
        slot: 'archive',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = tabItems.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { ...route.query, tab: tabItems[value]?.slot }, hash: '#' });
    },
});

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() =>
    threads.value?.pagination?.pageSize ? threads.value?.pagination?.pageSize * (page.value - 1) : 0,
);

const { pending: loading, refresh: refresh } = useLazyAsyncData(`mailer-thread:${page.value}`, () => loadThreads(), {
    immediate: false,
});

async function loadThreads(): Promise<ListThreadsResponse | undefined> {
    if (!selectedEmail.value?.id) {
        return;
    }

    if (selectedEmail.value.settings === undefined) {
        await mailerStore.getEmail(selectedEmail.value.id);
        unreadThreadIds.value = [];
    }

    const resp = await mailerStore.listThreads({
        pagination: {
            offset: offset.value,
        },
        emailIds: [selectedEmail.value.id],
        unread: selectedTab.value === 1 ? true : undefined,
        archived: selectedTab.value === 2 ? true : false,
    });

    return resp;
}

watch(selectedEmail, async () => await refresh());

// Reset selected thread if it's not in the filtered mails
watch(threads, () => {
    if (!threads.value?.threads.find((thread) => thread.id === selectedThread.value?.id)) {
        selectedThread.value = undefined;
    }
});

const selectedThreadId = useRouteQuery('thread', '', { transform: Number });
watch(selectedThreadId, async () => {
    if (selectedThreadId.value <= 0) {
        return;
    }

    const thread = await mailerStore.getThread(selectedThreadId.value);
    selectedThread.value = thread;
});

// Refresh threads when unread tab is selected
watch(selectedTab, async () => await refresh());

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

// Set thread as query param for persistence between reloads
function updateQuery(): void {
    if (!selectedThread.value || !selectedEmail.value) {
        router.replace({
            query: route.query.tab
                ? {
                      ...route.query,
                      tab: route.query.tab,
                  }
                : {
                      ...route.query,
                  },
        });
    } else {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({
            query: { ...route.query, thread: selectedThread.value.id },
            hash: '#',
        });
    }
}

watch(selectedThread, () => updateQuery());

onBeforeMount(async () => {
    await mailerStore.listEmails();

    if (route.query.thread) {
        selectedThread.value = await mailerStore.getThread(parseInt(route.query.thread as string));
        updateQuery();
    }

    if (route.query.to) {
        const to = (route.query.to as string).trim().toLowerCase();
        if (!draft.value.recipients.find((r) => r.label === to)) {
            draft.value.recipients.push({
                label: to,
            });
        }

        modal.open(ThreadCreateOrUpdateModal, {});
    }
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel id="mailerthreadlist" :width="450" :resizable="{ min: 325, max: 550 }">
            <UDashboardNavbar :title="$t('common.mail')" :badge="threads?.pagination?.totalCount ?? 0">
                <template #center>
                    <MessagesSearch />
                </template>

                <template #right>
                    <UButton
                        v-if="
                            selectedEmail &&
                            !selectedEmail.deactivated &&
                            canAccess(selectedEmail.access, selectedEmail.userId, AccessLevel.WRITE)
                        "
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
                :ui="{
                    wrapper: 'p-0 gap-x-0',
                    container:
                        'gap-x-0 gap-y-1 justify-stretch items-stretch h-full inline-flex flex-col bg-gray-100 p-0 px-1 dark:bg-gray-800 min-w-0',
                }"
            >
                <ClientOnly>
                    <UInput
                        v-if="emails.length === 1"
                        type="text"
                        disabled
                        :model-value="
                            (selectedEmail?.label && selectedEmail?.label !== ''
                                ? selectedEmail?.label + ' (' + selectedEmail.email + ')'
                                : undefined) ??
                            selectedEmail?.email ??
                            $t('common.none')
                        "
                        class="pt-1"
                    />
                    <USelectMenu
                        v-else
                        v-model="selectedEmail"
                        :options="emails"
                        :placeholder="$t('common.mail')"
                        searchable
                        :searchable-placeholder="$t('common.search_field')"
                        :search-attributes="['label', 'email']"
                        trailing
                        class="pt-1"
                        by="id"
                    >
                        <template #label>
                            <span class="truncate">
                                {{
                                    (selectedEmail?.label && selectedEmail?.label !== ''
                                        ? selectedEmail?.label + ' (' + selectedEmail.email + ')'
                                        : undefined) ??
                                    selectedEmail?.email ??
                                    $t('common.none')
                                }}

                                <UBadge
                                    v-if="selectedEmail?.deactivated"
                                    color="red"
                                    size="xs"
                                    :label="$t('common.disabled')"
                                />
                            </span>
                        </template>

                        <template #option="{ option }">
                            <span class="truncate">
                                {{
                                    (option?.label && option?.label !== ''
                                        ? option?.label + ' (' + option.email + ')'
                                        : undefined) ??
                                    (option?.userId
                                        ? $t('common.personal_email') + (isSuperuser ? ' (' + option.email + ')' : '')
                                        : undefined) ??
                                    option?.email ??
                                    $t('common.none')
                                }}
                            </span>

                            <UBadge v-if="option?.deactivated" color="red" size="xs" :label="$t('common.disabled')" />
                        </template>

                        <template #option-empty="{ query: search }">
                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                        </template>

                        <template #empty> {{ $t('common.not_found', [$t('common.mail', 2)]) }} </template>
                    </USelectMenu>
                </ClientOnly>

                <UTabs
                    v-if="!selectedEmail?.deactivated"
                    v-model="selectedTab"
                    :items="tabItems"
                    :ui="{ wrapper: 'w-full h-full space-y-0', list: { rounded: '' } }"
                />
            </UDashboardToolbar>

            <template v-if="selectedEmail">
                <div class="relative flex flex-1 overflow-x-auto">
                    <DataErrorBlock
                        v-if="selectedEmail.deactivated"
                        :title="$t('errors.MailerService.ErrEmailDisabled.title')"
                        :message="$t('errors.MailerService.ErrEmailDisabled.content')"
                    />

                    <ThreadList v-else v-model="selectedThread" :threads="threads?.threads ?? []" :loaded="true">
                        <template #after>
                            <div class="flex-1" />

                            <Pagination
                                v-model="page"
                                :pagination="threads?.pagination"
                                :loading="loading"
                                :refresh="refresh"
                                hide-text
                            />
                        </template>
                    </ThreadList>
                </div>

                <UDashboardToolbar
                    class="flex justify-between border-b-0 border-t border-gray-200 px-3 py-3.5 dark:border-gray-700"
                >
                    <template #left>
                        <UTooltip :text="$t('common.settings')">
                            <UButton color="gray" trailing-icon="i-mdi-cog" @click="() => modal.open(EmailSettingsModal, {})">
                                <span class="hidden truncate md:block"> {{ $t('common.settings') }} </span>
                            </UButton>
                        </UTooltip>

                        <UTooltip :text="$t('pages.mailer.manage.title')">
                            <UButton icon="i-mdi-email-plus" color="gray" to="/mail/manage" />
                        </UTooltip>
                    </template>

                    <template #right>
                        <UTooltip :text="$t('common.template', 2)">
                            <UButton
                                color="gray"
                                trailing-icon="i-mdi-file-outline"
                                @click="() => modal.open(TemplatesModal, {})"
                            >
                                <span class="hidden truncate md:block">{{ $t('common.template', 2) }}</span>
                            </UButton>
                        </UTooltip>
                    </template>
                </UDashboardToolbar>
            </template>
        </UDashboardPanel>

        <DashboardPanel v-if="selectedEmail" id="mailerthreadview" v-model="isMailerPanelOpen" collapsible grow side="right">
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
                                        (selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                unread: !threadState?.unread,
                                            },
                                            true,
                                        ))
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
                                        (selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                important: !threadState?.important,
                                            },
                                            true,
                                        ))
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
                                        (selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                favorite: !threadState?.favorite,
                                            },
                                            true,
                                        ))
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
                                        (selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                muted: !threadState?.muted,
                                            },
                                            true,
                                        ))
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
                                            (selectedThread!.state = await mailerStore.setThreadState(
                                                {
                                                    threadId: selectedThread!.id,
                                                    archived: !threadState?.archived,
                                                },
                                                true,
                                            )),
                                    })
                                "
                            />
                        </UTooltip>

                        <UTooltip
                            v-if="isSuperuser"
                            :text="!selectedThread.deletedAt ? $t('common.delete') : $t('common.restore')"
                        >
                            <UButton
                                :color="!selectedThread.deletedAt ? 'red' : 'green'"
                                :icon="!selectedThread.deletedAt ? 'i-mdi-trash-can-outline' : 'i-mdi-restore'"
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
        </DashboardPanel>
    </UDashboardPage>
</template>
