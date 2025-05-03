package storages

import (
	"log"
)

// SchemaUpgrader is a callable to perform schema upgrade.
type SchemaUpgrader func(existedRev int32) (schemaChanged bool, err error)

// SchemaUpgradeStatus tracks the status of schema upgrade operation.
type SchemaUpgradeStatus struct {
	Changed   bool
	LastError error
}

// RunUpgrade calls given upgrader with given revision information and update
// status with returned result.
// If there are previous error, the upgrade callable will not be invoke.
func (st *SchemaUpgradeStatus) RunUpgrade(schemaName string, upgrader SchemaUpgrader, existedRev int32) (schemaChanged bool, err error) {
	if nil != st.LastError {
		return false, st.LastError
	}
	if schemaChanged, err = upgrader(existedRev); nil != err {
		st.LastError = err
		log.Printf("ERR: failed on upgrading schma [%s] from [%d]: %v", schemaName, existedRev, err)
	}
	if schemaChanged {
		st.Changed = true
	}
	return
}

// PushUpgradeResult updates state with given upgrade result.
//
// Caller should make sure current state is error free before attempt to run
// upgrade without RunUpgrade.
func (st *SchemaUpgradeStatus) PushUpgradeResult(schemaName string, schemaChanged bool, err error) {
	if nil != st.LastError {
		return
	}
	if nil != err {
		st.LastError = err
		log.Printf("ERR: failed on upgrading schma [%s]: %v", schemaName, err)
	}
	if schemaChanged {
		st.Changed = true
	}
	return
}
