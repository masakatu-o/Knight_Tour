package hamiltonpath

// 2024/9/1  zero, oddone 段階までほぼ完成
//      22X22　seed = 16,29 now is not in odd segment を残している。
// 202/9/3 now-tube の前処理。baCkはisolationにまとめた。
// 2024/9/5 oddsegment-one を区別しない。baCkはsegment-zeroのときだけ。
// 		indefinite 導入
// 2024/9/7 enD-two 問題解決？ 孤島の処理を残すのみ。
// 2024/9/11 まだまだ。通過可能問題発生。
// 2024/9/14 isolation return 時のbaCk処理
// 2024/9/30 next回りend回りlen(next) == 1 でのstNum設定
// 2024/10/9 34X34 s=30 5095224 抜け、孤島、あり。
// 2024/10/16 34X34 s=30 消える。
// 2024/10/17 next回り,enD回り,errorなし。
// 2024/10/19 next回り,enD回り,enDsect errorなし。
// 2024/10/20 maxN 修正
// 2024/10/20 enDsegment　対応以前。
// 2024/10/26 branch segment Check 詳細バージョン
// 2024/12/8 34X34 s=30 54407209

type (
	seCt int
)

type seG struct {
	brdCnt  int        // 島の橋の数
	brdList [][2]*node // 橋のこちら側の端点とあちら側の端点
	sample  *node      // 島の代表
}

// mode
const (
	forWard  = 1 + iota // 変化なし、次のstepへ
	conTinue            // 問題あり、次の候補へ
	exiT                // 親の島から出る
	inTo                // 分裂はあったが、親の島内にとどまる
	knoT                // 結び目ができた。
	baCk                // 指定のretPまで戻る
)

type pset struct {
	stp  int   // isstackの書き込みポインター
	idp  int   // psでは現在の島、newpsでは次の島
	gate *node // 一つ島に蓋をした。通常はnil　exiT inTo　knoT
}
type island struct {
	// isstackの要素
	area    []*node    // 島に含まれるノードの集合
	brdCnt  int        // 島の橋の数。渡橋時に再カウントされる。
	brdList [][2]*node // 橋のこちら側の端点とあちら側の端点
	parent  seCt       // 親の島
	birth   int        // 島の誕生stepNum
}

