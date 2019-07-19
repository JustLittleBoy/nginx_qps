package main // import "stubstatus"

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var server string
var tick *time.Ticker
var oldFields map[string]int
var maxQps, maxLine int
var maxQpsTime, maxLineTime string

func main() {
	app := cli.NewApp()
	app.Name = "nginx QPS watch"
	app.Usage = "nginx QPS watch"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		{
			Name:  "wangchao",
			Email: "wangchaoxiaoban@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "url",
			Value:       "http://127.0.0.1:9100",
			Usage:       "nginx server status url",
			Destination: &server,
		},
	}

	app.Action = func(_ *cli.Context) error {
		tick = time.NewTicker(time.Second)
		return getStubStatus()
	}

	err := app.Run(os.Args)
	if nil != err {
		log.Printf("MAIN Run error; " + err.Error())
	}
}

func getStubStatus() error {
	oldFields = map[string]int{
		"active":   int(0),
		"accepts":  int(0),
		"handled":  int(0),
		"requests": int(0),
		"reading":  int(0),
		"writing":  int(0),
		"waiting":  int(0),
	}
	for {
		select {
		case <-tick.C:
			err := getNginxData()
			if nil != err {
				fmt.Println("get data error:" + err.Error())
			}
		}
	}
	return nil
}

func getNginxData() error {
	resp, err := http.Get(server)
	if nil != err {
		fmt.Println("get content error:" + server)
		return err
	}

	r := bufio.NewReader(resp.Body)

	// Active connections
	_, err = r.ReadString(':')
	if err != nil {
		return err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	active, err := strconv.ParseUint(strings.TrimSpace(line), 10, 64)
	if err != nil {
		return err
	}

	// Server accepts handled requests
	_, err = r.ReadString('\n')
	if err != nil {
		return err
	}
	line, err = r.ReadString('\n')
	if err != nil {
		return err
	}
	data := strings.Fields(line)
	accepts, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		return err
	}

	handled, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return err
	}
	requests, err := strconv.ParseUint(data[2], 10, 64)
	if err != nil {
		return err
	}

	// Reading/Writing/Waiting
	line, err = r.ReadString('\n')
	if err != nil {
		return err
	}
	data = strings.Fields(line)
	reading, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return err
	}
	writing, err := strconv.ParseUint(data[3], 10, 64)
	if err != nil {
		return err
	}
	waiting, err := strconv.ParseUint(data[5], 10, 64)
	if err != nil {
		return err
	}

	fields := map[string]int{
		"active":   int(active),
		"accepts":  int(accepts),
		"handled":  int(handled),
		"requests": int(requests),
		"reading":  int(reading),
		"writing":  int(writing),
		"waiting":  int(waiting),
	}
	if 0 == oldFields["requests"] {
		oldFields = fields
		maxQps = 0
		maxLine = 0
		return nil
	}
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	qps := fields["requests"] - oldFields["requests"]
	onLine := fields["active"]
	oldFields = fields
	if qps > maxQps {
		maxQps = qps
		maxQpsTime = nowTime
	}

	if onLine > maxLine {
		maxLine = onLine
		maxLineTime = nowTime
	}
	fmt.Println("时间:" + time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("当前QPS: " + strconv.Itoa(qps))
	fmt.Println("当前连接数: " + strconv.Itoa(fields["active"]))
	fmt.Println("最大QPS:" + strconv.Itoa(maxQps) + " 发生时间:" + maxQpsTime)
	fmt.Println("最大连接数:" + strconv.Itoa(maxLine) + " 发生时间:" + maxLineTime + "\n")
	return nil
}
