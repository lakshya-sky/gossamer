// Copyright 2019 ChainSafe Systems (ON) Corp.
// This file is part of gossamer.
//
// The gossamer library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The gossamer library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the gossamer library. If not, see <http://www.gnu.org/licenses/>.

package modules

import (
	"encoding/hex"
	"net/http"

	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/lib/scale"
)

// StateCallRequest holds json fields
type StateCallRequest struct {
	Method string       `json:"method"`
	Data   []byte       `json:"data"`
	Block  *common.Hash `json:"block"`
}

// StateChildStorageRequest holds json fields
type StateChildStorageRequest struct {
	ChildStorageKey []byte       `json:"childStorageKey"`
	Key             []byte       `json:"key"`
	Block           *common.Hash `json:"block"`
}

// StateStorageKeyRequest holds json fields
type StateStorageKeyRequest struct {
	Key   []byte       `json:"key"`
	Block *common.Hash `json:"block"`
}

// StateRuntimeMetadataQuery is a hash value
type StateRuntimeMetadataQuery struct {
	Bhash *common.Hash
}

// StateRuntimeVersionRequest is hash value
type StateRuntimeVersionRequest struct {
	Bhash *common.Hash
}

// StatePairRequest holds json field
type StatePairRequest struct {
	Prefix string
	Bhash  *common.Hash
}

// StateStorageSizeRequest holds json field
type StateStorageSizeRequest struct {
	Key   string
	Bhash *common.Hash
}

// StateStorageHashRequest holds json field
type StateStorageHashRequest struct {
	Key   string
	Bhash *common.Hash
}

// StateStorageRequest holds json field
type StateStorageRequest struct {
	Key   string
	Bhash *common.Hash
}

// StateStorageQueryRangeRequest holds json fields
type StateStorageQueryRangeRequest struct {
	Keys       []*common.Hash `json:"keys"`
	StartBlock *common.Hash   `json:"startBlock"`
	Block      *common.Hash   `json:"block"`
}

// StateStorageKeysQuery field to store storage keys
type StateStorageKeysQuery [][]byte

// StateCallResponse holds json fields
type StateCallResponse struct {
	StateCallResponse []byte `json:"stateCallResponse"`
}

// StateKeysResponse field to store the state keys
type StateKeysResponse [][]byte

// StateStorageDataResponse field to store data response
type StateStorageDataResponse string

// StateStorageHashResponse is a hash value
type StateStorageHashResponse struct {
	StorageHash string
}

// StateChildStorageResponse is a hash value
type StateChildStorageResponse struct {
	StorageHash string
}

// StateChildStorageSizeResponse is a unint value
type StateChildStorageSizeResponse struct {
	Size uint64
}

// StateStorageSizeResponse the default size for response
type StateStorageSizeResponse struct {
	StorageEntrySize uint64
}

// StateStorageResponse storage hash value
type StateStorageResponse struct {
	StorageValue string
}

// StatePairResponse is a key values
type StatePairResponse struct {
	keys []interface{}
}

// StateStorageKeysResponse field for storage keys
type StateStorageKeysResponse [][]byte

// StateMetadataResponse holds the metadata
//TODO: Determine actual type
type StateMetadataResponse struct {
	Metadata string
}

// StorageChangeSetResponse is the struct that holds the block and changes
type StorageChangeSetResponse struct {
	Block   *common.Hash
	Changes []KeyValueOption
}

// KeyValueOption struct holds json fields
type KeyValueOption struct {
	StorageKey  []byte `json:"storageKey"`
	StorageData []byte `json:"storageData"`
}

// StorageKey is the key for the storage
type StorageKey []byte

// StateRuntimeVersionResponse is the runtime version response
type StateRuntimeVersionResponse struct {
	SpecName         string        `json:"specName"`
	ImplName         string        `json:"implName"`
	AuthoringVersion int32         `json:"authoringVersion"`
	SpecVersion      int32         `json:"specVersion"`
	ImplVersion      int32         `json:"implVersion"`
	Apis             []interface{} `json:"apis"`
}

// StateModule is an RPC module providing access to storage API points.
type StateModule struct {
	networkAPI NetworkAPI
	storageAPI StorageAPI
	coreAPI    CoreAPI
}

