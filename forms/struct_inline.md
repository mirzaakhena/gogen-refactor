```go
type MyStruct struct {
    TheStruct struct {
        Name string
        Age  int
    }
}
```

```
     5  .  .  .  0: *ast.Field {
     6  .  .  .  .  Names: []*ast.Ident (len = 1) {
     7  .  .  .  .  .  0: *ast.Ident {
     8  .  .  .  .  .  .  NamePos: model/a3.go:4:2
     9  .  .  .  .  .  .  Name: "TheStruct"
    10  .  .  .  .  .  .  Obj: *ast.Object {
    11  .  .  .  .  .  .  .  Kind: var
    12  .  .  .  .  .  .  .  Name: "TheStruct"
    13  .  .  .  .  .  .  .  Decl: *(obj @ 5)
    14  .  .  .  .  .  .  }
    15  .  .  .  .  .  }
    16  .  .  .  .  }
    17  .  .  .  .  Type: *ast.StructType {
    18  .  .  .  .  .  Struct: model/a3.go:4:12
    19  .  .  .  .  .  Fields: *ast.FieldList {
    20  .  .  .  .  .  .  Opening: model/a3.go:4:19
    21  .  .  .  .  .  .  List: []*ast.Field (len = 2) {
    22  .  .  .  .  .  .  .  0: *ast.Field {
    23  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    24  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    25  .  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:5:3
    26  .  .  .  .  .  .  .  .  .  .  Name: "Name"
    27  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    28  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    29  .  .  .  .  .  .  .  .  .  .  .  Name: "Name"
    30  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 22)
    31  .  .  .  .  .  .  .  .  .  .  }
    32  .  .  .  .  .  .  .  .  .  }
    33  .  .  .  .  .  .  .  .  }
    34  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    35  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:5:8
    36  .  .  .  .  .  .  .  .  .  Name: "string"
    37  .  .  .  .  .  .  .  .  }
    38  .  .  .  .  .  .  .  }
    39  .  .  .  .  .  .  .  1: *ast.Field {
    40  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    41  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    42  .  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:6:3
    43  .  .  .  .  .  .  .  .  .  .  Name: "Age"
    44  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    45  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    46  .  .  .  .  .  .  .  .  .  .  .  Name: "Age"
    47  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 39)
    48  .  .  .  .  .  .  .  .  .  .  }
    49  .  .  .  .  .  .  .  .  .  }
    50  .  .  .  .  .  .  .  .  }
    51  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    52  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:6:8
    53  .  .  .  .  .  .  .  .  .  Name: "int"
    54  .  .  .  .  .  .  .  .  }
    55  .  .  .  .  .  .  .  }
    56  .  .  .  .  .  .  }
    57  .  .  .  .  .  .  Closing: model/a3.go:7:2
    58  .  .  .  .  .  }
    59  .  .  .  .  .  Incomplete: false
    60  .  .  .  .  }
    61  .  .  .  }
```