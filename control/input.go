package control

// Input represents an external input to the TV
type Input struct {
	Label string
	ID    string
	tv    *LgTv
}

// Switch switches the TV to this input
func (i *Input) Switch() error {
	return i.tv.SwitchInput(i.ID)
}