// NewStateModule creates a new State module.
func NewStateModule(net NetworkAPI, storage StorageAPI, core CoreAPI) *StateModule {
	return &StateModule{
		networkAPI: net,
		storageAPI: storage,
		coreAPI:    core,
	}
}

// GetPairs returns the keys with prefix, leave empty to get all the keys.
func (sm *StateModule) GetPairs(r *http.Request, req *StatePairRequest, res *StatePairResponse) error {
	// TODO implement change storage trie so that block hash parameter works (See issue #834)
	var (
		stateRootHash *common.Hash
		err           error
	)

	reqBytes, _ := common.HexToBytes(req.Prefix)
	if req.Bhash != nil {
		stateRootHash, err = sm.storageAPI.GetStateRootFromBlock(req.Bhash)
		if err != nil {
			return err
		}
	}

	if len(reqBytes) < 1 {
		pairs, err := sm.storageAPI.Entries(stateRootHash)
		if err != nil {
			return err
		}
		for k, v := range pairs {
			res.keys = append(res.keys, []string{"0x" + hex.EncodeToString([]byte(k)), "0x" + hex.EncodeToString(v)})
		}
	} else {
		// TODO this should return all keys with same prefix, currently only returning
		//  matches.  Implement when #837 is done.
		resI, err := sm.storageAPI.GetStorage(stateRootHash, reqBytes)
		if err != nil {
			return err
		}
		if resI != nil {
			res.keys = append(res.keys, []string{"0x" + hex.EncodeToString(reqBytes), "0x" + hex.EncodeToString(resI)})
		} else {
			res.keys = []interface{}{}
		}
	}

	return nil
}

// Call isn't implemented properly yet.
func (sm *StateModule) Call(r *http.Request, req *StateCallRequest, res *StateCallResponse) {
	_ = sm.networkAPI
	_ = sm.storageAPI
}

// GetChildKeys isn't implemented properly yet.
func (sm *StateModule) GetChildKeys(r *http.Request, req *StateChildStorageRequest, res *StateKeysResponse) {
	// TODO implement change storage trie so that block hash parameter works (See issue #834)
}

// GetChildStorage isn't implemented properly yet.
func (sm *StateModule) GetChildStorage(r *http.Request, req *StateChildStorageRequest, res *StateStorageDataResponse) {
	// TODO implement change storage trie so that block hash parameter works (See issue #834)
}

// GetChildStorageHash isn't implemented properly yet.
func (sm *StateModule) GetChildStorageHash(r *http.Request, req *StateChildStorageRequest, res *StateChildStorageResponse) {
	// TODO implement change storage trie so that block hash parameter works (See issue #834)
}

// GetChildStorageSize isn't implemented properly yet.
func (sm *StateModule) GetChildStorageSize(r *http.Request, req *StateChildStorageRequest, res *StateChildStorageSizeResponse) {
	// TODO implement change storage trie so that block hash parameter works (See issue #834)
}

// GetKeys isn't implemented properly yet.
func (sm *StateModule) GetKeys(r *http.Request, req *StateStorageKeyRequest, res *StateStorageKeysResponse) {
	// TODO implement change storage trie so that block hash parameter works (See issue #834)
}

// GetMetadata calls runtime Metadata_metadata function
func (sm *StateModule) GetMetadata(r *http.Request, req *StateRuntimeMetadataQuery, res *StateMetadataResponse) error {
	// TODO implement change storage trie so that block hash parameter works (See issue #834)
	metadata, err := sm.coreAPI.GetMetadata(req.Bhash)
	if err != nil {
		return err
	}

	decoded, err := scale.Decode(metadata, []byte{})
	res.Metadata = common.BytesToHex(decoded.([]byte))
	return err
}

// GetRuntimeVersion Get the runtime version at a given block.
//  If no block hash is provided, the latest version gets returned.
// TODO currently only returns latest version, add functionality to lookup runtime by block hash (see issue #834)
func (sm *StateModule) GetRuntimeVersion(r *http.Request, req *StateRuntimeVersionRequest, res *StateRuntimeVersionResponse) error {
	rtVersion, err := sm.coreAPI.GetRuntimeVersion(req.Bhash)
	if err != nil {
		return err
	}

	res.SpecName = string(rtVersion.RuntimeVersion.Spec_name)
	res.ImplName = string(rtVersion.RuntimeVersion.Impl_name)
	res.AuthoringVersion = rtVersion.RuntimeVersion.Authoring_version
	res.SpecVersion = rtVersion.RuntimeVersion.Spec_version
	res.ImplVersion = rtVersion.RuntimeVersion.Impl_version
	res.Apis = convertAPIs(rtVersion.API)

	return nil
}

