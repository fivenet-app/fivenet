<script lang="ts" setup>
import { PropertyDisplayMode, type Property } from '~~/gen/ts/resources/properties/properties';
import type { User } from '~~/gen/ts/resources/users/user';
import WantedBadge from '~/components/partials/WantedBadge.vue';
import { propertyHasDisplayValue, propertyHasValue } from '~/composables/usePropertyRenderer';

const props = defineProps<{ properties?: Property[]; user?: User; keyOrder?: string[] }>();
const { t } = useI18n();
const { game } = useAppConfig();
const { propertyLabel, propertyScalarValue, renderProperty } = usePropertyRenderer();

const badges = computed(() => {
    const badges = (props.properties ?? []).filter(
        (property) =>
            property.definition?.displayMode === PropertyDisplayMode.PROFILE_BADGE && propertyHasDisplayValue(property),
    );
    if (!props.keyOrder?.length) return badges;

    const positions = new Map(props.keyOrder.map((key, index) => [key, index]));
    return badges.sort(
        (a, b) => (positions.get(a.key) ?? Number.MAX_SAFE_INTEGER) - (positions.get(b.key) ?? Number.MAX_SAFE_INTEGER),
    );
});

function isWanted(property: Property): boolean {
    return property.key === 'user.wanted' && Boolean(propertyScalarValue(property));
}

function hasJobPropertyOverride(property: Property): boolean {
    return property.key === 'user.job' && propertyHasValue(property);
}

function valueFor(property: Property): string {
    if (property.key === 'user.wanted') return isWanted(property) ? t('common.yes') : t('common.no');

    if (property.key === 'user.job') {
        const job = props.user?.props?.job?.label || props.user?.jobLabel || String(propertyScalarValue(property) ?? '');
        const jobName = props.user?.props?.jobName ?? props.user?.job;
        const grade =
            props.user?.props?.jobGrade?.label ||
            (props.user?.props?.jobGradeNumber !== undefined ? String(props.user.props.jobGradeNumber) : '') ||
            props.user?.jobGradeLabel ||
            (props.user?.jobGrade !== undefined ? String(props.user.jobGrade) : '');
        return grade && jobName !== game.unemployedJobName ? `${job} (${grade})` : job;
    }

    return renderProperty(property);
}
</script>

<template>
    <div v-if="badges.length" class="flex flex-wrap items-center gap-2">
        <template v-for="property in badges" :key="property.key">
            <WantedBadge v-if="property.key === 'user.wanted'" :wanted="isWanted(property)" variant="long" size="lg" />
            <UBadge v-else color="neutral" variant="subtle" size="lg">
                <UIcon v-if="property.key === 'user.job'" name="i-mdi-briefcase" />
                <UIcon v-else name="i-mdi-tag-outline" />
                <span class="font-medium">{{ propertyLabel(property) }}:</span>
                <span>{{ valueFor(property) }}</span>
                <span v-if="hasJobPropertyOverride(property)">*</span>
            </UBadge>
        </template>
    </div>
</template>
