<script setup lang="ts">
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import { UserActivityType, type UserActivity } from '~~/gen/ts/resources/users/activity';
import type { CitizenLabels } from '~~/gen/ts/resources/users/labels';

const props = defineProps<{
    activity: UserActivity;
}>();
</script>

<template>
    <template v-if="activity.data">
        <template v-if="activity.type === UserActivityType.NAME && activity.data?.data.oneofKind === 'nameChange'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon name="i-mdi-identification-card" class="size-full text-info-600" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            <I18nT keypath="components.citizens.CitizenInfoActivityFeedEntry.name_change">
                                <template #old>
                                    <span class="font-semibold">
                                        {{ activity.data.data.nameChange.old }}
                                    </span>
                                </template>
                                <template #new>
                                    <span class="font-semibold">
                                        {{ activity.data.data.nameChange.new }}
                                    </span>
                                </template>
                            </I18nT>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template
            v-else-if="activity.type === UserActivityType.DOCUMENT && activity.data?.data.oneofKind === 'documentRelation'"
        >
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon
                        v-if="activity.data.data.documentRelation.added"
                        name="i-mdi-file-account"
                        class="size-full text-info-600"
                    />
                    <UIcon v-else name="i-mdi-file-account-outline" class="size-full text-base-600" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex items-center gap-1 text-sm font-medium">
                            <template v-if="activity.data.data.documentRelation.added">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.document_relation.added') }}
                            </template>
                            <template v-else>
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.document_relation.removed') }}
                            </template>

                            <DocumentInfoPopover
                                :document-id="activity.data.data.documentRelation.documentId"
                                hide-trailing
                                load-on-open
                            >
                                <template #title>
                                    <IDCopyBadge
                                        :id="activity.data.data.documentRelation.documentId"
                                        prefix="DOC"
                                        size="xs"
                                        disable-tooltip
                                        hide-icon
                                    />
                                </template>
                            </DocumentInfoPopover>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.type') }}:</span>
                            <span>
                                {{
                                    $t(
                                        `enums.docstore.DocRelation.${DocRelation[activity.data.data.documentRelation.relation]}`,
                                    )
                                }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.type === UserActivityType.WANTED && activity.data?.data.oneofKind === 'wantedChange'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon
                        v-if="activity.data.data.wantedChange.wanted"
                        name="i-mdi-bell-alert"
                        class="size-full text-error-400"
                    />
                    <UIcon v-else name="i-mdi-bell-sleep" class="size-full text-success-400" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.bool_set_citizen') }}
                            <span class="font-semibold">
                                {{
                                    activity.data.data.wantedChange.wanted
                                        ? $t('common.wanted')
                                        : `${$t('common.not')} ${$t('common.wanted')}`
                                }}
                            </span>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.type === UserActivityType.JOB && activity.data?.data.oneofKind === 'jobChange'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon name="i-mdi-briefcase" class="size-full text-gray-400" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.userprops_job_set') }}
                            <span class="font-semibold">
                                {{ activity.data.data.jobChange.jobLabel }}
                                <span v-if="activity.data.data.jobChange.grade">
                                    ({{ $t('common.grade') }}: {{ activity.data.data.jobChange.gradeLabel }})</span
                                >
                            </span>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template
            v-else-if="
                activity.type === UserActivityType.TRAFFIC_INFRACTION_POINTS &&
                activity.data.data.oneofKind === 'trafficInfractionPointsChange'
            "
        >
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon
                        name="i-mdi-traffic-cone"
                        class="size-full"
                        :class="
                            activity.data.data.trafficInfractionPointsChange.old >
                            activity.data.data.trafficInfractionPointsChange.new
                                ? 'text-gray-400'
                                : 'text-orange-400'
                        "
                    />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            <I18nT
                                keypath="components.citizens.CitizenInfoActivityFeedEntry.traffic_infraction_points.action_text"
                            >
                                <template #old>
                                    <span class="font-semibold">{{
                                        activity.data.data.trafficInfractionPointsChange.old
                                    }}</span>
                                </template>
                                <template #new>
                                    <span class="font-semibold">{{
                                        activity.data.data.trafficInfractionPointsChange.new
                                    }}</span>
                                </template>
                            </I18nT>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.type === UserActivityType.MUGSHOT && activity.data?.data.oneofKind === 'mugshotChange'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon
                        name="i-mdi-camera-account"
                        class="size-full text-amber-400"
                        :class="activity.data.data.mugshotChange.new ? 'text-gray-400' : 'text-amber-400'"
                    />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            <template v-if="activity.data.data.mugshotChange.new">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.userprops_mug_shot_set') }}
                            </template>
                            <template v-else>
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.userprops_mug_shot_removed') }}
                            </template>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.type === UserActivityType.LABELS && activity.data?.data.oneofKind === 'labelsChange'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon name="i-mdi-tag" class="size-full text-amber-200" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex flex-col gap-1 text-sm font-medium">
                            <span>
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.userprops_labels_updated') }}
                            </span>

                            <div class="inline-flex gap-1">
                                <UBadge
                                    v-for="attribute in activity.data.data.labelsChange.removed"
                                    :key="attribute.name"
                                    :style="{ backgroundColor: attribute.color }"
                                    class="justify-between gap-2 line-through"
                                    :class="
                                        isColourBright(hexToRgb(attribute.color, RGBBlack)!) ? '!text-black' : '!text-white'
                                    "
                                    size="xs"
                                >
                                    {{ attribute.name }}
                                </UBadge>

                                <UBadge
                                    v-for="attribute in activity.data.data.labelsChange.added"
                                    :key="attribute.name"
                                    :style="{ backgroundColor: attribute.color }"
                                    class="justify-between gap-2"
                                    :class="
                                        isColourBright(hexToRgb(attribute.color, RGBBlack)!) ? '!text-black' : '!text-white'
                                    "
                                    size="xs"
                                >
                                    {{ attribute.name }}
                                </UBadge>
                            </div>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.type === UserActivityType.LICENSES && activity.data?.data.oneofKind === 'licensesChange'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon
                        name="i-mdi-license"
                        class="size-full"
                        :class="activity.data.data.licensesChange.added ? 'text-info-600' : 'text-amber-600'"
                    />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            {{
                                activity.data.data.licensesChange.added
                                    ? $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.added')
                                    : $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.removed')
                            }}: {{ activity.data.data.licensesChange.licenses.map((l) => l.label).join(', ') }}
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.type === UserActivityType.JAIL && activity.data?.data.oneofKind === 'jailChange'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon v-if="activity.data.data.jailChange.seconds > 0" name="i-mdi-handcuffs" class="size-full" />
                    <UIcon v-else-if="activity.data.data.jailChange.seconds === 0" name="i-mdi-door-open" class="size-full" />
                    <UIcon v-else name="i-mdi-run-fast" class="size-full" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            <template v-if="activity.data.data.jailChange.seconds > 0">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.jailed') }}
                                {{ fromSecondsToFormattedDuration(activity.data.data.jailChange.seconds) }}
                            </template>
                            <template v-else-if="activity.data.data.jailChange.seconds === 0">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.unjailed') }}
                            </template>
                            <template v-else>
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.escaped') }}
                            </template>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <template v-if="activity.data.data.jailChange.seconds >= 0">
                                <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                                <span>
                                    {{ !activity.reason ? $t('common.na') : activity.reason }}
                                </span>
                            </template>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.type === UserActivityType.FINE && activity.data?.data.oneofKind === 'fineChange'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon
                        v-if="activity.data.data.fineChange.removed"
                        name="i-mdi-receipt-text-remove"
                        class="size-full text-red-400"
                    />
                    <UIcon
                        v-else-if="activity.data.data.fineChange.amount < 0"
                        name="i-mdi-receipt-text-check"
                        class="size-full text-success-400"
                    />
                    <UIcon v-else name="i-mdi-receipt-text-plus" class="size-full text-info-400" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex gap-1 text-sm font-medium">
                            <template v-if="activity.data.data.fineChange.removed">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.removed') }}
                            </template>
                            <template v-else-if="activity.data.data.fineChange.amount < 0">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.paid') }}
                            </template>
                            <template v-else>
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.created') }}
                            </template>

                            <span>
                                {{ $n(Math.abs(activity.data.data.fineChange.amount), 'currency') }}
                            </span>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else>
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon name="i-mdi-help-circle" class="size-full" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            {{ UserActivityType[props.activity.type] }}
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
    </template>

    <!-- Old `key` based system that uses oldValue and newValue -->
    <template v-else>
        <template v-if="activity.key === 'DocStore.Relation'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon v-if="activity.newValue !== ''" name="i-mdi-file-account" class="size-full text-info-600" />
                    <UIcon v-else name="i-mdi-file-account-outline" class="size-full text-base-600" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex gap-1 text-sm font-medium">
                            <span v-if="activity.newValue !== ''">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.document_relation.added') }}
                            </span>
                            <span v-else>
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.document_relation.removed') }}
                            </span>

                            <DocumentInfoPopover
                                :document-id="activity.newValue !== '' ? activity.newValue : activity.oldValue"
                                hide-trailing
                                load-on-open
                            >
                                <template #title>
                                    <IDCopyBadge
                                        :id="activity.newValue !== '' ? activity.newValue : activity.oldValue"
                                        prefix="DOC"
                                        size="xs"
                                        disable-tooltip
                                        hide-icon
                                    />
                                </template>
                            </DocumentInfoPopover>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ $t(`enums.docstore.DocRelation.${activity.reason.replace('DOC_RELATION_', '')}`) }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.key === 'UserProps.Wanted'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon v-if="activity.newValue === 'true'" name="i-mdi-bell-alert" class="size-full text-error-400" />
                    <UIcon v-else name="i-mdi-bell-sleep" class="size-full text-success-400" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.bool_set_citizen') }}
                            <span class="font-semibold">
                                {{
                                    activity.newValue === 'true'
                                        ? $t('common.wanted')
                                        : `${$t('common.not').toLowerCase()} ${$t('common.wanted')}`
                                }}
                            </span>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.key === 'UserProps.Job'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon name="i-mdi-briefcase" class="size-full text-gray-400" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.userprops_job_set') }}
                            <span class="font-semibold">
                                {{ activity.newValue }}
                            </span>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.key === 'UserProps.TrafficInfractionPoints'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon name="i-mdi-traffic-cone" class="size-full text-gray-400" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            <I18nT
                                keypath="components.citizens.CitizenInfoActivityFeedEntry.traffic_infraction_points.action_text"
                            >
                                <template #old>
                                    <span class="font-semibold">{{ activity.oldValue }}</span>
                                </template>
                                <template #new>
                                    <span class="font-semibold">{{ activity.newValue }}</span>
                                </template>
                            </I18nT>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.key === 'UserProps.MugShot'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon name="i-mdi-camera-account" class="size-full text-gray-400" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.userprops_mug_shot_set') }}
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.key === 'UserProps.Labels'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon name="i-mdi-tag" class="size-full text-amber-200" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex flex-col gap-1 text-sm font-medium">
                            <span>
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.userprops_labels_updated') }}
                            </span>

                            <div class="inline-flex gap-1">
                                <UBadge
                                    v-for="attribute in (JSON.parse(activity.oldValue) as CitizenLabels)?.list"
                                    :key="attribute.name"
                                    :style="{ backgroundColor: attribute.color }"
                                    class="justify-between gap-2 line-through"
                                    :class="
                                        isColourBright(hexToRgb(attribute.color, RGBBlack)!) ? '!text-black' : '!text-white'
                                    "
                                    size="xs"
                                >
                                    {{ attribute.name }}
                                </UBadge>

                                <UBadge
                                    v-for="attribute in (JSON.parse(activity.newValue) as CitizenLabels)?.list"
                                    :key="attribute.name"
                                    :style="{ backgroundColor: attribute.color }"
                                    class="justify-between gap-2"
                                    :class="
                                        isColourBright(hexToRgb(attribute.color, RGBBlack)!) ? '!text-black' : '!text-white'
                                    "
                                    size="xs"
                                >
                                    {{ attribute.name }}
                                </UBadge>
                            </div>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.key === 'Plugin.Licenses'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon
                        name="i-mdi-license"
                        class="size-full"
                        :class="activity.newValue !== '' ? 'text-info-600' : 'text-amber-600'"
                    />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            {{
                                activity.newValue !== ''
                                    ? $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.added')
                                    : $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.removed')
                            }}: {{ activity.reason }}
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm"></p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.key === 'Plugin.Jail'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon
                        v-if="activity.oldValue === '' && activity.newValue !== '0'"
                        name="i-mdi-handcuffs"
                        class="size-full"
                    />
                    <UIcon v-else-if="activity.newValue === '0'" name="i-mdi-door-open" class="size-full" />
                    <UIcon v-else name="i-mdi-run-fast" class="size-full" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            <template v-if="activity.oldValue === '' && activity.newValue !== '0'">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.jailed') }}
                                {{ fromSecondsToFormattedDuration(parseInt(props.activity.newValue)) }}
                            </template>
                            <template v-else-if="activity.newValue === '0'">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.unjailed') }}
                            </template>
                            <template v-else>
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.escaped') }}
                            </template>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <template v-if="activity.oldValue === '' && activity.newValue !== '0'">
                                <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                                <span>
                                    {{ !activity.reason ? $t('common.na') : activity.reason }}
                                </span>
                            </template>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else-if="activity.key === 'Plugin.Billing.Fines'">
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon
                        v-if="activity.newValue === '0'"
                        name="i-mdi-receipt-text-check"
                        class="size-full text-success-400"
                    />
                    <UIcon
                        v-else-if="activity.newValue === activity.oldValue"
                        name="i-mdi-receipt-text-remove"
                        class="size-full text-gray-400"
                    />
                    <UIcon v-else name="i-mdi-receipt-text-plus" class="size-full text-info-400" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="inline-flex gap-1 text-sm font-medium">
                            <template v-if="activity.newValue === '0'">
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.paid') }}
                            </template>
                            <template v-else-if="activity.newValue === activity.oldValue">
                                {{
                                    $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.removed')
                                }}</template
                            >
                            <template v-else>
                                {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.created') }}
                            </template>

                            <span>
                                {{ $n(parseInt(props.activity.newValue), 'currency') }}
                            </span>
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
        <template v-else>
            <div class="flex space-x-3">
                <div class="my-auto flex size-10 items-center justify-center rounded-full">
                    <UIcon name="i-mdi-help-circle" class="size-full" />
                </div>

                <div class="flex-1 space-y-1">
                    <div class="flex items-center justify-between">
                        <h3 class="text-sm font-medium">
                            {{ `${props.activity.oldValue} -> ${props.activity.newValue}` }}
                        </h3>

                        <p class="text-sm">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

                    <div class="flex items-center justify-between">
                        <p class="inline-flex gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </p>

                        <p v-if="activity.sourceUser" class="inline-flex text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
    </template>
</template>
