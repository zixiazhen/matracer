package main

import (
	"sort"
	"fmt"
	"strings"
)

func main() {
	candidates := []string{"_c_bbbb-le_0000000000", "_c_aaaa-le_0000000001", "_c_cccc-le_0000000002"}
	sort.Sort(CandidatesByZKSequenceNum(candidates))
	fmt.Printf(" New sort result: %+v \n ", candidates)

	getShortCandidateID("_c_aaaa-le_0000000001")
	fmt.Printf(" extract seq num: %+v \n ", getShortCandidateID("_c_aaaa-le_0000000001"))
}

type CandidatesByZKSequenceNum []string

func (s CandidatesByZKSequenceNum) Len() int {
	return len(s)
}

func (s CandidatesByZKSequenceNum) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s CandidatesByZKSequenceNum) Less(i, j int) bool {
	return getShortCandidateID(s[i]) < getShortCandidateID(s[j])
}

func getShortCandidateID(fullpath string) string {
	if len(fullpath) == 0 {
		return ""
	}
	index := strings.Index(fullpath, "le_")
	if index < 0 {
		return ""
	}
	return fullpath[index:]
}