// GetStorage Returns a storage entry at a specific block's state. If not block hash is provided, the latest value is returned.
func (sm *StateModule) GetStorage(r *http.Request, req *StateStorageRequest, res *StateStorageResponse) error {
	var (
		item []byte
		err  error
	)

	reqBytes, _ := common.HexToBytes(req.Key) // no need to catch error here
	if req.Bhash != nil {
		item, err = sm.storageAPI.GetStorageByBlockHash(*req.Bhash, reqBytes)
		if err != nil {
			return err
		}
	} else {
		item, err = sm.storageAPI.GetStorage(nil, reqBytes)
		if err != nil {
			return err
		}
	}

	if len(item) > 0 {
		res.StorageValue = common.BytesToHex(item)
	} else {
		*res = StateStorageResponse{}
	}

	return nil
}

// GetStorageHash returns the hash of a storage entry at a block's state.
//  If no block hash is provided, the latest value is returned.
//  TODO implement change storage trie so that block hash parameter works (See issue #834)
func (sm *StateModule) GetStorageHash(r *http.Request, req *StateStorageHashRequest, res *StateStorageHashResponse) error {
	var (
		item []byte
		err  error
	)

	reqBytes, _ := common.HexToBytes(req.Key)

	if req.Bhash != nil {
		item, err = sm.storageAPI.GetStorageByBlockHash(*req.Bhash, reqBytes)
		if err != nil {
			return err
		}
	} else {
		item, err = sm.storageAPI.GetStorage(nil, reqBytes)
		if err != nil {
			return err
		}
	}

	if len(item) > 0 {
		res.StorageHash = common.BytesToHash(item).String()
	} else {
		*res = StateStorageHashResponse{}
	}

	return nil
}

// GetStorageSize returns the size of a storage entry at a block's state.
//  If no block hash is provided, the latest value is used.
// TODO implement change storage trie so that block hash parameter works (See issue #834)
func (sm *StateModule) GetStorageSize(r *http.Request, req *StateStorageSizeRequest, res *StateStorageSizeResponse) error {
	var (
		item []byte
		err  error
	)

	reqBytes, _ := common.HexToBytes(req.Key)

	if req.Bhash != nil {
		item, err = sm.storageAPI.GetStorageByBlockHash(*req.Bhash, reqBytes)
		if err != nil {
			return err
		}
	} else {
		item, err = sm.storageAPI.GetStorage(nil, reqBytes)
		if err != nil {
			return err
		}
	}

	if len(item) > 0 {
		*res = StateStorageSizeResponse{uint64(len(item))}
	} else {
		*res = StateStorageSizeResponse{}
	}

	return nil
}

// QueryStorage isn't implemented properly yet.
func (sm *StateModule) QueryStorage(r *http.Request, req *StateStorageQueryRangeRequest, res *StorageChangeSetResponse) error {
	// TODO implement change storage trie so that block hash parameter works (See issue #834)
	return nil
}

// SubscribeRuntimeVersion isn't implemented properly yet.
// TODO make this actually a subscription that pushes data
func (sm *StateModule) SubscribeRuntimeVersion(r *http.Request, req *StateStorageQueryRangeRequest, res *StateRuntimeVersionResponse) error {
	// TODO implement change storage trie so that block hash parameter works (See issue #834)
	return sm.GetRuntimeVersion(r, nil, res)
}

// SubscribeStorage Storage subscription. If storage keys are specified, it creates a message for each block which
//  changes the specified storage keys. If none are specified, then it creates a message for every block.
//  This endpoint communicates over the Websocket protocol, but this func should remain here so it's added to rpc_methods list
func (sm *StateModule) SubscribeStorage(r *http.Request, req *StateStorageQueryRangeRequest, res *StorageChangeSetResponse) error {
	return nil
}

func convertAPIs(in []*runtime.API_Item) []interface{} {
	ret := make([]interface{}, 0)
	for _, item := range in {
		encStr := hex.EncodeToString(item.Name)
		ret = append(ret, []interface{}{"0x" + encStr, item.Ver})
	}
	return ret
}
