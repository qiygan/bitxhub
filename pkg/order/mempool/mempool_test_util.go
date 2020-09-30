package mempool

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/log"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	raftproto "github.com/meshplus/bitxhub/pkg/order/etcdraft/proto"
	"github.com/meshplus/bitxhub/pkg/storage/leveldb"
)

var (
	InterchainContractAddr = types.String2Address("000000000000000000000000000000000000000a")
)

const (
	DefaultTestBatchSize = uint64(4)
	DefaultTestTxSetSize = uint64(2)
	DefaultTestTxSetTick = 1 * time.Millisecond
	LevelDBDir           = "test-db"
)

func mockMempool() (MemPool, chan *raftproto.Ready) {
	config := &Config{
		ID:             1,
		ChainHeight:    1,
		BatchSize:      DefaultTestBatchSize,
		PoolSize:       DefaultPoolSize,
		TxSliceSize:    DefaultTestTxSetSize,
		BatchTick:      DefaultBatchTick,
		FetchTimeout:   DefaultFetchTxnTimeout,
		TxSliceTimeout: DefaultTestTxSetTick,
		Logger:         log.NewWithModule("consensus"),
	}
	db, _ := leveldb.New(LevelDBDir)
	proposalC := make(chan *raftproto.Ready)
	mempool := newMempoolImpl(config, db, proposalC)
	return mempool, proposalC
}

func mockMempoolImpl() (*mempoolImpl, chan *raftproto.Ready) {
	config := &Config{
		ID:             1,
		ChainHeight:    1,
		BatchSize:      DefaultTestBatchSize,
		PoolSize:       DefaultPoolSize,
		TxSliceSize:    DefaultTestTxSetSize,
		BatchTick:      DefaultBatchTick,
		FetchTimeout:   DefaultFetchTxnTimeout,
		TxSliceTimeout: DefaultTxSetTick,
		Logger:         log.NewWithModule("consensus"),
	}
	db, _ := leveldb.New(LevelDBDir)
	proposalC := make(chan *raftproto.Ready)
	mempool := newMempoolImpl(config, db, proposalC)
	return mempool, proposalC
}

func genPrivKey() crypto.PrivateKey {
	privKey, _ := asym.GenerateKeyPair(crypto.Secp256k1)
	return privKey
}

func constructTx(nonce uint64, privKey *crypto.PrivateKey) *pb.Transaction {
	var privK crypto.PrivateKey
	if privKey == nil {
		privK = genPrivKey()
	}
	privK = *privKey
	pubKey := privK.PublicKey()
	addr, _ := pubKey.Address()
	tx := &pb.Transaction{Nonce: nonce}
	tx.Timestamp = time.Now().UnixNano()
	tx.From = addr
	sig, _ := privK.Sign(tx.SignHash().Bytes())
	tx.Signature = sig
	tx.TransactionHash = tx.Hash()
	return tx
}

func constructIBTPTx(nonce uint64, privKey *crypto.PrivateKey) *pb.Transaction {
	var privK crypto.PrivateKey
	if privKey == nil {
		privK = genPrivKey()
	}
	privK = *privKey
	pubKey := privK.PublicKey()
	from, _ := pubKey.Address()
	to := from.Hex()
	ibtp := mockIBTP(from.Hex(), to ,nonce)
	tx := mockInterChainTx(ibtp)
	tx.Timestamp = time.Now().UnixNano()
	sig, _ := privK.Sign(tx.SignHash().Bytes())
	tx.Signature = sig
	tx.TransactionHash = tx.Hash()
	return tx
}

func cleanTestData() bool {
	err := os.RemoveAll(LevelDBDir)
	if err != nil {
		return false
	}
	return true
}

func mockInterChainTx(ibtp *pb.IBTP) *pb.Transaction {
	ib, _ := ibtp.Marshal()
	ipd := &pb.InvokePayload{
		Method: "HandleIBTP",
		Args:   []*pb.Arg{{Value: ib}},
	}
	pd, _ := ipd.Marshal()
	data := &pb.TransactionData{
		VmType:  pb.TransactionData_BVM,
		Type:    pb.TransactionData_INVOKE,
		Payload: pd,
	}
	return &pb.Transaction{
		To:    InterchainContractAddr,
		Nonce: ibtp.Index,
		Data:  data,
		Extra: []byte(fmt.Sprintf("%s-%s-%d", ibtp.From, ibtp.To, ibtp.Type)),
	}
}

func mockIBTP(from, to string, nonce uint64) *pb.IBTP {
	content := pb.Content{
		SrcContractId: from,
		DstContractId: from,
		Func:          "interchainget",
		Args:          [][]byte{[]byte("Alice"), []byte("10")},
	}
	bytes, _ := content.Marshal()
	ibtppd, _ := json.Marshal(pb.Payload{
		Encrypted: false,
		Content:   bytes,
	})
	return &pb.IBTP{
		From:      from,
		To:        to,
		Payload:   ibtppd,
		Index:     nonce,
		Type:      pb.IBTP_INTERCHAIN,
		Timestamp: time.Now().UnixNano(),
	}
}
