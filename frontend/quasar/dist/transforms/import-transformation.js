// Quasar auto-import shim.
//
// The Quasar webpack loader asks this helper which module path should be used
// for a given component/directive/composable name (for example "QBtn").
// The default project scaffold shipped a stub that always returned the
// top-level "quasar" entry, which prevents tree-shaking and triggers runtime
// warnings such as "export 'default' was not found in 'quasar'" because the
// loader ends up generating default imports from the root package.
//
// Here we leverage Quasar's published import-map so every auto-imported symbol
// resolves to its specific source file (e.g. "QBtn" -> "quasar/src/components/btn/QBtn.js").
// This silences the warnings and restores proper tree-shaking behaviour.

let importMap

try {
  importMap = require('quasar/dist/transforms/import-map.json')
}
catch (error) {
  importMap = null
}

module.exports = function (name) {
  if (!name || importMap === null) {
    return 'quasar'
  }

  const mappedPath = importMap[name]

  if (typeof mappedPath !== 'string' || mappedPath.length === 0) {
    return 'quasar'
  }

  // Ensure we return a resolvable path inside the quasar package
  const normalized = mappedPath.replace(/^\.\//, '')
  return `quasar/${normalized}`
}
