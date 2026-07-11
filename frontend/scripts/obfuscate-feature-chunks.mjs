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
