package hamiltonpath

import (
	"fmt"
	"os"
	"time"

	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"

	"gonum.org/v1/plot/plotutil"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func line(px plotter.XYs) (pl *plotter.Line) {

	var err error
	pl, err = plotter.NewLine(px)
	if err != nil {
		panic(err)
	}

	pl.LineStyle.Width = vg.Points(1)

	return
}

func scatter(px plotter.XYs, adj float64) (ps *plotter.Scatter) {

	var err error
	ps, err = plotter.NewScatter(px)
	if err != nil {
		panic(err)
	}

	ps.GlyphStyle.Radius = vg.Points(6 + adj)

	return
}

func pttoXYs(ptL []int) plotter.XYs {

	pts := make(plotter.XYs, len(ptL))

	for i, pos := range ptL {

		pts[i].X = float64(pos%xsize) + 0.5
		//pts[i].Y = float64(ysize-pos/xsize) + 0.5
		pts[i].Y = float64(ysize-pos/xsize) - .5
	}
	return pts
}

func initPlot(title string) (p *plot.Plot) {

	p = plot.New()

	p.HideAxes()

	var recSize float64

	if xsize >= ysize {
		recSize = float64(xsize)
	} else {
		recSize = float64(ysize)
	}

	boxL := make(plotter.XYs, 5)

	boxL[0].X = 0
	boxL[0].Y = 0
	boxL[1].X = float64(xsize)
	boxL[1].Y = 0
	boxL[2].X = float64(xsize)
	boxL[2].Y = float64(ysize)
	boxL[3].X = 0
	boxL[3].Y = float64(ysize)
	boxL[4].X = 0
	boxL[4].Y = 0

	q := line(boxL)
	q.LineStyle.Color = color.RGBA{R: 255, A: 255}
	p.Add(q)

	sqL := make(plotter.XYs, 1)

	sqL[0].X = recSize
	sqL[0].Y = recSize

	s := scatter(sqL, 0)
	s.GlyphStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Add(s)

	p.Title.Text = title

	return
}

func setpassed(p *plot.Plot, hg *hGlobal) {

	sequence := hg.sequence
	stepNum := hg.stepNum

	q := line(pttoXYs(sequence[:stepNum+1]))
	q.LineStyle.Color = color.RGBA{R: 255, A: 255}
	p.Add(q)

	n := scatter(pttoXYs([]int{sequence[stepNum]}), 0)
	n.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	n.Shape = draw.BoxGlyph{}
	p.Add(n)

	s := scatter(pttoXYs([]int{sequence[0]}), 0)
	s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	s.Shape = draw.CircleGlyph{}
	p.Add(s)
}

func setunpassed(p *plot.Plot, hg *hGlobal) {

	//
	//hg.enD.dec()

	//
	//defer hg.enD.inc()

	var (
		graph = hg.graph

		secnum int
		secL   []int
	)

	var (
		island func(*node)
		//tubes  func(*node, *node)
	)

	island = func(sq *node) {

		sq.mark = secnum
		secL = append(secL, sq.pos)

		for _, sqq := range sq.neighborList() {

			if sqq.mark != 0 {
				continue
			}

			switch sqq.neighborLen() {
			case 1:
			case 2:
				//tubes(sq, sqq)
			default:
				island(sqq)
			}
		}
	}

	for pos := 0; pos < size; pos++ {
		graph[pos].mark = 0
	}

	var tL [3]int

	for pos := 0; pos < size; pos++ {

		sq := &graph[pos]

		if !sq.unpassed || sq.mark != 0 {
			continue
		}

		switch sq.neighborLen() {
		case 1:
		case 2:
			nl := sq.neighborList()
			tL[0] = nl[0].pos
			tL[1] = sq.pos
			tL[2] = nl[1].pos

			q := line(pttoXYs(tL[:]))
			q.LineStyle.Color = color.RGBA{B: 255, A: 255}
			p.Add(q)

		default:
			secL = nil
			secnum++

			island(sq)

			q := scatter(pttoXYs(secL), 0)
			q.Shape = draw.CircleGlyph{}
			q.GlyphStyle.Color = plotutil.Color(int(secnum))
			p.Add(q)
		}
	}

	//endP := make(plotter.XYs, 1)
	//endP[0].X = float64(hg.enD.x)
	//endP[0].Y = float64(hg.enD.y)

	endL := []int{hg.enD.pos}
	q := scatter(pttoXYs(endL), 0)
	q.Shape = draw.BoxGlyph{}
	q.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
	p.Add(q)
}

func setnbL(p *plot.Plot, now *node, next *node) {

	var spl [2]int

	for _, sq := range now.neighborList() {

		spl[0] = now.pos
		spl[1] = sq.pos

		q, err := plotter.NewLine(pttoXYs(spl[:]))
		if err != nil {
			panic(err)
		}

		if sq == next {
			q.LineStyle.Dashes = []vg.Length{vg.Points(3), vg.Points(3)}
			q.LineStyle.Color = color.RGBA{R: 255, A: 255}
			p.Add(q)

			u := scatter(pttoXYs([]int{next.pos}), 0)
			u.Shape = draw.TriangleGlyph{}
			u.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
			p.Add(u)

		} else {
			q.LineStyle.Dashes = []vg.Length{vg.Points(3), vg.Points(3)}
			q.LineStyle.Color = color.RGBA{B: 255, A: 255}
			p.Add(q)
		}
	}
}

func setreapList(p *plot.Plot, reapList twins) {

	var spl [2]int

	for _, sql := range reapList {

		spl[0] = sql[0].pos
		spl[1] = sql[1].pos

		q, err := plotter.NewLine(pttoXYs(spl[:]))
		if err != nil {
			panic(err)
		}

		q.LineStyle.Dashes = []vg.Length{vg.Points(3), vg.Points(3)}
		q.LineStyle.Color = color.RGBA{G: 255, A: 255}
		p.Add(q)
	}
}

func fileout(p *plot.Plot, filename string) {

	if err := p.Save(8*vg.Inch, 8*vg.Inch, filename+".png"); err != nil {
		panic(err)
	}
}

func drawR(title string, filename string, hg *hGlobal) {

	//return

	p := initPlot(title)

	setpassed(p, hg)

	fileout(p, filename)

}

func title(head string, hg *hGlobal) string {

	return fmt.Sprintf("%v %vx%v start=(%v,%v) end=(%v,%v) \n seed=%v try=%v solutions=%v step=%v Elapsed time=%v",
		head,
		xsize,
		ysize,
		hg.start%xsize,
		hg.start/xsize,
		hg.end%xsize,
		hg.end/xsize,
		hg.seed,
		hg.steplimitCt,
		hg.solutionCt,
		hg.stepNum,
		time.Since(hg.ti),
	)

}

func drawD(filenum int, titLe string, now *node, next *node, reapList twins, hg *hGlobal) {

	var filename string

	if titLe == "" {
		titLe = title("Dbug", hg)
		filename = fmt.Sprintf("debug_%v", filenum)
	} else {
		filename = fmt.Sprintf("%v_%v", titLe, filenum)
	}

	p := initPlot(titLe)

	setpassed(p, hg)

	setunpassed(p, hg)

	setnbL(p, now, next)

	setreapList(p, reapList)

	//filename := fmt.Sprintf("debug_%v", filenum)

	fileout(p, filename)
}

func drawO(hg *hGlobal) {

	return

	title := title("Over", hg)

	p := initPlot(title)

	setpassed(p, hg)

	setunpassed(p, hg)

	filename := "Over"

	fileout(p, filename)
}

func mPrint(hg *hGlobal) {

	if !fswitch {
		return
	}

	//for i := 0; i < ysize; i++ {
	for i := range ysize {
		s := i * xsize
		line := hg.board[s : s+xsize]
		hr := hg.row[i]
		fmt.Printf("    %4d%3d%4d%4d \n", line, hr.cnt, hr.rest, hr.threshold)
	}

	if cross {
		fmt.Printf("%4d", hg.diagonal[1].cnt)
	} else {
		fmt.Printf("    ")
	}

	//for i := 0; i < xsize; i++ {
	for i := range xsize {
		fmt.Printf("%5d", hg.columns[i].cnt)
	}

	if cross {
		fmt.Printf("%4d", hg.diagonal[0].cnt)
	}

	fmt.Printf("\n")

	if cross {
		fmt.Printf("%4d", hg.diagonal[1].rest)
	} else {
		fmt.Printf("    ")
	}

	//for i := 0; i < xsize; i++ {
	for i := range xsize {
		fmt.Printf("%5d", hg.columns[i].rest)
	}

	if cross {
		fmt.Printf("%4d", hg.diagonal[0].rest)
	}

	fmt.Printf("\n")

	if cross {
		fmt.Printf("%4d", hg.diagonal[1].threshold)
	} else {
		fmt.Printf("    ")
	}

	//for i := 0; i < xsize; i++ {
	for i := range xsize {
		fmt.Printf("%5d", hg.columns[i].threshold)
	}

	if cross {
		fmt.Printf("%4d", hg.diagonal[0].threshold)
	}

	fmt.Printf("\n")
}

func mFile(hg *hGlobal) {

	if tracefile != "" || !fswitch {
		return
	}

	//
	total := (size + 1) * size / 2

	avg := total / ysize
	for y := range ysize {
		sum := 0
		for x := range xsize {
			sum += hg.board[y*xsize+x]
		}
		if sum != avg {
			stoperr("total err", hg)
		}
	}

	avg = total / xsize
	for x := range xsize {
		sum := 0
		for y := range ysize {
			sum += hg.board[y*xsize+x]
		}
		if sum != avg {
			stoperr("total err", hg)
		}
	}
	//

	mu.Lock()
	solutionNo++
	mu.Unlock()

	fmt.Printf("solution No. %4d\n", solutionNo)

	fname := fmt.Sprintf("%v_%v_%v", hg.start, hg.end, solutionNo)

	//drawR(fname,fname, hg)

	outSequence(fname, hg)

	f, err := os.Create(fname + ".txt")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	for i := range ysize {
		s := i * xsize
		line := hg.board[s : s+xsize]
		//hr := hg.row[i]
		fmt.Fprintf(f, "    %4d \n", line)
	}

	err = f.Sync()
	if err != nil {
		fmt.Println(err)
		return
	}
}
