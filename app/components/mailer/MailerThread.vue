<script setup lang="ts">
import { breakpointsTailwind } from '@vueuse/core';
import type { DropdownMenuItem, FormSubmitEvent } from '@nuxt/ui';
import type { JSONContent } from '@tiptap/core';
import { isSameDay } from 'date-fns';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useMailerStore } from '~/stores/mailer';
import { contentToTiptapValue, tiptapToContent } from '~/utils/content';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access/access';
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/messages/message';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import CustomContentRenderer from '../partials/content/CustomContentRenderer.vue';
import DocumentInfoPopover from '../partials/documents/DocumentInfoPopover.vue';
import EmailInfoPopover from './EmailInfoPopover.vue';
import { canAccess, generateResponseTitle } from './helpers';
import TemplateSelector from './TemplateSelector.vue';
import ThreadAttachmentsModal from './ThreadAttachmentsModal.vue';

const props = withDefaults(
    defineProps<{
        threadId: number;
        selected?: boolean;
    }>(),
    {
        selected: false,
    },
);

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'refresh'): void;
}>();

const minuteClock = useMinuteClock();
const breakpoints = useBreakpoints(breakpointsTailwind);
const isMobile = breakpoints.smaller('lg');

function isToday(date: Date): boolean {
    return isSameDay(date, minuteClock.value);
}

function startsNewDay(index: number): boolean {
    if (index === 0) return false;

    const current = messages.value?.messages[index];
    const previous = messages.value?.messages[index - 1];
    if (!current || !previous) return false;

    return !isSameDay(toDate(current.createdAt), toDate(previous.createdAt));
}

const overlay = useOverlay();

const { can, isSuperuser } = useAuth();
const { t } = useI18n();

const notifications = useNotificationsStore();

const { maxContentLength } = useAppConfig();

const mailerStore = useMailerStore();
const { draft: state, addressBook, messages, selectedEmail, selectedThread } = storeToRefs(mailerStore);

const schema = z.object({
    title: z.coerce.string().min(1).max(255),
    content: z.custom<JSONContent | string>().optional(),
    recipients: z
        .object({ label: z.coerce.string().min(6).max(80) })
        .array()
        .max(20)
        .default([]),
    attachments: z.custom<MessageAttachment>().array().max(3).default([]),
});

type Schema = z.output<typeof schema>;

const formSnapshot = computed(() => ({
    title: state.value.title,
    content: state.value.content,
    recipients: state.value.recipients.map((recipient) => recipient.label.trim().toLowerCase()),
    attachments: state.value.attachments.map((attachment) => ({
        type: attachment.data.oneofKind,
        documentId: attachment.data.oneofKind === 'document' ? attachment.data.document.id : null,
    })),
    selectedEmailId: selectedEmail.value?.id ?? null,
}));

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(formSnapshot);

async function resetForm(): Promise<void> {
    if (selectedThread.value) {
        if (state.value.title === '') {
            state.value.title = generateResponseTitle(selectedThread.value);
        }

        if ((!state.value.content || state.value.content === '<p><br></p>') && !!selectedEmail.value?.settings?.signature) {
            const end = editorRef.value?.editor?.$doc.content.size || 0;

            editorRef.value?.editor?.commands.insertContentAt(
                end,
                contentToTiptapValue(selectedEmail.value.settings.signature),
            );
        }
    }

    // Tiptap updates the v-model from its transaction asynchronously. Wait
    // until that update has been applied before recording the clean baseline;
    // otherwise opening a thread can be mistaken for an edit.
    await nextTick();
    syncSnapshot();
}

const { data: thread, status } = useAuthedLazyAsyncData(
    'userState',
    `mailer-thread:${props.threadId}`,
    ({ signal }) => mailerStore.getThread(props.threadId, { abort: signal }),
    {
        watch: [() => props.threadId],
    },
);

const messagePage = useRouteQuery('messagePage', '1', { transform: Number });

const { status: messagesStatus, refresh: refreshMessages } = useAuthedLazyAsyncData(
    'userState',
    () => `mailer-thread:${props.threadId}-messages:${messagePage.value}`,
    async ({ signal }) => {
        const response = await mailerStore.listThreadMessages(
            {
                pagination: {
                    offset: calculateOffset(messagePage.value, messages.value?.pagination),
                },
                emailId: selectedEmail.value!.id,
                threadId: props.threadId,
            },
            { abort: signal },
        );

        await resetForm();

        return response;
    },
);

