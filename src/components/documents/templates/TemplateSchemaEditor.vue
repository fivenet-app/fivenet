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
                <UCheckbox
                    v-model="modelValue.users.req"
                    name="users"
                    class="hover:bg-background hover:text-primary-500 size-8 rounded-md border-0 bg-base-700 transition-colors focus:ring-1 focus:ring-inset focus:ring-base-300"
                />
                <div class="flex flex-row">
                    <span
                        :class="[
                            modelValue.users.req ? 'text-neutral' : 'text-base-400',
                            'inline-flex h-9 items-center rounded-l-md border-0 bg-base-800 px-3 transition-colors sm:text-sm',
                        ]"
                    >
                        {{ $t('common.min') }}
                    </span>
                    <UInput
                        v-model="modelValue.users.min"
                        :disabled="!modelValue.users.req"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="users"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 transition-colors focus:ring-1 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
                <div class="flex flex-row">
                    <span
                        class="inline-flex h-9 items-center rounded-l-md border-0 bg-base-800 px-3 transition-colors sm:text-sm"
                    >
                        {{ $t('common.max') }}
                    </span>
                    <UInput
                        v-model="modelValue.users.max"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="users"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 transition-colors focus:ring-1 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
            </div>
        </div>
        <div class="flex-1">
            <label for="documents" class="block text-sm font-medium leading-6">
                {{ $t('common.document', 2) }}
            </label>
            <div class="relative flex items-center gap-1">
                <UCheckbox
                    v-model="modelValue.documents.req"
                    name="documents"
                    class="hover:bg-background hover:text-primary-500 size-8 rounded-md border-0 bg-base-700 transition-colors focus:ring-1 focus:ring-inset focus:ring-base-300"
                />
                <div class="flex flex-row">
                    <span
                        :class="[
                            modelValue.documents.req ? 'text-neutral' : 'text-base-400',
                            'inline-flex h-9 items-center rounded-l-md border-0 bg-base-800 px-3 transition-colors sm:text-sm',
                        ]"
                    >
                        {{ $t('common.min') }}
                    </span>
                    <UInput
                        v-model="modelValue.documents.min"
                        :disabled="!modelValue.documents.req"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="documents"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 transition-colors focus:ring-1 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
                <div class="flex flex-row">
                    <span
                        class="inline-flex h-9 items-center rounded-l-md border-0 bg-base-800 px-3 transition-colors sm:text-sm"
                    >
                        {{ $t('common.max') }}
                    </span>
                    <UInput
                        v-model="modelValue.documents.max"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="documents"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 transition-colors focus:ring-1 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
            </div>
        </div>
        <div class="flex-1">
            <label for="vehicles" class="block text-sm font-medium leading-6">
                {{ $t('common.vehicle', 2) }}
            </label>
            <div class="relative flex items-center gap-1">
                <UCheckbox
                    v-model="modelValue.vehicles.req"
                    name="vehicles"
                    class="hover:bg-background hover:text-primary-500 size-8 rounded-md border-0 bg-base-700 transition-colors focus:ring-1 focus:ring-inset focus:ring-base-300"
                />
                <div class="flex flex-row">
                    <span
                        :class="[
                            modelValue.vehicles.req ? 'text-neutral' : 'text-base-400',
                            'inline-flex h-9 items-center rounded-l-md border-0 bg-base-800 px-3 transition-colors sm:text-sm',
                        ]"
                    >
                        {{ $t('common.min') }}
                    </span>
                    <UInput
                        v-model="modelValue.vehicles.min"
                        :disabled="!modelValue.vehicles.req"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="vehicles"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 transition-colors focus:ring-1 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
                <div class="flex flex-row">
                    <span
                        class="inline-flex h-9 items-center rounded-l-md border-0 bg-base-800 px-3 transition-colors sm:text-sm"
                    >
                        {{ $t('common.max') }}
                    </span>
                    <UInput
                        v-model="modelValue.vehicles.max"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="vehicles"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 transition-colors focus:ring-1 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
            </div>
        </div>
    </div>
</template>
