<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import type { zWorkflowSchema } from '../types';

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
    <UPageCard :title="`${$t('common.workflow')}: ${$t('common.auto_close')}`">
        <!-- Auto Close -->
        <UFormField name="workflow.autoClose.autoClose" :label="$t('common.enabled')">
            <USwitch v-model="workflow.autoClose.autoClose" />
        </UFormField>

        <UFormField
            name="workflow.autoClose.autoCloseSettings"
            :description="$t('components.documents.template_workflow_editor.auto_close.description')"
        >
            <div class="flex items-center gap-1">
                <UFormField name="workflow.autoClose.autoCloseSettings.duration" :label="$t('common.time_ago.day', 2)">
                    <UInputNumber
                        v-model="workflow.autoClose.autoCloseSettings.duration"
                        :min="1"
                        :max="60"
                        :step="1"
                        :placeholder="$t('common.duration')"
                    />
                </UFormField>

                <UFormField
                    class="grid flex-1 grid-cols-1 items-center"
                    name="workflow.autoClose.autoCloseSettings.message"
                    :label="$t('common.message')"
                >
                    <UInput
                        v-model="workflow.autoClose.autoCloseSettings.message"
                        class="w-full flex-1"
                        type="text"
                        :placeholder="$t('common.message')"
                    />
                </UFormField>
            </div>
        </UFormField>
    </UPageCard>

    <UPageCard :title="`${$t('common.workflow')}: ${$t('common.reminder', 2)}`">
        <!-- Reminders -->
        <UFormField name="workflow.reminders.reminder" :label="$t('common.enabled')">
            <USwitch v-model="workflow.reminders.reminder" />
        </UFormField>

        <UFormField
            name="workflow.reminders.reminders"
            :label="$t('common.reminder', 2)"
            :description="$t('components.documents.template_workflow_editor.reminder.description')"
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

                            <UFieldGroup>
                                <UButton size="xs" variant="link" icon="i-mdi-arrow-up" @click="moveUp(idx)" />
                                <UButton size="xs" variant="link" icon="i-mdi-arrow-down" @click="moveDown(idx)" />
                            </UFieldGroup>
                        </div>

                        <UFormField
                            class="grid grid-cols-1 items-center"
                            :name="`workflow.reminders.reminders.${idx}.duration`"
                            :label="$t('common.time_ago.day', 2)"
                        >
                            <UInputNumber
                                v-model="workflow.reminders.reminderSettings.reminders[idx]!.duration"
                                :min="1"
                                :max="60"
                                :step="1"
                                :placeholder="$t('common.duration')"
                            />
                        </UFormField>

                        <UFormField
                            class="grid flex-1 grid-cols-1 items-center"
                            :name="`workflow.reminders.reminders.${idx}.message`"
                            :label="$t('common.message')"
                        >
                            <UInput
                                v-model="workflow.reminders.reminderSettings.reminders[idx]!.message"
                                class="w-full flex-1"
                                type="text"
                                :placeholder="$t('common.message')"
                            />
                        </UFormField>

                        <UFormField label="&nbsp;">
                            <UButton icon="i-mdi-close" @click="workflow.reminders.reminderSettings.reminders.splice(idx, 1)" />
                        </UFormField>
                    </div>
                </VueDraggable>
            </div>

            <UButton
                :class="workflow.reminders.reminderSettings.reminders.length ? 'mt-2' : ''"
                icon="i-mdi-plus"
                :disabled="workflow.reminders.reminderSettings.reminders.length >= 3"
                @click="
                    workflow.reminders.reminderSettings.reminders.push({
                        duration: 7,
                        message: '',
                    })
                "
            />
        </UFormField>

        <UFormField
            name="workflow.reminders.reminderSettings.maxReminderCount"
            :label="$t('components.documents.template_workflow_editor.max_reminder_count.title')"
            :description="$t('components.documents.template_workflow_editor.max_reminder_count.description')"
        >
            <UInputNumber v-model="workflow.reminders.reminderSettings.maxReminderCount" :min="1" :max="10" :step="1" />
        </UFormField>
    </UPageCard>
</template>
