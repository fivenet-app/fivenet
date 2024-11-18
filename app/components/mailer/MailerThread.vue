<script setup lang="ts">
import type { FormSubmitEvent } from '#ui/types';
import { isToday } from 'date-fns';
import { z } from 'zod';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { mailerDB, useMailerStore } from '~/store/mailer';
import ConfirmModal from '../partials/ConfirmModal.vue';
import DocEditor from '../partials/DocEditor.vue';
import Pagination from '../partials/Pagination.vue';

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

const mailerStore = useMailerStore();
const { draft: state, selectedEmail, selectedThread } = storeToRefs(mailerStore);

const schema = z.object({
    title: z.string().min(1).max(255),
    content: z.string().min(1).max(2048),
});

type Schema = z.output<typeof schema>;

const { data: thread, pending: loading } = useLazyAsyncData(
    `mailer-thread:${props.threadId}`,
    () => mailerStore.getThread(props.threadId),
    {
        watch: [() => props.threadId],
    },
);

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const messages = useDexieLiveQuery(
    () =>
        mailerDB.messages
            .where({ threadId: props.threadId })
            .reverse()
            .toArray()
            .then((messages) => ({ messages: messages, loaded: true })),
    {
        initialValue: { messages: [], loaded: false },
    },
);

const {
    data,
    pending: messagesLoading,
    refresh: refreshMessages,
} = useLazyAsyncData(
    `mailer-thread:${props.threadId}-messages:${page.value}`,
    async () => {
        const count = await mailerDB.threads.count();

        const response = await mailerStore.listThreadMessages({
            pagination: {
                offset: offset.value,
            },
            emailId: selectedEmail.value!.id,
            threadId: props.threadId,
            after: count > 0 ? undefined : toTimestamp(),
        });

        if (selectedThread.value) {
            if (state.value.title === '') {
                state.value.title = 'RE: ' + selectedThread.value.title;
            }

            if ((!state.value.content || state.value.content === '<p><br></p>') && !!selectedEmail.value?.settings?.signature) {
                state.value.content = '<p><br></p><p><br></p>' + selectedEmail.value?.settings?.signature;
            }
        }

        return response;
    },
    { watch: [() => props.threadId] },
);

const { start } = useTimeoutFn(
    async () => {
        return mailerStore.setThreadState({
            threadId: props.threadId,
            unread: false,
        });
    },
    1250,
    { immediate: false },
);

onMounted(() => start());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (!selectedEmail.value?.id) {
        return;
    }

    canSubmit.value = false;
    await mailerStore
        .postMessage({
            message: {
                id: '0',
                senderId: selectedEmail.value?.id,
                threadId: props.threadId,
                title: event.data.title,
                content: event.data.content,
                data: {
                    entry: [],
                },
            },
            recipients: [],
        })
        .then(() => {
            state.value.title = '';
            state.value.content = '';
            state.value.recipients = [];
        })
        .finally(() => useTimeoutFn(() => (canSubmit.value = true), 1000));
}, 1000);
</script>

<template>
    <UDashboardToolbar>
        <USkeleton v-if="loading" class="h-12 w-full" />

        <div v-else-if="thread" class="flex w-full">
            <div class="flex w-full flex-col gap-2">
                <div class="flex flex-1 items-center justify-between gap-1">
                    <p class="font-semibold text-gray-900 dark:text-white">
                        {{ thread.title }}
                    </p>

                    <p class="font-medium text-gray-900 dark:text-white">
                        {{
                            isToday(toDate(thread.createdAt))
                                ? $d(toDate(thread.createdAt), 'time')
                                : $d(toDate(thread.createdAt), 'date')
                        }}
                    </p>
                </div>

                <div class="min-w-0 text-sm">
                    <div class="flex snap-x flex-row flex-wrap gap-1 overflow-x-auto text-gray-500 dark:text-gray-400">
                        <span class="text-sm font-semibold">{{ $t('common.participant', 2) }}:</span>

                        <template v-for="(recipient, idx) in thread.recipients" :key="recipient.emailId">
                            <UButton variant="link" :padded="false" :label="recipient.email?.email" />

                            <span v-if="thread.recipients.length - 1 !== idx">, </span>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </UDashboardToolbar>

    <UDashboardPanelContent>
        <div class="relative -mx-4 flex-1 overflow-x-auto">
            <template v-if="messagesLoading">
                <div class="flex-1 space-y-2">
                    <USkeleton class="h-32 w-full" />
                    <USkeleton class="h-32 w-full" />
                </div>
            </template>
            <template v-else>
                <template v-for="message in messages?.messages" :key="message.id">
                    <div
                        class="hover:border-primary-500 hover:dark:border-primary-400 border-l-2 border-white px-2 pb-3 hover:bg-base-800 sm:pb-2 dark:border-gray-900"
                    >
                        <UDivider class="relative">
                            <GenericTime :value="message.createdAt" :type="'short'" />

                            <UTooltip v-if="isSuperuser" :text="$t('common.delete')" square class="absolute right-0">
                                <UButton
                                    icon="i-mdi-trash-can-outline"
                                    color="gray"
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

                                <div class="flex justify-between">{{ message.sender?.email ?? $t('common.unknown') }}</div>
                            </div>

                            <div class="inline-flex items-center gap-1">
                                <span class="text-sm font-semibold">{{ $t('common.title') }}:</span>
                                <h3 class="truncate text-xl font-bold">{{ message.title }}</h3>
                            </div>
                        </div>

                        <div class="mx-auto max-w-screen-xl break-words rounded-lg bg-base-900">
                            <!-- eslint-disable vue/no-v-html -->
                            <div class="prose dark:prose-invert min-w-full px-4 py-2" v-html="message.content"></div>
                        </div>
                    </div>
                </template>
            </template>
        </div>

        <Pagination
            v-if="data?.pagination && data?.pagination?.totalCount / data?.pagination?.pageSize > 1"
            v-model="page"
            :pagination="data?.pagination"
        />

        <UDivider class="my-2" />

        <UAccordion
            v-if="thread"
            variant="outline"
            :items="[{ slot: 'compose', label: $t('components.mailer.reply'), icon: 'i-mdi-paper-airplane' }]"
        >
            <template #compose>
                <UForm :schema="schema" :state="state" class="flex flex-col gap-2" @submit="onSubmitThrottle">
                    <!-- TODO add "add recipients" field -->
                    <UFormGroup name="title" class="w-full flex-1">
                        <UInput
                            v-model="state.title"
                            type="text"
                            size="xl"
                            class="font-semibold text-gray-900 dark:text-white"
                            :placeholder="$t('common.title')"
                            :disabled="!canSubmit"
                        />
                    </UFormGroup>

                    <UFormGroup name="message">
                        <ClientOnly>
                            <DocEditor v-model="state.content" :disabled="!canSubmit" :min-height="250" />
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
    </UDashboardPanelContent>
</template>
