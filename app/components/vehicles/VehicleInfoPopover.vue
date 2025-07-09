<script lang="ts" setup>
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import GenericTime from '../partials/elements/GenericTime.vue';
import VehiclePropsWantedModal from './VehiclePropsWantedModal.vue';

const vehicle = defineModel<Vehicle>({ required: true });

const modal = useModal();

const { can } = useAuth();
</script>

<template>
    <UPopover>
        <UButton variant="link" icon="i-mdi-car-settings" color="primary" />

        <template #panel>
            <div class="p-4">
                <div class="grid grid-cols-2 gap-2">
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
                                modal.open(VehiclePropsWantedModal, {
                                    vehicleProps: vehicle.props,
                                    plate: vehicle.plate,
                                    'onUpdate:vehicleProps': ($event) => {
                                        vehicle.props = $event;
                                    },
                                })
                            "
                        />
                    </UTooltip>
                </div>

                <ul role="list">
                    <li v-if="vehicle.props?.updatedAt">
                        <span class="font-semibold">{{ $t('common.last_updated') }}:</span>
                        <GenericTime class="ml-1" :value="vehicle.props?.updatedAt" />
                    </li>

                    <li v-if="vehicle.props?.wanted">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        {{ vehicle.props?.wantedReason ?? $t('common.na') }}
                    </li>
                </ul>
            </div>
        </template>
    </UPopover>
</template>
