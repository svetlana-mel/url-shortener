package alias

import (
	"math/rand"
)

func GenerateAlias(aliasLen int) string {
	symbols := []rune("AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789")

	res := make([]rune, aliasLen)

	var generator rand.Rand

	for i := 0; i < aliasLen; i++ {
		res[i] = symbols[generator.Intn(aliasLen)]
	}

	return string(res)
}
