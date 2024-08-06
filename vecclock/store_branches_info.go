package vecclock

//func (vi *Index) setRlp(table kvdb.Store, key []byte, val interface{}) {
//	buf, err := rlp.EncodeToBytes(val)
//	if err != nil {
//		vi.crit(err)
//	}
//
//	if err := table.Put(key, buf); err != nil {
//		vi.crit(err)
//	}
//}
//
//func (vi *Index) getRlp(table kvdb.Store, key []byte, to interface{}) interface{} {
//	buf, err := table.Get(key)
//	if err != nil {
//		vi.crit(err)
//	}
//	if buf == nil {
//		return nil
//	}
//
//	err = rlp.DecodeBytes(buf, to)
//	if err != nil {
//		vi.crit(err)
//	}
//	return to
//}
//
//func (vi *Index) getBytes(table kvdb.Store, id hash.Event) []byte {
//	key := id.Bytes()
//	b, err := table.Get(key)
//	if err != nil {
//		vi.crit(err)
//	}
//	return b
//}
//
//func (vi *Index) setBytes(table kvdb.Store, id hash.Event, b []byte) {
//	key := id.Bytes()
//	err := table.Put(key, b)
//	if err != nil {
//		vi.crit(err)
//	}
//}
//
//func (vi *Index) setBranchesInfo(info *BranchesInfo) {
//	key := []byte("c")
//
//	vi.setRlp(vi.table.BranchesInfo, key, info)
//}
//
//func (vi *Index) getBranchesInfo() *BranchesInfo {
//	key := []byte("c")
//
//	w, exists := vi.getRlp(vi.table.BranchesInfo, key, &BranchesInfo{}).(*BranchesInfo)
//	if !exists {
//		return nil
//	}
//
//	return w
//}
//
//// SetEventBranchID stores the event's global branch ID
//func (vi *Index) SetEventBranchID(id hash.Event, branchID idx.Validator) {
//	vi.setBytes(vi.table.EventBranch, id, branchID.Bytes())
//}
//
//// GetEventBranchID reads the event's global branch ID
//func (vi *Index) GetEventBranchID(id hash.Event) idx.Validator {
//	b := vi.getBytes(vi.table.EventBranch, id)
//	if b == nil {
//		vi.crit(errors.New("failed to read event's branch ID (inconsistent DB)"))
//		return 0
//	}
//	branchID := idx.BytesToValidator(b)
//	return branchID
//}
