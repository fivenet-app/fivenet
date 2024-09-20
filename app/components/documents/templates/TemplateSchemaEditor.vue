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

const emits = defineEmits<{
    (e: 'update:modelValue', payload: SchemaEditorValue): void;
}>();

const templateSchema = useVModel(props, 'modelValue', emits);
</script>

<template>
    <div class="flex flex-row flex-wrap gap-3">
        <div class="flex-1">
            <UFormGroup :label="$t('common.citizen', 2)">
                <div class="flex items-center gap-1">
                    <UCheckbox v-model="templateSchema.users.req" name="users" />

                    <div class="flex flex-row">
                        <span
                            :class="[
                                templateSchema.users.req
                                    ? 'dark:text-neutral text-base-100'
                                    : 'text-base-900 dark:text-base-400',
                                'inline-flex h-8 items-center px-3 transition-colors sm:text-sm',
                            ]"
                        >
                            {{ $t('common.min') }}
                        </span>

                        <UInput
                            v-model="templateSchema.users.min"
                            :disabled="!templateSchema.users.req"
                            type="number"
                            :min="0"
                            oninput="validity.valid||(value='');"
                            name="users"
                        />
                    </div>

                    <div class="flex flex-row">
                        <span class="inline-flex h-8 items-center rounded-l-md px-3 transition-colors sm:text-sm">
                            {{ $t('common.max') }}
                        </span>

                        <UInput
                            v-model="templateSchema.users.max"
                            type="number"
                            :min="0"
                            oninput="validity.valid||(value='');"
                            name="users"
                        />
                    </div>
                </div>
            </UFormGroup>
        </div>

        <div class="flex-1">
            <UFormGroup :label="$t('common.document', 2)">
                <div class="flex items-center gap-1">
                    <UCheckbox v-model="templateSchema.documents.req" name="documents" />

                    <div class="flex flex-row">
                        <span
                            :class="[
                                templateSchema.documents.req
                                    ? 'dark:text-neutral text-base-100'
                                    : 'text-base-900 dark:text-base-400',
                                'inline-flex h-8 items-center px-3 transition-colors sm:text-sm',
                            ]"
                        >
                            {{ $t('common.min') }}
                        </span>

                        <UInput
                            v-model="templateSchema.documents.min"
                            :disabled="!templateSchema.documents.req"
                            type="number"
                            :min="0"
                            oninput="validity.valid||(value='');"
                            name="documents"
                        />
                    </div>

                    <div class="flex flex-row">
                        <span class="inline-flex h-8 items-center rounded-l-md px-3 transition-colors sm:text-sm">
                            {{ $t('common.max') }}
                        </span>

                        <UInput
                            v-model="templateSchema.documents.max"
                            type="number"
                            :min="0"
                            oninput="validity.valid||(value='');"
                            name="documents"
                        />
                    </div>
                </div>
            </UFormGroup>
        </div>

        <div class="flex-1">
            <UFormGroup :label="$t('common.vehicle', 2)">
                <div class="flex items-center gap-1">
                    <UCheckbox v-model="templateSchema.vehicles.req" name="vehicles" />

                    <div class="flex flex-row">
                        <span
                            :class="[
                                templateSchema.vehicles.req
                                    ? 'dark:text-neutral text-base-100'
                                    : 'text-base-900 dark:text-base-400',
                                'inline-flex h-8 items-center px-3 transition-colors sm:text-sm',
                            ]"
                        >
                            {{ $t('common.min') }}
                        </span>

                        <UInput
                            v-model="templateSchema.vehicles.min"
                            :disabled="!templateSchema.vehicles.req"
                            type="number"
                            :min="0"
                            oninput="validity.valid||(value='');"
                            name="vehicles"
                        />
                    </div>
                    <div class="flex flex-row">
                        <span class="inline-flex h-8 items-center rounded-l-md px-3 transition-colors sm:text-sm">
                            {{ $t('common.max') }}
                        </span>

                        <UInput
                            v-model="templateSchema.vehicles.max"
                            type="number"
                            :min="0"
                            oninput="validity.valid||(value='');"
                            name="vehicles"
                        />
                    </div>
                </div>
            </UFormGroup>
        </div>
    </div>
</template>
