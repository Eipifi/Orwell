package hash
import "sort"

type IdList []ID

func (s IdList) Len() int {
    return len(s)
}

func (s IdList) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s IdList) Less(i, j int) bool {
    return Compare(s[i], s[j]) < 0
}

func (s IdList) Get(i int) ID {
    i = i % s.Len()
    if i < 0 { i += s.Len() }
    return s[i]
}

func (s IdList) ClosestNotSmaller(id ID) int {
    return sort.Search(s.Len(), func(i int) bool {
        return Compare(id, s[i]) <= 0
    })
}

func (s IdList) Contains(id ID) bool {
    return Compare(s.Get(s.ClosestNotSmaller(id)), id) == 0
}

func (s IdList) NClosest(id ID, n int) []ID {
    if s.Len() <= n {
        ret := make([]ID, s.Len())
        copy(ret, s)
        return ret
    }
    ret := make([]ID, n)
    ptrR := s.ClosestNotSmaller(id)
    ptrL := ptrR - 1
    for i := 0; i < n; i++ {
        if LeftCloser(s.Get(ptrL), id, s.Get(ptrR)) {
            ret[i] = s.Get(ptrL)
            ptrL -= 1
        } else {
            ret[i] = s.Get(ptrR)
            ptrR += 1
        }
    }
    return ret
}