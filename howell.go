package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"os"
	"strconv"
	"sync"
	"time"
)

type Spliter struct {
	Start []int
	End   []int // included
}

type Seat struct {
	NS    int `json:"ns"`
	Board int `json:"board"`
	EW    int `json:"ew"`
}

var Solution [][]Seat
var BoardSolution [][]Seat
var Choice []int

type CmpRes int

const (
	Less CmpRes = -1 + iota
	Eq
	More
)

func SliceCmp(a *[]int, b *[]int) CmpRes {
	for i := 0; i < len(*a); i++ {
		if (*a)[i] < (*b)[i] {
			return Less
		}
		if (*a)[i] > (*b)[i] {
			return More
		}
	}
	return Eq
}

func NextIter(prev *[]int, lastOne []int) bool {
	if SliceCmp(prev, &lastOne) != Less {
		return false
	}
	for i := len(*prev) - 1; i >= 0; i-- {
		mx := i*2 + 2 // 这一位的最大数字 2 4 6 8....
		if (*prev)[i] < mx {
			(*prev)[i] += 1
			return true
		} else {
			(*prev)[i] = 0
		}
	}
	panic("cannot go here")
}

func GetSlice(length int, num int) []int {
	res := make([]int, length)
	for i := length - 1; i >= 0; i-- {
		div := i*2 + 3 // 这一位的除数 3 5 7 9....，余数就是这一位的值
		res[i], num = num%div, num/div
	}
	return res
}

func SplitMargin(length int, gocount int, totalPlayerComb int) (res []Spliter) {
	res = make([]Spliter, gocount)
	margin := totalPlayerComb / gocount
	for i, _ := range res {
		res[i].Start = GetSlice(length, i*margin)
		if i < gocount-1 {
			res[i].End = GetSlice(length, i*margin+margin-1)
		} else {
			res[i].End = GetSlice(length, totalPlayerComb-1)
		}
	}
	return
}

func TotalPlayerComb(realPlayers int) int {
	totalComb := 1 // 首轮牌手捉对情况的总组合数。只调换座位或桌号不重复计入
	for i := 1; i < realPlayers; i += 2 {
		totalComb *= i
	}
	return totalComb
}

func CheckSeat(seatOrder []int) (res [][]Seat, ok bool) {
	seatOrder = append([]int{0}, seatOrder...)
	ls := len(seatOrder)
	maxN := ls * 2
	unused := make([]int, maxN)
	first := make([]Seat, ls)
	checker := make([]map[int]bool, maxN)

	for i, _ := range unused {
		unused[i] = 1 + i
	}
	for i, _ := range first {
		first[i].NS = unused[0]
		ew := seatOrder[ls-i-1] + 1
		first[i].EW = unused[ew]
		checker[unused[0]-1] = map[int]bool{unused[ew]: true}
		checker[unused[ew]-1] = map[int]bool{unused[0]: true}
		unused = append(unused[1:ew], unused[ew+1:]...)
	}
	res = append(res, first)

	for round := 1; round < maxN-1; round ++ {
		next := make([]Seat, ls)
		copy(next, res[len(res)-1])
		for i, _ := range next {
			switch next[i].NS {
			case maxN:
			case maxN - 1:
				next[i].NS = 1
			default:
				next[i].NS += 1
			}
			switch next[i].EW {
			case maxN:
			case maxN - 1:
				next[i].EW = 1
			default:
				next[i].EW += 1
			}
			if _, ok := checker[next[i].NS-1][next[i].EW]; ok {
				return nil, false
			}
			if _, ok := checker[next[i].EW-1][next[i].NS]; ok {
				return nil, false
			}
			checker[next[i].NS-1][next[i].EW] = true
			checker[next[i].EW-1][next[i].NS] = true
		}
		res = append(res, next)
	}
	return res, true
}

func FindASeat(start []int, end []int, wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0
	for {
		if seat, ok := CheckSeat(start); ok {
			Solution = seat
			return
		}
		i ++
		if (i%10000 == 0 && Solution != nil) || !NextIter(&start, end) {
			return
		}
	}
}

// then boards
func GetBoard(length int, boardId int) (choiceBoard []int) {
	total := length*2 - 1
	choiceBoard = make([]int, length)

	var i int
	for div := length; div <= total; div++ {
		boardId, i = boardId/div, boardId%div
		choiceBoard[total-div] = i
	}
	return
}

