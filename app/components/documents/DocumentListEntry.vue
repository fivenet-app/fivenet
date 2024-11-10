<script lang="ts" setup>
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import DocumentCategoryBadge from '../partials/documents/DocumentCategoryBadge.vue';

defineProps<{
    document: DocumentShort;
}>();
</script>

<template>
    <li
        :key="document.id"
        class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 flex-initial border-white dark:border-gray-900"
        :class="[document.deletedAt ? 'bg-warn-100 hover:bg-warn-200 dark:bg-warn-800 dark:hover:bg-warn-700' : '']"
    >
        <ULink
            :to="{
                name: 'documents-id',
                params: { id: document.id },
            }"
        >
            <div class="m-2">
                <div class="flex flex-row justify-between gap-2">
                    <div class="flex flex-row items-center">
                        <IDCopyBadge
                            :id="document.id"
                            prefix="DOC"
                            :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                            :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                            size="xs"
                        />
                    </div>

                    <UBadge v-if="document.state" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-note-check" class="size-5" />
                        <span>
                            {{ document.state }}
                        </span>
                    </UBadge>

                    <div class="flex flex-row items-center gap-1">
                        <OpenClosedBadge :closed="document.closed" />
                    </div>

                    <div v-if="document.deletedAt" class="flex flex-1 flex-row items-center justify-center gap-1.5 font-bold">
                        <UIcon name="i-mdi-trash-can" class="size-5 shrink-0" />
                        {{ $t('common.deleted') }}
                    </div>
                </div>

                <div class="flex max-w-full shrink flex-row gap-2">
                    <div class="flex items-center gap-1">
                        <DocumentCategoryBadge :category="document.category" />

                        <h2 class="my-2 mr-2 line-clamp-1 flex-1 break-all text-xl font-medium hover:line-clamp-3">
                            {{ document.title }}
                        </h2>
                    </div>

                    <div v-if="document.updatedAt" class="flex flex-1 flex-row items-center justify-end gap-1.5">
                        <UIcon name="i-mdi-update" class="size-5 shrink-0" />
                        <p class="text-nowrap">
                            {{ $t('common.updated') }}
                            <GenericTime :value="document.updatedAt" :ago="true" />
                        </p>
                    </div>
                </div>

                <div class="flex flex-row gap-2">
                    <div class="flex flex-1 flex-row items-center justify-start">
                        <CitizenInfoPopover :user="document.creator" />
                    </div>

                    <div class="flex flex-1 flex-row items-center justify-center gap-1.5">
                        <UIcon name="i-mdi-briefcase" class="size-5 shrink-0" />
                        {{ document.creatorJobLabel }}
                    </div>

                    <div class="flex flex-1 flex-initial flex-row items-center justify-end gap-1.5">
                        <UIcon name="i-mdi-calendar" class="size-5 shrink-0" />
                        <p>
                            {{ $t('common.created_at') }}
                            <GenericTime :value="document.createdAt" />
                        </p>
                    </div>
                </div>
            </div>
        </ULink>
    </li>
</template>
