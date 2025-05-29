<script setup lang="ts">
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { DocRelation } from '~~/gen/ts/resources/documents/documents';
import { UserActivityType, type UserActivity } from '~~/gen/ts/resources/users/activity';

const props = defineProps<{
    activity: UserActivity;
}>();
</script>

<template>
    <template v-if="activity.type === UserActivityType.NAME && activity.data?.data.oneofKind === 'nameChange'">
        <div class="flex space-x-3">
            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                <UIcon class="size-full text-info-600" name="i-mdi-identification-card" />
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
    <template v-else-if="activity.type === UserActivityType.DOCUMENT && activity.data?.data.oneofKind === 'documentRelation'">
        <div class="flex space-x-3">
            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                <UIcon
                    v-if="activity.data.data.documentRelation.added"
                    class="size-full text-info-600"
                    name="i-mdi-file-account"
                />
                <UIcon v-else class="size-full text-base-600" name="i-mdi-file-account-outline" />
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

                        <DocumentInfoPopover :document-id="activity.data.data.documentRelation.documentId" load-on-open>
                            <template #title>
                                <IDCopyBadge
                                    :id="activity.data.data.documentRelation.documentId"
                                    prefix="DOC"
                                    size="xs"
                                    disable-tooltip
                                    variant="link"
                                    :padded="false"
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
                            {{ $t(`enums.documents.DocRelation.${DocRelation[activity.data.data.documentRelation.relation]}`) }}
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
                <UIcon v-if="activity.data.data.wantedChange.wanted" class="size-full text-error-400" name="i-mdi-bell-alert" />
                <UIcon v-else class="size-full text-success-400" name="i-mdi-bell-sleep" />
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
                <UIcon class="size-full text-gray-400" name="i-mdi-briefcase" />
            </div>

            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        {{ $t('components.citizens.CitizenInfoActivityFeedEntry.user_props_job_set') }}
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
            activity.data?.data.oneofKind === 'trafficInfractionPointsChange'
        "
    >
        <div class="flex space-x-3">
            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                <UIcon
                    class="size-full"
                    :class="
                        activity.data.data.trafficInfractionPointsChange.old >
                        activity.data.data.trafficInfractionPointsChange.new
                            ? 'text-gray-400'
                            : 'text-orange-400'
                    "
                    name="i-mdi-traffic-cone"
                />
            </div>

            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        <I18nT keypath="components.citizens.CitizenInfoActivityFeedEntry.traffic_infraction_points.action_text">
                            <template #old>
                                <span class="font-semibold">{{ activity.data.data.trafficInfractionPointsChange.old }}</span>
                            </template>
                            <template #new>
                                <span class="font-semibold">{{ activity.data.data.trafficInfractionPointsChange.new }}</span>
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
                    class="size-full text-amber-400"
                    :class="activity.data.data.mugshotChange.new ? 'text-gray-400' : 'text-amber-400'"
                    name="i-mdi-camera-account"
                />
            </div>

            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        <template v-if="activity.data.data.mugshotChange.new">
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.user_props_mugshot_set') }}
                        </template>
                        <template v-else>
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.user_props_mugshot_removed') }}
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
                <UIcon class="size-full text-amber-200" name="i-mdi-tag" />
            </div>

            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="inline-flex flex-col gap-1 text-sm font-medium">
                        <span>
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.user_props_labels_updated') }}
                        </span>

                        <div class="inline-flex gap-1">
                            <UBadge
                                v-for="attribute in activity.data.data.labelsChange.removed"
                                :key="attribute.name"
                                class="justify-between gap-2 line-through"
                                :class="isColourBright(hexToRgb(attribute.color, RGBBlack)!) ? '!text-black' : '!text-white'"
                                :style="{ backgroundColor: attribute.color }"
                                size="xs"
                            >
                                {{ attribute.name }}
                            </UBadge>

                            <UBadge
                                v-for="attribute in activity.data.data.labelsChange.added"
                                :key="attribute.name"
                                class="justify-between gap-2"
                                :class="isColourBright(hexToRgb(attribute.color, RGBBlack)!) ? '!text-black' : '!text-white'"
                                :style="{ backgroundColor: attribute.color }"
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
                    class="size-full"
                    :class="activity.data.data.licensesChange.added ? 'text-info-600' : 'text-amber-600'"
                    name="i-mdi-license"
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
                <UIcon v-if="activity.data.data.jailChange.seconds > 0" class="size-full" name="i-mdi-handcuffs" />
                <UIcon v-else-if="activity.data.data.jailChange.seconds === 0" class="size-full" name="i-mdi-door-open" />
                <UIcon v-else class="size-full" name="i-mdi-run-fast" />
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
                    class="size-full text-red-400"
                    name="i-mdi-receipt-text-remove"
                />
                <UIcon
                    v-else-if="activity.data.data.fineChange.amount < 0"
                    class="size-full text-success-400"
                    name="i-mdi-receipt-text-check"
                />
                <UIcon v-else class="size-full text-info-400" name="i-mdi-receipt-text-plus" />
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
                <UIcon class="size-full" name="i-mdi-help-circle" />
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
