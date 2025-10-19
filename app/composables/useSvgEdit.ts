/* eslint-disable @typescript-eslint/no-explicit-any */
import { onBeforeUnmount, onMounted, ref } from 'vue';

export interface UseSvgEditOptions {
    container: Ref<HTMLElement | null>;
    workarea: Ref<HTMLElement | null>;
    config?: Record<string, any>;
    initialSvg?: string;
}

export function useSvgEdit(opts: UseSvgEditOptions) {
    const containerEl = opts.container;
    const workareaEl = opts.workarea;
    const canvas = ref<any | null>(null);
    const ready = ref(false);
    const error = ref<unknown | null>(null);
    const selection = ref<any[]>([]);

    let SvgCanvasCtor: any;

    const containerSize = useElementSize(containerEl);

    onMounted(async () => {
        try {
            const mod = await import('@svgedit/svgcanvas');
            SvgCanvasCtor = mod.default || (mod as any);

            if (!containerEl.value) return;

            canvas.value = new SvgCanvasCtor(containerEl.value, {
                initFill: { color: '000000', opacity: 1 },
                initStroke: { color: '000000', width: 1, opacity: 1 },
                dimensions: [containerSize.width.value, containerSize.height.value],
                baseUnit: 'px',
                ...opts.config,
            });

            watch(
                () => [containerSize.width.value, containerSize.height.value],
                ([width, height]) => canvas.value?.updateCanvas(width, height),
            );

            workareaEl.value?.addEventListener('wheel', (e) => {
                if (!e.altKey || !canvas.value) return;

                e.preventDefault();
                canvas.value.setZoom(e.deltaY > 0 ? canvas.value.getZoom() * 0.9 : canvas.value.getZoom() * 1.1, true);
                canvas.value.call('updateCanvas', { center: true });
            });

            // Example event hooks (svgcanvas exposes custom events via callbacks)
            canvas.value.bind('selected', (els: any[]) => {
                selection.value = els;
            });
            canvas.value.bind('changed', () => {
                /* noop: consumers can listen */
            });

            canvas.value.setMode('pan');

            if (opts.initialSvg) {
                canvas.value.setSvgString(opts.initialSvg, true);
            }

            ready.value = true;
        } catch (e) {
            error.value = e;
        }
    });

    onBeforeUnmount(() => {
        try {
            canvas.value?.destroy?.();
        } catch {
            // Ignore cleanup errors
        }
    });

    // Commands (wrap the common API surface)
    const commands = {
        setMode(mode: string) {
            canvas.value?.setMode(mode);
        },
        addRect(attrs: Record<string, any> = {}) {
            return canvas.value?.addSvgElementFromJson?.({ element: 'rect', attr: attrs });
        },
        addEllipse(attrs: Record<string, any> = {}) {
            return canvas.value?.addSvgElementFromJson?.({ element: 'ellipse', attr: attrs });
        },
        addLine(attrs: Record<string, any> = {}) {
            return canvas.value?.addSvgElementFromJson?.({ element: 'line', attr: attrs });
        },
        addText(text: string, attrs: Record<string, any> = {}) {
            return canvas.value?.addSvgElementFromJson?.({
                element: 'text',
                attr: { 'xml:space': 'preserve', text, ...attrs },
            });
        },
        deleteSelected() {
            canvas.value?.deleteSelectedElements?.();
        },
        setFill(color: string) {
            canvas.value?.setColor?.('fill', color);
        },
        setStroke(color: string) {
            canvas.value?.setColor?.('stroke', color);
        },
        setStrokeWidth(w: number) {
            canvas.value?.setStrokeWidth?.(w);
        },
        setOpacity(v: number) {
            canvas.value?.setOpacity?.(v);
        },

        toString(): Promise<string> {
            return Promise.resolve(canvas.value?.getSvgString?.());
        },
        load(svg: string) {
            canvas.value?.setSvgString?.(svg);
        },
        clear() {
            canvas.value?.clear?.();
        },
        setZoom(v: number) {
            canvas.value?.setZoom?.(v);
        },
        getZoom(): number {
            return canvas.value?.zoom ?? 1;
        },
        undo() {
            canvas.value?.history?.undoMgr?.undo?.();
        },
        redo() {
            canvas.value?.history?.undoMgr?.redo?.();
        },
    };

    return { containerEl, workareaEl, canvas, ready, error, selection, ...commands };
}
