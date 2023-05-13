<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'
import { RoleAttribute } from '@fivenet/gen/resources/permissions/permissions_pb';
import { ChevronDownIcon } from '@heroicons/vue/24/solid';

defineProps<{
    attribute?: RoleAttribute,
    disabled?: boolean,
    states: Map<number, (string | number)[]>,
}>();

defineEmits<{
    (e: 'update:states', payload: Map<number, (string | number)[]>): void,
}>();

</script>

<style scoped>
.upsidedown {
    transform: rotate(180deg);
}
</style>

<template>
    <div v-if="$props.attribute">
        <Disclosure as="div"
            :class="[$props.disabled ? 'border-neutral/10 text-base-300' : 'hover:border-neutral/70 border-neutral/20 text-neutral']"
            v-slot="{ open }">
            <DisclosureButton :disabled="$props.disabled"
                :class="[open ? 'rounded-t-lg border-b-0' : 'rounded-lg', $props.disabled ? 'cursor-not-allowed' : '', ' flex w-full items-start justify-between text-left border-2 p-2 border-inherit transition-colors']">
                <span class="text-base leading-7 transition-colors">
                    Options
                </span>
                <span class="ml-6 flex h-7 items-center">
                    <ChevronDownIcon :class="[open ? 'upsidedown' : '', 'h-6 w-6 transition-transform']"
                        aria-hidden="true" />
                </span>
            </DisclosureButton>
            <DisclosurePanel class="px-4 pb-2 border-2 border-t-0 rounded-b-lg transition-colors border-inherit -mt-2">
                <div class="flex flex-col gap-2 max-w-4xl mx-auto my-2">
                    {{ $props.attribute.getType() }} {{ $props.attribute.getValidValuesList() }}
                </div>
            </DisclosurePanel>
        </Disclosure>
    </div>
</template>
