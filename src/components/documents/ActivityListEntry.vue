<script lang="ts" setup>
import type { DefineComponent } from 'vue';
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import {
    AccountMultipleIcon,
    ChevronDownIcon,
    CommentEditIcon,
    CommentPlusIcon,
    DeleteCircleIcon,
    FileDocumentIcon,
    FileSwapIcon,
    HelpIcon,
    LockCheckIcon,
    LockIcon,
    LockOpenIcon,
    LockPlusIcon,
    LockQuestionIcon,
    NewBoxIcon,
    RefreshIcon,
    TrashCanIcon,
    UpdateIcon,
} from 'mdi-vue3';
import { type DocActivity, DocActivityType } from '~~/gen/ts/resources/documents/activity';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import Time from '~/components/partials/elements/Time.vue';

defineProps<{
    entry: DocActivity;
}>();

function disclosureNeeded(activityType: DocActivityType): boolean {
    switch (activityType) {
        case DocActivityType.UPDATED:
            return true;

        default:
            return false;
    }
}

function getDocAtivityIcon(activityType: DocActivityType): DefineComponent {
    switch (activityType) {
        // Base
        case DocActivityType.CREATED:
            return NewBoxIcon;
        case DocActivityType.STATUS_OPEN:
            return LockOpenIcon;
        case DocActivityType.STATUS_CLOSED:
            return LockIcon;
        case DocActivityType.UPDATED:
            return UpdateIcon;
        case DocActivityType.RELATIONS_UPDATED:
            return AccountMultipleIcon;
        case DocActivityType.REFERENCES_UPDATED:
            return FileDocumentIcon;
        case DocActivityType.ACCESS_UPDATED:
            return LockCheckIcon;
        case DocActivityType.OWNER_CHANGED:
            return FileSwapIcon;
        case DocActivityType.DELETED:
            return TrashCanIcon;

        // Requests
        case DocActivityType.REQUESTED_ACCESS:
            return LockPlusIcon;
        case DocActivityType.REQUESTED_CLOSURE:
            return LockQuestionIcon;
        case DocActivityType.REQUESTED_UPDATE:
            return RefreshIcon;
        case DocActivityType.REQUESTED_DELETION:
            return DeleteCircleIcon;

        // Comments
        case DocActivityType.COMMENT_ADDED:
            return CommentPlusIcon;
        case DocActivityType.COMMENT_UPDATED:
            return CommentEditIcon;
        case DocActivityType.COMMENT_DELETED:
            return TrashCanIcon;

        default:
            return HelpIcon;
    }
}
</script>

<template>
    <div class="p-1">
        <div v-if="!disclosureNeeded(entry.activityType)" class="flex space-x-3">
            <div class="h-10 w-10 rounded-full flex items-center justify-center my-auto">
                <component :is="getDocAtivityIcon(entry.activityType)" class="w-7 h-7" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium text-neutral">
                        <span class="font-bold">
                            {{ $t(`enums.docstore.DocActivityType.${DocActivityType[entry.activityType]}`) }}
                        </span>
                    </h3>
                    <p class="text-sm text-gray-400">
                        <Time :value="entry.createdAt" type="long" />
                    </p>
                </div>
                <p class="text-sm text-gray-300 inline-flex">
                    {{ $t('common.created_by') }}
                    <CitizenInfoPopover class="ml-1" text-class="underline" :user="entry.creator" />
                </p>
            </div>
        </div>

        <Disclosure v-else v-slot="{ open }" as="div">
            <DisclosureButton class="flex w-full items-start justify-between text-left transition">
                <div class="w-full flex space-x-3">
                    <div class="h-10 w-10 rounded-full flex items-center justify-center my-auto">
                        <component :is="getDocAtivityIcon(entry.activityType)" class="w-7 h-7" />
                    </div>
                    <div class="flex-1 space-y-1">
                        <div class="flex items-center justify-between">
                            <h3 class="inline-flex items-center text-sm font-medium text-neutral">
                                <span class="font-bold">
                                    {{ $t(`enums.docstore.DocActivityType.${DocActivityType[entry.activityType]}`) }}
                                </span>
                                <span class="ml-6 flex h-7 items-center">
                                    <ChevronDownIcon
                                        :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                                        aria-hidden="true"
                                    />
                                </span>
                            </h3>
                            <p class="text-sm text-gray-400">
                                <Time :value="entry.createdAt" type="long" />
                            </p>
                        </div>
                        <p class="text-sm text-gray-300 inline-flex">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" text-class="underline" :user="entry.creator" />
                        </p>
                    </div>
                </div>
            </DisclosureButton>
            <DisclosurePanel class="px-4 pt-2 pb-2">
                <template v-if="entry.activityType === DocActivityType.UPDATED">
                    <p class="text-base font-semibold">{{ $t('components.documents.activity_list.difference') }}:</p>
                    <!-- eslint-disable vue/no-v-html -->
                    <div
                        v-if="entry.data?.data.oneofKind === 'updated'"
                        class="mt-2 mb-2 rounded-lg text-neutral bg-base-800 break-words"
                    >
                        <span v-if="entry.data.data.updated.diff.length === 0">
                            {{ $t('common.na') }}
                        </span>
                        <div v-else class="px-4 py-4" v-html="entry.data.data.updated.diff"></div>
                    </div>

                    <span class="inline-flex gap-2">
                        <span class="text-base font-semibold">{{ $t('common.legend') }}:</span>
                        <span class="bg-success-600">{{ $t('components.documents.activity_list.legend.added') }}</span>
                        <span class="bg-error-600">{{ $t('components.documents.activity_list.legend.removed') }}</span>
                        <span class="bg-info-600">{{ $t('components.documents.activity_list.legend.changed') }}</span>
                    </span>
                </template>
            </DisclosurePanel>
        </Disclosure>
    </div>
</template>
