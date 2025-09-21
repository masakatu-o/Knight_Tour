package hamiltonpath

import (
	"fmt"
	"sync"
	"time"
)

const (
	Normal = iota
	Euler
	Magic
)

const (
	Single = iota
	Competition
	Multipath
)

// config
var (
	xsize int
	ysize int
	size  int

	startc int
	endc   int
	seedc  int

	solutions  int64
	solutionNo int
	steplimit  int64

	Totalsolutions int64
	Totalsteps     int64
	threadcnt      int

	mu sync.Mutex

	fswitch bool

	slct int

	threadmode int

	competethreads int

	cross bool

	tracefile string

	wg *sync.WaitGroup

	ch chan int
)

// global
type hGlobal struct {
	start int
	end   int
	seed  int

	graph    []node
	sequence []int
	enD      *node

	solutionCt  int64
	steplimitCt int64

	ti time.Time

	stepNum int

	// euler
	phase int8
	retP  int
	level int

	isstack [50]island

	// magic
	row      []rod
	columns  []rod
	diagonal []rod
	board    []int
	th       [5]int
}

func global(startP int, endP int, seedP int) (hg *hGlobal) {

	hg = &hGlobal{}

	hg.solutionCt = 0
	hg.steplimitCt = 0

	hg.ti = time.Now()

	if startP > size || endP > size {
		//
	}
	hg.start = startP
	hg.end = endP
	hg.seed = seedP

	// knight tour の場合
	kngraphInit(hg)

	hg.enD = &hg.graph[hg.end]

	sequence := make([]int, size)
	for i := 0; i < size; i++ {
		sequence[i] = -1
	}
	hg.sequence = sequence

	hg.stepNum = -1

	return
}

func thread(start int, th [5]int, hg *hGlobal) {

	if tracefile != "" {
		start = knalign(tracefile, hg)
	}

	switch slct {
	case Normal:
		//stepN(&(hg.graph[hg.start]), hg)
		stepN(&(hg.graph[start]), hg)
	case Euler:
		area := make([]*node, size)
		for i := 0; i < size; i++ {
			area[i] = &hg.graph[i]
		}

		ps := pset{1, 0, nil}
		hg.isstack[0].area = area

		hg.enD.gate = true

		stepE(&hg.graph[start], ps, hg)
	case Magic:
		initMagic(th, hg)
		//stepM(&hg.graph[start], tail{nil, nil, 64}, hg)
		//stepM(&hg.graph[start], paratail{nil, nil, size}, hg)
		stepM(&hg.graph[start], paratail{nil, hg.enD, size}, hg)
	default:
		//
	}

	stopenumerated(hg)
}

func threadR() {

	wg.Add(competethreads)

	for r := seedc; r < seedc+competethreads; r++ {

		hg := global(startc, endc, r)
		go thread(startc, [5]int{0, 0, 0, 0}, hg)
	}

	wg.Wait()
}

func threadM() {

	// 仮のグローバル構造体
	hgg := global(startc, endc, 1)

	enD := hgg.enD
	starT := &hgg.graph[hgg.start]

	var thL [][5]int

	for _, sq := range starT.neighborList() {

		if sq == enD {
			continue
		}

		for _, sqq := range sq.neighborList() {

			if sqq == enD || sqq == starT || sqq == sq {
				continue
			}

			for _, sqqq := range sqq.neighborList() {

				if sqqq == enD || sqqq == starT || sqqq == sq || sqqq == sqq {
					continue
				}

				for _, sqqqq := range sqqq.neighborList() {

					if sqqqq == enD || sqqqq == starT || sqqqq == sq || sqqqq == sqq || sqqqq == sqqq {
						continue
					}

					for _, sqqqqq := range sqqqq.neighborList() {
						if sqqqqq == enD || sqqqqq == starT || sqqqqq == sq || sqqqqq == sqq || sqqqqq == sqqq || sqqqqq == sqqqq {
							continue
						}

						thL = append(thL, [5]int{sq.pos, sqq.pos, sqqq.pos, sqqqq.pos, sqqqqq.pos})
					}
				}
			}
		}
	}

	thcnt := 0

label:
	for _, th := range thL {

		hg := global(startc, endc, seedc)

		ps0 := &hg.graph[hg.start]

		var ps1 *node

		for _, ps := range th {

			ps1 = &hg.graph[ps]

			hg.stepNum++
			ps0.leaVe(hg)
			hg.sequence[hg.stepNum] = ps0.pos
			ps0.unpassed = false

			de, _ := gather(ps0, ps1, hg)

			if de != 0 {
				//drawD(3, "", ps0, ps1, reapL, hg)
				//go stopfound(hg)
				continue label
			}

			ps0 = ps1
		}

		thcnt++

		hg.th = th
		wg.Add(1)
		go thread(th[4], th, hg)

	}

	mu.Lock()
	threadcnt = thcnt
	mu.Unlock()

	fmt.Printf("Multi Thresd Start %v Threads seed = %v\n", thcnt, seedc)

	wg.Wait()

	fmt.Printf("Total solutions: %v Total try: %v\n", Totalsolutions, Totalsteps)
}

func Hamiltonpath(
	xsizeC int,
	ysizeC int,
	startC int,
	endC int,
	seedC int,
	solutionsC int64,
	steplimitC int64,
	fswitchC bool,
	slctC int,
	threadmodeC int,
	threadsC int,
	crossC bool,
	tracefileC string,

) {

	//config
	if xsizeC <= 0 {
		stopparamerr("xsize")
	}
	xsize = xsizeC

	if ysizeC <= 0 {
		stopparamerr("ysize")
	}
	ysize = ysizeC

	startc = startC
	endc = endC
	seedc = seedC

	size = xsize * ysize

	if solutionsC <= 0 {
		stopparamerr("solutions")
	}
	solutions = solutionsC

	if steplimitC <= 0 {
		stopparamerr("steplimit")
	}
	steplimit = steplimitC

	fswitch = fswitchC
	//filepass = filepassC
	//filecaption = filecaptionC

	Totalsolutions = 0
	Totalsteps = 0
	solutionNo = 0

	if slct > Magic {
		// Mgic to be contruct
		stopparamerr("slct")
	}
	slct = slctC

	threadmode = threadmodeC
	competethreads = threadsC

	if cross && (xsize != ysize) {
		stopparamerr("cross")
	}
	cross = crossC

	tracefile = tracefileC

	var w sync.WaitGroup
	wg = &w

	ch = make(chan int)

	switch threadmode {
	case Single:
		competethreads = 1
		fmt.Printf("Single Thread Start seed = %v\n", seedc)
		threadR()
	case Competition:
		fmt.Printf("Compete Thread Start %v Threads seed = %v-%v\n", competethreads, seedc, seedc+competethreads)
		threadR()
	case Multipath:
		threadM()
	}
}
