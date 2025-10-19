<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import { DispatchAttribute, type DispatchAttributes } from '~~/gen/ts/resources/centrum/attributes';

interface Props extends /* @vue-ignore */ BadgeProps {
    attributes?: DispatchAttributes;
}

defineProps<Props>();
</script>

<template>
    <template v-if="attributes !== undefined && attributes?.list.length > 0">
        <div class="inline-flex gap-1">
            <UBadge
                v-for="attribute in attributes?.list"
                :key="attribute"
                color="warning"
                variant="outline"
                :label="$t(`enums.centrum.DispatchAttribute.${DispatchAttribute[attribute]}`)"
                v-bind="$attrs"
            />
        </div>
    </template>
    <span v-else>
        {{ $t('common.none', [$t('common.attributes', 2)]) }}
    </span>
</template>