async function refreshFirstMessagePage(): Promise<void> {
    if (messagePage.value !== 1) {
        messagePage.value = 1;
        return;
    }

    await refreshMessages();
}

watch(messagePage, async () => await refreshMessages());
watch(() => props.threadId, refreshFirstMessagePage);

watchDebounced(
    () => props.threadId,
    async () => {
        if (!thread.value?.state?.unread) return;

        if (!canAccess(selectedEmail.value?.access, selectedEmail.value?.userId, AccessLevel.WRITE)) return;

        await mailerStore.setThreadState({
            threadId: props.threadId,
            unread: false,
        });
    },
    {
        debounce: 500,
        maxWait: 2500,
    },
);

const threadState = computed(() => selectedThread.value?.state);
const visibleRecipients = computed(() => thread.value?.recipients.slice(0, isMobile.value ? 2 : 4) ?? []);
const hiddenRecipients = computed(() => thread.value?.recipients.slice(isMobile.value ? 2 : 4) ?? []);

async function toggleThreadState(field: 'unread' | 'important' | 'favorite' | 'muted'): Promise<void> {
    if (!selectedThread.value) return;

    selectedThread.value.state = await mailerStore.setThreadState(
        {
            threadId: selectedThread.value.id,
            [field]: !threadState.value?.[field],
        },
        true,
    );
}

function archiveThread(): void {
    if (!selectedThread.value) return;

    confirmModal.open({
        confirm: async () => {
            selectedThread.value!.state = await mailerStore.setThreadState(
                {
                    threadId: selectedThread.value!.id,
                    archived: !threadState.value?.archived,
                },
                true,
            );
            emit('refresh');
        },
    });
}

function deleteThread(): void {
    if (!selectedEmail.value?.id || !selectedThread.value) return;

    confirmModal.open({
        confirm: async () =>
            mailerStore.deleteThread({
                emailId: selectedEmail.value!.id,
                threadId: selectedThread.value!.id,
            }),
    });
}

const mobileActionItems = computed<DropdownMenuItem[][]>(() => [
    [
        {
            label: !threadState.value?.unread ? t('components.mailer.mark_unread') : t('components.mailer.mark_read'),
            icon: !threadState.value?.unread ? 'i-mdi-check-circle-outline' : 'i-mdi-check-circle',
            onSelect: () => toggleThreadState('unread'),
        },
        {
            label: t('components.mailer.mark_important'),
            icon: !threadState.value?.important ? 'i-mdi-alert-circle-outline' : 'i-mdi-alert-circle',
            onSelect: () => toggleThreadState('important'),
        },
        {
            label: t('components.mailer.star_thread'),
            icon: !threadState.value?.favorite ? 'i-mdi-star-circle-outline' : 'i-mdi-star-circle',
            onSelect: () => toggleThreadState('favorite'),
        },
        {
            label: t('components.mailer.mute_thread'),
            icon: !threadState.value?.muted ? 'i-mdi-pause-circle-outline' : 'i-mdi-pause-circle',
            onSelect: () => toggleThreadState('muted'),
        },
        {
            label: threadState.value?.archived ? t('common.unarchive') : t('common.archive'),
            icon: threadState.value?.archived ? 'i-mdi-archive' : 'i-mdi-archive-outline',
            onSelect: archiveThread,
        },
        ...(isSuperuser.value
            ? [
                  {
                      label: !selectedThread.value?.deletedAt ? t('common.delete') : t('common.restore'),
                      icon: !selectedThread.value?.deletedAt ? 'i-mdi-delete-outline' : 'i-mdi-restore',
                      color: !selectedThread.value?.deletedAt ? ('error' as const) : ('success' as const),
                      onSelect: deleteThread,
                  },
              ]
            : []),
    ],
]);

