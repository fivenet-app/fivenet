<template>
    <div class="relative h-full w-full text-[10px] text-gray-500 select-none dark:text-white">
        <div v-for="mm in ticks" :key="mm" :style="getTickStyle(mm)">
            <div v-if="orientation === 'horizontal'" class="h-3 w-px bg-transparent"></div>
            <div v-if="orientation === 'horizontal'" class="absolute top-3 -translate-x-1/2">{{ mm }}mm</div>
            <div v-if="orientation === 'vertical'" class="h-px w-3 bg-transparent"></div>
            <div v-if="orientation === 'vertical'" class="absolute left-3 -translate-y-1/2">{{ mm }}mm</div>
        </div>
    </div>
</template>

<script lang="ts">
import { computed, defineComponent, type StyleValue } from 'vue';

function mmToPx(mm: number): number {
    // Replace this with your actual mm to px conversion logic
    return mm * 3.7795275591; // Example: 1mm = 3.7795275591px
}

export default defineComponent({
    name: 'Ruler',
    props: {
        lengthPx: { type: Number, required: true },
        orientation: { type: String as () => 'horizontal' | 'vertical', required: true },
        zoom: { type: Number, required: true },
    },
    setup(props) {
        const ticks = computed(() => {
            const result: number[] = [];
            const totalMm = Math.ceil(props.lengthPx / props.zoom / mmToPx(1));
            for (let mm = 0; mm <= totalMm; mm += 10) {
                result.push(mm);
            }
            return result;
        });

        const getTickStyle = (mm: number): StyleValue => {
            if (props.orientation === 'horizontal') {
                return {
                    position: 'absolute',
                    left: `${mmToPx(mm) * props.zoom}px`,
                    top: '0',
                    height: '100%',
                };
            } else {
                return {
                    position: 'absolute',
                    top: `${mmToPx(mm) * props.zoom}px`,
                    left: '0',
                    width: '100%',
                };
            }
        };

        return {
            ticks,
            getTickStyle,
        };
    },
});
</script>
