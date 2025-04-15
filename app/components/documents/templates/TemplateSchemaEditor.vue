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

            <UFormField name="users.min" :label="$t('common.min')" class="flex-1">
                <UInput
                    v-model="templateSchema.users.min"
                    type="number"
                    :min="0"
                    :disabled="!templateSchema.users.req"
                    oninput="validity.valid||(value='');"
                    name="users"
                />
            </UFormField>

            <UFormField name="users.max" :label="$t('common.max')" class="flex-1">
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

            <UFormField name="documents.min" :label="$t('common.min')" class="flex-1">
                <UInput
                    v-model="templateSchema.documents.min"
                    type="number"
                    :min="0"
                    :disabled="!templateSchema.documents.req"
                    oninput="validity.valid||(value='');"
                    name="documents"
                />
            </UFormField>

            <UFormField name="documents.max" :label="$t('common.max')" class="flex-1">
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

            <UFormField name="vehicles.min" :label="$t('common.min')" class="flex-1">
                <UInput
                    v-model="templateSchema.vehicles.min"
                    type="number"
                    :min="0"
                    :disabled="!templateSchema.vehicles.req"
                    oninput="validity.valid||(value='');"
                />
            </UFormField>

            <UFormField name="vehicles.max" :label="$t('common.max')" class="flex-1">
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
