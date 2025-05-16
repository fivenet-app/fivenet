<script lang="ts" setup>
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import OpenClosedBadge from '~/components/partials/OpenClosedBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DocumentCategoryBadge from '~/components/partials/documents/DocumentCategoryBadge.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { DocumentShort } from '~~/gen/ts/resources/documents/documents';

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
            <div class="m-2 flex flex-col gap-1">
                <div class="flex flex-row justify-between gap-2">
                    <div class="flex items-center">
                        <IDCopyBadge
                            :id="document.id"
                            prefix="DOC"
                            :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                            :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                            size="xs"
                        />
                    </div>

                    <UBadge v-if="document.state" class="inline-flex gap-1" size="md">
                        <UIcon class="size-4" name="i-mdi-note-check" />
                        <span>
                            {{ document.state }}
                        </span>
                    </UBadge>

                    <div v-if="document.deletedAt" class="flex flex-1 flex-row items-center justify-center gap-1.5 font-bold">
                        <UIcon class="size-4 shrink-0" name="i-mdi-delete" />
                        {{ $t('common.deleted') }}
                    </div>

                    <div class="flex flex-row items-center gap-1">
                        <OpenClosedBadge :closed="document.closed" />
                    </div>
                </div>

                <div class="flex max-w-full shrink flex-col gap-2">
                    <div class="flex flex-col gap-1 md:flex-row">
                        <div>
                            <DocumentCategoryBadge :category="document.category" />
                        </div>

                        <h2
                            class="line-clamp-2 flex-1 break-all text-lg font-medium hover:line-clamp-3 sm:text-xl md:line-clamp-1"
                        >
                            {{ document.title }}
                        </h2>
                    </div>
                </div>

                <div class="flex gap-2">
                    <div class="flex flex-1 items-center justify-start gap-1.5">
                        <UIcon class="size-4 shrink-0" name="i-mdi-calendar" />
                        <p class="inline-flex gap-1 text-nowrap">
                            <span class="hidden truncate md:block">
                                {{ $t('common.created_at') }}
                            </span>
                            <GenericTime :value="document.createdAt" />
                        </p>
                    </div>

                    <div v-if="document.workflowState?.autoCloseTime" class="flex flex-1 items-center justify-start gap-1.5">
                        <UIcon class="size-4 shrink-0" name="i-mdi-lock-clock" />
                        <p class="inline-flex gap-1 text-nowrap">
                            <span class="hidden truncate lg:block">
                                {{ $t('common.auto_close', 2) }}
                            </span>
                            <GenericTime :value="document.workflowState.autoCloseTime" ago />
                        </p>
                    </div>
                    <div
                        v-else-if="document.workflowState?.nextReminderTime"
                        class="flex flex-1 items-center justify-start gap-1.5"
                    >
                        <UIcon class="size-4 shrink-0" name="i-mdi-lock-clock" />
                        <p class="inline-flex gap-1 text-nowrap">
                            <span class="hidden truncate lg:block">
                                {{ $t('common.reminder') }}
                            </span>
                            <GenericTime :value="document.workflowState.nextReminderTime" ago />
                        </p>
                    </div>

                    <div v-if="document.updatedAt" class="flex flex-1 items-center justify-end gap-1.5">
                        <p class="inline-flex gap-1 truncate">
                            <span class="hidden md:block">
                                {{ $t('common.updated') }}
                            </span>
                            <GenericTime :value="document.updatedAt" ago />
                        </p>
                        <UIcon class="size-4 shrink-0" name="i-mdi-update" />
                    </div>
                </div>

                <div class="flex gap-2">
                    <CitizenInfoPopover :user="document.creator" />

                    <div class="flex flex-1 flex-row items-center justify-end gap-1.5">
                        <span>{{ document.creatorJobLabel }}</span>
                        <UIcon class="size-4 shrink-0" name="i-mdi-briefcase" />
                    </div>
                </div>
            </div>
        </ULink>
    </li>
</template>
