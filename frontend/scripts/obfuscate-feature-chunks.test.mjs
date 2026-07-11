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

test('obfuscates feature chunks without dynamic imports', () => {
  const plugin = createFeatureChunkObfuscator()
  const source = 'export const page = "protected"'
  const result = plugin.renderChunk(source, featureChunk)

  assert.notEqual(result, null)
  assert.notEqual(result.code, source)
})

test('emits a manifest for obfuscated feature chunks', () => {
  const plugin = createFeatureChunkObfuscator()
  const emitted = []

  plugin.renderChunk('export const page = "protected"', featureChunk)
  plugin.generateBundle.call(
    {
      emitFile(asset) {
        emitted.push(asset)
      },
    },
    {},
    {
      'js/Timeline-final.js': {
        type: 'chunk',
        facadeModuleId: featureChunk.facadeModuleId,
        fileName: 'js/Timeline-final.js',
      },
    },
  )

  assert.deepEqual(emitted, [{
    type: 'asset',
    fileName: 'obfuscation-manifest.json',
    source: JSON.stringify({ files: ['js/Timeline-final.js'] }),
  }])
})

test('leaves feature chunks with dynamic imports unmodified', () => {
  const plugin = createFeatureChunkObfuscator()
  const source = 'export const load = () => import("./Timeline-abc123.js")'
  const result = plugin.renderChunk(source, featureChunk)

  assert.equal(result, null)
})
