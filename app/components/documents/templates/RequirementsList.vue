<script lang="ts" setup>
import { computed } from 'vue';
import type { ObjectSpecs } from '~~/gen/ts/resources/documents/templates/templates';

const props = defineProps<{
    name: string;
    plural?: string;
    specs: ObjectSpecs;
    fulfilled?: boolean;
}>();

const displayName = computed(() => props.plural ?? `${props.name}(s)`);

const isRequired = computed(() => props.specs.required);
const hasMin = computed(() => props.specs.min && props.specs.min > 0);
const hasMax = computed(() => props.specs.max && props.specs.max > 0);
const minEqualsMax = computed(() => hasMin.value && hasMax.value && props.specs.max === props.specs.min);
const showRequirement = computed(() => isRequired.value || hasMin.value || hasMax.value);
const badgeColor = computed(() => {
    if (props.fulfilled === true) return 'success';
    if (props.fulfilled === false) return isRequired.value ? 'error' : 'warning';
    return isRequired.value ? 'primary' : 'neutral';
});
const badgeVariant = computed(() => {
    if (props.fulfilled === false && isRequired.value) return 'solid';
    return 'soft';
});

const requirementLabel = computed(() => {
    const count = hasMin.value ? props.specs.min! : isRequired.value ? 1 : props.specs.max!;
    const kind = count === 1 ? props.name : displayName.value;

    let label: string;
    if (minEqualsMax.value) {
        label = isRequired.value
            ? $t('common.requirement.exact_required', { count, kind })
            : $t('common.requirement.exact', { count, kind });
    } else if (isRequired.value && hasMin.value) {
        label = $t('common.requirement.minimum_required', { count, kind });
    } else if (hasMin.value) {
        label = $t('common.requirement.minimum', { count, kind });
    } else if (isRequired.value) {
        label = $t('common.requirement.required', { count, kind });
    } else {
        label = $t('common.requirement.maximum', { count, kind });
    }

    return hasMax.value && !minEqualsMax.value ? $t('common.requirement.with_max', { label, maximum: props.specs.max }) : label;
});
</script>

<template>
    <UBadge v-if="showRequirement" :color="badgeColor" :variant="badgeVariant">
        <span :class="{ 'font-bold': isRequired }">{{ requirementLabel }}</span>
    </UBadge>
    <div v-else class="text-muted">
        {{ $t('common.requirement.not_required', { kind: displayName }) }}
    </div>
</template>
