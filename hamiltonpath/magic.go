package hamiltonpath

//import "golang.org/x/tools/go/analysis/passes/bools"

// tail や reserveは後半で有効だが、前半では効果がない。
// コクーンを利用することも後半では有効。
//　前半で有効なのはaからbまでの最短ステップ数を割り出すこと。
// 2025/2/12 v0.2 8X8のセミマジックの解を得る。
//
// 2025/4/14 tailの組み込み完成？ではなかった。
// 2025/5/2 - 5/5
//   tailの廃止
//   reserveのthresholdの動的書き換え廃止。
//   threshold の計算を最後に回す。
// 2025/5/17 tail復活
// 6/17 reserve putの順序を先頭に。
type rod struct {
	cnt       int
	rest      int
	threshold int
	diff      []int
	flag      bool
}

type parath struct {
	rd *rod
	th int
}

type paratail struct {
	bp    *node
	fp    *node
	order int
}

type param struct {
	record     [][3]int
	threashold []parath
}

type respara struct {
	x    int
	y    int
	rest int
	rd   *rod
}

// stepを縦横斜めの軸に記録。
// 確定した予約を記録
// magicの条件から外れればreTがエラー番号を返す。
// 回復用パラメーターpaと状態維持用パラメータtaNを返す。

func mark(now *node, next *node, ta paratail, hg *hGlobal, reapList twins) (reT int, pa param, nta paratail) {

	footprint := hg.stepNum + 2

	var flist []*rod

	monitor := func(sw bool) {
		if !sw {
			return
		}
		drawD(2, "", now, next, reapList, hg)
		mPrint(hg)
	}

	b := func(sw bool) {
		if !sw {
			stoperr("b", hg)
		}
	}

	if false {
		monitor(false)
		b(false)
	}

	flagclear := func() {

		for _, rd := range flist {
			rd.flag = false
		}
	}

	defer func() {
		flagclear()

		//if len(pa.record) > 2 {
		//	monitor(true)
		//	b(true)
		//}

		if reT != 0 {
			//if reT >= 1000 {
			//if nta.order == 0 {
			if hg != nil {
			}
			//	}
			//monitor(true)
			//b(true)
			//}
		}
	}()

	//各rodのthresholdがfootprintより小さければfalse
	checkth := func() bool {

		rDthreshold := func(rd *rod, rV int, nV int, tV int) (th int) {

			th = rd.threshold

			tos := func(dis int) (step int) {

				if dis < 0 {
					dis = -dis
				}

				if dis%2 == 0 {
					step = dis / 2
				} else {
					step = dis/2 + 1
				}
				return
			}

			d0 := tos(rV - nV)
			d1 := tos(rV - tV)

			bothends := func(rd *rod) int {

				var (
					e1 int
					e2 int
					s  int
				)

				for e1, s = range rd.diff {
					if hg.board[s] == 0 {
						break
					}
				}

				for e2 = len(rd.diff) - 1; e2 > 0; e2-- {
					if hg.board[rd.diff[e2]] == 0 {
						break
					}
				}

				return e2 - e1
			}

			dis := bothends(rd)
			if dis == 0 {

			}

			avg := rd.rest / rd.cnt

			var (
				cr  int
				nth int
				ret int
			)

			switch rd.cnt {
			case 0, 1:
				stoperr("case zero or one", hg)
			case 2:

				switch dis {
				case 1:
					cr = 2
				case 2:
					cr = 1
				case 3:
					cr = 2
				case 4:
					cr = 2
				case 5:
					cr = 2
				case 6:
					cr = 2
				default:
					cr = 4
				}

				ret = 2
			case 3:

				switch dis {
				case 2:
					cr = 3
				case 3:
					cr = 3
				case 4:
					cr = 3
				case 5:
					cr = 3
				case 6:
					cr = 3
				default:
					cr = 5
				}

				ret = 3
			case 4:

				switch dis {
				case 3:
					cr = 7
				case 4:
					cr = 6
				case 5:
					cr = 5
				case 6:
					cr = 4
				default:
					cr = 6
				}

				ret = 4
			case 5:

				switch dis {
				case 4:
					cr = 11
				case 5:
					cr = 6
				case 6:
					cr = 7
				default:
					cr = 7
				}

				ret = 5
			default:

				switch dis {
				case 5:
					cr = 10
				case 6:
					cr = 10
				default:
					cr = 9
				}

				ret = 6
			}

			nth = avg - cr

			if avg-footprint+d0 < cr {
				//monitor(true)
				th = 0
				reT = 200 + ret
				return
			}

			if nta.order-avg-d1 < cr {
				//monitor(true)
				th = 0
				reT = 300 + ret
				return
			}

			rd.threshold = nth
			return
		}

		var (
			max  int
			mini int
		)

		check := func(rods []rod, nV int, tV int) bool {

			var th int

			for i := range len(rods) {

				rd := &rods[i]

				if rd.flag {
					if rd.cnt != 0 {
						if th = rDthreshold(rd, i, nV, tV); th == 0 {
							//reT = 101
							return false
						}

						pa.threashold = append(pa.threashold, parath{rd, th})
					} else {
						//monitor(true)
						//b()
					}

				} else {

					if footprint > rd.threshold {
						//monitor(true)
						reT = 102
						return false
					}
				}

				if rd.cnt < mini {
					mini = rd.cnt
				}

				if rd.cnt > max {
					max = rd.cnt
				}
			}

			return true
		}

		var r int

		if footprint > 40 {
			r = 4
		} else {
			r = 5
		}

		mini = 100
		max = 0
		if !check(hg.row, next.y, nta.fp.y) {
			return false
		}

		if max-mini > r {
			reT = 103
			//monitor(true)
			return false
		}

		mini = 100
		max = 0
		if !check(hg.columns, next.x, nta.fp.x) {
			return false
		}

		if max-mini > r {
			//if ta.order == size {
			reT = 104
			//monitor(true)
			return false
			//}
		}

		if cross {
			if !check(hg.diagonal, next.x, nta.fp.y) {
				return false
			}
		}

		return true
	}

	// (x,y)にorderを書き込む。
	// 書き込みの結果が無効ならreTにエラー番号を設定。
	// 書き込みの復元はパラメータpaを媒介にunmarkが行う。

	var (
		put     func(int, int, int)
		reserve func(int, int, int, *rod)
	)

	reserve = func(x int, y int, order int, rd *rod) {

		searchterm := func(sq *node, sqq *node) (bool, int) {

			bp := sq
			fp := sqq
			dis := 1

			var (
				ar    *node
				bor   int
				pflag bool
				mflag bool
			)

			for {

				bor = hg.board[fp.pos]

				// 予約に当たった。
				if bor != 0 {

					if order-bor == dis {
						//monitor(true)
						// order > bor
						mflag = true
						//goto L
						break
					}

					if bor-order == dis {
						//monitor(true)
						// order < bor
						pflag = true
						//goto L
						break
					}
					// ふたつの予約に整合性なし。
					//monitor(true)
					return false, 0

					//L:
				}

				if ar = bp.arm(fp); ar == nil {
					break
				}

				bp = fp
				fp = ar
				dis++
			}

			// searchterm main

			// 方向が対立
			//if pflag && mflag {
			//monitor(true)
			//	return false, 0, nil, nil
			//}

			// order方向へ pflagと一致すべし。
			//tf0 := order-dis < footprint
			// orderと反対方向へ mflagと一致すべし。
			//tf1 := order+dis > ta.order

			switch {
			case pflag && mflag:
				// 方向が対立
				//monitor(true)
				return false, 0
			case pflag:
				//monitor(true)
				//if order+dis > ta.order {
				if order+dis > size {
					return false, 0
				}
				return true, 1
			case mflag:
				//monitor(true)
				if order-dis < footprint {
					return false, 0
				}
				return true, -1
			default:
				//monitor(true)
				return true, 0
			}

			// searchterm
		}

		fill := func(dir int, sq *node, sqq *node) bool {

			bp := sq
			fp := sqq
			dis := dir

			var (
				ar  *node
				bor int
			)

			for {

				bor = hg.board[fp.pos]

				// 予約に当たった。
				if bor != 0 {

					if bor-order != dis {
						return false
					}
				} else {

					if put(fp.x, fp.y, order+dis); reT != 0 {
						return false
					}
				}

				if ar = bp.arm(fp); ar == nil {
					break
				}

				bp = fp
				fp = ar
				dis += dir
			}

			return true
		}

		// reverse main

		if order < footprint || order > ta.order {
			//monitor(true)
			reT = 900
			return
		}

		//if hg.board[x+y*xsize] == 0 {
		pa.threashold = append(pa.threashold, parath{rd, rd.threshold})
		rd.threshold = order

		if put(x, y, order); reT != 0 {
			return
		}

		//}

		var (
			flag bool
			dir0 int
			dir1 int
		)

		pos := x + y*xsize
		sq := &hg.graph[pos]

		//monitor(true)
		//b(true)

		if next.neighborLen() == 2 {
			if sq == next {
				return
			} else {
				if sq.include(next.neighborList()) {
					//monitor(true)
					return
				}
			}
		}

		switch {
		case sq.neighborLen() == 2:
			l := sq.neighborList()
			//monitor(true)
			if flag, dir0 = searchterm(sq, l[0]); !flag {
				//monitor(true)
				reT = 1000
				return
			}

			if flag, dir1 = searchterm(sq, l[1]); !flag {
				//monitor(true)
				reT = 1001
				return
			}

			switch dir0 * dir1 {
			case 0:
				if dir0 == 0 && dir1 == 0 {
					return
				}

				if dir0 != 0 {
					if !fill(dir0, sq, l[0]) {
						//monitor(true)
						reT = 1002
						return
					}

					if !fill(-dir0, sq, l[1]) {
						//monitor(true)
						reT = 1003
						return
					}
				} else {
					if !fill(-dir1, sq, l[0]) {
						//monitor(true)
						reT = 1004
						return
					}

					if !fill(dir1, sq, l[1]) {
						//monitor(true)
						reT = 1005
						return
					}
				}
			case 1:
				reT = 1006
				return
			case -1:
				if !fill(dir0, sq, l[0]) {
					//monitor(true)
					reT = 1007
					return
				}

				if !fill(dir1, sq, l[1]) {
					//monitor(true)
					reT = 1008
					return
				}
			}

		case sq.handsCt() == 1:

			l := sq.handsList()

			if l[0] == hg.enD {
				//monitor(true)
				return
			}

			if flag, dir0 = searchterm(sq, l[0]); !flag {
				//monitor(true)
				reT = 1009
				return
			}

			if dir0 != 0 {
				if !fill(dir0, sq, l[0]) {
					//monitor(true)
					reT = 1010
					return
				}
			}

		default:
			return
		}

		//monitor(true)
		if false {
			//
		}

		//L1:

		//reserve
	}

	put = func(x int, y int, order int) {

		write := func() {

			hg.board[x+y*xsize] = order
			pa.record = append(pa.record, [3]int{x, y, order})

			rDwrite := func(rd *rod) {
				rd.cnt--
				rd.rest -= order
				rd.flag = true
				flist = append(flist, rd)
			}

			rDwrite(&hg.row[y])

			rDwrite(&hg.columns[x])

			if cross {
				if x == y {
					rDwrite(&hg.diagonal[0])
				}
				if x == ysize-y-1 {
					rDwrite(&hg.diagonal[1])
				}
			}

			// write
		}

		// put main

		if hg.board[x+y*xsize] != 0 {
			monitor(true)
			stoperr("put", hg)
			if hg.board[x+y*xsize] != order {
				reT = 20
			} else {
				reT = 0
			}
			return
		}

		// rodへの書き込み
		write()

		search := func(rd *rod) (pos int) {

			for pos, s := range rd.diff {
				if hg.board[s] == 0 {
					return pos
				}
			}
			stoperr("search err", hg)
			return
		}

		//　残りがひとつなら予約

		var resparaL []respara

		if hg.row[y].cnt == 1 {
			rst := hg.row[y].rest
			sx := search(&hg.row[y])
			resparaL = append(resparaL, respara{sx, y, rst, &hg.row[y]})
		}

		if hg.columns[x].cnt == 1 {
			rst := hg.columns[x].rest
			sy := search(&hg.columns[x])
			resparaL = append(resparaL, respara{x, sy, rst, &hg.columns[x]})
		}

		if cross {
			if x == y {
				if hg.diagonal[0].cnt == 1 {
					rst := hg.diagonal[0].rest
					sy := search(&hg.diagonal[0])
					resparaL = append(resparaL, respara{sy, sy, rst, &hg.diagonal[0]})
				}
			}

			if x == ysize-y-1 {
				if hg.diagonal[1].cnt == 1 {
					rst := hg.diagonal[1].rest
					sy := search(&hg.diagonal[1])
					resparaL = append(resparaL, respara{ysize - sy - 1, sy, rst, &hg.diagonal[1]})
				}
			}
		}

		for _, rp := range resparaL {

			x := rp.x
			y := rp.y
			rest := rp.rest
			rd := rp.rd

			if hg.board[y*xsize+x] != 0 {

				if hg.board[y*xsize+x] == rest {
					rd.threshold = rest
					//monitor(true)
					continue
				} else {
					//monitor(true)
					reT = 30
					return
				}
			} else {
				if reserve(x, y, rest, rd); reT != 0 {
					return
				}
			}
		}
		//put
	}

	tail := func() {

		if hg.enD.handsCt() != 1 {
			return
		}

		var (
			bp    *node
			fp    *node
			order int

			ar *node
		)

		if ta.bp == nil {
			bp = hg.enD
			fp = hg.enD.neighborList()[0]
			order = size - 1

			bor := hg.board[fp.pos]

			if bor != 0 {
				if bor != order {
					reT = 500
					return
				}
			} else {
				if put(fp.x, fp.y, order); reT != 0 {
					return
				}
			}

		} else {
			if ta.fp.neighborLen() != 2 {
				return
			}

			bp = ta.bp
			fp = ta.fp
			order = ta.order
		}

		for {

			if ar = bp.arm(fp); ar == nil {
				break
			}

			bp = fp
			fp = ar
			order--

			bor := hg.board[fp.pos]

			if bor != 0 {
				if bor != order {
					reT = 501
					return
				} else {
					continue
				}
			}

			if put(fp.x, fp.y, order); reT != 0 {
				return
			}
		}

		nta.bp = bp
		nta.fp = fp
		nta.order = order

		// tail
	}

	// mark main

	nta = ta

	//if tail(); reT != 0 {
	//	return
	//}

	x := next.x
	y := next.y

	// 予約されている場合
	if hg.board[next.pos] != 0 {
		if hg.board[next.pos] != footprint {
			// 予約と違う
			reT = 10
			return

		} else {
			//予約と一致したので当該rodのthreasholdをsizeに変更

			if hg.row[y].cnt == 0 {
				pa.threashold = append(pa.threashold, parath{&hg.row[y], hg.row[y].threshold})
				hg.row[y].threshold = size
			}

			if hg.columns[x].cnt == 0 {
				pa.threashold = append(pa.threashold, parath{&hg.columns[x], hg.columns[x].threshold})
				hg.columns[x].threshold = size
			}

			if cross {

				if x == y {
					if hg.diagonal[0].cnt == 0 {
						pa.threashold = append(pa.threashold, parath{&hg.diagonal[0], hg.diagonal[0].threshold})
						hg.diagonal[0].threshold = size
					}
				}

				if x == ysize-y-1 {
					if hg.diagonal[1].cnt == 0 {
						pa.threashold = append(pa.threashold, parath{&hg.diagonal[1], hg.diagonal[1].threshold})
						hg.diagonal[1].threshold = size
					}
				}
			}
		}
		//return
	} else {

		if put(next.x, next.y, footprint); reT != 0 {
			return
		}
	}

	if tail(); reT != 0 {
		return
	}

	//if !checkth() {
	//	return
	//}

	checkth()

	return
	//mark
}

