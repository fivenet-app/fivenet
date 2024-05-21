<script setup lang="ts">
import { format, isToday } from 'date-fns';
import type { Thread } from '~~/gen/ts/resources/messenger/thread';

const props = withDefaults(
    defineProps<{
        modelValue?: Thread;
        threads: Thread[];
        loaded: boolean;
    }>(),
    {
        threads: () => [],
    },
);

const emit = defineEmits<{
    (e: 'update:model-value', value: Thread | undefined): void;
}>();

const threadRefs = ref<Map<string, Element>>(new Map());

const selectedThread = computed({
    get() {
        return props.modelValue;
    },
    set(value: Thread | undefined) {
        emit('update:model-value', value);
    },
});

watch(selectedThread, () => {
    if (!selectedThread.value) {
        return;
    }

    const ref = threadRefs.value.get(selectedThread.value?.id);
    if (ref) {
        ref.scrollIntoView({ block: 'nearest' });
    }
});

defineShortcuts({
    arrowdown: () => {
        const index = props.threads.findIndex((thread) => thread.id === selectedThread.value?.id);

        if (index === -1) {
            selectedThread.value = props.threads[0];
        } else if (index < props.threads.length - 1) {
            selectedThread.value = props.threads[index + 1];
        }
    },
    arrowup: () => {
        const index = props.threads.findIndex((mail) => mail.id === selectedThread.value?.id);

        if (index === -1) {
            selectedThread.value = props.threads[props.threads.length - 1];
        } else if (index > 0) {
            selectedThread.value = props.threads[index - 1];
        }
    },
});
</script>

<template>
    <UDashboardPanelContent class="p-0">
        <div v-if="!loaded" class="space-y-2">
            <USkeleton :ui="{ rounded: '' }" class="h-[73px] w-full" />
            <USkeleton :ui="{ rounded: '' }" class="h-[73px] w-full" />
            <USkeleton :ui="{ rounded: '' }" class="h-[73px] w-full" />
        </div>
        <template v-else>
            <div
                v-for="(thread, index) in threads"
                :key="index"
                :ref="
                    (el) => {
                        threadRefs.set(thread.id, el as Element);
                    }
                "
            >
                <div
                    class="cursor-pointer border-l-2 p-4 text-sm"
                    :class="[
                        thread.userState?.unread ? 'text-gray-900 dark:text-white' : 'text-gray-600 dark:text-gray-300',
                        selectedThread && selectedThread.id === thread.id
                            ? 'border-primary-500 dark:border-primary-400 bg-primary-100 dark:bg-primary-900/25'
                            : 'hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 border-white dark:border-gray-900',
                    ]"
                    @click="selectedThread = thread"
                >
                    <div class="flex items-center justify-between" :class="[thread.userState?.unread && 'font-semibold']">
                        <div class="flex items-center gap-3">
                            {{ thread.title }}

                            <UChip v-if="thread.userState?.unread" />
                        </div>

                        <span>{{
                            isToday(toDate(thread.createdAt))
                                ? format(toDate(thread.createdAt), 'HH:mm')
                                : format(toDate(thread.createdAt), 'dd MMM')
                        }}</span>
                    </div>
                    <div class="flex items-center justify-between">
                        <p>{{ thread.creator?.firstname }} {{ thread.creator?.lastname }}</p>

                        <div class="inline-flex gap-1">
                            <UIcon
                                v-if="thread.userState?.important"
                                name="i-mdi-exclamation-thick"
                                class="size-5 text-red-500"
                            />
                            <UIcon v-if="thread.userState?.favorite" name="i-mdi-star" class="size-5 text-yellow-500" />
                        </div>
                    </div>
                    <p class="line-clamp-1 text-gray-400 dark:text-gray-500">
                        {{ thread.lastMessage?.message }}
                    </p>
                </div>

                <UDivider />
            </div>
        </template>
    </UDashboardPanelContent>
</template>
