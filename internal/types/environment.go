package types

type Environment int

const (
	Development Environment = iota
	CICD
	Production
)

func StringToEnv(s string) Environment {
	switch s {
	case "development":
		return Development
	case "cicd":
		return CICD
	case "production":
		return Production
	default:
		return Development
	}
}

func (e Environment) IsProduction() bool {
	return e == Production
}

func (e Environment) IsCICD() bool {
	return e == CICD
}

func (e Environment) IsDevelopment() bool {
	return e == Development
}
