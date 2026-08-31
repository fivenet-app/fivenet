<script lang="ts" setup>
import type { TemplateRequirements } from '~~/gen/ts/resources/documents/templates/templates';
import { CLIPBOARD_MAX_ITEMS } from '~/stores/clipboard';

defineProps<{
    disabled?: boolean;
}>();

const templateSchema = defineModel<TemplateRequirements>({ required: true });
</script>

<template>
    <div class="flex flex-col gap-1 divide-y divide-default">
        <UFormField
            v-if="templateSchema.users"
            class="pb-2"
            :label="$t('common.citizen', 2)"
            :ui="{ label: 'font-bold', container: 'flex flex-1 flex-row gap-1' }"
        >
            <UFormField class="flex-1" name="requirements.users.required" :label="$t('common.required')">
                <USwitch v-model="templateSchema.users.required" :disabled="disabled" />
            </UFormField>

            <UFormField class="flex-1" name="requirements.users.min" :label="$t('common.min')">
                <UInputNumber
                    v-model="templateSchema.users.min"
                    :min="0"
                    :max="CLIPBOARD_MAX_ITEMS"
                    :disabled="disabled || !templateSchema.users.required"
                />
            </UFormField>

            <UFormField class="flex-1" name="requirements.users.max" :label="$t('common.max')">
                <UInputNumber v-model="templateSchema.users.max" :min="0" :max="CLIPBOARD_MAX_ITEMS" :disabled="disabled" />
            </UFormField>
        </UFormField>

        <UFormField
            v-if="templateSchema.documents"
            class="pb-2"
            :label="$t('common.document', 2)"
            :ui="{ label: 'font-bold', container: 'flex flex-1 flex-row gap-1' }"
        >
            <UFormField class="flex-1" name="requirements.documents.required" :label="$t('common.required')">
                <USwitch v-model="templateSchema.documents.required" :disabled="disabled" />
            </UFormField>

            <UFormField class="flex-1" name="requirements.documents.min" :label="$t('common.min')">
                <UInputNumber
                    v-model="templateSchema.documents.min"
                    :min="0"
                    :max="CLIPBOARD_MAX_ITEMS"
                    :disabled="disabled || !templateSchema.documents.required"
                />
            </UFormField>

            <UFormField class="flex-1" name="requirements.documents.max" :label="$t('common.max')">
                <UInputNumber v-model="templateSchema.documents.max" :min="0" :max="CLIPBOARD_MAX_ITEMS" :disabled="disabled" />
            </UFormField>
        </UFormField>

        <UFormField
            v-if="templateSchema.vehicles"
            class="pb-2"
            :label="$t('common.vehicle', 2)"
            :ui="{ label: 'font-bold', container: 'flex flex-1 flex-row gap-1 justify-between' }"
        >
            <UFormField class="flex-1" name="requirements.vehicles.required" :label="$t('common.required')">
                <USwitch v-model="templateSchema.vehicles.required" :disabled="disabled" />
            </UFormField>

            <UFormField class="flex-1" name="requirements.vehicles.min" :label="$t('common.min')">
                <UInputNumber
                    v-model="templateSchema.vehicles.min"
                    :min="0"
                    :max="CLIPBOARD_MAX_ITEMS"
                    :disabled="disabled || !templateSchema.vehicles.required"
                />
            </UFormField>

            <UFormField class="flex-1" name="requirements.vehicles.max" :label="$t('common.max')">
                <UInputNumber v-model="templateSchema.vehicles.max" :min="0" :max="CLIPBOARD_MAX_ITEMS" :disabled="disabled" />
            </UFormField>
        </UFormField>
    </div>
</template>
