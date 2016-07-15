package main


import (
	"fmt"
	"log"
	"os/exec"
	"bufio"
	"io"
	"strings"
	"github.com/influxdata/influxdb/client/v2"
	"time"
	"os"
	"net"
)

const (
	MyDB = "introspector-1"
	username = ""
	password = ""
)
type ServerStats struct {
	hostname, ip_public, ip_address, ip_subnet, ip_gateway string
	processes, file_handles, server_ports, client_ports int
	cpu_load, disk_load, net_load, mem_util float32

}


func processHandles(stats *ServerStats, processMap map[string]map[string][]string) {
	handle := exec.Command("./etc/Handle.exe")
	handlesPipe,err := handle.StdoutPipe()
	handle.Start()
	handlesProcessReader := bufio.NewReader(handlesPipe)

	var currentPIdMap = make(map[string][]string)


	for {
		var line string
		line, err = handlesProcessReader.ReadString('\n')
		if err == io.EOF  {
			// Good end of file with no partial line
			break
		}
		if err == io.EOF {
			err := fmt.Errorf("Last line not terminated: %q", line)
			panic(err)
		}
		if err != nil {
			panic(err)
		}

		//fmt.Print(line)
		if (strings.Contains(line, "pid:")) {

			if (!strings.Contains(line, "unable to open process")) {

				var parts = strings.Fields(line)
				currentPIdMap = make(map[string][]string)

				currentPIdMap["name"] = []string{parts[0] }
				currentPIdMap["pid"] = []string{parts[2] }
				currentPIdMap["owner"] = []string{parts[3] }
				currentPIdMap["files"] = []string{}
				processMap[parts[2]] = currentPIdMap
			}

		} else if (strings.Contains(line, " File ")){
			var parts = strings.Fields(line)
			var sofile = parts[3]
			var fileIndex = strings.LastIndex(line, sofile)
			var ff = line[fileIndex:len(line)]
			currentPIdMap["files"] = append(currentPIdMap["files"],ff)
			stats.file_handles += 1
		}
	}


	// Wait for the result of the command; also closes our end of the pipe
	err = handle.Wait()
	fmt.Printf("HANDLES 2 %s\n", len(processMap))

}

func processNetstat(stats *ServerStats, processMap map[string]map[string][]string) {
	netout := exec.Command("netstat", "-a", "-n", "-o")
	netstatPipe,err2 := netout.StdoutPipe()

	netout.Start()
	netstatProcessReader := bufio.NewReader(netstatPipe)
	for {
		var line string
		line, err2 = netstatProcessReader.ReadString('\n')
		if err2 == io.EOF  {
			// Good end of file with no partial line
			break
		}
		if err2 == io.EOF {
			err2 := fmt.Errorf("Last line not terminated: %q", line)
			panic(err2)
		}
		if err2 != nil {
			panic(err2)
		}
		if (strings.Contains(line, "LISTENING")) {
			stats.server_ports += 1
		}
		if (strings.Contains(line, "ESTABLISHED")) {
			stats.client_ports += 1
		}

		fmt.Print(line)
	}
}




func writeServerStats(serverStats *ServerStats)  {


	/**
	* InfluxDB connection
	*
	 */

	c, err := client.NewHTTPClient(client.HTTPConfig{
		//		Addr: "http://localhost:8086",
		Addr: "http://192.168.99.100:32768",
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	tags := map[string]string{ "server": serverStats.hostname, "ip_address": serverStats.ip_address, "ip_subnet": serverStats.ip_subnet, "ip_gateway" : serverStats.ip_gateway}
	fields := map[string]interface{}{
		"processes": serverStats.processes,
		"cpu_load":  serverStats.cpu_load,
		"disk_load":   serverStats.disk_load,
		"net_load":   serverStats.net_load,
		"mem_util":   serverStats.mem_util,
		"ports_server": serverStats.server_ports,
		"ports_client":  serverStats.client_ports,
		"file_handles": serverStats.file_handles,
		"ip_address": serverStats.ip_address,
		"ip_public": serverStats.ip_public,
		"ip_gateway": serverStats.ip_gateway,
		"ip_subnet": serverStats.ip_subnet,
	}
	// + subnet, ip
	pt, err := client.NewPoint("server_stats", tags, fields, time.Now())


	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)

	// Write the batch
	c.Write(bp)
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}


func main() {


	stats := new(ServerStats)

	var processMap = make(map[string]map[string][]string)


	fmt.Printf("handles ")
	processHandles(stats, processMap)


	fmt.Printf("netstat ")

	processNetstat(stats, processMap)



	hhh, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}


	stats.processes = len(processMap)
	stats.hostname = hhh
	stats.net_load = 0.12
	stats.disk_load = 0.5
	stats.cpu_load = 25.5
	stats.ip_address = GetOutboundIP()

	writeServerStats(stats)



//	if err2 != nil {
//		log.Fatal(err)
//	}
	fmt.Printf("output host:%s handles:%d\n", stats.hostname, stats.file_handles)
}
