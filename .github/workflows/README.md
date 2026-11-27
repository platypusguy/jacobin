Workflows:

| `YML File` | `Runs on push?` | `Runs on demand?` | `Description` |
| :------------ | :--- | :--- | :--- |
|<img width=90/>|<img width=60/>|<img width=600/>|
| build_executables | no | yes | Builds executables for 3 O/Ses. If a release tag is supplied, these files will be added to that release. |
| go  | yes | no | Unit testing on 3 AMD/Intel O/Ses |
| run_platypusguy_jacobin | no | yes | Run Jacotest against Jacobin on 6 AMD/Intel O/Ses |

