<script lang="ts" setup>
import { useCentrumStore } from '~/store/centrum';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import ListEntry from './ListEntry.vue';

const centrumStore = useCentrumStore();
const { units } = storeToRefs(centrumStore);

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
    (e: 'details', unit: Unit): void;
}>();
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-scroll">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h2 class="text-base font-semibold leading-6 text-gray-100">{{ $t('common.units') }}</h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 -my-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <ul role="list" class="mt-3 grid grid-cols-1 gap-5 sm:grid-cols-2 sm:gap-2 lg:grid-cols-3">
                        <ListEntry
                            v-for="unit in units"
                            :key="unit.id.toString()"
                            :unit="unit"
                            @goto="$emit('goto', $event)"
                            @details="$emit('details', $event)"
                        />
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
