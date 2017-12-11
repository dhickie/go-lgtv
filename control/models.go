package control

import "time"

// Channel represents a tv channel
type Channel struct {
	ChannelName   string
	ChannelNumber int
	IsScrambled   bool
	IsHdtv        bool
}

// ChannelProgramList is the list of programs broadcast on a particular channel
type ChannelProgramList struct {
	Channel  Channel
	Programs []Program
}

// Program is a particular program broadcast on tv
type Program struct {
	Name      string
	Genre     string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}