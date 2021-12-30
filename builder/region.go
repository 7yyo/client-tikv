package builder

type Regions struct {
	Count   int
	Regions []Region
}

type Region struct {
	Id               int
	Start_key        string
	End_key          string
	Peers            []Peer
	Leader           Leader
	Read_bytes       int
	Written_bytes    int
	Read_keys        int
	Approximate_size int
	Approximate_keys int
}

type Leader struct {
	Id        int
	Store_id  int
	Role_name string
}

type Peer struct {
	Id        int
	Store_id  int
	Role_name string
}

type RegionTable struct {
	REGION_ID        int
	START_KEY        string
	END_KEY          string
	LEADER_ID        int
	LEADER_STORE_ID  int
	PEERS            string
	WRITTEN_BYTES    int
	READ_BYTES       int
	APPROXIMATE_SIZE int
	APPROXIMATE_KEYS int
}
