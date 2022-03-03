package database

import "github.com/uuosio/chain"

type StateData struct {
	op    int
	data  interface{}
	payer uint64
}

type StateDataDBI64 struct {
	db     *DBI64
	states map[uint64]*StateData
}

type StateDataIdxDB struct {
	db     SecondaryDB
	states map[uint64]*StateData
}

type StateManager struct {
	DBs    map[chain.Uint128]*StateDataDBI64
	IdxDBs map[chain.Uint128]*StateDataIdxDB
}

var StateManagerInstance *StateManager

func GetStateManager() *StateManager {
	if StateManagerInstance == nil {
		StateManagerInstance = &StateManager{
			DBs:    make(map[chain.Uint128]*StateDataDBI64),
			IdxDBs: make(map[chain.Uint128]*StateDataIdxDB),
		}
	}
	return StateManagerInstance
}

func (s *StateManager) OnStore(db *DBI64, primary uint64) {
	if !chain.IsRevertEnabled() {
		return
	}

	s.AddState(db, OperationStore, primary, nil, 0)
}

func (s *StateManager) OnUpdate(db *DBI64, it Iterator, payer uint64) {
	if !chain.IsRevertEnabled() {
		return
	}

	rawValue := db.GetByIterator(it)
	value := db.unpacker(rawValue)
	primary := value.GetPrimary()
	s.AddState(db, OperationUpdate, primary, rawValue, payer)
}

func (s *StateManager) OnRemove(db *DBI64, it Iterator) {
	if !chain.IsRevertEnabled() {
		return
	}

	rawValue := db.GetByIterator(it)
	value := db.unpacker(rawValue)
	primary := value.GetPrimary()
	s.AddState(db, OperationRemove, primary, value, chain.CurrentReceiver().N)
}

func (s *StateManager) AddState(db *DBI64, op int, primary uint64, value interface{}, payer uint64) {
	if !chain.IsRevertEnabled() {
		return
	}

	code, scope, table := db.GetTable()
	if code != chain.CurrentReceiver().N {
		return
	}

	key := *chain.NewUint128(scope, table)
	if _, ok := s.DBs[key]; !ok {
		s.DBs[key] = &StateDataDBI64{db, make(map[uint64]*StateData)}
	}

	state := &StateData{op, value, payer}
	s.DBs[key].states[primary] = state
}

func (s *StateManager) OnIdxDBStore(db SecondaryDB, primary uint64) {
	s.AddIdxState(db, OperationStore, primary, 0)
}

func (s *StateManager) OnIdxDBUpdate(db SecondaryDB, it SecondaryIterator, payer uint64) {
	s.AddIdxState(db, OperationUpdate, it.Primary, payer)
}

func (s *StateManager) OnIdxDBRemove(db SecondaryDB, it SecondaryIterator) {
	s.AddIdxState(db, OperationRemove, it.Primary, chain.CurrentReceiver().N)
}

func (s *StateManager) AddIdxState(db SecondaryDB, op int, primary uint64, payer uint64) {
	if !chain.IsRevertEnabled() {
		return
	}

	code, scope, table := db.GetTable()
	if code != chain.CurrentReceiver().N {
		return
	}

	key := *chain.NewUint128(scope, table)
	if _, ok := s.IdxDBs[key]; !ok {
		s.IdxDBs[key] = &StateDataIdxDB{db, make(map[uint64]*StateData)}
	}

	var value interface{}
	value = nil

	if op == OperationStore {
		payer = chain.CurrentReceiver().N
	} else {
		it2, secondaryValue := db.FindByPrimary(primary)
		chain.Check(it2.IsOk(), "bad iterator")
		value = secondaryValue
	}

	state := &StateData{op, value, payer}
	s.IdxDBs[key].states[primary] = state
}

func (s *StateManager) Revert() {
	if !chain.IsRevertEnabled() {
		return
	}

	for _, v := range s.DBs {
		for primary, state := range v.states {
			if state.op == OperationStore {
				it := v.db.Find(primary)
				if it.IsOk() {
					v.db.Remove(it)
				}
			} else if state.op == OperationUpdate {
				it := v.db.Find(primary)
				if it.IsOk() {
					v.db.Update(it, state.data.([]byte), chain.Name{state.payer})
				}
			} else if state.op == OperationRemove {
				if it := v.db.Find(primary); it.IsOk() {
					v.db.Update(it, state.data.([]byte), chain.Name{state.payer})
				} else {
					v.db.Store(primary, state.data.([]byte), chain.Name{state.payer})
				}
			}
		}
	}

	for _, v := range s.IdxDBs {
		for primary, state := range v.states {
			switch state.op {
			case OperationStore:
				it, _ := v.db.FindByPrimary(primary)
				if it.IsOk() {
					v.db.Remove(it)
				}
			case OperationUpdate, OperationRemove:
				it, _ := v.db.FindByPrimary(primary)
				if it.IsOk() {
					v.db.UpdateEx(it, state.data, state.payer)
				} else {
					v.db.StoreEx(primary, state.data, state.payer)
				}
				break
			}
		}
	}
}
