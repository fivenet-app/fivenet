export type Kind = 'text' | 'field' | 'image' | 'repeat' | 'checkbox' | 'grid' | 'line' | 'section';

export interface BaseFrame {
    id: string;
    name?: string;
    xMm: number;
    yMm: number;
    wMm: number;
    hMm: number;
    kind: Kind;
    strokeColor?: string;
    strokeWidth?: number;
    strokeEnabled?: boolean;
    fill?: string;
}

export interface TextFrame extends BaseFrame {
    kind: 'text';
    text: string;
    fontSize?: number; // pt
    bold?: boolean;
    italic?: boolean;
    underline?: boolean;
    align?: 'left' | 'center' | 'right';
    rotateDeg?: number;
    style?: string; // custom CSS or style string
}

export interface FieldFrame extends BaseFrame {
    kind: 'field';
    path: string;
    fallback?: string;
}

export interface ImageFrame extends BaseFrame {
    kind: 'image';
    src?: string;
    fit: 'contain' | 'cover' | 'stretch';
}

export interface RepeatFrame extends BaseFrame {
    kind: 'repeat';
    path: string;
}

export interface CheckboxFrame extends BaseFrame {
    kind: 'checkbox';
    path?: string;
    label?: string;
    checked?: boolean;
}

export interface GridFrame extends BaseFrame {
    kind: 'grid';
    cols?: number;
    rows?: number;
    gapMm?: number;
}

export interface LineFrame extends BaseFrame {
    kind: 'line';
}

export interface SectionFrame extends BaseFrame {
    kind: 'section';
    title?: string;
}

export type Frame = TextFrame | FieldFrame | ImageFrame | RepeatFrame | CheckboxFrame | GridFrame | LineFrame | SectionFrame;
