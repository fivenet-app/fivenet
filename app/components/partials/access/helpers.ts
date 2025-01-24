import { listEnumValues } from '@protobuf-ts/runtime';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';

export type JobAccessEntry = {
    id: number;
    targetId: number;
    job?: string;
    minimumGrade?: number;
    access: number;
    required?: boolean;
    jobLabel?: string;
    jobGradeLabel?: string;
};

export type UserAccessEntry = {
    id: number;
    targetId: number;
    userId?: number;
    access: number;
    required?: boolean;
    user?: UserShort;
};

export type QualificationAccessEntry = {
    id: number;
    targetId: number;
    qualificationId?: number;
    access: number;
    required?: boolean;
    qualification?: QualificationShort;
};

export type AccessEntryType = 'user' | 'job' | 'qualification';

export type MixedAccessEntry = {
    id: number;
    type: AccessEntryType;

    userId?: number;
    user?: UserShort;

    job?: string;
    minimumGrade?: number;

    qualificationId?: number;
    qualification?: QualificationShort;

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
