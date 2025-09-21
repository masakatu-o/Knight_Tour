package hamiltonpath

type node struct {
	sect seCt
	//level    int
	unpassed    bool
	neighborCnt int
	gate        bool

	x   int
	y   int
	pos int

	neighbor       []*node
	neighborLinked []bool
	neighborHist   []int
	//neighborCnt    int
	level int
	mark  int
	//gate bool
}

func (s *node) unlink(p *node, hg *hGlobal) {

	for i, sq := range s.neighbor {

		if sq == p && s.neighborLinked[i] {

			s.neighborLinked[i] = false
			s.neighborHist[i] = hg.stepNum
			s.neighborCnt--
			return
		}
	}

	//stoperr("unlink error", hg)
}

func (s *node) link(p *node) {

	for i, sq := range s.neighbor {

		if sq == p && !s.neighborLinked[i] {

			s.neighborLinked[i] = true
			s.neighborHist[i] = -1
			s.neighborCnt++
			return
		}
	}

	//stoperr("link error", hg)
}

func (s *node) neighborList() []*node {

	ret := make([]*node, s.neighborCnt)
	j := 0

	for i, flag := range s.neighborLinked {

		if flag {
			ret[j] = s.neighbor[i]
			j++
		}
	}
	return ret[:j]
}

func (s *node) neighborLen() int {

	return s.neighborCnt
}

func (s *node) leaVe(hg *hGlobal) {

	for _, sq := range s.neighborList() {

		sq.unlink(s, hg)
	}
}

func (s *node) baCk() {

	for _, sq := range s.neighborList() {

		sq.link(s)
	}
}

func (s *node) handsCt() (hCt int) {

	for _, sq := range s.neighborList() {

		if sq.neighborCnt == 2 {
			hCt++
		}
	}
	return
}

func (s *node) handsList() (ret []*node) {

	for _, sq := range s.neighborList() {

		if sq.neighborCnt == 2 {
			ret = append(ret, sq)
		}
	}
	return
}

func (s *node) maxHist() (rmax int) {

	//rmax = max

	for _, hist := range s.neighborHist {

		if hist > rmax {
			rmax = hist
		}
	}

	return
}

func (s *node) include(sl []*node) bool {

	for _, sq := range sl {
		if s == sq {
			return true
		}
	}

	return false
}

func (bs *node) arm(fs *node) *node {

	if fs.neighborCnt != 2 {
		return nil
	}

	for i, flag := range fs.neighborLinked {

		if flag && fs.neighbor[i] != bs {
			return fs.neighbor[i]
		}
	}

	return nil
}

func (bs *node) closed(fs *node) bool {

	bsp := bs
	fsp := fs

	for {

		ar := bsp.arm(fsp)

		if ar != nil {

			if ar == bs {
				return true
			}

		} else {
			return fsp == bs
		}

		bsp = fsp
		fsp = ar
	}
}

func (bs *node) inc() {
	bs.neighborCnt++
}

func (bs *node) dec() {
	bs.neighborCnt--
}

func (sq1 *node) neighBor(sq2 *node) bool {

	for i, sq := range sq1.neighbor {

		if !sq1.neighborLinked[i] {
			continue
		}

		if sq == sq2 {
			return true
		}
	}

	return false
}

func front(mlist []*node, mind int, llist []*node, des *node, level int) (
	ret bool, lind int) {

	//var hcnt int

	var tind int

	for ind := range mind {

		sq := mlist[ind]

		tind = lind

		for i, sqq := range sq.neighbor {

			if !sq.neighborLinked[i] || sqq.level == level {
				continue
			}

			if sqq.neighborCnt == 2 {
				lind = tind
				llist[tind] = sqq
				lind++
				break
			}

			llist[lind] = sqq
			lind++
		}

		for i := tind; i < lind; i++ {
			sqq := llist[i]
			if sqq == des {
				ret = true
				return
			}
			sqq.level = level
		}
	}

	return
}

/*
type addList struct {
	size int
	i    int
	list []*node
}

func mkaddList(size int) *addList {

	var al addList
	al.size = size
	al.i = 0
	al.list = make([]*node, size)
	return &al
}

func (al *addList) add(sq *node) {

	if al.i == size {
		stop("in addList add")
	}

	al.list[al.i] = sq
	al.i++
}
*/
