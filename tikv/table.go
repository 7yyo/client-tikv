package tikv

const (
	TikvTable_ = "t"
)

type table struct {
	name string
}

func checkTblName(i Executor) error {
	switch e := i.(type) {
	case deleteExecutor:
	case insertExecutor:
	case selectExecutor:
	case truncateExecutor:
		if e.table.name != TikvTable_ {
			return errTableName(e.table.name)
		}
	default:
	}
	return nil
}
