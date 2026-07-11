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
  assert.match(result.code, /import\s*\(\s*["']\.\/Timeline-abc123\.js["']\s*\)/)
})
