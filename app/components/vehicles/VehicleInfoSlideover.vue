<script lang="ts" setup>
import type { TabsItem } from '@nuxt/ui';
import { titleCase } from 'scule';
import ActivityFeed from '~/components/vehicles/ActivityFeed.vue';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';
import GenericTime from '../partials/elements/GenericTime.vue';
import LicensePlate from '../partials/LicensePlate.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import WantedBadge from '~/components/partials/WantedBadge.vue';
import VehicleFiles from './files/VehicleFiles.vue';

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const vehicle = defineModel<Vehicle>({ required: true });

const { t } = useI18n();
const { attr } = useAuth();

const canViewFiles = attr('vehicles.VehiclesService/ListVehicles', 'Properties', 'vehicle.files');

const activityFeed = useTemplateRef<InstanceType<typeof ActivityFeed>>('activityFeed');

const { canViewWanted, canSetWanted, dynamicProperties, openWantedDrawer, openDynamicPropertyDrawer } =
    useVehiclePropertyActions(vehicle, () => activityFeed.value?.refresh());

const { propertyLabel } = usePropertyRenderer();

const selectedTab = ref('activity');

const tabs = computed<TabsItem[]>(() => [
    {
        slot: 'activity',
        label: t('common.activity'),
        icon: 'i-mdi-history',
        value: 'activity',
    },
    ...(canViewFiles.value
        ? [
              {
                  slot: 'files',
                  label: t('common.image', 2),
                  icon: 'i-mdi-folder-multiple-image',
                  value: 'files',
              },
          ]
        : []),
]);
</script>

<template>
    <USlideover
        :title="$t('components.vehicles.VehicleInfoSlideover.title')"
        :overlay="false"
        :ui="{ body: 'flex flex-col p-0 sm:p-0 overflow-y-hidden', content: 'w-full max-w-6xl sm:max-w-5xl' }"
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

                            <div v-if="canSetWanted || dynamicProperties.length" class="space-y-1">
                                <p class="text-sm font-semibold">{{ $t('common.action', 2) }}</p>
                                <div class="flex flex-wrap gap-2">
                                    <UButton
                                        v-if="canSetWanted"
                                        :color="vehicle.props?.wanted ? 'error' : 'primary'"
                                        :icon="vehicle.props?.wanted ? 'i-mdi-account-alert' : 'i-mdi-account-cancel'"
                                        :label="vehicle.props?.wanted ? $t('common.revoke_wanted') : $t('common.set_wanted')"
                                        @click="() => openWantedDrawer()"
                                    />

                                    <UButton
                                        v-for="property in dynamicProperties"
                                        :key="property.key"
                                        color="neutral"
                                        variant="outline"
                                        icon="i-mdi-pencil"
                                        :label="`${$t(canActOnDynamicProperty(property) ? 'common.edit' : 'common.set')}: ${propertyLabel(property)}`"
                                        @click="openDynamicPropertyDrawer(property)"
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
                                    <WantedBadge :wanted="vehicle.props?.wanted ?? false" variant="long" />
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

                <div class="min-h-0 flex-1">
                    <UTabs
                        v-model="selectedTab"
                        class="flex h-full min-h-0 flex-col"
                        :items="tabs"
                        :unmount-on-hide="true"
                        variant="link"
                        :ui="{ root: 'flex min-h-0 flex-1 flex-col gap-0', content: 'min-h-0 flex-1 overflow-y-auto' }"
                    >
                        <template #activity>
                            <div class="min-h-full">
                                <ActivityFeed ref="activityFeed" :plate="vehicle.plate" :owner-id="vehicle.ownerId" />
                            </div>
                        </template>

                        <template #files>
                            <div class="min-h-full">
                                <VehicleFiles v-model="vehicle" />
                            </div>
                        </template>
                    </UTabs>
                </div>
            </div>
        </template>

        <template #footer>
            <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />
        </template>
    </USlideover>
</template>
