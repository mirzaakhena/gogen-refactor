```go
type MyStruct struct {
    StructSamePackageDiffFile
}
```

```
     5  .  .  .  0: *ast.Field {
     6  .  .  .  .  Type: *ast.Ident {
     7  .  .  .  .  .  NamePos: model/a3.go:4:2
     8  .  .  .  .  .  Name: "StructSamePackageDiffFile"
     9  .  .  .  .  }
    10  .  .  .  }
```