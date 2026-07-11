import JavaScriptObfuscator from 'javascript-obfuscator'

const dynamicImport = /\bimport\s*\(/
export const obfuscationManifestFile = 'obfuscation-manifest.json'

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
  const obfuscatedFacadeModuleIds = new Set()

  return {
    name: 'feature-chunk-obfuscator',
    apply: 'build',
    renderChunk(code, chunk) {
      if (!shouldObfuscateChunk(chunk) || dynamicImport.test(code)) {
        return null
      }

      const result = JavaScriptObfuscator.obfuscate(code, {
        ...obfuscationOptions,
        inputFileName: chunk.fileName,
      })

      obfuscatedFacadeModuleIds.add(chunk.facadeModuleId)

      return {
        code: result.getObfuscatedCode(),
        map: null,
      }
    },
    generateBundle(_, bundle) {
      const files = Object.values(bundle)
        .filter(
          (asset) =>
            asset.type === 'chunk' &&
            obfuscatedFacadeModuleIds.has(asset.facadeModuleId),
        )
        .map((asset) => asset.fileName)
        .sort()

      this.emitFile({
        type: 'asset',
        fileName: obfuscationManifestFile,
        source: JSON.stringify({
          files,
        }),
      })
    },
  }
}
