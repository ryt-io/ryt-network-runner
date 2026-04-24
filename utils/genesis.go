// Copyright (C) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	avago_upgrade "github.com/ryt-io/ryt-v2/upgrade"
	"github.com/ryt-io/ryt-v2/utils/crypto/bls/signer/localsigner"
	"github.com/ryt-io/ryt-v2/utils/formatting"
	"github.com/ryt-io/ryt-v2/runtime/platformvm/signer"
	"github.com/ryt-io/libevm/common"
	"github.com/ryt-io/libevm/core"
	"github.com/ryt-io/subnet-evm/params"
)

const (
	// difference between unlock schedule locktime and startime in original genesis
	genesisLocktimeStartimeDelta    = 2836800
	hexa0Str                        = "0x0"
	defaultLocalCChainFundedAddress = "8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
	defaultLocalCChainFundedBalance = "0x295BE96E64066972000000"
	allocationCommonEthAddress      = "0xb3d82b1367d362de99ab59a658165aff520cbd4d"
	stakingAddr                     = "X-custom1g65uqn6t77p656w64023nh8nd9updzmxwd59gh"
	walletAddr                      = "X-custom18jma8ppw3nhx5r4ap8clazz0dps7rv5u9xde7p"
	defaultGasLimit                 = uint64(100_000_000) // Gas limit is arbitrary
)

var defaultFundedKeyCChainAmount = new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil)

func generateCchainGenesis() ([]byte, error) {
	cChainBalances := make(core.GenesisAlloc, 1)
	cChainBalances[common.HexToAddress(defaultLocalCChainFundedAddress)] = core.GenesisAccount{
		Balance: defaultFundedKeyCChainAmount,
	}
	chainID := big.NewInt(43112)
	cChainGenesis := &core.Genesis{
		Config:     &params.ChainConfig{ChainID: chainID},            // The rest of the config is set in coreth on VM initialization
		Difficulty: big.NewInt(0),                                    // Difficulty is a mandatory field
		Timestamp:  uint64(avago_upgrade.InitiallyActiveTime.Unix()), // #nosec G115 // This time enables Avalanche upgrades by default
		GasLimit:   defaultGasLimit,
		Alloc:      cChainBalances,
	}
	return json.Marshal(cChainGenesis)
}

func GenerateGenesis(
	networkID uint32,
	nodeKeys []*NodeKeys,
) ([]byte, error) {
	genesisMap := map[string]interface{}{}

	// cchain
	cChainGenesisBytes, err := generateCchainGenesis()
	if err != nil {
		return nil, err
	}
	genesisMap["cChainGenesis"] = string(cChainGenesisBytes)

	// pchain genesis
	genesisMap["networkID"] = networkID
	startTime := time.Now().Unix()
	genesisMap["startTime"] = startTime
	initialStakers := []map[string]interface{}{}

	for _, keys := range nodeKeys {
		nodeID, err := ToNodeID(keys.StakingKey, keys.StakingCert)
		if err != nil {
			return nil, fmt.Errorf("couldn't get node ID: %w", err)
		}
		blsSk, err := localsigner.FromBytes(keys.BlsKey)
		if err != nil {
			return nil, err
		}
		p, err := signer.NewProofOfPossession(blsSk)
		if err != nil {
			return nil, err
		}
		pk, err := formatting.Encode(formatting.HexNC, p.PublicKey[:])
		if err != nil {
			return nil, err
		}
		pop, err := formatting.Encode(formatting.HexNC, p.ProofOfPossession[:])
		if err != nil {
			return nil, err
		}
		initialStaker := map[string]interface{}{
			"delegationFee": 1000000,
			"nodeID":        nodeID,
			"rewardAddress": walletAddr,
			"signer": map[string]interface{}{
				"proofOfPossession": pop,
				"publicKey":         pk,
			},
		}
		initialStakers = append(initialStakers, initialStaker)
	}

	genesisMap["initialStakeDuration"] = 31536000
	genesisMap["initialStakeDurationOffset"] = 5400
	genesisMap["initialStakers"] = initialStakers
	lockTime := startTime + genesisLocktimeStartimeDelta
	allocations := []interface{}{}
	alloc := map[string]interface{}{
		"avaxAddr":      walletAddr,
		"ethAddr":       allocationCommonEthAddress,
		"initialAmount": 300000000000000000,
		"unlockSchedule": []interface{}{
			map[string]interface{}{"amount": 20000000000000000},
			map[string]interface{}{"amount": 10000000000000000, "locktime": lockTime},
		},
	}
	allocations = append(allocations, alloc)
	alloc = map[string]interface{}{
		"avaxAddr":      stakingAddr,
		"ethAddr":       allocationCommonEthAddress,
		"initialAmount": 0,
		"unlockSchedule": []interface{}{
			map[string]interface{}{"amount": 10000000000000000, "locktime": lockTime},
		},
	}
	allocations = append(allocations, alloc)
	genesisMap["allocations"] = allocations
	genesisMap["initialStakedFunds"] = []interface{}{
		stakingAddr,
	}
	genesisMap["message"] = "{{ fun_quote }}"

	return json.MarshalIndent(genesisMap, "", " ")
}
