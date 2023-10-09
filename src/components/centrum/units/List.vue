<script lang="ts" setup>
import { CogIcon } from 'mdi-vue3';
import { useCentrumStore } from '~/store/centrum';
import ListEntry from './ListEntry.vue';

const centrumStore = useCentrumStore();
const { getSortedUnits } = storeToRefs(centrumStore);

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const sortedUnits = computed(() => getSortedUnits.value);
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-auto">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto inline-flex items-center">
                <h2 class="text-base font-semibold leading-6 text-gray-100 inline-flex">
                    {{ $t('common.units') }}
                    <NuxtLink :to="{ name: 'centrum-units' }" :title="$t('common.units')" class="ml-2">
                        <CogIcon class="h-6 w-6" />
                    </NuxtLink>
                </h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 -my-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <ul role="list" class="mt-3 grid grid-cols-1 gap-5 sm:grid-cols-2 sm:gap-2 lg:grid-cols-3">
                        <ListEntry v-for="unit in sortedUnits.reverse()" :unit="unit" @goto="$emit('goto', $event)" />
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
