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
    <span v-if="required">
        Requires
    </span>
    <span v-if="max == min">
        1 {{ name }}
    </span>
    <span v-else>
        {{ (min === 0 &&
            required) ?
            max : min }}
        {{ plural ?? name + "(s)" }}
        <span v-if="max! > 0">
        &nbsp;(Max: {{ max }})
        </span>
    </span>
</template>
