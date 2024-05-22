<script setup lang="ts">
import { isToday } from 'date-fns';
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import { messengerDB, useMessengerStore } from '~/store/messenger';
import CitizenInfoPopover from '../partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '../partials/elements/GenericTime.vue';
import { canAccessThread } from './helpers';
import { AccessLevel } from '~~/gen/ts/resources/messenger/access';
import type { Message } from '~~/gen/ts/resources/messenger/message';

const props = withDefaults(
    defineProps<{
        threadId: string;
        selected?: boolean;
    }>(),
    {
        selected: false,
    },
);

const messengerStore = useMessengerStore();

const schema = z.object({
    message: z.string().min(1).max(2048),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    message: '',
});

const {
    data: thread,
    pending: loading,
    error,
} = useLazyAsyncData(`messenger-thread:${props.threadId}`, async () => messengerStore.getThread(props.threadId));

onBeforeMount(async () => {
    const count = await messengerDB.threads.count();
    const call = getGRPCMessengerClient().getThreadMessages({
        threadId: props.threadId,
        after: count > 0 ? undefined : toTimestamp(),
    });
    const { response } = await call;

    await messengerDB.messages.bulkPut(response.messages);
});

watchDebounced(
    () => props.threadId,
    async () =>
        messengerStore.setThreadUserState({
            threadId: props.threadId,
            unread: false,
        }),
);

const messages = useDexieLiveQueryWithDeps(
    () => props.threadId,
    () =>
        messengerDB.messages
            .where('threadId')
            .equals(props.threadId)
            .limit(2500)
            .sortBy('id')
            .then((messages) => ({ messages, loaded: true })),
    {
        initialValue: { messages: [], loaded: false },
    },
);

const groupedMessages = computed(() => {
    return messages.value.messages
        .sort((a, b) => toDate(a.createdAt).getTime() - toDate(b.createdAt).getTime())
        .reduce((acc: { [key: string]: Message[] }, msg) => {
            const k = toDate(msg.createdAt).toDateString();

            acc[k] = acc[k] || [];
            if (acc[k].length > 0 && acc[k][acc[k].length - 1].creatorId === msg.creatorId) {
                msg.creator = undefined;
            }
            acc[k].push(msg);

            return acc;
        }, {});
});

const messageRef = ref<Element | undefined>();
watchDebounced(messages, () => messageRef.value?.scrollIntoView({ behavior: 'smooth' }), {
    debounce: 100,
    maxWait: 350,
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await messengerStore
        .postMessage({
            message: {
                id: '0',
                threadId: props.threadId,
                message: event.data.message,
                data: {},
            },
        })
        .then(() => (state.message = ''))
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

                                    <ProfilePictureImg
                                        v-for="user in thread.access?.users"
                                        :src="user.user?.avatar?.url"
                                        :name="`${user.user?.firstname} ${user.user?.lastname}`"
                                        disable-blur-toggle
                                    />
                                </UAvatarGroup>
                            </UButton>

                            <template #panel>
                                <div class="p-4 text-gray-900 dark:text-white">
                                    <ul role="list">
                                        <li v-if="thread.creator">
                                            <CitizenInfoPopover :user="thread.creator" show-avatar-in-name />
                                        </li>
                                        <li v-for="ua in thread.access?.users">
                                            <CitizenInfoPopover :user="ua.user" show-avatar-in-name />
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
                <template v-for="msgs in groupedMessages">
                    <UDivider class="text-xs">
                        <GenericTime :value="msgs[0]?.createdAt" :type="'date'" />
                    </UDivider>

                    <div
                        v-for="message in msgs"
                        :key="message.id"
                        class="hover:border-primary-500 hover:dark:border-primary-400 border-l-2 border-white p-0.5 px-2 hover:bg-base-800 dark:border-gray-900"
                    >
                        <div v-if="message.creator" class="flex justify-between text-xs">
                            <CitizenInfoPopover :user="message.creator" show-avatar-in-name />

                            <GenericTime :value="message.createdAt" type="time" />
                        </div>

                        <div class="flex justify-between text-xs">
                            <p
                                :ref="
                                    (el) => {
                                        if (messages.messages.length) {
                                            messageRef = el as Element;
                                        }
                                    }
                                "
                                class="prose dark:prose-invert"
                            >
                                {{ message.message }}
                            </p>

                            <GenericTime :value="message.createdAt" type="time" />
                        </div>
                    </div>
                </template>
            </template>
        </div>

        <UDivider class="my-2" />

        <UForm
            v-if="thread && canAccessThread(thread.access, thread.creator, AccessLevel.MESSAGE)"
            :schema="schema"
            :state="state"
            @submit="onSubmitThrottle"
        >
            <UFormGroup name="message">
                <UTextarea
                    v-model="state.message"
                    name="message"
                    color="gray"
                    required
                    size="xl"
                    :rows="4"
                    :placeholder="$t('components.messenger.reply')"
                    @focusin="focusTablet(true)"
                    @focusout="focusTablet(false)"
                >
                    <UButton
                        type="submit"
                        color="black"
                        :label="$t('components.messenger.send')"
                        icon="i-mdi-paper-airplane"
                        class="absolute bottom-2.5 right-3.5"
                    />
                </UTextarea>
            </UFormGroup>
        </UForm>
    </UDashboardPanelContent>
</template>
