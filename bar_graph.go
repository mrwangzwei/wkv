package main

import (
	"bufio"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"wkv/bar_graph"
)

func main() {
	table := bar_graph.NewBarTable([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	table.Draw(false)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	stdinTxt := make(chan string)

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			stdinTxt <- scanner.Text()
		}
	}()

	for {
		select {
		case <-quit:
			goto end
		case v, ok := <-stdinTxt:
			if ok {
				sli := strings.Split(v, " ")
				if len(sli) > 1 {
					pre, _ := strconv.Atoi(sli[0])
					next, _ := strconv.Atoi(sli[1])
					table.SwapAndDraw(pre, next)
				}
			}
		}
	}
end:
}
