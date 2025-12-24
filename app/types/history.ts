import type { File } from '~~/gen/ts/resources/file/file';

export interface Version<TContent> {
    id: number;
    date: string;
    type: string;
    name?: string;
    content: TContent;
}

export interface Content {
    title?: string;
    content: string;
    files: File[]; // Associated files, if any
}
