package api

import (
	"github.com/ryt-io/ryt-v2/api/admin"
	"github.com/ryt-io/ryt-v2/api/info"
	"github.com/ryt-io/ryt-v2/indexer"
	"github.com/ryt-io/ryt-v2/runtime/avm"
	"github.com/ryt-io/ryt-v2/runtime/platformvm"
	evmclient "github.com/ryt-io/ryt-v2/graft/ethereum/plugin/evm/client"
)

// Issues API calls to a node
// TODO: byzantine api. check if appropriate. improve implementation.
type Client interface {
	PChainAPI() *platformvm.Client
	XChainAPI() *avm.Client
	XChainWalletAPI() *avm.WalletClient
	CChainAPI() evmclient.Client
	CChainEthAPI() EthClient // ethclient websocket wrapper that adds mutexed calls, and lazy conn init (on first call)
	InfoAPI() *info.Client
	HealthAPI() HealthClient
	AdminAPI() *admin.Client
	PChainIndexAPI() *indexer.Client
	CChainIndexAPI() *indexer.Client
	// TODO add methods
}
