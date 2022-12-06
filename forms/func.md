```go
type MyStruct struct {
    Name func (a int, b string) (float64, error)
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
    17  .  .  .  .  Type: *ast.FuncType {
    18  .  .  .  .  .  Func: model/a3.go:4:7
    19  .  .  .  .  .  Params: *ast.FieldList {
    20  .  .  .  .  .  .  Opening: model/a3.go:4:11
    21  .  .  .  .  .  .  List: []*ast.Field (len = 2) {
    22  .  .  .  .  .  .  .  0: *ast.Field {
    23  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    24  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    25  .  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:4:12
    26  .  .  .  .  .  .  .  .  .  .  Name: "a"
    27  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    28  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    29  .  .  .  .  .  .  .  .  .  .  .  Name: "a"
    30  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 22)
    31  .  .  .  .  .  .  .  .  .  .  }
    32  .  .  .  .  .  .  .  .  .  }
    33  .  .  .  .  .  .  .  .  }
    34  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    35  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:4:14
    36  .  .  .  .  .  .  .  .  .  Name: "int"
    37  .  .  .  .  .  .  .  .  }
    38  .  .  .  .  .  .  .  }
    39  .  .  .  .  .  .  .  1: *ast.Field {
    40  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    41  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    42  .  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:4:19
    43  .  .  .  .  .  .  .  .  .  .  Name: "b"
    44  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    45  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    46  .  .  .  .  .  .  .  .  .  .  .  Name: "b"
    47  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 39)
    48  .  .  .  .  .  .  .  .  .  .  }
    49  .  .  .  .  .  .  .  .  .  }
    50  .  .  .  .  .  .  .  .  }
    51  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    52  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:4:21
    53  .  .  .  .  .  .  .  .  .  Name: "string"
    54  .  .  .  .  .  .  .  .  }
    55  .  .  .  .  .  .  .  }
    56  .  .  .  .  .  .  }
    57  .  .  .  .  .  .  Closing: model/a3.go:4:27
    58  .  .  .  .  .  }
    59  .  .  .  .  .  Results: *ast.FieldList {
    60  .  .  .  .  .  .  Opening: model/a3.go:4:29
    61  .  .  .  .  .  .  List: []*ast.Field (len = 2) {
    62  .  .  .  .  .  .  .  0: *ast.Field {
    63  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    64  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:4:30
    65  .  .  .  .  .  .  .  .  .  Name: "float64"
    66  .  .  .  .  .  .  .  .  }
    67  .  .  .  .  .  .  .  }
    68  .  .  .  .  .  .  .  1: *ast.Field {
    69  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    70  .  .  .  .  .  .  .  .  .  NamePos: model/a3.go:4:39
    71  .  .  .  .  .  .  .  .  .  Name: "error"
    72  .  .  .  .  .  .  .  .  }
    73  .  .  .  .  .  .  .  }
    74  .  .  .  .  .  .  }
    75  .  .  .  .  .  .  Closing: model/a3.go:4:44
    76  .  .  .  .  .  }
    77  .  .  .  .  }
    78  .  .  .  }

```