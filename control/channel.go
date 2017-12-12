package control

// Channel represents a tv channel
type Channel struct {
	ChannelName   string
	ChannelNumber int
	IsScrambled   bool
	IsHdtv        bool
	tv            *LgTv
}

// Watch switches the TV to this channel
func (ch *Channel) Watch() error {
	return ch.tv.SetChannel(ch.ChannelNumber)
}
