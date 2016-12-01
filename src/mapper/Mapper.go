package mapper


import (
	"fmt"
	"log"
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
	fileHandles int
	listening []string
	clients []string
}

func finaliseServerStats(stats *ServerStats, processMap map[int]PidMap){


	hhh, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}


	stats.processes = len(processMap)
	stats.hostname = hhh
	stats.ip_address = GetOutboundIP()

	stats.net_load = 0.0
	stats.disk_load = 0.0
	stats.cpu_load = Win_CpuLoad()



}

func writeServerStatsToInflux(serverStats *ServerStats, influxURL string, influxU string, influxP string)  {


	/**
	* InfluxDB connection
	*
	 */

	c, err := client.NewHTTPClient(client.HTTPConfig{
		//		Addr: "http://localhost:8086",
		Addr: influxURL,
		Username: influxU,
		Password: influxP,
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

	println("writing output")


	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)


	// Write the batch
	err = c.Write(bp)
	if err != nil {
		log.Fatalln("Failed to write Error: ", err)
	}

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


func MainGo(path string,  influxURL string, influxU string, influxP string) {

	fmt.Printf("------------- Collecting stats ")

	stats := new(ServerStats)

	var processMap = make(map[int]PidMap)


	fmt.Printf("collecting handles ")
	processHandles(path, stats, processMap)


	fmt.Printf("collecting netstat ")

	processNetstat(stats, processMap)

	fmt.Printf("write to influx")
	writeServerStatsToInflux(stats, influxURL, influxU, influxP)

//	if err2 != nil {
//		log.Fatal(err)
//	}
	fmt.Printf("output host:%s handles:%d\n", stats.hostname, stats.file_handles)
}
