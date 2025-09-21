package hamiltonpath

import (
	"fmt"
	"os"
	"time"
)

func stoperr(msg string, hg *hGlobal) {

	fmt.Printf("Error %v limitCt=%v\n", msg, hg.steplimitCt)

	//fmt.Printf("process time: %s\n", time.Since(hg.ti))

	os.Exit(0)

	//wg.Done()
	//<-ch
}

func stopunresolve(st string, hg *hGlobal) {
	fmt.Printf("Unresolve!! %v Seed=%v Try=%v Process time: %s\n", st, hg.seed, hg.steplimitCt, time.Since(hg.ti))

	wg.Done()
	<-ch
}

func stopfound(hg *hGlobal) {

	fmt.Printf("Found!! %v Seed=%v Try=%v Process time: %s\n", hg.solutionCt, hg.seed, hg.steplimitCt, time.Since(hg.ti))

	//fmt.Printf("process time: %s\n", time.Since(hg.ti))

	title := title("Knight Tour", hg)

	var filename string
	if fswitch {
		filename = fmt.Sprintf("Kn%vX%v_%v", xsize, ysize, hg.seed)

		drawR(title, filename, hg)

		outSequence(filename, hg)

	} else {
		//filename = "Result"

		//if threads == 1 {
		//	drawR(title, filename, hg)

		//	outSequence(filename, hg)

		//}
	}

	wg.Done()
	<-ch
}

func stopovr(hg *hGlobal) {

	fmt.Printf("\033[31mOverLimit!! \033[0mSeed = %v try=%v Process time: %s\n", hg.seed, hg.steplimitCt, time.Since(hg.ti))

	//fmt.Printf("process time: %s\n", time.Since(hg.ti))

	//if threads == 1 {
	//	drawO(hg)
	//}

	wg.Done()
	<-ch

}

func stopenumerated(hg *hGlobal) {

	fmt.Printf("%v Seed = %v ThreadHead %v Enumerated!! %v Try = %v Process time: %s\n",
		threadcnt, hg.seed, hg.th, hg.solutionCt, hg.steplimitCt, time.Since(hg.ti))

	mu.Lock()
	Totalsolutions += hg.solutionCt
	Totalsteps += hg.steplimitCt
	threadcnt--
	mu.Unlock()

	wg.Done()
	//<-ch
}

func stopparamerr(str string) {

	fmt.Printf("%v Parameter error", str)

	os.Exit(0)
}
