package blockstore
import (
    "code.google.com/p/leveldb-go/leveldb"
    "log"
    "orwell/lib/logging"
    "github.com/mitchellh/go-homedir"
    "code.google.com/p/leveldb-go/leveldb/db"
    "orwell/lib/butils"
    "orwell/lib/protocol/orchain"
    "bytes"
    "io"
)

/*
    Database structure:

    "head"                      : <last block number> <last block ID>
    "h" <ID>                    : <Header>
    "t" <ID>                    : <Transaction>
    "b" <BillNumber>            : <Bill>
    "l" <ID>                    : [<ID>]
    "n" <uint64>                : <ID>

*/
var (
    key_HEAD = []byte("head")
)
const (
    prefix_H = 'h'
    prefix_T = 't'
    prefix_B = 'b'
    prefix_N = 'n'
    prefix_L = 'l'
)

type LevelDB struct {
    db *leveldb.DB
    log *log.Logger
    batch leveldb.Batch
}

func Open(path string) (s *LevelDB, err error) {
    s = &LevelDB{}
    s.log = logging.GetLogger("")
    // Make sure the relative paths are properly handled
    if path, err = homedir.Expand(path); err != nil { return }
    // Attempt opening a database
    if s.db, err = leveldb.Open(path, nil); err != nil { return }
    // Initialize database if necessary
    if _, err := s.get(key_HEAD); err == db.ErrNotFound {
        s.write(key_HEAD, &head_data{})
        ensure(s.flush())
    }
    return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *LevelDB) PutBlock(b *orchain.Block) error {
    s.batch = leveldb.Batch{}
    hd := &head_data{}
    hd.Num = s.Length() + 1
    hd.ID = b.Header.ID()
    ensure(s.write(key_HEAD, hd))
    if err := s.write(uint256Key(prefix_H, hd.ID), &b.Header); err != nil { return err }
    s.write(uint64Key(prefix_N, s.Length()), &hd.ID)

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
    s.set(uint256Key(prefix_L, hd.ID), buf.Bytes())

    return s.flush()
}

func (s *LevelDB) PopBlock() error {
    s.batch = leveldb.Batch{}
    h := s.GetHeaderByID(s.Head())
    hid := h.ID()
    hd := &head_data{}
    hd.Num = s.Length() - 1
    hd.ID = h.Previous
    ensure(s.write(key_HEAD, hd))
    s.del(uint256Key(prefix_H, hid))
    s.del(uint64Key(prefix_N, hd.Num))

    for _, tid := range s.GetTransactions(hid) {
        txn := s.GetTransaction(tid)
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

type head_data struct {
    Num uint64
    ID butils.Uint256
}

func (h *head_data) Read(r io.Reader) (err error) {
    if h.Num, err = butils.ReadUint64(r); err != nil { return }
    if err = h.ID.Read(r); err != nil { return }
    return
}

func (h *head_data) Write(w io.Writer) (err error) {
    if err = butils.WriteUint64(w, h.Num); err != nil { return }
    if err = h.ID.Write(w); err != nil { return }
    return
}

func (s *LevelDB) getHead() (h *head_data) {
    h = &head_data{}
    ensure(s.read(key_HEAD, h))
    return
}

func (s *LevelDB) Length() uint64 {
    h := s.getHead()
    return h.Num
}

func (s *LevelDB) Head() butils.Uint256 {
    h := s.getHead()
    return h.ID
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *LevelDB) GetHeaderByID(id butils.Uint256) (h *orchain.Header) {
    h = &orchain.Header{}
    err := s.read(uint256Key(prefix_H, id), h)
    if err == nil { return }
    if err == db.ErrNotFound { return nil }
    panic(err)
}

func (s *LevelDB) GetHeaderByNum(num uint64) (h *orchain.Header) {
    id := butils.Uint256{}
    err := s.read(uint64Key(prefix_N, num), &id)
    if err == nil { return s.GetHeaderByID(id) }
    if err == db.ErrNotFound { return nil }
    panic(err)
}

func (s *LevelDB) GetTransaction(id butils.Uint256) (t *orchain.Transaction) {
    t = &orchain.Transaction{}
    err := s.read(uint256Key(prefix_T, id), t)
    if err == nil { return }
    if err == db.ErrNotFound { return nil }
    panic(err)
}

func (s *LevelDB) GetBill(number orchain.BillNumber) (b *orchain.Bill) {
    b = &orchain.Bill{}
    err := s.read(billKey(number), b)
    if err == nil { return }
    if err == db.ErrNotFound { return nil }
    panic(err)
}

func (s *LevelDB) GetTransactions(id butils.Uint256) []butils.Uint256 {
    data, err := s.get(uint256Key(prefix_L, id))
    if err == db.ErrNotFound { return nil }
    ensure(err)
    r := bytes.NewBuffer(data)
    num := len(data) / butils.UINT256_LENGTH_BYTES
    res := make([]butils.Uint256, int(num))
    for i := 0; i < num; i += 1 {
        ensure(res[i].Read(r))
    }
    return res
}


///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Reads the given key, and tries to unpack it into the given readable structure
func (s *LevelDB) read(key []byte, r butils.Readable) error {
    buf, err := s.get(key)
    if err != nil { return err }
    return butils.ReadAllInto(r, buf)
}

// Writes the given writable structure to a given key
func (s *LevelDB) write(key []byte, w butils.Writable) error {
    buf, err := butils.WriteToBytes(w)
    if err != nil { return err }
    s.set(key, buf)
    return nil
}

// Returns the raw value associated with the key, if any
func (s *LevelDB) get(key []byte) ([]byte, error) {
    return s.db.Get(key, nil)
}

// Writes the specified key-value pair to write buffer
func (s *LevelDB) set(key, value []byte) {
    s.batch.Set(key, value)
}

// Writes the remove-operation to the write buffer
func (s *LevelDB) del(key []byte) {
    s.batch.Delete(key)
}

// Applies all operations in the buffer
func (s *LevelDB) flush() error {
    return s.db.Apply(s.batch, &db.WriteOptions{true})
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func uint256Key(prefix byte, id butils.Uint256) []byte {
    key := &bytes.Buffer{}
    ensure(butils.WriteByte(key, prefix))
    ensure(butils.WriteFull(key, id[:]))
    return key.Bytes()
}

func uint64Key(prefix byte, num uint64) []byte {
    key := &bytes.Buffer{}
    ensure(butils.WriteByte(key, prefix))
    ensure(butils.WriteUint64(key, num))
    return key.Bytes()
}

func billKey(number orchain.BillNumber) []byte {
    buf := &bytes.Buffer{}
    ensure(butils.WriteByte(buf, prefix_B))
    ensure(number.Write(buf))
    return buf.Bytes()
}