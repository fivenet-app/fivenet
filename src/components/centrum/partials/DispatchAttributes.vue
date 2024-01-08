<script lang="ts" setup>
import type { Attributes } from '~~/gen/ts/resources/centrum/general';

defineProps<{
    attributes?: Attributes;
}>();
</script>

<template>
    <template v-if="attributes !== undefined && attributes?.list.length > 0">
        <span
            v-for="attribute in attributes?.list"
            :key="attribute"
            class="inline-flex items-center rounded-md bg-warn-400/10 px-2 py-1 text-xs font-medium text-warn-400 ring-1 ring-inset ring-warn-400/20"
        >
            {{ $t(`components.centrum.dispatches.attributes.${attribute}`) }}
        </span>
    </template>
    <span v-else>
        {{ $t('common.none_selected', [$t('common.attributes', 2)]) }}
    </span>
</template>
