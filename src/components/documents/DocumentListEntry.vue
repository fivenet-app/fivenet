<script lang="ts" setup>
import { AccountIcon, BriefcaseIcon, CalendarIcon, LockIcon, LockOpenVariantIcon, TrashCanIcon, UpdateIcon } from 'mdi-vue3';
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
            doc.deletedAt ? 'hover:bg-warn-700 bg-warn-800' : 'hover:bg-base-700 bg-base-800',
            'flex-initial my-1 rounded-lg',
        ]"
    >
        <NuxtLink
            :to="{
                name: 'documents-id',
                params: { id: doc.id },
            }"
        >
            <div class="m-2">
                <div class="flex flex-row gap-2 text-base-300 truncate">
                    <div class="flex flex-row items-center justify-start flex-1">
                        <IDCopyBadge
                            :id="doc.id"
                            prefix="DOC"
                            :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                            :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                        />
                    </div>
                    <p
                        v-if="doc.state"
                        class="inline-flex px-2 py-1 text-xs font-semibold leading-5 rounded-full bg-info-100 text-info-800 my-auto"
                    >
                        {{ doc.state }}
                    </p>
                    <div class="flex flex-row items-center justify-end flex-1">
                        <div v-if="doc?.closed" class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-error-100">
                            <LockIcon class="w-5 h-5 text-error-400" aria-hidden="true" />
                            <span class="text-sm font-medium text-error-700">
                                {{ $t('common.close', 2) }}
                            </span>
                        </div>
                        <div v-else class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-success-100">
                            <LockOpenVariantIcon class="w-5 h-5 text-success-500" aria-hidden="true" />
                            <span class="text-sm font-medium text-success-700">
                                {{ $t('common.open', 2) }}
                            </span>
                        </div>
                    </div>
                </div>

                <div class="flex flex-row gap-2 text-base-200 truncate">
                    <h2 class="inline-flex items-center gap-1 text-xl font-medium text-neutral truncate">
                        <span
                            v-if="doc.category"
                            class="flex flex-row flex-initial gap-1 px-2 py-1 rounded-full bg-primary-100 text-primary-500 break-words"
                        >
                            <span
                                class="text-xs font-medium text-primary-800 inline-flex items-center"
                                :title="doc.category.description ?? $t('common.na')"
                            >
                                {{ doc.category.name }}
                            </span>
                        </span>
                        <span class="py-2 pr-3">
                            {{ doc.title }}
                        </span>
                    </h2>
                    <div v-if="doc.deletedAt" type="button" class="flex flex-row items-center justify-center flex-1 font-bold">
                        <TrashCanIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                        {{ $t('common.deleted') }}
                    </div>
                    <div v-if="doc.updatedAt" class="flex flex-row items-center justify-end flex-1">
                        <UpdateIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                        <p>
                            {{ $t('common.updated') }}
                            <GenericTime :value="doc.updatedAt" :ago="true" />
                        </p>
                    </div>
                </div>

                <div class="flex flex-row gap-2 text-base-200">
                    <div class="flex flex-row items-center justify-start flex-1">
                        <CitizenInfoPopover :user="doc.creator">
                            <template #before>
                                <AccountIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                            </template>
                        </CitizenInfoPopover>
                    </div>
                    <div class="flex flex-row items-center justify-center flex-1">
                        <BriefcaseIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
                        {{ doc.creator?.jobLabel }}
                    </div>
                    <div class="flex flex-row items-center justify-end flex-1">
                        <CalendarIcon class="mr-1.5 h-5 w-5 flex-shrink-0 text-base-400" aria-hidden="true" />
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
