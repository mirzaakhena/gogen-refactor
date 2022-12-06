```go
type MyStruct struct {
    Name []*string
}
```

```
     5  .  .  .  0: *ast.Field {
     6  .  .  .  .  Names: []*ast.Ident (len = 1) {
     7  .  .  .  .  .  0: *ast.Ident {
     8  .  .  .  .  .  .  NamePos: model/a3.go:4:2
     9  .  .  .  .  .  .  Name: "Name"
    10  .  .  .  .  .  .  Obj: *ast.Object {
    11  .  .  .  .  .  .  .  Kind: var
    12  .  .  .  .  .  .  .  Name: "Name"
    13  .  .  .  .  .  .  .  Decl: *(obj @ 5)
    14  .  .  .  .  .  .  }
    15  .  .  .  .  .  }
    16  .  .  .  .  }
    17  .  .  .  .  Type: *ast.ArrayType {
    18  .  .  .  .  .  Lbrack: model/a3.go:4:7
    19  .  .  .  .  .  Elt: *ast.StarExpr {
    20  .  .  .  .  .  .  Star: model/a3.go:4:9
    21  .  .  .  .  .  .  X: *ast.Ident {
    22  .  .  .  .  .  .  .  NamePos: model/a3.go:4:10
    23  .  .  .  .  .  .  .  Name: "string"
    24  .  .  .  .  .  .  }
    25  .  .  .  .  .  }
    26  .  .  .  .  }
    27  .  .  .  }
```