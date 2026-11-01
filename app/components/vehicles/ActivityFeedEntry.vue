<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import PropertyActivityChange from '~/components/partials/properties/PropertyActivityChange.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { VehicleActivityType, type VehicleActivity } from '~~/gen/ts/resources/vehicles/activity/activity';

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
    <div class="flex min-w-0 flex-auto justify-between gap-2 text-xs leading-5 text-muted">
        <div class="flex min-w-0 flex-col gap-0.5">
            <span class="text-sm font-medium text-highlighted">
                <template
                    v-if="
                        activity.activityType === VehicleActivityType.WANTED && activity.data?.data.oneofKind === 'wantedChange'
                    "
                >
                    {{ $t('components.vehicles.VehicleActivityFeedEntry.wanted_set') }}
                    <span class="font-semibold">
                        {{
                            activity.data.data.wantedChange.wanted
                                ? $t('common.wanted')
                                : `${$t('common.not')} ${$t('common.wanted')}`
                        }}
                    </span>
                </template>

                <PropertyActivityChange
                    v-else-if="
                        activity.activityType === VehicleActivityType.PROPERTY &&
                        activity.data?.data.oneofKind === 'propertyChange'
                    "
                    :property-key="activity.data.data.propertyChange.key"
                    :old-value="activity.data.data.propertyChange.oldValue"
                    :new-value="activity.data.data.propertyChange.newValue"
                />
                <template
                    v-else-if="
                        activity.activityType === VehicleActivityType.FILE && activity.data?.data.oneofKind === 'fileChange'
                    "
                >
                    {{ $t(activity.data.data.fileChange.added ? 'common.image_uploaded' : 'common.image_deleted') }}
                </template>
                <template v-else>{{
                    $t(`enums.vehicles.VehicleActivityType.${VehicleActivityType[activity.activityType]}`)
                }}</template>
            </span>

            <div
                v-if="activity.activityType === VehicleActivityType.WANTED && activity.data?.data.oneofKind === 'wantedChange'"
                class="grid gap-0.5 text-sm md:grid-cols-2"
            >
                <p v-if="!activity.data.data.wantedChange.auto" class="inline-flex min-w-0 gap-1">
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
            </div>
        </div>

        <div class="flex flex-none flex-col items-end gap-0.5 text-xs">
            <span class="text-xs/5 text-dimmed"><GenericTime :value="activity.createdAt" type="short" /></span>

            <span v-if="activity.creator" class="inline-flex items-center gap-1">
                <span class="text-dimmed">{{ $t('common.created_by') }}:</span>
                <CitizenInfoPopover :user="activity.creator" />
            </span>
        </div>
    </div>
</template>
