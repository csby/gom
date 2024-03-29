package main

import (
	"fmt"
	"github.com/csby/gwsf/gpkg"
	"path"
	"path/filepath"
)

const (
	srcPath = "src/github.com/csby/gom"
)

var (
	goIgnore = []string{
		".git",
		".idea",
		".gitignore",
		"go.sum",
	}
	vueIgnore = []string{
		"node_modules",
		"dist",
		".git",
		".idea",
		".gitignore",
		".editorconfig",
		"package-lock.json",
	}
)

type Pkg struct {
	binPath string
}

func (s *Pkg) Run() {
	// app
	fmt.Println("binary file path: ", s.binPath)
	binFolder := filepath.Dir(s.binPath)
	fmt.Println("binary folder path: ", binFolder)
	_, binName := filepath.Split(s.binPath)
	fmt.Println("binary file name: ", binName)
	binExt := path.Ext(binName)
	fmt.Println("binary file ext: ", binExt)
	appName := moduleName
	appFileName := fmt.Sprintf("%s%s", appName, binExt)
	fmt.Println("app file name: ", appFileName)

	// shell
	shellBinName := fmt.Sprintf("%s%s", "go_build_github_com_csby_gtool_gshell", binExt)
	fmt.Println("shell bin name: ", shellBinName)
	shellFileName := fmt.Sprintf("%s%s", "gshell", binExt)
	fmt.Println("shell file name: ", shellFileName)

	tmpFolder := filepath.Dir(binFolder)
	srcFolder := filepath.Join(filepath.Dir(tmpFolder), srcPath)
	fmt.Println("source folder path: ", srcFolder)

	relFolder := filepath.Join(tmpFolder, "rel", moduleName)
	fmt.Println("output folder path: ", relFolder)

	// site
	vueFolder := filepath.Join(filepath.Dir(filepath.Dir(tmpFolder)), "vue")
	fmt.Println("vue folder path: ", vueFolder)
	docFolder := filepath.Join(vueFolder, "gwsf-doc")
	fmt.Println("doc folder path: ", docFolder)
	optFolder := filepath.Join(vueFolder, "gom-opt")
	fmt.Println("opt folder path: ", optFolder)

	c := &gpkg.Config{
		Version:     moduleVersion,
		Destination: relFolder,
		Source:      true,
		Apps: []gpkg.Application{
			{
				Enable: true,
				Name:   appName,
				Bin: gpkg.Binary{
					Root: binFolder,
					Files: map[string]string{
						binName:      appFileName,
						shellBinName: shellFileName,
					},
				},
				Src: gpkg.Source{
					Enable: false,
					Root:   srcFolder,
					Ignore: goIgnore,
				},
				Webs: []gpkg.Website{
					{
						Enable: false,
						Name:   "doc",
						Src: gpkg.Source{
							Root:   docFolder,
							Ignore: vueIgnore,
						},
					},
					{
						Enable: true,
						Name:   "opt",
						Src: gpkg.Source{
							Root:   optFolder,
							Ignore: vueIgnore,
						},
					},
				},
			},
			{
				Enable: true,
				Name:   "gshell",
				Bin: gpkg.Binary{
					Root: binFolder,
					Files: map[string]string{
						shellBinName: shellFileName,
					},
				},
			},
		},
	}

	cfgFolder := filepath.Join(tmpFolder, "cfg")
	cfgName := fmt.Sprintf("%s-pkg.json", appName)
	cfgPath := filepath.Join(cfgFolder, cfgName)
	fmt.Println("config file path: ", cfgPath)
	c.LoadFromFile(cfgPath)

	p := gpkg.NewPacker(c)
	e := p.Pack()
	if e != nil {
		fmt.Println("打包失败: ", e)
	} else {
		fmt.Println("打包成功: ", p.OutputFolder())
		c.SaveToFile(cfgPath)
	}
}
