<script lang="ts" setup>
import { computed } from 'vue';
import type { ObjectSpecs } from '~~/gen/ts/resources/documents/templates';

const props = defineProps<{
    name: string;
    plural?: string;
    specs: ObjectSpecs;
}>();

const displayName = computed(() => props.plural ?? `${props.name}(s)`);

const isRequired = computed(() => props.specs.required);
const hasMin = computed(() => props.specs.min && props.specs.min > 0);
const hasMax = computed(() => props.specs.max && props.specs.max > 0);
const minEqualsMax = computed(() => props.specs.max === props.specs.min);

const showRequirement = computed(() => isRequired.value || (hasMin.value && hasMax.value));
</script>

<template>
    <div v-if="showRequirement" class="inline-flex items-center gap-1">
        <span v-if="isRequired" class="font-bold">{{ $t('common.require', 2) }} </span>
        <span v-if="hasMin">{{ $t('common.min') }} </span>
        <span v-if="minEqualsMax"> {{ props.specs.max }} {{ props.name }} </span>
        <span v-else>
            {{ props.specs.min === 0 && isRequired ? props.specs.max : props.specs.min }}
            {{ displayName }}
            <span v-if="hasMax">&nbsp;({{ $t('common.max') }}: {{ props.specs.max }})</span>
        </span>
    </div>
    <div v-else>{{ displayName }} {{ $t('common.not').toLocaleLowerCase() }} {{ $t('common.require', 2) }}</div>
</template>
