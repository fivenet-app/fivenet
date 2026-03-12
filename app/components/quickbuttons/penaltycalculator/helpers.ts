import type { PenaltyCalculatorData } from '~~/gen/ts/resources/documents/data/data';
import type { Law, LawBook } from '~~/gen/ts/resources/laws/laws';

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

export function toPenaltyCalculatorData(selectedPenalties: SelectedPenalty[], reduction: number): PenaltyCalculatorData {
    return {
        reduction: reduction,
        selected: selectedPenalties.map((item) => ({
            lawId: item.law.id,
            count: item.count,
        })),
        total: {
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
        },
    };
}

export function resolvePenaltyCalculatorSelection(
    data: PenaltyCalculatorData | undefined,
    lawBooks: LawBook[] | undefined,
): SelectedPenalty[] {
    if (!data?.selected || !lawBooks || lawBooks.length === 0) return [];

    const laws = lawBooks.flatMap((book) => book.laws);

    return data.selected
        .map((item) => {
            const law = laws.find((l) => l.id === item.lawId);
            if (!law || item.count <= 0) return undefined;

            return {
                law,
                count: item.count,
            };
        })
        .filter((entry): entry is SelectedPenalty => entry !== undefined);
}
