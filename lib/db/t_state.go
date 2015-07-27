package db

var BUCKET_INFO = []byte("info")
var KEY_STATE = []byte("state")

func (t *Tx) GetState() (s *State) {
    s = &State{}
    if t.Read(BUCKET_INFO, KEY_STATE, s) { return }
    return nil
}

func (t *Tx) PutState(s *State) {
    t.Write(BUCKET_INFO, KEY_STATE, s)
}