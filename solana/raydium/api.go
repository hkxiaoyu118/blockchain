package raydium

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PoolInfo struct {
	Name       string `json:"name"`
	Official   []Pool `json:"official"`
	UnOfficial []Pool `json:"unOfficial"`
}

//type Pool struct {
//	ID                 string `json:"id"`
//	BaseMint           string `json:"baseMint"`
//	QuoteMint          string `json:"quoteMint"`
//	LpMint             string `json:"lpMint"`
//	BaseDecimals       int    `json:"baseDecimals"`
//	QuoteDecimals      int    `json:"quoteDecimals"`
//	LpDecimals         int    `json:"lpDecimals"`
//	Version            int    `json:"version"`
//	ProgramId          string `json:"programId"`
//	Authority          string `json:"authority"`
//	OpenOrders         string `json:"openOrders"`
//	TargetOrders       string `json:"targetOrders"`
//	BaseVault          string `json:"baseVault"`
//	QuoteVault         string `json:"quoteVault"`
//	WithdrawQueue      string `json:"withdrawQueue"`
//	LpVault            string `json:"lpVault"`
//	MarketVersion      int    `json:"marketVersion"`
//	MarketProgramId    string `json:"marketProgramId"`
//	MarketId           string `json:"marketId"`
//	MarketAuthority    string `json:"marketAuthority"`
//	MarketBaseVault    string `json:"marketBaseVault"`
//	MarketQuoteVault   string `json:"marketQuoteVault"`
//	MarketBids         string `json:"marketBids"`
//	MarketAsks         string `json:"marketAsks"`
//	MarketEventQueue   string `json:"marketEventQueue"`
//	LookupTableAccount string `json:"lookupTableAccount"`
//}

type Pool struct {
	ID                 string `gorm:"type:varchar(255);primaryKey" json:"id"`
	BaseMint           string `gorm:"type:varchar(255)" json:"baseMint"`
	QuoteMint          string `gorm:"type:varchar(255)" json:"quoteMint"`
	LpMint             string `gorm:"type:varchar(255)" json:"lpMint"`
	BaseDecimals       int    `gorm:"type:bigint" json:"baseDecimals"`
	QuoteDecimals      int    `gorm:"type:bigint" json:"quoteDecimals"`
	LpDecimals         int    `gorm:"type:bigint" json:"lpDecimals"`
	Version            int    `gorm:"type:bigint" json:"version"`
	ProgramId          string `gorm:"type:varchar(255)" json:"programId"`
	Authority          string `gorm:"type:varchar(255)" json:"authority"`
	OpenOrders         string `gorm:"type:varchar(255)" json:"openOrders"`
	TargetOrders       string `gorm:"type:varchar(255)" json:"targetOrders"`
	BaseVault          string `gorm:"type:varchar(255)" json:"baseVault"`
	QuoteVault         string `gorm:"type:varchar(255)" json:"quoteVault"`
	WithdrawQueue      string `gorm:"type:varchar(255)" json:"withdrawQueue"`
	LpVault            string `gorm:"type:varchar(255)" json:"lpVault"`
	MarketVersion      int    `gorm:"type:bigint" json:"marketVersion"`
	MarketProgramId    string `gorm:"type:varchar(255)" json:"marketProgramId"`
	MarketId           string `gorm:"type:varchar(255)" json:"marketId"`
	MarketAuthority    string `gorm:"type:varchar(255)" json:"marketAuthority"`
	MarketBaseVault    string `gorm:"type:varchar(255)" json:"marketBaseVault"`
	MarketQuoteVault   string `gorm:"type:varchar(255)" json:"marketQuoteVault"`
	MarketBids         string `gorm:"type:varchar(255)" json:"marketBids"`
	MarketAsks         string `gorm:"type:varchar(255)" json:"marketAsks"`
	MarketEventQueue   string `gorm:"type:varchar(255)" json:"marketEventQueue"`
	LookupTableAccount string `gorm:"type:varchar(255)" json:"lookupTableAccount"`
}

// GetMainnetPoolsInfo todo: 这个破方法，非常非常非常慢，所以最好是使用mysql/sqlite/redis等数据库进行保存，方面查询，不然就是噩梦
func GetMainnetPoolsInfo() (*PoolInfo, error) {
	response, err := http.Get("https://api.raydium.io/v2/sdk/liquidity/mainnet.json")
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get raydium pools; status_code=%d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var poolInfo PoolInfo
	err = json.Unmarshal(body, &poolInfo)
	if err != nil {
		return nil, err
	}

	return &poolInfo, nil
}

func FindPoolInfoByID(id string) (*Pool, error) {
	poolsInfo, err := GetMainnetPoolsInfo()
	if err != nil {
		return nil, err
	}
	pools := append(poolsInfo.Official, poolsInfo.UnOfficial...)
	for _, pool := range pools {
		if pool.ID == id {
			return &pool, nil
		}
	}
	return nil, fmt.Errorf("pool was not found by id [%s]", id)
}
