<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import type { zWorkflowSchema } from './types';

const props = defineProps<{
    modelValue: zWorkflowSchema;
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', payload: zWorkflowSchema): void;
}>();

const workflow = useVModel(props, 'modelValue', emit);

const { moveUp, moveDown } = useListReorder(toRef(() => workflow.value.reminders.reminderSettings.reminders));
</script>

<template>
    <div class="flex flex-col gap-1">
        <!-- Auto Close -->
        <UFormGroup name="workflow.autoClose" :label="$t('common.auto_close')">
            <UToggle v-model="workflow.autoClose.autoClose" />
        </UFormGroup>

        <UFormGroup
            name="workflow.autoClose.autoCloseSettings"
            :description="$t('components.documents.TemplateWorkflowEditor.auto_close.description')"
        >
            <div class="flex items-center gap-1">
                <UFormGroup name="workflow.autoClose.autoCloseSettings.duration" :label="$t('common.time_ago.day', 2)">
                    <UInput
                        v-model="workflow.autoClose.autoCloseSettings.duration"
                        type="number"
                        :min="1"
                        :max="60"
                        :step="1"
                        :placeholder="$t('common.duration')"
                    >
                        <template #trailing>
                            <span class="text-xs text-gray-500 dark:text-gray-400">{{ $t('common.time_ago.day', 2) }}</span>
                        </template>
                    </UInput>
                </UFormGroup>

                <UFormGroup
                    class="grid flex-1 grid-cols-1 items-center"
                    name="workflow.autoClose.autoCloseSettings.message"
                    :label="$t('common.message')"
                    :ui="{ container: '' }"
                >
                    <UInput
                        v-model="workflow.autoClose.autoCloseSettings.message"
                        class="w-full flex-1"
                        type="text"
                        :placeholder="$t('common.message')"
                    />
                </UFormGroup>
            </div>
        </UFormGroup>

        <!-- Reminders -->
        <UFormGroup name="workflow.reminders.reminder" :label="$t('common.reminder')">
            <UToggle v-model="workflow.reminders.reminder" />
        </UFormGroup>

        <UFormGroup
            name="workflow.reminders.reminders"
            :label="$t('common.reminder', 2)"
            :description="$t('components.documents.TemplateWorkflowEditor.reminder.description')"
        >
            <div class="flex flex-col gap-1">
                <VueDraggable
                    v-model="workflow.reminders.reminderSettings.reminders"
                    :disabled="disabled"
                    class="flex flex-col gap-2"
                    handle=".handle"
                >
                    <div
                        v-for="(_, idx) in workflow.reminders.reminderSettings.reminders"
                        :key="idx"
                        class="flex items-center gap-1"
                    >
                        <div class="inline-flex items-center gap-1">
                            <UTooltip :text="$t('common.draggable')">
                                <UIcon class="handle size-6 cursor-move" name="i-mdi-drag-horizontal" />
                            </UTooltip>

                            <UButtonGroup>
                                <UButton size="xs" variant="link" :padded="false" icon="i-mdi-arrow-up" @click="moveUp(idx)" />
                                <UButton
                                    size="xs"
                                    variant="link"
                                    :padded="false"
                                    icon="i-mdi-arrow-down"
                                    @click="moveDown(idx)"
                                />
                            </UButtonGroup>
                        </div>

                        <UFormGroup
                            class="grid grid-cols-1 items-center"
                            :name="`workflow.reminders.reminders.${idx}.duration`"
                            :label="$t('common.time_ago.day', 2)"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                v-model="workflow.reminders.reminderSettings.reminders[idx]!.duration"
                                type="number"
                                :min="1"
                                :max="60"
                                :step="1"
                                :placeholder="$t('common.duration')"
                            >
                                <template #trailing>
                                    <span class="text-xs text-gray-500 dark:text-gray-400">{{
                                        $t('common.time_ago.day', 2)
                                    }}</span>
                                </template>
                            </UInput>
                        </UFormGroup>

                        <UFormGroup
                            class="grid flex-1 grid-cols-1 items-center"
                            :name="`workflow.reminders.reminders.${idx}.message`"
                            :label="$t('common.message')"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                v-model="workflow.reminders.reminderSettings.reminders[idx]!.message"
                                class="w-full flex-1"
                                type="text"
                                :placeholder="$t('common.message')"
                            />
                        </UFormGroup>

                        <UFormGroup label="&nbsp;">
                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                icon="i-mdi-close"
                                @click="workflow.reminders.reminderSettings.reminders.splice(idx, 1)"
                            />
                        </UFormGroup>
                    </div>
                </VueDraggable>
            </div>

            <UButton
                :class="workflow.reminders.reminderSettings.reminders.length ? 'mt-2' : ''"
                :ui="{ rounded: 'rounded-full' }"
                icon="i-mdi-plus"
                :disabled="workflow.reminders.reminderSettings.reminders.length >= 3"
                @click="
                    workflow.reminders.reminderSettings.reminders.push({
                        duration: 7,
                        message: '',
                    })
                "
            />
        </UFormGroup>
    </div>
</template>
