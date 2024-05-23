<script lang="ts" setup>
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import OpenClosedBadge from '../partials/OpenClosedBadge.vue';

defineProps<{
    document: DocumentShort;
}>();
</script>

<template>
    <li
        :key="document.id"
        class="flex-initial"
        :class="[document.deletedAt ? 'bg-warn-100 hover:bg-warn-200 dark:bg-warn-800 dark:hover:bg-warn-700' : '']"
    >
        <ULink
            :to="{
                name: 'documents-id',
                params: { id: document.id },
            }"
        >
            <div class="m-2">
                <div class="flex flex-row gap-2 truncate">
                    <div class="flex flex-1 flex-row items-center justify-start">
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

                    <div class="flex flex-1 flex-row items-center justify-end gap-1">
                        <OpenClosedBadge :closed="document.closed" />
                    </div>
                </div>

                <div class="flex flex-row gap-2 truncate">
                    <div class="inline-flex items-center gap-1 truncate">
                        <UBadge v-if="document.category" class="inline-flex gap-1" size="md">
                            <UIcon name="i-mdi-shape" class="size-5" />
                            <span :title="document.category.description ?? $t('common.na')">
                                {{ document.category.name }}
                            </span>
                        </UBadge>

                        <h2 class="truncate py-2 pr-3 text-xl font-medium">
                            {{ document.title }}
                        </h2>
                    </div>

                    <div v-if="document.deletedAt" class="flex flex-1 flex-row items-center justify-center font-bold">
                        <UIcon name="i-mdi-trash-can" class="mr-1.5 size-5 shrink-0" />
                        {{ $t('common.deleted') }}
                    </div>

                    <div v-if="document.updatedAt" class="flex flex-1 flex-row items-center justify-end">
                        <UIcon name="i-mdi-update" class="mr-1.5 size-5 shrink-0" />
                        <p>
                            {{ $t('common.updated') }}
                            <GenericTime :value="document.updatedAt" :ago="true" />
                        </p>
                    </div>
                </div>

                <div class="flex flex-row gap-2">
                    <div class="flex flex-1 flex-row items-center justify-start">
                        <CitizenInfoPopover :user="document.creator" />
                    </div>

                    <div class="flex flex-1 flex-row items-center justify-center">
                        <UIcon name="i-mdi-briefcase" class="mr-1.5 size-5 shrink-0" />
                        {{ document.creator?.jobLabel }}
                    </div>

                    <div class="flex flex-1 flex-row items-center justify-end">
                        <UIcon name="i-mdi-calendar" class="mr-1.5 size-5 shrink-0" />
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
