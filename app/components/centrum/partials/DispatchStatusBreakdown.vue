<script lang="ts" setup>
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import { useCentrumStore } from '~/store/centrum';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

const centrumStore = useCentrumStore();
const { dispatches } = storeToRefs(centrumStore);

const counts = computedAsync(() => {
    const count = {
        unassigned: 0,
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
                    count.enRoute++;
                }
                break;

            case StatusDispatch.UNIT_ASSIGNED:
            case StatusDispatch.UNIT_ACCEPTED:
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

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UPopover class="flex-1">
        <UButton
            v-bind="$attrs"
            :ui="{ icon: { base: 'max-md:!hidden' } }"
            variant="ghost"
            class="items-center"
            trailing-icon="i-mdi-chevron-down"
            block
        >
            {{ $t('components.centrum.livemap.total_dispatches') }}: {{ dispatches.size }}
        </UButton>

        <template #panel>
            <div class="p-4">
                <ul class="text-nowrap text-sm font-normal">
                    <li>
                        <span class="text-black" :class="dispatchStatusToBGColor(StatusDispatch.NEW)"
                            >{{ $t('enums.centrum.StatusDispatch.UNASSIGNED') }}:</span
                        >
                        {{ counts.unassigned }}
                    </li>
                    <li>
                        <span class="text-black" :class="dispatchStatusToBGColor(StatusDispatch.EN_ROUTE)">{{
                            $t('enums.centrum.StatusDispatch.EN_ROUTE')
                        }}</span
                        >: {{ counts.enRoute }}
                    </li>
                    <li>
                        <span class="text-black" :class="dispatchStatusToBGColor(StatusDispatch.ON_SCENE)">{{
                            $t('enums.centrum.StatusDispatch.ON_SCENE')
                        }}</span
                        >: {{ counts.onScene }}
                    </li>
                    <li>
                        <span class="text-black" :class="dispatchStatusToBGColor(StatusDispatch.NEED_ASSISTANCE)">{{
                            $t('enums.centrum.StatusDispatch.NEED_ASSISTANCE')
                        }}</span
                        >: {{ counts.needAssistance }}
                    </li>
                    <li>
                        <span class="text-black" :class="dispatchStatusToBGColor(StatusDispatch.COMPLETED)"
                            >{{ $t('enums.centrum.StatusDispatch.COMPLETED') }}:</span
                        >
                        {{ counts.completed }}
                    </li>
                    <li class="underline">{{ $t('common.total_count') }}: {{ dispatches.size }}</li>
                </ul>
            </div>
        </template>
    </UPopover>
</template>
