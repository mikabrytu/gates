package systems

type Skill struct {
	STR int
	INT int
	SPD int
}

func NewSkill() *Skill {
	return &Skill{}
}

func (s *Skill) LevelUp(skill Skill) {
	s.STR += skill.STR
	s.INT += skill.INT
	s.SPD += skill.SPD
}

func (s *Skill) IncreaseSTR(str int) {
	s.STR += str
}

func (s *Skill) IncreaseINT(intl int) {
	s.INT += intl
}

func (s *Skill) IncreaseSPD(spd int) {
	s.SPD += spd
}

func (s *Skill) DecreaseSTR(str int) {
	s.STR -= str
}

func (s *Skill) DecreaseINT(intl int) {
	s.INT -= intl
}

func (s *Skill) DecreaseSPD(spd int) {
	s.SPD -= spd
}

func (s *Skill) Reset() {
	s.STR = 0
	s.INT = 0
	s.SPD = 0
}
