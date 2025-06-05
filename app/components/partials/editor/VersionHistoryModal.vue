<script setup lang="ts">
import type { Content, Version } from '~/types/history';
import VersinoDiffModal from './VersionDiffModal.vue';

const props = defineProps<{
    history: Version<Content>[];
    title?: string;
    currentContent: Content;
}>();

const emit = defineEmits<{
    (e: 'apply', version: Version<unknown>): void;
}>();

const { isOpen } = useModal();

const sortedHistory = computed(() => [...props.history].sort((a, b) => b.id.localeCompare(a.id)));

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
            :ui="{
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
                base: 'flex flex-1 flex-col',
                body: { base: 'flex flex-1 flex-col' },
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ title || $t('common.version_history') }}
                        <span v-if="history && history.length" class="text-xs text-gray-400">({{ history.length }})</span>
                    </h3>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div v-if="history.length">
                <ul class="list">
                    <li v-for="version in sortedHistory" :key="version.id" class="py-2">
                        <div class="flex w-full items-center justify-between">
                            <div>
                                <p class="font-semibold">
                                    {{ version.name || $t('common.untitled') }}
                                </p>
                                <p class="text-sm text-gray-500">{{ date(version.id) }}</p>
                                <p v-if="version.name" class="text-xs italic">{{ version.name }}</p>
                            </div>

                            <UButton size="sm" color="primary" variant="soft" @click="emitApply(version)">{{
                                $t('common.compare')
                            }}</UButton>
                        </div>
                    </li>
                </ul>
            </div>
            <div v-else class="py-4 text-center text-gray-400">No versions available.</div>

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
