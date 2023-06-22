---
title: "Dispatch Center (Leitstelle)"
---

## Definition

* Leitstelle
    * Kontext:
        * Eine für exekutive (PD, FIB, DOJ)
        * LSMD kriegt eine separate Leitstelle.
    * "Gruppierungen" von Jobs
    * Zwei Modi
        * Ein simpler "jeder" kann Dispatches annehmen und anfahren
        * (Automatischer) Disponent
            * Personal das verfügbar ist, werden automatisch naheliegende Dispatches "vorgeschlagen"/ "assigned".
* Squad
    * Eine oder mehrere Personen aus dem gleichen Job (e.g., ein Wagen mit 2 PDler ist ein Squad).
    * Von der Leitung "einmalig" die Squads erstellt werden können (damit ist nicht das zuweisen von Leuten gemeint).
* Aktion
    * Relation: Many Squads to Many Aktionen
    * Beispiele:
        * (Wers möchte) "Ich verbringe TV ins SG"
        * Verkehrskontrolle
        * Großschießerei
        *
* Dispatch
    * Relation: Many Squads to Many Dispatches
    * Beispiele:
        * Generell: "Positiongebundene Aktion"
        * PD:
            * Schiererein
            * Razzia
        * LSMD:
            * Notruf
            * Manueller Notruf (Freitext; Kommentar an mich, Emojis automatisch entfernen)
* Status
    * Squad
        * Verfügbar
        * Beschäftigt => "10-6"
        * Pause => "Code 7"
        * Abgemeldet
    * Dispatch + Aktion
        * Offen
        * In Bearbeitung
            * "Status" Update - Freitext (vs Code)
        * Abgeschlossen (Grundabfrage)
* Dispatch Type
        * Broadcast -> Alle sehen diesen mit Warnung aufpoppen.
            * Egal ob in Pause oder nicht.
* Codes
    * Squad
        * Code 7 => Pause
    * Dispatch
        * Code 4 => Dispatch/Aktion ist abgeschlossen.
        * Code 10 => "Verstärkung benötigt"
        * 11-99/ Code 99 => "Dringend Verstärkung" aka Panickbutton
* Livemap
    * Sperrzonen sichtbar machen
        * In die Datenbank schreiben.
    * Panickbutton sichtbar machen
        * "Broadcast" Dispatch der allen angezeigt wird.
        * Icon: [https://pictogrammers.com/library/mdi/icon/car-brake-alert/](https://pictogrammers.com/library/mdi/icon/car-brake-alert/)
