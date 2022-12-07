```go
type TargetInterface1 interface {
	func(a int, b float64) (string, error)
	InterfaceInSameFile
	InterfaceInDiffFile
	p2.InterfaceInDiffPackage
}
```


```
     0  *ast.TypeSpec {
     1  .  Name: *ast.Ident {
     2  .  .  NamePos: p1/file1.go:55:6
     3  .  .  Name: "TargetInterface1"
     4  .  .  Obj: *ast.Object {
     5  .  .  .  Kind: type
     6  .  .  .  Name: "TargetInterface1"
     7  .  .  .  Decl: *(obj @ 0)
     8  .  .  }
     9  .  }
    10  .  Assign: -
    11  .  Type: *ast.InterfaceType {
    12  .  .  Interface: p1/file1.go:55:23
    13  .  .  Methods: *ast.FieldList {
    14  .  .  .  Opening: p1/file1.go:55:33
    15  .  .  .  List: []*ast.Field (len = 4) {
    16  .  .  .  .  0: *ast.Field {
    17  .  .  .  .  .  Type: *ast.FuncType {
    18  .  .  .  .  .  .  Func: p1/file1.go:56:2
    19  .  .  .  .  .  .  Params: *ast.FieldList {
    20  .  .  .  .  .  .  .  Opening: p1/file1.go:56:6
    21  .  .  .  .  .  .  .  List: []*ast.Field (len = 2) {
    22  .  .  .  .  .  .  .  .  0: *ast.Field {
    23  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    24  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    25  .  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:56:7
    26  .  .  .  .  .  .  .  .  .  .  .  Name: "a"
    27  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    28  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    29  .  .  .  .  .  .  .  .  .  .  .  .  Name: "a"
    30  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 22)
    31  .  .  .  .  .  .  .  .  .  .  .  }
    32  .  .  .  .  .  .  .  .  .  .  }
    33  .  .  .  .  .  .  .  .  .  }
    34  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    35  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:56:9
    36  .  .  .  .  .  .  .  .  .  .  Name: "int"
    37  .  .  .  .  .  .  .  .  .  }
    38  .  .  .  .  .  .  .  .  }
    39  .  .  .  .  .  .  .  .  1: *ast.Field {
    40  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    41  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    42  .  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:56:14
    43  .  .  .  .  .  .  .  .  .  .  .  Name: "b"
    44  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    45  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    46  .  .  .  .  .  .  .  .  .  .  .  .  Name: "b"
    47  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 39)
    48  .  .  .  .  .  .  .  .  .  .  .  }
    49  .  .  .  .  .  .  .  .  .  .  }
    50  .  .  .  .  .  .  .  .  .  }
    51  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    52  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:56:16
    53  .  .  .  .  .  .  .  .  .  .  Name: "float64"
    54  .  .  .  .  .  .  .  .  .  }
    55  .  .  .  .  .  .  .  .  }
    56  .  .  .  .  .  .  .  }
    57  .  .  .  .  .  .  .  Closing: p1/file1.go:56:23
    58  .  .  .  .  .  .  }
    59  .  .  .  .  .  .  Results: *ast.FieldList {
    60  .  .  .  .  .  .  .  Opening: p1/file1.go:56:25
    61  .  .  .  .  .  .  .  List: []*ast.Field (len = 2) {
    62  .  .  .  .  .  .  .  .  0: *ast.Field {
    63  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    64  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:56:26
    65  .  .  .  .  .  .  .  .  .  .  Name: "string"
    66  .  .  .  .  .  .  .  .  .  }
    67  .  .  .  .  .  .  .  .  }
    68  .  .  .  .  .  .  .  .  1: *ast.Field {
    69  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    70  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:56:34
    71  .  .  .  .  .  .  .  .  .  .  Name: "error"
    72  .  .  .  .  .  .  .  .  .  }
    73  .  .  .  .  .  .  .  .  }
    74  .  .  .  .  .  .  .  }
    75  .  .  .  .  .  .  .  Closing: p1/file1.go:56:39
    76  .  .  .  .  .  .  }
    77  .  .  .  .  .  }
    78  .  .  .  .  }
    79  .  .  .  .  1: *ast.Field {
    80  .  .  .  .  .  Type: *ast.Ident {
    81  .  .  .  .  .  .  NamePos: p1/file1.go:57:2
    82  .  .  .  .  .  .  Name: "InterfaceInSameFile"
    83  .  .  .  .  .  .  Obj: *ast.Object {
    84  .  .  .  .  .  .  .  Kind: type
    85  .  .  .  .  .  .  .  Name: "InterfaceInSameFile"
    86  .  .  .  .  .  .  .  Decl: *ast.TypeSpec {
    87  .  .  .  .  .  .  .  .  Name: *ast.Ident {
    88  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:62:6
    89  .  .  .  .  .  .  .  .  .  Name: "InterfaceInSameFile"
    90  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 83)
    91  .  .  .  .  .  .  .  .  }
    92  .  .  .  .  .  .  .  .  Assign: -
    93  .  .  .  .  .  .  .  .  Type: *ast.InterfaceType {
    94  .  .  .  .  .  .  .  .  .  Interface: p1/file1.go:62:26
    95  .  .  .  .  .  .  .  .  .  Methods: *ast.FieldList {
    96  .  .  .  .  .  .  .  .  .  .  Opening: p1/file1.go:62:36
    97  .  .  .  .  .  .  .  .  .  .  List: []*ast.Field (len = 1) {
    98  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Field {
    99  .  .  .  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
   100  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
   101  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:63:2
   102  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "MethodInSameFile"
   103  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
   104  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Kind: func
   105  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "MethodInSameFile"
   106  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 98)
   107  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   108  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   109  .  .  .  .  .  .  .  .  .  .  .  .  }
   110  .  .  .  .  .  .  .  .  .  .  .  .  Type: *ast.FuncType {
   111  .  .  .  .  .  .  .  .  .  .  .  .  .  Func: -
   112  .  .  .  .  .  .  .  .  .  .  .  .  .  Params: *ast.FieldList {
   113  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Opening: p1/file1.go:63:18
   114  .  .  .  .  .  .  .  .  .  .  .  .  .  .  List: []*ast.Field (len = 1) {
   115  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Field {
   116  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
   117  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
   118  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:63:19
   119  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "ctx"
   120  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
   121  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
   122  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "ctx"
   123  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 115)
   124  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   125  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   126  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   127  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Type: *ast.SelectorExpr {
   128  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
   129  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:63:23
   130  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "context"
   131  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   132  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
   133  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:63:31
   134  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Context"
   135  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   136  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   137  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   138  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   139  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Closing: p1/file1.go:63:38
   140  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   141  .  .  .  .  .  .  .  .  .  .  .  .  .  Results: *ast.FieldList {
   142  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Opening: -
   143  .  .  .  .  .  .  .  .  .  .  .  .  .  .  List: []*ast.Field (len = 1) {
   144  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  0: *ast.Field {
   145  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
   146  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  NamePos: p1/file1.go:63:40
   147  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Name: "int"
   148  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   149  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   150  .  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   151  .  .  .  .  .  .  .  .  .  .  .  .  .  .  Closing: -
   152  .  .  .  .  .  .  .  .  .  .  .  .  .  }
   153  .  .  .  .  .  .  .  .  .  .  .  .  }
   154  .  .  .  .  .  .  .  .  .  .  .  }
   155  .  .  .  .  .  .  .  .  .  .  }
   156  .  .  .  .  .  .  .  .  .  .  Closing: p1/file1.go:64:1
   157  .  .  .  .  .  .  .  .  .  }
   158  .  .  .  .  .  .  .  .  .  Incomplete: false
   159  .  .  .  .  .  .  .  .  }
   160  .  .  .  .  .  .  .  }
   161  .  .  .  .  .  .  }
   162  .  .  .  .  .  }
   163  .  .  .  .  }
   164  .  .  .  .  2: *ast.Field {
   165  .  .  .  .  .  Type: *ast.Ident {
   166  .  .  .  .  .  .  NamePos: p1/file1.go:58:2
   167  .  .  .  .  .  .  Name: "InterfaceInDiffFile"
   168  .  .  .  .  .  }
   169  .  .  .  .  }
   170  .  .  .  .  3: *ast.Field {
   171  .  .  .  .  .  Type: *ast.SelectorExpr {
   172  .  .  .  .  .  .  X: *ast.Ident {
   173  .  .  .  .  .  .  .  NamePos: p1/file1.go:59:2
   174  .  .  .  .  .  .  .  Name: "p2"
   175  .  .  .  .  .  .  }
   176  .  .  .  .  .  .  Sel: *ast.Ident {
   177  .  .  .  .  .  .  .  NamePos: p1/file1.go:59:5
   178  .  .  .  .  .  .  .  Name: "InterfaceInDiffPackage"
   179  .  .  .  .  .  .  }
   180  .  .  .  .  .  }
   181  .  .  .  .  }
   182  .  .  .  }
   183  .  .  .  Closing: p1/file1.go:60:1
   184  .  .  }
   185  .  .  Incomplete: false
   186  .  }
   187  }
```