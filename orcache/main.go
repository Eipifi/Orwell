package main


func main() {
    m, err := NewManager(":1984")
    if err == nil {
        m.Lifecycle()
    }
}