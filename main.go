package main

import (
	"flag"
	"fmt"
	"os"

	"leftclick.io/geolocation/geoip"
)

func intToIP(ipInt uint32) string {
	ip := make([]byte, 4)

	// Extract octets using bitwise operations
	ip[0] = byte(ipInt & 0xFF)
	ip[1] = byte((ipInt >> 8) & 0xFF)
	ip[2] = byte((ipInt >> 16) & 0xFF)
	ip[3] = byte((ipInt >> 24) & 0xFF)

	return fmt.Sprintf("%d.%d.%d.%d", ip[3], ip[2], ip[1], ip[0])
}

type City struct {
	Name string
	IPs  []string
}

func bruteForceSearch(provider geoip.Provider, citiesByCountry map[string][]string) map[string]*City {
	incomplete := make(map[string]*City)
	completed := make(map[string]*City)

	// Initialize all cities as incomplete
	for _, cities := range citiesByCountry {
		for _, city := range cities {
			incomplete[city] = &City{Name: city, IPs: []string{}}
		}
	}

	var ip uint32 = 0
	var step uint32 = 1000
	// Start from 10000 and increase by 10000 for each iteration
	for len(incomplete) > 0 {
		ip += step
		myIp := intToIP(ip)
		info, err := provider.Lookup(myIp)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		// If the city is in the incomplete map, add the IP to its list
		if city, ok := incomplete[info.City]; ok {
			step = 1
			city.IPs = append(city.IPs, myIp)

			// If we have collected 100 IPs for this city, move it to the completed map
			if len(city.IPs) == 100 {
				completed[info.City] = city
				delete(incomplete, info.City)
				fmt.Printf("Completed: %s, Remained: %d\n", info.City, len(incomplete))
				if len(incomplete) < 12 {
					for city, val := range incomplete {
						fmt.Printf("city %s - %d\n", city, len(val.IPs))
					}
					fmt.Printf("current ip: %d", ip)
				}
			}
		} else {
			step = 1000
		}
	}

	return completed
}

func main() {
	db_path := flag.String("db", "IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ZIPCODE-TIMEZONE-ISP-DOMAIN-NETSPEED-AREACODE-WEATHER-MOBILE-ELEVATION-USAGETYPE.BIN", "path to db")
	flag.Parse()

	geo, err := geoip.New(*db_path)
	if err != nil {
		fmt.Println(err)
	}
	var citiesByCountry = map[string][]string{
		"US": {"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose", "Seattle", "Miami"},
		"RU": {"Moscow", "Saint Petersburg", "Novosibirsk", "Yekaterinburg", "Kazan", "Krasnoyarsk", "Chelyabinsk", "Samara", "Omsk", "Rostov"},
		"CN": {"Shanghai", "Beijing", "Chongqing", "Guangzhou", "Shenzhen", "Tianjin", "Chengdu", "Wuhan", "Hangzhou", "Hong Kong"},
		"IN": {"Mumbai", "Delhi", "Bengaluru", "Hyderabad", "Ahmedabad", "Chennai", "Kolkata", "Surat", "Pune", "Jaipur"},
		"FR": {"Paris", "Marseille", "Lyon", "Toulouse", "Nice", "Nantes", "Strasbourg", "Montpellier", "Bordeaux", "Lille"},
		"DE": {"Berlin", "Hamburg", "Munchen", "Frankfurt am Main", "Stuttgart", "Dusseldorf", "Dortmund", "Essen", "Leipzig"},
		"JP": {"Tokyo", "Yokohama", "Osaka", "Nagoya", "Sapporo", "Fukuoka", "Kobe", "Kyoto", "Kawasaki", "Saitama"},
		"GB": {"London", "Birmingham", "Manchester", "Glasgow", "Newcastle", "Sheffield", "Liverpool", "Leeds", "Bristol", "Edinburgh"},
		"KR": {"Seoul", "Busan", "Gwangmyeong", "Daegu", "Daejeon", "Gwangju", "Ulsan", "Suwon", "Seongnam", "Cheongju"},
		"SA": {"Riyadh", "Jeddah", "Mecca", "Medina", "Dammam", "Ta'if", "Tabuk", "Abha", "Buraydah"},
	}
	completed := bruteForceSearch(geo, citiesByCountry)

	file, err := os.Create("out.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for _, city := range completed {
		output := fmt.Sprintf("City: %s, IPs: %v\n", city.Name, city.IPs)
		_, err := file.WriteString(output)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	fmt.Println("Output written to file: out.txt")
}
