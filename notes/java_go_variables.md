| Java Variable Type | Class | Golang  Representation |
| --- | --- | --- |
| byte | n/a | types.JavaByte (int8) |
| byte[] | [B | []types.JavaByte in an Object ([]int8 in an Object) |
| boolean | n/a | int64(1) or int64(0) |
| char, int, long, short | n/a | int64 |
| float, double | n/a | float64 |
| boolean[] | [Z | []types.JavaByte in an Object |
| char[], int[], long[], short[] | [C, [I, [J, [S | []int64 in an Object |
| String | Ljava/lang/String; | []types.JavaByte in an Object |

