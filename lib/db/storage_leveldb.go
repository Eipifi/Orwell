package db
import (
    "log"
    "bytes"
    "code.google.com/p/leveldb-go/leveldb"
    "code.google.com/p/leveldb-go/leveldb/db"
    "orwell/lib/foo"
    "orwell/lib/utils"
    "orwell/lib/butils"
    "orwell/lib/logging"
    "orwell/lib/protocol/orchain"
)

/*
    Database structure:

    "state"                     : <State>
    "h" <ID>                    : <Header>
    "t" <ID>                    : <Transaction>
    "b" <BillNumber>            : <Bill>
    "l" <ID>                    : [<ID>]
    "n" <uint64>                : <ID>
    "r" <ID>                    : <uint64>

*/
var (
    key_State = []byte("state")
)
const (
    prefix_H = 'h'
    prefix_T = 't'
    prefix_B = 'b'
    prefix_N = 'n'
    prefix_L = 'l'
    prefix_R = 'r'
)

type LDBStorage struct {
    db *leveldb.DB
    log *log.Logger
    batch leveldb.Batch
}

func OpenLDBStorage(path string) (s *LDBStorage) {
    var err error
    s = &LDBStorage{}
    s.log = logging.GetLogger("")
    // Attempt opening a database
    if s.db, err = leveldb.Open(path, nil); err != nil {
        log.Fatalln(err)
    }
    // Initialize database if necessary
    if _, err := s.get(key_State); err == db.ErrNotFound {
        utils.Ensure(s.write(key_State, &State{}))
        utils.Ensure(s.flush())
    }
    return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *LDBStorage) PutBlock(b *orchain.Block) error {
    s.batch = leveldb.Batch{}
    state := s.State()
    state.Length += 1
    state.Head = b.Header.ID()
    state.Work.Add(b.Header.Difficulty)
    utils.Ensure(s.write(key_State, state))
    if err := s.write(uint256Key(prefix_H, state.Head), &b.Header); err != nil { return err }
    s.write(uint64Key(prefix_N, state.Length - 1), &state.Head)
    tmp := foo.U64(state.Length - 1)
    s.write(uint256Key(prefix_R, state.Head), &tmp)

    // Insert the transactions
    buf := &bytes.Buffer{}
    for _, txn := range b.Transactions {
        tid, err := txn.ID()
        if err != nil { return err }
        tid.Write(buf)
        if err := s.write(uint256Key(prefix_T, tid), &txn); err != nil { return err }
        for _, inp := range txn.Inputs {
            s.del(billKey(inp))
        }
        for i, out := range txn.Outputs {
            if err := s.write(billKey(orchain.BillNumber{tid, uint64(i)}), &out); err != nil { return err }
        }
    }
    // Assign the transactions to a header
    s.set(uint256Key(prefix_L, state.Head), buf.Bytes())

    return s.flush()
}

func (s *LDBStorage) PopBlock() error {
    s.batch = leveldb.Batch{}
    state := s.State()
    h := s.GetHeaderByID(state.Head)
    hid := h.ID()
    state.Length -= 1
    state.Head = h.Previous
    state.Work.Sub(h.Difficulty)
    utils.Ensure(s.write(key_State, state))
    s.del(uint256Key(prefix_H, hid))
    s.del(uint64Key(prefix_N, state.Length))
    s.del(uint256Key(prefix_R, hid))

    for _, tid := range s.GetTransactions(hid) {
        txn := s.GetTransaction(tid)
        utils.Assert(txn != nil)
        s.del(uint256Key(prefix_T, tid))
        for _, inp := range txn.Inputs {
            referred_txn := s.GetTransaction(inp.Txn)
            s.write(billKey(inp), &referred_txn.Outputs[inp.Index])
        }
        for i, _ := range txn.Outputs {
            s.del(billKey(orchain.BillNumber{tid, uint64(i)}))
        }
    }
    s.del(uint256Key(prefix_L, hid))
    return s.flush()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *LDBStorage) State() (h *State) {
    h = &State{}
    utils.Ensure(s.read(key_State, h))
    return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *LDBStorage) GetHeaderByID(id foo.U256) (h *orchain.Header) {
    h = &orchain.Header{}
    err := s.read(uint256Key(prefix_H, id), h)
    if err == nil { return }
    if err == db.ErrNotFound { return nil }
    panic(err)
}

