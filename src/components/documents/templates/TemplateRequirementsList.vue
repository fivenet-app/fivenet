<script lang="ts" setup>
import { ObjectSpecs } from '@fivenet/gen/resources/documents/templates_pb';

defineProps({
    name: {
        type: String,
        required: true,
    },
    plural: {
        type: String,
        required: false,
    },
    specs: {
        type: ObjectSpecs,
        required: true,
    }
});
</script>

<template>
    <span v-if="specs.getRequired() || (specs.getMin() > 0 && specs.getMax() > 0)">
        <span v-if="specs.getRequired()">
            {{ $t('common.require', 2) }}{{ ' ' }}
        </span>
        <span v-if="specs.getMin() > 0">
            {{ $t('common.min') }}{{ ' ' }}
        </span>
        <span v-if="specs.getMax() == specs.getMin()">
            {{ specs.getMax() }} {{ name }}
        </span>
        <span v-else>
            {{ (specs.getMin() === 0 &&
                specs.getRequired()) ?
                specs.getMax() : specs.getMin() }}
            {{ plural ?? name + "(s)" }}
            <span v-if="specs.getMax() > 0">
                &nbsp;({{ $t('common.max') }}: {{ specs.getMax() }})
            </span>
        </span>
    </span>
    <span v-else>
        {{ plural ?? name + "(s)" }} {{ $t('common.not').toLocaleLowerCase() }} {{ $t('common.require', 2) }}
    </span>
</template>
