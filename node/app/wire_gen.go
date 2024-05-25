// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"source.quilibrium.com/quilibrium/monorepo/node/config"
	"source.quilibrium.com/quilibrium/monorepo/node/consensus"
	"source.quilibrium.com/quilibrium/monorepo/node/consensus/master"
	"source.quilibrium.com/quilibrium/monorepo/node/consensus/time"
	"source.quilibrium.com/quilibrium/monorepo/node/crypto"
	"source.quilibrium.com/quilibrium/monorepo/node/execution/intrinsics/ceremony"
	"source.quilibrium.com/quilibrium/monorepo/node/keys"
	"source.quilibrium.com/quilibrium/monorepo/node/p2p"
	"source.quilibrium.com/quilibrium/monorepo/node/protobufs"
	"source.quilibrium.com/quilibrium/monorepo/node/store"
)

// Injectors from wire.go:

func NewDHTNode(configConfig *config.Config) (*DHTNode, error) {
	p2PConfig := configConfig.P2P
	zapLogger := debugLogger()
	blossomSub := p2p.NewBlossomSub(p2PConfig, zapLogger)
	dhtNode, err := newDHTNode(blossomSub)
	if err != nil {
		return nil, err
	}
	return dhtNode, nil
}

func NewDebugNode(configConfig *config.Config, selfTestReport *protobufs.SelfTestReport) (*Node, error) {
	zapLogger := debugLogger()
	dbConfig := configConfig.DB
	pebbleDB := store.NewPebbleDB(dbConfig)
	pebbleClockStore := store.NewPebbleClockStore(pebbleDB, zapLogger)
	keyConfig := configConfig.Key
	fileKeyManager := keys.NewFileKeyManager(keyConfig, zapLogger)
	p2PConfig := configConfig.P2P
	blossomSub := p2p.NewBlossomSub(p2PConfig, zapLogger)
	engineConfig := configConfig.Engine
	wesolowskiFrameProver := crypto.NewWesolowskiFrameProver(zapLogger)
	masterTimeReel := time.NewMasterTimeReel(zapLogger, pebbleClockStore, engineConfig, wesolowskiFrameProver)
	inMemoryPeerInfoManager := p2p.NewInMemoryPeerInfoManager(zapLogger)
	masterClockConsensusEngine := master.NewMasterClockConsensusEngine(engineConfig, zapLogger, pebbleClockStore, fileKeyManager, blossomSub, wesolowskiFrameProver, masterTimeReel, inMemoryPeerInfoManager, selfTestReport)
	node, err := newNode(zapLogger, pebbleClockStore, fileKeyManager, blossomSub, masterClockConsensusEngine)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func NewNode(configConfig *config.Config, selfTestReport *protobufs.SelfTestReport) (*Node, error) {
	zapLogger := logger()
	dbConfig := configConfig.DB
	pebbleDB := store.NewPebbleDB(dbConfig)
	pebbleClockStore := store.NewPebbleClockStore(pebbleDB, zapLogger)
	keyConfig := configConfig.Key
	fileKeyManager := keys.NewFileKeyManager(keyConfig, zapLogger)
	p2PConfig := configConfig.P2P
	blossomSub := p2p.NewBlossomSub(p2PConfig, zapLogger)
	engineConfig := configConfig.Engine
	wesolowskiFrameProver := crypto.NewWesolowskiFrameProver(zapLogger)
	masterTimeReel := time.NewMasterTimeReel(zapLogger, pebbleClockStore, engineConfig, wesolowskiFrameProver)
	inMemoryPeerInfoManager := p2p.NewInMemoryPeerInfoManager(zapLogger)
	masterClockConsensusEngine := master.NewMasterClockConsensusEngine(engineConfig, zapLogger, pebbleClockStore, fileKeyManager, blossomSub, wesolowskiFrameProver, masterTimeReel, inMemoryPeerInfoManager, selfTestReport)
	node, err := newNode(zapLogger, pebbleClockStore, fileKeyManager, blossomSub, masterClockConsensusEngine)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func NewDBConsole(configConfig *config.Config) (*DBConsole, error) {
	dbConsole, err := newDBConsole(configConfig)
	if err != nil {
		return nil, err
	}
	return dbConsole, nil
}

func NewClockStore(configConfig *config.Config) (store.ClockStore, error) {
	dbConfig := configConfig.DB
	pebbleDB := store.NewPebbleDB(dbConfig)
	zapLogger := logger()
	pebbleClockStore := store.NewPebbleClockStore(pebbleDB, zapLogger)
	return pebbleClockStore, nil
}

// wire.go:

func logger() *zap.Logger {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return log
}

func debugLogger() *zap.Logger {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return log
}

var loggerSet = wire.NewSet(
	logger,
)

var debugLoggerSet = wire.NewSet(
	debugLogger,
)

var keyManagerSet = wire.NewSet(wire.FieldsOf(new(*config.Config), "Key"), keys.NewFileKeyManager, wire.Bind(new(keys.KeyManager), new(*keys.FileKeyManager)))

var storeSet = wire.NewSet(wire.FieldsOf(new(*config.Config), "DB"), store.NewPebbleDB, wire.Bind(new(store.KVDB), new(*store.PebbleDB)), store.NewPebbleClockStore, store.NewPebbleKeyStore, store.NewPebbleDataProofStore, wire.Bind(new(store.ClockStore), new(*store.PebbleClockStore)), wire.Bind(new(store.KeyStore), new(*store.PebbleKeyStore)), wire.Bind(new(store.DataProofStore), new(*store.PebbleDataProofStore)))

var pubSubSet = wire.NewSet(wire.FieldsOf(new(*config.Config), "P2P"), p2p.NewInMemoryPeerInfoManager, p2p.NewBlossomSub, wire.Bind(new(p2p.PubSub), new(*p2p.BlossomSub)), wire.Bind(new(p2p.PeerInfoManager), new(*p2p.InMemoryPeerInfoManager)))

var engineSet = wire.NewSet(wire.FieldsOf(new(*config.Config), "Engine"), crypto.NewWesolowskiFrameProver, wire.Bind(new(crypto.FrameProver), new(*crypto.WesolowskiFrameProver)), crypto.NewKZGInclusionProver, wire.Bind(new(crypto.InclusionProver), new(*crypto.KZGInclusionProver)), time.NewMasterTimeReel, ceremony.NewCeremonyExecutionEngine)

var consensusSet = wire.NewSet(master.NewMasterClockConsensusEngine, wire.Bind(
	new(consensus.ConsensusEngine),
	new(*master.MasterClockConsensusEngine),
),
)
