/* eslint-disable @typescript-eslint/no-explicit-any */
declare module '@svgedit/svgcanvas' {
    // Type for event handler callbacks used in SvgCanvas events
    export type EventHandler = (win: Window, arg: any) => any;

    // Interface for shape style properties (fill, stroke, opacity, etc.)
    export interface StyleOptions {
        fill: string;
        fill_opacity: number;
        stroke: string;
        stroke_width: number;
        stroke_dasharray: string;
        stroke_linejoin: string;
        stroke_linecap: string;
        stroke_opacity: number;
        opacity: number;
        // Additional style properties can be included as needed
    }

    // The main SVG canvas class responsible for SVG editing operations
    export default class SvgCanvas {
        /**
         * Creates a new SvgCanvas.
         * @param container The container HTML element that should hold the SVG root element.
         * @param config An object that contains configuration data for the editor (SVG-Edit settings).
         */
        constructor(container: HTMLElement, config: any);

        // Float displaying the current zoom level (1 = 100%, 0.5 = 50%, etc).
        zoom: number;

        // Pointer to current group (for in-group editing). Null if no group context.
        currentGroup: SVGElement | null;

        // Object containing data for the currently selected styles (e.g., fill, stroke of selected element).
        curProperties: Partial<StyleOptions>;

        // Current text style properties (for newly created text elements).
        curText: StyleOptions;

        // Current shape style properties (for newly created shape elements).
        curShape: StyleOptions;

        // Array with all the currently selected elements. Default size of 1 until it needs to grow.
        selectedElements: Element[];

        // The root <svg> element of the drawing canvas.
        svgroot: SVGSVGElement;
        // The container SVG element (often same as svgroot, or a child of container HTML element).
        svgContent: SVGSVGElement;
        // The SVG document (for convenience).
        svgdoc: Document;
        // Container HTML element passed in (for reference).
        container: HTMLElement;

        // Various other internal state properties:
        started: boolean;
        startTransform: string;
        // Object containing save/export options (like round_digits).
        saveOptions: { [key: string]: any };
        // Undo/redo manager (history object with undo/redo methods).
        undoMgr: any;
        // The current Drawing (layers) object.
        current_drawing_: any;
        // Selector manager for selection outlines/handles.
        selectorManager: any;
        // Map of namespace URIs (e.g., NS.SVG, NS.HTML, etc.)
        static NS: Record<string, string>;

        // *** Editor control methods ***

        // Clears the current document. This is not an undoable action.
        clear(): void;

        /**
         * Attaches a callback function to an event.
         * @param ev  The name of the event (e.g., "selected", "changed", etc).
         * @param f   The callback function to bind to the event.
         * @returns   The previous event handler for this event, if any.
         */
        bind(ev: string, f: EventHandler): EventHandler;

        // Triggers an event’s callback, if it exists. Used internally to handle custom events.
        call(ev: string, arg: any): any;

        // Adds/updates an extension. If an extension with the given name exists, it throws an error.
        // `extInitFunc` is a function that initializes the extension and returns an object with extension data.
        addExtension(name: string, extInitFunc: (args: any) => any, options: { importLocale: any }): any;

        // Undo/Redo history methods
        addCommandToHistory(cmd: any): void;

        setMode(mode: string): void;

        // Returns the array of currently selected elements.
        getSelectedElements(): Element[];
        // Manually sets an entry in the selectedElements array (used internally for selection logic).
        setSelectedElements(index: number, element: Element | null): void;
        // Empties the selection array (deselects all elements).
        setEmptySelectedElements(): void;

        // Gets the root <svg> element of the canvas.
        getSvgRoot(): SVGSVGElement;
        // Gets the SVG document (SVG DOM root).
        getDOMDocument(): Document;
        // Gets the container HTML element.
        getDOMContainer(): HTMLElement;
        // (Alias) Gets the container HTML element.
        getContainer(): HTMLElement;

        // Returns the current configuration object (curConfig).
        getCurConfig(): any;
        // Update configuration options with given values.
        setConfig(opts: any): void;

        // Returns the current drawing (layer manager) object.
        getCurrentDrawing(): any;

        // Gets the current shape style properties.
        getCurShape(): StyleOptions;
        // Gets the currently active group (for in-group editing), or null if none.
        getCurrentGroup(): SVGElement | null;

        // Gets the base unit of the editor (e.g., "px", "mm").
        getBaseUnit(): string;

        // Returns an object with the current document's width and height.
        getResolution(): { w: number; h: number };

        // Returns the current height of the SVG canvas, accounting for zoom.
        getHeight(): number;
        // Returns the current width of the SVG canvas, accounting for zoom.
        getWidth(): number;

        // Returns the number of digits to round to for certain operations (from saveOptions.round_digits).
        getRoundDigits(): number;
        // Returns the grid snapping value.
        getSnappingStep(): number;
        // Returns whether grid snapping is enabled.
        getGridSnapping(): boolean;

        // Get or set the initial transform of the current element (used when moving elements).
        getStartTransform(): string;
        setStartTransform(transform: string): void;

        // Gets the current zoom level (1 = 100%).
        getZoom(): number;
        // Sets the zoom level for the canvas.
        setZoom(value: number): void;

