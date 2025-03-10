<script lang="ts" setup>
import { DispatchAttribute, type DispatchAttributes } from '~~/gen/ts/resources/centrum/attributes';

defineProps<{
    attributes?: DispatchAttributes;
}>();
</script>

<template>
    <template v-if="attributes !== undefined && attributes?.list.length > 0">
        <div class="inline-flex gap-1">
            <span
                v-for="attribute in attributes?.list"
                :key="attribute"
                class="inline-flex items-center rounded-md bg-warn-400/10 px-2 py-1 text-xs font-medium text-warn-400 ring-1 ring-inset ring-warn-400/20"
            >
                {{ $t(`enums.centrum.DispatchAttribute.${DispatchAttribute[attribute]}`) }}
            </span>
        </div>
    </template>
    <span v-else>
        {{ $t('common.none', [$t('common.attributes', 2)]) }}
    </span>
</template>
