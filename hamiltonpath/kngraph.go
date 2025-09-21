package hamiltonpath

import "math/rand"

func kngraphInit(hg *hGlobal) {

	var nodeInit func(*node, int, int, *rand.Rand)

	graph := make([]node, size)

	var dir = [8][2]int{{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}
	nodeInit = func(s *node, x int, y int, r *rand.Rand) {

		s.x = x
		s.y = y
		s.pos = y*xsize + x

		s.unpassed = true
		//s.order = -1
		s.neighborCnt = 0

		for _, d := range dir {

			xx := x + d[0]
			yy := y + d[1]

			if xx >= 0 && xx < xsize && yy >= 0 && yy < ysize {

				s.neighbor = append(s.neighbor, &graph[yy*xsize+xx])
				s.neighborLinked = append(s.neighborLinked, true)
				s.neighborHist = append(s.neighborHist, -1)
				s.neighborCnt += 1
			}
		}

		if hg.seed != 0 {

			//rand.Shuffle(len(s.neighbor), func(i, j int) {
			r.Shuffle(len(s.neighbor), func(i, j int) {
				s.neighbor[i], s.neighbor[j] = s.neighbor[j], s.neighbor[i]
			})
		}
	}

	pos := 0
	for x := 0; x < xsize; x++ {
		for y := 0; y < ysize; y++ {
			graph[pos] = node{}
			//sequence[pos] = -1
			pos++
		}
	}

	r := rand.New(rand.NewSource(int64(hg.seed)))

	//rand.Seed(int64(hg.seed))

	//rand.New(rand.NewSource(seed))

	for x := 0; x < xsize; x++ {
		for y := 0; y < ysize; y++ {
			pos := y*xsize + x
			nodeInit(&graph[pos], x, y, r)
		}
	}

	hg.graph = graph

}

func knalign(filename string, hg *hGlobal) (start int) {

	if threadmode != Single {
		stoperr("kngraphTrance Not Single Mode", hg)
	}

	sequence := inSequence(filename, hg)

	start = sequence[0]
	hg.start = start
	hg.end = sequence[size-1]
	hg.enD = &hg.graph[hg.end]

	solutions = 1
	steplimit = int64(size)

	var (
		b int
		f int
		//sq *node
		sb *node
		sf *node
		sa *node
	)

	b = sequence[0]

	for i := 1; i < len(sequence); i++ {

		f = sequence[i]

		sb = &hg.graph[b]
		sf = &hg.graph[f]

		for j, sqq := range sb.neighbor {

			if sqq == sf {
				sa = sb.neighbor[0]
				sb.neighbor[0] = sf
				sb.neighbor[j] = sa
				goto label
			}
		}
		stoperr("kngraphTrance can't exchange", hg)
	label:
		b = f
	}

	return
}

/*
var dir = [8][2]int{{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}

//var dir = [8][2]int{{2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}, {1, 2}, {2, 1}}

func (s *node) nodeInit(x int, y int) {

	s.x = x
	s.y = y
	s.pos = y*xsize + x

	s.unpassed = true
	//s.order = -1
	s.neighborCnt = 0

	for _, d := range dir {

		xx := x + d[0]
		yy := y + d[1]

		if xx >= 0 && xx < xsize && yy >= 0 && yy < ysize {

			s.neighbor = append(s.neighbor, &graph[yy*xsize+xx])
			s.neighborLinked = append(s.neighborLinked, true)
			s.neighborHist = append(s.neighborHist, -1)
			s.neighborCnt += 1
		}
	}

	if seed != 0 {

		rand.Shuffle(len(s.neighbor), func(i, j int) {
			s.neighbor[i], s.neighbor[j] = s.neighbor[j], s.neighbor[i]
		})
	}
}
*/
