<script setup lang="ts">
import { isToday } from 'date-fns';
import type { Thread } from '~~/gen/ts/resources/mailer/thread';

const props = withDefaults(
    defineProps<{
        modelValue?: Thread;
        threads: Thread[];
        loaded: boolean;
    }>(),
    {
        modelValue: undefined,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: Thread | undefined): void;
}>();

const threadRefs = ref(new Map<number, Element>());

const selectedThread = computed({
    get() {
        return props.modelValue;
    },
    set(value: Thread | undefined) {
        emit('update:modelValue', value);
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
    <div>
        <div v-if="!loaded" class="space-y-2">
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
        </div>

        <template v-else>
            <div v-for="(thread, index) in threads" :key="index" :ref="(el) => threadRefs.set(thread.id, el as Element)">
                <div
                    class="cursor-pointer border-l-2 p-4 text-sm"
                    :class="[
                        !!thread.state?.unread ? 'text-highlighted' : 'text-gray-600 dark:text-gray-300',
                        selectedThread && selectedThread.id === thread.id
                            ? 'border-primary-500 bg-primary-100 dark:border-primary-400 dark:bg-primary-900/25'
                            : 'border-white hover:border-primary-500/25 hover:bg-primary-100/50 dark:border-neutral-900 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10',
                    ]"
                    @click="selectedThread = thread"
                >
                    <div class="flex items-center justify-between gap-1" :class="[thread.state?.unread && 'font-semibold']">
                        <div class="flex items-center gap-3 truncate font-semibold">
                            <span class="truncate">
                                {{ thread.title }}
                            </span>

                            <UChip v-if="!!thread.state?.unread" class="mr-1" />
                        </div>

                        <div
                            v-if="thread.deletedAt"
                            class="flex shrink-0 flex-row items-center justify-center gap-1.5 font-bold"
                        >
                            <UIcon class="size-4 shrink-0" name="i-mdi-delete" />
                            {{ $t('common.deleted') }}
                        </div>
                        <UTooltip v-else class="shrink-0" :text="$d(toDate(thread.updatedAt ?? thread.createdAt), 'long')">
                            {{
                                isToday(toDate(thread.updatedAt ?? thread.createdAt))
                                    ? $d(toDate(thread.updatedAt ?? thread.createdAt), 'time')
                                    : $d(toDate(thread.updatedAt ?? thread.createdAt), 'date')
                            }}
                        </UTooltip>
                    </div>
                    <div class="flex items-center justify-between">
                        <p>{{ thread.creatorEmail?.email }}</p>

                        <div class="inline-flex gap-1">
                            <UIcon v-if="thread.state?.important" class="size-5 text-red-500" name="i-mdi-exclamation-thick" />
                            <UIcon v-if="thread.state?.favorite" class="size-5 text-yellow-500" name="i-mdi-star" />
                        </div>
                    </div>
                </div>

                <USeparator v-if="index < threads.length" />
            </div>

            <slot name="after" />
        </template>
    </div>
</template>