func (s *LDBStorage) GetHeaderByNum(num uint64) (h *orchain.Header) {
    id := s.GetIDByNum(num)
    if id == nil { return nil }
    return s.GetHeaderByID(*id)
}

func (s *LDBStorage) GetIDByNum(num uint64) (h *foo.U256) {
    id := foo.U256{}
    err := s.read(uint64Key(prefix_N, num), &id)
    if err == nil { return &id }
    if err == db.ErrNotFound { return nil }
    panic(err)
}

func (s *LDBStorage) GetNumByID(id foo.U256) *uint64 {
    var val foo.U64
    err := s.read(uint256Key(prefix_R, id), &val)
    if err == nil {
        tmp := uint64(val)
        return &tmp
    }
    if err == db.ErrNotFound { return nil }
    panic(err)
}

func (s *LDBStorage) GetTransaction(id foo.U256) (t *orchain.Transaction) {
    t = &orchain.Transaction{}
    err := s.read(uint256Key(prefix_T, id), t)
    if err == nil { return }
    if err == db.ErrNotFound { return nil }
    panic(err)
}

func (s *LDBStorage) GetBill(number orchain.BillNumber) (b *orchain.Bill) {
    b = &orchain.Bill{}
    err := s.read(billKey(number), b)
    if err == nil { return }
    if err == db.ErrNotFound { return nil }
    panic(err)
}

func (s *LDBStorage) GetTransactions(id foo.U256) []foo.U256 {
    data, err := s.get(uint256Key(prefix_L, id))
    if err == db.ErrNotFound { return nil }
    utils.Ensure(err)
    r := bytes.NewBuffer(data)
    num := len(data) / foo.U256_BYTES
    res := make([]foo.U256, int(num))
    for i := 0; i < num; i += 1 {
        utils.Ensure(res[i].Read(r))
    }
    return res
}


///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Reads the given key, and tries to unpack it into the given readable structure
func (s *LDBStorage) read(key []byte, r butils.Readable) error {
    buf, err := s.get(key)
    if err != nil { return err }
    return butils.ReadAllInto(r, buf)
}

// Writes the given writable structure to a given key
func (s *LDBStorage) write(key []byte, w butils.Writable) error {
    buf, err := butils.WriteToBytes(w)
    if err != nil { return err }
    s.set(key, buf)
    return nil
}

// Returns the raw value associated with the key, if any
func (s *LDBStorage) get(key []byte) ([]byte, error) {
    return s.db.Get(key, nil)
}

// Writes the specified key-value pair to write buffer
func (s *LDBStorage) set(key, value []byte) {
    s.batch.Set(key, value)
}

// Writes the remove-operation to the write buffer
func (s *LDBStorage) del(key []byte) {
    s.batch.Delete(key)
}

// Applies all operations in the buffer
func (s *LDBStorage) flush() error {
    return s.db.Apply(s.batch, &db.WriteOptions{true})
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func uint256Key(prefix byte, id foo.U256) []byte {
    key := &bytes.Buffer{}
    utils.Ensure(butils.WriteByte(key, prefix))
    utils.Ensure(butils.WriteFull(key, id[:]))
    return key.Bytes()
}

func uint64Key(prefix byte, num uint64) []byte {
    key := &bytes.Buffer{}
    utils.Ensure(butils.WriteByte(key, prefix))
    utils.Ensure(butils.WriteUint64(key, num))
    return key.Bytes()
}

func billKey(number orchain.BillNumber) []byte {
    buf := &bytes.Buffer{}
    utils.Ensure(butils.WriteByte(buf, prefix_B))
    utils.Ensure(number.Write(buf))
    return buf.Bytes()
}