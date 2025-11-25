import { describe, expect, it } from 'vitest';
import { QualificationExamMode, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { checkAccess } from './acl';

describe('checkAccess', () => {
    const activeChar = {
        userId: 1,
        job: 'police',
        jobGrade: 3,
        firstname: 'John',
        lastname: 'Doe',
    };
    const creator = {
        userId: 2,
        job: 'police',
        jobGrade: 2,
        firstname: 'Jane',
        lastname: 'Smith',
    };

    function getQualificationAccess(result = ResultStatus.SUCCESSFUL) {
        return {
            qualifications: [
                {
                    qualificationId: 1,
                    qualification: {
                        id: 1,
                        job: 'police',
                        weight: 10,
                        closed: false,
                        draft: false,
                        public: true,
                        abbreviation: 'POL',
                        title: 'Police Qualification',
                        result: {
                            id: 1,
                            qualificationId: 1,
                            userId: 1,
                            summary: '',
                            status: result,
                            creatorJob: 'police',
                            creatorId: 1,
                        },
                        creatorJob: 'police',
                        requirements: [],
                        examMode: QualificationExamMode.DISABLED,
                    },
                    access: 1,
                },
            ],
        };
    }

    it('should return false if access is undefined', () => {
        const result = checkAccess(activeChar, undefined, creator, 1);
        expect(result).toBe(false);
    });

    it('should return true if activeChar is the creator with the same job', () => {
        const result = checkAccess(activeChar, {}, activeChar, 1, 'police');
        expect(result).toBe(true);
    });

    it('should return false if activeChar is the creator but with a different job', () => {
        const result = checkAccess(activeChar, {}, activeChar, 1, 'firefighter');
        expect(result).toBe(false);
    });

    it('should return true if user access matches', () => {
        const access = { users: [{ userId: 1, access: 1 }] };
        const result = checkAccess(activeChar, access, creator, 1);
        expect(result).toBe(true);
    });

    it('should return true if job access matches', () => {
        const access = { jobs: [{ job: 'police', minimumGrade: 2, access: 1 }] };
        const result = checkAccess(activeChar, access, creator, 1);
        expect(result).toBe(true);
    });

    it('should return true because no access matches', () => {
        const access = {
            users: [{ userId: 2, access: 1 }],
            jobs: [{ job: 'police', minimumGrade: 5, access: 1 }],
        };
        const result = checkAccess(activeChar, access, creator, 1);
        expect(result).toBe(false);
    });

    it('should return true if qualification access matches', () => {
        const access = getQualificationAccess();
        const result = checkAccess(activeChar, access, creator, 1);
        expect(result).toBe(true);
    });

    it("should return false if qualification access isn't successful ", () => {
        const access = getQualificationAccess(ResultStatus.FAILED);
        const result = checkAccess(activeChar, access, creator, 1);
        expect(result).toBe(false);
    });

    it('should return false if no access matches', () => {
        const access = {
            users: [{ userId: 3, access: 1 }],
            jobs: [{ job: 'firefighter', minimumGrade: 2, access: 1 }],
        };
        const result = checkAccess(activeChar, access, creator, 1);
        expect(result).toBe(false);
    });
});
