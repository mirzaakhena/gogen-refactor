```go
type MyStruct struct {
    MultiDirectional chan int
    SendDirectional  chan<- int
    RecvDirectional  <-chan int
}
```


```
     5  .  .  .  0: *ast.Field {
     6  .  .  .  .  Names: []*ast.Ident (len = 1) {
     7  .  .  .  .  .  0: *ast.Ident {
     8  .  .  .  .  .  .  NamePos: model/a3.go:4:2
     9  .  .  .  .  .  .  Name: "MultiDirectional"
    10  .  .  .  .  .  .  Obj: *ast.Object {
    11  .  .  .  .  .  .  .  Kind: var
    12  .  .  .  .  .  .  .  Name: "MultiDirectional"
    13  .  .  .  .  .  .  .  Decl: *(obj @ 5)
    14  .  .  .  .  .  .  }
    15  .  .  .  .  .  }
    16  .  .  .  .  }
    17  .  .  .  .  Type: *ast.ChanType {
    18  .  .  .  .  .  Begin: model/a3.go:4:19
    19  .  .  .  .  .  Arrow: -
    20  .  .  .  .  .  Dir: 3
    21  .  .  .  .  .  Value: *ast.Ident {
    22  .  .  .  .  .  .  NamePos: model/a3.go:4:24
    23  .  .  .  .  .  .  Name: "int"
    24  .  .  .  .  .  }
    25  .  .  .  .  }
    26  .  .  .  }
    27  .  .  .  1: *ast.Field {
    28  .  .  .  .  Names: []*ast.Ident (len = 1) {
    29  .  .  .  .  .  0: *ast.Ident {
    30  .  .  .  .  .  .  NamePos: model/a3.go:5:2
    31  .  .  .  .  .  .  Name: "SendDirectional"
    32  .  .  .  .  .  .  Obj: *ast.Object {
    33  .  .  .  .  .  .  .  Kind: var
    34  .  .  .  .  .  .  .  Name: "SendDirectional"
    35  .  .  .  .  .  .  .  Decl: *(obj @ 27)
    36  .  .  .  .  .  .  }
    37  .  .  .  .  .  }
    38  .  .  .  .  }
    39  .  .  .  .  Type: *ast.ChanType {
    40  .  .  .  .  .  Begin: model/a3.go:5:19
    41  .  .  .  .  .  Arrow: model/a3.go:5:23
    42  .  .  .  .  .  Dir: 1
    43  .  .  .  .  .  Value: *ast.Ident {
    44  .  .  .  .  .  .  NamePos: model/a3.go:5:26
    45  .  .  .  .  .  .  Name: "int"
    46  .  .  .  .  .  }
    47  .  .  .  .  }
    48  .  .  .  }
    49  .  .  .  2: *ast.Field {
    50  .  .  .  .  Names: []*ast.Ident (len = 1) {
    51  .  .  .  .  .  0: *ast.Ident {
    52  .  .  .  .  .  .  NamePos: model/a3.go:6:2
    53  .  .  .  .  .  .  Name: "RecvDirectional"
    54  .  .  .  .  .  .  Obj: *ast.Object {
    55  .  .  .  .  .  .  .  Kind: var
    56  .  .  .  .  .  .  .  Name: "RecvDirectional"
    57  .  .  .  .  .  .  .  Decl: *(obj @ 49)
    58  .  .  .  .  .  .  }
    59  .  .  .  .  .  }
    60  .  .  .  .  }
    61  .  .  .  .  Type: *ast.ChanType {
    62  .  .  .  .  .  Begin: model/a3.go:6:19
    63  .  .  .  .  .  Arrow: model/a3.go:6:19
    64  .  .  .  .  .  Dir: 2
    65  .  .  .  .  .  Value: *ast.Ident {
    66  .  .  .  .  .  .  NamePos: model/a3.go:6:26
    67  .  .  .  .  .  .  Name: "int"
    68  .  .  .  .  .  }
    69  .  .  .  .  }
    70  .  .  .  }
```