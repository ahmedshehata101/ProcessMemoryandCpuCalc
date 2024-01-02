package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func checkerror(err error) {
	if err != nil {
		panic(err)
	}
}

func getfromcmd() string {
	fmt.Print("Enter the Process Id ")
	var input string
	fmt.Scanln(&input)
	return input
}
func prcesscpucalc(processid string) float32 {

	var utime string
	var sime string
	var starttime string
	var uptime string

	data, err := os.Open("/proc/" + processid + "/stat")
	checkerror(err)
	defer data.Close()
	//tik, err := exec.Command("getconf CLK_TCK").Output()
	//tik_int := binary.BigEndian.Uint64(tik)
	tik := 100
	checkerror(err)
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanWords)
	count := 0

	for scanner.Scan() {
		count++
		line := scanner.Text()

		if count == 14 {
			utime = line
		} else if count == 15 {
			sime = line

		} else if count == 22 {
			starttime = line

		}
	}

	datanew, err := os.Open("/proc/uptime")
	checkerror(err)
	scanner2 := bufio.NewScanner(datanew)
	for scanner2.Scan() {
		text := scanner2.Text()
		splitvari := strings.Split(text, " ")
		uptime = splitvari[0]
	}

	utime_int, err := strconv.Atoi(utime)
	uptime_int, err := strconv.ParseFloat(uptime, 32)
	sime_int, err := strconv.Atoi(sime)
	starttime_int, err := strconv.Atoi(starttime)

	process_utime_sec := utime_int / tik
	process_sime_sec := sime_int / tik
	process_starttime_sec := starttime_int / tik
	process_starttime_sec_float := float32(process_starttime_sec)

	process_elapsed_sec := float32(uptime_int) - process_starttime_sec_float
	process_usage_sec := process_utime_sec + process_sime_sec
	process_usage := float32(process_usage_sec) * 100 / process_elapsed_sec

	return process_usage

}
func checkrunninguser() {
	userid := os.Getuid()

	if userid != 0 {
		fmt.Println("ERROR : THIS Program needs to be run with root user or with sudo")
		os.Exit(2)
	}
}

func processmemorycalc(processid string) float32 {

	var vmrss string
	var totalem string
	var memory_precentage_usage float32
	data, err := os.Open("/proc/" + processid + "/status")
	checkerror(err)
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "VmRSS") {
			vmrss2 := strings.Split(line, "	  ")
			finalvmrss := strings.Split(vmrss2[1], " ")
			vmrss = finalvmrss[0]
		}
	}
	datanew, err := os.Open("/proc/meminfo")
	checkerror(err)
	scanner2 := bufio.NewScanner(datanew)
	for scanner2.Scan() {
		line := scanner2.Text()
		if strings.Contains(line, "MemTotal") {
			totmemory := strings.Split(line, "        ")
			//fmt.Println(totmemory[1])
			finaltotal := strings.Split(totmemory[1], " ")
			totalem = finaltotal[0]
		}

	}

	vmrss_float, err := strconv.ParseFloat(strings.ReplaceAll(vmrss, "kB", ""), 32)
	totalmem_float, err := strconv.ParseFloat(strings.ReplaceAll(totalem, "kB", ""), 32)
	//fmt.Println(vmrss_float, totalmem_float)
	memory_precentage_usage = (float32(vmrss_float) / float32(totalmem_float)) * 100
	return memory_precentage_usage
}

func main() {

	checkrunninguser()
	processid := getfromcmd()

	process_cpu_usage := prcesscpucalc(processid)
	physical_memory_usage := processmemorycalc(processid)
	fmt.Println("Process ID : "+processid+" and its CPU usage is : ", process_cpu_usage, "%")
	fmt.Println("Process ID : "+processid+" and its Memory usage is : ", physical_memory_usage, "%")
}
