import type { AttributeValues } from '~~/gen/ts/resources/permissions/permissions';

export function isEmptyAttributes(val?: AttributeValues): boolean {
    if (!val) {
        return false;
    }

    switch (val.validValues.oneofKind) {
        case 'stringList':
            return val.validValues.stringList.strings.length === 0;
        case 'jobGradeList':
            return val.validValues.jobGradeList.jobs.length === 0;
        case 'jobList':
            return val.validValues.jobList.strings.length === 0;
    }

    return false;
}