func GetRealBoard(choiceBoard []int) (board []int) {
	choice := make([]int, len(Choice))
	copy(choice, Choice)
	for _, v := range choiceBoard {
		board = append(board, choice[v])
		choice = append(choice[:v], choice[v+1:]...)
	}
	return
}

func TotalBoardComb(round int, tb int) (res int) {
	res = 1
	for round >= tb {
		res *= round
		round--
	}
	return res
}

func NextBoard(prev *[]int, last []int) bool {
	if SliceCmp(prev, &last) != Less {
		return false
	}
	for i := len(last) - 1; i >= 0; i-- {
		if (*prev)[i] < len(last)*2-i-2 { // max number of this digit
			(*prev)[i]++
			return true
		} else {
			(*prev)[i] = 0
		}
	}
	panic("cannot go here")
}

func SplitPerm(length int, gocount int, totalBoardComb int) (res []Spliter) {
	res = make([]Spliter, gocount)
	margin := totalBoardComb / gocount
	fmt.Println("counts in thread: ", margin)
	for i, _ := range res {
		res[i].Start = GetBoard(length, i*margin)
		if i < gocount-1 {
			res[i].End = GetBoard(length, i*margin+margin-1)
		} else {
			res[i].End = GetBoard(length, totalBoardComb-1)
		}
	}
	return
}

func doubleCopy(src [][]Seat) (dst [][]Seat) {
	for _, sub := range src {
		sd := make([]Seat, len(sub))
		copy(sd, sub)
		dst = append(dst, sd)
	}
	return
}

func CheckBoard(choiceBoard []int, seat [][]Seat) (ok bool) {
	realBoard := GetRealBoard(choiceBoard)
	first := seat[0]
	maxR := len(seat)
	checker := make([]map[int]bool, maxR+1)
	for i, v := range realBoard {
		first[i].Board = v
		checker[first[i].NS-1] = map[int]bool{v: true}
		checker[first[i].EW-1] = map[int]bool{v: true}
	}
	for i := 1; i < maxR; i++ {
		for j, _ := range seat[i] {
			var bd int
			if seat[i-1][j].Board < maxR {
				bd = seat[i-1][j].Board + 1
			} else {
				bd = 1
			}
			if _, ok := checker[seat[i][j].NS-1][bd]; ok {
				return false
			}
			if _, ok := checker[seat[i][j].EW-1][bd]; ok {
				return false
			}
			checker[seat[i][j].NS-1][bd] = true
			checker[seat[i][j].EW-1][bd] = true
			seat[i][j].Board = bd
		}
	}
	return true
}

func FindABoard(start []int, end []int, wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0
	res := doubleCopy(Solution)
	for {
		if ok := CheckBoard(start, res); ok {
			BoardSolution = res
			return
		}
		i++
		if i%1000000 == 0 {
			fmt.Println(i)
			if BoardSolution != nil {
				return
			}
		}
		if !NextBoard(&start, end) {
			return
		}
	}
}

func main() {
	goCount, err := cpu.Counts(false)
	if err != nil {
		panic(err)
	}
	fmt.Println("go routines", goCount)
	var tbStr string
	fmt.Println("please input the number of tables")
	st := time.Now()
	_, err = fmt.Scanln(&tbStr)
	if err != nil {
		panic(err)
	}
	tb, err := strconv.Atoi(tbStr)
	if err != nil {
		panic(err)
	}
	if tb < 4 {
		fmt.Println("tables for howell is at least 4")
		os.Exit(-1)
	}

	realPlayers := tb * 2
	totalPlayerComb := TotalPlayerComb(realPlayers)

	playRoutes := SplitMargin(tb-1, goCount, totalPlayerComb)
	wg := &sync.WaitGroup{}
	wg.Add(goCount)
	for _, playRoute := range playRoutes {
		go FindASeat(playRoute.Start, playRoute.End, wg)
	}
	wg.Wait()

	fmt.Println("players solution found...")
	for i := 1; i < realPlayers; i++ {
		Choice = append(Choice, i)
	}

	totalBoardComb := TotalBoardComb(realPlayers-1, tb)

	boardRoutes := SplitPerm(tb, goCount, totalBoardComb)
	wg = &sync.WaitGroup{}
	wg.Add(goCount)
	for _, boardRoute := range boardRoutes {
		go FindABoard(boardRoute.Start, boardRoute.End, wg)
	}
	wg.Wait()
	for i, v := range BoardSolution {
		fmt.Printf("%d: %+v\n", i+1, v)
	}
	res, _ := json.Marshal(BoardSolution)
	fmt.Println(string(res))
	fmt.Println(time.Now().Sub(st))
}
