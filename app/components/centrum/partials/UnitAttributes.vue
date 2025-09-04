<script lang="ts" setup>
import { UnitAttribute, type UnitAttributes } from '~~/gen/ts/resources/centrum/attributes';

defineProps<{
    attributes?: UnitAttributes;
}>();
</script>

<template>
    <template v-if="attributes !== undefined && attributes?.list.length > 0">
        <div class="inline-flex gap-1">
            <span
                v-for="attribute in attributes?.list"
                :key="attribute"
                class="bg-warn-400/10 text-warn-400 ring-warn-400/20 inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset"
            >
                {{ $t(`enums.centrum.UnitAttribute.${UnitAttribute[attribute]}`) }}
            </span>
        </div>
    </template>
    <span v-else>
        {{ $t('common.none', [$t('common.attributes', 2)]) }}
    </span>
</template>
