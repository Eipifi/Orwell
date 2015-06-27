package blockstore
import (
    "code.google.com/p/leveldb-go/leveldb"
    "github.com/mitchellh/go-homedir"
    "orwell/lib/protocol/orchain"
    "bytes"
    "orwell/lib/butils"
    "log"
    "orwell/lib/logging"
    "code.google.com/p/leveldb-go/leveldb/db"
    "errors"
)

type LevelDB struct {
    db *leveldb.DB
    log *log.Logger
}

var key_HEAD = []byte("head")
const prefix_HEADER = 'h'           // 0x68
const prefix_TXN = 't'              // 0x74
const prefix_BILL = 'b'             // 0x62
const prefix_TXN_LIST = 'l'         // 0x6c
const prefix_HEADER_NUM = 'n'       // 0x6e

func Open(path string) (s *LevelDB, err error) {
    s = &LevelDB{}
    s.log = logging.GetLogger("")
    // Make sure the relative paths are properly handled
    if path, err = homedir.Expand(path); err != nil { return }
    // Attempt opening a database
    if s.db, err = leveldb.Open(path, nil); err != nil { return }
    // Initialize database if necessary
    if s.get(key_HEAD) == nil {
        s.StoreHead(butils.Uint256{}, 0)
    }
    return
}

func (s *LevelDB) StoreHead(head butils.Uint256, num uint64) {
    buf := &bytes.Buffer{}
    ensure(head.Write(buf))
    ensure(butils.WriteUint64(buf, num))
    ensure(s.set(key_HEAD, buf.Bytes()))
}

func (s *LevelDB) FetchHead() (head butils.Uint256, num uint64) {
    buf := s.get(key_HEAD)
    if buf == nil {
        panic(errors.New("No head entry"))
    }
    r := bytes.NewBuffer(buf)
    ensure(head.Read(r))
    num, err := butils.ReadUint64(r)
    ensure(err)
    return
}

func (s *LevelDB) StoreHeader(h *orchain.Header, num uint64) error {
    buf := &bytes.Buffer{}
    hid := h.ID()
    ensure(butils.WriteUint64(buf, num))
    if err := h.Write(buf); err != nil { return err }
    ensure(s.set(uint256Key(prefix_HEADER, hid), buf.Bytes()))
    ensure(s.set(uint64Key(prefix_HEADER_NUM, num), hid[:]))
    return nil
}

func (s *LevelDB) FetchHeader(hash butils.Uint256) (h *orchain.Header) {
    data := s.get(uint256Key(prefix_HEADER, hash))
    if data == nil { return nil }
    buf := bytes.NewBuffer(data)
    _, err := butils.ReadUint64(buf)
    ensure(err)
    h = &orchain.Header{}
    ensure(h.Read(buf))
    return
}

func (s *LevelDB) FetchHeaderByNum(num uint64) (h *orchain.Header) {
    data := s.get(uint64Key(prefix_HEADER_NUM, num))
    assert(data != nil)
    hid := butils.Uint256{}
    ensure(butils.ReadAllInto(&hid, data))
    return s.FetchHeader(hid)
}

func (s *LevelDB) RemoveHeader(hid butils.Uint256) {
    data := s.get(uint256Key(prefix_HEADER, hid))
    assert(data != nil)
    buf := bytes.NewBuffer(data)
    num, err := butils.ReadUint64(buf)
    ensure(err)
    ensure(s.del(uint256Key(prefix_HEADER, hid)))
    ensure(s.del(uint64Key(prefix_HEADER_NUM, num)))
}

func (s *LevelDB) StoreBlockTransactionIDs(bid butils.Uint256, tids []butils.Uint256) {
    buf := &bytes.Buffer{}
    for _, tid := range tids {
        ensure(tid.Write(buf))
    }
    s.set(uint256Key(prefix_TXN_LIST, bid), buf.Bytes())
}

func (s *LevelDB) FetchBlockTransactionIDs(bid butils.Uint256) []butils.Uint256 {
    buf := s.get(uint256Key(prefix_TXN_LIST, bid))
    if buf == nil { return nil }
    num := len(buf) / butils.UINT256_LENGTH_BYTES
    tids := make([]butils.Uint256, num)
    r := bytes.NewBuffer(buf)
    for i := 0; i < num; i += 1 {
        ensure(tids[i].Read(r))
    }
    return tids
}

func (s *LevelDB) RemoveBlockTransactionIDs(bid butils.Uint256) {
    ensure(s.del(uint256Key(prefix_TXN_LIST, bid)))
}

func (s *LevelDB) StoreTransaction(t *orchain.Transaction) error {
    buf, err := butils.WriteToBytes(t)
    if err != nil { return err }
    tid, err := t.ID()
    if err != nil { return err }
    ensure(s.set(uint256Key(prefix_TXN, tid), buf))
    return nil
}

func (s *LevelDB) FetchTransaction(tid butils.Uint256) *orchain.Transaction {
    buf := s.get(uint256Key(prefix_TXN, tid))
    if buf == nil { return nil }
    t := &orchain.Transaction{}
    ensure(butils.ReadAllInto(t, buf))
    return t
}

func (s *LevelDB) RemoveTransaction(tid butils.Uint256) {
    ensure(s.del(uint256Key(prefix_TXN, tid)))
}

func (s *LevelDB) StoreUnspentBill(number orchain.BillNumber, bill orchain.Bill) {
    buf, err := butils.WriteToBytes(&bill)
    ensure(err)
    ensure(s.set(billKey(number), buf))
}

func (s *LevelDB) FetchUnspentBill(number orchain.BillNumber) *orchain.Bill {
    buf := s.get(billKey(number))
    if buf == nil { return nil }
    bill := &orchain.Bill{}
    ensure(butils.ReadAllInto(bill, buf))
    return bill
}

func (s *LevelDB) SpendBill(number orchain.BillNumber) {
    ensure(s.del(billKey(number)))
}

////////////////////////////

func (s *LevelDB) get(key []byte) []byte {
    buf, err := s.db.Get(key, nil)
    if err == nil {
        s.log.Printf("Get [%x] => [%x]", key, buf)
    } else {
        s.log.Printf("Get [%x] => [%x] [Cause: %v]", key, buf, err)
    }
    return buf
}

func (s *LevelDB) set(key []byte, val []byte) error {
    s.log.Printf("Set [%x] => [%x]", key, val)
    err := s.db.Set(key, val, &db.WriteOptions{true})
    if err != nil {
        s.log.Printf("Cause: %v", err)
    }
    return err
}

func (s *LevelDB) del(key []byte) error {
    s.log.Printf("Del [%x]", key)
    err := s.db.Delete(key, &db.WriteOptions{true})
    if err != nil {
        s.log.Printf("Cause: %v", err)
    }
    return err
}

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
    ensure(butils.WriteByte(buf, prefix_BILL))
    ensure(number.Write(buf))
    return buf.Bytes()
}