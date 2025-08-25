<script setup lang="ts">
import type { FormSubmitEvent } from '@nuxt/ui';
import { isToday } from 'date-fns';
import { z } from 'zod';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import HTMLContent from '~/components/partials/content/HTMLContent.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useMailerStore } from '~/stores/mailer';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/message';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
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

const overlay = useOverlay();

const { can, isSuperuser } = useAuth();

const notifications = useNotificationsStore();

const mailerStore = useMailerStore();
const { draft: state, addressBook, messages, selectedEmail, selectedThread } = storeToRefs(mailerStore);

const schema = z.object({
    title: z.string().min(1).max(255),
    content: z.string().min(1).max(2048),
    recipients: z
        .object({ label: z.string().min(6).max(80) })
        .array()
        .max(20)
        .default([]),
    attachments: z.custom<MessageAttachment>().array().max(3).default([]),
});

type Schema = z.output<typeof schema>;

function resetForm(): void {
    if (selectedThread.value) {
        if (state.value.title === '') {
            state.value.title = generateResponseTitle(selectedThread.value);
        }

        if ((!state.value.content || state.value.content === '<p><br></p>') && !!selectedEmail.value?.settings?.signature) {
            state.value.content = '<p><br></p><p><br></p>' + selectedEmail.value?.settings?.signature;
        }
    }
}

const { data: thread, status } = useLazyAsyncData(
    `mailer-thread:${props.threadId}`,
    () => mailerStore.getThread(props.threadId),
    {
        watch: [() => props.threadId],
    },
);

const page = useRouteQuery('page', '1', { transform: Number });

const { status: messagesStatus, refresh: refreshMessages } = useLazyAsyncData(
    `mailer-thread:${props.threadId}-messages:${page.value}`,
    async () => {
        const response = await mailerStore.listThreadMessages({
            pagination: {
                offset: calculateOffset(page.value, messages.value?.pagination),
            },
            emailId: selectedEmail.value!.id,
            threadId: props.threadId,
        });

        resetForm();

        return response;
    },
    { watch: [() => props.threadId] },
);

