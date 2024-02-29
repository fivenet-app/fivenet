<script lang="ts" setup>
import type { DefineComponent } from 'vue';
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import {
    AccountMultipleIcon,
    ChevronDownIcon,
    CommentEditIcon,
    CommentPlusIcon,
    DeleteCircleIcon,
    DeleteCircleOutlineIcon,
    FileAccountIcon,
    FileMultipleIcon,
    FileSwapOutlineIcon,
    HelpIcon,
    LockCheckIcon,
    LockIcon,
    LockOpenIcon,
    LockOpenOutlineIcon,
    LockPlusOutlineIcon,
    LockQuestionIcon,
    NewBoxIcon,
    RefreshCircleIcon,
    TrashCanIcon,
    UpdateIcon,
} from 'mdi-vue3';
import { type DocActivity, DocActivityType } from '~~/gen/ts/resources/documents/activity';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import ActivityDocUpdatedDiff from '~/components/documents/activity/ActivityDocUpdatedDiff.vue';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';

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
            return FileMultipleIcon;
        case DocActivityType.ACCESS_UPDATED:
            return LockCheckIcon;
        case DocActivityType.OWNER_CHANGED:
            return FileAccountIcon;
        case DocActivityType.DELETED:
            return DeleteCircleIcon;

        // Requests
        case DocActivityType.REQUESTED_ACCESS:
            return LockPlusOutlineIcon;
        case DocActivityType.REQUESTED_CLOSURE:
            return LockQuestionIcon;
        case DocActivityType.REQUESTED_OPENING:
            return LockOpenOutlineIcon;
        case DocActivityType.REQUESTED_UPDATE:
            return RefreshCircleIcon;
        case DocActivityType.REQUESTED_OWNER_CHANGE:
            return FileSwapOutlineIcon;
        case DocActivityType.REQUESTED_DELETION:
            return DeleteCircleOutlineIcon;

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
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                <component :is="getDocAtivityIcon(entry.activityType)" class="h-7 w-7" aria-hidden="true" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="inline-flex items-center gap-2 text-sm font-medium text-neutral">
                        <span class="font-bold">
                            {{ $t(`enums.docstore.DocActivityType.${DocActivityType[entry.activityType]}`) }}
                        </span>
                        <span v-if="entry.data">
                            <span v-if="entry.data?.data.oneofKind === 'accessRequested'">
                                ({{ $t('common.access') }}:
                                {{
                                    $t(
                                        `enums.docstore.AccessLevel.${AccessLevel[entry.data?.data.accessRequested.level ?? 0]}`,
                                    )
                                }})
                            </span>
                        </span>
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="entry.createdAt" type="long" />
                    </p>
                </div>
                <p class="inline-flex text-sm text-gray-300">
                    {{ $t('common.created_by') }}
                    <CitizenInfoPopover class="ml-1" text-class="underline" :user="entry.creator" />
                </p>
            </div>
        </div>

        <Disclosure v-else v-slot="{ open }" as="div">
            <DisclosureButton class="flex w-full items-start justify-between text-left transition">
                <div class="flex w-full space-x-3">
                    <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                        <component :is="getDocAtivityIcon(entry.activityType)" class="h-7 w-7" aria-hidden="true" />
                    </div>
                    <div class="flex-1 space-y-1">
                        <div class="flex items-center justify-between">
                            <h3 class="inline-flex items-center text-sm font-medium text-neutral">
                                <span class="font-bold">
                                    {{ $t(`enums.docstore.DocActivityType.${DocActivityType[entry.activityType]}`) }}
                                </span>
                                <span class="ml-6 flex h-7 items-center">
                                    <ChevronDownIcon
                                        :class="[open ? 'upsidedown' : '', 'h-5 w-5 transition-transform']"
                                        aria-hidden="true"
                                    />
                                </span>
                            </h3>
                            <p class="text-sm text-gray-400">
                                <GenericTime :value="entry.createdAt" type="long" />
                            </p>
                        </div>
                        <p class="inline-flex text-sm text-gray-300">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" text-class="underline" :user="entry.creator" />
                        </p>
                    </div>
                </div>
            </DisclosureButton>
            <DisclosurePanel class="px-4 pb-2 pt-2">
                <template v-if="entry.activityType === DocActivityType.UPDATED">
                    <ActivityDocUpdatedDiff
                        v-if="entry.data?.data.oneofKind === 'updated'"
                        :update="entry.data?.data.updated"
                    />
                </template>
            </DisclosurePanel>
        </Disclosure>
    </div>
</template>
