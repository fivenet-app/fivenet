import { describe, expect, it } from 'vitest';
import { isRoute } from './route';

describe('isRoute', () => {
    it('matches the route itself and child routes', () => {
        expect(isRoute('/jobs', '/jobs')).toBe(true);
        expect(isRoute('/jobs/overview', '/jobs')).toBe(true);
    });

    it('does not match unrelated path prefixes', () => {
        expect(isRoute('/jobs-old', '/jobs')).toBe(false);
        expect(isRoute('/job', '/jobs')).toBe(false);
    });
});
