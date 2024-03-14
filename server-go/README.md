## Build Local
- To install the required packages
```bash
go mod download
```
- To run local
```bash
go run .
```
- To commit
```bash
make go-publish
```
### Go Package Adder Script
The script takes one argument, a directory, and iterates over all `.go` files in that directory. For each file, it checks if there is a line starting with `package`. If not, it adds a line `package [directory]` at the beginning of the file.


```bash
./declare-pkg.sh [directory]