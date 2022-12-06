```go
type MyStruct struct {
    VarExternalStruct service.ExternalStruct
}
```

```
    17  .  .  .  1: *ast.Field {
    18  .  .  .  .  Names: []*ast.Ident (len = 1) {
    19  .  .  .  .  .  0: *ast.Ident {
    20  .  .  .  .  .  .  NamePos: model/a3.go:13:2
    21  .  .  .  .  .  .  Name: "VarExternalStruct"
    22  .  .  .  .  .  .  Obj: *ast.Object {
    23  .  .  .  .  .  .  .  Kind: var
    24  .  .  .  .  .  .  .  Name: "VarExternalStruct"
    25  .  .  .  .  .  .  .  Decl: *(obj @ 17)
    26  .  .  .  .  .  .  }
    27  .  .  .  .  .  }
    28  .  .  .  .  }
    29  .  .  .  .  Type: *ast.SelectorExpr {
    30  .  .  .  .  .  X: *ast.Ident {
    31  .  .  .  .  .  .  NamePos: model/a3.go:13:20
    32  .  .  .  .  .  .  Name: "service"
    33  .  .  .  .  .  }
    34  .  .  .  .  .  Sel: *ast.Ident {
    35  .  .  .  .  .  .  NamePos: model/a3.go:13:28
    36  .  .  .  .  .  .  Name: "ExternalStruct"
    37  .  .  .  .  .  }
    38  .  .  .  .  }
    39  .  .  .  }

```