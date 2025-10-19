<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import { ref } from 'vue';
import { useSvgEdit } from '~/composables/useSvgEdit';

const props = withDefaults(
    defineProps<{
        config?: Record<string, any>;
        initialSvg?: string;
    }>(),
    {
        config: () => ({}),
        initialSvg: '',
    },
);

const emit = defineEmits<{
    (e: 'ready', api: any): void;
    (e: 'error', err: unknown): void;
    (e: 'changed'): void;
    (e: 'selected', els: any[]): void;
}>();

const host = ref<HTMLElement | null>(null);
const workarea = ref<HTMLElement | null>(null);
const textInput = ref<HTMLInputElement | null>(null);
const editor = useSvgEdit({
    container: host,
    workarea: workarea,
    config: props.config,
    initialSvg: props.initialSvg,
});

const enableToolCancel = ref(true);

// Attach container after mount so the composable can create the canvas
watchEffect(() => {
    if (host.value && !editor.containerEl.value) editor.containerEl.value = host.value;
});

// Extracted event listener callbacks

function preventBrowserScaling(e: WheelEvent) {
    if (e.ctrlKey) {
        e.preventDefault();
    }
}

function handleWindowMouseUp(e: MouseEvent) {
    enableToolCancel.value = true;
    if (e.button === 1) {
        editor.canvas.value.setMode(previousMode.value ?? 'select');
    }
    panning.value = false;
}

function handleDocumentKeyDown(e: KeyboardEvent) {
    if (!e.target || (e.target as HTMLElement).nodeName !== 'BODY') return;

    if (e.code.toLowerCase() === 'space') {
        editor.canvas.value.spaceKey = keypan.value = true;
        e.preventDefault();
    } else if (e.key.toLowerCase() === 'shift' && editor.canvas.value.getMode() === 'zoom') {
        //this.workarea.style.cursor = zoomOutIcon;
        e.preventDefault();
    } else if (e.code.toLowerCase() === 'delete') {
        editor.canvas.value.deleteSelectedElements();
        e.preventDefault();
    }
}

function handleWorkareaWheel(e: WheelEvent) {
    if (e.altKey) {
        e.preventDefault();
        editor.setZoom(e.deltaY > 0 ? editor.canvas.value.getZoom() * 0.9 : editor.canvas.value.getZoom() * 1.1, true);
        editor.canvas.value.updateCanvas(true);
    }
}

function handleDocumentKeyUp(e: KeyboardEvent) {
    if (!e.target || (e.target as HTMLElement).nodeName !== 'BODY') return;
    if (e.code.toLowerCase() === 'space') {
        editor.canvas.value.spaceKey = keypan.value = false;
        editor.canvas.value.setMode(previousMode.value === 'ext-panning' ? 'select' : (previousMode.value ?? 'select'));
        e.preventDefault();
    } else if (e.key.toLowerCase() === 'shift' && editor.canvas.value.getMode() === 'zoom') {
        //this.workarea.style.cursor = zoomInIcon;
        e.preventDefault();
    }
}

const lastX = ref<number>(0);
const lastY = ref<number>(0);
const panning = ref(false);
const keypan = ref(false);
const previousMode = ref('select');

// Updated watch callback
watch(
    () => editor.ready.value,
    (r) => {
        if (!r) return;

        if (textInput.value) editor.canvas.value.textActions.setInputElem(textInput.value);

        const addListenerMulti = (element: HTMLElement | null, eventNames: string, listener: (e: Event) => void) => {
            eventNames.split(' ').forEach((eventName) => element?.addEventListener(eventName, listener, false));
        };
        addListenerMulti(textInput.value, 'keyup input', (e) => {
            const event = e as KeyboardEvent;
            const input = event.target as HTMLInputElement;
            e.currentTarget && editor.canvas.value.setTextContent(input.value);
        });

        host.value?.addEventListener('mouseup', (e) => {
            if (panning.value === false) {
                return true;
            }

            if (workarea.value) {
                workarea.value.scrollLeft -= e.clientX - lastX.value;
                workarea.value.scrollTop -= e.clientY - lastY.value;
            }

            lastX.value = e.clientX;
            lastY.value = e.clientY;

            if (e.type === 'mouseup') {
                panning.value = false;
            }
            return false;
        });
        host.value?.addEventListener('mousemove', (e) => {
            if (panning.value === false) {
                return true;
            }

            if (workarea.value) {
                workarea.value.scrollLeft -= e.clientX - lastX.value;
                workarea.value.scrollTop -= e.clientY - lastY.value;
            }

            lastX.value = e.clientX;
            lastY.value = e.clientY;

            if (e.type === 'mouseup') {
                panning.value = false;
            }
            return false;
        });
        host.value?.addEventListener('mousedown', (e) => {
            enableToolCancel.value = false;
            if (e.button === 1 || keypan.value === true) {
                // prDefault to avoid firing of browser's panning on mousewheel
                e.preventDefault();
                panning.value = true;
                previousMode.value = editor.canvas.value.getMode();
                editor.canvas.value.setMode('ext-panning');
                //this.workarea.style.cursor = 'grab';
                lastX.value = e.clientX;
                lastY.value = e.clientY;
                return false;
            }
            return true;
        });

        document.addEventListener('wheel', (e) => preventBrowserScaling(e));
        window.addEventListener('mouseup', (e) => handleWindowMouseUp(e));
        document.addEventListener('keydown', (e) => handleDocumentKeyDown(e));
        workarea.value?.addEventListener('wheel', (e) => handleWorkareaWheel(e));
        document.addEventListener('keyup', (e) => handleDocumentKeyUp(e));

        emit('ready', editor.canvas.value);
    },
);
watch(
    () => editor.error.value,
    (err) => {
        console.log('error?', err);
        if (err) emit('error', err);
    },
);
watch(
    () => editor.selection.value,
    (els) => emit('selected', els),
);
</script>

<template>
    <div id="workarea" ref="workarea" class="relative h-full w-full overflow-hidden rounded-xl bg-white/1 shadow-inner">
        <div ref="host" class="host"></div>
    </div>
    <input ref="textInput" style="width: 0; height: 0; opacity: 0" />
    <input id="zoom" type="hidden" style="width: 0; height: 0; opacity: 0" />
</template>

<style scoped>
div#workarea,
div.host {
    block-size: 100%;
    inline-size: 100%;
}
</style>
