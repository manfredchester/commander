package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/connext-cs/pub/zhlog"
)

func schedul() {
	// tm()
	// tc()
	// tmc()
	// tmstop()
	// zone()
	// tcstop()
}

func zone() {
	// 生成UTC的时间
	t1, err := time.Parse("2006-01-02 15:04:05", "2020-02-21 08:59:10")
	zhlog.Assert(err)
	fmt.Println(t1)
	// fmt.Println(t1.Unix())

	// 生成当地的时间带有CST
	local, err := time.LoadLocation("Local")
	zhlog.Assert(err)
	str, err := time.ParseInLocation("2006-01-02 15:04:05", "2020-02-21 08:59:10", local)
	zhlog.Assert(err)
	fmt.Println(str)
	// fmt.Println(str.Unix())

	timestamp := time.Now().Unix()
	// fmt.Println(timestamp)
	tm := time.Unix(timestamp, 0)
	fmt.Println(tm)
}

func tcstop() {
	t3 := time.NewTicker(2 * time.Second)
	go func() {
		for {
			<-t3.C
			fmt.Println("t3")
		}
	}()
	<-time.After(7 * time.Second)
	t3.Stop()
}

func tmstop() {
	t3 := time.NewTimer(3 * time.Second)
	go func() {
		<-t3.C
		fmt.Println("t3")
	}()
	stop := t3.Stop()
	if stop {
		fmt.Println("t3 stop")
	}
}

func tmc() {
	var wg sync.WaitGroup
	wg.Add(2)
	tim1 := time.NewTimer(1 * time.Second)
	tic1 := time.NewTicker(3 * time.Second)

	go func(s *time.Ticker) {
		defer wg.Done()
		for {
			<-s.C
			fmt.Println("tic1")
		}
	}(tic1)

	go func(t *time.Timer) {
		defer wg.Done()
		for {
			<-t.C
			fmt.Println("tim1")
		}
	}(tim1)

	wg.Wait()
}

func tc() {
	tc := time.NewTicker(time.Second)

	for {
		<-tc.C
		fmt.Println("tc1")
	}
}

func tm() {
	timer1 := time.NewTimer(time.Second * time.Duration(1))
	<-timer1.C
	fmt.Println("timer1")
	<-time.After(time.Duration(5) * time.Second)
	fmt.Println("timer2")
	timer1.Reset(time.Second * time.Duration(10))
	<-timer1.C
	fmt.Println("timer3")

}
