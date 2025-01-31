package config

const (
	PROD  Environment = "production"
	LOCAL Environment = "local"
)

type Environment string

func (e Environment) IsProduction() bool {
	return e == PROD
}

func (e Environment) IsLocal() bool {
	return e == LOCAL
}

func (e Environment) IsKnown() bool {
	return e.IsProduction() || e.IsLocal()
}

func (e Environment) ShortName() string {
	switch {
	case e.IsProduction():
		return string(PROD)
	case e.IsLocal():
		return string(LOCAL)
	default:
		return ""
	}
}
