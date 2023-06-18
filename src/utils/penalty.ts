import { Law, LawBook } from '~~/gen/ts/resources/laws/laws';

export type Penalty = Law & {
    show?: boolean;
};

export type PenaltyCategory = LawBook & {
    laws: Penalty[];
    show?: boolean;
};

export type Penalties = Array<PenaltyCategory>;

export type SelectedPenalty = {
    penalty: Penalty;
    count: bigint;
};

export type PenaltiesSummary = {
    fine: bigint;
    detentionTime: bigint;
    stvoPoints: bigint;
    count: bigint;
};
