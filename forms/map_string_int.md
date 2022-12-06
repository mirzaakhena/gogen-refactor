```go
type MyStruct struct {
    Name map[string]int
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
    17  .  .  .  .  Type: *ast.MapType {
    18  .  .  .  .  .  Map: model/a3.go:4:7
    19  .  .  .  .  .  Key: *ast.Ident {
    20  .  .  .  .  .  .  NamePos: model/a3.go:4:11
    21  .  .  .  .  .  .  Name: "string"
    22  .  .  .  .  .  }
    23  .  .  .  .  .  Value: *ast.Ident {
    24  .  .  .  .  .  .  NamePos: model/a3.go:4:18
    25  .  .  .  .  .  .  Name: "int"
    26  .  .  .  .  .  }
    27  .  .  .  .  }
    28  .  .  .  }
```