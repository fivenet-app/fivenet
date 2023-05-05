<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { MinusSmallIcon, PlusSmallIcon } from '@heroicons/vue/24/outline';
import ListEntry from '~/components/penaltycalculator/ListEntry.vue';
import { Penalties, PenaltiesSummary, SelectedPenalty } from '~/utils/penalty';
import Stats from '~/components/penaltycalculator/Stats.vue';
import SummaryTable from './SummaryTable.vue';

const penalties: Penalties = [
    {
        name: 'STGB',
        penalties: [
            {
                name: '§12 Mord',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 45,
                stvoPoints: 0,
            },
            {
                name: '§13 Totschlag',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0,
            },
        ],
    },
    {
        name: 'WaffG',
        penalties: [
            {
                name: '§7 Abs. 2 Nr. 1 Rechtswidriger Waffenbesitz',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0,
            },
            {
                name: '§7 Abs. 2 Nr. 2 Öffentliches Führen einer Schusswaffe',
                description: 'Freiheitstrafe nötig',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0,
            },
        ],
    },
    {
        name: 'BtMG',
        penalties: [
            {
                name: '§3 Straftat',
                description: 'Gesetz gültig für: Cannabis, Kokain, Methamphetamin, Lysergid, Opium | Geldstrafe möglich.',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            }
        ]
    },
    {
        name: 'LuftVO',
        penalties: [
            {
                name: '§1 Registrierung von Flugobjekten',
                description: 'Geldstrafe möglich',
                fine: 20000,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§2 Flugverbotszonen',
                description: 'Geldstrafe möglich',
                fine: 30000,
                detentionTime: 25,
                stvoPoints: 0
            },
            {
                name: '§3 Mindestflughöhe und Landung',
                description: 'Bei Gefährdung anderer, Haftzeit anwenden.',
                fine: 15000,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§4 Flugfunk',
                description: 'Geldstrafe möglich',
                fine: 15000,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§5 Mitführpflicht von Fallschirmen',
                description: 'Geldstrafe möglich',
                fine: 20000,
                detentionTime: 0,
                stvoPoints: 0
            }
        ]
    },
    {
        name: 'StVO',
        penalties: [
            {
                name: '§1 Allgemeine Vorsicht im Straßenverkehr',
                description: 'Geldstrafe',
                fine: 2200,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§2 Straßenbenutzung durch Fahrzeuge',
                description: 'Geldstrafe',
                fine: 2200,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§3 ab 10 km/h Geschwindigkeitsüberschreitung',
                description: '3 km/h Toleranz abziehen, Geldstrafe',
                fine: 2000,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§3 ab 30 km/h Geschwindigkeitsüberschreitung',
                description: '3 km/h Toleranz abziehen, Geldstrafe',
                fine: 4000,
                detentionTime: 0,
                stvoPoints: 1
            },
            {
                name: '§3 ab 50 km/h Geschwindigkeitsüberschreitung',
                description: '3 km/h Toleranz abziehen, Geldstrafe',
                fine: 7000,
                detentionTime: 0,
                stvoPoints: 2
            },
            {
                name: '§3 ab 80 km/h Geschwindigkeitsüberschreitung',
                description: '3 km/h Toleranz abziehen, Freiheitsstrafe möglich',
                fine: 12000,
                detentionTime: 5,
                stvoPoints: 3
            },
            {
                name: '§4 Abstand',
                description: 'Geldstrafe',
                fine: 500,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§5 Überholen',
                description: 'Geldstrafe',
                fine: 1000,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§8 Vorfahrt',
                description: 'Geldstrafe',
                fine: 1800,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§9 Abbiegen, Wenden und Rückwärtsfahren',
                description: 'Geldstrafe',
                fine: 1000,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§11 Besondere Verkehrslagen (Stau, Stockender Verkehr)',
                description: 'Geldstrafe',
                fine: 900,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§12 Halten und Parken',
                description: 'Geldstrafe, doppeltes Bußgeld und 1 StVO Punkt bei Behinderung von Einsatzkräften/Einsatzfahrzeugen von PD/FIB/LSMD/DoJ',
                fine: 1150,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§15 Liegenbleiben von Fahrzeugen',
                description: 'Geldstrafe',
                fine: 700,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§16 Abschleppen von Fahrzeugen mit Abschleppseilen',
                description: 'Geldstrafe',
                fine: 1000,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§17 Beleuchtung',
                description: 'Geldstrafe',
                fine: 1000,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§18 Autobahnen und Kraftfahrstraßen',
                description: 'Geldstrafe',
                fine: 1000,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§19 Bahnübergänge',
                description: 'Geldstrafe',
                fine: 1500,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§21 Personenbeförderung',
                description: 'Geldstrafe',
                fine: 750,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§21a Sicherheitsgurte, Schutzhelme',
                description: 'Geldstrafe',
                fine: 2500,
                detentionTime: 0,
                stvoPoints: 1
            },
            {
                name: '§22 Sonstige Pflichten von Fahrzeugführenden',
                description: 'Geldstrafe',
                fine: 1800,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§24 Übermäßige Straßenbenutzung',
                description: 'Geldstrafe',
                fine: 2250,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§25 Verkehrshindernisse',
                description: 'Geldstrafe',
                fine: 1350,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§26 Verkehrsbeeinträchtigungen',
                description: 'Geldstrafe',
                fine: 1000,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§27 Unfall',
                description: 'Geldstrafe',
                fine: 2000,
                detentionTime: 0,
                stvoPoints: 1
            }
        ]
    },
    {
        name: 'GewO',
        penalties: [
            {
                name: '§7b Illegaler gewerblicher Abbau, Transport, Handel und Verarbeitung von Rohstoffen',
                description: 'Freimenge ohne Lizenz: 200kg pro Rohstoff am Tag. Bei mehr Rohstoffen ohne Lizenz: Beschlagnahmung der zu viel mitgeführten Rohstoffe (Freimenge kann belassen werden). Bußgeld nötig: 200-300kg Minimalstrafe, bei Wiederholungstätern/bei 300kg und mehr bis zu 65.000$ Bußgeld. Fahrzeuge ab 200kg Kofferraum-Kapazität (LKW) dürfen zur Überprüfung bei Verkehrskontrollen durchsucht werden, wenn keine Rohstoff-Gewerbelizenz vorgelegt werden kann. Für Weintrauben/Muscheln keine Lizenz nötig!!',
                fine: 30000,
                detentionTime: 0,
                stvoPoints: 0
            }
        ]
    },
    {
        name: 'WirtG',
        penalties: [
            {
                name: '§12 Vermögensdiebstahl',
                description: 'Geldstrafe möglich',
                fine: 17500,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§13 Vermögensmissbrauch',
                description: 'Geldstrafe möglich',
                fine: 40000,
                detentionTime: 40,
                stvoPoints: 0
            },
            {
                name: '§14 Vermögensunterschlagung',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§15 Steuerhinterziehung',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 25,
                stvoPoints: 0
            },
            {
                name: '§16 Strafgeldhinterziehung',
                description: 'Geldstrafe nötig',
                fine: 50000,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§17 Besitz von Schwarzgeld',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§18 Besitz von Falschgeld',
                description: 'Freiheitsstrafe und Geldstrafe nötig',
                fine: 25000,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§19 Diebstahl von Schwarzgeld',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 25,
                stvoPoints: 0
            },
            {
                name: '§20 Herstellung von Falschgeld',
                description: 'Freiheitsstrafe und Geldstrafe nötig (2$ pro gefälschtem Dollar)',
                fine: 2,
                detentionTime: 40,
                stvoPoints: 0
            },
            {
                name: '§21 Handel mit illegalen Währungen',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§22 Geldwäsche',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 25,
                stvoPoints: 0
            },
            {
                name: '§23 Vertragsbruch',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            }
        ]
    },
];

