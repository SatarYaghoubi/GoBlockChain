# GoBlockchain

GoBlockchain is a basic blockchain application implemented in Go (Golang). It includes features like persistence using MongoDB, mining rewards, and a simple on-chain governance system. This application allows users to create, explore, and govern a blockchain network.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Blockchain Governance](#blockchain-governance)


## Features

- Basic blockchain implementation in Go.
- Persistence using MongoDB for data storage.
- Mining rewards for adding new blocks.
- On-chain governance system for proposing and voting on changes.

## Getting Started

### Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (Golang) installed on your system.
- MongoDB installed and running (or provide MongoDB connection details in the configuration).

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/SatarYaghoubi/GoBlockchain.git
   cd GoBlockchain
   nano the main.go
   configure the database
   and use you own secretkey
   go run main.go

# Blockchain Usage Guide

Once the application is running, you can use it to explore and interact with the blockchain. Here are some common tasks:

- **View Blockchain:** Access the blockchain data at [http://localhost:8080/blocks](http://localhost:8080/blocks).

- **Add a Block:** Add a new block to the blockchain by making a POST request to [http://localhost:8080/addBlock](http://localhost:8080/addBlock) with the appropriate JSON payload.

- **Create a Governance Proposal:** Propose changes to the blockchain using the `/propose` endpoint.

- **Vote on a Proposal:** Cast your vote on a proposal using the `/vote` endpoint.

Please refer to the API Endpoints section for more details on the available endpoints and request/response formats.

## API Endpoints

- **GET /blocks:** Retrieve the current blockchain.

- **POST /addBlock:** Add a new block to the blockchain.

- **POST /propose:** Create a new governance proposal.

- **POST /vote:** Cast your vote on a governance proposal.

Detailed information about each endpoint and the expected request/response formats can be found in the API documentation.

## Blockchain Governance

This blockchain application includes a basic on-chain governance system. Users can propose changes and vote on them. Please note that this governance system is simplified and may require further development for real-world use cases.
   