watchDebounced(
    () => props.threadId,
    async () => {
        if (!thread.value?.state?.unread) {
            return;
        }

        if (!canAccess(selectedEmail.value?.access, selectedEmail.value?.userId, AccessLevel.WRITE)) {
            return;
        }

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

async function postMessage(values: Schema): Promise<void> {
    if (!selectedEmail.value?.id) {
        return;
    }

    await mailerStore.postMessage({
        message: {
            id: 0,
            senderId: selectedEmail.value.id,
            threadId: props.threadId,
            title: values.title,
            content: {
                rawContent: values.content,
            },
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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!selectedEmail.value?.id) {
        return;
    }

    canSubmit.value = false;
    await postMessage(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 1000));
}, 1000);

const confirmModal = overlay.create(ConfirmModal);
const threadAttachmentsModal = overlay.create(ThreadAttachmentsModal);
</script>

<template>
    <UDashboardToolbar :ui="{ container: 'flex-col gap-y-2' }">
        <USkeleton v-if="isRequestPending(status)" class="h-12 w-full" />

        <template v-else-if="thread">
            <div class="flex w-full flex-1 items-center justify-between gap-1">
                <h3 class="line-clamp-2 text-left font-semibold break-all text-gray-900 hover:line-clamp-none dark:text-white">
                    {{ thread.title }}
                </h3>

                <p class="shrink-0 font-medium text-highlighted">
                    {{
                        isToday(toDate(thread.createdAt))
                            ? $d(toDate(thread.createdAt), 'time')
                            : $d(toDate(thread.createdAt), 'date')
                    }}
                </p>
            </div>

            <div class="w-full min-w-0 flex-1 text-sm">
                <div class="flex snap-x flex-row flex-wrap gap-1 overflow-x-auto text-muted">
                    <span class="text-sm font-semibold">{{ $t('common.participant', 2) }}:</span>

                    <EmailInfoPopover
                        v-for="recipient in thread.recipients"
                        :key="recipient.emailId"
                        :email="recipient.email?.email"
                    />
                </div>
            </div>
        </template>
    </UDashboardToolbar>

    <UDashboardPanelContent class="p-0 sm:pb-0">
        <div v-if="isRequestPending(messagesStatus)" class="flex-1 space-y-2">
            <USkeleton class="h-32 w-full" />
            <USkeleton class="h-48 w-full" />
            <USkeleton class="h-32 w-full" />
        </div>

        <div v-else class="flex flex-1 shrink-0 flex-col overflow-y-auto">
            <div
                v-for="message in messages?.messages"
                :key="message.id"
                :ref="
                    (el) => {
                        messageRefs[message.id] = el as Element;
                    }
                "
                class="dark:hover:bg-base-800 border-l-2 border-white px-2 pb-3 hover:border-primary-500 hover:bg-neutral-100 sm:pb-2 dark:border-gray-900 hover:dark:border-primary-400"
                :class="selectedMessage === message.id && '!border-primary-500'"
                @click="selectedMessageId = message.id"
            >
                <USeparator>
                    <GenericTime :value="message.createdAt" type="short" />

                    <UTooltip
                        v-if="isSuperuser"
                        class="ml-2"
                        :text="!message.deletedAt ? $t('common.delete') : $t('common.restore')"
                    >
                        <UButton
                            :color="!message.deletedAt ? 'error' : 'success'"
                            :icon="!message.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                            variant="ghost"
                            size="xs"
                            @click="
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
                </USeparator>

                <div class="flex flex-col gap-1">
                    <div class="inline-flex items-center gap-1 text-sm">
                        <span class="font-semibold">{{ $t('common.from') }}:</span>

                        <EmailInfoPopover :email="message.sender?.email" variant="link" truncate :trailing="false" />
                    </div>

                    <div class="inline-flex items-center gap-1">
                        <span class="text-sm font-semibold">{{ $t('common.title') }}:</span>
                        <h3 class="line-clamp-2 text-xl font-bold break-all hover:line-clamp-none">{{ message.title }}</h3>
                    </div>
                </div>

                <div class="dark:bg-base-900 mx-auto w-full max-w-(--breakpoint-xl) rounded-lg bg-neutral-100 break-words">
                    <HTMLContent v-if="message.content?.content" class="px-4 py-2" :value="message.content.content" />
                </div>

                <UAccordion
                    v-if="message.data?.attachments && message.data?.attachments.length > 0"
                    class="mt-1"
                    :items="[
                        {
                            slot: 'attachments' as const,
                            label: $t('common.attachment', 2),
                            color: 'neutral',
                            variant: 'outline',
                        },
                    ]"
                >
                    <template #attachments>
                        <div class="flex flex-col gap-1">
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
                </UAccordion>
            </div>
        </div>

        <Pagination
            v-if="messages?.pagination"
            v-model="page"
            :pagination="messages?.pagination"
            :loading="isRequestPending(messagesStatus)"
            :refresh="refreshMessages"
        />
    </UDashboardPanelContent>

    <UDashboardToolbar
        v-if="thread && canAccess(selectedEmail?.access, selectedEmail?.userId, AccessLevel.WRITE)"
        class="flex min-w-0 justify-between overflow-y-hidden border-t border-b-0 border-gray-200 dark:border-gray-700"
        :ui="{
            container: 'gap-x-0 justify-stretch items-stretch h-full inline-flex flex-col p-0 px-1 min-w-0',
        }"
    >
        <UAccordion
            class="mt-2 max-h-[50vh] overflow-y-auto"
            variant="outline"
            :items="[{ slot: 'compose' as const, label: $t('components.mailer.reply'), icon: 'i-mdi-paper-airplane' }]"
            :ui="{ default: { class: 'mb-0' } }"
        >
            <template #compose>
                <UForm
                    class="flex flex-1 grow-0 flex-col gap-2 px-1"
                    :schema="schema"
                    :state="state"
                    @submit="onSubmitThrottle"
                >
                    <UFormField name="recipients" :label="$t('common.additional_recipients')">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.recipients"
                                multiple
                                trailing
                                :items="[...state.recipients, ...addressBook]"
                                searchable
                                :searchable-placeholder="$t('common.recipient')"
                                :placeholder="$t('common.recipient')"
                                creatable
                                :disabled="!canSubmit"
                            >
                                <template #item-label>&nbsp;</template>

                                <template #option-create="{ option }">
                                    <span class="shrink-0">{{ $t('common.recipient') }}: {{ option.label }}</span>
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.recipient', 2)]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>

                        <div
                            v-if="state.recipients.length > 0"
                            class="mt-2 flex snap-x flex-row flex-wrap gap-2 overflow-x-auto"
                        >
                            <UButtonGroup
                                v-for="(recipient, idx) in state.recipients"
                                :key="idx"
                                size="sm"
                                orientation="horizontal"
                            >
                                <UButton variant="solid" color="neutral" :label="recipient.label" />

                                <UButton
                                    variant="outline"
                                    icon="i-mdi-close"
                                    color="error"
                                    @click="state.recipients.splice(idx, 1)"
                                />
                            </UButtonGroup>
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
                                        v-show="state.title !== ''"
                                        color="neutral"
                                        variant="link"
                                        icon="i-mdi-close"
                                        aria-controls="search"
                                        @click="state.title = generateResponseTitle(selectedThread)"
                                    />
                                </template>
                            </UInput>

                            <TemplateSelector v-model="state.content" class="ml-auto" size="lg" />
                        </div>
                    </UFormField>

                    <UFormField name="message">
                        <ClientOnly>
                            <TiptapEditor v-model="state.content" :disabled="!canSubmit" wrapper-class="min-h-44" />
                        </ClientOnly>
                    </UFormField>

                    <div class="inline-flex gap-1">
                        <UButton
                            class="flex-1"
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
                                color="neutral"
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
            </template>
        </UAccordion>
    </UDashboardToolbar>
</template>
