package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/v12"
)


type Block struct {
	Index     int
	Timestamp string
	Data      string
	PrevHash  string
	Hash      string
	Miner     string
	Reward    float64
}
var Blockchain []Block

const APIKey = "your_secret_key"

const (
	MongoDBHost = "localhost:27017"
	DBName      = "blockchain"
	BlockchainCollection = "blocks"
)

var session *mgo.Session

type Vote struct {
	ID        bson.ObjectId `bson:"_id"`
	Proposal  string        `json:"proposal"`
	Voter     string        `json:"voter"`
	Approved  bool          `json:"approved"`
	Timestamp string        `json:"timestamp"`
}
// calculateHash calculates the hash of a block based on its index, timestamp, data, and the previous block's hash.
func calculateHash(block Block) string {
	data := fmt.Sprintf("%d%s%s%s%s", block.Index, block.Timestamp, block.Data, block.PrevHash, block.Miner)
	hashBytes := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hashBytes[:])
}

// createNewBlock creates a new block for the blockchain.
func createNewBlock(prevBlock Block, data string, miner string, reward float64) Block {
	newBlock := Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
		Miner:     miner,
		Reward:    reward,
	}
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}
// Initialize MongoDB session
func initMongoDB() {
	var err error
	session, err = mgo.Dial(MongoDBHost)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	session.SetMode(mgo.Monotonic, true)
}

func main() {
	app := iris.New()

	// Middleware for API key authentication.
	apiKeyMiddleware := func(ctx iris.Context) {
		apiKey := ctx.GetHeader("X-API-Key")
		if apiKey != APIKey {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(iris.Map{"error": "Unauthorized"})
			return
		}
		ctx.Next()
	}
	initMongoDB()
	
	session.DB(DBName).C(BlockchainCollection).Find(nil).Sort("index").All(&Blockchain)

	// Create and add new blocks to the blockchain.
	app.Get("/blocks", func(ctx iris.Context) {
		ctx.JSON(Blockchain)
	})

	app.Post("/addBlock", apiKeyMiddleware, func(ctx iris.Context) {
		var data struct {
			Data   string  `json:"data"`
			Miner  string  `json:"miner"`
			Reward float64 `json:"reward"`
		}

		if err := ctx.ReadJSON(&data); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": "Invalid request body"})
			return
		}

		lastBlock := Blockchain[len(Blockchain)-1]
		newBlock := createNewBlock(lastBlock, data.Data, data.Miner, data.Reward)
		Blockchain = append(Blockchain, newBlock)

		// Save the new block to MongoDB.
		collection := session.DB(DBName).C(BlockchainCollection)
		if err := collection.Insert(newBlock); err != nil {
			log.Printf("Failed to insert block into MongoDB: %v", err)
		}

		ctx.JSON(iris.Map{"message": "Block added successfully"})
	})

	// Governance-related routes.
	app.Post("/propose", apiKeyMiddleware, func(ctx iris.Context) {
		var voteData struct {
			Proposal string `json:"proposal"`
		}

		if err := ctx.ReadJSON(&voteData); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": "Invalid request body"})
			return
		}

		// Create a new proposal and set it as unapproved by default.
		proposal := Vote{
			ID:        bson.NewObjectId(),
			Proposal:  voteData.Proposal,
			Voter:     "",
			Approved:  false,
			Timestamp: time.Now().String(),
		}

		// Save the proposal to MongoDB.
		collection := session.DB(DBName).C("proposals")
		if err := collection.Insert(proposal); err != nil {
			log.Printf("Failed to insert proposal into MongoDB: %v", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to create proposal"})
			return
		}

		ctx.JSON(iris.Map{"message": "Proposal created successfully"})
	})

	app.Post("/vote", apiKeyMiddleware, func(ctx iris.Context) {
		var voteData struct {
			ProposalID bson.ObjectId `json:"proposalID"`
			Voter      string        `json:"voter"`
			Approve    bool          `json:"approve"`
		}
		if err := ctx.ReadJSON(&voteData); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": "Invalid request body"})
			return
		}
		// Check if the proposal exists.
		var proposal Vote
		collection := session.DB(DBName).C("proposals")
		if err := collection.FindId(voteData.ProposalID).One(&proposal); err != nil {
			log.Printf("Failed to find proposal in MongoDB: %v", err)
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": "Proposal not found"})
			return
		}
		// Check if the voter has already voted on this proposal.
		if proposal.Voter != "" {
			ctx.StatusCode(iris.StatusConflict)
			ctx.JSON(iris.Map{"error": "You've already voted on this proposal"})
			return
		}
		// Update the proposal with the voter's decision.
		proposal.Voter = voteData.Voter
		proposal.Approved = voteData.Approve

		// Save the updated proposal to MongoDB.
		if err := collection.UpdateId(proposal.ID, proposal); err != nil {
			log.Printf("Failed to update proposal in MongoDB: %v", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to update proposal"})
			return
		}
		ctx.JSON(iris.Map{"message": "Vote recorded successfully"})
	})

	// Run the Iris application on port 8080.
	app.Run(iris.Addr(":8080"))
}
