package hamiltonpath

type twins [][2]*node

func stepN(now *node, hg *hGlobal) {

	if hg.steplimitCt++; hg.steplimitCt > steplimit {
		stopovr(hg)
	}

	hg.stepNum++
	now.leaVe(hg)
	hg.sequence[hg.stepNum] = now.pos
	now.unpassed = false

	defer func() {
		now.unpassed = true
		hg.sequence[hg.stepNum] = -1
		now.baCk()
		hg.stepNum--
	}()

	if now == hg.enD {
		if hg.stepNum == size-1 {
			if hg.solutionCt++; hg.solutionCt == solutions {
				stopfound(hg)
			}
		}
		return
	}

	nbL := now.neighborList()

	if len(nbL) == 1 && nbL[0].neighborLen() < 2 {
		stepN(nbL[0], hg)
		return
	}

	for _, next := range nbL {

		deadEnd, reapList := gather(now, next, hg)

		if deadEnd == 0 {
			stepN(next, hg)
		}

		unreaP(reapList)
	}
}

// func stepE(now *node, area []*node, hg *hGlobal) (gMark bool) {
func stepE(now *node, ps pset, hg *hGlobal) (gMark bool) {
	if hg.steplimitCt++; hg.steplimitCt > steplimit {
		//drawD(2, "test", now, nil, nil, hg)
		stopovr(hg)
	}

	hg.stepNum++
	now.leaVe(hg)
	hg.sequence[hg.stepNum] = now.pos
	now.unpassed = false

	defer func() {
		now.unpassed = true
		hg.sequence[hg.stepNum] = -1
		now.baCk()
		hg.stepNum--
	}()

	if now == hg.enD {
		if hg.stepNum == size-1 {
			if hg.solutionCt++; hg.solutionCt == solutions {
				stopfound(hg)
			} else {
				gMark = true
			}
		} else {
			//drawD(2, "test", now, nil, nil, hg)
			//stopunresolve("incomplete", hg)
		}
		return
	}

	nbL := now.neighborList()

	if len(nbL) == 1 && nbL[0].neighborLen() < 2 {
		//gMark = stepE(nbL[0], area, hg)
		gMark = stepE(nbL[0], ps, hg)
		return
	}

	for _, next := range nbL {

		deadEnd, reapList := gather(now, next, hg)

		if deadEnd == 0 {

			mode, newps, er := scan(now, next, ps, reapList, hg)
			switch mode {
			case forWard:
				gMark = stepE(next, newps, hg)
			case conTinue:
				if er == 13 || er == 23 {
					//	unreaP(reapList)

					//	drawD(2, "zero0", now, next, nil, hg)
					//	continue

				}
			case exiT:
				gMark = stepE(next, newps, hg)
				//drawD(2, "test", now, next, reapList, hg)
				unexit(newps, hg)

				if hg.phase == 0 {
					if !gMark {

						if catch(newps, hg, now, next, reapList) {
							//drawD(1, "test", now, next, reapList, hg)
							unreaP(reapList)
							return
						}
					}
				}

			case inTo:
				gMark = stepE(next, newps, hg)
				uninto(newps, hg)

			case knoT:
				gMark = stepE(next, newps, hg)
				unknot(newps, hg)

				//case reT:
			//	unreaP(reapList)
			//	return
			//case cRoss:

			//	gMark = stepE(next, newps, hg)

			/*
				if hg.phase == 0 {
					if !catch(now, next, newps, gMark, reapList, hg) {
						unreaP(reapList)
						return
					}
				}
			*/

			case baCk:
				unreaP(reapList)
				return
			}
		}

		unreaP(reapList)

		if hg.phase != 0 {
			stop := bottom(hg)

			if stop {
				//drawD(2, "test", now, next, reapList, hg)
				continue
			} else {
				return
			}
		}
	}

	return
}

func stepM(now *node, ta paratail, hg *hGlobal) {

	if hg.steplimitCt++; hg.steplimitCt > steplimit {
		mPrint(hg)
		drawD(3, "", now, nil, nil, hg)
		stopovr(hg)
	}

	hg.stepNum++
	now.leaVe(hg)
	hg.sequence[hg.stepNum] = now.pos
	now.unpassed = false

	defer func() {
		now.unpassed = true
		hg.sequence[hg.stepNum] = -1
		now.baCk()
		hg.stepNum--
	}()

	if now == hg.enD {
		if hg.stepNum == size-1 {
			//drawD(2, "", now, nil, nil, hg)
			mPrint(hg)
			mFile(hg)
			if hg.solutionCt++; hg.solutionCt == solutions {
				stopfound(hg)
			}
		}
		return
	}

	nbL := now.neighborList()

	//if len(nbL) == 1 {
	if len(nbL) == 1 && nbL[0].neighborLen() < 2 {
		reT, pa, nta := mark(now, nbL[0], ta, hg, nil)
		if reT == 0 {
			stepM(nbL[0], nta, hg)
		}
		unmark(pa, hg)
		return
	}

	for _, next := range nbL {

		deadEnd, reapList := gather(now, next, hg)

		if deadEnd == 0 {

			reT, pa, nta := mark(now, next, ta, hg, reapList)
			if reT == 0 {
				stepM(next, nta, hg)
			}
			unmark(pa, hg)
		}
		unreaP(reapList)
	}
}
