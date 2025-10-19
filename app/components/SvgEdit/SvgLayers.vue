<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
const props = defineProps<{
    api?: any | null;
}>();

const layers = ref<
    {
        name: string;
        visible: boolean;
        locked: boolean;
    }[]
>([]);

onMounted(() => {
    // naive refresh; real impl should sync with svgcanvas layer model
    const refresh = () => {
        const list = props.api?.getLayers?.() ?? [];
        layers.value = list.map((l: any) => ({ name: l.name, visible: !l.hidden, locked: !!l.locked }));
    };
    props.api?.bind('changed', refresh);
    refresh();
});
</script>

<template>
    <UCard class="h-full">
        <template #header>Layers</template>
        <div class="space-y-2">
            <div v-for="l in layers" :key="l.name" class="flex items-center gap-2">
                <USwitch :model-value="l.visible" @update:model-value="(v) => props.api?.setLayerVisibility?.(l.name, v)" />
                <USwitch :model-value="l.locked" @update:model-value="(v) => props.api?.setLayerLock?.(l.name, v)" />
                <span class="text-sm">{{ l.name }}</span>
            </div>
        </div>
    </UCard>
</template>
