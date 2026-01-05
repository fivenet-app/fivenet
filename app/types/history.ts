import type { File } from '~~/gen/ts/resources/file/file';

export interface Version<TContent> {
    id: number;
    date: string;
    type: string;
    name?: string;
    content: TContent;
}

export interface HistoryContent {
    title?: string;
    content: EditorDocument | string | undefined;
    files: File[]; // Associated files, if any
}
