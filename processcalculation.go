package main

// This program calculates the CPU and memory usage of a process given its ID
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// checkerror is a helper function that panics if there is an error
func checkerror(err error) {
	if err != nil {
		panic(err)
	}
}

// getfromcmd is a function that prompts the user to enter the process ID and returns it as a string
func getfromcmd() string {
	fmt.Print("Enter the Process Id ")
	var input string
	fmt.Scanln(&input)
	return input
}

// prcesscpucalc is a function that calculates the CPU usage percentage of a process given its ID
func prcesscpucalc(processid string) float32 {

	var utime string     // user mode time
	var sime string      // system mode time
	var starttime string // start time of the process
	var uptime string    // uptime of the system

	// open the /proc/<pid>/stat file and read the relevant fields
	data, err := os.Open("/proc/" + processid + "/stat")
	checkerror(err)
	defer data.Close()
	//tik, err := exec.Command("getconf CLK_TCK").Output()
	//tik_int := binary.BigEndian.Uint64(tik)
	tik := 100 // clock ticks per second
	checkerror(err)
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanWords)
	count := 0

	for scanner.Scan() {
		count++
		line := scanner.Text()

		if count == 14 {
			utime = line // 14th field is user mode time
		} else if count == 15 {
			sime = line // 15th field is system mode time

		} else if count == 22 {
			starttime = line // 22nd field is start time of the process

		}
	}

	// open the /proc/uptime file and read the first field
	datanew, err := os.Open("/proc/uptime")
	checkerror(err)
	scanner2 := bufio.NewScanner(datanew)
	for scanner2.Scan() {
		text := scanner2.Text()
		splitvari := strings.Split(text, " ")
		uptime = splitvari[0] // first field is uptime of the system
	}

	// convert the strings to integers or floats
	utime_int, err := strconv.Atoi(utime)
	uptime_int, err := strconv.ParseFloat(uptime, 32)
	sime_int, err := strconv.Atoi(sime)
	starttime_int, err := strconv.Atoi(starttime)

	// calculate the CPU usage percentage using the formula
	// CPU_usage = (process_utime_sec + process_sime_sec) * 100 / process_elapsed_sec
	process_utime_sec := utime_int / tik
	process_sime_sec := sime_int / tik
	process_starttime_sec := starttime_int / tik
	process_starttime_sec_float := float32(process_starttime_sec)

	process_elapsed_sec := float32(uptime_int) - process_starttime_sec_float
	process_usage_sec := process_utime_sec + process_sime_sec
	process_usage := float32(process_usage_sec) * 100 / process_elapsed_sec

	return process_usage

}

// checkrunninguser is a function that checks if the program is run by the root user or with sudo
func checkrunninguser() {
	userid := os.Getuid()

	if userid != 0 {
		fmt.Println("ERROR : THIS Program needs to be run with root user or with sudo")
		os.Exit(2)
	}
}

// processmemorycalc is a function that calculates the memory usage percentage of a process given its ID
func processmemorycalc(processid string) float32 {

	var vmrss string    // resident set size
	var totalem string  // total memory
	var memory_precentage_usage float32 // memory usage percentage

	// open the /proc/<pid>/status file and read the VmRSS field
	data, err := os.Open("/proc/" + processid + "/status")
	checkerror(err)
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "VmRSS") {
			vmrss2 := strings.Split(line, "          ")
			finalvmrss := strings.Split(vmrss2[1], " ")
			vmrss = finalvmrss[0] // VmRSS is the resident set size
		}
	}

	// open the /proc/meminfo file and read the MemTotal field
	datanew, err := os.Open("/proc/meminfo")
	checkerror(err)
	scanner2 := bufio.NewScanner(datanew)
	for scanner2.Scan() {
		line := scanner2.Text()
		if strings.Contains(line, "MemTotal") {
			totmemory := strings.Split(line, "        ")
			//fmt.Println(totmemory[1])
			finaltotal := strings.Split(totmemory[1], " ")
			totalem = finaltotal[0] // MemTotal is the total memory
		}

	}

	// convert the strings to floats and remove the "kB" suffix
	vmrss_float, err := strconv.ParseFloat(strings.ReplaceAll(vmrss, "kB", ""), 32)
	totalmem_float, err := strconv.ParseFloat(strings.ReplaceAll(totalem, "kB", ""), 32)
	//fmt.Println(vmrss_float, totalmem_float)

	// calculate the memory usage percentage using the formula
	// memory_precentage_usage = (vmrss_float / totalmem_float) * 100
	memory_precentage_usage = (float32(vmrss_float) / float32(totalmem_float)) * 100
	return memory_precentage_usage
}

func main() {

	checkrunninguser() // check if the program is run by the root user or with sudo
	processid := getfromcmd() // get the process ID from the user

	process_cpu_usage := prcesscpucalc(processid) // calculate the CPU usage percentage of the process
	process_memory_usage := processmemorycalc(processid) // calculate the memory usage percentage of the process

	// print the results
	fmt.Println("The CPU usage percentage of the process is", process_cpu_usage)
	fmt.Println("The memory usage percentage of the process is", process_memory_usage)
}
