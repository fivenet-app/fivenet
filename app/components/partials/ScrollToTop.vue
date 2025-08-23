<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        element?: HTMLElement;
        showHeight?: number;
    }>(),
    {
        element: undefined,
        showHeight: 200,
    },
);

const { y } = useScroll(() => props.element);

const show = ref(false);

watchDebounced(
    y,
    () => {
        if (y.value > props.showHeight) {
            show.value = true;
        } else {
            show.value = false;
        }
    },
    { debounce: 75 },
);

function scrollToTop() {
    props.element?.scrollTo({
        top: 0,
        behavior: 'smooth',
    });
}
</script>

<template>
    <Transition
        enter-active-class="transition-opacity duration-100 sm:duration-200"
        enter-from-class="opacity-0"
        leave-to-class="opacity-0"
    >
        <UTooltip
            v-if="show"
            class="z-100 fixed bottom-32 right-6 transition delay-150 duration-300 ease-in-out"
            :text="$t('common.scroll_to_top')"
            :prevent="!show"
        >
            <UButton icon="i-mdi-arrow-top" @click="scrollToTop" />
        </UTooltip>
    </Transition>
</template>
