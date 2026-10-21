<script lang="ts" setup>
import type { Property } from '~~/gen/ts/resources/properties/properties';
import type { User } from '~~/gen/ts/resources/users/user';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { pinnedProperties } from '~/composables/usePropertyRenderer';

const props = defineProps<{ properties?: Property[]; user?: User }>();
const { t } = useI18n();
const { propertyLabel, propertyScalarValue, renderProperty } = usePropertyRenderer();

const pinned = computed(() => pinnedProperties(props.properties));

const statCardClass = computed(() => {
    if (pinned.value.length === 1) return 'xl:col-span-4';
    if (pinned.value.length === 2) return 'xl:col-span-2';
    return 'xl:col-span-1';
});

function isWanted(property: Property): boolean {
    return property.key === 'user.wanted' && Boolean(propertyScalarValue(property));
}

function valueFor(property: Property): string {
    if (property.key === 'user.job' && props.user?.jobLabel) return props.user.jobLabel;
    if (property.key === 'user.traffic_infraction_points' && props.user?.props?.trafficInfractionPoints !== undefined) {
        return t('common.point', props.user.props.trafficInfractionPoints);
    }
    return renderProperty(property);
}
</script>

<template>
    <div v-if="pinned.length" class="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
        <UCard
            v-for="property in pinned"
            :key="property.key"
            :class="[statCardClass, property.key === 'user.wanted' && isWanted(property) ? 'border-error-500/50' : '']"
            :ui="{ body: 'p-3 sm:p-4' }"
        >
            <div class="flex items-center gap-3">
                <UIcon
                    v-if="property.key === 'user.wanted'"
                    class="size-5 shrink-0"
                    :class="isWanted(property) ? 'text-error-500' : 'text-success-500'"
                    :name="isWanted(property) ? 'i-mdi-account-alert' : 'i-mdi-account-check'"
                />
                <UIcon
                    v-else-if="property.key === 'user.traffic_infraction_points'"
                    class="size-5 shrink-0 text-muted"
                    name="i-mdi-counter"
                />
                <UIcon
                    v-else-if="property.key === 'user.open_fines'"
                    class="size-5 shrink-0 text-muted"
                    name="i-mdi-cash-remove"
                />
                <UIcon v-else-if="property.key === 'user.job'" class="size-5 shrink-0 text-muted" name="i-mdi-briefcase" />
                <UIcon v-else class="size-5 shrink-0 text-muted" name="i-mdi-tag-outline" />
                <div class="min-w-0">
                    <p class="text-xs font-medium text-muted">{{ propertyLabel(property) }}</p>
                    <p
                        class="truncate font-semibold text-highlighted"
                        :class="
                            property.key === 'user.traffic_infraction_points' && Number(propertyScalarValue(property)) >= 10
                                ? 'text-error-500'
                                : property.key === 'user.open_fines' && Number(propertyScalarValue(property)) > 0
                                  ? 'text-error-500'
                                  : ''
                        "
                    >
                        <PropertyRender
                            v-if="['user.email', 'user.phone_number', 'user.sex', 'user.height'].includes(property.key)"
                            :property="property"
                        />
                        <template v-else>{{ valueFor(property) }}</template>
                    </p>
                    <p
                        v-if="property.key === 'user.wanted' && isWanted(property) && props.user?.props?.wantedTill"
                        class="text-xs text-muted"
                    >
                        {{ $t('common.until') }}
                        <GenericTime :value="props.user.props.wantedTill" />
                    </p>
                </div>
            </div>
        </UCard>
    </div>
</template>
