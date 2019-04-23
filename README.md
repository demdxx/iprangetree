# iprangetree collection

    @License MIT

The simple IP range-tree collection which provides fast access to inMemory storage of searching by IP range.
Usually can be used in GeoIP services to fast IP lookup in a big amount of IP addresses.
The project used in the real high-load project with DB of >13,000,000 records and average time of response ~10µs

## Example

```go
tree := iprangetree.New(2)

err := tree.AddRangeByString("86.100.32.0-86.100.32.255", []int{10})
if err != nil {
  log.Fatal(err)
}

fmt.Println(tree.LookupByString("86.100.32.10"))
```
