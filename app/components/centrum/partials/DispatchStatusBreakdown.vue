<script lang="ts" setup>
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import { useCentrumStore } from '~/stores/centrum';
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
                <UIcon v-if="!counts" name="i-mdi-refresh" class="size-4 animate-spin" />
                <div v-else class="flex flex-col gap-1 text-nowrap text-sm font-normal">
                    <div class="inline-flex justify-between gap-1.5">
                        <UBadge class="px-2 py-1" :class="dispatchStatusToBGColor(StatusDispatch.UNASSIGNED)" size="sm">
                            {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[StatusDispatch.UNASSIGNED]}`) }}
                        </UBadge>
                        <p class="font-semibold">{{ counts?.unassigned }}</p>
                    </div>
                    <div class="inline-flex justify-between gap-1.5">
                        <UBadge class="px-2 py-1" :class="dispatchStatusToBGColor(StatusDispatch.EN_ROUTE)" size="sm">
                            {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[StatusDispatch.EN_ROUTE]}`) }}
                        </UBadge>
                        <p class="font-semibold">{{ counts.enRoute }}</p>
                    </div>
                    <div class="inline-flex justify-between gap-1.5">
                        <UBadge class="px-2 py-1" :class="dispatchStatusToBGColor(StatusDispatch.ON_SCENE)" size="sm">
                            {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[StatusDispatch.ON_SCENE]}`) }}
                        </UBadge>
                        <p class="font-semibold">{{ counts.onScene }}</p>
                    </div>
                    <div class="inline-flex justify-between gap-1.5">
                        <UBadge class="px-2 py-1" :class="dispatchStatusToBGColor(StatusDispatch.NEED_ASSISTANCE)" size="sm">
                            {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[StatusDispatch.NEED_ASSISTANCE]}`) }}
                        </UBadge>
                        <p class="font-semibold">{{ counts.needAssistance }}</p>
                    </div>
                    <div class="inline-flex justify-between gap-1.5">
                        <UBadge class="px-2 py-1" :class="dispatchStatusToBGColor(StatusDispatch.COMPLETED)" size="sm">
                            {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[StatusDispatch.COMPLETED]}`) }}
                        </UBadge>
                        <p class="font-semibold">{{ counts.completed }}</p>
                    </div>
                    <div class="flex justify-between font-semibold">
                        <span>{{ $t('common.total_count') }}</span>
                        <span>{{ dispatches.size }}</span>
                    </div>
                </div>
            </div>
        </template>
    </UPopover>
</template>
