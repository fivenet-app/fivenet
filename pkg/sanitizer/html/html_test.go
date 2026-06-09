package htmlsanitizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	t.Parallel()
	for _, run := range []struct {
		input  string
		result string
		msg    string
	}{
		{
			input:  "",
			result: "",
			msg:    "Empty, stays empty",
		},
		{
			input:  "<h2 class=\"test\">HELLO WORLD!</h2> This is a test.",
			result: "<h2>HELLO WORLD!</h2> This is a test.",
			msg:    "Simple tags with class attribute (unallowed class)",
		},
		{
			input:  "<h2 class=\"ql-link\">HELLO WORLD!</h2> This is a test.",
			result: "<h2>HELLO WORLD!</h2> This is a test.",
			msg:    "Simple tags with class attributes (allowed class)",
		},
		{
			input:  "<h2 style=\"color: #ffffff\" class=\"ql-link\">HELLO WORLD!</h2> <p><span>This is a test.</span></p>",
			result: "<h2 style=\"color: #ffffff\">HELLO WORLD!</h2> <p><span>This is a test.</span></p>",
			msg:    "Simple tags with style attribute",
		},
		{
			input:  "</h2>HELLO WORLD!</p> This<table is a> test</table.",
			result: "</h2>HELLO WORLD!</p> This<table> test",
			msg:    "Broken tags",
		},
		{
			input:  "<script src=\"example.com\"></script>HELLO WORLD!</p> This<table is a> test</table.",
			result: "HELLO WORLD!</p> This<table> test",
			msg:    "Make sure bad tag is removed, even with broken tags",
		},
		{
			input:  "FiveNet Fehler // \"Troll&#34;",
			result: "FiveNet Fehler // &#34;Troll&#34;",
			msg:    "Make sure that HTML entities are escaped",
		},
	} {
		assert.Equal(t, Sanitize(run.input), run.result, run.msg)
	}
}

func TestStripHTMLTags(t *testing.T) {
	t.Parallel()
	for _, run := range []struct {
		input  string
		result string
		msg    string
	}{
		{
			input:  "",
			result: "",
			msg:    "Empty, stays empty",
		},
		{
			input:  "<h2 class=\"test\">HELLO WORLD!</h2> This is a test.",
			result: "HELLO WORLD! This is a test.",
			msg:    "Simple tags with attributes",
		},
		{
			input:  "</h2>HELLO WORLD!</p> This<table is a> test</table.",
			result: "HELLO WORLD! This test",
			msg:    "Broken tags",
		},
		{
			input:  "<script src=\"example.com\"></script>HELLO WORLD!</p> This<table is a> test</table.",
			result: "HELLO WORLD! This test",
			msg:    "Make sure bad tag is removed, even with broken tags",
		},
		{
			input:  "FiveNet Fehler // &#34;Troll&#34;",
			result: "FiveNet Fehler // \"Troll\"",
			msg:    "Make sure that HTML entities are unescaped",
		},
	} {
		assert.Equal(t, StripHTMLTags(run.input), run.result, run.msg)
	}
}

func TestSanitizePreservesMapBlockAndPenaltyCalculator(t *testing.T) {
	t.Parallel()

	mapHTML := `<span data-embed="map" data-map-x="12.34" data-map-y="56.78" data-map-zoom="3" data-map-postal="12345" data-map-layer="postal"></span>`
	mapOut := Sanitize(mapHTML)
	assert.Contains(t, mapOut, `<span`, "map block tag should survive sanitization: %s", mapOut)
	assert.Contains(t, mapOut, `data-embed="map"`)
	assert.Contains(t, mapOut, `data-map-x="12.34"`)
	assert.Contains(t, mapOut, `data-map-y="56.78"`)
	assert.Contains(t, mapOut, `data-map-zoom="3"`)
	assert.Contains(t, mapOut, `data-map-postal="12345"`)
	assert.Contains(t, mapOut, `data-map-layer="postal"`)

	penaltyHTML := `<div data-type="penalty-calculator" data-embed="penalty-calculator"></div>`
	penaltyOut := Sanitize(penaltyHTML)
	assert.Contains(
		t,
		penaltyOut,
		`<div`,
		"penalty calculator tag should survive sanitization: %s",
		penaltyOut,
	)
	assert.Contains(t, penaltyOut, `data-type="penalty-calculator"`)
	assert.Contains(t, penaltyOut, `data-embed="penalty-calculator"`)
}

func TestSanitizeRejectsNonExportedMapBlockAndPenaltyCalculatorValues(t *testing.T) {
	t.Parallel()

	mapOut := Sanitize(
		`<span data-embed="map-legacy" data-map-x="12.34" data-map-y="56.78" data-map-zoom="3"></span>`,
	)
	assert.NotContains(t, mapOut, `data-embed="map-legacy"`)
	assert.Contains(
		t,
		mapOut,
		`<span`,
		"map block wrapper should still survive sanitization: %s",
		mapOut,
	)

	penaltyOut := Sanitize(
		`<div data-type="penaltyCalculator" data-embed="penaltyCalculator"></div>`,
	)
	assert.NotContains(t, penaltyOut, `data-type="penaltyCalculator"`)
	assert.NotContains(t, penaltyOut, `data-embed="penaltyCalculator"`)
	assert.Contains(
		t,
		penaltyOut,
		`<div`,
		"penalty calculator wrapper should still survive sanitization: %s",
		penaltyOut,
	)
}
