<script setup lang="ts">
import { format, isToday } from 'date-fns';
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import type { Thread } from '~~/gen/ts/resources/messenger/thread';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';
import { messengerStore } from '~/store/messenger';
import CitizenInfoPopover from '../partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '../partials/elements/GenericTime.vue';

const props = withDefaults(
    defineProps<{
        thread: Thread;
        selected?: boolean;
    }>(),
    {
        selected: false,
    },
);

const schema = z.object({
    message: z.string().min(1).max(2048),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    message: '',
});

onBeforeMount(async () => {
    const count = await messengerStore.threads.count();
    const call = getGRPCMessengerClient().getThreadMessages({
        threadId: props.thread.id,
        after: count > 0 ? undefined : toTimestamp(),
    });
    const { response } = await call;

    await messengerStore.messages.bulkPut(response.messages);
});

watchDebounced(
    () => props.thread,
    async () =>
        messengerStore.setThreadUserState({
            threadId: props.thread.id,
            unread: false,
        }),
);

const messages = useDexieLiveQueryWithDeps(
    () => props.thread.id,
    () =>
        messengerStore.messages
            .where('threadId')
            .equals(props.thread.id)
            .toArray()
            .then((messages) => ({ messages, loaded: true })),
    {
        initialValue: { messages: [], loaded: false },
    },
);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await messengerStore
        .postMessage({
            message: {
                id: '0',
                threadId: props.thread.id,
                message: event.data.message,
                data: {},
            },
        })
        .finally(() => useTimeoutFn(() => (canSubmit.value = true), 1000));
}, 1000);
</script>

<template>
    <UDashboardPanelContent>
        <div class="flex justify-between">
            <div class="flex items-center gap-4">
                <div class="min-w-0">
                    <p class="font-semibold text-gray-900 dark:text-white">
                        {{ thread.title }}
                    </p>
                    <p class="mt-2 font-medium text-gray-500 dark:text-gray-400">
                        <UAvatarGroup size="sm" :max="10">
                            <ProfilePictureImg
                                :src="thread.creator?.avatar?.url"
                                :name="`${thread.creator?.firstname} ${thread.creator?.lastname}`"
                            />

                            <ProfilePictureImg
                                v-for="user in thread.access?.users"
                                :src="user.user?.avatar?.url"
                                :name="`${user.user?.firstname} ${user.user?.lastname}`"
                            />
                        </UAvatarGroup>
                    </p>
                </div>
            </div>

            <p class="font-medium text-gray-900 dark:text-white">
                {{
                    isToday(toDate(thread.createdAt))
                        ? $d(toDate(thread.createdAt), 'time')
                        : $d(toDate(thread.createdAt), 'date')
                }}
            </p>
        </div>

        <UDivider class="my-4" />

        <div class="flex-1 flex-col-reverse">
            <template v-if="!messages.loaded">
                <div class="space-y-2">
                    <USkeleton class="h-6 w-full" />
                    <USkeleton class="h-6 w-full" />
                </div>
            </template>
            <template v-else>
                <div v-for="message in messages.messages" :key="message.id">
                    <div class="flex justify-between">
                        <CitizenInfoPopover :user="message.creator" show-avatar-in-name />

                        <GenericTime :value="message.createdAt" :type="isToday(toDate(message.createdAt)) ? 'time' : 'long'" />
                    </div>

                    <p class="text-lg">
                        {{ message.message }}
                    </p>
                </div>
            </template>
        </div>

        <UDivider class="my-4" />

        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UFormGroup name="message">
                <UTextarea
                    v-model="state.message"
                    name="message"
                    color="gray"
                    required
                    size="xl"
                    :rows="4"
                    :placeholder="$t('components.messenger.reply')"
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
