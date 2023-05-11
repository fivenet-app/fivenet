<script lang="ts" setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue';
import { MinusSmallIcon, PlusSmallIcon } from '@heroicons/vue/24/outline';
import ListEntry from '~/components/penaltycalculator/ListEntry.vue';
import { Penalties, PenaltiesSummary, SelectedPenalty } from '~/utils/penalty';
import Stats from '~/components/penaltycalculator/Stats.vue';
import SummaryTable from './SummaryTable.vue';
import { useClipboard } from '@vueuse/core';
import { useNotificationsStore } from '~/store/notifications';

const { t, d } = useI18n();
const clipboard = useClipboard();

const notifications = useNotificationsStore();

const penalties: Penalties = [
    {
        name: "StGB",
        penalties: [
            {
                name: '§12 Mord',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 45,
                stvoPoints: 0
            },
            {
                name: '§13 Totschlag',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§13a Fahrlässige Tötung',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§14 Körperverletzung',
                description: 'Geldstrafe möglich',
                fine: 5000,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§14a Fahrlässige Körperverletzung',
                description: 'Geldstrafe möglich',
                fine: 5000,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§15 Gefährliche Körperverletzung',
                description: 'Geldstrafe möglich',
                fine: 20000,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§16 Schwere Körperverletzung',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§17 Unterlassene Hilfeleistung',
                description: 'Geldstrafe möglich',
                fine: 5000,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§18 Missbrauch von Notrufen',
                description: 'Geldstrafe möglich',
                fine: 5000,
                detentionTime: 5,
                stvoPoints: 0
            },
            {
                name: '§19 Diebstahl',
                description: 'Geldstrafe möglich, an Wert des Diebesgut orientieren',
                fine: 15000,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§20 Unterschlagung',
                description: 'Geldstrafe möglich',
                fine: 5000,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§21 Besonders schwerer Fall des Diebstahls',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§22 Wohnungseinbruchsdiebstahl',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§23 Hausfriedensbruch',
                description: 'Gelstrafe möglich',
                fine: 5000,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§24 Raub',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 5,
                stvoPoints: 0
            },
            {
                name: '§25 Schwerer Raub',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§26 Betrug',
                description: 'Geldstrafe möglich',
                fine: 10000,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§27 Menschenhandel',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§28 Freiheitsberaubung',
                description: 'Geldstrafe möglich',
                fine: 13000,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§29 Geiselnahme',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§30 Beleidigung, üble Nachrede und Verleumdunng',
                description: 'Geldstrafe möglich',
                fine: 1500,
                detentionTime: 5,
                stvoPoints: 0
            },
            {
                name: '§30a Abs. 1 Bedrohung',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§30a Abs. 2 Bedrohung',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§30a Abs. 4 Bedrohung',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 40,
                stvoPoints: 0
            },
            {
                name: '§31 Nötigung',
                description: 'Geldstrafe möglich',
                fine: 10000,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§32 Gefährlicher Eingriff in den Straßenverkehr',
                description: 'Geldstrafe möglich',
                fine: 10000,
                detentionTime: 15,
                stvoPoints: 2
            },
            {
                name: '§33 Gefährdung des Straßenverkehrs',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§34 Verbotenes Kraftfahrzeugsrennen',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§35 Abs. 1 Vorteilsgewährung',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§35 Abs. 2 Vorteilsgewährung',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 35,
                stvoPoints: 0
            },
            {
                name: '§35a Abs. 1 Vorteilsannahme',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 60,
                stvoPoints: 0
            },
            {
                name: '§35 Abs. 2 Vorteilsgewährung',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 35,
                stvoPoints: 0
            },
            {
                name: '§35a Abs. 1 Vorteilsannahme',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 60,
                stvoPoints: 0
            },
            {
                name: '§35a Abs. 2 Vorteilsannahme',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 60,
                stvoPoints: 0
            },
            {
                name: '§36 Bildung krimineller Vereinigungen',
                description: 'Geldstrafe möglich',
                fine: 30000,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§37 Amtsanmaßung',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§38 Unerlaubtes Entfernen vom Unfallort',
                description: 'Freiheitsstrafe nötig',
                fine: 5000,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§39 Fahren ohne Fahrerlaubnis',
                description: 'Geldstrafe möglich',
                fine: 15000,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§40 Strafvereitelung',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§41 Strafvereitelung im Amt',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§42 Landfriedensbruch',
                description: 'Geldstrafe möglich',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§43 Falsche uneidliche Aussage',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§45 Sachbeschädigung',
                description: 'Geldstrafe nötig, abhängig vom Sachwert',
                fine: 0,
                detentionTime: 0,
                stvoPoints: 0
            },
            {
                name: '§45a Abs. 1 Brandstiftung',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 35,
                stvoPoints: 0
            },
            {
                name: '§45a Abs. 2 Brandstiftung',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 20,
                stvoPoints: 0
            },
            {
                name: '§46 Terroristische Straftaten',
                description: 'NUR VON STAATSANWALTSCHAFT/RICHTER',
                fine: 0,
                detentionTime: 120,
                stvoPoints: 0
            },
            {
                name: '§47 Aufforderung zu terroristischen Straftaten und Gutheißung terroristischer Straftaten',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 80,
                stvoPoints: 0
            },
            {
                name: '§48 Mißbrauch der Amtsgewalt',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 40,
                stvoPoints: 0
            },
            {
                name: '§49 Sperrzonen',
                description: 'Bußgeld möglich. Haftzeit nötig bei Wiederholungstätern/bei Behinderung der Einsatzkräfte.',
                fine: 25000,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§50a Verbot von Vermummung im öffentlichen Raum',
                description: 'Bußgeld möglich. Freiheitsstrafe nötig bei Wiederholungstätern/wenn die Vermummung im Zusammenhang mit einer Straftat (wie Raub, Körperverletzung, Diebstahl, etc.) in der Öffentlichkeit getragen worden ist.',
                fine: 5000,
                detentionTime: 10,
                stvoPoints: 0
            },
            {
                name: '§50b Verbot von Vermummung auf staatlichem Gelände',
                description: 'Bußgeld möglich',
                fine: 10000,
                detentionTime: 15,
                stvoPoints: 0
            },
            {
                name: '§51 Verbotene Mitteilungen über Gerichtsverhandlungen',
                description: 'Bußgeld möglich',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§51a Weitergabe geheimer Informationen',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§52 Widerstand gegen Vollstreckungsbeamte',
                description: 'Freiheitsstrafe nötig. Wenn die Ausübung der Staatsgewalt unrechtmäßig ist, dann nicht strafbar | Bei einem besonders schwerem Fall bis zu 30 HE, ansonsten bis zu 25!',
                fine: 0,
                detentionTime: 30,
                stvoPoints: 0
            },
            {
                name: '§53 Umgehung der Haftzeit',
                description: 'Bußgeld möglich',
                fine: 0,
                detentionTime: 35,
                stvoPoints: 0
            },
            {
                name: '§54 Besitz von polizeilichen Mitteln',
                description: 'Bußgeld möglich',
                fine: 0,
                detentionTime: 25,
                stvoPoints: 0
            },
            {
                name: '§54a Besitz von illegalen Gegenständen',
                description: 'Freiheitsstrafe nötig',
                fine: 0,
                detentionTime: 25,
                stvoPoints: 0
            },
            {
                name: '§55 Unbefugter Gebrauch eines Fahrzeugs',
                description: 'Bußgeld möglich',
                fine: 0,
                detentionTime: 25,
                stvoPoints: 0
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
            },
        ],
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
            },
        ],
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
            },
        ],
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
            },
        ],
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
            },
        ],
    },
];

