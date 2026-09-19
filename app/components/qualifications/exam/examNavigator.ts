import type { ExamQuestion } from '~~/gen/ts/resources/qualifications/exam/exam';

export interface ExamNavigatorEntry {
    question: ExamQuestion;
    kind: string;
    answerable: boolean;
    number?: number;
}

export function createExamNavigatorEntries(questions: ExamQuestion[]): ExamNavigatorEntry[] {
    let questionNumber = 0;

    return questions.map((question) => {
        const kind = question.data?.data.oneofKind ?? 'separator';
        const answerable = kind !== 'separator' && kind !== 'image';
        if (answerable) questionNumber += 1;

        return { question, kind, answerable, number: answerable ? questionNumber : undefined };
    });
}

export function getExamNavigatorTargetIndex(key: string, currentIndex: number, length: number): number {
    if (key === 'Home') return 0;
    if (key === 'End') return length - 1;
    if (key === 'ArrowDown') return currentIndex + 1;
    if (key === 'ArrowUp') return currentIndex - 1;
    return -1;
}
