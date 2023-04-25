<script lang="ts" setup>
const props = defineProps({
    throttle: {
        type: Number,
        default: 200,
    },
    duration: {
        type: Number,
        default: 2000,
    },
    height: {
        type: Number,
        default: 3,
    },
    percent: {
        type: Number,
        default: 0,
        required: false,
    },
    show: {
        type: Boolean,
        required: false,
        default: false,
    },
    canSucceed: {
        type: Boolean,
        required: false,
        default: true,
    },
});

// Options & Data
const data = reactive({
    percent: props.percent,
    left: 0,
    show: props.show,
    canSucceed: props.canSucceed,
});

// Local variables
let _timer: undefined | number | NodeJS.Timeout = undefined;
let _throttle: undefined | NodeJS.Timeout = undefined;
const _cut = 10000 / Math.floor(props.duration);

// Functions
const clear = () => {
    _timer && clearInterval(_timer);
    _throttle && clearTimeout(_throttle);
    _timer = undefined;
};
const start = () => {
    clear();
    data.percent = 0;
    data.canSucceed = true;

    if (props.throttle) {
        _throttle = setTimeout(startTimer, props.throttle);
    } else {
        startTimer();
    }
};
const set = (num: number) => {
    data.show = true;
    data.canSucceed = true;
    data.percent = Math.min(100, Math.max(0, Math.floor(num)));
};
const increase = (num: number) => {
    data.percent = Math.min(100, Math.floor(data.percent + num));
};
const decrease = (num: number) => {
    data.percent = Math.max(0, Math.floor(data.percent - num));
};
const pause = () => clearInterval(_timer);
const resume = () => startTimer();
const finish = () => {
    data.percent = 100;
    hide();
};
const hide = () => {
    clear();
    setTimeout(() => {
        data.show = false;
        setTimeout(() => {
            data.percent = 0;
        }, 400);
    }, 500);
};
const startTimer = () => {
    data.show = true;
    _timer = setInterval(() => {
        increase(_cut);
    }, 100);
};

// Hooks
const nuxt = useNuxtApp();

nuxt.hook('page:start', start);
nuxt.hook('page:finish', () => {
    setTimeout(() => {
        finish();
    }, 250);
});
watch(data, () => {
    if (props.show) {
        start();
    } else {
        setTimeout(() => {
            finish();
        }, 250);
    }
});

onBeforeUnmount(() => clear);
</script>

<template>
    <div class="nuxt-progress" :class="{
            'nuxt-progress-failed': !data.canSucceed,
        }" :style="{
        width: data.percent + '%',
        left: data.left,
        height: props.height + 'px',
        opacity: data.show ? 1 : 0,
        backgroundSize: (100 / data.percent) * 100 + '% auto',
    }" />
</template>

<style scoped>
.nuxt-progress {
    position: fixed;
    top: 0px;
    left: 0px;
    right: 0px;
    width: 0%;
    opacity: 1;
    transition: width 0.1s, height 0.4s, opacity 0.4s;
    background: repeating-linear-gradient(to right,
            #55dde0 0%,
            #34cdfe 50%,
            #7161ef 100%);
    z-index: 999999;
}
</style>
