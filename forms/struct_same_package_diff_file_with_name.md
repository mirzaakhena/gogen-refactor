```go
type MyStruct struct {
    NameStructSamePackageDiffFile StructSamePackageDiffFile
}
```

```
    37  .  .  .  2: *ast.Field {
    38  .  .  .  .  Names: []*ast.Ident (len = 1) {
    39  .  .  .  .  .  0: *ast.Ident {
    40  .  .  .  .  .  .  NamePos: model/a3.go:6:2
    41  .  .  .  .  .  .  Name: "NameStructSamePackageDiffFile"
    42  .  .  .  .  .  .  Obj: *ast.Object {
    43  .  .  .  .  .  .  .  Kind: var
    44  .  .  .  .  .  .  .  Name: "NameStructSamePackageDiffFile"
    45  .  .  .  .  .  .  .  Decl: *(obj @ 37)
    46  .  .  .  .  .  .  }
    47  .  .  .  .  .  }
    48  .  .  .  .  }
    49  .  .  .  .  Type: *ast.Ident {
    50  .  .  .  .  .  NamePos: model/a3.go:6:32
    51  .  .  .  .  .  Name: "StructSamePackageDiffFile"
    52  .  .  .  .  }
    53  .  .  .  }
```