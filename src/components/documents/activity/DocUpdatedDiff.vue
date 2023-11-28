<script lang="ts" setup>
defineProps<{
    update: DocUpdated;
}>();
</script>

<template>
    <div v-if="update.titleDiff">
        <p class="text-base font-semibold">
            {{ $t('common.title') }} {{ $t('components.documents.activity_list.difference') }}:
        </p>
        <div class="mt-2 mb-2 rounded-lg text-neutral bg-base-800 break-words">
            <span v-if="update.titleDiff?.length === 0">
                {{ $t('common.na') }}
            </span>
            <!-- eslint-disable vue/no-v-html -->
            <div v-else class="px-4 py-4" v-html="update.titleDiff"></div>
        </div>
    </div>

    <div v-if="update.contentDiff">
        <p class="text-base font-semibold">
            {{ $t('common.content') }} {{ $t('components.documents.activity_list.difference') }}:
        </p>
        <div class="mt-2 mb-2 rounded-lg text-neutral bg-base-800 break-words">
            <span v-if="update.contentDiff?.length === 0">
                {{ $t('common.na') }}
            </span>
            <!-- eslint-disable vue/no-v-html -->
            <div v-else class="px-4 py-4" v-html="update.contentDiff"></div>
        </div>
    </div>

    <div v-if="update.stateDiff">
        <p class="text-base font-semibold">
            {{ $t('common.state') }} {{ $t('components.documents.activity_list.difference') }}:
        </p>
        <div class="mt-2 mb-2 rounded-lg text-neutral bg-base-800 break-words">
            <span v-if="update.stateDiff?.length === 0">
                {{ $t('common.na') }}
            </span>
            <!-- eslint-disable vue/no-v-html -->
            <div v-else class="px-4 py-4" v-html="update.stateDiff"></div>
        </div>
    </div>
</template>
