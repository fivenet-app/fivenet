export type Penalty = {
    category?: string;
    name: string;
    description: string;
    fine?: number;
    detentionTime?: number;
    stvoPoints?: number;
};

export type PenaltyCategory = { name: string; penalties: Array<Penalty> };

export type Penalties = Array<PenaltyCategory>;

export type SelectedPenalty = {
    penalty: Penalty;
    count: number;
};

export type PenaltiesSummary = {
    fine: number;
    detentionTime: number;
    stvoPoints: number;
    count: number;
};
