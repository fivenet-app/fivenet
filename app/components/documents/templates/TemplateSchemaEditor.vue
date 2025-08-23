<script lang="ts" setup>
import type { ObjectSpecsValue } from './types';

export interface SchemaEditorValue {
    users: ObjectSpecsValue;

    documents: ObjectSpecsValue;

    vehicles: ObjectSpecsValue;
}

const props = defineProps<{
    modelValue: SchemaEditorValue;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', payload: SchemaEditorValue): void;
}>();

const templateSchema = useVModel(props, 'modelValue', emit);
</script>

<template>
    <div class="flex flex-col gap-1">
        <UFormField :label="$t('common.citizen', 2)" :ui="{ container: 'flex flex-1 flex-row gap-1' }">
            <UFormField name="users.req" :label="$t('common.required')">
                <USwitch v-model="templateSchema.users.req" />
            </UFormField>

            <UFormField class="flex-1" name="users.min" :label="$t('common.min')">
                <UInput
                    v-model="templateSchema.users.min"
                    type="number"
                    :min="0"
                    :disabled="!templateSchema.users.req"
                    oninput="validity.valid||(value='');"
                    name="users"
                />
            </UFormField>

            <UFormField class="flex-1" name="users.max" :label="$t('common.max')">
                <UInput
                    v-model="templateSchema.users.max"
                    type="number"
                    :min="0"
                    oninput="validity.valid||(value='');"
                    name="users"
                />
            </UFormField>
        </UFormField>

        <UFormField name="documents.req" :label="$t('common.document', 2)" :ui="{ container: 'flex flex-1 flex-row gap-1' }">
            <UFormField :label="$t('common.required')">
                <USwitch v-model="templateSchema.documents.req" />
            </UFormField>

            <UFormField class="flex-1" name="documents.min" :label="$t('common.min')">
                <UInput
                    v-model="templateSchema.documents.min"
                    type="number"
                    :min="0"
                    :disabled="!templateSchema.documents.req"
                    oninput="validity.valid||(value='');"
                    name="documents"
                />
            </UFormField>

            <UFormField class="flex-1" name="documents.max" :label="$t('common.max')">
                <UInput
                    v-model="templateSchema.documents.max"
                    type="number"
                    :min="0"
                    oninput="validity.valid||(value='');"
                    name="documents"
                />
            </UFormField>
        </UFormField>

        <UFormField :label="$t('common.vehicle', 2)" :ui="{ container: 'flex flex-1 flex-row gap-1' }">
            <UFormField name="vehicles.req" :label="$t('common.required')">
                <USwitch v-model="templateSchema.vehicles.req" />
            </UFormField>

            <UFormField class="flex-1" name="vehicles.min" :label="$t('common.min')">
                <UInput
                    v-model="templateSchema.vehicles.min"
                    type="number"
                    :min="0"
                    :disabled="!templateSchema.vehicles.req"
                    oninput="validity.valid||(value='');"
                />
            </UFormField>

            <UFormField class="flex-1" name="vehicles.max" :label="$t('common.max')">
                <UInput
                    v-model="templateSchema.vehicles.max"
                    type="number"
                    :min="0"
                    oninput="validity.valid||(value='');"
                    name="vehicles"
                />
            </UFormField>
        </UFormField>
    </div>
</template>
