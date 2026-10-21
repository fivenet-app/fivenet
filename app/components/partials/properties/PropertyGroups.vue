<script lang="ts" setup>
import PropertyRender from './PropertyRender.vue';
import {
    PropertyDisplayMode,
    PropertyRenderer,
    type Property,
    type PropertyGroup,
} from '~~/gen/ts/resources/properties/properties';

const props = defineProps<{ properties?: Property[]; groups?: PropertyGroup[] }>();
const { t, te } = useI18n();
const { propertyLabel } = usePropertyRenderer();

const groupMetadata = computed(() => new Map((props.groups ?? []).map((group) => [group.key, group])));

const groups = computed(() => {
    const grouped = new Map<string, Property[]>();
    for (const property of props.properties ?? []) {
        if (property.definition?.displayMode && property.definition.displayMode !== PropertyDisplayMode.DEFAULT) continue;
        if (property.value?.value.oneofKind === undefined && property.definition?.renderer !== PropertyRenderer.EMAIL) continue;
        const key = property.definition?.groupKey || '';
        grouped.set(key, [...(grouped.get(key) ?? []), property]);
    }
    return [...grouped.entries()]
        .map(([key, properties]) => ({ key, properties }))
        .sort((left, right) => {
            const leftOrder =
                groupMetadata.value.get(left.key)?.sortOrder ??
                left.properties[0]?.definition?.groupSortOrder ??
                Number.MAX_SAFE_INTEGER;
            const rightOrder =
                groupMetadata.value.get(right.key)?.sortOrder ??
                right.properties[0]?.definition?.groupSortOrder ??
                Number.MAX_SAFE_INTEGER;
            return leftOrder - rightOrder || left.key.localeCompare(right.key);
        });
});

function groupLabel(key: string): string {
    if (!key) return t('components.settings.property_definitions.ungrouped');
    const group = groupMetadata.value.get(key);
    if (group?.label) return group.label;
    if (group?.displayKey && te(group.displayKey)) return t(group.displayKey);
    const translationKey = `components.settings.property_definitions.groups.${key}`;
    return te(translationKey) ? t(translationKey) : key;
}

function label(property: Property): string {
    return propertyLabel(property);
}
</script>

<template v-if="groups.length">
    <UCard v-for="group in groups" :key="group.key" :ui="{ body: 'p-0 sm:p-0' }">
        <template #header>
            <h2 class="font-semibold text-highlighted">{{ groupLabel(group.key) }}</h2>
        </template>

        <dl class="divide-y divide-default">
            <div v-for="property in group.properties" :key="property.key" class="flex items-center gap-4 px-4 py-3 sm:px-5">
                <dt class="w-36 shrink-0 text-sm font-medium text-muted">{{ label(property) }}</dt>
                <dd class="min-w-0 text-sm text-toned">
                    <PropertyRender :property="property" />
                </dd>
            </div>
        </dl>
    </UCard>
</template>
