package app

import "Gofra_Market/internal/docs"

// Переэкспортирует описание пакета из пакета docs для обратной совместимости.
type PackageDoc = docs.PackageDoc

// Переадресует вызов в пакет docs для обратной совместимости.
var NewPackageDocsHandler = docs.NewPackageDocsHandler
