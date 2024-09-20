<script lang="ts" setup>
import type { ObjectSpecs } from '~~/gen/ts/resources/documents/templates';

defineProps<{
    name: string;
    plural?: string;
    specs: ObjectSpecs;
}>();
</script>

<template>
    <span v-if="specs.required || (specs.min && specs.min > 0 && specs.max && specs.max > 0)">
        <span v-if="specs.required"> {{ $t('common.require', 2) }}{{ ' ' }} </span>
        <span v-if="specs.min && specs.min > 0"> {{ $t('common.min') }}{{ ' ' }} </span>
        <span v-if="specs.max === specs.min"> {{ specs.max }} {{ name }} </span>
        <span v-else>
            {{ specs.min === 0 && specs.required ? specs.max : specs.min }}
            {{ plural ?? name + '(s)' }}
            <span v-if="specs.max && specs.max > 0"> &nbsp;({{ $t('common.max') }}: {{ specs.max }}) </span>
        </span>
    </span>
    <span v-else>
        {{ plural ?? name + '(s)' }} {{ $t('common.not').toLocaleLowerCase() }}
        {{ $t('common.require', 2) }}
    </span>
</template>