func scan(now *node, next *node, ps pset, reapList twins, hg *hGlobal) (
	mode int, newps pset, er int) {
	// 現在いる島での、nextへ移る場合の状態変化をscanして、mode,newps,erを返す。

	enD := hg.enD
	level := hg.level

	stp := ps.stp
	selfIdp := ps.idp

	defer func() {

		hg.level = level

		// 出口調査
		if er != 0 {
			if ps.stp == newps.stp {
			}
			//drawD(2, "test", now, next, reapList, hg)
			//d := 0
			//d++
		}
		//

		//check

		if newps.idp >= stp {
			//drawD(2, "test", now, next, reapList, hg)
			stopunresolve("over idp", hg)
		}

		var (
		//live []seCt
		)

		stp := newps.stp
		table := make([]bool, stp)
		stack := hg.isstack

		if mode == 0 {
		}

		if ps.idp > 0 {
			for i := range stp {
				table[stack[i].parent] = true
			}

			for i, flag := range table {
				if flag {
					if i == newps.idp {
						//drawD(2, "test", now, next, reapList, hg)
						stopunresolve("unsuitable idp", hg)
					}
					continue
				}
				//live = append(live, seCt(i))
			}
		}

		/*
			cnt := 0

			for _, sect := range live {
				cnt += stack[sect].brdCnt
			}

			if cnt%2 != 0 {
				//drawD(2, "test", now, next, reapList, hg)
				d := 0
				d++
				//stopunresolve("unsuitable total cnt", hg)
			}
		*/

		//check

	}()

	var (
		term         func(*node, *node) (*node, *node)
		segmentation func([]*node, seCt) map[seCt]*seG
		decision     func(map[seCt]*seG)
	)

	// br 入り口
	// fr 最初のtube
	// brp 最後のtube
	// frp 出口
	term = func(br, fr *node) (*node, *node) {

		var (
			brp *node = br
			frp *node = fr
			ar  *node
		)

		for {

			if ar = brp.arm(frp); ar == nil {
				break
			}

			brp = frp
			frp = ar
		}

		return brp, frp
		//term
	}

	// グラフを島と橋に文節する。
	segmentation = func(area []*node, secnum seCt) map[seCt]*seG {

		var (
			segment    = make(map[seCt]*seG)
			bridgeList [][2]*node
			infiltrate func(*node)
		)

		infiltrate = func(sq *node) {
			//橋のたもとの候補を登録

			sq.level = level
			sq.sect = secnum

			for _, sqq := range sq.neighborList() {

				if sqq.level == level {
					continue
				}

				switch sqq.neighborLen() {
				case 0, 1:
				case 2:
					br, fr := term(sq, sqq)
					br.level = level
					//橋のたもとの候補
					bridgeList = append(bridgeList, [2]*node{sq, fr})
				default:
					infiltrate(sqq)
				}
			}
			//infiltrate
		}

		// segmenntation main

		// 使い捨てのmark
		level++

		// enDはどの島の要素でもない。
		enD.level = level

		for _, sq := range area {
			if sq.unpassed && sq.level != level && sq.neighborLen() > 2 {

				infiltrate(sq)
				isl := &seG{}
				isl.sample = sq
				segment[secnum] = isl
				secnum++
			}
		}

		// 島と島をつなぐ橋だけを抽出し、こちら側とあちら側のnode、島ごとの橋の数を記録する
		for _, four := range bridgeList {

			fr0 := four[0]
			fr1 := four[1]

			sl0 := fr0.sect
			sl1 := fr1.sect

			if sl0 != sl1 {

				if fr1 == enD || fr1 == next {
					continue
				}

				isl0 := segment[sl0]
				isl0.brdCnt++
				isl0.brdList = append(isl0.brdList, [2]*node{fr0, fr1})
				segment[sl0] = isl0

				if fr1.level == level {
					isl1 := segment[sl1]
					isl1.brdCnt++
					isl1.brdList = append(isl1.brdList, [2]*node{fr1, fr0})
					segment[sl1] = isl1
				}
			}
		}

		return segment
		//segmenntation
	}

	decision = func(segment map[seCt]*seG) {
		// nowからnextへの移動が有効か吟味する。

		var (
			nextT   *node //　次はどこへ行こうとしているか。
			zerocnt int
			oddcnt  int
		)

		// nexTは、next自身か、渡ろうとする橋の向こう。
		if next.neighborLen() == 1 {
			_, nextT = term(next, next.neighborList()[0])
		} else {
			nextT = next
		}

		//　新規の島の情報（seG構造体）をstackに積む。
		push := func(segment map[seCt]*seG) {

			makenewArea := func(areaP *node) (area []*node) {
				//新しい島のメンバー

				var infiltrate func(*node)

				level++
				enD.level = level

				infiltrate = func(sq *node) {

					area = append(area, sq)
					sq.level = level

					for _, sqq := range sq.neighborList() {

						if sqq.unpassed && sqq.level != level && sqq.neighborLen() > 2 {

							infiltrate(sqq)
						}
					}
				}

				infiltrate(areaP)

				return
				// makenewArea
			}

			if len(segment) == 0 {
				// 巡歴を完了した島のコピー　debug後は必要ない。
				hg.isstack[stp].brdCnt = 0
				hg.isstack[stp].parent = seCt(selfIdp)
				hg.isstack[stp].birth = hg.stepNum
				stp++
			} else {
				//分割した島あるいは通過して残された島あるいはknotで橋の数が変化した島
				for sect, seg := range segment {

					hg.isstack[sect].area = makenewArea(seg.sample)
					hg.isstack[sect].brdList = seg.brdList
					hg.isstack[sect].brdCnt = seg.brdCnt
					hg.isstack[sect].parent = seCt(selfIdp)
					hg.isstack[sect].birth = hg.stepNum
				}

				stp += len(segment)
			}
		}

		// enDはsectの島に触れているか
		touchE := func(sect seCt) bool {

			endn := enD.neighborList()

			if len(endn) == 1 {
				_, eterm := term(enD, endn[0])
				return eterm.sect == sect
			} else {
				for _, sq := range endn {

					hl := sq.handsList()
					switch len(hl) {
					case 0:
						if sq.sect == sect {
							return true
						}
					case 1:
						_, fr := term(sq, hl[0])
						if fr.sect == sect {
							return true
						}
					case 2:
						//drawD(2, "test", now, next, reapList, hg)
						for i, sqq := range hl {
							if sqq != enD {
								_, fr := term(sq, hl[i])
								if fr.sect == sect {
									return true
								}
							}
						}
					default:
						stoperr("touchE", hg)
					}
				}

				return false
			}
		}

		exit := func() {
			// 渡橋。行く先の島のbridgeCnt再設定、gateの設定
			sect := nextT.sect

			area := hg.isstack[sect].area

			seg := segmentation(area, sect)

			if len(seg) != 1 {
				// nextT.sectが正しくない？
				// 島がknotで分裂した？
				//　4つのノードの島を通過した？
				//drawD(2, "test", now, next, reapList, hg)
				stopunresolve("zero seg", hg)
			}

			//修正されたbrdCnt
			// gMarkがtureのときはcntも元に戻すこと。
			hg.isstack[sect].brdCnt = seg[sect].brdCnt

			if hg.isstack[sect].brdCnt == 1 {
				gt := seg[sect].brdList[0][1]
				newps.gate = gt
				gt.gate = true
			}

			// exit
		}

		into := func() {
			if segment[nextT.sect].brdCnt == 1 {
				//drawD(2, "test", now, next, reapList, hg)
				gt := segment[nextT.sect].brdList[0][1]
				newps.gate = gt
				gt.gate = true
			}
			// into
		}

		knot := func() {
			//gateの処理
			if segment[seCt(stp)].brdCnt == 1 {
				//drawD(2, "test", now, next, reapList, hg)
				gt := segment[seCt(stp)].brdList[0][1]
				newps.gate = gt
				gt.gate = true
			}
			// knot
		}

		knotcheck := func() bool {
			// 外側の島が孤立していないか
			var (
				ls    []seCt
				fr    *node
				sect0 seCt
				sect1 seCt
			)

			brl := hg.isstack[selfIdp].brdList

			for _, tw := range brl {

				if tw[0].unpassed && tw[0].neighborLen() == 2 {

					nb := tw[0].neighborList()

					_, fr = term(tw[0], nb[0])
					sect0 = fr.sect

					_, fr = term(tw[0], nb[1])
					sect1 = fr.sect

					if sect0 == sect1 {
						ls = append(ls, sect0)
					}
				}
			}

			sectN := nextT.sect

			if sectN == 0 {
			}

			if len(ls) != 0 {
				ls0 := ls[0]
				ls1 := ls[1]
				if ls0 == ls1 {
				}

				if hg.isstack[ls0].brdCnt-2 == 0 {
					//drawD(2, "test", now, next, reapList, hg)
					//d := 0
					//d++
					return false
				}

			} else {
				return true
				//drawD(2, "test", now, next, reapList, hg)
				//d := 0
				//d++

			}
			return true
			//knotcheck
		}

		oddzerocnt := func(segment map[seCt]*seG) {

			for _, isl := range segment {

				if isl.brdCnt == 0 {
					zerocnt++
				}

				if isl.brdCnt%2 == 1 {
					oddcnt++
				}
			}
			// oddzerocnt
		}

		//decision main

		//
		if ps == newps {
		}
		//

		selfCnt := hg.isstack[selfIdp].brdCnt

		if nextT.level != level {
			//親の島を出ようとしている。
			selfCnt--
		}

		switch len(segment) {
		case 0:
			//　この島の巡歴を終えた
			if selfCnt > 1 {
				//drawD(2, "test", now, next, reapList, hg)
				if !knotcheck() {
					mode = conTinue
					er = 1
					return
				}
			}

			if nextT == enD {
				mode = forWard
				return
			} else {
				mode = exiT
				exit()
			}

		case 1:
			// 分裂はなかった

			selfNCnt := segment[seCt(stp)].brdCnt

			if selfCnt-selfNCnt > 1 {
				//橋の数に1以上の変化があった。

				if selfNCnt == 0 {
					//　ゼロ島が残った
					//　内部状態がゼロ
					mode = conTinue
					er = 11
					return
				}

				if !knotcheck() {
					// 外部の島にゼロが生じた
					mode = conTinue
					er = 12
					return
				}
			}

			if nextT.level == level {
				//島の内へ

				if nextT == enD {
					mode = conTinue
					er = 14
					return
				}

				if selfCnt-selfNCnt != 0 {
					mode = knoT
					knot()
				} else {
					mode = forWard
					return
				}

			} else {
				//島の外へ
				if nextT.gate {
					//drawD(2, "test", now, next, reapList, hg)
					mode = conTinue
					er = 15
					return
				}

				mode = exiT
				exit()
				//}
			}

		//　島が分割された。
		default:

			selfNCnt := 0

			for _, seg := range segment {

				for _, tw := range seg.brdList {
					if tw[1].sect >= seCt(ps.stp) {
						continue
					}

					selfNCnt++
				}
			}

			if selfCnt-selfNCnt > 1 {
				//橋の数に1以上の変化があった。

				if selfNCnt == 0 {
					//　ゼロ島が残った
					//　内部状態がゼロ
					mode = conTinue
					er = 11
					return
				}

				if !knotcheck() {
					// 外部の島にゼロが生じた
					mode = conTinue
					er = 12
					return
				}
			}

			oddzerocnt(segment)

			if zerocnt > 0 {
				//　ゼロ島が残った
				//　分裂による内部状態
				//drawD(2, "test", now, next, reapList, hg)
				//stopunresolve("zero island", hg)

				if zerocnt > 2 {
					//drawD(2, "test", now, next, reapList, hg)
					stopunresolve("zerocnt > 2", hg)
				}

				var (
					sect seCt
				)

				for s, isl := range segment {

					if isl.brdCnt == 0 {
						sect = s
					}
				}

				max1 := 0
				max2 := 0

				for _, sq := range hg.isstack[seCt(ps.idp)].area {

					if sq.unpassed && sq.sect == sect {
						maxH := sq.maxHist()

						switch {
						case maxH >= max1:
							max1 = maxH
						case maxH > max2:
							max2 = maxH

						}
					}
				}

				if max2 == 0 {
					mode = conTinue
					er = 21
					return
					//drawD(2, "test", now, next, reapList, hg)
					//stopunresolve("max2 ==0", hg)
				}

				hg.phase = 2
				hg.retP = max2
				mode = baCk
				return

			}

			// 親島の奇偶
			//peven := hg.isstack[selfIdp].brdCnt%2 == 0
			peven := selfNCnt%2 == 0

			if nextT.level == level {
				// 親島の中

				if nextT == enD {
					mode = conTinue
					er = 24
					return
				}

				if peven {

					switch oddcnt {
					case 0:
						// 同居
						if !touchE(nextT.sect) {
							mode = conTinue
							er = 25
							return
						}
					case 2:
						// 別居

						var (
							nflag bool
							eflag bool
						)

						for sect, seg := range segment {

							if seg.brdCnt%2 == 0 {
								continue
							}

							if nextT.sect == sect {
								nflag = true
							} else {
								if touchE(sect) {
									eflag = true
								}
							}
						}

						if !nflag || !eflag {
							mode = conTinue
							er = 26
							return
						}

					default:
						mode = conTinue
						er = 27
						return
						// 4,6,,,
					}
				} else {
					switch oddcnt {
					case 1:
						// nextTは奇数島か？
						if segment[nextT.sect].brdCnt%2 == 0 {
							mode = conTinue
							er = 28
							return
						}
					default:
						mode = conTinue
						er = 29
						return
						// 3,5,,,
					}
				}
				//drawD(2, "test", now, next, reapList, hg)
				mode = inTo
				into()

			} else {
				// 親島から出る。

				if nextT.gate {
					mode = conTinue
					er = 31
					return
				}

				if peven {
					switch oddcnt {
					case 1:
						// end は奇数島に触れている？
						for sect, ist := range segment {
							if ist.brdCnt%2 == 1 {
								if !touchE(sect) {
									mode = conTinue
									er = 32
									return
								}
								break
							}
						}
					default:
						mode = conTinue
						er = 33
						return
						// 3,5,..
					}
				} else {
					if oddcnt != 0 {
						mode = conTinue
						er = 34
						return
					}
				}

				mode = exiT
				exit()
			}
		}

		// 分割あるいは残された島をスタックに積む

		push(segment)
		newps.stp = stp
		if mode == knoT {
			newps.idp = stp - 1
		} else {
			newps.idp = int(nextT.sect)
		}

		// decision
	}

	// scan main

	newps.idp = ps.idp
	newps.stp = ps.stp
	newps.gate = nil

	if next.neighborLen() == 2 {
		mode = forWard
		return
	}

	segment := segmentation(hg.isstack[ps.idp].area, seCt(ps.stp))

	decision(segment)

	return
	// scan
}

