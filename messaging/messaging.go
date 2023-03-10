package messaging

import (
	"encoding/json"
	"fmt"
	"proofofaccess/database"
	"proofofaccess/ipfs"
	"proofofaccess/localdata"
	"proofofaccess/pubsub"
	"proofofaccess/validation"
	"strings"
	"time"
)

type Request struct {
	Type string `json:"type"`
	Hash string `json:"hash"`
	CID  string `json:"CID"`
	Seed string `json:"seed"`
	User string `json:"user"`
}

var request Request

var Ping = map[string]bool{}
var ProofRequest = map[string]bool{}

// SendProof
// This is the function that sends the proof of access to the validation node
func SendProof(hash string, seed string, user string) {
	jsonString := `{"type": "ProofOfAccess", "hash":"` + hash + `", "seed":"` + seed + `"}`
	jsonString = strings.TrimSpace(jsonString)
	pubsub.Publish(jsonString, user)
}

// HandleMessage
// This is the function that handles the messages from the pubsub
func HandleMessage(message string, nodeType *int) {
	// JSON decode message
	err := json.Unmarshal([]byte(message), &request)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	fmt.Println("Message received:", message)
	//Handle proof of access response from storage node
	if *nodeType == 1 {
		if request.Type == "ProofOfAccess" {
			HandleProofOfAccess(request)
		}
		if request.Type == "PingPongPong" {
			Ping[request.Hash] = true
			fmt.Println("PingPongPong received")
		}
	}

	//Handle request for proof request from validation node
	if *nodeType == 2 {
		if request.Type == "RequestProof" {
			HandleRequestProof(request)
		}
		if request.Type == "PingPongPing" {
			fmt.Println("PingPongPing received")
			PingPongPong(request.Hash, request.User)
		}
	}
}

// HandleRequestProof
// This is the function that handles the request for proof from the validation node
func HandleRequestProof(request Request) {
	CID := request.CID
	hash := request.Hash
	user := request.User
	if ipfs.IsPinned(CID) == true {
		validationHash := validation.CreatProofHash(hash, CID)
		SendProof(validationHash, hash, user)
	} else {
		SendProof("", hash, user)
	}
}

// HandleProofOfAccess
// This is the function that handles the proof of access response from the storage node
func HandleProofOfAccess(request Request) {
	// Get the start time from the seed
	start := localdata.GetTime(request.Seed)
	fmt.Println("Start time:", start)

	// Get the current time
	elapsed := time.Since(start)

	// Set the elapsed time
	localdata.SetElapsed(request.Seed, elapsed)

	// Get the CID and Seed
	data := database.Read([]byte(request.Seed))
	var message Request
	err := json.Unmarshal([]byte(string(data)), &message)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	CID := message.CID
	ProofRequest[CID] = true
	fmt.Println("CID:", CID)
	Seed := request.Seed
	// Create the proof hash
	validationHash := validation.CreatProofHash(Seed, CID)

	// Check if the proof of access is valid
	if validationHash == request.Hash && elapsed < 500*time.Millisecond {
		fmt.Println("Proof of access is valid")
		fmt.Println(request.Seed)
		localdata.SetStatus(request.Seed, CID, "Valid")
	} else {
		fmt.Println("Proof of access is invalid took too long")
		localdata.SetStatus(request.Seed, CID, "Invalid")
	}
}

func PingPong(hash string, user string) {
	localdata.SaveTime(hash)
	PingPongPing(hash, user)
}

func PingPongPing(hash string, user string) {
	jsonString := `{"type": "PingPongPing", "hash":"` + hash + `", "user":"` + localdata.GetNodeName() + `"}`
	jsonString = strings.TrimSpace(jsonString)
	pubsub.Publish(jsonString, user)
}

func PingPongPong(hash string, user string) {
	jsonString := `{"type": "PingPongPong", "hash":"` + hash + `", "user":"` + user + `"}`
	jsonString = strings.TrimSpace(jsonString)
	pubsub.Publish(jsonString, user)
}
