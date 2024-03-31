<script lang="ts" setup>
import { Float } from '@headlessui-float/vue';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { computedAsync } from '@vueuse/core';
import { useCentrumStore } from '~/store/centrum';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { dispatchStatusToBGColor } from '../helpers';

const centrumStore = useCentrumStore();
const { dispatches } = storeToRefs(centrumStore);

const counts = computedAsync(() => {
    const count = {
        unassigned: 0,
        assigned: 0,
        enRoute: 0,
        onScene: 0,
        needAssistance: 0,
        completed: 0,
    };

    dispatches.value.forEach((dsp) => {
        switch (dsp.status?.status) {
            case StatusDispatch.NEW:
            case StatusDispatch.UNASSIGNED:
                count.unassigned++;
                break;

            case StatusDispatch.UNIT_UNASSIGNED:
            case StatusDispatch.UNIT_DECLINED:
                if (dsp.units.length <= 0) {
                    count.unassigned++;
                } else {
                    count.assigned++;
                }
                break;

            case StatusDispatch.UNIT_ASSIGNED:
            case StatusDispatch.UNIT_ACCEPTED:
                count.assigned++;
                break;

            case StatusDispatch.EN_ROUTE:
                count.enRoute++;
                break;

            case StatusDispatch.ON_SCENE:
                count.onScene++;
                break;
            case StatusDispatch.NEED_ASSISTANCE:
                count.needAssistance++;
                break;

            case StatusDispatch.COMPLETED:
            case StatusDispatch.CANCELLED:
            case StatusDispatch.ARCHIVED:
            default:
                count.completed++;
        }
    });

    return count;
});
</script>

<template>
    <Popover>
        <Float
            portal
            auto-placement
            :offset="40"
            enter="transition duration-150 ease-out"
            enter-from="scale-95 opacity-0"
            enter-to="scale-100 opacity-100"
            leave="transition duration-100 ease-in"
            leave-from="scale-100 opacity-100"
            leave-to="scale-95 opacity-0"
        >
            <PopoverButton class="inline-flex items-center">
                {{ $t('components.centrum.livemap.total_dispatches') }}: {{ dispatches.size }}
            </PopoverButton>

            <PopoverPanel
                class="absolute z-5 w-52 min-w-fit max-w-60 rounded-lg border border-gray-600 bg-gray-800 text-sm text-gray-400 shadow-sm transition-opacity"
            >
                <div class="p-3">
                    <ul class="text-nowrap text-sm font-normal text-neutral">
                        <li>
                            <span :class="dispatchStatusToBGColor(StatusDispatch.NEW)"
                                >{{ $t('enums.centrum.StatusDispatch.UNASSIGNED') }}:</span
                            >
                            {{ counts.unassigned }}
                        </li>
                        <li>
                            <span :class="dispatchStatusToBGColor(StatusDispatch.UNIT_ASSIGNED)">{{
                                $t('enums.centrum.StatusDispatch.UNIT_ASSIGNED')
                            }}</span
                            >: {{ counts.assigned }}
                        </li>
                        <li>
                            <span :class="dispatchStatusToBGColor(StatusDispatch.EN_ROUTE)">{{
                                $t('enums.centrum.StatusDispatch.EN_ROUTE')
                            }}</span
                            >: {{ counts.enRoute }}
                        </li>
                        <li>
                            <span :class="dispatchStatusToBGColor(StatusDispatch.ON_SCENE)">{{
                                $t('enums.centrum.StatusDispatch.ON_SCENE')
                            }}</span
                            >: {{ counts.onScene }}
                        </li>
                        <li>
                            <span :class="dispatchStatusToBGColor(StatusDispatch.NEED_ASSISTANCE)">{{
                                $t('enums.centrum.StatusDispatch.NEED_ASSISTANCE')
                            }}</span
                            >: {{ counts.needAssistance }}
                        </li>
                        <li>
                            <span :class="dispatchStatusToBGColor(StatusDispatch.COMPLETED)"
                                >{{ $t('enums.centrum.StatusDispatch.COMPLETED') }}:</span
                            >
                            {{ counts.completed }}
                        </li>
                        <li class="underline">{{ $t('common.total_count') }}: {{ dispatches.size }}</li>
                    </ul>
                </div>
            </PopoverPanel>
        </Float>
    </Popover>
</template>
