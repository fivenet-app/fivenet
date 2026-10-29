<script lang="ts" setup>
import type { PropertyValue } from '~~/gen/ts/resources/properties/properties';

const props = defineProps<{
    propertyKey: string;
    oldValue?: PropertyValue;
    newValue?: PropertyValue;
}>();

const { renderProperty } = usePropertyRenderer();

function render(value: PropertyValue | undefined): string {
    if (!value) return '';
    return renderProperty({ key: props.propertyKey, value });
}

const isClear = computed(() => props.oldValue !== undefined && props.newValue === undefined);
const isCreate = computed(() => props.oldValue === undefined && props.newValue !== undefined);
</script>

<template>
    <I18nT
        :keypath="
            isClear
                ? 'components.properties.activity.cleared'
                : isCreate
                  ? 'components.properties.activity.set'
                  : 'components.properties.activity.changed'
        "
    >
        <template #property
            ><code class="font-semibold">{{ propertyKey }}</code></template
        >
        <template #old
            ><span class="font-semibold line-through">{{ render(oldValue) }}</span></template
        >
        <template #new
            ><span class="font-semibold">{{ render(newValue) }}</span></template
        >
    </I18nT>
</template>
