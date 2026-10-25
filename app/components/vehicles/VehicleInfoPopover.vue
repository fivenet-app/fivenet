<script lang="ts" setup>
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import GenericTime from '../partials/elements/GenericTime.vue';
import WantedBadge from '~/components/partials/WantedBadge.vue';

const vehicle = defineModel<Vehicle>({ required: true });

const { canViewWanted, canSetWanted, dynamicProperties, openWantedDrawer, openDynamicPropertyDrawer } =
    useVehiclePropertyActions(vehicle);
const { propertyLabel } = usePropertyRenderer();
</script>

<template>
    <UPopover>
        <UButton variant="link" icon="i-mdi-car-settings" color="primary" />

        <template #content>
            <div class="p-4">
                <div class="grid grid-cols-1 gap-2">
                    <UTooltip
                        v-if="canSetWanted"
                        :text="vehicle?.props?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted')"
                    >
                        <UButton
                            variant="link"
                            :color="vehicle?.props?.wanted ? 'error' : 'primary'"
                            :icon="vehicle?.props?.wanted ? 'i-mdi-account-alert' : 'i-mdi-account-cancel'"
                            :label="vehicle?.props?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted')"
                            @click="openWantedDrawer"
                        />
                    </UTooltip>

                    <UButton
                        v-for="property in dynamicProperties"
                        :key="property.key"
                        variant="link"
                        color="neutral"
                        icon="i-mdi-pencil"
                        :label="`${$t(canActOnDynamicProperty(property) ? 'common.edit' : 'common.set')}: ${propertyLabel(property)}`"
                        @click="openDynamicPropertyDrawer(property)"
                    />

                    <p v-if="!canSetWanted && !dynamicProperties.length" class="font-semibold">
                        {{ $t('common.no_actions_available') }}
                    </p>
                </div>

                <ul class="mt-1" role="list">
                    <li v-if="vehicle.props?.updatedAt">
                        <span class="font-semibold">{{ $t('common.last_updated') }}:</span>
                        <GenericTime class="ml-1" :value="vehicle.props?.updatedAt" />
                    </li>

                    <li v-if="canViewWanted && vehicle.props?.wanted" class="inline-flex items-center gap-2">
                        <WantedBadge :wanted="true" variant="short" />

                        <span class="line-clamp-3 font-semibold">{{ $t('common.reason') }}:</span>
                        {{ vehicle.props?.wantedReason ?? $t('common.na') }}
                    </li>
                </ul>
            </div>
        </template>
    </UPopover>
</template>
