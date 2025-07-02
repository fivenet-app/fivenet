import { ResultStatus, type QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserLike } from './strings';

// Access check helper function

type JobAccess<L> = {
    job: string;
    minimumGrade: number;
    access: L;
};

type UserAccess<L> = {
    userId: number;
    access: L;
};

type QualificationAccess<L> = {
    qualificationId: number;
    qualification?: QualificationShort;
    access: L;
};

export function checkAccess<L = number>(
    activeChar: UserLike,
    access:
        | { jobs?: JobAccess<L>[]; users?: UserAccess<L>[]; qualifications?: QualificationAccess<L>[] }
        | { jobs?: JobAccess<L>[]; users?: UserAccess<L>[]; qualifications?: QualificationAccess<L>[] }
        | undefined,
    creator: UserLike | undefined,
    level: L,
    creatorJob?: string,
): boolean {
    if (access === undefined) {
        return false;
    }

    if (creator !== undefined && activeChar.userId === creator.userId && activeChar.job === (creatorJob ?? creator.job)) {
        return true;
    }

    const ju = access.users?.find((ua) => ua.userId === activeChar.userId && level <= ua.access);
    if (ju !== undefined) {
        return true;
    }

    if (access.jobs !== undefined) {
        let lowestAccess: L | undefined = undefined;
        for (let index = 0; index < access.jobs?.length; index++) {
            const ja = access.jobs[index]!;
            if (ja.job !== activeChar.job) {
                continue;
            }
            if (ja.minimumGrade > activeChar.jobGrade) {
                continue;
            }
            if (ja.access < level) {
                continue;
            }
            if (lowestAccess === undefined || ja.access < lowestAccess!) {
                lowestAccess = ja.access;
            }
        }

        if (level <= (lowestAccess ?? 0)) {
            return true;
        }
    }

    if (access.qualifications !== undefined) {
        for (let index = 0; index < access.qualifications.length; index++) {
            const jq = access.qualifications[index]!;

            if (jq.qualification === undefined || jq.qualification.result === undefined) {
                continue;
            }

            if (jq.qualification.result.status !== ResultStatus.SUCCESSFUL) {
                continue;
            }

            if (level <= jq?.access) {
                return true;
            }
        }
    }

    return false;
}
