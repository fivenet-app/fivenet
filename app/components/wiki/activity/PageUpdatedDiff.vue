<script lang="ts" setup>
import TextDiff from '~/components/partials/content/TextDiff.vue';
import DiffBlock from '~/components/partials/DiffBlock.vue';
import type { PageUpdated } from '~~/gen/ts/resources/wiki/activity/activity';

defineProps<{
    update: PageUpdated;
}>();
</script>

<template>
    <div>
        <div v-if="update.titleDiff">
            <p class="text-base font-semibold">
                {{ $t('common.title') }} {{ $t('components.documents.activity_list.difference') }}:
            </p>
            <div class="my-2 rounded-lg break-words">
                <TextDiff v-if="update.titleCdiff" :ops="update.titleCdiff.ops" />
                <span v-else-if="update.titleDiff?.length === 0">
                    {{ $t('common.na') }}
                </span>
                <!-- eslint-disable vue/no-v-html -->
                <div v-else class="p-4" v-html="update.titleDiff"></div>
            </div>
        </div>

        <div v-if="update.descriptionDiff || update.descriptionCdiff">
            <p class="text-base font-semibold">
                {{ $t('common.state') }} {{ $t('components.documents.activity_list.difference') }}:
            </p>
            <div class="my-2 rounded-lg break-words">
                <TextDiff v-if="update.descriptionCdiff" :ops="update.descriptionCdiff.ops" />
                <span v-else-if="update.contentDiff?.length === 0">
                    {{ $t('common.na') }}
                </span>
                <div v-else-if="update.descriptionDiff?.startsWith('---')" class="p-4">
                    <DiffBlock :diff="update.descriptionDiff" />
                </div>
                <!-- eslint-disable vue/no-v-html -->
                <div v-else class="p-4" v-html="update.descriptionDiff"></div>
            </div>
        </div>

        <div v-if="update.contentDiff || update.contentCdiff">
            <p class="text-base font-semibold">
                {{ $t('common.content') }} {{ $t('components.documents.activity_list.difference') }}:
            </p>
            <div class="my-2 rounded-lg break-words">
                <TextDiff v-if="update.contentCdiff" :ops="update.contentCdiff.ops" />
                <span v-else-if="update.contentDiff?.length === 0">
                    {{ $t('common.na') }}
                </span>
                <div v-else-if="update.contentDiff?.startsWith('---')" class="p-4">
                    <DiffBlock :diff="update.contentDiff" />
                </div>
                <!-- eslint-disable vue/no-v-html -->
                <div v-else class="p-4" v-html="update.contentDiff"></div>
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
