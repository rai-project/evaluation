package evaluation

func uptoIndex(arry []interface{}, idx int) int {
	if len(arry) <= idx {
		return len(arry) - 1
	}
	return idx
}
