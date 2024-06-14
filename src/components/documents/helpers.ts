import type { BadgeColor, ButtonColor } from '#ui/types';
import { useAuthStore } from '~/store/auth';
import type { Perms } from '~~/gen/ts/perms';
import { AccessLevel, DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import { DocReference, DocRelation } from '~~/gen/ts/resources/documents/documents';
import { User, UserShort } from '~~/gen/ts/resources/users/users';

export const logger = useLogger('📃 Docstore');

export function checkDocAccess(
    docAccess: DocumentAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
    perm?: Perms,
): boolean {
    const authStore = useAuthStore();
    if (authStore.isSuperuser) {
        return true;
    }

    const activeChar = authStore.activeChar;
    if (activeChar === null) {
        return false;
    }

    if (!checkBaseDocAccess(activeChar, docAccess, creator, level)) {
        return false;
    }

    if (perm !== undefined && creator !== undefined && creator?.job === activeChar.job) {
        return checkIfCanAccessOwnJobDocument(activeChar, creator, perm);
    }

    return true;
}

function checkBaseDocAccess(
    activeChar: User,
    docAccess: DocumentAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
): boolean {
    if (docAccess === undefined) {
        return false;
    }

    if (creator !== undefined && activeChar.userId === creator.userId) {
        return true;
    }

    const ju = docAccess.users.find((ua) => ua.userId === activeChar.userId && level <= ua.access);
    if (ju !== undefined) {
        return true;
    }

    const ja = docAccess.jobs.find(
        (ja) => ja.job === activeChar.job && ja.minimumGrade <= activeChar.jobGrade && level <= ja.access,
    );
    if (ja !== undefined) {
        return true;
    }

    return false;
}

function checkIfCanAccessOwnJobDocument(activeChar: User, creator: UserShort, perm: Perms): boolean {
    const authStore = useAuthStore();
    if (authStore.isSuperuser) {
        return true;
    }

    const fields = attrList(perm, 'Access').value;
    if (fields.length === 0) {
        return creator?.userId === activeChar.userId;
    }

    if (fields.includes('any')) {
        return true;
    }
    if (fields.includes('lower_rank')) {
        if (creator?.jobGrade < activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('same_rank')) {
        if (creator?.jobGrade <= activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('own')) {
        if (creator?.userId === activeChar.userId) {
            return true;
        }
    }

    return false;
}

// Document Refernces

export function refToIcon(ref: DocReference): string {
    switch (ref) {
        case DocReference.LINKED:
            return 'i-mdi-link';
        case DocReference.SOLVES:
            return 'i-mdi-check';
        case DocReference.CLOSES:
            return 'i-mdi-close-box';
        case DocReference.DEPRECATES:
            return 'i-mdi-lock-clock';
        default:
            return 'i-mdi-question-mark';
    }
}

export function refToColor(ref: DocReference): ButtonColor {
    switch (ref) {
        case DocReference.LINKED:
            return 'blue';
        case DocReference.SOLVES:
            return 'green';
        case DocReference.CLOSES:
            return 'red';
        case DocReference.DEPRECATES:
            return 'amber';
        default:
            return 'black';
    }
}

export function refToBadge(ref: DocReference): BadgeColor {
    switch (ref) {
        case DocReference.LINKED:
            return 'blue';
        case DocReference.SOLVES:
            return 'green';
        case DocReference.CLOSES:
            return 'red';
        case DocReference.DEPRECATES:
            return 'amber';
        default:
            return 'black';
    }
}

// Document Relations

export function relationToIcon(ref: DocRelation): string {
    switch (ref) {
        case DocRelation.MENTIONED:
            return 'i-mdi-at';
        case DocRelation.TARGETS:
            return 'i-mdi-target';
        case DocRelation.CAUSED:
            return 'i-mdi-source-commit-start';
        default:
            return 'i-mdi-question-mark';
    }
}

export function relationToColor(ref: DocRelation): ButtonColor {
    switch (ref) {
        case DocRelation.MENTIONED:
            return 'blue';
        case DocRelation.TARGETS:
            return 'amber';
        case DocRelation.CAUSED:
            return 'red';
        default:
            return 'black';
    }
}

export function relationToBadge(ref: DocRelation): BadgeColor {
    switch (ref) {
        case DocRelation.MENTIONED:
            return 'blue';
        case DocRelation.TARGETS:
            return 'amber';
        case DocRelation.CAUSED:
            return 'red';
        default:
            return 'black';
    }
}

// Document Activity

export function getDocAtivityIcon(activityType: DocActivityType): string {
    switch (activityType) {
        // Base
        case DocActivityType.CREATED:
            return 'i-mdi-new-box';
        case DocActivityType.STATUS_OPEN:
            return 'i-mdi-lock-open-variant';
        case DocActivityType.STATUS_CLOSED:
            return 'i-mdi-lock';
        case DocActivityType.UPDATED:
            return 'i-mdi-update';
        case DocActivityType.RELATIONS_UPDATED:
            return 'i-mdi-account-multiple';
        case DocActivityType.REFERENCES_UPDATED:
            return 'i-mdi-file-multiple';
        case DocActivityType.ACCESS_UPDATED:
            return 'i-mdi-lock-check';
        case DocActivityType.OWNER_CHANGED:
            return 'i-mdi-file-account';
        case DocActivityType.DELETED:
            return 'i-mdi-delete-circle';

        // Requests
        case DocActivityType.REQUESTED_ACCESS:
            return 'i-mdi-lock-plus-outline';
        case DocActivityType.REQUESTED_CLOSURE:
            return 'i-mdi-lock-question';
        case DocActivityType.REQUESTED_OPENING:
            return 'i-mdi-lock-open-outline';
        case DocActivityType.REQUESTED_UPDATE:
            return 'i-mdi-refresh-circle';
        case DocActivityType.REQUESTED_OWNER_CHANGE:
            return 'i-mdi-file-swap-outline';
        case DocActivityType.REQUESTED_DELETION:
            return 'i-mdi-delete-circle-outline';

        // Comments
        case DocActivityType.COMMENT_ADDED:
            return 'i-mdi-comment-plus';
        case DocActivityType.COMMENT_UPDATED:
            return 'i-mdi-comment-edit';
        case DocActivityType.COMMENT_DELETED:
            return 'i-mdi-trash-can';

        default:
            return 'i-mdi-help';
    }
}
