package hamiltonpath

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

var sequenceR []int

func outSequence(filename string, hg *hGlobal) {

	//file, err := os.Create(fmt.Sprintf("Sequence%v.csv", filenum))
	file, err := os.Create(filename + ".csv")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	var csvl [][]string

	csvl = append(csvl, []string{strconv.Itoa(xsize)})
	csvl = append(csvl, []string{strconv.Itoa(ysize)})
	csvl = append(csvl, []string{strconv.Itoa(size)})

	for _, n := range hg.sequence {

		str := []string{strconv.Itoa(n)}

		csvl = append(csvl, str)

	}

	w := csv.NewWriter(file)
	w.WriteAll(csvl) // 一度にすべて書き込む

	err = file.Sync()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func inSequence(filename string, hg *hGlobal) (sequence []int) {

	var reverse bool

	if filename[0] == '-' {
		filename = filename[1:]
		reverse = true
	}

	file, err := os.Open(filename + ".csv")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	r := csv.NewReader(file)
	csvl, err := r.ReadAll() // csvを一度に全て読み込む

	if err != nil {
		panic(err)
	}

	for i, n := range csvl {

		k, _ := strconv.Atoi(n[0])

		switch i {
		case 0:
			if k != xsize {
				stoperr("inSequence xsize err", hg)
			}
		case 1:
			if k != ysize {
				stoperr("inSequence ysize err", hg)
			}
		case 2:
			if k != size {
				stoperr("inSequence size err", hg)
			}
		default:
			sequence = append(sequence, k)
		}
	}

	if len(sequence) != size {
		stoperr("inSequence len(sequence) != size", hg)
	}

	var q int
	if reverse {
		for i := range size / 2 {
			q = sequence[i]
			sequence[i] = sequence[size-1-i]
			sequence[size-1-i] = q
		}
	}

	return
}

/*
func checkSequence(now, next *square, ret bool) bool {

	if !checkSflag {
		return false
	}

	nextP := sequenceR[stepNum]

	if next.pos == nextP && !ret {
		return true
	}

	//if next.pos != nextP && ret {
	//	return true
	//}

	return false
}
*/
