```go
type MyStruct struct {
    Name, Address string
}
```

```
    28  .  .  .  2: *ast.Field {
    29  .  .  .  .  Names: []*ast.Ident (len = 2) {
    30  .  .  .  .  .  0: *ast.Ident {
    31  .  .  .  .  .  .  NamePos: model/a3.go:6:2
    32  .  .  .  .  .  .  Name: "Address"
    33  .  .  .  .  .  .  Obj: *ast.Object {
    34  .  .  .  .  .  .  .  Kind: var
    35  .  .  .  .  .  .  .  Name: "Address"
    36  .  .  .  .  .  .  .  Decl: *(obj @ 28)
    37  .  .  .  .  .  .  }
    38  .  .  .  .  .  }
    39  .  .  .  .  .  1: *ast.Ident {
    40  .  .  .  .  .  .  NamePos: model/a3.go:6:11
    41  .  .  .  .  .  .  Name: "Phone"
    42  .  .  .  .  .  .  Obj: *ast.Object {
    43  .  .  .  .  .  .  .  Kind: var
    44  .  .  .  .  .  .  .  Name: "Phone"
    45  .  .  .  .  .  .  .  Decl: *(obj @ 28)
    46  .  .  .  .  .  .  }
    47  .  .  .  .  .  }
    48  .  .  .  .  }
    49  .  .  .  .  Type: *ast.Ident {
    50  .  .  .  .  .  NamePos: model/a3.go:6:17
    51  .  .  .  .  .  Name: "string"
    52  .  .  .  .  }
    53  .  .  .  }
```