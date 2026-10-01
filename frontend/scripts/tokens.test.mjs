// Token-layer integrity tests (node --test, zero deps).
// Guards the two failure modes the plan's Review Focus names: a role missing
// from one theme block (a light hex leaking into dark), and hard-coded color
// bypassing the palette (drift the spec forbids).
import test from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), "../src/styles/tokens");
const read = (f) => readFileSync(resolve(ROOT, f), "utf8");

/** Extract `--name: value;` pairs from the first block whose selector contains `sel`. */
function blockVars(css, sel) {
  const start = css.indexOf(sel);
  assert.ok(start >= 0, `selector not found: ${sel}`);
  const open = css.indexOf("{", start);
  const close = css.indexOf("}", open);
  assert.ok(open >= 0 && close > open, `unterminated block for ${sel}`);
  const body = css.slice(open + 1, close);
  const out = new Map();
  for (const m of body.matchAll(/(--[\w-]+)\s*:\s*([^;]+);/g)) out.set(m[1], m[2].trim());
  return out;
}

const light = () => blockVars(read("color.css"), ':root[data-theme="light"]');
const dark = () => blockVars(read("color.css"), ':root[data-theme="dark"]');

test("both theme blocks define exactly the same 49 M3 color roles", () => {
  const l = [...light().keys()].filter((k) => k.startsWith("--m3-color-"));
  const d = [...dark().keys()].filter((k) => k.startsWith("--m3-color-"));
  assert.equal(l.length, 49, `light has ${l.length} roles`);
  assert.equal(d.length, 49, `dark has ${d.length} roles`);
  assert.deepEqual(l.sort(), d.sort(), "theme blocks disagree on role names");
});

test("app semantic aliases exist in both themes", () => {
  const aliases = ["--or-success", "--or-success-container", "--or-on-success-container",
    "--or-warning", "--or-warning-container", "--or-on-warning-container"];
  for (const [name, block] of [["light", light()], ["dark", dark()]]) {
    for (const a of aliases) assert.ok(block.has(a), `${name} missing ${a}`);
  }
});

test("every color role resolves through the palette, never a literal hex", () => {
  const css = read("color.css");
  const palette = read("palette.css");
  const declared = new Set([...palette.matchAll(/(--m3-palette-[\w-]+)\s*:/g)].map((m) => m[1]));
  const used = new Set();
  for (const m of css.matchAll(/(--m3-color-[\w-]+)\s*:\s*([^;]+);/g)) {
    const [, name, value] = m;
    assert.ok(!/#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})\b/.test(value), `${name} hard-codes ${value}`);
    const ref = value.match(/^var\((--m3-palette-[\w-]+)\)$/);
    assert.ok(ref, `${name} must be a single var(--m3-palette-*), got: ${value}`);
    used.add(ref[1]);
  }
  const missing = [...used].filter((v) => !declared.has(v));
  assert.deepEqual(missing, [], `palette vars referenced but never declared: ${missing.join(", ")}`);
});

test("type scale ships 15 baseline and 15 emphasized roles", () => {
  const css = read("type.css");
  const styles = ["display-large", "display-medium", "display-small",
    "headline-large", "headline-medium", "headline-small",
    "title-large", "title-medium", "title-small",
    "body-large", "body-medium", "body-small",
    "label-large", "label-medium", "label-small"];
  for (const s of styles) {
    assert.ok(css.includes(`--m3-typescale-${s}:`), `missing baseline ${s}`);
    assert.ok(css.includes(`--m3-typescale-emphasized-${s}:`), `missing emphasized ${s}`);
  }
  // Emphasized = same metrics, one weight step heavier (spec §3.5).
  assert.match(css, /--m3-typescale-emphasized-label-large:\s*700/, "emphasized label-large must be bold");
  assert.match(css, /--m3-typescale-emphasized-body-large:\s*500/, "emphasized body-large must be medium");
});

test("shape scale ships the 10 M3/M3E corners", () => {
  const css = read("shape.css");
  const expected = { none: "0", xs: "4px", sm: "8px", md: "12px", lg: "16px",
    xl: "20px", "2xl": "28px", "3xl": "32px", "4xl": "48px", full: "9999px" };
  for (const [name, value] of Object.entries(expected)) {
    assert.ok(css.includes(`--m3-shape-${name}: ${value}`), `--m3-shape-${name} must be ${value}`);
  }
});

test("index.css is the single entry and imports every token file", () => {
  const css = read("index.css");
  const files = ["palette.css", "color.css", "type.css", "shape.css", "motion.css",
    "elevation.css", "space.css", "state.css"];
  for (const f of files) assert.ok(css.includes(`@import "./${f}"`), `index.css must import ${f}`);
  // Springs are reached through motion.css; index must not reach past it.
  assert.ok(!css.includes("./motion-springs.css"), "springs are imported by motion.css only");
});

test("reduced-motion gate is present and global", () => {
  const css = read("motion.css");
  assert.match(css, /@media \(prefers-reduced-motion: reduce\)/);
  assert.match(css, /transition-duration:\s*1ms\s*!important/);
  assert.match(css, /animation-iteration-count:\s*1\s*!important/);
});
