<script setup lang="ts">
import { isSameDay } from 'date-fns';
import { computed, ref, watch } from 'vue';
import DeletedAtBadge from '~/components/partials/DeletedAtBadge.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import type { Thread } from '~~/gen/ts/resources/mailer/threads/thread';

const props = withDefaults(
    defineProps<{
        modelValue?: Thread;
        threads: Thread[];
        loaded: boolean;
        emptyMessage?: string;
    }>(),
    {
        modelValue: undefined,
        emptyMessage: undefined,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: Thread | undefined): void;
}>();

const threadRefs = ref(new Map<number, Element>());
const minuteClock = useMinuteClock();

function isToday(date: Date): boolean {
    return isSameDay(date, minuteClock.value);
}

const routeParams = useRouteParams('thread');

const selectedThread = computed({
    get() {
        return props.modelValue;
    },
    set(value: Thread | undefined) {
        routeParams.value = value ? String(value.id) : '';

        emit('update:modelValue', value);
    },
});

watch(selectedThread, () => {
    if (!selectedThread.value) return;

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

function setThreadRef(threadId: number, el: Element | null): void {
    if (el) {
        threadRefs.value.set(threadId, el);
    } else {
        threadRefs.value.delete(threadId);
    }
}
</script>

<template>
    <div class="flex flex-1 flex-col">
        <div v-if="!loaded" class="space-y-2">
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
            <USkeleton class="h-[73px] w-full" />
        </div>

        <template v-else>
            <DataNoDataBlock
                v-if="threads.length === 0"
                class="m-4 flex flex-1 items-center justify-center rounded-lg border border-default p-4"
                :message="emptyMessage"
                :type="$t('common.mail', 2)"
                icon="i-mdi-email-outline"
                :padded="false"
            />

            <div v-else class="min-h-0 flex-1 divide-y divide-default overflow-y-auto">
                <div v-for="thread in threads" :key="thread.id" :ref="(el) => setThreadRef(thread.id, el as Element | null)">
                    <div
                        class="cursor-pointer border-l-2 p-4 text-sm transition-colors focus-visible:ring-2 focus-visible:ring-primary focus-visible:outline-none focus-visible:ring-inset sm:px-6"
                        role="button"
                        tabindex="0"
                        :aria-current="selectedThread?.id === thread.id ? 'true' : undefined"
                        :class="[
                            !!thread.state?.unread ? 'text-highlighted' : 'text-toned',
                            selectedThread && selectedThread.id === thread.id
                                ? 'border-primary bg-primary/10'
                                : 'border-(--ui-bg) hover:border-primary hover:bg-primary/5',
                        ]"
                        @click="selectedThread = thread"
                        @keydown.enter="selectedThread = thread"
                        @keydown.space.prevent="selectedThread = thread"
                    >
                        <div
                            class="flex min-w-0 items-center justify-between gap-1"
                            :class="[thread.state?.unread && 'font-semibold']"
                        >
                            <div class="flex min-w-0 flex-1 items-center gap-3 truncate font-semibold">
                                <span class="block truncate">
                                    {{ thread.title }}
                                </span>

                                <UChip v-if="!!thread.state?.unread" class="mr-1" />
                            </div>

                            <DeletedAtBadge
                                v-if="thread.deletedAt"
                                class="shrink-0"
                                hide-date
                                icon="i-mdi-delete"
                                :deleted-at="thread.deletedAt"
                            />
                            <UTooltip v-else class="shrink-0" :text="$d(toDate(thread.updatedAt ?? thread.createdAt), 'long')">
                                {{
                                    isToday(toDate(thread.updatedAt ?? thread.createdAt))
                                        ? $d(toDate(thread.updatedAt ?? thread.createdAt), 'time')
                                        : $d(toDate(thread.updatedAt ?? thread.createdAt), 'date')
                                }}
                            </UTooltip>
                        </div>
                        <div class="flex min-w-0 items-center justify-between gap-2 text-xs text-muted">
                            <p class="truncate">{{ thread.creatorEmail?.email }}</p>

                            <div class="inline-flex h-5 min-w-10 items-center justify-end gap-1">
                                <UIcon
                                    v-if="thread.state?.important"
                                    class="size-5 text-red-500"
                                    name="i-mdi-exclamation-thick"
                                />
                                <UIcon v-if="thread.state?.favorite" class="size-5 text-yellow-500" name="i-mdi-star" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <slot name="after" />
        </template>
    </div>
</template>
