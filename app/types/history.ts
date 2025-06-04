import type { DocumentAccess } from '~~/gen/ts/resources/documents/access';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { DocumentReference, DocumentRelation } from '~~/gen/ts/resources/documents/documents';
import type { File } from '~~/gen/ts/resources/file/file';
import type { QualificationAccess } from '~~/gen/ts/resources/qualifications/access';
import type { QualificationRequirement } from '~~/gen/ts/resources/qualifications/qualifications';
import type { PageAccess } from '~~/gen/ts/resources/wiki/access';

export interface Version<TContent, TMeta> {
    id: string;
    type: string;
    content: TContent;
    meta: TMeta;
    name?: string;
}

export interface DocumentContent {
    content: string;
    files: File[];
}

export interface DocumentMeta {
    title: string;
    state: string;
    category: Category | undefined;
    closed: boolean;
    access: DocumentAccess;
    references: DocumentReference[];
    relations: DocumentRelation[];
}

export interface QualificationContent {
    content: string;
    files: File[];
    requirements: QualificationRequirement[];
}

export interface QualificationMeta {
    weight: number;
    abbreviation: string;
    title: string;
    description: string | undefined;
    closed: boolean;
    public: boolean;
    access: QualificationAccess;
}

export interface WikiContent {
    content: string;
    files: File[];
}

export interface WikiMeta {
    parentId: number;
    meta: {
        title: string;
        description: string | undefined;
        public: boolean;
        toc: boolean | undefined;
    };
    access: PageAccess;
}
