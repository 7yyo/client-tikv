package syntax

type KvPair struct {
	Key   string
	Value string
}

func GetKeys(kvPairs []KvPair) [][]byte {
	var keys [][]byte
	for _, kv := range kvPairs {
		keys = append(keys, []byte(kv.Key))
	}
	return keys
}
