<script lang="ts" setup>
import { titleCase } from 'scule';
import ActivityFeed from '~/components/vehicles/ActivityFeed.vue';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import GenericTime from '../partials/elements/GenericTime.vue';
import LicensePlate from '../partials/LicensePlate.vue';
import SetWantedModal from './SetWantedModal.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const vehicle = defineModel<Vehicle>({ required: true });

const overlay = useOverlay();
const setWantedModal = overlay.create(SetWantedModal);

const { attr, can } = useAuth();

const activityFeed = useTemplateRef<InstanceType<typeof ActivityFeed>>('activityFeed');

const canViewWanted = computed(() => attr('vehicles.VehiclesService/ListVehicles', 'Fields', 'Wanted').value);
const canSetWanted = computed(
    () =>
        canViewWanted.value &&
        can('vehicles.VehiclesService/SetVehicleProps').value &&
        attr('vehicles.VehiclesService/SetVehicleProps', 'Fields', 'Wanted').value,
);
</script>

<template>
    <USlideover
        :title="$t('components.vehicles.VehicleInfoSlideover.title')"
        :overlay="false"
        :ui="{ body: 'flex flex-col p-0 sm:p-0 overflow-y-hidden', content: 'max-w-3xl' }"
    >
        <template #body>
            <div class="flex h-full flex-1 flex-col">
                <div class="flex flex-col gap-4 px-4 py-4 sm:px-6">
                    <div class="space-y-3">
                        <div class="flex w-full flex-col gap-3 sm:grid sm:grid-cols-2">
                            <div class="space-y-2 sm:min-w-40 md:min-w-48">
                                <LicensePlate :plate="vehicle.plate" />
                                <div class="text-sm text-muted">
                                    <span v-if="vehicle.model"
                                        ><span class="font-semibold">{{ $t('common.model') }}:</span> {{ vehicle.model }}</span
                                    >
                                    <span v-if="vehicle.model && vehicle.type"> · </span>
                                    <span v-if="vehicle.type"
                                        ><span class="font-semibold">{{ $t('common.type') }}:</span>
                                        {{ titleCase(vehicle.type) }}</span
                                    >
                                </div>
                            </div>

                            <div v-if="canSetWanted" class="space-y-1">
                                <p class="text-sm font-semibold">{{ $t('common.action', 2) }}</p>
                                <div>
                                    <UButton
                                        :color="vehicle.props?.wanted ? 'error' : 'primary'"
                                        :icon="vehicle.props?.wanted ? 'i-mdi-account-alert' : 'i-mdi-account-cancel'"
                                        :label="vehicle.props?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted')"
                                        @click="
                                            () =>
                                                setWantedModal.open({
                                                    vehicleProps: vehicle.props,
                                                    plate: vehicle.plate,
                                                    'onUpdate:vehicleProps': ($event) => {
                                                        vehicle.props = $event;
                                                        activityFeed?.refresh();
                                                    },
                                                })
                                        "
                                    />
                                </div>
                            </div>
                        </div>

                        <dl class="grid gap-3 text-sm sm:grid-cols-2">
                            <div v-if="vehicle.owner">
                                <dt class="font-semibold">{{ $t('common.owner') }}</dt>
                                <dd>
                                    <CitizenInfoPopover :user="vehicle.owner" :owner-id="vehicle.ownerId" />
                                </dd>
                            </div>

                            <div v-if="vehicle.jobLabel || vehicle.job">
                                <dt class="font-semibold">{{ $t('common.job') }}</dt>
                                <dd>{{ vehicle.jobLabel ?? vehicle.job }}</dd>
                            </div>

                            <div v-if="vehicle.props?.updatedAt">
                                <dt class="font-semibold">{{ $t('common.last_updated') }}</dt>
                                <dd><GenericTime :value="vehicle.props.updatedAt" /></dd>
                            </div>

                            <div v-if="canViewWanted">
                                <dt class="font-semibold">{{ $t('common.wanted') }}</dt>
                                <dd>
                                    <UBadge
                                        :color="vehicle.props?.wanted ? 'error' : 'neutral'"
                                        :label="
                                            vehicle.props?.wanted
                                                ? $t('common.wanted').toUpperCase()
                                                : `${$t('common.not')} ${$t('common.wanted')}`
                                        "
                                    />
                                </dd>
                            </div>

                            <div v-if="canViewWanted && vehicle.props?.wanted">
                                <dt class="font-semibold">{{ $t('common.reason') }}</dt>
                                <dd>{{ vehicle.props.wantedReason ?? $t('common.na') }}</dd>
                            </div>

                            <div v-if="canViewWanted && vehicle.props?.wantedTill">
                                <dt class="font-semibold">{{ $t('common.expiration') }}</dt>
                                <dd><GenericTime :value="vehicle.props.wantedTill" /></dd>
                            </div>
                        </dl>
                    </div>
                </div>

                <USeparator />

                <div class="min-h-0 flex-1">
                    <ActivityFeed ref="activityFeed" :plate="vehicle.plate" :owner-id="vehicle.ownerId" />
                </div>
            </div>
        </template>

        <template #footer>
            <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
        </template>
    </USlideover>
</template>
