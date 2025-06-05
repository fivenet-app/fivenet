import type { File } from '~~/gen/ts/resources/file/file';

export interface Version<TContent> {
    id: string;
    type: string;
    name?: string;
    content: TContent;
}

export interface Content {
    content: string;
    files: File[]; // Associated files, if any
}
