package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"gopkg.in/ns3777k/go-shodan.v2/shodan"
)

var version = "0.1.0"

var au aurora.Aurora

func parseArgs() (string, string, string, bool) {
	var query, net, ip string
	var compact, color bool

	au = aurora.NewAurora(true)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", au.Bold(os.Args[0]))
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n  Version: %s\n", au.Bold(version))
	}

	flag.StringVar(&query, "q", "", "query ['!http']")
	flag.StringVar(&net, "n", "", "net [192.168.0.0/24]")
	flag.StringVar(&ip, "i", "", "ip [192.168.0.1]")
	flag.BoolVar(&compact, "c", false, "compact, no detail")
	flag.BoolVar(&color, "b", false, "black & white, no color")

	flag.Parse()

	au = aurora.NewAurora(!color)

	return query, net, ip, compact
}

func printHost(j int, h *shodan.HostData) {
	udp := ""
	if h.Transport == "udp" {
		udp = "/udp"
	}

	t, _ := time.Parse(
		"2006-01-02T15:04:05.000000", h.Timestamp)

	fmt.Printf("%d> %s\t[%d%s]\t%s\t\t%s\n",
		j, au.Bold(h.IP), h.Port, udp,
		au.Green(h.Product), t.Format("02/01/2006 15h04"))
	if len(h.OS) > 0 {
		fmt.Printf("  (%s) ", au.Magenta(h.OS))
	}
	if len(h.Hostnames) > 0 {
		fmt.Printf(" ")
		for a := range h.Hostnames {
			fmt.Printf(" %s", au.Cyan(h.Hostnames[a]))
		}
	}
	if len(h.OS) > 0 || len(h.Hostnames) > 0 {
		fmt.Println()
	}
}

func readDefaultQuery(defPath string) (string, error) {
	defQuery := ""

	file, err := os.Open(defPath)
	if err != nil {
		fmt.Println("Enter default query [org:\"some Org\"] : ")
		reader := bufio.NewReader(os.Stdin)
		defQuery, _ = reader.ReadString('\n')
		defQuery = strings.TrimSuffix(defQuery, "\n")

		file, err = os.Create(defPath)
		if err != nil {
			log.Panic(err)
		}
		w := bufio.NewWriter(file)
		fmt.Fprintln(w, defQuery)
		w.Flush()

	} else {
		s := bufio.NewScanner(file)
		for s.Scan() {
			defQuery = s.Text()
		}
	}
	defer file.Close()

	return defQuery, nil

}

func main() {
	if os.Getenv("SHODAN") == "" {
		fmt.Printf("Missing $SHODAN API Key\n export SHODAN=xxxxx...\n")
		os.Exit(0)
	}

	defPath := ".shoddanrc"
	defQuery := ""

	query, net, ip, compact := parseArgs()

	if net == "" && ip == "" {
		defQuery, _ = readDefaultQuery(defPath)
	}
	if defQuery != "" {
		fmt.Printf("Default query from file %s:\n %s\n", defPath, defQuery)
		query += ` ` + defQuery
	}
	if net != "" {
		query = "net:" + net
	}

	client := shodan.NewClient(nil, os.Getenv("SHODAN"))

	// Print only one IP
	if ip != "" {
		h, er := client.GetServicesForHost(ip, nil)
		if er != nil {
			fmt.Println("Error:", er)
		} else {
			for j, hd := range h.Data {
				printHost(j, hd)
				if compact == false {
					fmt.Printf("--\n%+v\n--\n", hd)
				} else {
					fmt.Printf("--\n%+v\n--\n", hd.Data)
				}
			}
			fmt.Println("Maybe more details with:")
			fmt.Println(au.Bold(au.Sprintf(" curl https://api.shodan.io/shodan/host/%s?key=$SHODAN | jq '.'\n", ip)))
		}
		os.Exit(0)
	}

	// Build query
	a := &shodan.HostQueryOptions{Query: query} //first query "org: Company port: ....."
	log.Printf("%+v\n", a)

	// Count
	r, e := client.GetHostsCountForQuery(a)
	if e != nil {
		fmt.Println("Error GetHostsCountForQuery:", e)
		os.Exit(0)
	}
	log.Println(r.Total)

	// Query
	res, err := client.GetHostsForQuery(a)
	log.Println(res.Total)
	if err != nil {
		fmt.Println("Error HostsForQuery:", err)
		os.Exit(0)
	}
	// Print results
	for j := range res.Matches {
		printHost(j, res.Matches[j])
		if net != "" && compact == false { // for net query only
			fmt.Printf("--\n%+v\n--\n", res.Matches[j].Data)
		}
	}
}
