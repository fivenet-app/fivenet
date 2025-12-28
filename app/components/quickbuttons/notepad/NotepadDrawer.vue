<script lang="ts" setup>
defineEmits<{
    close: [boolean];
}>();
</script>

<template>
    <UDrawer
        :title="$t('components.notepad.title')"
        :overlay="false"
        :close="{ onClick: () => $emit('close', false) }"
        side="bottom"
        handle-only
        :ui="{ title: 'flex gap-2' }"
    >
        <template #title>
            <div class="inline-flex flex-1 items-center gap-1 font-medium">
                <span>{{ $t('components.notepad.title') }}</span>

                <UPopover :content="{ align: 'start' }" arrow :ui="{ content: 'max-w-120' }">
                    <UTooltip :text="$t('components.notepad.tooltip.title')">
                        <UButton icon="i-mdi-warning-circle-outline" color="warning" variant="ghost" />
                    </UTooltip>

                    <template #content>
                        <div class="flex flex-col items-center gap-2 p-4">
                            <UAlert
                                color="warning"
                                variant="subtle"
                                icon="i-mdi-warning-circle-outline"
                                :title="$t('common.experimental_feature.title')"
                                :description="$t('common.experimental_feature.description')"
                            />

                            <UAlert
                                color="info"
                                variant="subtle"
                                icon="i-mdi-information-outline"
                                :title="$t('components.notepad.tooltip.title')"
                                :description="$t('components.notepad.tooltip.description')"
                            />
                        </div>
                    </template>
                </UPopover>
            </div>

            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="flex justify-center overflow-y-hidden">
                <LazyQuickbuttonsNotepad class="min-h-[60lvh] w-full max-w-[80%] min-w-1/2" />
            </div>
        </template>
    </UDrawer>
</template>
