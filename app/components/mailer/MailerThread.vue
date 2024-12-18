<script setup lang="ts">
import type { FormSubmitEvent } from '#ui/types';
import { isToday } from 'date-fns';
import { z } from 'zod';
import { useMailerStore } from '~/store/mailer';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel } from '~~/gen/ts/resources/mailer/access';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import ConfirmModal from '../partials/ConfirmModal.vue';
import HTMLContent from '../partials/content/HTMLContent.vue';
import GenericTime from '../partials/elements/GenericTime.vue';
import Pagination from '../partials/Pagination.vue';
import TiptapEditor from '../partials/TiptapEditor.vue';
import EmailInfoPopover from './EmailInfoPopover.vue';
import { canAccess, generateResponseTitle } from './helpers';
import TemplateSelector from './TemplateSelector.vue';

const props = withDefaults(
    defineProps<{
        threadId: string;
        selected?: boolean;
    }>(),
    {
        selected: false,
    },
);

const modal = useModal();

const { isSuperuser } = useAuth();

const notifications = useNotificatorStore();

const mailerStore = useMailerStore();
const { draft: state, addressBook, messages, selectedEmail, selectedThread } = storeToRefs(mailerStore);

const schema = z.object({
    title: z.string().min(1).max(255),
    content: z.string().min(1).max(2048),
    recipients: z
        .object({ label: z.string().min(6).max(80) })
        .array()
        .max(20),
});

type Schema = z.output<typeof schema>;

const { data: thread, pending: loading } = useLazyAsyncData(
    `mailer-thread:${props.threadId}`,
    () => mailerStore.getThread(props.threadId),
    {
        watch: [() => props.threadId],
    },
);

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() =>
    messages.value?.pagination?.pageSize ? messages.value?.pagination?.pageSize * (page.value - 1) : 0,
);

const { pending: messagesLoading, refresh: refreshMessages } = useLazyAsyncData(
    `mailer-thread:${props.threadId}-messages:${page.value}`,
    async () => {
        const response = await mailerStore.listThreadMessages({
            pagination: {
                offset: offset.value,
            },
            emailId: selectedEmail.value!.id,
            threadId: props.threadId,
        });

        if (selectedThread.value) {
            if (state.value.title === '') {
                state.value.title = generateResponseTitle(selectedThread.value);
            }

            if ((!state.value.content || state.value.content === '<p><br></p>') && !!selectedEmail.value?.settings?.signature) {
                state.value.content = '<p><br></p><p><br></p>' + selectedEmail.value?.settings?.signature;
            }
        }

        return response;
    },
    { watch: [() => props.threadId] },
);

watch(offset, async () => refreshMessages());

watchDebounced(
    selectedThread,
    async () =>
        canAccess(selectedEmail.value?.access, selectedEmail.value?.userId, AccessLevel.WRITE) &&
        (await mailerStore.setThreadState({
            threadId: props.threadId,
            unread: false,
        })),
    {
        debounce: 1150,
        maxWait: 3250,
    },
);

