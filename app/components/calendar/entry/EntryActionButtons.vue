<script setup lang="ts">
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/entries/entries';
import { useCalendarEntryShortcutState } from '~/composables/useCalendarEntryShortcutState';

const props = withDefaults(
    defineProps<{
        entry: CalendarEntry;
        canEdit?: boolean;
        canShare?: boolean;
        canDelete?: boolean;
        showOpen?: boolean;
        mode?: 'popover' | 'header';
    }>(),
    {
        canEdit: false,
        canShare: false,
        canDelete: false,
        showOpen: false,
        mode: 'popover',
    },
);

const emit = defineEmits<{
    (e: 'open', entry: CalendarEntry): void;
    (e: 'edit', entry: CalendarEntry): void;
    (e: 'share', entry: CalendarEntry): void;
    (e: 'delete', entry: CalendarEntry): void;
}>();

const shortcutState = useCalendarEntryShortcutState();
const isShortcutLayerActive = computed(() => {
    if (props.mode === 'header') {
        return shortcutState.isModalOpen.value && !shortcutState.isPopoverOpen.value;
    }

    return shortcutState.isPopoverOpen.value && !shortcutState.isModalOpen.value;
});

defineShortcuts({
    o: () => {
        if (!isShortcutLayerActive.value || !props.showOpen) return;
        emit('open', props.entry);
    },
    e: () => {
        if (!isShortcutLayerActive.value || !props.canEdit) return;
        emit('edit', props.entry);
    },
    s: () => {
        if (!isShortcutLayerActive.value || !props.canShare) return;
        emit('share', props.entry);
    },
    d: () => {
        if (!isShortcutLayerActive.value || !props.canDelete) return;
        emit('delete', props.entry);
    },
});
</script>

<template>
    <div v-if="mode === 'header'" class="flex items-center gap-1">
        <UTooltip v-if="canEdit" :text="$t('common.edit')" :kbds="['E']">
            <UButton color="neutral" variant="ghost" icon="i-mdi-pencil" @click="emit('edit', entry)" />
        </UTooltip>

        <UTooltip v-if="canShare" :text="$t('common.share')" :kbds="['S']">
            <UButton color="neutral" variant="ghost" icon="i-mdi-share-variant-outline" @click="emit('share', entry)" />
        </UTooltip>

        <UTooltip v-if="canDelete" :text="$t('common.delete')" :kbds="['D']">
            <UButton color="error" variant="ghost" icon="i-mdi-delete-outline" @click="emit('delete', entry)" />
        </UTooltip>
    </div>

    <div v-else class="flex flex-wrap items-center gap-2">
        <UTooltip v-if="showOpen" :text="$t('common.open')" :kbds="['O']">
            <UButton color="neutral" size="sm" icon="i-mdi-eye" :label="$t('common.open')" @click="emit('open', entry)" />
        </UTooltip>

        <UTooltip v-if="canEdit" :text="$t('common.edit')" :kbds="['E']">
            <UButton
                color="neutral"
                size="sm"
                variant="outline"
                icon="i-mdi-pencil"
                :label="$t('common.edit')"
                @click="emit('edit', entry)"
            />
        </UTooltip>

        <div v-if="canShare || canDelete" class="ms-auto flex items-center gap-2">
            <UTooltip v-if="canShare" :text="$t('common.share')" :kbds="['S']">
                <UButton
                    color="neutral"
                    size="sm"
                    variant="ghost"
                    icon="i-mdi-share"
                    :label="$t('common.share')"
                    @click="emit('share', entry)"
                />
            </UTooltip>

            <UTooltip v-if="canDelete" :text="$t('common.delete')" :kbds="['D']">
                <UButton
                    color="error"
                    size="sm"
                    variant="ghost"
                    icon="i-mdi-delete"
                    :label="$t('common.delete')"
                    @click="emit('delete', entry)"
                />
            </UTooltip>
        </div>
    </div>
</template>
