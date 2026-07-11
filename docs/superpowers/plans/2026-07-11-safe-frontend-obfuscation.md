# Safe Frontend Obfuscation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Obfuscate production business-page chunks without breaking Vite lazy-route imports.

**Architecture:** Replace the source-stage `rollup-plugin-obfuscator` integration with a local Rollup `renderChunk` plugin. The plugin selects only dynamic chunks whose facade module is under `src/`, so the main router entry, preload runtime, and vendor chunks retain Vite's generated relative imports. A build verifier fails on source-alias dynamic imports and confirms the router entry retains the admin and timeline chunks.

**Tech Stack:** Vite 5, Rollup output hooks, `javascript-obfuscator`, Node.js built-in test runner.

**Implementation adjustment:** The repository ignores every `build/` directory,
so the helper and tests live in `frontend/scripts/`. Repeated tests established
that strong obfuscation can rewrite a dynamic-import specifier; the final helper
therefore skips any chunk containing `import()` and continues to obfuscate the
remaining business chunks.

---

### Task 1: Add and Test the Dynamic-Feature Chunk Obfuscator

**Files:**
- Create: `frontend/build/obfuscate-feature-chunks.mjs`
- Create: `frontend/build/obfuscate-feature-chunks.test.mjs`

- [ ] **Step 1: Write the failing selection and import-preservation tests**

```js
import assert from 'node:assert/strict'
import test from 'node:test'
import {
  createFeatureChunkObfuscator,
  shouldObfuscateChunk,
} from './obfuscate-feature-chunks.mjs'

const featureChunk = {
  type: 'chunk',
  fileName: 'js/Timeline-abc123.js',
  isEntry: false,
  isDynamicEntry: true,
  facadeModuleId: '/project/src/views/Timeline.vue',
}

test('selects only dynamic chunks whose facade is application source', () => {
  assert.equal(shouldObfuscateChunk(featureChunk), true)
  assert.equal(shouldObfuscateChunk({ ...featureChunk, isEntry: true }), false)
  assert.equal(shouldObfuscateChunk({
    ...featureChunk,
    facadeModuleId: '/project/node_modules/vue/dist/vue.runtime.esm-bundler.js',
  }), false)
  assert.equal(shouldObfuscateChunk({
    ...featureChunk,
    isDynamicEntry: false,
  }), false)
})

test('preserves relative dynamic import specifiers in an obfuscated feature chunk', () => {
  const plugin = createFeatureChunkObfuscator()
  const source = 'export const load = () => import("./Timeline-abc123.js")'
  const result = plugin.renderChunk(source, featureChunk)

  assert.notEqual(result, null)
  assert.notEqual(result.code, source)
  assert.match(result.code, /import\\s*\\(\\s*["']\\.\\/Timeline-abc123\\.js["']\\s*\\)/)
})
```

- [ ] **Step 2: Run the test and verify it fails because the module does not exist**

Run: `node --test build/obfuscate-feature-chunks.test.mjs`

Expected: `ERR_MODULE_NOT_FOUND` for `build/obfuscate-feature-chunks.mjs`.

- [ ] **Step 3: Implement the final-chunk-only plugin**

```js
import JavaScriptObfuscator from 'javascript-obfuscator'

export const obfuscationOptions = Object.freeze({
  compact: true,
  controlFlowFlattening: true,
  controlFlowFlatteningThreshold: 0.5,
  deadCodeInjection: true,
  deadCodeInjectionThreshold: 0.3,
  debugProtection: true,
  selfDefending: true,
  sourceMap: false,
  stringArray: false,
})

export function shouldObfuscateChunk(chunk) {
  return (
    chunk.type === 'chunk' &&
    chunk.isDynamicEntry &&
    !chunk.isEntry &&
    /[\\/]src[\\/]/.test(chunk.facadeModuleId || '')
  )
}

export function createFeatureChunkObfuscator() {
  return {
    name: 'feature-chunk-obfuscator',
    apply: 'build',
    renderChunk(code, chunk) {
      if (!shouldObfuscateChunk(chunk)) {
        return null
      }

      const result = JavaScriptObfuscator.obfuscate(code, {
        ...obfuscationOptions,
        inputFileName: chunk.fileName,
      })

      return {
        code: result.getObfuscatedCode(),
        map: null,
      }
    },
  }
}
```

This hook executes only after Vite has emitted final relative chunk specifiers. Do not use `rollup-plugin-obfuscator`, whose default `transform` hook runs on source modules.

- [ ] **Step 4: Run the test and verify it passes**

Run: `node --test build/obfuscate-feature-chunks.test.mjs`

Expected: 2 passing subtests and 0 failures.

- [ ] **Step 5: Commit the tested helper**