func catch(newps pset, hg *hGlobal, now *node, next *node, reapList twins) bool {
	if newps.idp == 0 {
		return false
	}

	hg.retP = hg.isstack[newps.idp].birth
	hg.phase = 1
	//drawD(2, "test", now, next, reapList, hg)

	return true
}

func bottom(hg *hGlobal) (stop bool) {

	if hg.stepNum > hg.retP {
		return
	}

	hg.phase = 0
	stop = true

	return
}

func ungate(newps pset, hg *hGlobal) {
	if newps.gate != nil {

		if newps.gate == hg.enD {
			d := 0
			d++
		}

		//if newps.gatecnt--; newps.gatecnt == 0 {
		newps.gate.gate = false
		newps.gate = nil
		//} else {
		//stoperr("gateの重なり処理OK", hg)
		//}
	}
}

func unexit(newps pset, hg *hGlobal) {

	sect := newps.idp

	//exit() の後始末
	//hg.isstack[sect].brdCnt++

	// 橋のこちら側のsectを島のsectに戻す。
	br := hg.isstack[sect].brdList

	for _, twin := range br {

		twin[0].sect = seCt(sect)
	}

	ungate(newps, hg)
}

func uninto(newps pset, hg *hGlobal) {
	ungate(newps, hg)
}

