
## checking whether the existing struct already implements the method

### Empty existing struct gateway

```
type gateway struct {
}

0  *ast.GenDecl {
1  .  TokPos: gateway/prod/gateway.go:10:1
2  .  Tok: type
3  .  Lparen: -
4  .  Specs: []ast.Spec (len = 1) {
5  .  .  0: *ast.TypeSpec {
6  .  .  .  Name: *ast.Ident {
7  .  .  .  .  NamePos: gateway/prod/gateway.go:10:6
8  .  .  .  .  Name: "gateway"
9  .  .  .  .  Obj: *ast.Object {
0  .  .  .  .  .  Kind: type
1  .  .  .  .  .  Name: "gateway"
2  .  .  .  .  .  Decl: *(obj @ 5)
3  .  .  .  .  }
4  .  .  .  }
5  .  .  .  Assign: -
6  .  .  .  Type: *ast.StructType {
7  .  .  .  .  Struct: gateway/prod/gateway.go:10:14
8  .  .  .  .  Fields: *ast.FieldList {
9  .  .  .  .  .  Opening: gateway/prod/gateway.go:10:21
0  .  .  .  .  .  Closing: gateway/prod/gateway.go:11:1
1  .  .  .  .  }
2  .  .  .  .  Incomplete: false
3  .  .  .  }
4  .  .  }
5  .  }
6  .  Rparen: -
7  }
```

### simple composition struct gateway 
```
type gateway struct {
	payment
}

0  *ast.GenDecl {
1  .  TokPos: gateway/prod/gateway.go:10:1
2  .  Tok: type
3  .  Lparen: -
4  .  Specs: []ast.Spec (len = 1) {
5  .  .  0: *ast.TypeSpec {
6  .  .  .  Name: *ast.Ident {
7  .  .  .  .  NamePos: gateway/prod/gateway.go:10:6
8  .  .  .  .  Name: "gateway"
9  .  .  .  .  Obj: *ast.Object {
0  .  .  .  .  .  Kind: type
1  .  .  .  .  .  Name: "gateway"
2  .  .  .  .  .  Decl: *(obj @ 5)
3  .  .  .  .  }
4  .  .  .  }
5  .  .  .  Assign: -
6  .  .  .  Type: *ast.StructType {
7  .  .  .  .  Struct: gateway/prod/gateway.go:10:14
8  .  .  .  .  Fields: *ast.FieldList {
9  .  .  .  .  .  Opening: gateway/prod/gateway.go:10:21
0  .  .  .  .  .  List: []*ast.Field (len = 1) {
1  .  .  .  .  .  .  0: *ast.Field {
2  .  .  .  .  .  .  .  Type: *ast.Ident {
3  .  .  .  .  .  .  .  .  NamePos: gateway/prod/gateway.go:11:2
4  .  .  .  .  .  .  .  .  Name: "payment"
5  .  .  .  .  .  .  .  }
6  .  .  .  .  .  .  }
7  .  .  .  .  .  }
8  .  .  .  .  .  Closing: gateway/prod/gateway.go:12:1
9  .  .  .  .  }
0  .  .  .  .  Incomplete: false
1  .  .  .  }
2  .  .  }
3  .  }
4  .  Rparen: -
5  }
```

### pointer composition struct gateway
```
type gateway struct {
	*payment
}

 0  *ast.GenDecl {
 1  .  TokPos: gateway/prod/gateway.go:10:1
 2  .  Tok: type
 3  .  Lparen: -
 4  .  Specs: []ast.Spec (len = 1) {
 5  .  .  0: *ast.TypeSpec {
 6  .  .  .  Name: *ast.Ident {
 7  .  .  .  .  NamePos: gateway/prod/gateway.go:10:6
 8  .  .  .  .  Name: "gateway"
 9  .  .  .  .  Obj: *ast.Object {
10  .  .  .  .  .  Kind: type
11  .  .  .  .  .  Name: "gateway"
12  .  .  .  .  .  Decl: *(obj @ 5)
13  .  .  .  .  }
14  .  .  .  }
15  .  .  .  Assign: -
16  .  .  .  Type: *ast.StructType {
17  .  .  .  .  Struct: gateway/prod/gateway.go:10:14
18  .  .  .  .  Fields: *ast.FieldList {
19  .  .  .  .  .  Opening: gateway/prod/gateway.go:10:21
20  .  .  .  .  .  List: []*ast.Field (len = 1) {
21  .  .  .  .  .  .  0: *ast.Field {
22  .  .  .  .  .  .  .  Type: *ast.StarExpr {
23  .  .  .  .  .  .  .  .  Star: gateway/prod/gateway.go:11:2
24  .  .  .  .  .  .  .  .  X: *ast.Ident {
25  .  .  .  .  .  .  .  .  .  NamePos: gateway/prod/gateway.go:11:3
26  .  .  .  .  .  .  .  .  .  Name: "payment"
27  .  .  .  .  .  .  .  .  }
28  .  .  .  .  .  .  .  }
29  .  .  .  .  .  .  }
30  .  .  .  .  .  }
31  .  .  .  .  .  Closing: gateway/prod/gateway.go:12:1
32  .  .  .  .  }
33  .  .  .  .  Incomplete: false
34  .  .  .  }
35  .  .  }
36  .  }
37  .  Rparen: -
38  }
```

