import { describe, expect, it } from 'vitest';
import { StatusDispatch, type Dispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';
import { selectOwnDispatch } from './centrum-dispatch-selection';

const dispatch = (id: number, status = StatusDispatch.NEW): Dispatch => ({
    id,
    job: 'ambulance',
    message: '',
    x: 0,
    y: 0,
    anon: false,
    units: [],
    status: { dispatchId: id, id: 1, status },
});

const dispatchMap = (...dispatches: Dispatch[]): Map<number, Dispatch> => new Map(dispatches.map((item) => [item.id, item]));

describe('selectOwnDispatch', () => {
    it('keeps the selected active dispatch', () => {
        const dispatches = dispatchMap(dispatch(101), dispatch(102), dispatch(103));

        expect(selectOwnDispatch([103, 102, 101], dispatches, 102)).toBe(102);
    });

    it('selects the next active dispatch when newer dispatches are completed', () => {
        const dispatches = dispatchMap(
            dispatch(103, StatusDispatch.COMPLETED),
            dispatch(102, StatusDispatch.COMPLETED),
            dispatch(101),
        );

        expect(selectOwnDispatch([103, 102, 101], dispatches, 103)).toBe(101);
    });

    it('selects the next dispatch when the selected dispatch was removed', () => {
        const dispatches = dispatchMap(dispatch(101), dispatch(102));

        expect(selectOwnDispatch([102, 101], dispatches, 103)).toBe(102);
    });

    it('skips missing and terminal dispatch projections', () => {
        const dispatches = dispatchMap(dispatch(101, StatusDispatch.CANCELLED), dispatch(100));

        expect(selectOwnDispatch([102, 101, 100], dispatches, 102)).toBe(100);
    });

    it('clears the selection when no active dispatch remains', () => {
        const dispatches = dispatchMap(dispatch(103, StatusDispatch.COMPLETED), dispatch(102, StatusDispatch.CANCELLED));

        expect(selectOwnDispatch([103, 102], dispatches, 103)).toBeUndefined();
    });

    it('clears the selection when there are no own dispatches', () => {
        expect(selectOwnDispatch([], new Map(), 103)).toBeUndefined();
    });
});
