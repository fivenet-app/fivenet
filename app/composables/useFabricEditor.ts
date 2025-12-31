import {
    type ActiveSelection,
    Canvas,
    type CanvasEvents,
    type CanvasOptions,
    Circle,
    FabricImage,
    type FabricObject,
    Group,
    Pattern,
    Point,
    Rect,
    Textbox,
    util,
} from 'fabric';
import { AligningGuidelines } from 'fabric/extensions';
import { ref, shallowRef } from 'vue';
import '~/composables/fabric/FabricHtmlInput';
import { FabricCurvedText, type FabricCurvedTextOptions } from './fabric/FabricCurvedText';

export const formatPresets = {
    // 72 DPI
    A4: { width: 595.28, height: 841.89 },
    // 300 DPI
    A4_300DPI: { width: 2480, height: 3508 },
    A5: { width: 419.53, height: 595.28 },
    // US Letter
    Letter: { width: 612, height: 792 },
};

export const svgPatterns = [
    {
        name: 'None',
        value: undefined,
    },
    {
        name: 'Checkerboard',
        value: 'checkerboard',
        url: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAwAAAAMCAYAAABWdVznAAAAV0lEQVR42mNggAIwEqmOgWjYGBgY+DDqACrEJrA2gUGRZDjCJ8U4oHqggIogxYMoNIBaGM6KPAkKQpQBaGYZNEILUwEgwULYGhBWhEAJmvMZ0ElspkAAAAASUVORK5CYII=',
    },
    {
        name: 'Noise Dots',
        value: 'noise',
        svg: (color: string) => `<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8">
<rect fill="#3b3b3b" height="8" width="8"/>
<circle cx="1" cy="1" fill="${color}" r="1"/>
<circle cx="6" cy="3" fill="${color}" r="1"/>
<circle cx="3" cy="6" fill="${color}" r="1"/>
</svg>`,
    },
    {
        name: 'Diagonal Stripes',
        value: 'stripes',
        svg: (color: string) => `<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10">
<line x1="0" y1="0" x2="10" y2="10" stroke="${color}" stroke-width="2"/>
</svg>`,
    },
    {
        name: 'Crosshatch',
        value: 'crosshatch',
        svg: (color: string) => `<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10">
<line x1="0" y1="0" x2="10" y2="10" stroke="${color}" stroke-width="2"/>
<line x1="0" y1="10" x2="10" y2="0" stroke="${color}" stroke-width="2"/>
</svg>`,
    },
    {
        name: 'Dotted Grid',
        value: 'dots',
        svg: (color: string) => `<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10">
<circle cx="5" cy="5" r="1.5" fill="${color}"/>
</svg>`,
    },
];

export const strokeDashes = [
    { name: 'Solid', value: null },
    { name: 'Large Dashes', value: [10, 5] },
    { name: 'Medium Dashes', value: [5, 3] },
    { name: 'Small Dashes', value: [1, 1] },
];

let instance: ReturnType<typeof createFabricEditor> | null = null;