penalties.forEach((ps) => {
    ps.penalties.forEach((p) => {
        p.category = ps.name;
    });
});

const queryPenalities = ref<string>('');
const filteredPenalities = ref<typeof penalties>([]);
const selectedPenalties = ref<Array<SelectedPenalty>>([]);

const summary = ref<PenaltiesSummary>({
    fine: 0,
    detentionTime: 0,
    stvoPoints: 0,
    count: 0,
});

async function applyQuery(): Promise<void> {
    let newPenalties = structuredClone(penalties);

    newPenalties = newPenalties.map(ps => {
        const penalties = ps.penalties.map(p => {
            const show = p.name.toLowerCase().includes(queryPenalities.value.toLowerCase()) || p.description.toLowerCase().includes(queryPenalities.value.toLowerCase()) ? true : false;
            return {
                ...p,
                show
            }
        });

        const show = penalties.some(p => p.show === true);

        return {
            ...ps,
            penalties,
            show
        }
    });

    filteredPenalities.value = newPenalties;
}

watch(queryPenalities, async () => applyQuery());

function calculate(e: SelectedPenalty): void {
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

async function copyToClipboard(): Promise<void> {
    let text = t('components.penaltycalculator.title') + ` (` + d(new Date(), 'long') + `)

${t('components.penaltycalculator.fine')}: $${summary.value.fine}
${t('components.penaltycalculator.detention_time')}: ${summary.value.detentionTime} ${t('common.time_ago.month', summary.value.detentionTime)}
${t('components.penaltycalculator.stvo_points', 2)}: ${summary.value.stvoPoints}
${t('common.total_count')}: ${summary.value.count}
`;

    if (selectedPenalties.value.length > 0) {
        text += `
${t('components.penaltycalculator.crime', selectedPenalties.value.length)}:
`;

        selectedPenalties.value.forEach((v) => {
            text += `* ${v.penalty.category} - ${v.penalty.name} (${v.count}x)
`;
        });
    }

    notifications.dispatchNotification({
        title: t('notifications.penaltycalculator.title'),
        content: t('notifications.penaltycalculator.content'),
        type: 'info',
    });

    return clipboard.copy(text);
}

onMounted(async () => {
    applyQuery();
})
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="relative">
                <h3 class="text-2xl font-semibold leading-6">
                    {{ $t('components.penaltycalculator.title') }}
                </h3>
            </div>
            <div class="sm:flex sm:items-center pb-4">
                <div class="sm:flex-auto">
                    <div class="divide-y divide-white/10">
                        <div class="mt-5">
                            <input v-model="queryPenalities" type="text" name="search" id="search"
                                :placeholder="$t('common.filter')"
                                class="block w-full rounded-md border-0 py-1.5 pr-14 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6" />
                        </div>
                        <dl class="mt-5 space-y-2 divide-y divide-white/10">
                            <Disclosure as="div" v-for="ps in filteredPenalities" :key="ps.name" class="pt-3"
                                v-slot="{ open }" v-show="ps.show">
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
                                                                {{ $t('components.penaltycalculator.crime') }}
                                                            </th>
                                                            <th scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                                {{ $t('components.penaltycalculator.fine') }}
                                                            </th>
                                                            <th scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                                {{ $t('components.penaltycalculator.detention_time') }}
                                                            </th>
                                                            <th scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                                {{ $t('components.penaltycalculator.stvo_points', 2) }}
                                                            </th>
                                                            <th scope="col"
                                                                class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                                                                {{ $t('common.other') }}
                                                            </th>
                                                            <th scope="col"
                                                                class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                                                                {{ $t('common.count') }}
                                                            </th>
                                                        </tr>
                                                    </thead>
                                                    <tbody class="divide-y divide-base-800">
                                                        <ListEntry v-for="penalty, idx in ps.penalties" :key="idx"
                                                            :penalty="penalty" @selected="calculate($event)"
                                                            v-show="penalty.show" />
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
                    <span class="bg-white px-3 text-base font-semibold leading-6 text-gray-900">
                        {{ $t('common.result') }}
                    </span>
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
            <div class="flow-root mt-2">
                <div class="flex items-center">
                    <button type="button" @click="copyToClipboard()"
                        class="flex-1 rounded-md bg-info-700 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-info-600">
                        {{ $t('common.copy') }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
