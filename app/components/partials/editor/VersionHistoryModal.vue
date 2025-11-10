<script setup lang="ts">
import type { Content, Version } from '~/types/history';
import VersionDiffModal from './VersionDiffModal.vue';

const props = defineProps<{
    historyType: string;
    currentContent: Content;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'apply', version: Version<unknown>): void;
}>();

const overlay = useOverlay();

const historyStore = useHistoryStore();

const history = historyStore.listHistory<Content>(props.historyType);

const sortedHistory = computed(() => (history.value ?? []).slice().sort((a, b) => b.date.localeCompare(a.date)));

const selectedVersion = ref<Version<Content> | undefined>(undefined);

function date(val: string) {
    return new Date(val).toLocaleString();
}

const versionDiffModal = overlay.create(VersionDiffModal);

function emitApply(version: Version<Content>) {
    selectedVersion.value = version;

    versionDiffModal.open({
        currentContent: props.currentContent.content,
        selectedVersion: version,
        onApply: (version: Version<Content>) => onConfirmDiff(version),
    });
}

function onConfirmDiff(version: Version<Content>) {
    emit('apply', version);
    selectedVersion.value = undefined;
    emit('close', false);
}

watch(
    () => versionDiffModal.isOpen,
    (val) => {
        if (!val) {
            selectedVersion.value = undefined;
        }
    },
);
</script>

<template>
    <UModal :title="$t('common.version_history')" fullscreen>
        <template #body>
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

                            <UFieldGroup>
                                <UButton
                                    size="sm"
                                    color="primary"
                                    variant="soft"
                                    :label="$t('common.compare')"
                                    @click="emitApply(version as Version<Content>)"
                                />
                                <UButton icon="i-mdi-trash" color="error" @click="historyStore.deleteVersion(version.date)" />
                            </UFieldGroup>
                        </div>
                    </li>
                </ul>
            </div>
            <div v-else class="py-4 text-center text-dimmed">{{ $t('common.no_versions') }}</div>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" block color="neutral" @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>
            </UFieldGroup>
        </template>
    </UModal>
</template>
