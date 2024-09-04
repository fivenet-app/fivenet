<script lang="ts" setup>
import type { ObjectSpecsValue } from './types';

export interface SchemaEditorValue {
    users: ObjectSpecsValue;

    documents: ObjectSpecsValue;

    vehicles: ObjectSpecsValue;
}

defineProps<{
    modelValue: SchemaEditorValue;
}>();

defineEmits<{
    (e: 'update:modelValue', payload: SchemaEditorValue): void;
}>();
</script>

<template>
    <div class="flex flex-row flex-wrap gap-3">
        <div class="flex-1">
            <label for="users" class="block text-sm font-medium leading-6">
                {{ $t('common.citizen', 2) }}
            </label>
            <div class="relative flex items-center gap-1">
                <UCheckbox v-model="modelValue.users.req" name="users" />

                <div class="flex flex-row">
                    <span
                        :class="[
                            modelValue.users.req ? 'dark:text-neutral text-base-100' : 'text-base-900 dark:text-base-400',
                            'inline-flex h-8 items-center px-3 transition-colors sm:text-sm',
                        ]"
                    >
                        {{ $t('common.min') }}
                    </span>

                    <UInput
                        v-model="modelValue.users.min"
                        :disabled="!modelValue.users.req"
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
                        v-model="modelValue.users.max"
                        type="number"
                        :min="0"
                        oninput="validity.valid||(value='');"
                        name="users"
                    />
                </div>
            </div>
        </div>
        <div class="flex-1">
            <label for="documents" class="block text-sm font-medium leading-6">
                {{ $t('common.document', 2) }}
            </label>
            <div class="relative flex items-center gap-1">
                <UCheckbox v-model="modelValue.documents.req" name="documents" />

                <div class="flex flex-row">
                    <span
                        :class="[
                            modelValue.documents.req ? 'dark:text-neutral text-base-100' : 'text-base-900 dark:text-base-400',
                            'inline-flex h-8 items-center px-3 transition-colors sm:text-sm',
                        ]"
                    >
                        {{ $t('common.min') }}
                    </span>

                    <UInput
                        v-model="modelValue.documents.min"
                        :disabled="!modelValue.documents.req"
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
                        v-model="modelValue.documents.max"
                        type="number"
                        :min="0"
                        oninput="validity.valid||(value='');"
                        name="documents"
                    />
                </div>
            </div>
        </div>
        <div class="flex-1">
            <label for="vehicles" class="block text-sm font-medium leading-6">
                {{ $t('common.vehicle', 2) }}
            </label>
            <div class="relative flex items-center gap-1">
                <UCheckbox v-model="modelValue.vehicles.req" name="vehicles" />

                <div class="flex flex-row">
                    <span
                        :class="[
                            modelValue.vehicles.req ? 'dark:text-neutral text-base-100' : 'text-base-900 dark:text-base-400',
                            'inline-flex h-8 items-center px-3 transition-colors sm:text-sm',
                        ]"
                    >
                        {{ $t('common.min') }}
                    </span>

                    <UInput
                        v-model="modelValue.vehicles.min"
                        :disabled="!modelValue.vehicles.req"
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
                        v-model="modelValue.vehicles.max"
                        type="number"
                        :min="0"
                        oninput="validity.valid||(value='');"
                        name="vehicles"
                    />
                </div>
            </div>
        </div>
    </div>
</template>
