import type { BadgeProps, ButtonProps } from '@nuxt/ui';
import type { CentrumAccessLevel } from '~~/gen/ts/resources/centrum/access';
import { type JobList, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { type Unit, StatusUnit } from '~~/gen/ts/resources/centrum/units';
import type { UnitAccess, UnitAccessLevel } from '~~/gen/ts/resources/centrum/units_access';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

export type GroupedUnits = { status: StatusUnit; key: string; units: Unit[] }[];

export function dispatchStatusToFillColor(status: StatusDispatch | undefined): string {
    switch (status) {
        case StatusDispatch.UNSPECIFIED:
        case StatusDispatch.NEW:
        case StatusDispatch.UNASSIGNED:
        case StatusDispatch.UNIT_DECLINED:
            return '!text-error-600';
        case StatusDispatch.EN_ROUTE:
            return '!text-info-500';
        case StatusDispatch.ON_SCENE:
            return '!text-info-700';
        case StatusDispatch.NEED_ASSISTANCE:
            return '!text-warning-600';
        case StatusDispatch.COMPLETED:
            return '!text-success-600';
        case StatusDispatch.CANCELLED:
            return '!text-success-800';
        case StatusDispatch.ARCHIVED:
            return '!text-base-600';
        case StatusDispatch.UNIT_ACCEPTED:
            return '!text-info-600';
        default:
            return '!text-info-500';
    }
}

export function dispatchStatusToBGColor(status: StatusDispatch | undefined): string {
    switch (status) {
        case StatusDispatch.UNSPECIFIED:
        case StatusDispatch.NEW:
        case StatusDispatch.UNASSIGNED:
        case StatusDispatch.UNIT_DECLINED:
            return '!bg-error-600 text-white!';
        case StatusDispatch.EN_ROUTE:
            return '!bg-info-600 text-white!';
        case StatusDispatch.ON_SCENE:
            return '!bg-info-700 text-white!';
        case StatusDispatch.NEED_ASSISTANCE:
            return '!bg-warning-600 text-white!';
        case StatusDispatch.COMPLETED:
            return '!bg-success-600 text-white!';
        case StatusDispatch.CANCELLED:
            return '!bg-success-800 text-white!';
        case StatusDispatch.ARCHIVED:
            return '!bg-default text-white!';
        case StatusDispatch.UNIT_ACCEPTED:
            return '!bg-info-600 text-white!';
        default:
            return '!bg-info-600 text-white!';
    }
}

export function dispatchStatusToBadgeColor(status: StatusDispatch | undefined): BadgeProps['color'] {
    switch (status) {
        case StatusDispatch.UNSPECIFIED:
        case StatusDispatch.NEW:
        case StatusDispatch.UNASSIGNED:
        case StatusDispatch.UNIT_DECLINED:
            return 'error';
        case StatusDispatch.EN_ROUTE:
            return 'info';
        case StatusDispatch.ON_SCENE:
            return 'info';
        case StatusDispatch.NEED_ASSISTANCE:
            return 'warning';
        case StatusDispatch.COMPLETED:
            return 'success';
        case StatusDispatch.CANCELLED:
            return 'success';
        case StatusDispatch.ARCHIVED:
            return 'white';
        case StatusDispatch.UNIT_ACCEPTED:
            return 'info';
        default:
            return 'info';
    }
}

export function dispatchStatusToButtonColor(status: StatusDispatch | undefined): ButtonProps['color'] {
    switch (status) {
        case StatusDispatch.UNSPECIFIED:
        case StatusDispatch.NEW:
        case StatusDispatch.UNASSIGNED:
        case StatusDispatch.UNIT_DECLINED:
            return 'error';
        case StatusDispatch.EN_ROUTE:
            return 'info';
        case StatusDispatch.ON_SCENE:
            return 'info';
        case StatusDispatch.NEED_ASSISTANCE:
            return 'warning';
        case StatusDispatch.COMPLETED:
            return 'success';
        case StatusDispatch.CANCELLED:
            return 'success';
        case StatusDispatch.ARCHIVED:
            return 'white';
        case StatusDispatch.UNIT_ACCEPTED:
            return 'info';
        default:
            return 'info';
    }
}

export function dispatchStatusToIcon(status: StatusDispatch | undefined): string {
    const found = dispatchStatuses.find((ds) => ds.status === status);
    if (found) {
        return found.icon;
    }

    return 'i-mdi-info-circle';
}

export const animateStates = [
    StatusDispatch.NEW.valueOf(),
    StatusDispatch.UNASSIGNED.valueOf(),
    StatusDispatch.NEED_ASSISTANCE.valueOf(),
];

export function dispatchStatusAnimate(status: StatusDispatch | undefined): boolean {
    return animateStates.includes((status ?? StatusDispatch.NEW).valueOf());
}

export function unitStatusToBGColor(status: StatusUnit | undefined): string {
    switch (status) {
        case StatusUnit.ON_BREAK:
        case StatusUnit.USER_ADDED:
        case StatusUnit.USER_REMOVED:
            return '!bg-info-600';
        case StatusUnit.AVAILABLE:
            return '!bg-success-600';
        case StatusUnit.BUSY:
            return '!bg-warning-600';
        case StatusUnit.UNSPECIFIED:
        case StatusUnit.UNKNOWN:
        case StatusUnit.UNAVAILABLE:
        default:
            return '!bg-error-600';
    }
}

export function unitStatusToBadgeColor(status: StatusUnit | undefined): BadgeProps['color'] {
    switch (status) {
        case StatusUnit.ON_BREAK:
        case StatusUnit.USER_ADDED:
        case StatusUnit.USER_REMOVED:
            return 'info';
        case StatusUnit.AVAILABLE:
            return 'success';
        case StatusUnit.BUSY:
            return 'warning';
        case StatusUnit.UNSPECIFIED:
        case StatusUnit.UNKNOWN:
        case StatusUnit.UNAVAILABLE:
        default:
            return 'error';
    }
}

export function unitStatusToIcon(status: StatusUnit | undefined): string {
    const found = unitStatuses.find((ds) => ds.status === status);
    if (found) {
        return found.icon;
    }

    return 'i-mdi-info-circle';
}

export const statusOrder = [
    StatusUnit.AVAILABLE,
    StatusUnit.ON_BREAK,
    StatusUnit.BUSY,
    StatusUnit.USER_ADDED,
    StatusUnit.USER_REMOVED,
    StatusUnit.UNAVAILABLE,
    StatusUnit.UNKNOWN,
    StatusUnit.UNSPECIFIED,
];

export const unitStatuses: {
    icon: string;
    name: string;
    status?: StatusUnit;
}[] = [
    { icon: 'i-mdi-cancel', name: 'Unavailable', status: StatusUnit.UNAVAILABLE },
    { icon: 'i-mdi-calendar-check', name: 'Available', status: StatusUnit.AVAILABLE },
    { icon: 'i-mdi-coffee', name: 'On Break', status: StatusUnit.ON_BREAK },
    { icon: 'i-mdi-calendar-remove', name: 'Busy', status: StatusUnit.BUSY },
];

export const dispatchStatuses: {
    icon: string;
    name: string;
    status?: StatusDispatch;
}[] = [
    { icon: 'i-mdi-car-back', name: 'En Route', status: StatusDispatch.EN_ROUTE },
    { icon: 'i-mdi-map-marker-check', name: 'On Scene', status: StatusDispatch.ON_SCENE },
    { icon: 'i-mdi-help-circle', name: 'Need Assistance', status: StatusDispatch.NEED_ASSISTANCE },
    { icon: 'i-mdi-check-bold', name: 'Completed', status: StatusDispatch.COMPLETED },
    { icon: 'i-mdi-cancel', name: 'Cancelled', status: StatusDispatch.CANCELLED },
];

export function isStatusDispatchCompleted(status: StatusDispatch): boolean {
    return status === StatusDispatch.ARCHIVED || status === StatusDispatch.CANCELLED || status === StatusDispatch.COMPLETED;
}

// "Color stops" with optional ping
const steps: {
    variant: BadgeProps['variant'];
    color?: BadgeProps['color'];
    class: string;
    animation?: string;
    ping: boolean;
}[] = [
    { variant: 'solid', class: '', ping: false }, // 0-10%
    { variant: 'solid', color: 'green', class: 'bg-green-200!', ping: false }, // 10-20%
    { variant: 'solid', color: 'warning', class: 'bg-yellow-200!', ping: false }, // 20-30%
    { variant: 'solid', color: 'warning', class: 'bg-yellow-400!', ping: false }, // 30-40%
    { variant: 'solid', color: 'orange', class: 'bg-orange-300!', ping: false }, // 40-50%
    { variant: 'solid', color: 'orange', class: 'bg-orange-500!', ping: false }, // 50-60%
    { variant: 'solid', color: 'red', class: 'bg-red-400!', ping: true }, // 60-70%
    { variant: 'solid', color: 'red', class: 'bg-red-600!', ping: true }, // 70-80%
    { variant: 'solid', color: 'red', class: 'bg-red-700!', ping: true }, // 80-90%
    { variant: 'solid', color: 'red', class: 'bg-red-800!', ping: true }, // 90-100%
    { variant: 'solid', color: 'red', class: 'bg-red-700!', animation: 'animate-bounce', ping: true }, // caps at 100%
];

export function dispatchTimeToBadge(
    date: Timestamp | undefined,
    status: StatusDispatch = StatusDispatch.UNSPECIFIED,
    maxTime: number = 600,
): {
    color?: BadgeProps['color'];
    variant?: BadgeProps['variant'];
    class?: string;
} {
    if (isStatusDispatchCompleted(status)) {
        return { variant: 'solid', color: 'neutral' };
    }

    // Elapsed time in seconds since dispatch
    const elapsed = (Date.now() - toDate(date).getTime()) / 1000;
    // Fraction of max elapsed (clamped 0…1)
    const over = Math.min(Math.max(elapsed / maxTime, 0), 1);

    // Pick one of the above based on over ∈ [0,1]
    const idx = Math.floor(over * (steps.length - 1));
    const step = steps[idx] ?? steps[steps.length - 1]!;
    if (step.animation) {
        return { variant: 'soft', class: `${step.class} ${step.animation}`, color: step.color };
    }
    return { variant: 'solid', class: step.class, color: step.color };
}

export function dispatchTimeToTextColorSidebar(
    date: Timestamp | undefined,
    status: StatusDispatch = StatusDispatch.UNSPECIFIED,
    maxTime: number = 900,
): { ping: boolean; class: string } {
    if (isStatusDispatchCompleted(status)) {
        return { ping: false, class: '' };
    }

    // elapsed time in seconds since dispatch
    const elapsed = (Date.now() - toDate(date).getTime()) / 1000;
    // fraction of max elapsed (clamped 0…1)
    const over = Math.min(Math.max(elapsed / maxTime, 0), 1);

    // Pick the right stop
    const idx = Math.floor(over * (steps.length - 1));
    return steps[idx] ?? steps[steps.length - 1]!;
}

export function checkUnitAccess(unitAccess: UnitAccess | undefined, level: UnitAccessLevel): boolean {
    if (unitAccess === undefined || (unitAccess.jobs.length === 0 && unitAccess.qualifications.length === 0)) {
        return true;
    }

    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    if (!checkAccess(activeChar.value, unitAccess, undefined, level)) {
        return false;
    }

    return true;
}

export function checkDispatchAccess(dispatchJobs: JobList | undefined, level: CentrumAccessLevel): boolean {
    const centrumStore = useCentrumStore();
    const { acls } = storeToRefs(centrumStore);

    // The dispatch has no ACLs at all
    if (dispatchJobs === undefined || dispatchJobs.jobs.length === 0) {
        return true;
    }

    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    // If the active character's job isn't in the dispatch's jobs list, check dispatchesJobs access list (contains the job to job access)
    if (!dispatchJobs?.jobs.find((j) => j.name === activeChar.value?.job)) {
        if (acls.value?.dispatches?.jobs) {
            for (const djob of dispatchJobs.jobs) {
                let lowestAccess: CentrumAccessLevel | undefined = undefined;
                for (let index = 0; index < acls.value.dispatches.jobs.length; index++) {
                    const ja = acls.value.dispatches.jobs[index]!;
                    if (ja.job !== djob.name) {
                        continue;
                    }
                    // Dispatch access doesn't use job grade for access check.
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
        }

        return false;
    }

    return true;
}

export function calculateDispatchZIndexOffset(status: StatusDispatch | undefined): number {
    switch (status) {
        case StatusDispatch.COMPLETED:
        case StatusDispatch.CANCELLED:
        case StatusDispatch.ARCHIVED:
            return 5;

        case StatusDispatch.NEW:
        case StatusDispatch.UNASSIGNED:
        case StatusDispatch.UNIT_DECLINED:
            return 15;

        default:
            return 10;
    }
}