        // Utility function to round a value to the nearest grid step or zoom unit.
        round(val: number): number;

        // Creates a new SVG element based on the given JSON map (convenience method).
        createSVGElement(jsonMap: object): Element;

        // Set or get editor "save options" (like whether to apply certain sanitizations on output).
        getSvgOption(): any;
        setSvgOption(key: string, value: any): void;

        // *** Element selection and editing methods ***

        // Selects only the given elements, replacing any current selection.
        selectOnly(elems: Element[], showGrips?: boolean): void;
        // Removes the given elements from the current selection.
        removeFromSelection(elemsToRemove: Element[]): void;
        // Selects all elements in the current layer.
        selectAllInCurrentLayer(): void;
        // Clears the selection (deselects all elements).
        clearSelection(force?: boolean): void;
        // Adds elements to the current selection.
        addToSelection(elems: Element[], showGrips?: boolean): void;

        // Returns the current fill color's opacity.
        getFillOpacity(): number;
        // Returns the current stroke's opacity.
        getStrokeOpacity(): number;
        // Gets whether "snap to grid" is enabled (same as getGridSnapping, provided for API).
        getSnapToGrid(): boolean;
        // Returns a string describing the revision or version of SvgCanvas.
        getVersion(): string;

        // Sets UI locale-specific strings (not generally needed by end-users of the API).
        setUiStrings(strs: { [key: string]: any }): void;

        // Returns the current document title (from the <title> element), or an empty string if not found.
        getDocumentTitle(): string | undefined;

        // Returns an object with x,y representing the SVG content’s offset (if any).
        getOffset(): { x: number; y: number };

        // Gets the current value of a style property (fill or stroke) for the current tool/element.
        getColor(type: string): string | number;
        // Sets the current stroke paint (color/gradient) for new strokes.
        setStrokePaint(paint: any): void;
        // Sets the current fill paint (color/gradient) for new fills.
        setFillPaint(paint: any): void;

        // Returns the current stroke width value.
        getStrokeWidth(): number | string;
        // Returns an object with the current style properties for new elements (e.g., fill, stroke, etc.).
        getStyle(): StyleOptions;

        // Sets the given overall opacity on the current selected elements.
        setOpacity(val: string): void;
        // Get the overall opacity of the first selected element.
        getOpacity(): number;

        // Sets the fill or stroke opacity to a new value (if preventUndo is false, it records an undo step).
        setPaintOpacity(type: 'fill' | 'stroke', val: number, preventUndo?: boolean): void;
        // Gets the fill or stroke opacity (depending on `type` argument).
        getPaintOpacity(type: 'fill' | 'stroke'): number;

        // Gets the blur value of the given element (returns the stdDeviation value of its filter).
        getBlur(elem: Element): string;
        // Sets a "last good image" URL (presumably for image elements that failed to load).
        setGoodImage(val: string): void;

        // Returns the current drawing as raw SVG XML text.
        getSvgString(): string;
        // Sets the current drawing from an SVG string (imports SVG XML into the canvas).
        setSvgString(xmlString: string): void;

        // Cuts the selected elements: removes them from the DOM and places them on the clipboard (undoable).
        cutSelectedElements(): void;
        // Copies the selected elements to the clipboard (without removing them).
        copySelectedElements(): void;
        // Deletes the selected elements from the DOM (undoable).
        deleteSelectedElements(): void;
        // Pastes elements from the clipboard to the current layer (at optionally a given location).
        pasteElements(): void;

        // Clears the SVG content element (empties the drawing area).
        clearSvgContentElement(): void;

        // Aligns selected elements (e.g., left, center, right, top, middle, bottom alignment relative to each other or artboard).
        alignSelectedElements(type: string, relativeTo?: string): void;
        // Group the selected elements into a single group element.
        groupSelectedElements(): SVGGElement;
        // Ungroup the selected group (if a single group is selected).
        ungroupSelectedElement(): void;

        // Additional methods for canvas operations (rotate, move, resize, etc.) would be listed here...
        // (For brevity, not all SvgCanvas methods are shown in this excerpt)

        // *** Static utility methods attached to SvgCanvas class ***

        // Get an element by its SVG id. Shortcut for document.getElementById.
        static $id(id: string): Element | null;
        // Query the DOM for a single element (within an optional parent).
        static $qq(selector: string, parent?: Element): Element | null;
        // Query the DOM for all elements matching the selector (within an optional parent). Returns a list of elements.
        static $qa(selector: string, parent?: Element): NodeListOf<Element>;
        // Attach a click event handler to an element (wrapper for addEventListener).
        static $click(element: Element, handler: (e: MouseEvent) => any): void;
        // Base64-encode a string.
        static encode64(input: string): string;
        // Base64-decode a string.
        static decode64(input: string): string;
        // Deep-merge multiple objects into the first object (similar to $.extend or lodash merge).
        static mergeDeep(target: object, ...sources: object[]): object;
        // Get the closest ancestor of an element with a given tag name.
        static getClosest(elem: Element, tagName: string): Element | null;
        // Get all ancestor elements of a given element (up to the root).
        static getParents(elem: Element): Element[];
        // A data URL string for a blank page (used for certain operations like exporting).
        static blankPageObjectURL: string;
    }
}
