<script lang="ts" setup>
import LabelBadge from '~/components/citizens/labels/LabelBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DocumentInfoPopover from '~/components/partials/documents/DocumentInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import PropertyActivityChange from '~/components/partials/properties/PropertyActivityChange.vue';
import { DocRelation } from '~~/gen/ts/resources/documents/relations/relations';
import { type UserActivity, UserActivityType } from '~~/gen/ts/resources/users/activity/activity';

defineProps<{
    activity: UserActivity;
}>();

const numberFormatter = useDisplayNumberFormat();
</script>

<template>
    <div class="flex min-w-0 flex-auto justify-between gap-2 text-xs leading-5 text-muted">
        <div class="flex min-w-0 flex-col gap-0.5">
            <span class="text-sm font-medium text-highlighted">
                <template v-if="activity.type === UserActivityType.NAME && activity.data?.data.oneofKind === 'nameChange'">
                    <I18nT keypath="components.citizens.CitizenInfoActivityFeedEntry.name_change">
                        <template #old
                            ><span class="font-semibold">{{ activity.data.data.nameChange.old }}</span></template
                        >
                        <template #new
                            ><span class="font-semibold">{{ activity.data.data.nameChange.new }}</span></template
                        >
                    </I18nT>
                </template>
                <template
                    v-else-if="
                        activity.type === UserActivityType.DOCUMENT && activity.data?.data.oneofKind === 'documentRelation'
                    "
                >
                    {{
                        $t(
                            `components.citizens.CitizenInfoActivityFeedEntry.document_relation.${activity.data.data.documentRelation.added ? 'added' : 'removed'}`,
                        )
                    }}
                    <DocumentInfoPopover :document-id="activity.data.data.documentRelation.documentId" load-on-open>
                        <template #title>
                            <IDCopyBadge
                                :id="activity.data.data.documentRelation.documentId"
                                prefix="DOC"
                                size="xs"
                                disable-tooltip
                                variant="link"
                                hide-icon
                            />
                        </template>
                    </DocumentInfoPopover>
                </template>
                <template
                    v-else-if="activity.type === UserActivityType.WANTED && activity.data?.data.oneofKind === 'wantedChange'"
                >
                    {{ $t('components.citizens.CitizenInfoActivityFeedEntry.bool_set_citizen') }}
                    <span class="font-semibold">{{
                        activity.data.data.wantedChange.wanted
                            ? $t('common.wanted')
                            : `${$t('common.not')} ${$t('common.wanted')}`
                    }}</span>
                    <span
                        v-if="activity.data.data.wantedChange.wanted && activity.data.data.wantedChange.wantedTill"
                        class="inline-flex gap-1"
                    >
                        ({{ $t('common.until') }} <GenericTime :value="activity.data.data.wantedChange.wantedTill" />)
                    </span>
                </template>
                <template v-else-if="activity.type === UserActivityType.JOB && activity.data?.data.oneofKind === 'jobChange'">
                    {{ $t('components.citizens.CitizenInfoActivityFeedEntry.user_props_job_set') }}
                    <span class="font-semibold">
                        {{ activity.data.data.jobChange.jobLabel }}
                        <span v-if="activity.data.data.jobChange.grade">
                            ({{ $t('common.rank') }}: {{ activity.data.data.jobChange.gradeLabel }})
                        </span>
                    </span>
                </template>
                <template
                    v-else-if="
                        activity.type === UserActivityType.TRAFFIC_INFRACTION_POINTS &&
                        activity.data?.data.oneofKind === 'trafficInfractionPointsChange'
                    "
                >
                    <I18nT keypath="components.citizens.CitizenInfoActivityFeedEntry.traffic_infraction_points.action_text">
                        <template #old
                            ><span class="font-semibold">{{
                                activity.data.data.trafficInfractionPointsChange.old
                            }}</span></template
                        >
                        <template #new
                            ><span class="font-semibold">{{
                                activity.data.data.trafficInfractionPointsChange.new
                            }}</span></template
                        >
                    </I18nT>
                </template>
                <template
                    v-else-if="activity.type === UserActivityType.MUGSHOT && activity.data?.data.oneofKind === 'mugshotChange'"
                >
                    {{
                        $t(
                            `components.citizens.CitizenInfoActivityFeedEntry.user_props_mugshot_${activity.data.data.mugshotChange.new ? 'set' : 'removed'}`,
                        )
                    }}
                </template>
                <template v-else-if="activity.type === UserActivityType.FILE && activity.data?.data.oneofKind === 'fileChange'">
                    {{ $t(activity.data.data.fileChange.added ? 'common.image_uploaded' : 'common.image_deleted') }}
                </template>
                <template
                    v-else-if="activity.type === UserActivityType.LABELS && activity.data?.data.oneofKind === 'labelsChange'"
                >
                    {{ $t('components.citizens.CitizenInfoActivityFeedEntry.user_props_labels_updated') }}
                </template>
                <template
                    v-else-if="
                        activity.type === UserActivityType.LICENSES && activity.data?.data.oneofKind === 'licensesChange'
                    "
                >
                    {{
                        activity.data.data.licensesChange.added
                            ? $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.added')
                            : $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.removed')
                    }}:
                    {{ activity.data.data.licensesChange.licenses.map((license) => license.label).join(', ') }}
                </template>
                <template v-else-if="activity.type === UserActivityType.JAIL && activity.data?.data.oneofKind === 'jailChange'">
                    <template v-if="activity.data.data.jailChange.seconds > 0">
                        {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.jailed') }}
                        {{ fromSecondsToFormattedDuration(activity.data.data.jailChange.seconds) }}
                    </template>
                    <template v-else-if="activity.data.data.jailChange.seconds === 0">{{
                        $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.unjailed')
                    }}</template>
                    <template v-else>{{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_jail.escaped') }}</template>
                </template>
                <template v-else-if="activity.type === UserActivityType.FINE && activity.data?.data.oneofKind === 'fineChange'">
                    <template v-if="activity.data.data.fineChange.removed">{{
                        $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.removed')
                    }}</template>
                    <template v-else-if="activity.data.data.fineChange.amount < 0">{{
                        $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.paid')
                    }}</template>
                    <template v-else>{{
                        $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.created')
                    }}</template>
                    {{ numberFormatter.format(Math.abs(activity.data.data.fineChange.amount)) }}
                </template>
                <PropertyActivityChange
                    v-else-if="
                        activity.type === UserActivityType.PROPERTY && activity.data?.data.oneofKind === 'propertyChange'
                    "
                    :property-key="activity.data.data.propertyChange.key"
                    :old-value="activity.data.data.propertyChange.oldValue"
                    :new-value="activity.data.data.propertyChange.newValue"
                />
                <template v-else>{{ UserActivityType[activity.type] }}</template>
            </span>

            <div
                v-if="activity.type === UserActivityType.LABELS && activity.data?.data.oneofKind === 'labelsChange'"
                class="flex flex-wrap gap-1"
            >
                <LabelBadge
                    v-for="label in activity.data.data.labelsChange.removed"
                    :key="label.id"
                    class="line-through"
                    :label="label"
                />
                <LabelBadge v-for="label in activity.data.data.labelsChange.added" :key="label.id" :label="label" />
                <LabelBadge
                    v-for="label in activity.data.data.labelsChange.removedIds"
                    :id="label"
                    :key="label"
                    :expired="activity.data.data.labelsChange.expired"
                />
                <LabelBadge
                    v-for="label in activity.data.data.labelsChange.addedIds"
                    :id="label.id"
                    :key="label.id"
                    :expires-at="label.expiresAt"
                />
            </div>

            <p
                v-if="
                    !(
                        activity.type === UserActivityType.WANTED &&
                        activity.data?.data.oneofKind === 'wantedChange' &&
                        activity.data.data.wantedChange.auto
                    )
                "
                class="inline-flex min-w-0 gap-1 text-sm break-words"
            >
                <template
                    v-if="activity.type === UserActivityType.DOCUMENT && activity.data?.data.oneofKind === 'documentRelation'"
                >
                    <span class="font-semibold">{{ $t('common.type') }}:</span>
                    <span>{{
                        $t(`enums.documents.DocRelation.${DocRelation[activity.data.data.documentRelation.relation]}`)
                    }}</span>
                </template>
                <template
                    v-else-if="
                        activity.type !== UserActivityType.JAIL ||
                        (activity.data?.data.oneofKind === 'jailChange' && activity.data.data.jailChange.seconds >= 0)
                    "
                >
                    <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                    <span
                        v-if="
                            activity.type === UserActivityType.LABELS &&
                            activity.data?.data.oneofKind === 'labelsChange' &&
                            activity.data.data.labelsChange.expired
                        "
                        >{{ $t('common.expired') }}</span
                    >
                    <!-- eslint-disable-next-line vue/no-v-html -->
                    <span v-else v-html="activity.reason" />
                </template>
            </p>
        </div>

        <div class="flex flex-none flex-col items-end gap-0.5 text-xs">
            <span class="text-xs/5 text-dimmed"><GenericTime :value="activity.createdAt" type="short" /></span>
            <span
                v-if="
                    activity.type === UserActivityType.WANTED &&
                    activity.data?.data.oneofKind === 'wantedChange' &&
                    activity.data.data.wantedChange.auto
                "
                class="inline-flex gap-1 text-sm"
            >
                <span class="font-semibold">{{ $t('components.citizens.CitizenInfoActivityFeedEntry.automatic') }}</span>
            </span>
            <span v-if="activity.sourceUser" class="inline-flex items-center gap-1">
                <span class="text-dimmed">{{ $t('common.created_by') }}:</span>
                <CitizenInfoPopover :user="activity.sourceUser" />
            </span>
        </div>
    </div>
</template>
