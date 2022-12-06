```go
type MyStruct struct {
    TheStruct *struct {
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
    17  .  .  .  .  Type: *ast.StarExpr {
    18  .  .  .  .  .  Star: model/a3.go:4:12
    19  .  .  .  .  .  X: *ast.StructType {
    20  .  .  .  .  .  .  Struct: model/a3.go:4:13
    21  .  .  .  .  .  .  Fields: *ast.FieldList {
    22  .  .  .  .  .  .  .  Opening: model/a3.go:4:20
    23  .  .  .  .  .  .  .  List: []*ast.Field (len = 2) {
    24  .  .  .  .  .  .  .  .  0: *ast.Field {
    25  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    26  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    27  .  .  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:5:3
    28  .  .  .  .  .  .  .  .  .  .  .  Name: "Name"
    29  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    30  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    31  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Name"
    32  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 24)
    33  .  .  .  .  .  .  .  .  .  .  .  }
    34  .  .  .  .  .  .  .  .  .  .  }
    35  .  .  .  .  .  .  .  .  .  }
    36  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    37  .  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:5:8
    38  .  .  .  .  .  .  .  .  .  .  Name: "string"
    39  .  .  .  .  .  .  .  .  .  }
    40  .  .  .  .  .  .  .  .  }
    41  .  .  .  .  .  .  .  .  1: *ast.Field {
    42  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    43  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    44  .  .  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:6:3
    45  .  .  .  .  .  .  .  .  .  .  .  Name: "Age"
    46  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    47  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    48  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Age"
    49  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 41)
    50  .  .  .  .  .  .  .  .  .  .  .  }
    51  .  .  .  .  .  .  .  .  .  .  }
    52  .  .  .  .  .  .  .  .  .  }
    53  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    54  .  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:6:8
    55  .  .  .  .  .  .  .  .  .  .  Name: "int"
    56  .  .  .  .  .  .  .  .  .  }
    57  .  .  .  .  .  .  .  .  }
    58  .  .  .  .  .  .  .  }
    59  .  .  .  .  .  .  .  Closing: model/a3.go:7:2
    60  .  .  .  .  .  .  }
    61  .  .  .  .  .  .  Incomplete: false
    62  .  .  .  .  .  }
    63  .  .  .  .  }
    64  .  .  .  }
```