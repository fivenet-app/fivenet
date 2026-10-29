import type { Ref } from 'vue';
import SetDynamicPropertyDrawer from '~/components/partials/properties/SetDynamicPropertyDrawer.vue';
import SetWantedDrawer from '~/components/vehicles/SetWantedDrawer.vue';
import { dynamicPropertiesForAction } from '~/composables/useDynamicPropertyAction';
import { PropertyScope, type Property } from '~~/gen/ts/resources/properties/properties';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';

export function useVehiclePropertyActions(vehicle: Ref<Vehicle>, onUpdated?: () => void) {
    const overlay = useOverlay();
    const setWantedDrawer = overlay.create(SetWantedDrawer);
    const setDynamicPropertyDrawer = overlay.create(SetDynamicPropertyDrawer);

    const { attr, can } = useAuth();
    const canManageVehicleProperties = computed(() => can('internal.Superuser/JobAdmin').value);

    const canViewWanted = computed(() => attr('vehicles.VehiclesService/ListVehicles', 'Properties', 'vehicle.wanted').value);
    const canSetWanted = computed(
        () =>
            canViewWanted.value &&
            can('vehicles.VehiclesService/SetVehicleProps').value &&
            attr('vehicles.VehiclesService/SetVehicleProps', 'Properties', 'vehicle.wanted').value,
    );
    const dynamicProperties = computed(() =>
        dynamicPropertiesForAction(
            vehicle.value.properties?.properties,
            (property) =>
                canManageVehicleProperties.value ||
                attr('vehicles.VehiclesService/SetVehicleProps', 'Properties', property.key).value,
        ),
    );

    function openWantedDrawer(): void {
        setWantedDrawer.open({
            vehicleProps: vehicle.value.props,
            plate: vehicle.value.plate,
            'onUpdate:vehicleProps': (value) => {
                vehicle.value.props = value;
                onUpdated?.();
            },
        });
    }

    function openDynamicPropertyDrawer(property: Property): void {
        setDynamicPropertyDrawer.open({
            property,
            scope: PropertyScope.VEHICLE,
            plate: vehicle.value.plate,
            'onUpdate:properties': (value) => {
                vehicle.value.properties = value;
                onUpdated?.();
            },
        });
    }

    return {
        canViewWanted,
        canSetWanted,
        dynamicProperties,
        openWantedDrawer,
        openDynamicPropertyDrawer,
    };
}