```bash
git add frontend/build/obfuscate-feature-chunks.mjs frontend/build/obfuscate-feature-chunks.test.mjs
git commit -m "feat: obfuscate lazy feature chunks after bundling"
```

### Task 2: Wire the Helper into the Production Build

**Files:**
- Modify: `frontend/vite.config.js:1-27`
- Modify: `frontend/package.json:5-41`
- Modify: `frontend/package-lock.json`

- [ ] **Step 1: Add the production dependency and remove the incompatible wrapper**

Replace the two relevant `devDependencies` entries with:

```json
"javascript-obfuscator": "^5.4.3",
"terser": "^5.44.1"
```

Remove `rollup-plugin-obfuscator`. Then run:

```bash
npm uninstall rollup-plugin-obfuscator
npm install --save-dev javascript-obfuscator@^5.4.3
```

The direct dependency is required because the local plugin imports it. The lockfile must no longer list `rollup-plugin-obfuscator` as a root dependency.

- [ ] **Step 2: Wire the `renderChunk` plugin into Vite**

Add this import after the existing Vite plugin imports:

```js
import { createFeatureChunkObfuscator } from './build/obfuscate-feature-chunks.mjs'
```

Add the plugin after `Components(...)` in the existing `plugins` array:

```js
createFeatureChunkObfuscator(),
```

Keep the current Terser options, including `sourcemap: false`, console/debugger removal, and top-level mangling. Do not re-add `rollup-plugin-obfuscator` or any plugin with a `transform` hook that obfuscates source modules.

- [ ] **Step 3: Run helper tests and a production build**

Run:

```bash
npm test
npm run build
```

Expected: the two helper tests pass; Vite builds without an unresolved `@/` import; `dist/js/index-*.js` remains the un-obfuscated route entry while lazy page chunks are processed by the `renderChunk` hook.

- [ ] **Step 4: Commit the build integration**

```bash
git add frontend/vite.config.js frontend/package.json frontend/package-lock.json
git commit -m "build: preserve Vite lazy imports during obfuscation"
```

### Task 3: Fail Builds that Reintroduce Source-Alias Dynamic Imports

**Files:**
- Create: `frontend/scripts/verify-production-build.mjs`
- Modify: `frontend/package.json:5-12`

- [ ] **Step 1: Write the build verifier**

```js
import { readFileSync, readdirSync } from 'node:fs'
import { resolve } from 'node:path'

const jsDirectory = resolve('dist/js')
const files = readdirSync(jsDirectory).filter((file) => file.endsWith('.js'))
const aliasDynamicImport = /\\bimport\\s*\\(\\s*["']@\\//
const invalidFiles = files.filter((file) =>
  aliasDynamicImport.test(readFileSync(resolve(jsDirectory, file), 'utf8')),
)

if (invalidFiles.length > 0) {
  throw new Error(
    `Production assets retain Vite alias dynamic imports: ${invalidFiles.join(', ')}`,
  )
}

const html = readFileSync(resolve("dist/index.html"), "utf8")
const entryMatch = html.match(/<script[^>]+src="\/js\/(index-[A-Za-z0-9_-]+\.js)"/)
const routerEntry = entryMatch?.[1]

if (!routerEntry) {
  throw new Error('Production router entry was not emitted')
}

const routerCode = readFileSync(resolve(jsDirectory, routerEntry), 'utf8')

for (const routeChunk of ['AdminLayout', 'Timeline']) {
  if (!new RegExp(`import\\(\\"\\./${routeChunk}-[A-Za-z0-9_-]+\\.js\\"\\)`).test(routerCode)) {
    throw new Error(`Router entry no longer references the ${routeChunk} lazy chunk`)
  }
}
```

- [ ] **Step 2: Make the verifier part of every production build**

Change the package scripts to:

```json
"build": "vite build && node scripts/verify-production-build.mjs",
"test": "node --test build/obfuscate-feature-chunks.test.mjs"
```

- [ ] **Step 3: Run final non-browser verification**

Run:

```bash
npm test
npm run build
```

Expected: the verifier exits 0, no generated JavaScript contains `import('@/`, and the router entry refers to relative hashed AdminLayout and Timeline chunk URLs.

- [ ] **Step 4: Verify browser lazy navigation against production assets**

Run: `npm run preview -- --host 127.0.0.1`

Open `http://127.0.0.1:4173/timeline` and `http://127.0.0.1:4173/admin` in a browser. Expected: both routes load without `Failed to resolve module specifier` in the console. Stop the preview server after recording the result.

- [ ] **Step 5: Commit the guardrail**

```bash
git add frontend/package.json frontend/scripts/verify-production-build.mjs
git commit -m "test: guard Vite lazy route imports in production assets"
```
