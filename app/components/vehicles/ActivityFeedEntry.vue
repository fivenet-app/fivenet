<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ActivityFeedItem from '~/components/partials/data/ActivityFeedItem.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { VehicleActivityType, type VehicleActivity } from '~~/gen/ts/resources/vehicles/activity/activity';
import { vehicleActivityIconColor, vehicleActivityTypeIcon } from './helpers';

const props = defineProps<{
    activity: VehicleActivity;
}>();

const { t } = useI18n();

const reasonHtml = computed(() => {
    if (props.activity.data?.data.oneofKind !== 'wantedChange') {
        return props.activity.reason ?? t('common.na');
    }

    return props.activity.reason ?? props.activity.data.data.wantedChange.wantedReason ?? t('common.na');
});
</script>

<template>
    <template v-if="activity.activityType === VehicleActivityType.WANTED && activity.data?.data.oneofKind === 'wantedChange'">
        <ActivityFeedItem
            :icon="vehicleActivityTypeIcon(activity.activityType)"
            :icon-class="vehicleActivityIconColor(activity)"
        >
            <template #title>
                {{ $t('components.vehicles.VehicleActivityFeedEntry.wanted_set') }}
                <span class="font-semibold">
                    {{
                        activity.data.data.wantedChange.wanted
                            ? $t('common.wanted')
                            : `${$t('common.not')} ${$t('common.wanted')}`
                    }}
                </span>
            </template>

            <template #timestamp>
                <GenericTime :value="activity.createdAt" type="long" />
            </template>

            <div class="grid gap-2 text-sm md:grid-cols-2">
                <p class="inline-flex min-w-0 gap-1">
                    <span class="font-semibold">{{ $t('common.reason', 1) }}:</span>
                    <!-- Reason text is sanitized by the backend and may contain HTML entities. -->
                    <!-- eslint-disable-next-line vue/no-v-html -->
                    <span class="truncate" v-html="reasonHtml" />
                </p>

                <p v-if="activity.data.data.wantedChange.wantedTill" class="inline-flex gap-1">
                    <span class="font-semibold">{{ $t('common.expiration') }}:</span>
                    <GenericTime :value="activity.data.data.wantedChange.wantedTill" type="long" />
                </p>

                <p v-if="activity.data.data.wantedChange.auto" class="inline-flex gap-1">
                    <span class="font-semibold">{{ $t('components.vehicles.VehicleActivityFeedEntry.automatic') }}</span>
                </p>

                <p v-if="activity.creator" class="inline-flex min-w-0 justify-end text-sm md:justify-self-end">
                    {{ $t('common.created_by') }}
                    <CitizenInfoPopover class="ml-1" :user="activity.creator" />
                </p>
            </div>
        </ActivityFeedItem>
    </template>

    <template v-else>
        <ActivityFeedItem :icon="vehicleActivityTypeIcon(activity.activityType)">
            <template #title>
                {{ $t(`enums.vehicles.VehicleActivityType.${VehicleActivityType[activity.activityType]}`) }}
            </template>

            <template #timestamp>
                <GenericTime :value="activity.createdAt" type="long" />
            </template>
        </ActivityFeedItem>
    </template>
</template>
