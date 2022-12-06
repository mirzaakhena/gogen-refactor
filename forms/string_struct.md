
```go
type MyStruct struct {
    Name [3]string
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
    19  .  .  .  .  .  Len: *ast.BasicLit {
    20  .  .  .  .  .  .  ValuePos: model/a3.go:4:8
    21  .  .  .  .  .  .  Kind: INT
    22  .  .  .  .  .  .  Value: "3"
    23  .  .  .  .  .  }
    24  .  .  .  .  .  Elt: *ast.Ident {
    25  .  .  .  .  .  .  NamePos: model/a3.go:4:10
    26  .  .  .  .  .  .  Name: "string"
    27  .  .  .  .  .  }
    28  .  .  .  .  }
    29  .  .  .  }
```