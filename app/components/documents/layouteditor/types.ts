// New: supported widgets
export type WidgetType = 'text' | 'textarea' | 'boxed-text';

// Base frame that already exists in your code
export interface BaseFrame {
    id: string;
    kind: 'widget'; // <- normalize to a single kind
    name: string;
    xMm: number;
    yMm: number;
    wMm: number;
    hMm: number;
    strokeColor?: string;
    strokeWidth?: number;
    strokeEnabled?: boolean;
    fill?: string;
}

// Binding & display options that a widget can carry
export interface DisplayBinding {
    path: string; // e.g. "patient.name"
    formatter?: string; // e.g. "upper", "date:DD.MM.YYYY"
    fit?: 'wrap' | 'truncate' | 'shrink';
    minFontSize?: number; // for shrink
}

// The one and only placeable “input” frame
export interface WidgetFrame extends BaseFrame {
    kind: 'widget';
    widget: WidgetType; // 'text' | 'textarea' | 'boxed-text'
    binding?: DisplayBinding; // optional until user binds
    // widget-specific options:
    rows?: number; // textarea
    boxCharCount?: number; // boxed-text
    letterSpacing?: number; // boxed-text visual gap (px)
    maxChars?: number; // logical cap for text/textarea/boxed
    placeholder?: string;
}

// If you had a big Frame union, replace input variants with WidgetFrame:
export type Frame = WidgetFrame /* | ImageFrame | LineFrame | etc... (non-input frames only) */;
