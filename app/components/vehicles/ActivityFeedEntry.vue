<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import LicensePlate from '~/components/partials/LicensePlate.vue';
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
    <li class="px-4 py-4 sm:px-6">
        <template
            v-if="activity.activityType === VehicleActivityType.WANTED && activity.data?.data.oneofKind === 'wantedChange'"
        >
            <div class="flex gap-3">
                <div class="my-auto flex size-10 shrink-0 items-center justify-center rounded-full">
                    <UIcon
                        :class="[vehicleActivityIconColor(activity), 'size-full']"
                        :name="vehicleActivityTypeIcon(activity.activityType)"
                    />
                </div>

                <div class="min-w-0 flex-1 space-y-2">
                    <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                        <div class="space-y-1">
                            <h3 class="text-sm font-medium">
                                {{ $t('components.vehicles.VehicleActivityFeedEntry.wanted_set') }}
                                <span class="font-semibold">
                                    {{
                                        activity.data.data.wantedChange.wanted
                                            ? $t('common.wanted')
                                            : `${$t('common.not')} ${$t('common.wanted')}`
                                    }}
                                </span>
                            </h3>

                            <LicensePlate :plate="activity.plate" />
                        </div>

                        <p class="shrink-0 text-sm text-dimmed">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>

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
                            <span class="font-semibold">{{
                                $t('components.vehicles.VehicleActivityFeedEntry.automatic')
                            }}</span>
                        </p>

                        <p v-if="activity.creator" class="inline-flex min-w-0 text-sm">
                            {{ $t('common.created_by') }}
                            <CitizenInfoPopover class="ml-1" :user="activity.creator" />
                        </p>
                    </div>
                </div>
            </div>
        </template>

        <template v-else>
            <div class="flex gap-3">
                <div class="my-auto flex size-10 shrink-0 items-center justify-center rounded-full">
                    <UIcon class="size-full" :name="vehicleActivityTypeIcon(activity.activityType)" />
                </div>

                <div class="min-w-0 flex-1">
                    <div class="flex items-center justify-between gap-3">
                        <h3 class="text-sm font-medium">
                            {{ VehicleActivityType[activity.activityType] }}
                        </h3>

                        <p class="shrink-0 text-sm text-dimmed">
                            <GenericTime :value="activity.createdAt" type="long" />
                        </p>
                    </div>
                </div>
            </div>
        </template>
    </li>
</template>
