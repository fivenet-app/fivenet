<script lang="ts" setup>
import { VueDraggable } from 'vue-draggable-plus';
import type { zWorkflowSchema } from './types';

const props = defineProps<{
    modelValue: zWorkflowSchema;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', payload: zWorkflowSchema): void;
}>();

const workflow = useVModel(props, 'modelValue', emit);
</script>

<template>
    <div class="flex flex-col gap-1">
        <!-- Auto Close -->
        <UFormField name="workflow.autoClose" :label="$t('common.auto_close')">
            <USwitch v-model="workflow.autoClose.autoClose" />
        </UFormField>

        <UFormField
            name="workflow.autoClose.autoCloseSettings"
            :description="$t('components.documents.TemplateWorkflowEditor.auto_close.description')"
        >
            <div class="flex items-center gap-1">
                <UFormField name="workflow.autoClose.autoCloseSettings.duration" :label="$t('common.time_ago.day', 2)">
                    <UInput
                        v-model="workflow.autoClose.autoCloseSettings.duration"
                        type="number"
                        :min="1"
                        :max="60"
                        :step="1"
                        :placeholder="$t('common.duration')"
                    >
                        <template #trailing>
                            <span class="text-xs text-neutral-500 dark:text-neutral-400">{{
                                $t('common.time_ago.day', 2)
                            }}</span>
                        </template>
                    </UInput>
                </UFormField>

                <UFormField
                    name="workflow.autoClose.autoCloseSettings.message"
                    :label="$t('common.message')"
                    class="grid flex-1 grid-cols-1 items-center"
                    :ui="{ container: '' }"
                >
                    <UInput
                        v-model="workflow.autoClose.autoCloseSettings.message"
                        type="text"
                        class="w-full flex-1"
                        :placeholder="$t('common.message')"
                    />
                </UFormField>
            </div>
        </UFormField>

        <!-- Reminders -->
        <UFormField name="workflow.reminders.reminder" :label="$t('common.reminder')">
            <USwitch v-model="workflow.reminders.reminder" />
        </UFormField>

        <UFormField
            name="workflow.reminders.reminders"
            :label="$t('common.reminder', 2)"
            :description="$t('components.documents.TemplateWorkflowEditor.reminder.description')"
        >
            <div class="flex flex-col gap-1">
                <VueDraggable v-model="workflow.reminders.reminderSettings.reminders" class="flex flex-col gap-2">
                    <div
                        v-for="(_, idx) in workflow.reminders.reminderSettings.reminders"
                        :key="idx"
                        class="flex items-center gap-1"
                    >
                        <UIcon name="i-mdi-drag-horizontal" class="size-6" />

                        <UFormField
                            :name="`workflow.reminders.reminders.${idx}.duration`"
                            :label="$t('common.time_ago.day', 2)"
                            class="grid grid-cols-1 items-center"
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
                                    <span class="text-xs text-neutral-500 dark:text-neutral-400">{{
                                        $t('common.time_ago.day', 2)
                                    }}</span>
                                </template>
                            </UInput>
                        </UFormField>

                        <UFormField
                            :name="`workflow.reminders.reminders.${idx}.message`"
                            :label="$t('common.message')"
                            class="grid flex-1 grid-cols-1 items-center"
                            :ui="{ container: '' }"
                        >
                            <UInput
                                v-model="workflow.reminders.reminderSettings.reminders[idx]!.message"
                                type="text"
                                class="w-full flex-1"
                                :placeholder="$t('common.message')"
                            />
                        </UFormField>

                        <UFormField label="&nbsp;">
                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                icon="i-mdi-close"
                                @click="workflow.reminders.reminderSettings.reminders.splice(idx, 1)"
                            />
                        </UFormField>
                    </div>
                </VueDraggable>
            </div>

            <UButton
                :ui="{ rounded: 'rounded-full' }"
                icon="i-mdi-plus"
                :disabled="workflow.reminders.reminderSettings.reminders.length >= 3"
                :class="workflow.reminders.reminderSettings.reminders.length ? 'mt-2' : ''"
                @click="
                    workflow.reminders.reminderSettings.reminders.push({
                        duration: 7,
                        message: '',
                    })
                "
            />
        </UFormField>
    </div>
</template>
