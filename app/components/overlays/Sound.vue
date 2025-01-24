<!-- eslint-disable vue/multi-word-component-names -->
<script lang="ts" setup>
import { Howl } from 'howler';

const props = withDefaults(
    defineProps<{
        id?: number;
        name: SoundBites;
        volume: number;
        rate?: number;
    }>(),
    {
        id: undefined,
        volume: 0.5,
        rate: 1.0,
    },
);

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const sound = new Howl({
    src: [`/sounds/${props.name}.mp3`],
    volume: props.volume,
    rate: props.rate,
});

sound.on('end', () => emit('close'));

sound.play();

onBeforeUnmount(async () => sound.stop());
</script>

<template>
    <div></div>
</template>
