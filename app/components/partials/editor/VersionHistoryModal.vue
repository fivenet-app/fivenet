<script setup lang="ts">
import type { Content, Version } from '~/types/history';
import VersinoDiffModal from './VersionDiffModal.vue';

const props = defineProps<{
    historyType: string;
    currentContent: Content;
}>();

const emit = defineEmits<{
    (e: 'apply', version: Version<unknown>): void;
}>();

const { isOpen } = useModal();

const historyStore = useHistoryStore();

const history = historyStore.listHistory<Content>(props.historyType);

const sortedHistory = computed(() => (history.value ?? []).slice().sort((a, b) => b.date.localeCompare(a.date)));

const showDiffModal = ref(false);
const selectedVersion = ref<Version<Content> | undefined>(undefined);

function date(val: string) {
    return new Date(val).toLocaleString();
}

function emitApply(version: Version<Content>) {
    selectedVersion.value = version;
    showDiffModal.value = true;
}

function onConfirmDiff(version: Version<Content>) {
    emit('apply', version);
    showDiffModal.value = false;
    selectedVersion.value = undefined;
    isOpen.value = false;
}

watch(showDiffModal, (val) => {
    if (!val) {
        selectedVersion.value = undefined;
        showDiffModal.value = false;
    }
});
</script>

<template>
    <UModal fullscreen>
        <VersinoDiffModal
            v-if="showDiffModal && selectedVersion"
            v-model="showDiffModal"
            :current-content="props.currentContent.content"
            :selected-version="selectedVersion"
            @apply="onConfirmDiff"
        />

        <UCard
            v-else
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.version_history') }}
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div v-if="sortedHistory.length">
                <ul class="list">
                    <li v-for="version in sortedHistory" :key="version.date" class="py-2">
                        <div class="flex w-full items-center justify-between">
                            <div>
                                <p class="font-semibold">
                                    {{ version.name || $t('common.untitled') }}
                                </p>
                                <p class="text-sm text-gray-500">{{ date(version.date) }}</p>
                                <p v-if="version.name" class="text-xs italic">{{ version.name }}</p>
                            </div>

                            <UButtonGroup>
                                <UButton
                                    size="sm"
                                    color="primary"
                                    variant="soft"
                                    :label="$t('common.compare')"
                                    @click="emitApply(version as Version<Content>)"
                                />
                                <UButton icon="i-mdi-trash" color="error" @click="historyStore.deleteVersion(version.date)" />
                            </UButtonGroup>
                        </div>
                    </li>
                </ul>
            </div>
            <div v-else class="py-4 text-center text-gray-400">{{ $t('common.no_versions') }}</div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton class="flex-1" block color="black" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
