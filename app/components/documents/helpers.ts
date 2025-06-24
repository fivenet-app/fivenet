import type { BadgeColor, ButtonColor } from '#ui/types';
import type { Perms } from '~~/gen/ts/perms';
import type { AccessLevel, DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import { DocReference, DocRelation } from '~~/gen/ts/resources/documents/documents';
import type { UserShort } from '~~/gen/ts/resources/users/users';

export const logger = useLogger('📃 Docstore');

export function checkDocAccess(
    docAccess: DocumentAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
    perm?: Perms,
): boolean {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    if (!checkBaseDocAccess(activeChar.value, docAccess, creator, level)) {
        return false;
    }

    if (perm !== undefined && creator !== undefined && creator?.job === activeChar.value.job) {
        return checkIfCanAccessOwnJobDocument(activeChar.value, creator, perm);
    }

    return true;
}

function checkBaseDocAccess(
    activeChar: UserShort,
    access: DocumentAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
): boolean {
    return checkAccess(activeChar, access, creator, level);
}

function checkIfCanAccessOwnJobDocument(activeChar: UserShort, creator: UserShort, perm: Perms): boolean {
    const { attrStringList } = useAuth();

    const fields = attrStringList(perm, 'Access').value;
    if (fields.length === 0) {
        return creator?.userId === activeChar.userId;
    }

    if (fields.includes('Any')) {
        return true;
    }
    if (fields.includes('Lower_Rank')) {
        if (creator?.jobGrade < activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('Same_Rank')) {
        if (creator?.jobGrade <= activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('Own')) {
        if (creator?.userId === activeChar.userId) {
            return true;
        }
    }

    return false;
}

// Document Refernces

export function docReferenceToIcon(ref: DocReference): string {
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

export function docReferenceToColor(ref: DocReference): ButtonColor {
    switch (ref) {
        case DocReference.LINKED:
            return 'blue';
        case DocReference.SOLVES:
            return 'success';
        case DocReference.CLOSES:
            return 'error';
        case DocReference.DEPRECATES:
            return 'amber';
        default:
            return 'white';
    }
}

export function docReferenceToBadge(ref: DocReference): BadgeColor {
    switch (ref) {
        case DocReference.LINKED:
            return 'blue';
        case DocReference.SOLVES:
            return 'success';
        case DocReference.CLOSES:
            return 'error';
        case DocReference.DEPRECATES:
            return 'amber';
        default:
            return 'white';
    }
}

// Document Relations

export function docRelationToIcon(ref: DocRelation): string {
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

export function docRelationToColor(ref: DocRelation): ButtonColor {
    switch (ref) {
        case DocRelation.MENTIONED:
            return 'blue';
        case DocRelation.TARGETS:
            return 'amber';
        case DocRelation.CAUSED:
            return 'error';
        default:
            return 'white';
    }
}

export function docRelationToBadge(ref: DocRelation): BadgeColor {
    switch (ref) {
        case DocRelation.MENTIONED:
            return 'blue';
        case DocRelation.TARGETS:
            return 'amber';
        case DocRelation.CAUSED:
            return 'error';
        default:
            return 'white';
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
        case DocActivityType.DRAFT_TOGGLED:
            return 'i-mdi-file-document-edit-outline';

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
            return 'i-mdi-delete';

        default:
            return 'i-mdi-help';
    }
}
