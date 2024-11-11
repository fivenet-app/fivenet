import { listEnumValues } from '@protobuf-ts/runtime';

export type JobAccessEntry = {
    id: string;
    userId?: number;
    job?: string;
    minimumGrade?: number;
    access: number;
    required?: boolean;
};

export type UserAccessEntry = {
    id: string;
    userId?: number;
    access: number;
    required?: boolean;
};

export type AccessEntryType = 'user' | 'job';

export type MixedAccessEntry = {
    id: string;
    type: AccessEntryType;
    userId?: number;
    job?: string;
    minimumGrade?: number;
    access: number;
    required?: boolean;
};

export type AccessType = {
    type: AccessEntryType;
    name: string;
};

export type AccessLevelEnum = {
    label: string;
    value: number;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function enumToAccessLevelEnums(accessLevel: any, translationPrefix: string): AccessLevelEnum[] {
    const { t } = useI18n();

    return [
        ...listEnumValues(accessLevel)
            .filter((e) => e.number !== 0)
            .map((e) => {
                return {
                    label: t(`${translationPrefix}.${e.name}`),
                    value: e.number,
                };
            }),
    ];
}
