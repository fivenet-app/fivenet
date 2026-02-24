import type { Law } from '~~/gen/ts/resources/laws/laws';

export type PenaltiesSummary = {
    fine: number;
    detentionTime: number;
    stvoPoints: number;
    count: number;
};

export type SelectedPenalty = {
    law: Law;
    count: number;
};

export function calculatePenaltySummary(selectedPenalties: SelectedPenalty[]): PenaltiesSummary {
    return {
        fine: selectedPenalties.reduce((acc, curr) => acc + (curr.law.fine ? curr.law.fine * curr.count : 0), 0),
        detentionTime: selectedPenalties.reduce(
            (acc, curr) => acc + (curr.law.detentionTime ? curr.law.detentionTime * curr.count : 0),
            0,
        ),
        stvoPoints: selectedPenalties.reduce(
            (acc, curr) => acc + (curr.law.stvoPoints ? curr.law.stvoPoints * curr.count : 0),
            0,
        ),
        count: selectedPenalties.reduce((acc, curr) => acc + curr.count, 0),
    };
}
