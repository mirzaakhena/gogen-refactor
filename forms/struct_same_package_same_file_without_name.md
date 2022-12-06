```go
type MyStruct struct {
    StructSamePackageSameFile
}
```

```
    11  .  .  .  1: *ast.Field {
    12  .  .  .  .  Type: *ast.Ident {
    13  .  .  .  .  .  NamePos: model/a3.go:5:2
    14  .  .  .  .  .  Name: "StructSamePackageSameFile"
    15  .  .  .  .  .  Obj: *ast.Object {
    16  .  .  .  .  .  .  Kind: type
    17  .  .  .  .  .  .  Name: "StructSamePackageSameFile"
    18  .  .  .  .  .  .  Decl: *ast.TypeSpec {
    19  .  .  .  .  .  .  .  Name: *ast.Ident {
    20  .  .  .  .  .  .  .  .  NamePos: model/a3.go:32:6
    21  .  .  .  .  .  .  .  .  Name: "StructSamePackageSameFile"
    22  .  .  .  .  .  .  .  .  Obj: *(obj @ 15)
    23  .  .  .  .  .  .  .  }
    24  .  .  .  .  .  .  .  Assign: -
    25  .  .  .  .  .  .  .  Type: *ast.StructType {
    26  .  .  .  .  .  .  .  .  Struct: model/a3.go:32:32
    27  .  .  .  .  .  .  .  .  Fields: *ast.FieldList {
    28  .  .  .  .  .  .  .  .  .  Opening: model/a3.go:32:39
    29  .  .  .  .  .  .  .  .  .  Closing: model/a3.go:33:1
    30  .  .  .  .  .  .  .  .  }
    31  .  .  .  .  .  .  .  .  Incomplete: false
    32  .  .  .  .  .  .  .  }
    33  .  .  .  .  .  .  }
    34  .  .  .  .  .  }
    35  .  .  .  .  }
    36  .  .  .  }
```