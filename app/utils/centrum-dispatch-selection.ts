import { isStatusDispatchCompleted } from '~/components/dispatch/helpers';
import { StatusDispatch, type Dispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';

/**
 * Keeps the selected own dispatch when it is still active, otherwise returns
 * the newest remaining active dispatch.
 */
export function selectOwnDispatch(
    ownDispatchIds: readonly number[],
    dispatches: ReadonlyMap<number, Dispatch>,
    selectedDispatchId?: number,
): number | undefined {
    if (selectedDispatchId !== undefined && ownDispatchIds.includes(selectedDispatchId)) {
        const selectedDispatch = dispatches.get(selectedDispatchId);
        if (
            selectedDispatch !== undefined &&
            !isStatusDispatchCompleted(selectedDispatch.status?.status ?? StatusDispatch.UNSPECIFIED)
        ) {
            return selectedDispatchId;
        }
    }

    return ownDispatchIds.find((dispatchId) => {
        const dispatch = dispatches.get(dispatchId);
        return dispatch !== undefined && !isStatusDispatchCompleted(dispatch.status?.status ?? StatusDispatch.UNSPECIFIED);
    });
}
