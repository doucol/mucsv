package main

import (
	"encoding/csv"
	"fmt"
	"github.com/doucol/meetup-client"
	"log"
	"os"
	"reflect"
	"strings"
)

func main() {
	groupNames := strings.Split(os.Getenv("MUCSV_GROUPS"), ",")
	if len(os.Args) > 1 {
		groupNames = strings.Split(os.Args[1], ",")
	}
	key := os.Getenv("MUCSV_APIKEY")

	if key == "" {
		log.Fatalln("The 'MUCSV_APIKEY' env variable is missing")
	}
	if len(groupNames) == 0 {
		log.Fatalln("No groups: You must set the MUCSV_GROUPS env variable, or pass as first argument a comma separated list of groups")
	}

	client := meetup.NewClient(&meetup.ClientOpts{
		APIKey: key,
	})
	groups, err := client.GroupByURLName(groupNames)
	if err != nil {
		log.Fatalln(err)
	}

	cw := csv.NewWriter(os.Stdout)
	defer cw.Flush()
	headerWritten := false

	for _, group := range groups.Groups {
		for i1 := 0; i1 < 10000; i1++ {
			members, err := client.MembersByPage(group.ID, i1)
			if err != nil {
				log.Fatalln(err)
			}
			for _, member := range members.Members {
				if !headerWritten {
					headers := convertToHeaders(member)
					headers = append([]string{"GroupName"}, headers...)
					err = cw.Write(headers)
					if err != nil {
						log.Fatalln(err)
					}
					headerWritten = true
				}
				values := append([]string{group.URLName}, convertToStringSlice(member)...)
				err = cw.Write(values)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
}

func convertToHeaders(data interface{}) []string {
	v := reflect.Indirect(reflect.ValueOf(data)).Type()
	n := v.NumField()
	if n <= 0 { return []string{} }
	headers := make([]string, n)
	for i1 := 0; i1 < n; i1++ {
		headers[i1] = v.Field(i1).Name
	}
	return headers
}

func convertToStringSlice(data interface{}) []string {
	v := reflect.ValueOf(data)
	n := v.NumField()
	if n <= 0 { return []string{} }
	rowContents := make([]string, n)
	for i := 0; i < n; i++ {
		x := v.Field(i)
		s := fmt.Sprintf("%v", x.Interface())
		rowContents[i] = s
	}
	return rowContents
}
