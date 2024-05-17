<script setup lang="ts">
import { format, isToday } from 'date-fns';
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import type { Thread } from '~~/gen/ts/resources/messenger/thread';
import ProfilePictureImg from '~/components/partials/citizens/ProfilePictureImg.vue';

withDefaults(
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

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    //await forgotPassword(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UDashboardPanelContent>
        <div class="flex justify-between">
            <div class="flex items-center gap-4">
                <ProfilePictureImg
                    :src="thread.creator?.avatar?.url"
                    :name="`${thread.creator?.firstname} ${thread.creator?.lastname}`"
                    size="lg"
                />

                <div class="min-w-0">
                    <p class="font-semibold text-gray-900 dark:text-white">
                        {{ thread.creator?.firstname }} {{ thread.creator?.lastname }}
                    </p>
                    <p class="font-medium text-gray-500 dark:text-gray-400">
                        {{ thread.title }}
                    </p>
                </div>
            </div>

            <p class="font-medium text-gray-900 dark:text-white">
                {{
                    isToday(toDate(thread.createdAt))
                        ? format(toDate(thread.createdAt), 'HH:mm')
                        : format(toDate(thread.createdAt), 'dd MMM')
                }}
            </p>
        </div>

        <UDivider class="my-5" />

        <div class="flex-1">
            <p class="text-lg">
                {{ thread.lastMessage }}
            </p>
        </div>

        <UDivider class="my-5" />

        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UTextarea
                name="message"
                color="gray"
                required
                size="xl"
                :rows="5"
                :placeholder="
                    $t('components.messenger.reply_to', { name: `${thread.creator?.firstname} ${thread.creator?.lastname}` })
                "
            >
                <UButton
                    type="submit"
                    color="black"
                    :label="$t('components.messenger.send')"
                    icon="i-mdi-paper-airplane"
                    class="absolute bottom-2.5 right-3.5"
                />
            </UTextarea>
        </UForm>
    </UDashboardPanelContent>
</template>