async function postMessage(values: Schema): Promise<void> {
    if (!selectedEmail.value?.id) {
        return;
    }

    await mailerStore.postMessage({
        message: {
            id: '0',
            senderId: selectedEmail.value.id,
            threadId: props.threadId,
            title: values.title,
            content: {
                rawContent: values.content,
            },
            data: {
                entry: [],
            },
        },
        recipients: [...new Set(values.recipients.map((r) => r.label.trim()))],
    });

    notifications.add({
        title: { key: 'notifications.action_successfull.title', parameters: {} },
        description: { key: 'notifications.action_successfull.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });

    // Clear draft data
    state.value.title = '';
    state.value.content = '';
    state.value.recipients = [];
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

const selectedMessageId = useRouteQuery('msg', '0', { transform: Number });
const selectedMessage = computed(() => selectedMessageId.value.toString());
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
</script>

<template>
    <UDashboardToolbar :ui="{ container: 'flex-col gap-y-2' }">
        <USkeleton v-if="loading" class="h-12 w-full" />

        <template v-else-if="thread">
            <div class="flex w-full flex-1 items-center justify-between gap-1">
                <h3 class="line-clamp-2 text-left font-semibold text-gray-900 hover:line-clamp-none dark:text-white">
                    {{ thread.title }}
                </h3>

                <p class="shrink-0 font-medium text-gray-900 dark:text-white">
                    {{
                        isToday(toDate(thread.createdAt))
                            ? $d(toDate(thread.createdAt), 'time')
                            : $d(toDate(thread.createdAt), 'date')
                    }}
                </p>
            </div>

            <div class="w-full min-w-0 flex-1 text-sm">
                <div class="flex snap-x flex-row flex-wrap gap-1 overflow-x-auto text-gray-500 dark:text-gray-400">
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
        <div v-if="messagesLoading" class="flex-1 space-y-2">
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
                        messageRefs[parseInt(message.id)] = el as Element;
                    }
                "
                class="hover:border-primary-500 hover:dark:border-primary-400 border-l-2 border-white px-2 pb-3 hover:bg-neutral-100 sm:pb-2 dark:border-gray-900 dark:hover:bg-base-800"
                :class="selectedMessage === message.id && '!border-primary-500'"
                @click="selectedMessageId = parseInt(message.id)"
            >
                <UDivider>
                    <GenericTime :value="message.createdAt" type="short" />

                    <UTooltip v-if="isSuperuser" :text="$t('common.delete')" square class="absolute right-0">
                        <UButton
                            icon="i-mdi-trash-can-outline"
                            color="red"
                            variant="ghost"
                            size="xs"
                            @click="
                                modal.open(ConfirmModal, {
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
                </UDivider>

                <div class="flex flex-col gap-1">
                    <div class="inline-flex items-center gap-1 text-sm">
                        <span class="font-semibold">{{ $t('common.from') }}:</span>

                        <EmailInfoPopover
                            :email="message.sender?.email"
                            variant="link"
                            truncate
                            :trailing="false"
                            :padded="false"
                        />
                    </div>

                    <div class="inline-flex items-center gap-1">
                        <span class="text-sm font-semibold">{{ $t('common.title') }}:</span>
                        <h3 class="line-clamp-2 text-xl font-bold hover:line-clamp-none">{{ message.title }}</h3>
                    </div>
                </div>

                <div class="mx-auto w-full max-w-screen-xl break-words rounded-lg bg-neutral-100 dark:bg-base-900">
                    <HTMLContent v-if="message.content?.content" class="px-4 py-2" :value="message.content.content" />
                </div>
            </div>
        </div>

        <Pagination v-if="messages?.pagination" v-model="page" :pagination="messages?.pagination" :refresh="refreshMessages" />

        <UDashboardToolbar
            v-if="thread && canAccess(selectedEmail?.access, selectedEmail?.userId, AccessLevel.WRITE)"
            class="flex justify-between border-t border-gray-200 px-3 py-3.5 dark:border-gray-700"
        >
            <UAccordion
                variant="outline"
                :items="[{ slot: 'compose', label: $t('components.mailer.reply'), icon: 'i-mdi-paper-airplane' }]"
                :ui="{ default: { class: 'mb-0' } }"
            >
                <template #compose>
                    <UForm :schema="schema" :state="state" class="flex flex-col gap-2" @submit="onSubmitThrottle">
                        <UFormGroup name="recipients" class="w-full flex-1" :label="$t('common.additional_recipients')">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.recipients"
                                    :placeholder="$t('common.recipient')"
                                    block
                                    multiple
                                    trailing
                                    searchable
                                    :options="[...state.recipients, ...addressBook]"
                                    :searchable-placeholder="$t('common.recipient')"
                                    creatable
                                    :disabled="!canSubmit"
                                >
                                    <template #label>&nbsp;</template>

                                    <template #option-create="{ option }">
                                        <span class="shrink-0">{{ $t('common.recipient') }}: {{ option.label }}</span>
                                    </template>

                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.recipient', 2)]) }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>

                            <div class="mt-2 flex snap-x flex-row flex-wrap gap-2 overflow-x-auto">
                                <UButtonGroup
                                    v-for="(recipient, idx) in state.recipients"
                                    :key="idx"
                                    size="sm"
                                    orientation="horizontal"
                                >
                                    <UButton variant="solid" color="gray" :label="recipient.label" />

                                    <UButton
                                        variant="outline"
                                        icon="i-mdi-close"
                                        color="red"
                                        @click="state.recipients.splice(idx, 1)"
                                    />
                                </UButtonGroup>
                            </div>
                        </UFormGroup>

                        <UFormGroup name="title" :label="$t('common.title')" class="w-full flex-1">
                            <div class="flex flex-1 flex-col items-center gap-2 sm:flex-row">
                                <UInput
                                    v-model="state.title"
                                    type="text"
                                    size="lg"
                                    class="w-full font-semibold text-gray-900 dark:text-white"
                                    :placeholder="$t('common.title')"
                                    :disabled="!canSubmit"
                                    :ui="{ icon: { trailing: { pointer: '' } } }"
                                >
                                    <template #trailing>
                                        <UButton
                                            v-show="state.title !== ''"
                                            color="gray"
                                            variant="link"
                                            icon="i-mdi-close"
                                            :padded="false"
                                            @click="state.title = generateResponseTitle(selectedThread)"
                                        />
                                    </template>
                                </UInput>

                                <TemplateSelector v-model="state.content" size="lg" class="ml-auto" />
                            </div>
                        </UFormGroup>

                        <UFormGroup name="message">
                            <ClientOnly>
                                <TiptapEditor v-model="state.content" :disabled="!canSubmit" wrapper-class="min-h-44" />
                            </ClientOnly>
                        </UFormGroup>

                        <UButton
                            type="submit"
                            :disabled="!canSubmit"
                            block
                            class="flex-1"
                            :label="$t('components.mailer.send')"
                            trailing-icon="i-mdi-paper-airplane"
                        />
                    </UForm>
                </template>
            </UAccordion>
        </UDashboardToolbar>
    </UDashboardPanelContent>
</template>
