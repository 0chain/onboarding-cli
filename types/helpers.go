package types

func NewMpks() *Mpks {
	return &Mpks{Mpks: make(map[string]*MPK)}
}
