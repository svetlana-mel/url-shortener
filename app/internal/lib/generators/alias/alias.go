package alias

import (
	"math/rand"
)

func GenerateAlias(aliasLen int) string {
	symbols := []rune("AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789")

	// rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	res := make([]rune, aliasLen)

	for i := 0; i < aliasLen; i++ {
		res[i] = symbols[rand.Intn(len(symbols))]
	}

	return string(res)
}
