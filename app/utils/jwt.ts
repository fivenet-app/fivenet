export function parseJWTPayload<T>(token: string): T {
    const parts = token.split('.');
    if (parts.length !== 3) {
        throw new Error('Invalid JWT token');
    }

    try {
        const parsed = JSON.parse(atob(parts[1]!)) as T;
        return parsed;
    } catch {
        throw new Error('Invalid JWT payload encoding');
    }
}

// Base JWT claims info (e.g., expiration), limited to what we need right now.
export type JWTBaseClaimsType = {
    exp?: number;
};

export type JWTUserInfoClaimsType = JWTBaseClaimsType & {
    aid: number;
    uid: number;
    jb?: string;
    jbg?: number;
    su?: boolean;
    imp?: {
        jb: string;
        jbg: number;
    };
};

// JWTUserInfoClaims class represents the "cleaned up" user information extracted from
// our JWT token structure.
export class JWTUserInfoClaims {
    readonly expiration?: Date;

    readonly accountId: number;
    readonly userId: number;

    readonly job?: string;
    readonly jobGrade?: number;

    readonly superAdmin?: boolean;

    readonly impersonation?: {
        readonly job: string;
        readonly jobGrade: number;
    };

    constructor(token: string) {
        const claims = parseJWTPayload<JWTUserInfoClaimsType>(token);

        this.expiration = claims.exp ? new Date(claims.exp * 1000) : undefined;

        this.accountId = claims.aid;
        this.userId = claims.uid;

        this.job = claims.jb;
        this.jobGrade = claims.jbg;

        this.superAdmin = claims.su;

        if (claims.imp) {
            this.impersonation = {
                job: claims.imp.jb,
                jobGrade: claims.imp.jbg,
            };
        }
    }
}
