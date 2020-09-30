package mempool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecvTransaction(t *testing.T){
	ast := assert.New(t)
	mempool, _ := mockMempool()
	err := mempool.Start()
	defer cleanTestData()
	ast.Nil(err)

	//privKey1 := genPrivKey()
	//account1,_ := privKey1.PublicKey().Address()
	//tx := constructTx(uint64(1),&privKey1)
	//go func() {
	//	err := mempool.RecvTransaction(tx)
	//	ast.Nil(err)
	//}()
	//time.Sleep(2* time.Millisecond)
	//pendingNonce := mempool.GetPendingNonceByAccount(account1.Hex())
	//ast.Equal(uint64(2), pendingNonce)
}



