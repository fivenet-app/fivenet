<script setup lang="ts">
import {
    AtIcon,
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
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                <FileAccountIcon v-if="activity.newValue !== ''" class="text-info-600 h-full w-full" aria-hidden="true" />
                <FileAccountOutlineIcon v-else class="text-base-600 h-full w-full" aria-hidden="true" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="inline-flex items-center gap-1 text-sm font-medium text-neutral">
                        <span
                            ><template v-if="activity.newValue !== ''">
                                {{ $t('components.citizens.citizen_info_activity_feed_entry.document_relation.added') }}
                            </template>
                            <template v-else>
                                {{ $t('components.citizens.citizen_info_activity_feed_entry.document_relation.removed') }}
                            </template>
                        </span>
                        <span class="font-semibold">
                            <NuxtLink
                                :to="{
                                    name: 'documents-id',
                                    params: { id: activity.newValue !== '' ? activity.newValue : activity.oldValue },
                                }"
                            >
                                {{ $t('common.document', 1) }}
                            </NuxtLink>
                        </span>
                        <IDCopyBadge :id="activity.newValue !== '' ? activity.newValue : activity.oldValue" prefix="DOC" />
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span>{{ $t('common.reason') }}:</span>
                        <span class="font-bold">
                            {{ $t(`enums.docstore.DocRelation.${activity.reason.replace('DOC_RELATION_', '')}`) }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
    <template v-else-if="activity.key === 'UserProps.Wanted'">
        <div class="flex space-x-3">
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                <BellAlertIcon v-if="activity.newValue === 'true'" class="text-error-400 h-full w-full" aria-hidden="true" />
                <BellSleepIcon v-else class="text-success-400 h-full w-full" aria-hidden="true" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium text-neutral">
                        {{ $t('components.citizens.citizen_info_activity_feed_entry.bool_set_citizen') }}
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
                        <span>{{ $t('common.reason') }}:</span>
                        <span class="font-semibold">
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
    <template v-else-if="activity.key === 'UserProps.Job'">
        <div class="flex space-x-3">
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                <BriefcaseIcon class="text-secondary-400 h-full w-full" aria-hidden="true" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium text-neutral">
                        {{ $t('components.citizens.citizen_info_activity_feed_entry.userprops_job_set') }}
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
                        <span>{{ $t('common.reason') }}:</span>
                        <span class="font-semibold">
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
    <template v-else-if="activity.key === 'UserProps.TrafficInfractionPoints'">
        <div class="flex space-x-3">
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                <TrafficConeIcon class="text-secondary-400 h-full w-full" aria-hidden="true" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium text-neutral">
                        {{ $t('components.citizens.citizen_info_activity_feed_entry.traffic_infraction_points.action_text') }}
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
                        <span>{{ $t('common.reason') }}:</span>
                        <span class="font-semibold">
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
    <template v-else-if="activity.key === 'UserProps.MugShot'">
        <div class="flex space-x-3">
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                <CameraAccountIcon class="text-secondary-400 h-full w-full" aria-hidden="true" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium text-neutral">
                        {{ $t('components.citizens.citizen_info_activity_feed_entry.userprops_mug_shot_set') }}
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span>{{ $t('common.reason') }}:</span>
                        <span class="font-semibold">
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
    <template v-else-if="activity.key === 'Plugin.Licenses'">
        <div class="flex space-x-3">
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                <LicenseIcon
                    class="h-full w-full"
                    :class="activity.newValue !== '' ? 'text-info-600' : 'text-warn-600'"
                    aria-hidden="true"
                />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium text-neutral">
                        {{
                            activity.newValue !== ''
                                ? $t('components.citizens.citizen_info_activity_feed_entry.plugin_license.added')
                                : $t('components.citizens.citizen_info_activity_feed_entry.plugin_license.removed')
                        }}
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span>{{ $t('common.reason') }}:</span>
                        <span class="font-semibold">
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
    <template v-else-if="activity.key === 'Plugin.Jail'">
        <div class="flex space-x-3">
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full text-neutral">
                <HandcuffsIcon
                    v-if="activity.oldValue === '' && activity.newValue !== '0'"
                    class="h-full w-full"
                    aria-hidden="true"
                />
                <DoorOpenIcon v-else-if="activity.newValue === '0'" class="h-full w-full" aria-hidden="true" />
                <RunFastIcon v-else class="h-full w-full" aria-hidden="true" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium text-neutral">
                        <template v-if="activity.oldValue === '' && activity.newValue !== '0'">
                            {{ $t('components.citizens.citizen_info_activity_feed_entry.plugin_jail.jailed') }}
                            {{ fromSecondsToFormattedDuration(parseInt(props.activity.newValue)) }}
                        </template>
                        <template v-else-if="activity.newValue === '0'">
                            {{ $t('components.citizens.citizen_info_activity_feed_entry.plugin_jail.unjailed') }}
                        </template>
                        <template v-else>
                            {{ $t('components.citizens.citizen_info_activity_feed_entry.plugin_jail.escaped') }}
                        </template>
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span>{{ $t('common.reason') }}:</span>
                        <span class="font-semibold">
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
    <template v-else-if="activity.key === 'Plugin.Billing.Fines'">
        <div class="flex space-x-3">
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                <ReceiptTextCheckIcon
                    v-if="activity.newValue === '0'"
                    class="text-success-400 h-full w-full"
                    aria-hidden="true"
                />
                <ReceiptTextRemoveIcon
                    v-else-if="activity.newValue === activity.oldValue"
                    class="text-secondary-400 h-full w-full"
                    aria-hidden="true"
                />
                <ReceiptTextPlusIcon v-else class="text-info-400 h-full w-full" aria-hidden="true" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium text-neutral">
                        <template v-if="activity.newValue === '0'">
                            {{ $t('components.citizens.citizen_info_activity_feed_entry.plugin_billing_fines.paid') }}
                        </template>
                        <template v-else-if="activity.newValue === activity.oldValue">
                            {{
                                $t('components.citizens.citizen_info_activity_feed_entry.plugin_billing_fines.removed')
                            }}</template
                        >
                        <template v-else>
                            {{ $t('components.citizens.citizen_info_activity_feed_entry.plugin_billing_fines.created') }}
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
                        <span>{{ $t('common.reason') }}:</span>
                        <span class="font-semibold">
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
    <template v-else>
        <div class="flex space-x-3">
            <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
                <HelpCircleIcon class="text-neutral h-full w-full" aria-hidden="true" />
            </div>
            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium text-neutral">
                        {{ `${props.activity.oldValue} -> ${props.activity.newValue}` }}
                    </h3>
                    <p class="text-sm text-gray-400">
                        <GenericTime :value="activity.createdAt" type="long" />
                    </p>
                </div>
                <div class="flex items-center justify-between">
                    <p class="inline-flex gap-1 text-sm text-gray-300">
                        <span>{{ $t('common.reason') }}:</span>
                        <span class="font-semibold">
                            {{ activity.reason }}
                        </span>
                    </p>
                    <p class="inline-flex text-sm text-gray-300">
                        {{ $t('common.created_by') }}
                        <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                    </p>
                </div>
            </div>
        </div>
    </template>
</template>
