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
        <UFormGroup :label="$t('common.citizen', 2)" :ui="{ container: 'flex flex-1 flex-row gap-1' }">
            <UFormGroup name="users.req" :label="$t('common.required')">
                <UToggle v-model="templateSchema.users.req" />
            </UFormGroup>

            <UFormGroup class="flex-1" name="users.min" :label="$t('common.min')">
                <UInput
                    v-model="templateSchema.users.min"
                    type="number"
                    :min="0"
                    :disabled="!templateSchema.users.req"
                    oninput="validity.valid||(value='');"
                    name="users"
                />
            </UFormGroup>

            <UFormGroup class="flex-1" name="users.max" :label="$t('common.max')">
                <UInput
                    v-model="templateSchema.users.max"
                    type="number"
                    :min="0"
                    oninput="validity.valid||(value='');"
                    name="users"
                />
            </UFormGroup>
        </UFormGroup>

        <UFormGroup name="documents.req" :label="$t('common.document', 2)" :ui="{ container: 'flex flex-1 flex-row gap-1' }">
            <UFormGroup :label="$t('common.required')">
                <UToggle v-model="templateSchema.documents.req" />
            </UFormGroup>

            <UFormGroup class="flex-1" name="documents.min" :label="$t('common.min')">
                <UInput
                    v-model="templateSchema.documents.min"
                    type="number"
                    :min="0"
                    :disabled="!templateSchema.documents.req"
                    oninput="validity.valid||(value='');"
                    name="documents"
                />
            </UFormGroup>

            <UFormGroup class="flex-1" name="documents.max" :label="$t('common.max')">
                <UInput
                    v-model="templateSchema.documents.max"
                    type="number"
                    :min="0"
                    oninput="validity.valid||(value='');"
                    name="documents"
                />
            </UFormGroup>
        </UFormGroup>

        <UFormGroup :label="$t('common.vehicle', 2)" :ui="{ container: 'flex flex-1 flex-row gap-1' }">
            <UFormGroup name="vehicles.req" :label="$t('common.required')">
                <UToggle v-model="templateSchema.vehicles.req" />
            </UFormGroup>

            <UFormGroup class="flex-1" name="vehicles.min" :label="$t('common.min')">
                <UInput
                    v-model="templateSchema.vehicles.min"
                    type="number"
                    :min="0"
                    :disabled="!templateSchema.vehicles.req"
                    oninput="validity.valid||(value='');"
                />
            </UFormGroup>

            <UFormGroup class="flex-1" name="vehicles.max" :label="$t('common.max')">
                <UInput
                    v-model="templateSchema.vehicles.max"
                    type="number"
                    :min="0"
                    oninput="validity.valid||(value='');"
                    name="vehicles"
                />
            </UFormGroup>
        </UFormGroup>
    </div>
</template>
