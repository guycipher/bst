# Lockless BST Package
A Go lang implementation of a lockless binary search tree.

## Features
- `Get`, `Put`, `Delete`, `Remove`, `Range`, `NGet`, `NRange`, `GreaterThan`, `GreaterThanEq`, `LessThan`, `LessThanEq` methods
- Lockless implementation
- Thread safe
- Very fast

## Usage

### Put
```go
tree := bst.New()
tree.Put([]byte("key"), []byte("value"))
```

### Get
```go
key := tree.Get([]byte("key"))
if key == nil {
    fmt.Println("Key not found")
}
```

### Delete
```go
tree.Delete([]byte("key"))
```

### Remove
```go
tree.Remove([]byte("key"), []byte("value to remove"))
```

### Range
```go
keys := tree.Range([]byte("key1"), []byte("key2"))
```

### NGet
```go
keys := tree.NGet([]byte("key"))
```

### NRange
```go
keys := tree.NRange([]byte("key1"), []byte("key2"))
```

### GreaterThan
```go
keys := tree.GreaterThan([]byte("key"))
```

### GreaterThanEq
```go
keys := tree.GreaterThanEq([]byte("key"))
```

### LessThan
```go
keys := tree.LessThan([]byte("key"))
```

### LessThanEq
```go
keys := tree.LessThanEq([]byte("key"))
```



