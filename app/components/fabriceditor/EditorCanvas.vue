<script setup lang="ts">
import type { FabricObject } from 'fabric';
import { useFabricEditor } from '~/composables/useFabricEditor';

const canvasEl = ref<HTMLCanvasElement | null>(null);
// Get composable state and methods
const { canvas, documentSize, initCanvas } = useFabricEditor();

onMounted(async () => {
    // Initialize Fabric canvas
    const canvasElement = canvasEl.value!;
    // Set canvas dimensions to container size
    const parent = canvasElement.parentElement;
    const width = parent?.clientWidth || 800;
    const height = parent?.clientHeight || 600;

    initCanvas(canvasElement, {
        width,
        height,
    });

    const fabric = await import('fabric');
    const loadedSvg = await fabric.loadSVGFromString(
        `<?xml version="1.0"?>
<svg width="1080" height="768.0000000000001" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg">
 <title>arbeitsunfähigkeitsbescheinigung</title>
 <defs>
  <pattern height="8" id="noise" patternUnits="userSpaceOnUse" width="8" x="0" y="0">
   <rect fill="#3b3b3b" height="8" id="svg_1" width="8"/>
   <circle cx="1" cy="1" fill="#f4f4f4" id="svg_2" r="1"/>
   <circle cx="6" cy="3" fill="#f4f4f4" id="svg_3" r="1"/>
   <circle cx="3" cy="6" fill="#f4f4f4" id="svg_4" r="1"/>
  </pattern>
 </defs>
 <g class="layer">
  <title>Layer 1</title>
  <rect fill="#f4e6b3" height="768" id="svg_65" transform="matrix(1, 0, 0, 1, 0, 0)" width="1080" x="-0.5" y="-0.5"/>
  <rect fill="none" height="300" id="svg_6" stroke="#cc2a2a" stroke-width="2" width="575" x="42" y="38"/>
  <line id="svg_7" stroke="#cc2a2a" stroke-width="2" x1="42" x2="617" y1="102" y2="102"/>
  <line id="svg_8" stroke="#cc2a2a" stroke-width="2" transform="matrix(1, 0, 0, 1, 0, 0)" x1="42" x2="617" y1="196" y2="196"/>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="10" id="svg_9" x="45" y="50">Krankenkasse bzw. Kostenträger</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="10" id="svg_10" x="45" y="115">Name, Vorname des Versicherten</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="10" id="svg_11" transform="matrix(1, 0, 0, 1, 0, 0)" x="485" y="140">geb. am</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="10" id="svg_12" transform="matrix(1, 0, 0, 1, 0, 0)" x="45" y="210">Kostenträgerkennung</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="10" id="svg_13" transform="matrix(1, 0, 0, 1, 0, 0)" x="246.37" y="210">Versicherten-Nr.</text>
  <rect fill="url(#noise)" height="69.58" id="svg_14" stroke="#cc2a2a" stroke-width="2" transform="matrix(1, 0, 0, 1, 0, 0)" width="200" x="42" y="268.25"/>
  <rect fill="url(#noise)" height="71.93" id="svg_15" stroke="#cc2a2a" stroke-width="2" transform="matrix(1, 0, 0, 1, 0, 0)" width="189" x="428" y="196.45"/>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="10" id="svg_16" transform="matrix(1, 0, 0, 1, 0, 0)" x="245" y="280">Arzt-Nr.</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="10" id="svg_17" transform="matrix(1, 0, 0, 1, 0, 0)" x="410" y="280">Datum</text>
  <line id="svg_18" stroke="#cc2a2a" stroke-width="2" transform="matrix(1, 0, 0, 1, 0, 0)" x1="241.69" x2="427.07" y1="268.35" y2="268.35"/>
  <line id="svg_19" stroke="#cc2a2a" stroke-width="2" x1="242.07" x2="242.07" y1="196.91" y2="267.31"/>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="38" font-weight="700" id="svg_20" x="664" y="76">Arbeitsunfähigkeits-</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="38" font-weight="700" id="svg_21" transform="matrix(1, 0, 0, 1, 0, 0)" x="664" y="114">bescheinigung</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="38" font-weight="700" id="svg_22" x="968" y="114">1</text>
  <rect fill="none" height="34" id="svg_23" stroke="#cc2a2a" stroke-width="2" width="34" x="664" y="140"/>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_24" transform="matrix(1, 0, 0, 1, 0, 0)" x="708" y="166">Erstbescheinigung</text>
  <rect fill="none" height="34" id="svg_25" stroke="#cc2a2a" stroke-width="2" width="34" x="664" y="194"/>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_26" x="708" y="220">Folgebescheinigung</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="17" id="svg_27" transform="matrix(1, 0, 0, 1, 0, 0)" x="665" y="260">
   <tspan dy="0" font-size="22.13" id="svg_28" x="663.77">Der angegebenen Krankenkasse wird</tspan>
   <tspan dy="24" font-size="22.13" id="svg_29" x="663.77">unverzüglich eine Bescheinigung über</tspan>
   <tspan dy="24" font-size="22.13" id="svg_30" x="663.77">die Arbeitsunfähigkeit mit Angaben</tspan>
   <tspan dy="24" font-size="22.13" id="svg_31" x="663.77">über die Diagnose sowie die</tspan>
   <tspan dy="24" font-size="22.13" id="svg_32" x="663.77">voraussichtliche Dauer der</tspan>
   <tspan dy="24" font-size="22.13" id="svg_33" x="663.77">Arbeitsunfähigkeit übersandt.</tspan>
  </text>
  <rect fill="none" height="34" id="svg_34" stroke="#cc2a2a" stroke-width="2" width="34" x="54" y="378"/>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_35" transform="matrix(1, 0, 0, 1, 0, 0)" x="100" y="395">Arbeitsunfall, Arbeitsunfall-</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_36" transform="matrix(1, 0, 0, 1, 0, 0)" x="100" y="420">folgen, Berufskrankheit</text>
  <rect fill="none" height="35" id="svg_37" stroke="#cc2a2a" stroke-width="2" transform="matrix(1, 0, 0, 1, 0, 0)" width="35" x="365" y="380"/>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_38" transform="matrix(1, 0, 0, 1, 0, 0)" x="410" y="395">dem Durchgangsarzt</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_39" transform="matrix(1, 0, 0, 1, 0, 0)" x="410" y="415">zugewiesen</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_40" transform="matrix(1, 0, 0, 1, 0, 0)" x="55" y="485">arbeitsunfähig seit</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_46" x="55" y="540">voraussichtlich arbeitsunfähig</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_47" transform="matrix(1, 0, 0, 1, 0, 0)" x="55" y="565">bis einschließlich oder letzter</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_48" transform="matrix(1, 0, 0, 1, 0, 0)" x="55" y="590">Tag der Arbeitsunfähigkeit</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" id="svg_54" transform="matrix(1, 0, 0, 1, 0, 0)" x="55" y="640">festgestellt am</text>
  <rect fill="none" height="45" id="svg_60" stroke="#cc2a2a" stroke-width="2" transform="matrix(1, 0, 0, 1, 0, 0)" width="570" x="45" y="690"/>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="20" font-weight="bold" id="svg_61" x="115" y="720">Ausfertigung zur Vorlage beim Arbeitgeber</text>
  <rect fill="none" height="282" id="svg_62" stroke="#cc2a2a" stroke-width="2" width="381" x="664" y="404"/>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="18" id="svg_63" transform="matrix(1, 0, 0, 1, 0, 0)" x="677" y="676">Vertragsarztstempel / Unterschrift des Arztes</text>
  <text fill="#cc2a2a" font-family="Arial, Helvetica, sans-serif" font-size="18" id="svg_64" transform="matrix(1, 0, 0, 1, 0, 0)" x="896" y="726">Muster 1b (1.2018)</text>
  <rect fill="none" height="45" id="svg_67" stroke="#cc2a2a" stroke-width="2" transform="matrix(1, 0, 0, 1, 0, 0)" width="260" x="355" y="455"/>
  <rect fill="none" height="45" id="svg_68" stroke="#cc2a2a" stroke-width="2" width="260" x="355" y="530"/>
  <rect fill="none" height="45" id="svg_69" stroke="#cc2a2a" stroke-width="2" transform="matrix(1, 0, 0, 1, 0, 0)" width="260" x="355" y="605"/>
  <line id="svg_70" stroke="#cc2a2a" stroke-width="2" transform="matrix(1, 0, 0, 1, 0, 0)" x1="406.15" x2="406.15" y1="269.36" y2="337.12"/>
 </g>
</svg>`,
    );

    if (loadedSvg.objects) {
        if (loadedSvg.options['width'] && loadedSvg.options['height']) {
            documentSize.value.width = parseInt(loadedSvg.options['width']);
            documentSize.value.height = parseInt(loadedSvg.options['height']);
        }

        canvas.value?.add(...loadedSvg.objects.filter((obj): obj is FabricObject => !!obj));
    }
});
</script>

<template>
    <!-- The canvas element for Fabric.js -->
    <canvas ref="canvasEl" class="h-full w-full"></canvas>
</template>
