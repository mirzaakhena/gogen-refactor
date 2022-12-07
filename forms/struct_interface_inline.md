```go
type MyStruct struct {
    Name interface{}
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
    17  .  .  .  .  Type: *ast.InterfaceType {
    18  .  .  .  .  .  Interface: model/a3.go:4:7
    19  .  .  .  .  .  Methods: *ast.FieldList {
    20  .  .  .  .  .  .  Opening: model/a3.go:4:16
    21  .  .  .  .  .  .  Closing: model/a3.go:4:17
    22  .  .  .  .  .  }
    23  .  .  .  .  .  Incomplete: false
    24  .  .  .  .  }
    25  .  .  .  }
```