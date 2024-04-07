<script setup lang="ts">
import {
    BellAlertIcon,
    BellSleepIcon,
    BriefcaseIcon,
    CameraAccountIcon,
    DoorOpenIcon,
    FileAccountIcon,
    FileAccountOutlineIcon,
    HandcuffsIcon,
    HelpCircleIcon,
    LicenseIcon,
    ReceiptTextCheckIcon,
    ReceiptTextPlusIcon,
    ReceiptTextRemoveIcon,
    RunFastIcon,
    TrafficConeIcon,
} from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { UserActivity } from '~~/gen/ts/resources/users/users';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';

const props = defineProps<{
    activity: UserActivity;
}>();
</script>

<template>
    <template v-if="activity.key === 'DocStore.Relation'">
        <div class="flex space-x-3">
            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                <FileAccountIcon v-if="activity.newValue !== ''" class="size-full text-info-600" />
                <FileAccountOutlineIcon v-else class="size-full text-base-600" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="inline-flex items-center gap-1 text-sm font-medium">
                        <template v-if="activity.newValue !== ''">
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.document_relation.added') }}
                        </template>
                        <template v-else>
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.document_relation.removed') }}
                        </template>

                        <UButton
                            variant="link"
                            :padded="false"
                            :to="{
                                name: 'documents-id',
                                params: { id: activity.newValue !== '' ? activity.newValue : activity.oldValue },
                            }"
                        >
                            {{ $t('common.document', 1) }}
                            <IDCopyBadge :id="activity.newValue !== '' ? activity.newValue : activity.oldValue" prefix="DOC" />
                        </UButton>
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>
                            {{ $t(`enums.docstore.DocRelation.${activity.reason.replace('DOC_RELATION_', '')}`) }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
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
                <BellAlertIcon v-if="activity.newValue === 'true'" class="size-full text-error-400" />
                <BellSleepIcon v-else class="size-full text-success-400" />
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
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
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
                <BriefcaseIcon class="text-secondary-400 size-full" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        {{ $t('components.citizens.CitizenInfoActivityFeedEntry.userprops_job_set') }}
                        <span class="font-semibold">
                            {{ activity.newValue }}
                        </span>
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
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
                <TrafficConeIcon class="text-secondary-400 size-full" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        {{ $t('components.citizens.CitizenInfoActivityFeedEntry.traffic_infraction_points.action_text') }}
                        <span>
                            <span class="font-semibold">{{ activity.oldValue }}</span>
                            {{ $t('common.to').toLocaleLowerCase() }}
                            <span class="font-semibold">{{ activity.newValue }}</span>
                        </span>
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
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
                <CameraAccountIcon class="text-secondary-400 size-full" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        {{ $t('components.citizens.CitizenInfoActivityFeedEntry.userprops_mug_shot_set') }}
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
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
                <LicenseIcon class="size-full" :class="activity.newValue !== '' ? 'text-info-600' : 'text-warn-600'" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        {{
                            activity.newValue !== ''
                                ? $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.added')
                                : $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_license.removed')
                        }}
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
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
                <HandcuffsIcon v-if="activity.oldValue === '' && activity.newValue !== '0'" class="size-full" />
                <DoorOpenIcon v-else-if="activity.newValue === '0'" class="size-full" />
                <RunFastIcon v-else class="size-full" />
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
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <template v-if="activity.oldValue === '' && activity.newValue !== '0'">
                            <span class="font-semibold">{{ $t('common.reason') }}:</span>
                            <span>
                                {{ !activity.reason ? $t('common.na') : activity.reason }}
                            </span>
                        </template>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
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
                <ReceiptTextCheckIcon v-if="activity.newValue === '0'" class="size-full text-success-400" />
                <ReceiptTextRemoveIcon
                    v-else-if="activity.newValue === activity.oldValue"
                    class="text-secondary-400 size-full"
                />
                <ReceiptTextPlusIcon v-else class="size-full text-info-400" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        <template v-if="activity.newValue === '0'">
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.paid') }}
                        </template>
                        <template v-else-if="activity.newValue === activity.oldValue">
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.removed') }}</template
                        >
                        <template v-else>
                            {{ $t('components.citizens.CitizenInfoActivityFeedEntry.plugin_billing_fines.created') }}
                        </template>
                        <span>
                            {{ $n(parseInt(props.activity.newValue), 'currency') }}
                        </span>
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
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
                <HelpCircleIcon class="size-full" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        {{ `${props.activity.oldValue} -> ${props.activity.newValue}` }}
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
</template>
