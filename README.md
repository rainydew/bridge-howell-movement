# bridge-howell-movement
A howell movement match table generator written in Go.
It can find any legal movement table for given numbers of table.
It becomes slower by tables increasing. In an Intel Xeon E5-2650 0 CPU with 8 gouroutines, calculating 8 tables will take about 7 secs, however 11 tables will take about 81 mins.
Any PR that boost it is appreciated.