### simple composition struct gateway in different packages

```
type gateway struct {
	other.Payment
}
    
 0  *ast.GenDecl {
 1  .  TokPos: gateway/prod/gateway.go:11:1
 2  .  Tok: type
 3  .  Lparen: -
 4  .  Specs: []ast.Spec (len = 1) {
 5  .  .  0: *ast.TypeSpec {
 6  .  .  .  Name: *ast.Ident {
 7  .  .  .  .  NamePos: gateway/prod/gateway.go:11:6
 8  .  .  .  .  Name: "gateway"
 9  .  .  .  .  Obj: *ast.Object {
10  .  .  .  .  .  Kind: type
11  .  .  .  .  .  Name: "gateway"
12  .  .  .  .  .  Decl: *(obj @ 5)
13  .  .  .  .  }
14  .  .  .  }
15  .  .  .  Assign: -
16  .  .  .  Type: *ast.StructType {
17  .  .  .  .  Struct: gateway/prod/gateway.go:11:14
18  .  .  .  .  Fields: *ast.FieldList {
19  .  .  .  .  .  Opening: gateway/prod/gateway.go:11:21
20  .  .  .  .  .  List: []*ast.Field (len = 1) {
21  .  .  .  .  .  .  0: *ast.Field {
22  .  .  .  .  .  .  .  Type: *ast.SelectorExpr {
23  .  .  .  .  .  .  .  .  X: *ast.Ident {
24  .  .  .  .  .  .  .  .  .  NamePos: gateway/prod/gateway.go:12:2
25  .  .  .  .  .  .  .  .  .  Name: "other"
26  .  .  .  .  .  .  .  .  }
27  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
28  .  .  .  .  .  .  .  .  .  NamePos: gateway/prod/gateway.go:12:8
29  .  .  .  .  .  .  .  .  .  Name: "Something"
30  .  .  .  .  .  .  .  .  }
31  .  .  .  .  .  .  .  }
32  .  .  .  .  .  .  }
33  .  .  .  .  .  }
34  .  .  .  .  .  Closing: gateway/prod/gateway.go:13:1
35  .  .  .  .  }
36  .  .  .  .  Incomplete: false
37  .  .  .  }
38  .  .  }
39  .  }
40  .  Rparen: -
41  }
```

### pointer composition struct gateway in different packages
```
type gateway struct {
	*other.Payment
}
    
 0  *ast.GenDecl {
 1  .  TokPos: gateway/prod/gateway.go:11:1
 2  .  Tok: type
 3  .  Lparen: -
 4  .  Specs: []ast.Spec (len = 1) {
 5  .  .  0: *ast.TypeSpec {
 6  .  .  .  Name: *ast.Ident {
 7  .  .  .  .  NamePos: gateway/prod/gateway.go:11:6
 8  .  .  .  .  Name: "gateway"
 9  .  .  .  .  Obj: *ast.Object {
10  .  .  .  .  .  Kind: type
11  .  .  .  .  .  Name: "gateway"
12  .  .  .  .  .  Decl: *(obj @ 5)
13  .  .  .  .  }
14  .  .  .  }
15  .  .  .  Assign: -
16  .  .  .  Type: *ast.StructType {
17  .  .  .  .  Struct: gateway/prod/gateway.go:11:14
18  .  .  .  .  Fields: *ast.FieldList {
19  .  .  .  .  .  Opening: gateway/prod/gateway.go:11:21
20  .  .  .  .  .  List: []*ast.Field (len = 1) {
21  .  .  .  .  .  .  0: *ast.Field {
22  .  .  .  .  .  .  .  Type: *ast.StarExpr {
23  .  .  .  .  .  .  .  .  Star: gateway/prod/gateway.go:12:2
24  .  .  .  .  .  .  .  .  X: *ast.SelectorExpr {
25  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
26  .  .  .  .  .  .  .  .  .  .  NamePos: gateway/prod/gateway.go:12:3
27  .  .  .  .  .  .  .  .  .  .  Name: "other"
28  .  .  .  .  .  .  .  .  .  }
29  .  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
30  .  .  .  .  .  .  .  .  .  .  NamePos: gateway/prod/gateway.go:12:9
31  .  .  .  .  .  .  .  .  .  .  Name: "Something"
32  .  .  .  .  .  .  .  .  .  }
33  .  .  .  .  .  .  .  .  }
34  .  .  .  .  .  .  .  }
35  .  .  .  .  .  .  }
36  .  .  .  .  .  }
37  .  .  .  .  .  Closing: gateway/prod/gateway.go:13:1
38  .  .  .  .  }
39  .  .  .  .  Incomplete: false
40  .  .  .  }
41  .  .  }
42  .  }
43  .  Rparen: -
44  }
```