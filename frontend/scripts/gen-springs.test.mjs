// Tests for the M3 Expressive spring generator (node --test, zero deps).
// Pins the M3 motion-physics matrix from the spec (§3.7): stiffness/damping
// pairs, overshoot behavior, and the shape of the emitted CSS linear() curve.
import test from "node:test";
import assert from "node:assert/strict";
import { SPRINGS, springResponse, springCurve, toLinear } from "./gen-springs.mjs";

// M3 Expressive scheme — spec §3.7 table.
const EXPRESSIVE = {
  "spatial-fast": [800, 0.6],
  "spatial-default": [380, 0.8],
  "spatial-slow": [200, 0.8],
  "effects-fast": [3800, 1],
  "effects-default": [1600, 1],
  "effects-slow": [800, 1],
};

const STANDARD = {
  "spatial-fast": [1400, 0.9],
  "spatial-default": [700, 0.9],
  "spatial-slow": [300, 0.9],
  "effects-fast": [3800, 1],
  "effects-default": [1600, 1],
  "effects-slow": [800, 1],
};

test("SPRINGS exposes the M3 matrix for both schemes", () => {
  assert.deepEqual(Object.keys(SPRINGS.expressive).sort(), Object.keys(EXPRESSIVE).sort());
  assert.deepEqual(Object.keys(SPRINGS.standard).sort(), Object.keys(STANDARD).sort());
  for (const [name, [k, r]] of Object.entries(EXPRESSIVE)) {
    assert.equal(SPRINGS.expressive[name].stiffness, k, `expressive/${name} stiffness`);
    assert.equal(SPRINGS.expressive[name].ratio, r, `expressive/${name} damping ratio`);
  }
  for (const [name, [k, r]] of Object.entries(STANDARD)) {
    assert.equal(SPRINGS.standard[name].stiffness, k, `standard/${name} stiffness`);
    assert.equal(SPRINGS.standard[name].ratio, r, `standard/${name} damping ratio`);
  }
});

test("springResponse starts at rest and settles within 0.1%", () => {
  for (const scheme of Object.values(SPRINGS)) {
    for (const [name, spec] of Object.entries(scheme)) {
      const duration = springCurve(spec.stiffness, spec.ratio).durationMs;
      assert.ok(Math.abs(springResponse(0, spec.stiffness, spec.ratio)) < 1e-9, `${name}: x(0)`);
      const end = springResponse(duration / 1000, spec.stiffness, spec.ratio);
      assert.ok(Math.abs(1 - end) < 0.001, `${name}: not settled at ${duration}ms (x=${end})`);
    }
  }
});

test("underdamped spatial springs overshoot; critically damped effects springs do not", () => {
  for (const [name, spec] of Object.entries(SPRINGS.expressive)) {
    const { points } = springCurve(spec.stiffness, spec.ratio);
    const peak = Math.max(...points);
    if (spec.ratio < 1) {
      // Analytic peak overshoot = exp(-pi*r/sqrt(1-r^2)).
      const expected = Math.exp((-Math.PI * spec.ratio) / Math.sqrt(1 - spec.ratio ** 2));
      assert.ok(peak > 1 + expected * 0.9, `${name}: expected overshoot ~${(expected * 100).toFixed(1)}%, got ${((peak - 1) * 100).toFixed(1)}%`);
    } else {
      assert.ok(peak <= 1.0001, `${name}: critically damped must not overshoot (peak=${peak})`);
    }
  }
});

test("springCurve emits at least 24 evenly spaced samples from 0 to 1", () => {
  const curve = springCurve(800, 0.6);
  assert.ok(curve.points.length >= 24, `only ${curve.points.length} samples`);
  assert.ok(Math.abs(curve.points[0]) < 1e-6, "first sample must be 0");
  assert.ok(Math.abs(curve.points[curve.points.length - 1] - 1) < 0.001, "last sample must settle at 1");
  assert.ok(curve.durationMs > 0, "duration must be positive");
});

test("toLinear renders a CSS linear() with 0 first and 1 last", () => {  const css = toLinear(springCurve(800, 0.6));
  assert.match(css, /^linear\(/, `must start with linear(: ${css.slice(0, 20)}`);
  assert.match(css, /\)$/, "must close paren");
  const nums = css.slice(7, -1).split(",").map((v) => Number(v.trim()));
  assert.equal(nums.length, springCurve(800, 0.6).points.length, "sample count preserved");
  assert.ok(Math.abs(nums[0]) < 1e-6, `first value ${nums[0]} must be 0`);
  assert.ok(Math.abs(nums[nums.length - 1] - 1) < 0.001, `last value ${nums[nums.length - 1]} must be 1`);
  // Must be parseable by a browser: no NaN, reasonable precision.
  assert.ok(nums.every((n) => Number.isFinite(n)), "no NaN/Infinity in curve");
  assert.ok(css.length < 4000, `curve too long for a CSS declaration (${css.length} chars)`);
});

test("generated stylesheet emits a matching duration token for every spring", async () => {
  const { build } = await import("./gen-springs.mjs");
  const css = build();
  // A spring curve must run to completion: truncating it mid-overshoot reads
  // as a stutter, so each curve carries its own settle duration.
  const open = css.indexOf(":root {");
  assert.ok(open >= 0, "declarations must be wrapped in a :root block");
  assert.equal(css.indexOf(":root {", open + 1), -1, "exactly one :root block");
  for (const [scheme, springs] of Object.entries(SPRINGS)) {
    const suffix = scheme === "expressive" ? "" : "-standard";
    for (const [name, spec] of Object.entries(springs)) {
      const curve = springCurve(spec.stiffness, spec.ratio);
      const decl = `--m3-motion-${name}${suffix}-duration: ${curve.durationMs}ms;`;
      assert.ok(css.includes(decl), `missing ${decl}`);
      assert.ok(
        css.includes(`--m3-motion-${name}${suffix}: linear(`),
        `missing curve --m3-motion-${name}${suffix}`,
      );
    }
  }
  assert.ok(css.trimEnd().endsWith("}"), "must close the :root block");
});
