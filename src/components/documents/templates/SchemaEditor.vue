<script lang="ts" setup>
export interface ObjectSpecsValue {
    req: boolean;
    min: number;
    max: number;
}

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
    <div class="flex flex-row gap-3 flex-wrap">
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
                    class="h-8 w-8 rounded-md transition-colors bg-base-700 text-primary-600 hover:text-primary-500 hover:bg-base-600 focus:ring-2 focus:ring-inset focus:ring-base-300 border-0"
                />
                <div class="flex flex-row">
                    <span
                        :class="[
                            modelValue.users.req ? 'text-neutral' : 'text-base-400',
                            'transition-colors inline-flex h-9 items-center rounded-l-md border-0 px-3 bg-base-800 sm:text-sm',
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
                        class="block w-full rounded-r-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6 disabled:text-base-400 transition-colors"
                    />
                </div>
                <div class="flex flex-row">
                    <span
                        class="text-neutral transition-colors inline-flex h-9 items-center rounded-l-md border-0 px-3 bg-base-800 sm:text-sm"
                    >
                        {{ $t('common.max') }}
                    </span>
                    <input
                        v-model="modelValue.users.max"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="users"
                        class="block w-full rounded-r-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6 disabled:text-base-400 transition-colors"
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
                    class="h-8 w-8 rounded-md transition-colors bg-base-700 text-primary-600 hover:text-primary-500 hover:bg-base-600 focus:ring-2 focus:ring-inset focus:ring-base-300 border-0"
                />
                <div class="flex flex-row">
                    <span
                        :class="[
                            modelValue.documents.req ? 'text-neutral' : 'text-base-400',
                            'transition-colors inline-flex h-9 items-center rounded-l-md border-0 px-3 bg-base-800 sm:text-sm',
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
                        class="block w-full rounded-r-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6 disabled:text-base-400 transition-colors"
                    />
                </div>
                <div class="flex flex-row">
                    <span
                        class="text-neutral transition-colors inline-flex h-9 items-center rounded-l-md border-0 px-3 bg-base-800 sm:text-sm"
                    >
                        {{ $t('common.max') }}
                    </span>
                    <input
                        v-model="modelValue.documents.max"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="documents"
                        class="block w-full rounded-r-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6 disabled:text-base-400 transition-colors"
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
                    class="h-8 w-8 rounded-md transition-colors bg-base-700 text-primary-600 hover:text-primary-500 hover:bg-base-600 focus:ring-2 focus:ring-inset focus:ring-base-300 border-0"
                />
                <div class="flex flex-row">
                    <span
                        :class="[
                            modelValue.vehicles.req ? 'text-neutral' : 'text-base-400',
                            'transition-colors inline-flex h-9 items-center rounded-l-md border-0 px-3 bg-base-800 sm:text-sm',
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
                        class="block w-full rounded-r-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6 disabled:text-base-400 transition-colors"
                    />
                </div>
                <div class="flex flex-row">
                    <span
                        class="text-neutral transition-colors inline-flex h-9 items-center rounded-l-md border-0 px-3 bg-base-800 sm:text-sm"
                    >
                        {{ $t('common.max') }}
                    </span>
                    <input
                        v-model="modelValue.vehicles.max"
                        type="number"
                        min="0"
                        oninput="validity.valid||(value='');"
                        name="vehicles"
                        class="block w-full rounded-r-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6 disabled:text-base-400 transition-colors"
                    />
                </div>
            </div>
        </div>
    </div>
</template>
