<script lang="ts" setup>
import type { DocUpdated } from '~~/gen/ts/resources/documents/activity';

defineProps<{
    update: DocUpdated;
}>();
</script>

<template>
    <div>
        <div v-if="update.titleDiff">
            <p class="text-base font-semibold">
                {{ $t('common.title') }} {{ $t('components.documents.activity_list.difference') }}:
            </p>
            <div class="my-2 break-words rounded-lg">
                <span v-if="update.titleDiff?.length === 0">
                    {{ $t('common.na') }}
                </span>
                <!-- eslint-disable vue/no-v-html -->
                <div v-else class="p-4" v-html="update.titleDiff"></div>
            </div>
        </div>

        <div v-if="update.contentDiff">
            <p class="text-base font-semibold">
                {{ $t('common.content') }} {{ $t('components.documents.activity_list.difference') }}:
            </p>
            <div class="my-2 break-words rounded-lg">
                <span v-if="update.contentDiff?.length === 0">
                    {{ $t('common.na') }}
                </span>
                <!-- eslint-disable vue/no-v-html -->
                <div v-else class="p-4" v-html="update.contentDiff"></div>
            </div>
        </div>

        <div v-if="update.stateDiff">
            <p class="text-base font-semibold">
                {{ $t('common.state') }} {{ $t('components.documents.activity_list.difference') }}:
            </p>
            <div class="my-2 break-words rounded-lg">
                <span v-if="update.stateDiff?.length === 0">
                    {{ $t('common.na') }}
                </span>
                <!-- eslint-disable vue/no-v-html -->
                <div v-else class="p-4" v-html="update.stateDiff"></div>
            </div>
        </div>

        <span class="inline-flex gap-2">
            <span class="text-base font-semibold text-white">{{ $t('common.legend') }}:</span>
            <span class="bg-success-600 text-white">{{ $t('components.documents.activity_list.legend.added') }}</span>
            <span class="bg-error-600 text-white">{{ $t('components.documents.activity_list.legend.removed') }}</span>
            <span class="bg-info-600 text-white">{{ $t('components.documents.activity_list.legend.changed') }}</span>
        </span>
    </div>
</template>