func unmark(pa param, hg *hGlobal) {

	erase := func(r [3]int) {

		x := r[0]
		y := r[1]
		order := r[2]

		rderase := func(rd *rod) {
			rd.cnt++
			rd.rest += order
		}

		hg.board[x+y*xsize] = 0

		rderase(&hg.row[y])

		rderase(&hg.columns[x])

		if cross {
			if x == y {
				rderase(&hg.diagonal[0])
			}
			if x == ysize-y-1 {
				rderase(&hg.diagonal[1])
			}
		}
	}

	for _, r := range pa.record {
		erase(r)
	}

	for _, pt := range pa.threashold {
		pt.rd.threshold = pt.th
	}

	// unmark
}

func initMagic(th [5]int, hg *hGlobal) {

	total := (size + 1) * size / 2

	makerodL := func(rods int, rodn int, div int) (rodL []rod) {

		sum := total / rodn
		rodL = make([]rod, div)

		for i := 0; i < div; i++ {
			rodL[i].cnt = rods
			rodL[i].rest = sum
			rodL[i].threshold = size
			rodL[i].diff = make([]int, rods)
		}

		return
	}

	line := func(rd *rod, s int, ds int) {
		for i := 0; i < rd.cnt; i++ {
			rd.diff[i] = s
			s += ds
		}
	}

	hg.row = makerodL(xsize, ysize, ysize)
	for i := range ysize {
		line(&hg.row[i], i*xsize, 1)
	}

	hg.columns = makerodL(ysize, xsize, xsize)
	for i := range xsize {
		line(&hg.columns[i], i, xsize)
	}

	if cross {
		hg.diagonal = makerodL(xsize, ysize, 2)
		line(&hg.diagonal[0], 0, xsize+1)
		line(&hg.diagonal[1], xsize*(ysize-1), -xsize+1)
	}

	hg.board = make([]int, size)

	put := func(x int, y int, order int, hg *hGlobal) {

		rdwrite := func(rd *rod) {
			rd.cnt--
			rd.rest -= order
		}

		rdwrite(&hg.row[y])

		rdwrite(&hg.columns[x])

		if cross {
			if x == y {
				rdwrite(&hg.diagonal[0])
			}
			if x == ysize-y-1 {
				rdwrite(&hg.diagonal[1])
			}
		}

		hg.board[x+y*xsize] = order
	}

	put(hg.start%xsize, hg.start/xsize, 1, hg)
	put(hg.enD.x, hg.enD.y, size, hg)

	if threadmode == Multipath {
		for i, pa := range th {
			put(pa%xsize, pa/xsize, i+2, hg)
		}
	}
}

/*
	distance := func(srt *node, des *node) (dis int) {

		//
		if true {
			return
		}
		//

		if srt == des {
			return 0
		}

		var mem [2][]*node

		mem[0] = make([]*node, size)
		mem[1] = make([]*node, size)

		var mind int

		hg.level++
		level := hg.level

		srt.level = level

		if srt == next {
			hg.enD.level = level
			mem[0][0] = srt
			mind = 1
		} else {
			next.level = level
			dis = 1

			for _, sq := range srt.neighborList() {
				if sq == des {
					return
				}
				if sq.level == level {
					continue
				}
				mem[0][mind] = sq
				mind++
				sq.level = level
			}
		}

		sw := 1

		for {

			sw = 1 - sw
			mlist := mem[sw]
			llist := mem[1-sw]

			ret, lind := front(mlist, mind, llist, des, level)

			dis++

			if ret {
				return
			} else {
				if lind == 0 {
					return
				}
			}

			mind = lind
		}

		// distance
		//return
	}
*/