func unknot(newps pset, hg *hGlobal) {
	ungate(newps, hg)
}

/*
	// cocoon が出来たとき、endP, chaiN の設定。
	//chainMark = func(sect seCt) (chain *node) {
	chainMark = func(sect seCt) {

		//drawSL(2, segment, now, next, area, reapList, "")

		brdL := segment[sect].brdList

		if sect == enD.sect {
			// test
			if len(brdL) != 1 {
				drawD(2, "", now, next, reapList, hg)
				stoperr("chainMark 1", hg)
			}
			//
			chain := brdL[0][0]
			chain.level = endP
			return
		}

		// test
		if len(brdL) != 2 {
			drawD(2, "", now, next, reapList, hg)
			stoperr("chainMark 2", hg)
		}
		//

		for i, twin := range brdL {

			if twin[1].level == endP {
				chain := brdL[1-i][0]
				chain.level = endP
				//return brdL[1-i][0]
				return
			}
		}

		brdL[0][0].level = chaiN
		brdL[1][0].level = chaiN
	}

	//
	isInSect = func() (sect seCt) {

		// 一色ならばその色を、多色ならば0を返す。
		surround := func(w *node) (sect seCt) {

			var s seCt

			for _, sq := range w.neighborList() {

				if sq.level != level {
					continue
				}

				if sq.neighborLen() == 2 {

					for _, sqq := range sq.neighborList() {

						if sqq.neighborLen() != 2 {
							s = sqq.sect
						}
					}

				} else {

					s = sq.sect
				}

				if sect == 0 {
					sect = s
				} else {
					if sect != s {
						return 0
					}
				}
			}

			return
		}

		// 接触する島がないので、直前の結果を返す。
		if now.neighborLen() == 1 {
			return surround(&graph[sequence[stepNum-1]])
		}

		if sect = surround(now); sect != 0 {
			return
		}

		sn := stepNum - 1

		var ww *node

		for {

			ww = &graph[sequence[sn]]

			if ww.neighborLen() > 1 {
				break
			}

			sn--
		}

		return surround(ww)
	}

	maxI = func(sect seCt) (max int, st bool) {

		var maxx int

		for _, sq := range graph {

			if sq.unpassed && sq.level != level || (sq.level == level && sq.sect != sect) {

				maxx = sq.maxHist()

				if maxx == stepNum {
					st = true
					continue
				}

				if maxx > max {
					max = maxx
				}
			}
		}
		return
	}

	maxO = func(sect seCt) (max int, st bool) {

		var maxx int

		for _, sq := range area {

			if sq.level == level && sq.sect == sect {

				maxx = sq.maxHist()

				if maxx == stepNum {
					st = true
					continue
				}

				if maxx > max {
					max = maxx
				}
			}
		}

		return
	}
*/
