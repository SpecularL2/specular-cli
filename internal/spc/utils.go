package spc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	TestGasLimit             = 21000      // 21 Kwei
	TestMaxPriorityFeePerGas = 2000000000 // 2 Gwei
	TestMaxFeePerGas         = 2000000000 // 20 Gwei
)

var (
	Wei      = big.NewInt(1)                   //nolint:nolintlint,gochecknoglobals,gomnd
	Kwei     = big.NewInt(1000)                //nolint:nolintlint,gochecknoglobals,gomnd
	Mwei     = big.NewInt(1000000)             //nolint:nolintlint,gochecknoglobals,gomnd
	Gwei     = big.NewInt(1000000000)          //nolint:nolintlint,gochecknoglobals,gomnd
	Microeth = big.NewInt(1000000000000)       //nolint:nolintlint,gochecknoglobals,gomnd
	Millieth = big.NewInt(1000000000000000)    //nolint:nolintlint,gochecknoglobals,gomnd
	Eth      = big.NewInt(1000000000000000000) //nolint:nolintlint,gochecknoglobals,gomnd
	Eth_1K   = Eth.Mul(Kwei, Eth)              //nolint:nolintlint,gochecknoglobals,gomnd
)

func PrivateKeyHexToAccountAddress(privateKeyHex string) (common.Address, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address, nil
}

func CheckBalance(client *ethclient.Client, address common.Address) (*big.Int, error) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func NewTransaction(
	client *ethclient.Client,
	privateKeyHex string,
	value *big.Int,
	toAddress common.Address,
) (*types.Transaction, error) {
	fromAddress, err := PrivateKeyHexToAccountAddress(privateKeyHex)
	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(TestGasLimit)
	tip := big.NewInt(TestMaxPriorityFeePerGas)
	feeCap := big.NewInt(TestMaxFeePerGas)

	var data []byte

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tip,
		GasFeeCap: feeCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func SendTransaction(client *ethclient.Client, signedTx *types.Transaction) (*types.Transaction, error) {
	err := client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