function createFabricEditor() {
    // State
    const canvas = shallowRef<Canvas | null>(null);
    const canvasEl = ref<HTMLCanvasElement | null>(null);
    const canvasContainer = ref<HTMLElement | null>(null);
    const activeObject = ref<FabricObject | null | undefined>(undefined);
    const zoom = ref(1);

    const canvasWidth = ref(800);
    const canvasHeight = ref(600);
    const documentSize = ref({
        width: 800,
        height: 600,
        fill: '#ffffff',
        disabled: false,
    });

    const history = ref<string[]>([]);
    const redoStack = ref<string[]>([]);

    // Alignment guide snapping threshold (in pixels)
    const snapThreshold = ref(5);

    const isDragging = ref(false);
    const lastPosX = ref(0);
    const lastPosY = ref(0);
    const exportedJSON = ref<string>('');
    const exportedSVG = ref<string>('');

    const pickingColor = ref(false);
    const pickedColor = ref<string | null>(null);

    const clipboard = ref<FabricObject | null>(null);

    const borderBox = shallowRef<Rect | null>(null);

    // Methods
    const createBorder = () => {
        const box = new Rect({
            left: 0,
            top: 0,
            width: documentSize.value.width,
            height: documentSize.value.height,
            stroke: '#999',
            strokeWidth: 1,
            fill: '',
            strokeDashOffset: 4,
            selectable: false,
            evented: false,
            excludeFromExport: true,
            name: 'document-border',
        });

        canvas.value?.add(box);
        canvas.value?.sendObjectToBack(box);

        borderBox.value = box;
    };

    const initCanvas = (canvasContainerElement: HTMLCanvasElement, opts: Partial<CanvasOptions>) => {
        const fabricCanvas = new Canvas(canvasContainerElement, {
            backgroundColor: '#ffffff00', // Transparent
            selection: true,
            preserveObjectStacking: true,
            ...opts,
        });
        canvas.value = fabricCanvas;
        canvasEl.value = canvasContainerElement;
        canvasContainer.value = canvasContainerElement.parentElement?.parentElement || null;

        const aligningGuidelines = new AligningGuidelines(fabricCanvas, {
            margin: 4,
            width: 1,
            color: 'rgba(0,0,0,0.85)',
        });

        const { width: containerWidth, height: containerHeight } = useElementSize(canvasContainer.value);

        const resizeCanvas = () => {
            if (!canvas.value) return;

            // Set the canvas element to match the full container size visually
            const containerEl = canvasContainer.value;
            if (!containerEl) return;

            // Keep track of element and window sizes
            const ratio = window.devicePixelRatio || 1;
            const width = containerWidth.value;
            const height = containerHeight.value;

            const canvasEl = canvas.value.getElement();
            canvasEl.width = width * ratio;
            canvasEl.height = height * ratio;
            canvasEl.style.width = width + 'px';
            canvasEl.style.height = height + 'px';

            canvas.value.renderAll();
        };

        // Enable alignment guide snapping and selection events
        fabricCanvas.on('object:moving', onObjectMoving);
        // Track active object
        fabricCanvas.on('selection:created', (e) => {
            activeObject.value = e.selected?.[0] || null;
        });
        fabricCanvas.on('selection:updated', (e) => {
            activeObject.value = e.selected?.[0] || null;
        });
        fabricCanvas.on('selection:cleared', () => {
            activeObject.value = null;
        });
        // History handling
        fabricCanvas.on('object:modified', saveHistory);
        fabricCanvas.on('object:added', saveHistory);
        fabricCanvas.on('object:removed', saveHistory);

        saveHistory();

        createBorder();

        // Zoom with mouse wheel
        canvasContainer.value?.addEventListener('wheel', (evt) => {
            const delta = (evt as WheelEvent).deltaY;
            let currentZoom = fabricCanvas.getZoom();
            currentZoom *= 0.999 ** delta;
            currentZoom = Math.min(3, Math.max(0.1, currentZoom));
            fabricCanvas.zoomToPoint(new Point({ x: evt.offsetX, y: evt.offsetY }), currentZoom);
            evt.preventDefault();
            evt.stopPropagation();
            zoom.value = currentZoom;
        });

        // Panning with middle-click or ctrl
        canvasContainer.value?.addEventListener('mousedown', (e) => {
            if (!canvas.value?.getActiveObject() && (e.button === 1 || e.ctrlKey)) {
                isDragging.value = true;
                lastPosX.value = e.clientX;
                lastPosY.value = e.clientY;
                // Disable selection
                if (canvas.value) canvas.value.selection = false;
                canvas.value?.setCursor('grab');
            } else if (canvas.value && pickingColor.value) {
                const { x, y } = canvas.value.getViewportPoint(e);
                const ctx = canvas.value.getContext();
                const pixel = ctx.getImageData(x, y, 1, 1).data;
                const hex = `#${[...pixel]
                    .slice(0, 3)
                    .map((n) => n.toString(16).padStart(2, '0'))
                    .join('')}`;
                pickedColor.value = hex;
                pickingColor.value = false;

                copyToClipboardWrapper(hex);
            }
        });

        canvasContainer.value?.addEventListener('mousemove', (e) => {
            if (isDragging.value && canvas.value) {
                const zoom = canvas.value.getZoom();
                const vpt = canvas.value.viewportTransform!;
                vpt[4] += (e.clientX - lastPosX.value) / zoom;
                vpt[5] += (e.clientY - lastPosY.value) / zoom;
                lastPosX.value = e.clientX;
                lastPosY.value = e.clientY;
                canvas.value.requestRenderAll();
            }
        });

        canvasContainer.value?.addEventListener('mouseup', () => {
            isDragging.value = false;
            // Re-enable selection
            if (canvas.value) canvas.value.selection = true;
            canvas.value?.setCursor('default');
        });

        // Register keyboard events for shortcuts
        window.addEventListener('keydown', handleKeydown);

        watch(canvas, (c) => {
            if (!c) return;

            canvasWidth.value = c.getWidth();
            canvasHeight.value = c.getHeight();
        });

        watchThrottled([containerWidth, containerHeight], () => resizeCanvas(), {
            leading: true,
            trailing: true,
            throttle: 200,
        });

        watchThrottled(
            documentSize,
            () => {
                if (!canvas.value) return;

                borderBox.value?.set({
                    width: documentSize.value.width,
                    height: documentSize.value.height,
                });
                if (borderBox.value) {
                    borderBox.value.fill = documentSize.value.fill;
                }
                canvas.value.renderAll();
            },
            { immediate: true, leading: true, trailing: true, throttle: 200, deep: true },
        );

        watch(pickingColor, (newVal) => {
            // Enable/Disable selection
            if (canvas.value) canvas.value.selection = !newVal;
        });

        onUnmounted(() => {
            if (!canvas.value) return;

            // Cleanup event listeners and canvas instance on destroy
            window.removeEventListener('keydown', handleKeydown);

            aligningGuidelines.dispose();
            canvas.value.dispose();
            canvas.value = null;
        });
    };

    const isInputFocused = () => {
        const tag = document.activeElement?.tagName.toLowerCase();
        return tag === 'input' || tag === 'textarea' || document.activeElement?.getAttribute('contenteditable') === 'true';
    };

    // Keyboard event handler for shortcuts
    async function handleKeydown(e: KeyboardEvent) {
        if (!canvas.value || isInputFocused()) return;

        const active = canvas.value.getActiveObject();
        // Delete or Backspace: remove selected object
        if ((e.key === 'Delete' || e.key === 'Backspace') && active) {
            canvas.value.remove(...canvas.value.getActiveObjects());
            canvas.value.discardActiveObject();
            canvas.value.renderAll();
            activeObject.value = null;
            e.preventDefault();
            return;
        }

        // Arrow keys: nudge object by 5px (with Shift for 1px for fine tuning)
        const step = e.shiftKey ? 1 : 5;
        if (active) {
            switch (e.key) {
                case 'ArrowLeft':
                    active.left! -= step;
                    active.setCoords();
                    canvas.value.renderAll();
                    e.preventDefault();
                    return;

                case 'ArrowRight':
                    active.left! += step;
                    active.setCoords();
                    canvas.value.renderAll();
                    e.preventDefault();
                    return;

                case 'ArrowUp':
                    active.top! -= step;
                    active.setCoords();
                    canvas.value.renderAll();
                    e.preventDefault();
                    return;

                case 'ArrowDown':
                    active.top! += step;
                    active.setCoords();
                    canvas.value.renderAll();
                    e.preventDefault();
                    return;
            }
        }

        const isMac = /Mac|iPod|iPhone|iPad/.test(navigator.platform);
        const ctrlKey = isMac ? e.metaKey : e.ctrlKey;

        if (ctrlKey) {
            // Copy (Ctrl/Cmd+C)
            if (e.key === 'c') {
                if (active) {
                    clipboard.value = await active.clone();
                    e.preventDefault();
                }
            }

            // Paste (Ctrl/Cmd+V)
            if (e.key === 'v' && clipboard.value) {
                const clonedObj = await clipboard.value.clone();
                canvas.value!.discardActiveObject();
                // Offset the pasted object slightly
                clonedObj.set({ left: clonedObj.left! + 20, top: clonedObj.top! + 20 });
                canvas.value!.add(clonedObj);
                canvas.value!.setActiveObject(clonedObj);
                e.preventDefault();
                return;
            }

            // Duplicate (Ctrl/Cmd+D)
            if (e.key === 'd' && active) {
                const clonedObj = await active.clone();
                clonedObj.set({ left: active.left! + 20, top: active.top! + 20 });
                canvas.value!.add(clonedObj);
                canvas.value!.setActiveObject(clonedObj);
                e.preventDefault();
                return;
            }

            if (e.key === 'z') {
                e.preventDefault();
                if (e.shiftKey) {
                    redo();
                } else {
                    undo();
                }
            } else if (e.key === 'y') {
                e.preventDefault();
                redo();
            }
        }
    }

    // Alignment guides: snap moving object to canvas center or other objects
    function onObjectMoving(e: CanvasEvents['object:moving']) {
        const obj = e.target as import('fabric').Object;
        if (!obj) return;
        let snapped = false;
        // Snap to canvas center (horizontal & vertical)
        const canvasWidth = canvas.value!.getWidth();
        const canvasHeight = canvas.value!.getHeight();
        const objCenterX = obj.left! + (obj.width! * obj.scaleX!) / 2;
        const objCenterY = obj.top! + (obj.height! * obj.scaleY!) / 2;
        const canvasCenterX = canvasWidth / 2;
        const canvasCenterY = canvasHeight / 2;
        if (Math.abs(objCenterX - canvasCenterX) < snapThreshold.value) {
            obj.left = canvasCenterX - (obj.width! * obj.scaleX!) / 2;
            snapped = true;
        }
        if (Math.abs(objCenterY - canvasCenterY) < snapThreshold.value) {
            obj.top = canvasCenterY - (obj.height! * obj.scaleY!) / 2;
            snapped = true;
        }
        // Snap to other objects' left/top edges and centers
        canvas.value!.getObjects().forEach((target) => {
            if (target === obj) return;
            // Snap left edges
            if (Math.abs(obj.left! - target.left!) < snapThreshold.value) {
                obj.left = target.left!;
                snapped = true;
            }
            // Snap top edges
            if (Math.abs(obj.top! - target.top!) < snapThreshold.value) {
                obj.top = target.top!;
                snapped = true;
            }
            // Snap horizontal centers
            const targetCenterX = target.left! + (target.width! * target.scaleX!) / 2;
            if (Math.abs(objCenterX - targetCenterX) < snapThreshold.value) {
                obj.left = targetCenterX - (obj.width! * obj.scaleX!) / 2;
                snapped = true;
            }
            // Snap vertical centers
            const targetCenterY = target.top! + (target.height! * target.scaleY!) / 2;
            if (Math.abs(objCenterY - targetCenterY) < snapThreshold.value) {
                obj.top = targetCenterY - (obj.height! * obj.scaleY!) / 2;
                snapped = true;
            }
        });
        if (snapped) {
            obj.setCoords(); // update object bounds
            canvas.value!.renderAll();
        }
    }

    const saveHistory = () => {
        if (canvas.value) {
            const json = canvas.value.toJSON();
            history.value.push(JSON.stringify(json));
            redoStack.value = [];
        }
    };

    const undo = async () => {
        if (!canvas.value || history.value.length < 2) return;
        const current = history.value.pop();
        if (current) redoStack.value.push(current);
        const previous = history.value[history.value.length - 1];
        if (previous) {
            await canvas.value.loadFromJSON(previous);
            canvas.value.requestRenderAll();
        }
    };

    const redo = async () => {
        if (!canvas.value || redoStack.value.length === 0) return;
        const json = redoStack.value.pop();
        if (json) {
            history.value.push(json);
            await canvas.value.loadFromJSON(json);
            canvas.value.requestRenderAll();
        }
    };

    function addText() {
        if (!canvas.value) return;
        const text = new Textbox('New Text', {
            left: 50,
            top: 50,
            width: 150,
            fontSize: 24,
            fill: '#000000',
        });
        canvas.value.add(text);
        canvas.value.setActiveObject(text);
        // Focus text for editing
        if (text.enterEditing) text.enterEditing();
    }

    function addCurvedText(text: string, radius: number = 100, opts: FabricCurvedTextOptions = {}) {
        if (!canvas.value) return;

        const curved = new FabricCurvedText(text, radius, opts);
        canvas.value.add(curved);
        canvas.value.setActiveObject(curved);
    }

    function addPlaceholder() {
        if (!canvas.value) return;
        const placeholderText = new Textbox('{{placeholder}}', {
            left: 50,
            top: 100,
            width: 180,
            fontSize: 24,
            fill: '#d97706' /* amber-600 */,
            fontStyle: 'italic',
        });
        // Mark as placeholder for identification
        (placeholderText as Textbox & { isPlaceholder?: boolean }).isPlaceholder = true;
        canvas.value.add(placeholderText);
        canvas.value.setActiveObject(placeholderText);
        if (placeholderText.enterEditing) placeholderText.enterEditing();
    }

    function addRectangle() {
        if (!canvas.value) return;
        const rect = new Rect({
            left: 80,
            top: 80,
            width: 100,
            height: 60,
            fill: '#93c5fd',
            stroke: '#3b82f6',
            strokeWidth: 1,
        });
        canvas.value.add(rect);
        canvas.value.setActiveObject(rect);
    }

    function addCircle() {
        if (!canvas.value) return;
        const circle = new Circle({
            left: 200,
            top: 80,
            radius: 50,
            fill: '#fca5a5',
            stroke: '#f87171',
            strokeWidth: 1,
        });
        canvas.value.add(circle);
        canvas.value.setActiveObject(circle);
    }

    function addImage(src: string) {
        if (!canvas.value) return;
        const img = new Image();
        img.src = src;
        img.onload = () => {
            const fi = new FabricImage(img, {
                left: 100,
                top: 100,
            });
            canvas.value?.add(fi);
            canvas.value?.setActiveObject(fi);
            canvas.value?.requestRenderAll();
        };
    }

    function exportJSON() {
        if (!canvas.value) return;
        // Include custom placeholder flag in JSON output
        const json = canvas.value.toDatalessJSON(['isPlaceholder']);
        exportedJSON.value = JSON.stringify(json);
        console.log('Canvas JSON:', exportedJSON.value); // (For debugging or external use)
    }

    function exportSVG() {
        if (!canvas.value) return;
        exportedSVG.value = canvas.value.toSVG({
            width: documentSize.value.width + 'px',
            height: documentSize.value.height + 'px',
        });
        console.log('Canvas SVG:', exportedSVG.value);
    }

    watch(zoom, (newZoom) => {
        if (canvas.value) {
            canvas.value.setZoom(newZoom);
            canvas.value.requestRenderAll();
        }
    });

    const updateCanvasSize = () => {
        if (canvas.value) {
            canvas.value.setDimensions({ width: canvasWidth.value, height: canvasHeight.value });
            canvas.value.renderAll();
            fitDocumentToView();
        }
    };

    const fitDocumentToView = () => {
        if (!canvas.value) return;

        const c = canvas.value;

        const container = c.getElement().parentElement;
        if (!container) return;

        const containerWidth = container.clientWidth;
        const containerHeight = container.clientHeight;

        // Use actual canvas doc size (not canvas element size!)
        const docWidth = c.getWidth();
        const docHeight = c.getHeight();

        const z = Math.min(containerWidth / docWidth, containerHeight / docHeight) * 0.9; // add some margin

        const offsetX = (containerWidth - docWidth * z) / 2;
        const offsetY = (containerHeight - docHeight * z) / 2;

        c.setZoom(z);
        c.viewportTransform = [z, 0, 0, z, offsetX, offsetY];

        c.requestRenderAll();
        zoom.value = z;
    };

    function resetZoom() {
        if (canvas.value) {
            zoom.value = 1;
        }
    }

    const copySelected = async () => {
        if (!canvas.value) return;

        const obj = canvas.value.getActiveObject();
        if (obj) {
            clipboard.value = await obj.clone();
        }
    };

    const cutSelected = async () => {
        if (!canvas.value) return;

        const obj = canvas.value.getActiveObject();
        if (obj) {
            clipboard.value = await obj.clone();
            canvas.value?.remove(obj);
            canvas.value?.discardActiveObject();
            canvas.value?.requestRenderAll();
        }
    };

    const paste = async () => {
        if (!canvas.value || !clipboard.value) return;

        const clonedObj = await clipboard.value.clone();
        canvas.value?.discardActiveObject();
        clonedObj.set({
            left: (clonedObj.left || 0) + 20,
            top: (clonedObj.top || 0) + 20,
            evented: true,
        });

        if (clonedObj instanceof Group) {
            const newGroup = await clonedObj.clone();
            canvas.value?.add(newGroup);
            canvas.value?.setActiveObject(newGroup);
            canvas.value?.requestRenderAll();
        } else {
            canvas.value?.add(clonedObj);
            canvas.value?.setActiveObject(clonedObj);
            canvas.value?.requestRenderAll();
        }
    };

    const applyPatternFill = async (patternType: string, color = '#000000') => {
        console.debug('Applying pattern fill:', patternType, activeObject.value?.type);
        if (!canvas.value || !activeObject.value || activeObject.value.type !== 'rect') return;

        const svgPattern = svgPatterns.find((p) => p.value === patternType);
        if (!svgPattern) return;

        if (svgPattern.url) {
            const img = new Image();
            img.src = svgPattern.url;
            img.onload = () => {
                const pattern = new Pattern({ source: img, repeat: 'repeat' });
                activeObject.value?.set('fill', pattern);
                canvas.value?.requestRenderAll();
            };
        } else if (svgPattern.svg) {
            const svg = svgPattern.svg(color);

            const blob = new Blob([svg], { type: 'image/svg+xml' });
            const url = URL.createObjectURL(blob);
            const img = await util.loadImage(url);
            const pattern = new Pattern({
                source: img,
                repeat: 'repeat',
            });
            activeObject.value?.set('fill', pattern);
            canvas.value?.requestRenderAll();
            URL.revokeObjectURL(url);
        }
    };

    function groupObject() {
        if (!activeObject.value || !canvas.value) return;

        if (activeObject.value.type !== 'activeSelection' && activeObject.value.type !== 'activeselection') return;

        const ao = activeObject.value as ActiveSelection;
        const group = new Group(ao.removeAll());
        canvas.value?.add(group);
        canvas.value?.setActiveObject(group);
        canvas.value?.requestRenderAll();
    }

    function ungroupObject() {
        if (!activeObject.value || !canvas.value) return;

        const ao = activeObject.value;
        if (!ao || !ao.group) return;

        canvas.value?.add(...(ao.group as Group).removeAll());

        canvas.value?.remove(ao.group as Group);
        ao.dispose();
        canvas.value?.requestRenderAll();
    }

    const bringForward = () => {
        if (canvas.value && activeObject.value) {
            canvas.value.bringObjectForward(activeObject.value as FabricObject);
            canvas.value.requestRenderAll();
        }
    };

    const sendBackward = () => {
        if (canvas.value && activeObject.value) {
            canvas.value.sendObjectBackwards(activeObject.value as FabricObject);
            canvas.value.requestRenderAll();
        }
    };

    const bringToFront = () => {
        if (canvas.value && activeObject.value) {
            canvas.value.bringObjectToFront(activeObject.value as FabricObject);
            canvas.value.requestRenderAll();
        }
    };

    const sendToBack = () => {
        if (canvas.value && activeObject.value) {
            canvas.value.sendObjectToBack(activeObject.value as FabricObject);
            canvas.value.requestRenderAll();
        }
    };

    // Return state and methods
    return {
        canvas,
        canvasContainer,
        activeObject,
        zoom,
        documentSize,
        snapThreshold,
        exportedJSON,
        exportedSVG,
        pickingColor,
        pickedColor,

        initCanvas,
        addText,
        addCurvedText,
        addPlaceholder,
        addRectangle,
        addCircle,
        addImage,
        exportJSON,
        exportSVG,

        updateCanvasSize,
        fitDocumentToView,
        resetZoom,

        undo,
        redo,

        copySelected,
        cutSelected,
        paste,

        applyPatternFill,

        groupObject,
        ungroupObject,

        bringForward,
        sendBackward,
        bringToFront,
        sendToBack,
    };
}

export function useFabricEditor() {
    if (!instance) {
        instance = createFabricEditor();
    }
    return instance;
}
