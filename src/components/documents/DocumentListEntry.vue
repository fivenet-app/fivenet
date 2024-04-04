<script lang="ts" setup>
import { BriefcaseIcon, CalendarIcon, LockIcon, LockOpenVariantIcon, TrashCanIcon, UpdateIcon } from 'mdi-vue3';
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
        :class="[
            doc.deletedAt ? 'bg-warn-800 hover:bg-warn-700' : 'bg-base-900 hover:bg-base-800',
            'my-1 flex-initial rounded-lg',
        ]"
    >
        <NuxtLink
            :to="{
                name: 'documents-id',
                params: { id: doc.id },
            }"
        >
            <div class="m-2">
                <div class="flex flex-row gap-2 truncate text-base-300">
                    <div class="flex flex-1 flex-row items-center justify-start">
                        <IDCopyBadge
                            :id="doc.id"
                            prefix="DOC"
                            :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                            :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                        />
                    </div>
                    <p
                        v-if="doc.state"
                        class="my-auto inline-flex rounded-full bg-info-100 px-2 py-1 text-xs font-semibold leading-5 text-info-800"
                    >
                        {{ doc.state }}
                    </p>
                    <div class="flex flex-1 flex-row items-center justify-end">
                        <div v-if="doc?.closed" class="inline-flex flex-initial gap-1">
                            <LockIcon class="size-5 text-error-400" />
                            <span class="text-sm font-medium text-error-400">
                                {{ $t('common.close', 2) }}
                            </span>
                        </div>
                        <div v-else class="inline-flex flex-initial gap-1">
                            <LockOpenVariantIcon class="size-5 text-success-400" />
                            <span class="text-sm font-medium text-success-400">
                                {{ $t('common.open', 2) }}
                            </span>
                        </div>
                    </div>
                </div>

                <div class="flex flex-row gap-2 truncate">
                    <h2 class="inline-flex items-center gap-1 truncate text-xl font-medium">
                        <span
                            v-if="doc.category"
                            class="flex flex-initial flex-row gap-1 break-words rounded-full bg-primary-100 px-2 py-1 text-primary-500"
                        >
                            <span
                                class="inline-flex items-center text-xs font-medium text-primary-800"
                                :title="doc.category.description ?? $t('common.na')"
                            >
                                {{ doc.category.name }}
                            </span>
                        </span>
                        <span class="py-2 pr-3">
                            {{ doc.title }}
                        </span>
                    </h2>
                    <div v-if="doc.deletedAt" class="flex flex-1 flex-row items-center justify-center font-bold">
                        <TrashCanIcon class="mr-1.5 size-5 shrink-0 text-base-300" />
                        {{ $t('common.deleted') }}
                    </div>
                    <div v-if="doc.updatedAt" class="flex flex-1 flex-row items-center justify-end">
                        <UpdateIcon class="mr-1.5 size-5 shrink-0 text-base-300" />
                        <p>
                            {{ $t('common.updated') }}
                            <GenericTime :value="doc.updatedAt" :ago="true" />
                        </p>
                    </div>
                </div>

                <div class="flex flex-row gap-2 text-accent-200">
                    <div class="flex flex-1 flex-row items-center justify-start">
                        <CitizenInfoPopover :user="doc.creator" />
                    </div>
                    <div class="flex flex-1 flex-row items-center justify-center">
                        <BriefcaseIcon class="mr-1.5 size-5 shrink-0 text-base-300" />
                        {{ doc.creator?.jobLabel }}
                    </div>
                    <div class="flex flex-1 flex-row items-center justify-end">
                        <CalendarIcon class="mr-1.5 size-5 shrink-0 text-base-300" />
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
