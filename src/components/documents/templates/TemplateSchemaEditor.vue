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
            <label for="users" class="block text-sm font-medium leading-6 text-neutral">
                {{ $t('common.citizen', 2) }}
            </label>
            <div class="relative flex items-center gap-1">
                <input
                    v-model="modelValue.users.req"
                    type="checkbox"
                    name="users"
                    data-te-toggle="tooltip"
                    class="size-8 rounded-md border-0 bg-base-700 text-primary-600 transition-colors hover:bg-base-600 hover:text-primary-500 focus:ring-2 focus:ring-inset focus:ring-base-300"
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
                    <input
                        v-model="modelValue.users.min"
                        :disabled="!modelValue.users.req"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="users"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 text-neutral transition-colors placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
                <div class="flex flex-row">
                    <span
                        class="inline-flex h-9 items-center rounded-l-md border-0 bg-base-800 px-3 text-neutral transition-colors sm:text-sm"
                    >
                        {{ $t('common.max') }}
                    </span>
                    <input
                        v-model="modelValue.users.max"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="users"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 text-neutral transition-colors placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
            </div>
        </div>
        <div class="flex-1">
            <label for="documents" class="block text-sm font-medium leading-6 text-neutral">
                {{ $t('common.document', 2) }}
            </label>
            <div class="relative flex items-center gap-1">
                <input
                    v-model="modelValue.documents.req"
                    type="checkbox"
                    name="documents"
                    data-te-toggle="tooltip"
                    class="size-8 rounded-md border-0 bg-base-700 text-primary-600 transition-colors hover:bg-base-600 hover:text-primary-500 focus:ring-2 focus:ring-inset focus:ring-base-300"
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
                    <input
                        v-model="modelValue.documents.min"
                        :disabled="!modelValue.documents.req"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="documents"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 text-neutral transition-colors placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
                <div class="flex flex-row">
                    <span
                        class="inline-flex h-9 items-center rounded-l-md border-0 bg-base-800 px-3 text-neutral transition-colors sm:text-sm"
                    >
                        {{ $t('common.max') }}
                    </span>
                    <input
                        v-model="modelValue.documents.max"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="documents"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 text-neutral transition-colors placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
            </div>
        </div>
        <div class="flex-1">
            <label for="vehicles" class="block text-sm font-medium leading-6 text-neutral">
                {{ $t('common.vehicle', 2) }}
            </label>
            <div class="relative flex items-center gap-1">
                <input
                    v-model="modelValue.vehicles.req"
                    type="checkbox"
                    name="vehicles"
                    data-te-toggle="tooltip"
                    class="size-8 rounded-md border-0 bg-base-700 text-primary-600 transition-colors hover:bg-base-600 hover:text-primary-500 focus:ring-2 focus:ring-inset focus:ring-base-300"
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
                    <input
                        v-model="modelValue.vehicles.min"
                        :disabled="!modelValue.vehicles.req"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="vehicles"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 text-neutral transition-colors placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
                <div class="flex flex-row">
                    <span
                        class="inline-flex h-9 items-center rounded-l-md border-0 bg-base-800 px-3 text-neutral transition-colors sm:text-sm"
                    >
                        {{ $t('common.max') }}
                    </span>
                    <input
                        v-model="modelValue.vehicles.max"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="vehicles"
                        class="block w-full rounded-r-md border-0 bg-base-700 py-1.5 text-neutral transition-colors placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 disabled:text-base-400 sm:text-sm sm:leading-6"
                        @focusin="focusTablet(true)"
                        @focusout="focusTablet(false)"
                    />
                </div>
            </div>
        </div>
    </div>
</template>
