<script lang="ts" setup>
import type { TemplateRequirements } from '~~/gen/ts/resources/documents/templates';

const props = defineProps<{
    modelValue: TemplateRequirements;
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', payload: TemplateRequirements): void;
}>();

const templateSchema = useVModel(props, 'modelValue', emit);

onBeforeMount(() => {
    if (!templateSchema.value.users) {
        templateSchema.value.users = { required: false, min: 0, max: 0 };
    }

    if (!templateSchema.value.documents) {
        templateSchema.value.documents = { required: false, min: 0, max: 0 };
    }

    if (!templateSchema.value.vehicles) {
        templateSchema.value.vehicles = { required: false, min: 0, max: 0 };
    }
});
</script>

<template>
    <div class="flex flex-col gap-1 divide-y divide-default">
        <UFormField
            v-if="templateSchema.users"
            :label="$t('common.citizen', 2)"
            class="pb-2"
            :ui="{ label: 'font-bold', container: 'flex flex-1 flex-row gap-1' }"
        >
            <UFormField class="flex-1" name="users.required" :label="$t('common.required')">
                <USwitch v-model="templateSchema.users.required" :disabled="disabled" />
            </UFormField>

            <UFormField class="flex-1" name="users.min" :label="$t('common.min')">
                <UInputNumber
                    v-model="templateSchema.users.min"
                    :min="0"
                    :max="100"
                    :disabled="disabled || !templateSchema.users.required"
                />
            </UFormField>

            <UFormField class="flex-1" name="users.max" :label="$t('common.max')">
                <UInputNumber v-model="templateSchema.users.max" :min="0" :max="100" :disabled="disabled" />
            </UFormField>
        </UFormField>

        <UFormField
            v-if="templateSchema.documents"
            :label="$t('common.document', 2)"
            class="pb-2"
            :ui="{ label: 'font-bold', container: 'flex flex-1 flex-row gap-1' }"
        >
            <UFormField class="flex-1" name="documents.required" :label="$t('common.required')">
                <USwitch v-model="templateSchema.documents.required" :disabled="disabled" />
            </UFormField>

            <UFormField class="flex-1" name="documents.min" :label="$t('common.min')">
                <UInputNumber
                    v-model="templateSchema.documents.min"
                    :min="0"
                    :max="100"
                    :disabled="disabled || !templateSchema.documents.required"
                />
            </UFormField>

            <UFormField class="flex-1" name="documents.max" :label="$t('common.max')">
                <UInputNumber v-model="templateSchema.documents.max" :min="0" :max="100" :disabled="disabled" />
            </UFormField>
        </UFormField>

        <UFormField
            v-if="templateSchema.vehicles"
            :label="$t('common.vehicle', 2)"
            class="pb-2"
            :ui="{ label: 'font-bold', container: 'flex flex-1 flex-row gap-1 justify-between' }"
        >
            <UFormField class="flex-1" name="vehicles.required" :label="$t('common.required')">
                <USwitch v-model="templateSchema.vehicles.required" :disabled="disabled" />
            </UFormField>

            <UFormField class="flex-1" name="vehicles.min" :label="$t('common.min')">
                <UInputNumber
                    v-model="templateSchema.vehicles.min"
                    :min="0"
                    :max="100"
                    :disabled="disabled || !templateSchema.vehicles.required"
                />
            </UFormField>

            <UFormField class="flex-1" name="vehicles.max" :label="$t('common.max')">
                <UInputNumber v-model="templateSchema.vehicles.max" :min="0" :max="100" :disabled="disabled" />
            </UFormField>
        </UFormField>
    </div>
</template>
