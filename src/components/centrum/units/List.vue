<script lang="ts" setup>
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import ListEntry from './ListEntry.vue';

const props = defineProps<{
    units: Unit[];
}>();

defineEmits<{
    (e: 'goto', loc: { x: number; y: number }): void;
}>();

const sortedUnits = computed(
    () =>
        props.units?.sort((a, b) => {
            if (!b.status) return -1;
            if (!a.status) return 1;
            if (a.status < b.status) return -1;
            if (a.status > b.status) return 1;
            return 0;
        }),
);
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-scroll">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h2 class="text-base font-semibold leading-6 text-gray-100">{{ $t('common.unit', 2) }}</h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 -my-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <ul role="list" class="mt-3 grid grid-cols-1 gap-5 sm:grid-cols-2 sm:gap-2 lg:grid-cols-3">
                        <ListEntry
                            v-for="unit in sortedUnits"
                            :key="unit.id.toString()"
                            :unit="unit"
                            @goto="$emit('goto', $event)"
                        />
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
