<script lang="ts" setup>
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import GenericTime from '../partials/elements/GenericTime.vue';
import VehiclePropsWantedModal from './VehiclePropsWantedModal.vue';

const vehicle = defineModel<Vehicle>({ required: true });

const overlay = useOverlay();

const vehiclePropsWantedModal = overlay.create(VehiclePropsWantedModal);

const { can } = useAuth();
</script>

<template>
    <UPopover>
        <UButton variant="link" icon="i-mdi-car-settings" color="primary" />

        <template #content>
            <div class="p-4">
                <div class="grid grid-cols-1 gap-2">
                    <UTooltip
                        v-if="can('vehicles.VehiclesService/SetVehicleProps').value"
                        :text="vehicle?.props?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted')"
                    >
                        <UButton
                            variant="link"
                            :color="vehicle?.props?.wanted ? 'error' : 'primary'"
                            :icon="vehicle?.props?.wanted ? 'i-mdi-account-alert' : 'i-mdi-account-cancel'"
                            :label="vehicle?.props?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted')"
                            @click="
                                vehiclePropsWantedModal.open({
                                    vehicleProps: vehicle.props,
                                    plate: vehicle.plate,
                                    'onUpdate:vehicleProps': ($event) => {
                                        vehicle.props = $event;
                                    },
                                })
                            "
                        />
                    </UTooltip>

                    <p v-else class="font-semibold">{{ $t('common.no_actions_available') }}</p>
                </div>

                <ul role="list" class="mt-1">
                    <li v-if="vehicle.props?.updatedAt">
                        <span class="font-semibold">{{ $t('common.last_updated') }}:</span>
                        <GenericTime class="ml-1" :value="vehicle.props?.updatedAt" />
                    </li>

                    <li v-if="vehicle.props?.wanted" class="inline-flex items-center gap-2">
                        <UBadge color="error">
                            {{ $t('common.wanted').toUpperCase() }}
                        </UBadge>

                        <span class="line-clamp-3 font-semibold">{{ $t('common.reason') }}:</span>
                        {{ vehicle.props?.wantedReason ?? $t('common.na') }}
                    </li>
                </ul>
            </div>
        </template>
    </UPopover>
</template>
