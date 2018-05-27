package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"github.com/pkg/errors"
	"fmt"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Remark    string
}

var Blockchain []Block

// 计算hash值
func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 根据旧值 生成新值
func generateBlock(oldBlock Block, BPM int) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func checkBlock(oldBlock Block, newBlock Block) (bool, error) {
	if oldBlock.Hash != newBlock.PrevHash {
		return false, errors.New("区块无法连接");
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false, errors.New("hash 被篡改");
	}
	return true, nil;
}

func process(newBlock Block) ([]Block, error) {

	oldBlock := Blockchain[len(Blockchain)-1]
	ret, err := checkBlock(oldBlock, newBlock);
	if !ret && err != nil {
		fmt.Println(err)
		return Blockchain, err;
	}
	Blockchain = append(Blockchain, newBlock)
	return Blockchain, nil;
}

func Run(BPM int) ([]Block, error) {
	block, err := generateBlock(Blockchain[len(Blockchain)-1], BPM);
	if err != nil {
		fmt.Println(err)
		return Blockchain, err;
	}
	return process(block);
}

// 获取最新的区块
func GetLatestBlock() Block {
	return Blockchain[len(Blockchain)-1];
}

// 生成创世区块
func BuildGenesisBlock() Block {
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", "", "创世区块"}
	genesisBlock.Hash = calculateHash(genesisBlock);
	Blockchain = append(Blockchain, genesisBlock)
	return genesisBlock;
}
