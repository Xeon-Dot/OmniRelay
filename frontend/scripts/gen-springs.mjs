#!/usr/bin/env node
/**
 * gen-springs.mjs — Material 3 Expressive motion-physics springs → CSS linear().
 *
 * M3 Expressive replaces duration+easing pairs with springs for component
 * interaction. CSS has no spring solver, so we sample the analytic response
 * of each spring and emit it as a `linear()` timing function, which browsers
 * interpolate as a piecewise-linear curve across the element's transition
 * duration.
 *
 * Source of the matrix: M3 motion scheme (spatial/effects × fast/default/slow),
 * stiffness + damping-ratio pairs, mass = 1. Spatial springs are under-damped
 * (ratio < 1 → overshoot); effects springs are critically damped (ratio = 1).
 *
 * Run: `node scripts/gen-springs.mjs` — writes src/styles/tokens/motion-springs.css.
 * Test: `node --test scripts/gen-springs.test.mjs`.
 */
import { writeFileSync, mkdirSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));

/** The M3 spring matrix. `stiffness` is k, `ratio` is the damping ratio. */
export const SPRINGS = {
  expressive: {
    "spatial-fast": { stiffness: 800, ratio: 0.6 },
    "spatial-default": { stiffness: 380, ratio: 0.8 },
    "spatial-slow": { stiffness: 200, ratio: 0.8 },
    "effects-fast": { stiffness: 3800, ratio: 1 },
    "effects-default": { stiffness: 1600, ratio: 1 },
    "effects-slow": { stiffness: 800, ratio: 1 },
  },
  standard: {
    "spatial-fast": { stiffness: 1400, ratio: 0.9 },
    "spatial-default": { stiffness: 700, ratio: 0.9 },
    "spatial-slow": { stiffness: 300, ratio: 0.9 },
    "effects-fast": { stiffness: 3800, ratio: 1 },
    "effects-default": { stiffness: 1600, ratio: 1 },
    "effects-slow": { stiffness: 800, ratio: 1 },
  },
};

/**
 * Analytic step response of a unit-mass spring released from rest at 0 toward 1.
 * Under-damped (ratio < 1): decaying cosine with a sine correction term.
 * Critically damped (ratio === 1): (1 + w t) e^{-w t} decay envelope.
 */
export function springResponse(t, stiffness, ratio) {
  if (t <= 0) return 0;
  const w = Math.sqrt(stiffness);
  if (ratio < 1) {
    const wd = w * Math.sqrt(1 - ratio * ratio);
    const decay = Math.exp(-ratio * w * t);
    return 1 - decay * (Math.cos(wd * t) + ((ratio * w) / wd) * Math.sin(wd * t));
  }
  if (ratio === 1) {
    return 1 - (1 + w * t) * Math.exp(-w * t);
  }
  // Over-damped — not part of the M3 matrix, but keep the math total.
  const wd = w * Math.sqrt(ratio * ratio - 1);
  const a = w * (ratio - Math.sqrt(ratio * ratio - 1));
  const b = w * (ratio + Math.sqrt(ratio * ratio - 1));
  return 1 - (b * Math.exp(-a * t) - a * Math.exp(-b * t)) / (b - a) + 0 * wd;
}

/** Seconds until the response stays within 0.1% of its target. */
function settleTime(stiffness, ratio) {
  const dt = 0.002;
  let t = 0;
  let insideSince = null;
  while (t < 5) {
    if (Math.abs(1 - springResponse(t, stiffness, ratio)) < 0.001) {
      if (insideSince === null) insideSince = t;
      // Must hold, not merely touch — an under-damped spring crosses 1 on the way.
      if (t - insideSince >= 0.05) return t;
    } else {
      insideSince = null;
    }
    t += dt;
  }
  return 5;
}

/**
 * Sample a spring's response at evenly spaced times — `linear()` distributes
 * its arguments uniformly across the timeline, so uniform sampling is required.
 * @returns {{durationMs: number, points: number[]}}
 */
export function springCurve(stiffness, ratio, { samples = 32 } = {}) {
  const duration = settleTime(stiffness, ratio);
  const points = [];
  for (let i = 0; i < samples; i++) {
    const t = (i / (samples - 1)) * duration;
    points.push(springResponse(t, stiffness, ratio));
  }
  // Pin the endpoints: linear() must start at rest and end on target.
  points[0] = 0;
  points[points.length - 1] = 1;
  return { durationMs: Math.max(1, Math.round(duration * 1000)), points };
}

const fmt = (n) => String(Number(n.toFixed(4)));

/** Render a sampled curve as a CSS `linear()` timing function. */
export function toLinear({ points }) {
  return `linear(${points.map(fmt).join(", ")})`;
}

export function build() {
  const lines = [
    "/*",
    " * motion-springs.css — GENERATED FILE, DO NOT EDIT BY HAND.",
    " * Source: scripts/gen-springs.mjs  (test: scripts/gen-springs.test.mjs)",
    " * M3 Expressive motion-physics springs sampled into CSS linear() curves.",
    " * spatial = under-damped (overshoots); effects = critically damped (never overshoots).",
    " */",
    "",
    ":root {",
  ];
  for (const [scheme, springs] of Object.entries(SPRINGS)) {
    const suffix = scheme === "expressive" ? "" : "-standard";
    lines.push(`/* ── ${scheme} scheme ── */`);
    for (const [name, spec] of Object.entries(springs)) {
      const curve = springCurve(spec.stiffness, spec.ratio);
      lines.push(
        `--m3-motion-${name}${suffix}: ${toLinear(curve)}; /* k=${spec.stiffness} r=${spec.ratio} */`,
      );
      lines.push(
        `--m3-motion-${name}${suffix}-duration: ${curve.durationMs}ms;`,
      );
    }
    lines.push("");
  }
  lines.push("}");
  return lines.join("\n");
}

const invokedDirectly =
  process.argv[1] && import.meta.url === pathToFileURL(resolve(process.argv[1])).href;
if (invokedDirectly) {
  const out = resolve(__dirname, "../src/styles/tokens/motion-springs.css");
  mkdirSync(dirname(out), { recursive: true });
  writeFileSync(out, build() + "\n", "utf8");
  for (const [scheme, springs] of Object.entries(SPRINGS)) {
    for (const [name, spec] of Object.entries(springs)) {
      const c = springCurve(spec.stiffness, spec.ratio);
      const peak = Math.max(...c.points);
      process.stdout.write(
        `${scheme}/${name}: ${c.durationMs}ms, peak ${peak.toFixed(4)}, ${c.points.length} samples\n`,
      );
    }
  }
  process.stdout.write(`→ ${out}\n`);
}
