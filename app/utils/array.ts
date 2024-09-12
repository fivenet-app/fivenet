import type { UserShort } from '~~/gen/ts/resources/users/users';

// GRPC Websocket helper functions

export function shuffle<T extends any[]>(arr: T): T {
    for (let i = arr.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [arr[i], arr[j]] = [arr[j], arr[i]];
    }

    return arr;
}

export function range(size: number, startAt?: number): number[] {
    return [...Array(size).keys()].map((i) => i + (startAt ?? 0));
}

export function writeUInt32BE(arr: Uint8Array, value: number, offset: number) {
    value = +value;
    offset = offset | 0;
    arr[offset] = value >>> 24;
    arr[offset + 1] = value >>> 16;
    arr[offset + 2] = value >>> 8;
    arr[offset + 3] = value & 0xff;
    return offset + 4;
}

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

export function checkAccess<L = number>(
    activeChar: UserShort,
    access: { jobs?: JobAccess<L>[]; users?: UserAccess<L>[] } | undefined,
    creator: UserShort | undefined,
    level: L,
): boolean {
    if (access === undefined) {
        return false;
    }

    if (creator !== undefined && activeChar.userId === creator.userId) {
        return true;
    }

    const ju = access.users?.find((ua) => ua.userId === activeChar.userId && level <= ua.access);
    if (ju !== undefined) {
        return true;
    }

    let lowestAccess: L | undefined = undefined;
    if (access.jobs === undefined) {
        return false;
    }
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

    console.log('level', level, 'lowestAccess', lowestAccess, activeChar.firstname, activeChar.jobGrade);

    return level <= (lowestAccess ?? 0);
}
