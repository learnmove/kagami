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

// Package client contains various data structures and utilities related to individual
// maplestory clients that are currently connected to the login server
package client

import (
	"errors"
	"fmt"
	"net"
	"time"
)

import "github.com/Francesco149/kagami/common"

// Possible player statuses used by SetStatus and Status
const (
	LoggedIn    = 0x00
	SetPin      = 0x01
	CheckPin    = 0x04
	SetGender   = 0x0A
	PinSelect   = 0x0B
	AskPin      = 0x0C
	NotLoggedIn = 0x0D
)

// client.Connection represent a MapleStory client connected to the login server.
// It's a wrapper around EncryptedConnection specialized for MapleStory clients.
// It caches various data from the database such as gm level, char delete password and so on.
type Connection struct {
	*common.EncryptedConnection        // underlying encrypted connection
	admin                       bool   // true if the user is an admin
	worldid                     byte   // numeric world id
	channel                     byte   // channel number
	userid                      uint32 // user id in the database
	invalidLogins               uint32 // invalid login counter
	gmLevel                     int32  // gm level
	charDeletePassword          uint32 // char deletion password (birthday YYYYMMDD)
	accountCreationTime         int64  // unix timestamp of when then account was created
	status                      byte   // current status (offline, logged in, ...)
}

// NewConnection initializes and returns an encrypted connection to a MapleStory client
func NewConnection(basecon net.Conn, testserver bool) *Connection {
	return &Connection{
		EncryptedConnection: common.NewEncryptedConnection(basecon, testserver), // base class
		admin:               false,
		worldid:             0xFF,
		channel:             0,
		userid:              0,
		invalidLogins:       0,
		gmLevel:             0,
		charDeletePassword:  11111111,
		accountCreationTime: 0,
		status:              NotLoggedIn,
	}
}

func (c *Connection) String() string {
	return fmt.Sprintf(""+
		"\n%v:{\n"+
		"\tadmin: %v\n"+
		"\tworldid: %v\n"+
		"\tchannel: %v\n"+
		"\tuserid: %v\n"+
		"\tinvalidLogins: %v\n"+
		"\tgmLevel: %v\n"+
		"\tcharDeletePassword: %v\n"+
		"\taccountCreationTime: %v\n"+
		"\tstatus: %s\n"+
		"}\n",
		c.Conn().RemoteAddr(), c.Admin(), c.WorldId(), c.Channel(), c.Id(),
		c.InvalidLogins(), c.GmLevel(), c.CharDeletePassword(),
		time.Unix(c.AccountCreationTime(), 0), c.PlayerStatusString())
}

// Admin returns whether the user is an admin or not
func (c *Connection) Admin() bool {
	return c.admin
}

// SetAdmin sets the client's admin status
func (c *Connection) SetAdmin(admin bool) {
	c.admin = admin
}

// WorldId returns the client's current world id
func (c *Connection) WorldId() byte {
	return c.worldid
}

// SetWorldId sets the client's current world id
func (c *Connection) SetWorldId(worldid byte) {
	c.worldid = worldid
}

// Channel returns the client's current channel
func (c *Connection) Channel() byte {
	return c.channel
}

// SetChannel sets the client's current channel
func (c *Connection) SetChannel(channel byte) {
	c.channel = channel
}

// Id returns the user's id
func (c *Connection) Id() uint32 {
	return c.userid
}

// SetId sets the user id for this connection
func (c *Connection) SetId(id uint32) {
	c.userid = id
}

// InvalidLogins() returns the failed login count
func (c *Connection) InvalidLogins() uint32 {
	return c.invalidLogins
}

// RegisterInvalidLogin increases the invalid login counter
func (c *Connection) RegisterInvalidLogin() {
	c.invalidLogins++
}

// GmLevel returns the user's gm level
func (c *Connection) GmLevel() int32 {
	return c.gmLevel
}

// SetGmLevel sets the user's gm level
func (c *Connection) SetGmLevel(gmLevel int32) {
	c.gmLevel = gmLevel
}

// CharDeletePassword returns the player's birthday code
func (c *Connection) CharDeletePassword() uint32 {
	return c.charDeletePassword
}

// SetCharDeletePassword sets the player's birthday code
func (c *Connection) SetCharDeletePassword(charDeletePassword uint32) {
	c.charDeletePassword = charDeletePassword
}

// AccountCreationTime returns the timestamp of when the account was created
func (c *Connection) AccountCreationTime() int64 {
	return c.accountCreationTime
}

// SetAccountCreationTime sets the timestamp of when the account was created
func (c *Connection) SetAccountCreationTime(accountCreationTime int64) {
	c.accountCreationTime = accountCreationTime
}

/*
   PlayerStatus returns the current status of this client (logged in, offline...)

   Possible values:
   LoggedIn = 0x00
   SetPin = 0x01
   CheckPin = 0x04
   SetGender = 0x0A
   PinSelect = 0x0B
   AskPin = 0x0C
   NotLoggedIn = 0x0D
*/
func (c *Connection) PlayerStatus() byte {
	return c.status
}

/*
   SetPlayerStatus sets the current status of this client (logged in, offline...)

   Possible values:
   LoggedIn = 0x00
   SetPin = 0x01
   CheckPin = 0x04
   SetGender = 0x0A
   PinSelect = 0x0B
   AskPin = 0x0C
   NotLoggedIn = 0x0D
*/
func (c *Connection) SetPlayerStatus(status byte) {
	c.status = status
}

func (c *Connection) PlayerStatusString() string {
	switch c.PlayerStatus() {
	case LoggedIn:
		return "LoggedIn"
	case SetPin:
		return "SetPin"
	case CheckPin:
		return "CheckPin"
	case SetGender:
		return "SetGender"
	case PinSelect:
		return "PinSelect"
	case AskPin:
		return "AskPin"
	case NotLoggedIn:
		return "NotLoggedIn"
	}

	panic(errors.New(fmt.Sprintf("Invalid player status %v", c.PlayerStatus())))
	return "Invalid status!"
}
