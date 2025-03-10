package stocks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseTickets(t *testing.T) {
	file, err := os.Open("./testdata/nasdaq_full_tickers.json")
	require.NoError(t, err)
	defer file.Close()

	stocks, err := DecodeTickers(file)
	require.NoError(t, err)
	assert.Len(t, stocks, 3922)
	assert.Equal(t, stocks[0].Symbol, "AACBU")
	assert.Equal(t, stocks[0].Lastsale, "$10.04")
}
