export default function colorMixTransparencyFallback(opts = {}) {
    const {
        // Keep the original color-mix line
        preserve = true,
        rgbSuffix = '-rgb',
        props = [
            'color',
            'background',
            'background-color',
            'border-color',
            'outline-color',
            'text-decoration-color',
            'caret-color',
            'fill',
            'stroke',
        ],
    } = opts;

    const propSet = new Set(props.map((p) => p.toLowerCase()));

    // color-mix(in <space>, var(--name)N%, transparent) OR transparent, var(--name)N%
    // - captures var name and percentage regardless of spacing
    const MIX_RE =
        /color-mix\(\s*in\s+[a-z0-9-]+\s*,\s*(?:var\(\s*(--[\w-]+)\s*\)\s*([0-9.]+%)\s*,\s*transparent|transparent\s*,\s*var\(\s*(--[\w-]+)\s*\)\s*([0-9.]+%))\s*\)/gi;

    function computeFallbackValue(value) {
        let val = value;

        if (val.includes('rgb(255, 255, 255)')) {
            val = val.replace('rgb(255, 255, 255)', 'var(--color-white)');
        } else if (val.includes('rgb(0, 0, 0)')) {
            val = val.replace('rgb(0, 0, 0)', 'var(--color-black)');
        }

        val = val.replace('-info-', '-blue-');
        val = val.replace('-warning-', '-yellow-');
        val = val.replace('-error-', '-red-');
        val = val.replace('-success-', '-green-');
        val = val.replace('-neutral', '-old-neutral');

        let matched = false;
        const out = val.replace(MIX_RE, (_, n1, p1, n2, p2) => {
            matched = true;
            const name = n1 || n2;
            const pct = p1 || p2;
            return `rgb(var(${name}${rgbSuffix}), ${pct})`;
        });

        return matched ? out : null;
    }

    function findTopmostSupports(node) {
        // returns {top: AtRule, container: Parent} or null
        let cur = node.parent;
        let topSupports = null;
        while (cur) {
            if (cur.type === 'atrule' && cur.name === 'supports') {
                topSupports = cur; // keep walking to find the topmost
            }
            cur = cur.parent;
        }
        if (!topSupports) return null;
        return { top: topSupports, container: topSupports.parent };
    }

    function ensureRuleBefore(at, rule, prop, value) {
        // Avoid dupes: if an identical selector+prop exists just before, skip
        const parent = at.parent || at; // usually at.parent exists
        let hasSame = false;

        // Build a minimal clone of the rule with just the fallback decl
        const clone = rule.clone({ nodes: [] });
        clone.append({ prop, value });

        // Optional light duplicate check among previous siblings
        let prev = at.prev();
        while (prev && prev.type === 'comment') prev = prev.prev(); // skip comments
        if (prev && prev.type === 'rule' && prev.selector === clone.selector) {
            prev.walkDecls(prop, (d) => {
                if (d.value === value) hasSame = true;
            });
        }
        if (!hasSame) parent.insertBefore(at, clone);
    }

    return {
        postcssPlugin: 'postcss-color-mix-transparency-fallback',
        Declaration(decl) {
            if (!decl.value || decl.value.indexOf('color-mix(') === -1) return;
            if (decl.value.indexOf('transparent') === -1) return;

            if (!propSet.has(decl.prop.toLowerCase())) return;

            const fallbackVal = computeFallbackValue(decl.value);
            if (!fallbackVal) return;

            const sup = findTopmostSupports(decl);
            if (sup && sup.container) {
                // Insert a sibling rule (with same selector) before the topmost @supports
                if (decl.parent.type !== 'rule') return; // safety: we need a rule to copy
                ensureRuleBefore(sup.top, decl.parent, decl.prop, fallbackVal);

                if (!preserve) {
                    // If not preserving, also replace the original value inside @supports
                    decl.value = fallbackVal;
                }
            } else {
                // No @supports ancestors → insert fallback right before this decl
                if (preserve) decl.cloneBefore({ value: fallbackVal });
                else decl.value = fallbackVal;
            }
        },
    };
}

export const postcss = true;
