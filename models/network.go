// Project gonder.
// Author Supme
// Copyright Supme 2016
// License http://opensource.org/licenses/MIT MIT License	
//
//  THE SOFTWARE AND DOCUMENTATION ARE PROVIDED "AS IS" WITHOUT WARRANTY OF
//  ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
//  IMPLIED WARRANTIES OF MERCHANTABILITY AND/OR FITNESS FOR A PARTICULAR
//  PURPOSE.
//
// Please see the License.txt file for more information.
//
package models

import (
	"sync"
	"net"
	"time"
	"log"
	"fmt"
	"strings"
	"strconv"
)

type mxStor struct {
	records []*net.MX
	update time.Time
}

var mx = struct {
	stor map[string]mxStor
	sync.Mutex
} {
	stor: make(map[string]mxStor),
}

func DomainGetMX(domain string) ([]*net.MX, error) {
	var (
		record []*net.MX
		err error
	)

	if Config.DnsCache {
		mx.Lock()
		defer mx.Unlock()
		if _, ok := mx.stor[domain]; !ok || time.Since(mx.stor[domain].update) > 15 * time.Minute {
			record, err = net.LookupMX(domain)
			if err == nil {
				mx.stor[domain] = mxStor{
					records: record,
					update:time.Now(),
				}
			}
		} else {
			record = mx.stor[domain].records
		}
	} else {
		record, err = net.LookupMX(domain)
	}

	return record, err
}


type (
	profileData struct {
		iface, host string
		streamNow, streamMax int
		lastUpdate time.Time
	}
)

var (
	profileStor = map[int]profileData{}
	profileGroup = map[int]int{}
	profileMutex sync.Mutex
)

func ProfileNext(id int) (int, string, string){
	var res profileData

	profileMutex.Lock()

	// Если есть в массиве, недавно обновлялось
	_, ok := profileStor[id]
	if ok && time.Since(profileStor[id].lastUpdate) < 60 * time.Second {

		// Если это группа кампаний
		if strings.ToLower(strings.TrimSpace(profileStor[id].host)) == "group" {
			if _, gok := profileGroup[id]; !gok {
				profileGroup[id] = 0
			}
			gIfaces := strings.Split(profileStor[id].iface, ",")
			if profileGroup[id] + 1 > len(gIfaces) {
				profileGroup[id] = 0
			}
			i, e := strconv.Atoi(gIfaces[profileGroup[id]])
			if e != nil {
				log.Print(e)
			}
			profileGroup[id]++

			profileMutex.Unlock()
			return ProfileNext(i)
		}

		// Не достигли максимума потоков
		if profileStor[id].streamNow < profileStor[id].streamMax{
			res = profileStor[id]
			res.streamNow++
			profileStor[id] = res
			profileMutex.Unlock()
			return id, res.iface, res.host
		} else {
			// достигли максимума потоков, ждём освобождения
			profileMutex.Unlock()
			for !profileCheck(id) {}
			return ProfileNext(id)
		}
	}

	// В остальных случаях обновляем данные
	err := Db.QueryRow("SELECT `iface`,`host`,`stream` FROM `profile` WHERE `id`=?", id).Scan(&res.iface, &res.host, &res.streamMax)
	if err != nil {
		log.Print(err)
	}

	// если уже существовало, сохраним
	if ok {
		res.streamNow = profileStor[id].streamNow
	}

	res.lastUpdate = time.Now()

	profileStor[id] = res

	// и повторяем действие
	profileMutex.Unlock()
	return ProfileNext(id)

}

func profileCheck(id int) bool {
	var free bool
	profileMutex.Lock()
	free = profileStor[id].streamNow < profileStor[id].streamMax
	profileMutex.Unlock()
	return free
}

func ProfileFree(id int)  {
	var res profileData

	profileMutex.Lock()
	res = profileStor[id]
	res.streamNow--
	profileStor[id] = res
	profileMutex.Unlock()
	fmt.Println("Profile id =", id, " connection count =", res.streamNow)
}


/*
func ProfileUpdate() {

	var profile profileType
	var id int

	query, err := Db.Query("SELECT `id`,`iface`,`host`,`stream` FROM `profile`")
	if err != nil {
		log.Print(err)
	}
	defer query.Close()

	profileMutex.Lock()
	defer profileMutex.Unlock()
	for query.Next() {
		profile.iface, profile.host = "", ""
		id, profile.stream  = 0, 0
		err = query.Scan(&id, &profile.iface, &profile.host, &profile.stream)
		if err != nil {
			log.Print(err)
		}
		profileData[id] = profile
		if _, ok := profileCount[id]; !ok {
			profileCount[id] = 0
		}
	}

	fmt.Println(profileData)
}

func ProfileNext(id int) (iface, host string) {
	var free bool
	for iface, host, free = profileNext(id); free;  {
	}
	profileCount[id]++
	fmt.Println("Profile id =", id, " get next connection count=", profileCount[id], " host = ", profileData[id].host)
	return
}

func profileNext(id int) (iface, host string, free bool) {
	profileMutex.Lock()
	defer profileMutex.Unlock()
	iface, host, free = profileData[id].iface, profileData[id].host, profileCount[id] < profileData[id].stream
	return
}

func ProfileFree(id int)  {
	profileMutex.Lock()
	defer profileMutex.Unlock()
	profileCount[id]--
	fmt.Println("Profile id =", id, " free connection count =", profileCount[id])
}
*/