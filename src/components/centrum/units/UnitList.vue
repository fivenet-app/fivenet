<script lang="ts" setup>
import { CogIcon } from 'mdi-vue3';
import { computedAsync } from '@vueuse/core';
import { useCentrumStore } from '~/store/centrum';
import UnitListEntry from '~/components/centrum/units/UnitListEntry.vue';
import { StatusUnit } from '~~/gen/ts/resources/centrum/units';
import { statusOrder, type GroupedUnits } from '~/components/centrum/helpers';

const centrumStore = useCentrumStore();
const { getSortedUnits } = storeToRefs(centrumStore);

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const grouped = computedAsync(async () => {
    const groups: GroupedUnits = [];
    getSortedUnits.value.forEach((e) => {
        const idx = groups.findIndex((g) => g.key === e.status?.status.toString());
        if (idx === -1) {
            groups.push({
                status: e.status?.status ?? 0,
                units: [e],
                key: e.status?.status.toString() ?? '',
            });
        } else {
            groups[idx].units.push(e);
        }
    });

    groups
        .sort((a, b) => statusOrder.indexOf(a.status) - statusOrder.indexOf(b.status))
        .forEach((group) =>
            group.units.sort((a, b) => {
                if (a.users.length === b.users.length) {
                    return 0;
                } else if (a.users.length === 0) {
                    return 1;
                } else if (b.users.length === 0) {
                    return -1;
                } else {
                    return a.name.localeCompare(b.name);
                }
            }),
        );

    return groups;
});
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-auto">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto inline-flex items-center">
                <h2 class="text-base font-semibold leading-6 text-gray-100 inline-flex items-center">
                    {{ $t('common.units') }}
                    <NuxtLink
                        v-if="can('CentrumService.CreateOrUpdateUnit')"
                        :to="{ name: 'centrum-units' }"
                        :title="$t('common.units')"
                        class="ml-2"
                    >
                        <CogIcon class="h-6 w-6" />
                    </NuxtLink>
                </h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 -my-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <template v-for="group in grouped" :key="group.key">
                        <p class="-mb-1.5 text-sm text-neutral">
                            {{ $t(`enums.centrum.StatusUnit.${StatusUnit[group.status]}`) }}
                        </p>
                        <ul role="list" class="mt-3 grid grid-cols-1 gap-4 sm:grid-cols-2 sm:gap-1.5 lg:grid-cols-3">
                            <UnitListEntry
                                v-for="unit in group.units"
                                :key="unit.id"
                                :unit="unit"
                                @goto="$emit('goto', $event)"
                            />
                        </ul>
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>
