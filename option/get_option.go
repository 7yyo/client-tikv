package option

//func GetOption(args []string) []prompt.Suggest {
//
//	var s []prompt.Suggest
//
//	switch len(args) {
//	case 2:
//		s = []prompt.Suggest{
//			{Text: syntax.Kv, Description: "If you want to directly query kv pairs, please use `kv` as the column name."},
//			{Text: " ", Description: "If the value you want to query is standard json, you can query according to the label of json, enter the label you want to query here."},
//		}
//	case 3:
//		s = []prompt.Suggest{
//			{Text: syntax.From, Description: "No need to explain it."},
//		}
//	case 4:
//		s = []prompt.Suggest{
//			{Text: syntax.Tikv, Description: "We have defined that the table must be `tikv`."},
//		}
//	case 5:
//		s = []prompt.Suggest{
//			{Text: syntax.Where, Description: "No need to explain it."},
//		}
//	case 6:
//		s = []prompt.Suggest{
//			{Text: syntax.Key, Description: "Enter the key you want to query."},
//		}
//	case 7:
//		s = []prompt.Suggest{
//			{Text: syntax.Eq, Description: "Equal sign."},
//			{Text: syntax.In, Description: "Scope sign."},
//		}
//	case 8:
//		if args[6] == syntax.Eq {
//			s = []prompt.Suggest{
//				{Text: syntax.Apostrophe, Description: ""},
//			}
//		} else {
//			s = []prompt.Suggest{
//				{Text: syntax.BracketsIn, Description: ""},
//			}
//		}
//	}
//	return s
//}
