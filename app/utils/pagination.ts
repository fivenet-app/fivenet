import type { PaginationResponse } from '~~/gen/ts/resources/common/database/database';

// User-facing client query page sizes. Keep these centralized until they become configurable.
// Backend page sizes are in gen/go/proto/resources/common/database/database.go
export const PAGINATION_PAGE_SIZE = {
    notificationPopover: 6,
    qualificationResults: 10,
    qualificationRequests: 10,
    commandSearch: 10,
    helpGuideEligibility: 1,
    mailerThreadAttachments: 10,
    mailerMessageSearch: 6,
    wikiPageSearch: 6,
    documentList: 20,
    documentReferenceManager: 8,
} as const;

export function calculateOffset(page: number, pagination?: PaginationResponse): number {
    return pagination?.pageSize ? pagination?.pageSize * (page - 1) : 0;
}
