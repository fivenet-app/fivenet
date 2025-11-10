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

defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'refresh'): void;
}>();

const overlay = useOverlay();

const { can, isSuperuser } = useAuth();

const notifications = useNotificationsStore();

const mailerStore = useMailerStore();
const { draft: state, addressBook, messages, selectedEmail, selectedThread } = storeToRefs(mailerStore);

const schema = z.object({
    title: z.coerce.string().min(1).max(255),
    content: z.coerce.string().min(1).max(2048),
    recipients: z
        .object({ label: z.coerce.string().min(6).max(80) })
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

async function postMessage(values: Schema): Promise<void> {
    if (!selectedEmail.value?.id) return;

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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!selectedEmail.value?.id) return;

    canSubmit.value = false;
    await postMessage(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 1000));
}, 1000);

const confirmModal = overlay.create(ConfirmModal);
const threadAttachmentsModal = overlay.create(ThreadAttachmentsModal);
</script>

<template>
    <UDashboardPanel id="mail-thread-view" :ui="{ root: 'min-h-full', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="thread?.title" :toggle="false">
                <template #title>
                    <h3 class="line-clamp-2 text-left font-semibold break-all text-highlighted hover:line-clamp-none">
                        {{ thread?.title }}
                    </h3>
                </template>

                <template #leading>
                    <UButton
                        icon="i-mdi-close"
                        color="neutral"
                        variant="ghost"
                        class="-ms-1.5"
                        @click="$emit('close', false)"
                    />
                </template>
            </UDashboardNavbar>

            <UDashboardNavbar :ui="{ toggle: 'hidden' }">
                <template #left>
                    <UTooltip
                        :text="!threadState?.unread ? $t('components.mailer.mark_unread') : $t('components.mailer.mark_read')"
                    >
                        <UButton
                            :icon="!threadState?.unread ? 'i-mdi-check-circle-outline' : 'i-mdi-check-circle'"
                            :color="!threadState?.unread ? 'neutral' : 'green'"
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

                <template #right>
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
                            @click="
                                confirmModal.open({
                                    confirm: async () => {
                                        selectedThread!.state = await mailerStore.setThreadState(
                                            {
                                                threadId: selectedThread!.id,
                                                archived: !threadState?.archived,
                                            },
                                            true,
                                        );
                                        $emit('refresh');
                                    },
                                })
                            "
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
                            @click="
                                confirmModal.open({
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

            <div class="flex flex-col justify-between gap-1 border-b border-default p-4 sm:flex-row sm:px-6">
                <USkeleton v-if="isRequestPending(status)" class="h-12 w-full" />

                <template v-else-if="thread">
                    <div class="flex items-start gap-4 sm:my-1.5">
                        <div class="min-w-0">
                            <span class="text-sm font-semibold">{{ $t('common.participant', 2) }}:</span>

                            <EmailInfoPopover
                                v-for="recipient in thread.recipients"
                                :key="recipient.emailId"
                                :email="recipient.email?.email"
                                variant="link"
                                color="primary"
                                :padded="false"
                                :ui="{ base: 'px-2 py-1' }"
                            />
                        </div>
                    </div>

                    <p class="text-sm text-muted max-sm:pl-16 sm:mt-2">
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
                    v-for="message in messages?.messages"
                    :key="message.id"
                    :ref="
                        (el) => {
                            messageRefs[message.id] = el as Element;
                        }
                    "
                    class="border-l-2 border-default px-2 pb-3 hover:border-primary-500 hover:bg-neutral-100 sm:pb-2 hover:dark:border-primary-400 dark:hover:bg-neutral-800"
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

                    <div
                        class="mx-auto w-full max-w-(--breakpoint-xl) rounded-lg bg-neutral-100 p-4 break-words dark:bg-neutral-800"
                    >
                        <HTMLContent v-if="message.content?.content" :value="message.content.content" />
                    </div>

                    <UCollapsible v-if="message.data?.attachments && message.data?.attachments.length > 0" class="my-2">
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
                </div>
            </template>
        </template>

        <template #footer>
            <Pagination
                v-if="messages?.pagination"
                v-model="page"
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
                    class="my-1 flex flex-1 flex-col gap-1"
                    :unmount-on-hide="false"
                    :ui="{ content: 'max-h-[50vh] overflow-y-auto' }"
                >
                    <UButton
                        :label="$t('components.mailer.reply')"
                        icon="i-mdi-paper-airplane"
                        variant="subtle"
                        color="neutral"
                        class="w-full"
                        block
                        truncate
                    />

                    <template #content>
                        <UCard variant="subtle" class="mt-auto" :ui="{ body: 'min-w-0 p-2 sm:p-2' }">
                            <UForm
                                class="flex flex-1 grow-0 flex-col gap-2 px-1"
                                :schema="schema"
                                :state="state"
                                @submit="onSubmitThrottle"
                            >
                                <UFormField name="recipients" :label="$t('common.additional_recipients')" class="flex-1">
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

                                    <div
                                        v-if="state.recipients.length > 0"
                                        class="mt-2 flex snap-x flex-row flex-wrap gap-2 overflow-x-auto"
                                    >
                                        <UFieldGroup
                                            v-for="(recipient, idx) in state.recipients"
                                            :key="idx"
                                            size="sm"
                                            orientation="horizontal"
                                        >
                                            <UButton variant="solid" color="neutral" :label="recipient.label" />

                                            <UButton
                                                variant="outline"
                                                icon="i-mdi-clear"
                                                color="error"
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

                                        <TemplateSelector v-model="state.content" class="ml-auto" size="lg" />
                                    </div>
                                </UFormField>

                                <UFormField name="content" :ui="{ error: 'hidden' }">
                                    <ClientOnly>
                                        <TiptapEditor
                                            v-model="state.content"
                                            name="content"
                                            :disabled="!canSubmit"
                                            wrapper-class="min-h-44"
                                        />
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
                        </UCard>
                    </template>
                </UCollapsible>
            </UDashboardToolbar>
        </template>
    </UDashboardPanel>
</template>
