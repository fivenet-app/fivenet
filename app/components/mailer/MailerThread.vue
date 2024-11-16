<script setup lang="ts">
import type { FormSubmitEvent } from '#ui/types';
import { isToday } from 'date-fns';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { mailerDB, useMailerStore } from '~/store/mailer';
import DocEditor from '../partials/DocEditor.vue';

const props = withDefaults(
    defineProps<{
        threadId: string;
        selected?: boolean;
    }>(),
    {
        selected: false,
    },
);

const mailerStore = useMailerStore();

const schema = z.object({
    title: z.string().min(1).max(255),
    content: z.string().min(1).max(2048),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    title: '',
    content: '',
});

const { data: thread, pending: loading } = useLazyAsyncData(`mailer-thread:${props.threadId}`, async () =>
    mailerStore.getThread(props.threadId),
);

onBeforeMount(async () => {
    const count = await mailerDB.threads.count();
    const call = getGRPCMailerClient().listThreadMessages({
        threadId: props.threadId,
        after: count > 0 ? undefined : toTimestamp(),
    });
    const { response } = await call;

    await mailerDB.messages.bulkPut(response.messages);
});

watchDebounced(
    () => props.threadId,
    async () =>
        mailerStore.setThreadState({
            threadId: props.threadId,
            unread: false,
        }),
);

const messages = useDexieLiveQueryWithDeps(
    () => props.threadId,
    () =>
        mailerDB.messages
            .where('threadId')
            .equals(props.threadId)
            .limit(2500)
            .sortBy('id')
            .then((messages) => ({ messages, loaded: true })),
    {
        initialValue: { messages: [], loaded: false },
    },
);

const messageRef = ref<Element | undefined>();
watchDebounced(messages, () => messageRef.value?.scrollIntoView({ behavior: 'smooth' }), {
    debounce: 100,
    maxWait: 350,
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await mailerStore
        .postMessage({
            message: {
                id: '0',
                threadId: props.threadId,
                title: event.data.title,
                content: event.data.content,
                data: {
                    entry: [],
                },
            },
        })
        .then(() => (state.content = ''))
        .finally(() => useTimeoutFn(() => (canSubmit.value = true), 1000));
}, 1000);
</script>

<template>
    <UDashboardPanelContent>
        <USkeleton v-if="!thread && loading" class="h-12 w-full" />

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

                <div class="min-w-0">
                    <div class="font-medium text-gray-500 dark:text-gray-400">
                        <UPopover>
                            <UButton block variant="link" :padded="false">
                                <UAvatarGroup size="sm" :max="5">
                                    <ProfilePictureImg
                                        v-if="thread.creator"
                                        :src="thread.creator?.avatar?.url"
                                        :name="`${thread.creator?.firstname} ${thread.creator?.lastname}`"
                                        disable-blur-toggle
                                    />

                                    <div v-for="recipient in thread.recipients" :key="recipient.emailId">
                                        {{ recipient.emailId }}
                                    </div>
                                </UAvatarGroup>
                            </UButton>

                            <template #panel>
                                <div class="p-4 text-gray-900 dark:text-white">
                                    <ul role="list">
                                        <li v-if="thread.creator">
                                            <CitizenInfoPopover :user="thread.creator" show-avatar-in-name />
                                        </li>
                                        <li v-for="ua in thread.recipients" :key="ua.emailId">
                                            {{ ua.emailId }}
                                        </li>
                                    </ul>
                                </div>
                            </template>
                        </UPopover>
                    </div>
                </div>
            </div>
        </div>

        <UDivider class="my-2" />

        <div class="relative -mx-4 flex-1 overflow-x-auto">
            <template v-if="!messages.loaded">
                <div class="space-y-2">
                    <USkeleton class="h-6 w-full" />
                    <USkeleton class="h-6 w-full" />
                </div>
            </template>
            <template v-else>
                <template v-for="message in messages.messages" :key="message.id">
                    <div
                        class="hover:border-primary-500 hover:dark:border-primary-400 border-l-2 border-white px-2 hover:bg-base-800 dark:border-gray-900"
                    >
                        <UDivider class="text-xs">
                            <GenericTime :value="message.createdAt" :type="'date'" />
                        </UDivider>

                        <div v-if="message.creator" class="flex justify-between text-xs">
                            <CitizenInfoPopover :user="message.creator" show-avatar-in-name />

                            <GenericTime :value="message.createdAt" type="time" />
                        </div>

                        <div class="flex justify-between text-xs">
                            <!-- eslint-disable vue/no-v-html -->
                            <div
                                :ref="
                                    (el) => {
                                        if (messages.messages.length) {
                                            messageRef = el as Element;
                                        }
                                    }
                                "
                                class="prose dark:prose-invert min-w-full px-4 py-2"
                                v-html="message.content"
                            ></div>
                        </div>
                    </div>
                </template>
            </template>
        </div>

        <UDivider class="my-2" />

        <UForm v-if="thread" :schema="schema" :state="state" class="flex flex-col gap-2" @submit="onSubmitThrottle">
            <!-- TODO add recipients field -->

            <UFormGroup name="message">
                <ClientOnly>
                    <DocEditor v-model="state.content" :disabled="!canSubmit" />
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
    </UDashboardPanelContent>
</template>
