package systems

type Stats struct {
	STR int
	INT int
	SPD int
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) LevelUp(stats Stats) {
	s.STR += stats.STR
	s.INT += stats.INT
	s.SPD += stats.SPD
}

func (s *Stats) IncreaseSTR(str int) {
	s.STR += str
}

func (s *Stats) IncreaseINT(intl int) {
	s.INT += intl
}

func (s *Stats) IncreaseSPD(spd int) {
	s.SPD += spd
}

func (s *Stats) DecreaseSTR(str int) {
	s.STR -= str
}

func (s *Stats) DecreaseINT(intl int) {
	s.INT -= intl
}

func (s *Stats) DecreaseSPD(spd int) {
	s.SPD -= spd
}

func (s *Stats) Reset() {
	s.STR = 0
	s.INT = 0
	s.SPD = 0
}
