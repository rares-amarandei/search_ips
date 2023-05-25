package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"leftclick.io/geolocation/geoip"
)

func scalePopulationsToQuota(cities map[string]map[string]float64, quota float64) map[string]map[string]int {
	// Calculate total population
	totalPopulation := 0.0
	for _, countryCities := range cities {
		for _, population := range countryCities {
			totalPopulation += population
		}
	}

	// Calculate the scaling factor
	scalingFactor := quota / totalPopulation

	// Create the new map
	citiesInt := make(map[string]map[string]int)

	for country, countryCities := range cities {
		citiesInt[country] = make(map[string]int)
		for city, population := range countryCities {
			// Scale the population and round to the nearest integer
			scaledPopulation := int(math.Round(population * scalingFactor))

			citiesInt[country][city] = scaledPopulation
		}
	}

	return citiesInt
}

func writeIntMapToJsonFile(data map[string]map[string]int, filename string) {
	// Convert the data to JSON
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding to JSON: %v\n", err)
		return
	}

	// Write the JSON data to the file
	err = ioutil.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Data written to %s\n", filename)
}

func parseInputCsv(input string) map[string]map[string]float64 {
	// Open the input file
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Initialize the citiesByCountry map
	citiesByCountry := make(map[string]map[string]float64)

	// Read the records from the input file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	// Iterate over the records
	for _, record := range records {
		// Normalize country and city names
		country := strings.Split(record[0], " (")[0]
		city := strings.Split(record[1], " (")[0]

		// Ignore names with a comma
		if strings.Contains(country, ",") || strings.Contains(city, ",") {
			continue
		}

		// Capitalize the first letter of each word
		country = strings.Title(strings.ToLower(country))
		city = strings.Title(strings.ToLower(city))

		// Parse the population
		population, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Println("Error parsing population:", err)
			continue
		}

		// Convert the population to millions
		population = population / 1_000_000

		// If this country hasn't been seen before, initialize it in the map
		if _, ok := citiesByCountry[country]; !ok {
			citiesByCountry[country] = make(map[string]float64)
		}

		// Add this city to the country's map
		citiesByCountry[country][city] = population
	}

	return citiesByCountry
}

func intToIP(ipInt uint32) string {
	ip := make([]byte, 4)

	// Extract octets using bitwise operations
	ip[0] = byte(ipInt & 0xFF)
	ip[1] = byte((ipInt >> 8) & 0xFF)
	ip[2] = byte((ipInt >> 16) & 0xFF)
	ip[3] = byte((ipInt >> 24) & 0xFF)

	return fmt.Sprintf("%d.%d.%d.%d", ip[3], ip[2], ip[1], ip[0])
}

type Country struct {
	Cities []*City
}

type City struct {
	Name  string
	IPs   []string
	Quota int
}

func bruteForceSearch(provider geoip.Provider, citiesByCountry map[string]map[string]int) map[string]*Country {
	incomplete := make(map[string]*City)
	completed := make(map[string]*City)

	countries := make(map[string]*Country)

	// Initialize all cities as incomplete
	for _, cities := range citiesByCountry {
		for city, quota := range cities {
			incomplete[city] = &City{Name: city, IPs: []string{}, Quota: quota}
		}
	}

	var ip uint32 = 1000
	var step uint32 = 1000
	// Start from 10000 and increase by 10000 for each iteration
	for len(incomplete) > 0 {
		if ip < 1000 {
			break
		}
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

			// If we have collected enough IPs for this city, move it to the completed map
			if len(city.IPs) == city.Quota {
				completed[info.City] = city
				delete(incomplete, info.City)
				fmt.Printf("Completed: %s, Remained: %d\n", info.City, len(incomplete))
			}
		} else {
			step = 1000
		}
	}

	// If reached uint32(-1), move cities with any IPs collected to the completed map
	if ip < 1000 {
		for key, city := range incomplete {
			if len(city.IPs) > 0 {
				completed[key] = city
				delete(incomplete, key)
			}
		}
	}

	file, err := os.Create("incomplete.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
	}
	defer file.Close()
	for key, _ := range incomplete {
		_, err := file.WriteString(key + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}

	// Loop over the completed cities and group them by country
	for country, cities := range citiesByCountry {
		if _, ok := countries[country]; !ok {
			countries[country] = &Country{Cities: []*City{}}
		}
		for cityName, _ := range cities {
			if city, ok := completed[cityName]; ok {
				countries[country].Cities = append(countries[country].Cities, city)
			}
		}
	}

	return countries
}

func main() {
	db_path := flag.String("db", "IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ZIPCODE-TIMEZONE-ISP-DOMAIN-NETSPEED-AREACODE-WEATHER-MOBILE-ELEVATION-USAGETYPE.BIN", "path to db")
	flag.Parse()

	geo, err := geoip.New(*db_path)
	if err != nil {
		fmt.Println(err)
	}
	inputCsv := parseInputCsv("input.csv")
	citiesByCountry := scalePopulationsToQuota(inputCsv, 200000)
	writeIntMapToJsonFile(citiesByCountry, "out.json")
	countries := bruteForceSearch(geo, citiesByCountry)

	foundIps := 0
	for _, country := range countries {
		for _, city := range country.Cities {
			foundIps += len(city.IPs)
		}
	}

	jsonString, err := json.MarshalIndent(countries, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	// Write the JSON to a file
	file, err := os.Create("countries.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonString)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("Found ips: %d", foundIps)

}
