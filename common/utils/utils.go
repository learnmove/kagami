/*
   Copyright 2014 Franc[e]sco (lolisamurai@tfwno.gf)
   This file is part of kagami.
   kagami is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   kagami is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with kagami. If not, see <http://www.gnu.org/licenses/>.
*/

// Package utils contains various utility functions shared by multiple files
package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"strconv"
	"strings"
)

import "github.com/Francesco149/kagami/common/consts"

// HashPassword returns a salted sha-512 hash of the given password
func HashPassword(password, salt string) string {
	hasher := sha512.New()
	saltedpassword := fmt.Sprintf("%sIREALLYLIKELOLIS%s", password, salt)
	hasher.Write([]byte(saltedpassword))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// MakeSalt generates a random string of fixed length that will be used as a password salt
func MakeSalt() string {
	var salt [consts.SaltLength]byte
	rand.Read(salt[:])

	// make it a valid string
	for i := 0; i < consts.SaltLength; i++ {
		salt[i] %= 93 // characters will be between ascii 33 and ascii 126
		salt[i] += 33
	}

	return string(salt[:])
}

// UnixToTempBanTimestamp converts a unix timestamp (in seconds) to a temp ban timestamp
// (number of 100-ns intervals since 1/1/1601)
func UnixToTempBanTimestamp(unixSeconds int64) uint64 {
	// this should be the offset between the unix timestamp and this weird korean timestamp
	const offset = 116444736000000000
	millisecs := uint64(unixSeconds * 1000)
	nano100 := millisecs * 10000 // number of 100-ns intervals
	return nano100 + offset
}

// UnixToTempBanTimestamp converts a unix timestamp (in seconds) to a item timestamp
func UnixToItemTimestamp(unixSeconds int64) uint64 {
	const realYear2000 = 946681229830
	const itemYear2000 = 1085019342
	millisecs := uint64(unixSeconds * 1000)
	time := (millisecs - realYear2000) / 1000 / 60
	// what the fuck
	return uint64(float64(time)*35.762787) - itemYear2000
}

// UnixToQuestTimestamp converts a unix timestamp (in seconds) to a quest timestamp
func UnixToQuestTimestamp(unixSeconds int64) uint64 {
	const questUnixAge = 27111908
	millisecs := uint64(unixSeconds * 1000)
	time := millisecs / 1000 / 60
	// what the fuck
	return uint64(float64(time)*0.1396987) + questUnixAge
}

// RemoteAddrToIp returns the ip of a ip:port string
func RemoteAddrToIp(addr string) string {
	return strings.Split(addr, ":")[0]
}

// RemoteAddrToPort returns the port of a ip:port string
func RemoteAddrToPort(addr string) string {
	return strings.Split(addr, ":")[1]
}

// RemoteAddrToBytes converts a xxx.xxx.xxx.xxx:port string to an array of bytes that contains the ip
func RemoteAddrToBytes(addr string) (res []byte) {
	addr = RemoteAddrToIp(addr)
	split := strings.Split(addr, ".")
	res = make([]byte, len(split))

	for i := 0; i < len(split); i++ {
		tmp, err := strconv.Atoi(split[i])
		if err != nil {
			return nil
		}

		if tmp > 255 || tmp < 0 {
			return nil
		}

		res[i] = byte(tmp)
	}

	return
}

// BytesToIpString converts a byte array that represents an ip to a string
func BytesToIpString(ip []byte) string {
	if len(ip) != 4 {
		return "ipv6 not supported"
	}
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// MakeWarning formats a warning message
func MakeWarning(msg string) string {
	return fmt.Sprint("\n",
		"******* /!\\ WARNING /!\\ *******\n",
		msg, "\n",
		"*******************************\n",
		"\n")
}

// MakeError formats an error message
func MakeError(msg string) string {
	return fmt.Sprint("\n",
		"******** /!\\ ERROR /!\\ ********\n",
		msg, "\n",
		"*******************************\n",
		"\n")
}

// AnyNil returns true if any of the given values is nil
func AnyNil(a ...interface{}) bool {
	for _, val := range a {
		if val == nil {
			return true
		}
	}

	return false
}

// A Pair holds two values of any type that can be
// casted back to their original type through type assertions
type Pair struct {
	First, Second interface{}
}