penalties.forEach((ps) => {
    ps.penalties.forEach((p) => {
        p.category = ps.name;
    });
});

const selectedPenalties = ref<Array<SelectedPenalty>>([]);

const summary = ref<PenaltiesSummary>({
    fine: 0,
    detentionTime: 0,
    stvoPoints: 0,
    count: 0,
});

function calculate(e: SelectedPenalty) {
    const idx = selectedPenalties.value.findIndex((v) => v.penalty.category == e.penalty.category && v.penalty.name == e.penalty.name);
    let count = e.count;
    if (idx > -1) {
        const existing = selectedPenalties.value.at(idx)!;
        selectedPenalties.value[idx] = e;
        if (existing.count != e.count) {
            count = e.count - existing.count;
        }
    } else {
        selectedPenalties.value.push(e);
    }

    if (e.penalty.fine) {
        summary.value.fine += (count * e.penalty.fine);
    }
    if (e.penalty.detentionTime) {
        summary.value.detentionTime += (count * e.penalty.detentionTime);
    }
    if (e.penalty.stvoPoints) {
        summary.value.stvoPoints += (count * e.penalty.stvoPoints);
    }
    summary.value.count = +summary.value.count + +count;
}
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center pb-4">
                <div class="sm:flex-auto">
                    <div class="divide-y divide-white/10">
                        <dl class="mt-1 space-y-2 divide-y divide-white/10">
                            <Disclosure as="div" v-for="ps in penalties" :key="ps.name" class="pt-3" v-slot="{ open }">
                                <dt>
                                    <DisclosureButton class="flex w-full items-start justify-between text-left text-white">
                                        <span class="text-base font-semibold leading-7">{{ ps.name }}</span>
                                        <span class="ml-6 flex h-7 items-center">
                                            <PlusSmallIcon v-if="!open" class="h-6 w-6" aria-hidden="true" />
                                            <MinusSmallIcon v-else class="h-6 w-6" aria-hidden="true" />
                                        </span>
                                    </DisclosureButton>
                                </dt>
                                <DisclosurePanel as="dd" class="mt-2 px-4">
                                    <div class="flow-root mt-2">
                                        <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                                            <div class="inline-block min-w-full align-middle sm:px-6 lg:px-8">
                                                <table class="min-w-full divide-y divide-base-600">
                                                    <thead>
                                                        <tr>
                                                            <th scope="col"
                                                                class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                                                                Straftat
                                                            </th>
                                                            <th scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                                Geldstrafe
                                                            </th>
                                                            <th scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                                Haftzeit
                                                            </th>
                                                            <th scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                                StVo-Punkte
                                                            </th>
                                                            <th scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                                Sonstige
                                                            </th>
                                                            <th scope="col"
                                                                class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                                                Anzahl
                                                            </th>
                                                        </tr>
                                                    </thead>
                                                    <tbody class="divide-y divide-base-800">
                                                        <ListEntry v-for="penalty, idx in ps.penalties" :key="idx"
                                                            :penalty="penalty" @selected="calculate($event)" />
                                                    </tbody>
                                                </table>
                                            </div>
                                        </div>
                                    </div>
                                </DisclosurePanel>
                            </Disclosure>
                        </dl>
                    </div>
                </div>
            </div>
            <div class="relative">
                <div class="absolute inset-0 flex items-center" aria-hidden="true">
                    <div class="w-full border-t border-gray-300" />
                </div>
                <div class="relative flex justify-center">
                    <span class="bg-white px-3 text-base font-semibold leading-6 text-gray-900">Ergebnis</span>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <div class="text-neutral text-xl">
                            <Stats :summary="summary" />
                            <div class="mt-4">
                                <SummaryTable :selected-penalties="selectedPenalties" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
