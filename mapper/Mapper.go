package mapper


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
	MyDB2 = "introspector-1"
)
type ServerStats struct {
	hostname, ip_public, ip_address, ip_subnet, ip_gateway string
	processes, file_handles, server_ports, client_ports int
	cpu_load, disk_load, net_load, mem_util float32
}

type PidMap struct {
	name, owner string
	pid, ppid int
	files []string
}



func processNetstat(stats *ServerStats, processMap map[int]PidMap) {
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


func main2() {


	stats := new(ServerStats)

	var processMap = make(map[int]PidMap)


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
