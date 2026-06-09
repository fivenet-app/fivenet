import { describe, expect, it } from 'vitest';
import { isCalendarCreator } from './helpers';

describe('isCalendarCreator', () => {
    it('returns true when the active character owns the calendar', () => {
        expect(
            isCalendarCreator(
                {
                    userId: 1,
                    firstname: 'John',
                    lastname: 'Doe',
                    dateofbirth: '',
                    job: 'ambulance',
                    jobGrade: 2,
                },
                {
                    userId: 1,
                    firstname: 'John',
                    lastname: 'Doe',
                    dateofbirth: '',
                    job: 'police',
                    jobGrade: 4,
                },
                undefined,
            ),
        ).toBe(true);
    });

    it('returns false when the active character owns a job calendar', () => {
        expect(
            isCalendarCreator(
                {
                    userId: 1,
                    firstname: 'John',
                    lastname: 'Doe',
                    dateofbirth: '',
                    job: 'ambulance',
                    jobGrade: 2,
                },
                {
                    userId: 1,
                    firstname: 'John',
                    lastname: 'Doe',
                    dateofbirth: '',
                    job: 'police',
                    jobGrade: 4,
                },
                'ambulance',
            ),
        ).toBe(false);
    });

    it('returns false when the active character does not own the calendar', () => {
        expect(
            isCalendarCreator(
                {
                    userId: 1,
                    firstname: 'John',
                    lastname: 'Doe',
                    dateofbirth: '',
                    job: 'ambulance',
                    jobGrade: 2,
                },
                {
                    userId: 2,
                    firstname: 'Jane',
                    lastname: 'Smith',
                    dateofbirth: '',
                    job: 'police',
                    jobGrade: 4,
                },
                undefined,
            ),
        ).toBe(false);
    });

    it('returns false when no creator is present', () => {
        expect(
            isCalendarCreator({
                userId: 1,
                firstname: 'John',
                lastname: 'Doe',
                dateofbirth: '',
                job: 'ambulance',
                jobGrade: 2,
            }),
        ).toBe(false);
    });
});
