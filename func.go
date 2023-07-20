package astX

import (
    "errors"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "io/fs"
    "path/filepath"
    "strings"

    "github.com/xzf/strX"
)

type FuncInfo struct {
    FilePath    string
    FileName    string
    FuncName    string
    PkgName     string
    PkgPath     string
    FuncComment []string //方法注释
}

type WalkPackageFuncReq struct {
    RootPkgPath     string
    CallBack        func(FuncInfo)
    SkipGoTestFile  bool
    SkipPrivateFunc bool
}

func WalkPackageFunc(req WalkPackageFuncReq) error {
    if req.RootPkgPath == "" {
        return errors.New(`[ez8jeqyx58] req.RootPkgPath == ""`)
    }
    if req.CallBack == nil {
        return errors.New(`[2skb2i5dr7] req.CallBack == nil`)
    }
    var getAtLeastOneFunc bool
    filepath.Walk(req.RootPkgPath, func(filePath string, fileInfo fs.FileInfo, fileErr error) error {
        //skip os.Stat return err
        if fileErr != nil {
            fmt.Println("[fzvlhvyp86] path ["+filePath+"] read err :", fileErr)
            return nil
        }
        //skip dir
        if fileInfo.IsDir() {
            return nil
        }
        fileName := fileInfo.Name()
        if strings.HasSuffix(fileName, ".go") == false {
            return nil
        }
        if req.SkipGoTestFile {
            if strings.HasSuffix(fileName, "_test.go") {
                return nil
            }
        }
        astFileSet := token.NewFileSet()
        astInfo, err := parser.ParseFile(astFileSet, filePath, nil, parser.ParseComments)
        if err != nil {
            fmt.Println("[82xisnz235] path", "["+filePath+"]", "code parse failed:", err)
            return nil
        }
        for _, item := range astInfo.Decls {
            switch item.(type) {
            case *ast.FuncDecl:
                funcItem := item.(*ast.FuncDecl)
                funcName := funcItem.Name.String()
                if funcName == "" {
                    break
                }
                if req.SkipPrivateFunc {
                    start := funcName[0]
                    if ('A' <= start && start <= 'Z') == false {
                        continue
                    }
                }
                getAtLeastOneFunc = true
                var commentSlice []string
                if funcItem.Doc != nil {
                    for _, commentObj := range funcItem.Doc.List {
                        if commentObj == nil {
                            continue
                        }
                        commentSlice = append(commentSlice, commentObj.Text)
                    }
                }
                pkgPath := strX.SubBeforeLast(filePath, "/")
                pkgPath = strings.TrimLeft(pkgPath, "./")
                req.CallBack(FuncInfo{
                    FilePath:    filePath,
                    FileName:    fileName,
                    FuncName:    funcName,
                    PkgName:     astInfo.Name.Name,
                    PkgPath:     pkgPath,
                    FuncComment: commentSlice,
                })
            }
        }
        return nil
    })
    if getAtLeastOneFunc == false {
        return errors.New("[fizgffgv4s] not func match found")
    }
    return nil
}
