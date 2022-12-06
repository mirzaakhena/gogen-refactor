```go
type MyStruct struct {
    Name string
}
```

```
    11  .  .  .  1: *ast.Field {
    12  .  .  .  .  Names: []*ast.Ident (len = 1) {
    13  .  .  .  .  .  0: *ast.Ident {
    14  .  .  .  .  .  .  NamePos: model/a3.go:5:2
    15  .  .  .  .  .  .  Name: "Name"
    16  .  .  .  .  .  .  Obj: *ast.Object {
    17  .  .  .  .  .  .  .  Kind: var
    18  .  .  .  .  .  .  .  Name: "Name"
    19  .  .  .  .  .  .  .  Decl: *(obj @ 11)
    20  .  .  .  .  .  .  }
    21  .  .  .  .  .  }
    22  .  .  .  .  }
    23  .  .  .  .  Type: *ast.Ident {
    24  .  .  .  .  .  NamePos: model/a3.go:5:17
    25  .  .  .  .  .  Name: "string"
    26  .  .  .  .  }
    27  .  .  .  }
```