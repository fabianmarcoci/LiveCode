package config

var JWTSecret string

func Init(jwtSecret string) {
	JWTSecret = jwtSecret
}
