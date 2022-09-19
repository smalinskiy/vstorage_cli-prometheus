package tools

import (
	"fmt"
	"sync"
	"time"

	cpu "github.com/shirou/gopsutil/cpu"
	mem "github.com/shirou/gopsutil/mem"
)

func getStats() {
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	c, _ := cpu.Percent(1*time.Second, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Cpu usage: %v\n", c)
}

func localTelemetry(wg *sync.WaitGroup, instance int) {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			getStats()
		case <-quit:
			ticker.Stop()
		}
		wg.Done()
	}

}

func RunTelemetry() {
	fmt.Println("main started..")
	var wg sync.WaitGroup

	for i := 1; i <= 1; i++ {
		wg.Add(1)
		go localTelemetry(&wg, 1)
	}
	wg.Wait()
}
