# bridge-howell-movement
A howell movement match table generator written in Go.
It can find any legal movement table for given numbers of table.
It becomes slower by tables increasing. In an Intel Xeon E5-2650 0 CPU with 8 gouroutines, with go 1.13.7 for windows amd64, calculating 8 tables will take about 7 secs, however 11 tables will take about 81 mins.
Any PR that boost it is appreciated.

# Usage
```batch
set GO111MODULE=on
go build howell.go
```
Then run the executable file directly.

Example
```batch
>howell.exe

go routines 8
please input the number of tables
4
players solution found...
counts in thread:  105
1: [{NS:1 Board:4 EW:3} {NS:2 Board:6 EW:6} {NS:4 Board:3 EW:5} {NS:7 Board:2 EW:8}]
2: [{NS:2 Board:5 EW:4} {NS:3 Board:7 EW:7} {NS:5 Board:4 EW:6} {NS:1 Board:3 EW:8}]
3: [{NS:3 Board:6 EW:5} {NS:4 Board:1 EW:1} {NS:6 Board:5 EW:7} {NS:2 Board:4 EW:8}]
4: [{NS:4 Board:7 EW:6} {NS:5 Board:2 EW:2} {NS:7 Board:6 EW:1} {NS:3 Board:5 EW:8}]
5: [{NS:5 Board:1 EW:7} {NS:6 Board:3 EW:3} {NS:1 Board:7 EW:2} {NS:4 Board:6 EW:8}]
6: [{NS:6 Board:2 EW:1} {NS:7 Board:4 EW:4} {NS:2 Board:1 EW:3} {NS:5 Board:7 EW:8}]
7: [{NS:7 Board:3 EW:2} {NS:1 Board:5 EW:5} {NS:3 Board:2 EW:4} {NS:6 Board:1 EW:8}]
[[{"ns":1,"board":4,"ew":3},{"ns":2,"board":6,"ew":6},{"ns":4,"board":3,"ew":5},
{"ns":7,"board":2,"ew":8}],[{"ns":2,"board":5,"ew":4},{"ns":3,"board":7,"ew":7},
{"ns":5,"board":4,"ew":6},{"ns":1,"board":3,"ew":8}],[{"ns":3,"board":6,"ew":5},
{"ns":4,"board":1,"ew":1},{"ns":6,"board":5,"ew":7},{"ns":2,"board":4,"ew":8}],[
{"ns":4,"board":7,"ew":6},{"ns":5,"board":2,"ew":2},{"ns":7,"board":6,"ew":1},{"ns":
3,"board":5,"ew":8}],[{"ns":5,"board":1,"ew":7},{"ns":6,"board":3,"ew":3},{"ns":
1,"board":7,"ew":2},{"ns":4,"board":6,"ew":8}],[{"ns":6,"board":2,"ew":1},{"ns":
7,"board":4,"ew":4},{"ns":2,"board":1,"ew":3},{"ns":5,"board":7,"ew":8}],[{"ns":
7,"board":3,"ew":2},{"ns":1,"board":5,"ew":5},{"ns":3,"board":2,"ew":4},{"ns":6,
"board":1,"ew":8}]]
648.037ms
```
