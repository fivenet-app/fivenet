<script lang="ts" setup>
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { DocumentShort } from '~~/gen/ts/resources/documents/documents';

defineProps<{
    doc: DocumentShort;
}>();
</script>

<template>
    <li
        :key="doc.id"
        class="flex-initial rounded-lg"
        :class="[
            doc.deletedAt
                ? 'bg-warn-100 hover:bg-warn-200 dark:bg-warn-800 dark:hover:bg-warn-700'
                : 'bg-base-100 hover:bg-base-200 dark:bg-base-900 dark:hover:bg-base-700',
        ]"
    >
        <NuxtLink
            :to="{
                name: 'documents-id',
                params: { id: doc.id },
            }"
        >
            <div class="m-2">
                <div class="flex flex-row gap-2 truncate">
                    <div class="flex flex-1 flex-row items-center justify-start">
                        <IDCopyBadge
                            :id="doc.id"
                            prefix="DOC"
                            :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                            :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                            size="xs"
                        />
                    </div>

                    <UBadge v-if="doc.state" class="inline-flex gap-1" size="md">
                        <UIcon name="i-mdi-note-check" class="h-auto w-5" />
                        <span>
                            {{ doc.state }}
                        </span>
                    </UBadge>

                    <div class="flex flex-1 flex-row items-center justify-end gap-1">
                        <UBadge v-if="doc.closed" color="red" class="inline-flex gap-1" size="md">
                            <UIcon name="i-mdi-lock" color="red" class="h-auto w-5" />
                            <span>
                                {{ $t('common.close', 2) }}
                            </span>
                        </UBadge>
                        <UBadge v-else color="green" class="inline-flex gap-1" size="md">
                            <UIcon name="i-mdi-lock-open-variant" color="green" class="h-auto w-5" />
                            <span>
                                {{ $t('common.open', 2) }}
                            </span>
                        </UBadge>
                    </div>
                </div>

                <div class="flex flex-row gap-2 truncate">
                    <div class="inline-flex items-center gap-1">
                        <UBadge v-if="doc.category" class="inline-flex gap-1" size="md">
                            <UIcon name="i-mdi-shape" class="h-auto w-5" />
                            <span :title="doc.category.description ?? $t('common.na')">
                                {{ doc.category.name }}
                            </span>
                        </UBadge>

                        <h2 class="truncate py-2 pr-3 text-xl font-medium">
                            {{ doc.title }}
                        </h2>
                    </div>

                    <div v-if="doc.deletedAt" class="flex flex-1 flex-row items-center justify-center font-bold">
                        <UIcon name="i-mdi-trash-can" class="mr-1.5 size-5 shrink-0" />
                        {{ $t('common.deleted') }}
                    </div>

                    <div v-if="doc.updatedAt" class="flex flex-1 flex-row items-center justify-end">
                        <UIcon name="i-mdi-update" class="mr-1.5 size-5 shrink-0" />
                        <p>
                            {{ $t('common.updated') }}
                            <GenericTime :value="doc.updatedAt" :ago="true" />
                        </p>
                    </div>
                </div>

                <div class="flex flex-row gap-2">
                    <div class="flex flex-1 flex-row items-center justify-start">
                        <CitizenInfoPopover :user="doc.creator" />
                    </div>

                    <div class="flex flex-1 flex-row items-center justify-center">
                        <UIcon name="i-mdi-briefcase" class="mr-1.5 size-5 shrink-0" />
                        {{ doc.creator?.jobLabel }}
                    </div>

                    <div class="flex flex-1 flex-row items-center justify-end">
                        <UIcon name="i-mdi-calendar" class="mr-1.5 size-5 shrink-0" />
                        <p>
                            {{ $t('common.created_at') }}
                            <GenericTime :value="doc.createdAt" />
                        </p>
                    </div>
                </div>
            </div>
        </NuxtLink>
    </li>
</template>
