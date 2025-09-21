package hamiltonpath

// ver 1 2023/10/2
// ver 1.1 2023/10/3
// ver 1.2 2023/10/4
// ver 1.3 2023/10/12
// ver 1.4 2023/10/15
// ver 1.5 2023/11/14 掟破りの解決を残している。
// ver 2.0 2023/11/24 6X6 test pass 6X8 test pass
// ver 3.0 2024/2/7 ver 2.0 の役2倍の効率。
// 2024/7/27 reapList を返り値に変更
// 2024/7/29 loopを残すbugの修正。
// 2024/8/22 成長tubeの反対側nodeへのlinkをreap
//           まだ完成していない。
// 2024/8/31 next-tube の場合のreapによる解決
// 2024/9/20 上記を元に戻す。
// 2025/4/8 next とenDのlinkをOnereapする。
// 2025/4/9 enDsur(enD)
// 2025/4/10 nexTsur()
// 2025/4/13 enDsur nexTsur廃止　next.neighBor(enD)
//2025/4/16 next.opposite()
//2025/4/18 nextの完全スリム化
//2025/4/22 enDの完全スリム化
//2025/4/23 nextの完全スリム化 enDの完全スリム化

func unreaP(reapList twins) {

	for _, sql := range reapList {
		sql[0].link(sql[1])
		sql[1].link(sql[0])
	}
}

func gather(now *node, next *node, hg *hGlobal) (
	deadEnd int, reapList twins) {

	//enD := hg.enD

	var (
		h3pt    *node
		tip     func(*node)
		tube    func(*node)
		reap    func(*node)
		oneReap func(*node, *node)
	)

	tip = func(sq *node) {

		switch sq.neighborLen() {

		case 0:
			deadEnd = 1
			return
		case 1:
			//if sq != next && sq != enD {
			if sq != next {
				deadEnd = 2
				return
			}
		case 2:
			if tube(sq); deadEnd != 0 {
				return
			}
		default:
			// 探索打ち切り。
		}
	}

	tube = func(sq *node) {

		nbL := sq.neighborList()

		for _, sqq := range nbL {

			handsL := sqq.handsList()

			switch len(handsL) {
			case 0:
			case 1:
			case 2:

				if !next.include(handsL) {
					if reap(sqq); deadEnd != 0 {
						return
					}
				}

			case 3:

				// h3pt はひとつのみ。
				if h3pt != nil {
					if h3pt == sqq {
						continue
					} else {
						// ３翼が複数
						deadEnd = 4
						return
					}
				} else {

					if !next.include(handsL) {
						// 3翼を突破できない。
						deadEnd = 5
						return
					}

					// ３翼を登録
					h3pt = sqq

					if reap(sqq); deadEnd != 0 {
						return
					}
				}
			default:
				// ３本以上の腕がある。
				deadEnd = 6
				return
			}
		}

		// loopのチェック
		if sq.closed(nbL[0]) {
			deadEnd = 7
			return
		}
	}

	reap = func(sq *node) {

		var twin []*node

		for _, sqq := range sq.neighborList() {

			if sqq.neighborLen() == 2 {
				continue
			}

			sq.unlink(sqq, hg)
			sqq.unlink(sq, hg)
			reapList = append(reapList, [2]*node{sq, sqq})
			twin = append(twin, sqq)
		}

		for _, sqq := range twin {

			if tip(sqq); deadEnd != 0 {
				return
			}
		}
	}

	oneReap = func(sq1, sq2 *node) {

		sq1.unlink(sq2, hg)
		sq2.unlink(sq1, hg)
		reapList = append(reapList, [2]*node{sq1, sq2})
	}

	// main code

	if next == hg.enD {
		deadEnd = 10
		return
	}

	//
	hg.enD.inc()

	//
	defer hg.enD.dec()

	for _, sq := range now.neighborList() {

		if sq == next {
			continue
		}

		if tip(sq); deadEnd != 0 {
			return
		}
	}

	if next.neighBor(hg.enD) {
		oneReap(next, hg.enD)
		if hg.enD.neighborLen() == 1 {
			deadEnd = 11
			return
		}
	}

	hl := next.handsList()

	switch len(hl) {
	case 0:
	case 1:
		for _, sq := range next.neighborList() {

			if sq == hl[0] {
				continue
			}

			oneReap(next, sq)

			if tip(sq); deadEnd != 0 {
				return
			}
		}
	default:
		deadEnd = 12
		return
	}

	hl = hg.enD.handsList()

	switch len(hl) {
	case 0:
	case 1:
		if hg.enD.neighborLen() != 2 {
			for _, sq := range hg.enD.neighborList() {

				if sq == hl[0] {
					continue
				}

				oneReap(hg.enD, sq)

				if tip(sq); deadEnd != 0 {
					return
				}
			}
		}
	default:
		deadEnd = 13
		return
	}

	if h3pt != nil {

		oneReap(h3pt, next)

		if tube(h3pt); deadEnd != 0 {
			return
		}

	}

	//if hg.enD.handsCt() > 1 {
	//	deadEnd = 20
	//	return
	//drawD(2, "", now, next, reapList, hg)
	//stoperr("enD hands", hg)
	//}

	// enD回りで未発見の次数ゼロが残る
	if next.neighborLen() == 0 {
		deadEnd = 21
		return
	}

	/*
		if next.handsCt() > 1 {
			//drawD(2, "", now, next, reapList, hg)
			stoperr("next hands", hg)
		}

		if hg.enD.neighborLen() == 1 {
			//drawD(2, "", now, next, reapList, hg)
			stoperr("end only", hg)
		}

		if next.neighBor(hg.enD) {
			//drawD(2, "", now, next, reapList, hg)
			stoperr("end neighbor next", hg)
		}
	*/

	return
	// gather
}
