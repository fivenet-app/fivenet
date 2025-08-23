<script lang="ts" setup>
import EmailSettingsModal from '~/components/mailer/EmailSettingsModal.vue';
import { canAccess } from '~/components/mailer/helpers';
import MailerThread from '~/components/mailer/MailerThread.vue';
import MessageSearch from '~/components/mailer/MessageSearch.vue';
import TemplateModal from '~/components/mailer/TemplateModal.vue';
import ThreadCreateOrUpdateModal from '~/components/mailer/ThreadCreateOrUpdateModal.vue';
import ThreadList from '~/components/mailer/ThreadList.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useMailerStore } from '~/stores/mailer';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';
import type { ListThreadsResponse } from '~~/gen/ts/services/mailer/mailer';

useHead({
    title: 'common.mail',
});

definePageMeta({
    title: 'common.mail',
    requiresAuth: true,
    permission: 'mailer.MailerService/ListEmails',
});

const { t } = useI18n();

const { isSuperuser } = useAuth();

const modal = useOverlay();

const mailerStore = useMailerStore();
const { draft, emails, selectedEmail, selectedThread, threads, unreadThreadIds } = storeToRefs(mailerStore);

const items = [
    {
        label: t('common.all'),
        slot: 'all' as const,
        value: 'all',
    },
    {
        label: t('common.unread'),
        slot: 'unread' as const,
        value: 'unread',
    },
    {
        label: t('common.archive'),
        slot: 'archive' as const,
        value: 'archive',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'all';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

const page = useRouteQuery('page', '1', { transform: Number });

const { status, refresh } = useLazyAsyncData(`mailer-thread:${page.value}`, () => loadThreads(), {
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
            offset: calculateOffset(page.value, threads.value?.pagination),
        },
        emailIds: [selectedEmail.value.id],
        unread: selectedTab.value === 'unread' ? true : undefined,
        archived: selectedTab.value === 'archive' ? true : false,
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
    <UDashboardPanel id="mailerthreadlist" resizable :width="40" :min-size="35" :max-size="60">
        <UDashboardNavbar :title="$t('common.mail')" :badge="threads?.pagination?.totalCount ?? 0">
            <template #center>
                <MessageSearch />
            </template>

            <template #right>
                <UTooltip
                    v-if="
                        selectedEmail &&
                        !selectedEmail.deactivated &&
                        canAccess(selectedEmail.access, selectedEmail.userId, AccessLevel.WRITE)
                    "
                    :text="$t('components.mailer.create_thread')"
                >
                    <UButton color="neutral" trailing-icon="i-mdi-plus" @click="modal.open(ThreadCreateOrUpdateModal, {})">
                        {{ $t('components.mailer.create_thread') }}
                    </UButton>
                </UTooltip>
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
                    class="pt-1"
                    type="text"
                    disabled
                    :model-value="
                        (selectedEmail?.label && selectedEmail?.label !== ''
                            ? selectedEmail?.label + ' (' + selectedEmail.email + ')'
                            : undefined) ??
                        selectedEmail?.email ??
                        $t('common.none')
                    "
                />
                <USelectMenu
                    v-else
                    v-model="selectedEmail"
                    class="pt-1"
                    :items="emails"
                    :placeholder="$t('common.mail')"
                    searchable
                    :searchable-placeholder="$t('common.search_field')"
                    :search-attributes="['label', 'email']"
                    trailing
                    by="id"
                >
                    <template #item-label>
                        <span class="truncate">
                            {{
                                (selectedEmail?.label && selectedEmail?.label !== ''
                                    ? selectedEmail?.label + ' (' + selectedEmail.email + ')'
                                    : undefined) ??
                                selectedEmail?.email ??
                                $t('common.none')
                            }}

                            <UBadge v-if="selectedEmail?.deactivated" color="error" size="xs" :label="$t('common.disabled')" />
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

                        <UBadge v-if="option?.deactivated" color="error" size="xs" :label="$t('common.disabled')" />
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
                :items="items"
                :ui="{ wrapper: 'w-full h-full space-y-0' }"
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
                            :status="status"
                            :refresh="refresh"
                            hide-text
                        />
                    </template>
                </ThreadList>
            </div>

            <UDashboardToolbar
                class="flex justify-between border-t border-b-0 border-gray-200 px-3 py-3.5 dark:border-gray-700"
            >
                <template #left>
                    <UTooltip :text="$t('common.settings')">
                        <UButton color="neutral" trailing-icon="i-mdi-cog" @click="() => modal.open(EmailSettingsModal, {})">
                            <span class="hidden truncate md:block"> {{ $t('common.settings') }} </span>
                        </UButton>
                    </UTooltip>

                    <UTooltip :text="$t('pages.mailer.manage.title')">
                        <UButton icon="i-mdi-email-plus" color="neutral" to="/mail/manage" />
                    </UTooltip>
                </template>

                <template #right>
                    <UTooltip :text="$t('common.template', 2)">
                        <UButton
                            color="neutral"
                            trailing-icon="i-mdi-file-outline"
                            @click="() => modal.open(TemplateModal, {})"
                        >
                            <span class="hidden truncate md:block">{{ $t('common.template', 2) }}</span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardToolbar>
        </template>
    </UDashboardPanel>

    <UDashboardPanel v-if="selectedEmail" id="mailerthreadview" v-model="isMailerPanelOpen" side="right">
        <template v-if="selectedThread">
            <UDashboardNavbar>
                <template #toggle>
                    <UDashboardSidebarToggle icon="i-mdi-close" />

                    <USeparator class="mx-1.5 lg:hidden" orientation="vertical" />
                </template>

                <template #left>
                    <UTooltip :text="$t('components.mailer.mark_unread')">
                        <UButton
                            :icon="!threadState?.unread ? 'i-mdi-check-circle-outline' : 'i-mdi-check-circle'"
                            color="neutral"
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
                            color="neutral"
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
                            color="neutral"
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
                            color="neutral"
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
                            color="neutral"
                            variant="ghost"
                            @click="
                                modal.open(ConfirmModal, {
                                    confirm: async () => {
                                        selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                archived: !threadState?.archived,
                                            },
                                            true,
                                        );
                                        await refresh();
                                    },
                                })
                            "
                        />
                    </UTooltip>

                    <UTooltip v-if="isSuperuser" :text="!selectedThread.deletedAt ? $t('common.delete') : $t('common.restore')">
                        <UButton
                            :color="!selectedThread.deletedAt ? 'error' : 'success'"
                            :icon="!selectedThread.deletedAt ? 'i-mdi-delete-outline' : 'i-mdi-restore'"
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

        <div v-else class="hidden flex-1 flex-col items-center justify-center gap-2 text-gray-400 lg:flex dark:text-gray-500">
            <UIcon class="h-32 w-32" name="i-mdi-email-multiple" />
            <p>{{ $t('common.none_selected', [$t('common.mail')]) }}</p>
        </div>
    </UDashboardPanel>
</template>
