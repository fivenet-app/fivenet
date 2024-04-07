<script lang="ts" setup>
import { ChevronDownIcon } from 'mdi-vue3';
import { type DocActivity, DocActivityType } from '~~/gen/ts/resources/documents/activity';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import ActivityDocUpdatedDiff from '~/components/documents/activity/ActivityDocUpdatedDiff.vue';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';

defineProps<{
    entry: DocActivity;
}>();

function spoilerNeeded(activityType: DocActivityType): boolean {
    switch (activityType) {
        case DocActivityType.UPDATED:
            return true;

        default:
            return false;
    }
}

function getDocAtivityIcon(activityType: DocActivityType): string {
    switch (activityType) {
        // Base
        case DocActivityType.CREATED:
            return 'i-mdi-new-box';
        case DocActivityType.STATUS_OPEN:
            return 'i-mdi-lock-open-variant';
        case DocActivityType.STATUS_CLOSED:
            return 'i-mdi-lock';
        case DocActivityType.UPDATED:
            return 'i-mdi-update';
        case DocActivityType.RELATIONS_UPDATED:
            return 'i-mdi-account-multiple';
        case DocActivityType.REFERENCES_UPDATED:
            return 'i-mdi-file-multiple';
        case DocActivityType.ACCESS_UPDATED:
            return 'i-mdi-lock-check';
        case DocActivityType.OWNER_CHANGED:
            return 'i-mdi-file-account';
        case DocActivityType.DELETED:
            return 'i-mdi-delete-circle';

        // Requests
        case DocActivityType.REQUESTED_ACCESS:
            return 'i-mdi-lock-plus-outline';
        case DocActivityType.REQUESTED_CLOSURE:
            return 'i-mdi-lock-question';
        case DocActivityType.REQUESTED_OPENING:
            return 'i-mdi-lock-open-outline';
        case DocActivityType.REQUESTED_UPDATE:
            return 'i-mdi-refresh-circle';
        case DocActivityType.REQUESTED_OWNER_CHANGE:
            return 'i-mdi-file-swap-outline';
        case DocActivityType.REQUESTED_DELETION:
            return 'i-mdi-delete-circle-outline';

        // Comments
        case DocActivityType.COMMENT_ADDED:
            return 'i-mdi-comment-plus';
        case DocActivityType.COMMENT_UPDATED:
            return 'i-mdi-comment-edit';
        case DocActivityType.COMMENT_DELETED:
            return 'i-mdi-trash-can';

        default:
            return 'i-mdi-help';
    }
}
</script>

<template>
    <div class="p-1">
        <div v-if="!spoilerNeeded(entry.activityType)" class="flex space-x-3">
            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                <UIcon :name="getDocAtivityIcon(entry.activityType)" class="size-7" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="inline-flex items-center gap-2 text-sm font-medium">
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
                <p class="inline-flex text-sm">
                    {{ $t('common.created_by') }}
                    <CitizenInfoPopover class="ml-1" :user="entry.creator" />
                </p>
            </div>
        </div>

        <UAccordion v-else :items="[{}]">
            <template #default="{ open }">
                <div class="flex space-x-3">
                    <div class="my-auto flex size-10 items-center justify-center rounded-full">
                        <UIcon :name="getDocAtivityIcon(entry.activityType)" class="size-7" />
                    </div>
                    <div class="flex-1 space-y-1">
                        <div class="flex items-center justify-between">
                            <h3 class="inline-flex items-center text-sm font-medium">
                                <span class="font-bold">
                                    {{ $t(`enums.docstore.DocActivityType.${DocActivityType[entry.activityType]}`) }}
                                </span>
                                <span class="ml-6 flex h-7 items-center">
                                    <ChevronDownIcon :class="[open ? 'upsidedown' : '', 'size-5 transition-transform']" />
                                </span>
                            </h3>
                            <p class="text-sm text-gray-400">
                                <GenericTime :value="entry.createdAt" type="long" />
                            </p>
                        </div>
                        <p class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="entry.creator" />
                        </p>
                    </div>
                </div>
            </template>

            <template v-if="entry.activityType === DocActivityType.UPDATED" #item>
                <div class="bg-background rounded-md p-2">
                    <ActivityDocUpdatedDiff
                        v-if="entry.data?.data.oneofKind === 'updated'"
                        :update="entry.data?.data.updated"
                    />
                </div>
            </template>
        </UAccordion>
    </div>
</template>
