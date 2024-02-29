<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        throttle?: number;
        duration?: number;
        height?: number;
    }>(),
    {
        throttle: 175,
        duration: 2150,
        height: 3,
    },
);

// Options & Data
const data = reactive({
    progress: 0,
    isLoading: true,
    canSucceed: true,
});

// Local variables
let _timer: undefined | number | NodeJS.Timeout;
let _throttle: undefined | number | NodeJS.Timeout;
const _cut = 10000 / Math.floor(props.duration);

// Functions
function clear() {
    _timer && clearInterval(_timer);
    _throttle && clearTimeout(_throttle);
    _timer = undefined;
}
function start() {
    clear();
    data.progress = 0;

    if (props.throttle) {
        _throttle = setTimeout(startTimer, props.throttle);
    } else {
        startTimer();
    }
}
function increase(num: number) {
    data.progress = Math.min(100, Math.floor(data.progress + num));
}
function finish() {
    data.progress = 100;
    hide();
}
function hide() {
    clear();
    setTimeout(() => {
        data.isLoading = false;
        setTimeout(() => {
            data.progress = 0;
        }, 550);
    }, 500);
}
function startTimer() {
    data.isLoading = true;
    _timer = setInterval(() => {
        increase(_cut);
    }, 100);
}
function delayedFinish() {
    data.progress = 65;
    setTimeout(() => {
        finish();
    }, 500);
}

// Hooks
const nuxt = useNuxtApp();

// @ts-ignore we are currently unable to add custom event types to the typings
nuxt.hook('data:loading:start', start);
// @ts-ignore we are currently unable to add custom event types to the typings
nuxt.hook('data:loading:finish', delayedFinish);
// @ts-ignore we are currently unable to add custom event types to the typings
nuxt.hook('data:loading:finish_error', () => {
    data.canSucceed = false;
    delayedFinish();
    setTimeout(() => {
        data.canSucceed = true;
    }, 1250);
});

onBeforeUnmount(() => clear);
</script>

<template>
    <div
        :class="['nuxt-loading-indicator', data.canSucceed ? '' : 'nuxt-loading-indicator-failed']"
        :style="{
            position: 'fixed',
            top: 0,
            right: 0,
            left: 0,
            pointerEvents: 'none',
            width: 'auto',
            height: `${props.height}px`,
            opacity: data.isLoading ? 1 : 0,
            backgroundSize: `${(100 / data.progress) * 100}% auto`,
            transform: `scaleX(${data.progress}%)`,
            transformOrigin: 'left',
            transition: 'transform 0.1s, height 0.4s, opacity 0.4s',
            zIndex: 999999,
        }"
    />
</template>

<style scoped>
.nuxt-loading-indicator {
    background: repeating-linear-gradient(to right, #55dde0 0%, #34cdfe 50%, #7161ef 100%);
}

.nuxt-loading-indicator-failed {
    background: repeating-linear-gradient(to right, #d72638 0%, #ac1e2d 50%, #d72638 100%);
}
</style>