async function postMessage(values: Schema): Promise<void> {
    if (!selectedEmail.value?.id) return;

    await mailerStore.postMessage({
        message: {
            id: 0,
            senderId: selectedEmail.value.id,
            threadId: props.threadId,
            title: values.title,
            content: tiptapToContent(values.content),
            data: {
                attachments: values.attachments.filter((a) => {
                    if (a.data.oneofKind === 'document') {
                        return a.data.document.id > 0;
                    }

                    return false;
                }),
            },
        },
        recipients: [...new Set(values.recipients.map((r) => r.label.trim()))],
    });

    notifications.add({
        title: { key: 'notifications.action_successful.title', parameters: {} },
        description: { key: 'notifications.action_successful.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });

    // Clear draft data
    state.value.title = '';
    state.value.content = '';
    state.value.recipients = [];
    state.value.attachments = [];

    resetForm();
}

watch(
    state.value.recipients,
    () =>
        (state.value.recipients = state.value.recipients.filter(
            (item, idx) => state.value.recipients.findIndex((r) => r.label.toLowerCase() === item.label.toLowerCase()) === idx,
        )),
);

const messageRefs = ref<Element[]>([]);

function scrollToMessage(messageId: number): void {
    const ref = messageRefs.value[messageId];
    if (ref) {
        ref.scrollIntoView({ block: 'start' });
    }
}

const selectedMessageId = useRouteQuery('msg', 0, { transform: Number });
const selectedMessage = computed(() => selectedMessageId.value);
watch(selectedMessageId, () => scrollToMessage(selectedMessageId.value));

watch(messages, () => {
    if (selectedMessageId.value !== 0) {
        scrollToMessage(selectedMessageId.value);
    }
});

function onCreate(item: string): void {
    const email = item.trim();
    if (
        email.length < 6 ||
        email.length > 80 ||
        state.value.recipients.find((r) => r.label.toLowerCase() === email.toLowerCase())
    )
        return;

    state.value.recipients.push({ label: email });
}

const { submit, canSubmit } = useSubmitGuard(async (event: FormSubmitEvent<Schema>) => {
    if (!selectedEmail.value?.id) return;
    await postMessage(event.data);
});

const editorRef = useTemplateRef('editorRef');

const confirmModal = overlay.create(ConfirmModal);
const threadAttachmentsModal = overlay.create(ThreadAttachmentsModal);

async function closeThread(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UDashboardPanel
        id="mail-thread-view"
        :ui="{ root: 'min-h-full pb-(--page-content-bottom-offset)', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }"
    >
        <template #header>
            <UDashboardNavbar :title="thread?.title" :toggle="false">
                <template #title>
                    <h3 class="line-clamp-2 text-left font-semibold break-words text-highlighted hover:line-clamp-none">
                        {{ thread?.title }}
                    </h3>
                </template>

                <template #leading>
                    <UButton
                        class="-ms-1.5"
                        icon="i-mdi-close"
                        color="neutral"
                        variant="ghost"
                        :aria-label="$t('common.close', 1)"
                        @click="closeThread"
                    />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <template #left>
                    <template v-if="!isMobile">
                        <UTooltip
                            :text="
                                !threadState?.unread ? $t('components.mailer.mark_unread') : $t('components.mailer.mark_read')
                            "
                        >
                            <UButton
                                :icon="!threadState?.unread ? 'i-mdi-check-circle-outline' : 'i-mdi-check-circle'"
                                :color="!threadState?.unread ? 'neutral' : 'success'"
                                variant="ghost"
                                @click="
                                    async () => {
                                        selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                unread: !threadState?.unread,
                                            },
                                            true,
                                        );
                                    }
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('components.mailer.mark_important')">
                            <UButton
                                :icon="!threadState?.important ? 'i-mdi-alert-circle-outline' : 'i-mdi-alert-circle'"
                                :color="!threadState?.important ? 'neutral' : 'red'"
                                variant="ghost"
                                @click="
                                    async () => {
                                        selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                important: !threadState?.important,
                                            },
                                            true,
                                        );
                                    }
                                "
                            />
                        </UTooltip>
                    </template>
                </template>

                <template #right>
                    <template v-if="!isMobile">
                        <UTooltip :text="$t('components.mailer.star_thread')">
                            <UButton
                                :icon="!threadState?.favorite ? 'i-mdi-star-circle-outline' : 'i-mdi-star-circle'"
                                :color="!threadState?.favorite ? 'neutral' : 'amber'"
                                variant="ghost"
                                @click="
                                    async () => {
                                        selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                favorite: !threadState?.favorite,
                                            },
                                            true,
                                        );
                                    }
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="$t('components.mailer.mute_thread')">
                            <UButton
                                :icon="!threadState?.muted ? 'i-mdi-pause-circle-outline' : 'i-mdi-pause-circle'"
                                :color="!threadState?.muted ? 'neutral' : 'orange'"
                                variant="ghost"
                                @click="
                                    async () => {
                                        selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                muted: !threadState?.muted,
                                            },
                                            true,
                                        );
                                    }
                                "
                            />
                        </UTooltip>

                        <UTooltip :text="threadState?.archived ? $t('common.unarchive') : $t('common.archive')">
                            <UButton
                                :icon="threadState?.archived ? 'i-mdi-archive' : 'i-mdi-archive-outline'"
                                :color="threadState?.archived ? 'neutral' : 'gray'"
                                variant="ghost"
                                @click="archiveThread"
                            />
                        </UTooltip>

                        <UTooltip
                            v-if="isSuperuser && selectedThread"
                            :text="!selectedThread?.deletedAt ? $t('common.delete') : $t('common.restore')"
                        >
                            <UButton
                                :color="!selectedThread.deletedAt ? 'error' : 'success'"
                                :icon="!selectedThread.deletedAt ? 'i-mdi-delete-outline' : 'i-mdi-restore'"
                                variant="ghost"
                                @click="deleteThread"
                            />
                        </UTooltip>
                    </template>

                    <UDropdownMenu v-else :items="mobileActionItems" :content="{ align: 'end' }" :ui="{ content: 'w-64' }">
                        <UButton
                            color="neutral"
                            variant="ghost"
                            icon="i-mdi-dots-vertical"
                            :aria-label="$t('common.action', 2)"
                        />
                    </UDropdownMenu>
                </template>
            </UDashboardToolbar>

            <div class="flex flex-col justify-between gap-0.5 border-b border-default p-2 sm:flex-row sm:px-4">
                <USkeleton v-if="isRequestPending(status)" class="h-12 w-full" />

                <template v-else-if="thread">
                    <div class="flex min-w-0 flex-1 items-start gap-3 sm:my-1">
                        <div class="min-w-0">
                            <span class="text-sm font-semibold">{{ $t('common.participant', 2) }}:</span>

                            <span class="inline-flex flex-wrap items-center gap-1 align-middle">
                                <EmailInfoPopover
                                    v-for="recipient in visibleRecipients"
                                    :key="recipient.emailId"
                                    :email="recipient.email?.email"
                                    variant="link"
                                    color="primary"
                                    :padded="false"
                                    :ui="{ base: 'px-1.5 py-0.5' }"
                                />

                                <UPopover v-if="hiddenRecipients.length > 0">
                                    <UButton
                                        size="xs"
                                        color="neutral"
                                        variant="subtle"
                                        :label="`+${hiddenRecipients.length}`"
                                        :aria-label="
                                            $t('components.mailer.more_recipients', { count: hiddenRecipients.length })
                                        "
                                    />

                                    <template #content>
                                        <div class="flex max-w-xs flex-col gap-1 p-3">
                                            <EmailInfoPopover
                                                v-for="recipient in hiddenRecipients"
                                                :key="recipient.emailId"
                                                :email="recipient.email?.email"
                                                variant="link"
                                                color="primary"
                                            />
                                        </div>
                                    </template>
                                </UPopover>
                            </span>
                        </div>
                    </div>

                    <p class="text-sm text-muted max-sm:pl-16 sm:mt-1">
                        {{
                            isToday(toDate(thread.createdAt))
                                ? $d(toDate(thread.createdAt), 'time')
                                : $d(toDate(thread.createdAt), 'date')
                        }}
                    </p>
                </template>
            </div>
        </template>

        <template #body>
            <div v-if="isRequestPending(messagesStatus)" class="flex-1 space-y-2">
                <USkeleton class="h-32 w-full" />
                <USkeleton class="h-48 w-full" />
                <USkeleton class="h-32 w-full" />
            </div>

            <template v-else>
                <div
                    v-for="(message, index) in messages?.messages"
                    :key="message.id"
                    :ref="
                        (el) => {
                            messageRefs[message.id] = el as Element;
                        }
                    "
                    class="px-2 pt-2 pb-4 sm:px-4 sm:pb-5"
                    @click="selectedMessageId = message.id"
                >
                    <USeparator v-if="startsNewDay(index)" class="mx-auto mb-3 w-full max-w-(--breakpoint-xl)">
                        <span class="text-sm text-muted">{{ $d(toDate(message.createdAt), 'long') }}</span>
                    </USeparator>

                    <UCard
                        class="mx-auto w-full max-w-(--breakpoint-xl) transition-[background-color,border-color,box-shadow] hover:border-primary-500/50 hover:bg-neutral-100 hover:shadow-md hover:shadow-primary/10 dark:hover:bg-neutral-800 dark:hover:shadow-primary/10"
                        :class="selectedMessage === message.id && 'ring-2 ring-primary'"
                        :ui="{ header: 'p-4 sm:p-5', body: 'p-4 sm:p-5', footer: 'p-3 sm:p-4' }"
                    >
                        <template #header>
                            <div class="flex items-start justify-between gap-3">
                                <div class="min-w-0 flex-1">
                                    <div class="flex min-w-0 items-center gap-1 text-sm">
                                        <span class="shrink-0 font-semibold">{{ $t('common.from') }}:</span>

                                        <EmailInfoPopover
                                            :email="message.sender?.email"
                                            variant="link"
                                            truncate
                                            :trailing="false"
                                        />
                                    </div>

                                    <h3 class="mt-1 line-clamp-2 text-xl font-bold break-all hover:line-clamp-none">
                                        {{ message.title }}
                                    </h3>
                                </div>

                                <div class="flex shrink-0 items-center gap-1 text-sm text-muted">
                                    <span v-if="index === 0 && !isToday(toDate(message.createdAt))">
                                        {{ $d(toDate(message.createdAt), 'long') }}
                                    </span>
                                    <GenericTime v-else :value="message.createdAt" type="short" />

                                    <UTooltip
                                        v-if="isSuperuser"
                                        :text="!message.deletedAt ? $t('common.delete') : $t('common.restore')"
                                    >
                                        <UButton
                                            :color="!message.deletedAt ? 'error' : 'success'"
                                            :icon="!message.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                                            :aria-label="!message.deletedAt ? $t('common.delete') : $t('common.restore')"
                                            variant="ghost"
                                            size="xs"
                                            @click.stop="
                                                confirmModal.open({
                                                    confirm: async () =>
                                                        selectedEmail?.id &&
                                                        selectedThread &&
                                                        (await mailerStore.deleteMessage({
                                                            emailId: selectedEmail.id,
                                                            threadId: selectedThread.id,
                                                            messageId: message.id,
                                                        })) &&
                                                        (await refreshMessages()),
                                                })
                                            "
                                        />
                                    </UTooltip>
                                </div>
                            </div>
                        </template>

                        <div class="min-w-0 overflow-x-auto break-words">
                            <CustomContentRenderer v-if="message.content" :value="message.content" />
                        </div>

                        <template v-if="message.data?.attachments && message.data?.attachments.length > 0" #footer>
                            <UCollapsible>
                                <UButton
                                    :label="`${$t('common.attachment', 2)} (${message.data.attachments.length})`"
                                    color="neutral"
                                    variant="outline"
                                    trailing-icon="i-mdi-chevron-down"
                                    block
                                />

                                <template #content>
                                    <div class="flex flex-col gap-1 p-1">
                                        <template v-for="(attachment, idx) in message.data.attachments" :key="idx">
                                            <DocumentInfoPopover
                                                v-if="attachment.data.oneofKind === 'document'"
                                                class="flex-1"
                                                :document-id="attachment.data.document.id"
                                                button-class="flex-1 items-center"
                                                show-id
                                                load-on-open
                                                disable-tooltip
                                            />
                                        </template>
                                    </div>
                                </template>
                            </UCollapsible>
                        </template>
                    </UCard>
                </div>
            </template>
        </template>

        <template #footer>
            <Pagination
                v-if="messages?.pagination"
                v-model="messagePage"
                :pagination="messages?.pagination"
                :status="messagesStatus"
                :refresh="refreshMessages"
                compact
            />

            <UDashboardToolbar
                v-if="thread && canAccess(selectedEmail?.access, selectedEmail?.userId, AccessLevel.WRITE)"
                class="flex justify-between overflow-y-hidden border-t border-b-0 border-default"
            >
                <UCollapsible
                    class="mx-auto my-1 flex w-full max-w-(--breakpoint-xl) flex-1 flex-col gap-1"
                    :default-open="!isMobile"
                    :unmount-on-hide="false"
                    :ui="{ content: 'max-h-[50vh] overflow-y-auto' }"
                >
                    <UButton
                        class="w-full"
                        :label="$t('components.mailer.reply')"
                        icon="i-mdi-paper-airplane"
                        variant="subtle"
                        color="neutral"
                        block
                        truncate
                    />

                    <template #content>
                        <UCard class="mt-auto" variant="subtle" :ui="{ body: 'min-w-0 p-2 sm:p-2' }">
                            <UForm
                                class="flex flex-1 grow-0 flex-col gap-2 px-1"
                                :schema="schema"
                                :state="state"
                                @submit="submit"
                            >
                                <UFormField class="flex-1" name="recipients" :label="$t('common.additional_recipients')">
                                    <ClientOnly>
                                        <USelectMenu
                                            v-model="state.recipients"
                                            class="w-full"
                                            multiple
                                            trailing
                                            :items="[...state.recipients, ...addressBook]"
                                            :search-input="{ placeholder: $t('common.recipient', 1) }"
                                            :placeholder="$t('common.recipient')"
                                            creatable
                                            :disabled="!canSubmit"
                                            @create="(item: string) => onCreate(item)"
                                        >
                                            <template #default>&nbsp;</template>

                                            <template #empty>
                                                {{ $t('common.not_found', [$t('common.recipient', 2)]) }}
                                            </template>
                                        </USelectMenu>
                                    </ClientOnly>

                                    <div v-if="state.recipients.length > 0" class="mt-2 flex flex-row flex-wrap gap-2">
                                        <UFieldGroup
                                            v-for="(recipient, idx) in state.recipients"
                                            :key="recipient.label"
                                            class="max-w-full"
                                            size="sm"
                                            orientation="horizontal"
                                        >
                                            <UButton
                                                class="max-w-[calc(100%-2rem)] truncate"
                                                variant="solid"
                                                color="neutral"
                                                :label="recipient.label"
                                            />

                                            <UButton
                                                variant="outline"
                                                icon="i-mdi-clear"
                                                color="error"
                                                :aria-label="$t('common.remove')"
                                                @click="state.recipients.splice(idx, 1)"
                                            />
                                        </UFieldGroup>
                                    </div>
                                </UFormField>

                                <UFormField name="title" :label="$t('common.title')">
                                    <div class="flex flex-1 flex-col items-center gap-2 sm:flex-row">
                                        <UInput
                                            v-model="state.title"
                                            class="w-full font-semibold text-highlighted"
                                            type="text"
                                            size="lg"
                                            :placeholder="$t('common.title')"
                                            :disabled="!canSubmit"
                                            :ui="{ trailing: 'pe-1' }"
                                        >
                                            <template #trailing>
                                                <UButton
                                                    v-if="state.title !== ''"
                                                    color="neutral"
                                                    variant="link"
                                                    icon="i-mdi-close"
                                                    aria-controls="search"
                                                    @click="state.title = generateResponseTitle(selectedThread)"
                                                />
                                            </template>
                                        </UInput>

                                        <TemplateSelector
                                            v-if="editorRef"
                                            class="ml-auto"
                                            :editor="unref(editorRef).editor"
                                            size="lg"
                                        />
                                    </div>
                                </UFormField>

                                <UFormField name="content" :ui="{ error: 'hidden' }">
                                    <ClientOnly>
                                        <TiptapEditor
                                            ref="editorRef"
                                            v-model="state.content"
                                            name="content"
                                            :disabled="!canSubmit"
                                            :limit="maxContentLength"
                                            wrapper-class="min-h-44"
                                        />
                                    </ClientOnly>
                                </UFormField>

                                <div class="flex flex-col gap-2 sm:flex-row">
                                    <UButton
                                        class="w-full sm:flex-1"
                                        type="submit"
                                        :disabled="!canSubmit"
                                        :label="$t('components.mailer.send')"
                                        trailing-icon="i-mdi-paper-airplane"
                                    />

                                    <UTooltip
                                        v-if="can('documents.DocumentsService/ListDocuments').value"
                                        :text="$t('common.attachment', 2)"
                                    >
                                        <UButton
                                            class="w-full sm:w-auto"
                                            color="neutral"
                                            :label="$t('common.attachment', 2)"
                                            trailing-icon="i-mdi-attach-file"
                                            @click="
                                                threadAttachmentsModal.open({
                                                    attachments: state.attachments,
                                                    canSubmit: canSubmit,
                                                    'onUpdate:attachments': ($event) => (state.attachments = $event),
                                                })
                                            "
                                        />
                                    </UTooltip>
                                </div>
                            </UForm>
                        </UCard>
                    </template>
                </UCollapsible>
            </UDashboardToolbar>
        </template>
    </UDashboardPanel>
</template>
