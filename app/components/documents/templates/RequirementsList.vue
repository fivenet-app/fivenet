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
const minimum = computed(() => (hasMin.value ? props.specs.min : isRequired.value ? 1 : (props.specs.min ?? 0)));

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
</script>

<template>
    <UBadge v-if="showRequirement" :color="badgeColor" :variant="badgeVariant">
        <span v-if="isRequired" class="font-bold">{{ $t('common.require', 2) }} </span>
        <span v-if="hasMin">{{ $t('common.min') }} </span>
        <span v-if="minEqualsMax"> {{ props.specs.max }} {{ props.name }} </span>
        <span v-else>
            {{ minimum }}
            {{ displayName }}
            <span v-if="hasMax">&nbsp;({{ $t('common.max') }}: {{ props.specs.max }})</span>
        </span>
    </UBadge>
    <div v-else class="text-muted">
        {{ displayName }} {{ $t('common.not').toLocaleLowerCase() }} {{ $t('common.require', 2) }}
    </div>
</template>
