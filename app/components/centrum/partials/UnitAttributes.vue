<script lang="ts" setup>
import { UnitAttribute, type UnitAttributes } from '~~/gen/ts/resources/centrum/attributes';

defineProps<{
    attributes?: UnitAttributes;
}>();
</script>

<template>
    <template v-if="attributes !== undefined && attributes?.list.length > 0">
        <div class="inline-flex gap-1">
            <UBadge
                v-for="attribute in attributes?.list"
                :key="attribute"
                color="warning"
                variant="subtle"
                :label="$t(`enums.centrum.UnitAttribute.${UnitAttribute[attribute]}`)"
            />
        </div>
    </template>
    <span v-else>
        {{ $t('common.none', [$t('common.attributes', 2)]) }}
    </span>
</template>
