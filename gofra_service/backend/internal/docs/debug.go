package docs

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"go/doc"
	"golang.org/x/tools/go/packages"
)

// Описывает краткую документацию пакета для debug-эндпоинта.
type PackageDoc struct {
	ImportPath string `json:"import_path"`
	Name       string `json:"name"`
	Synopsis   string `json:"synopsis"`
	Doc        string `json:"doc"`
}

// Возвращает хэндлер, отдающий документацию пакетов в формате JSON.
func NewPackageDocsHandler(root string) (gin.HandlerFunc, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedFiles | packages.NeedModule,
		Dir:  root,
	}

	pkgs, err := packages.Load(cfg, "./internal/...", "./cmd/api")
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	docs := make([]PackageDoc, 0, len(pkgs))
	for _, pkg := range pkgs {
		if len(pkg.Syntax) == 0 {
			continue
		}
		if _, ok := seen[pkg.PkgPath]; ok {
			continue
		}
		seen[pkg.PkgPath] = struct{}{}

		dpkg, err := doc.NewFromFiles(pkg.Fset, pkg.Syntax, pkg.PkgPath, doc.AllDecls)
		if err != nil {
			return nil, err
		}

		docs = append(docs, PackageDoc{
			ImportPath: pkg.PkgPath,
			Name:       dpkg.Name,
			Synopsis:   doc.Synopsis(dpkg.Doc),
			Doc:        strings.TrimSpace(dpkg.Doc),
		})
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].ImportPath < docs[j].ImportPath
	})

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"packages": docs})
	}, nil
}
