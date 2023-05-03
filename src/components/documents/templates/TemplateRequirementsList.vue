<script lang="ts" setup>
defineProps({
    name: {
        type: String,
        required: true,
    },
    plural: {
        type: String,
        required: false,
    },
    required: {
        type: Boolean,
        required: true,
    },
    min: {
        type: Number,
        required: true,
    },
    max: {
        type: Number,
        required: true,
    },
});
</script>

<template>
    <span v-tag v-if="required || (min > 0 && max > 0)">
        <span v-if="required">
            {{ $t('common.require', 2) }}{{ ' ' }}
        </span>
        <span v-if="min > 0">
            {{ $t('common.min') }}{{ ' ' }}
        </span>
        <span v-if="max == min">
            {{ max }} {{ name }}
        </span>
        <span v-else>
            {{ (min === 0 &&
                required) ?
                max : min }}
            {{ plural ?? name + "(s)" }}
            <span v-if="max! > 0">
                &nbsp;({{ $t('common.max') }}: {{ max }})
            </span>
        </span>
    </span>
    <span v-else>
        {{ plural ?? name + "(s)" }} {{ $t('common.not').toLocaleLowerCase() }} {{ $t('common.require', 2) }}
    </span>
</template>
