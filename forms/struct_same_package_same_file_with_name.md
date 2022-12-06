```go
type MyStruct struct {
    NameStructSamePackageSameFile StructSamePackageSameFile
}
```

```
     5  .  .  .  0: *ast.Field {
     6  .  .  .  .  Names: []*ast.Ident (len = 1) {
     7  .  .  .  .  .  0: *ast.Ident {
     8  .  .  .  .  .  .  NamePos: model/a3.go:4:2
     9  .  .  .  .  .  .  Name: "NameStructSamePackageSameFile"
    10  .  .  .  .  .  .  Obj: *ast.Object {
    11  .  .  .  .  .  .  .  Kind: var
    12  .  .  .  .  .  .  .  Name: "NameStructSamePackageSameFile"
    13  .  .  .  .  .  .  .  Decl: *(obj @ 5)
    14  .  .  .  .  .  .  }
    15  .  .  .  .  .  }
    16  .  .  .  .  }
    17  .  .  .  .  Type: *ast.Ident {
    18  .  .  .  .  .  NamePos: model/a3.go:4:32
    19  .  .  .  .  .  Name: "StructSamePackageSameFile"
    20  .  .  .  .  .  Obj: *ast.Object {
    21  .  .  .  .  .  .  Kind: type
    22  .  .  .  .  .  .  Name: "StructSamePackageSameFile"
    23  .  .  .  .  .  .  Decl: *ast.TypeSpec {
    24  .  .  .  .  .  .  .  Name: *ast.Ident {
    25  .  .  .  .  .  .  .  .  NamePos: model/a3.go:29:6
    26  .  .  .  .  .  .  .  .  Name: "StructSamePackageSameFile"
    27  .  .  .  .  .  .  .  .  Obj: *(obj @ 20)
    28  .  .  .  .  .  .  .  }
    29  .  .  .  .  .  .  .  Assign: -
    30  .  .  .  .  .  .  .  Type: *ast.StructType {
    31  .  .  .  .  .  .  .  .  Struct: model/a3.go:29:32
    32  .  .  .  .  .  .  .  .  Fields: *ast.FieldList {
    33  .  .  .  .  .  .  .  .  .  Opening: model/a3.go:29:39
    34  .  .  .  .  .  .  .  .  .  Closing: model/a3.go:30:1
    35  .  .  .  .  .  .  .  .  }
    36  .  .  .  .  .  .  .  .  Incomplete: false
    37  .  .  .  .  .  .  .  }
    38  .  .  .  .  .  .  }
    39  .  .  .  .  .  }
    40  .  .  .  .  }
    41  .  .  .  }

```