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

test('leaves feature chunks with dynamic imports unmodified', () => {
  const plugin = createFeatureChunkObfuscator()
  const source = 'export const load = () => import("./Timeline-abc123.js")'
  const result = plugin.renderChunk(source, featureChunk)

  assert.equal(result, null)
})
