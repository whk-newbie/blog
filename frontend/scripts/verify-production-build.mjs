import { readFileSync, readdirSync } from 'node:fs'
import { resolve } from 'node:path'

const distDirectory = resolve('dist')
const jsDirectory = resolve(distDirectory, 'js')
const files = readdirSync(jsDirectory).filter((file) => file.endsWith('.js'))
const aliasDynamicImport = /\bimport\s*\(\s*["']@\//
const invalidFiles = files.filter((file) =>
  aliasDynamicImport.test(readFileSync(resolve(jsDirectory, file), 'utf8')),
)

if (invalidFiles.length > 0) {
  throw new Error(
    `Production assets retain Vite alias dynamic imports: ${invalidFiles.join(', ')}`,
  )
}

const html = readFileSync(resolve(distDirectory, 'index.html'), 'utf8')
const entryMatch = html.match(
  /<script[^>]+src="\/js\/(index-[A-Za-z0-9_-]+\.js)"/,
)
const routerEntry = entryMatch?.[1]

if (!routerEntry) {
  throw new Error('Production router entry was not emitted')
}

const routerCode = readFileSync(resolve(jsDirectory, routerEntry), 'utf8')

for (const routeChunk of ['AdminLayout', 'Timeline']) {
  const lazyImport = new RegExp(
    `import\\(\\s*["']\\./${routeChunk}-[A-Za-z0-9_-]+\\.js["']\\s*\\)`,
  )

  if (!lazyImport.test(routerCode)) {
    throw new Error(
      `Router entry no longer references the ${routeChunk} lazy chunk`,
    )
  }
}
