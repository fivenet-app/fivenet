<script lang="ts" setup>
import type { Property } from '~~/gen/ts/resources/properties/properties';
import { pinnedProperties } from '~/composables/usePropertyRenderer';

const props = defineProps<{ properties?: Property[] }>();
const { propertyLabel } = usePropertyRenderer();

const pinned = computed(() => pinnedProperties(props.properties));

function label(property: Property): string {
    return propertyLabel(property);
}
</script>

<template>
    <UCard v-if="pinned.length" :ui="{ body: 'p-0 sm:p-0' }">
        <template #header>
            <h2 class="font-semibold text-highlighted">{{ $t('common.pinned_properties') }}</h2>
        </template>

        <dl class="divide-y divide-default">
            <div v-for="property in pinned" :key="property.key" class="flex items-center gap-4 px-4 py-3 sm:px-5">
                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ label(property) }}</dt>
                <dd class="min-w-0 text-sm text-toned"><PropertyRender :property="property" /></dd>
            </div>
        </dl>
    </UCard>
</template>
