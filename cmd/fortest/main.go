package main

import (
	//"sort"
	"fmt"
	"strings"
	//"time"
	//"regexp"
	//"github.comcast.com/viper-cog/clog"
	//"sort"
	"regexp"
	//"reflect"
	//"github.comcast.com/viper-cog/clog"
	"sort"
	//"time"
)

type LongIDs []string

func (s LongIDs) Len() int {
	return len(s)
}

func (s LongIDs) Less(i, j int) bool {
	return getCandidateSequentialID(s[i]) < getCandidateSequentialID(s[j])
}

func (s LongIDs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}



func checkCandidatesIDFormat(candidateID string) bool {
	candidateRegex := CandidateLongIDRegex
	r := regexp.MustCompile(candidateRegex)
	match := r.MatchString(candidateID)
	if !match {
		//TODO: Do some logging here, can we use log globally?
		return false
	}
	return true
}


const (
	Candidate_Long_ID_Sample  = "_c_2053be8377d8405788aff3c9578e45f3-le_0000000005"
	Candidate_Short_ID_Sample = "le_0000000005"
	Candidate_Prefix          = "_c_"

	Cadidate_ID_Sample = "_c_2053be8377d8405788aff3c9578e45f3-le_0000000005"
	Sequential_ID_Prefix = "le_"
)


func main() {



	candidates := []string{"_c_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb-le_0000000000", "_c_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-le_0000000001", "le_0000000002", "_c_cccccccccccccccccccccccccccccccc-le_0000000003", "le_0000000004", }

	sort.Strings(candidates)

	fmt.Printf("==== %+v \n", candidates)

	sort.Sort(SortByCandidateSequentialIDASC(candidates))

	fmt.Printf("==== %+v \n", candidates)


	//shortCndtID1 := "_c_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-le_0000000001"
	//shortCndtID2 := "_c_cccccccccccccccccccccccccccccccc-le_0000000002"
	//shortCndtID3 := "_c_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb-le_0000000000"
	//shortCndtID4 := "_c_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb-le_0000000004"

	//candidates := []string{"_c_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb-le_0000000000", "_c_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-le_0000000001","_c_cccccccccccccccccccccccccccccccc-le_0000000002"}
	//candidates := []string{"_c_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb-le_0000000000","_c_cccccccccccccccccccccccccccccccc-le_0000000002", "_c_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-le_0000000001"}
/*	idx1 := sort.Search(len(candidates), func(i int) bool {
				return getCandidateSequentialID(candidates[i]) >= getCandidateSequentialID(shortCndtID1)
			})
	idx2 := sort.Search(len(candidates), func(i int) bool {
		return getCandidateSequentialID(candidates[i]) >= getCandidateSequentialID(shortCndtID2)
	})
	idx3 := sort.Search(len(candidates), func(i int) bool {
		return getCandidateSequentialID(candidates[i]) >= getCandidateSequentialID(shortCndtID3)
	})
	idx4 := sort.Search(len(candidates), func(i int) bool {
		return getCandidateSequentialID(candidates[i]) >= getCandidateSequentialID(shortCndtID4)
	})
	fmt.Print(idx1)
	fmt.Print(idx2)
	fmt.Print(idx3)
	fmt.Print(idx4)

	fmt.Printf("@@@@ %v", strings.HasPrefix(shortCndtID1, Candidate_Prefix))*/


}

func getID(input string)string{
	return input[4:]
}

const (
	electionOver = "DONE"
	CandidateLongIDRegex = "^_c_[a-z0-9]{32}-(le_[0-9]{10})$"
)

//Sort the candidates slice by short candidate ID (ASC order).
type SortByCandidateSequentialIDASC []string

func (s SortByCandidateSequentialIDASC) Len() int {
	return len(s)
}

func (s SortByCandidateSequentialIDASC) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortByCandidateSequentialIDASC) Less(i, j int) bool {
	return getCandidateSequentialID(s[i]) < getCandidateSequentialID(s[j])
}

//Get candidate sequential ID
//For example, input: "_c_2053be8377d8405788aff3c9578e45f3-le_0000000005" output: "le_0000000005"
//This API simply get the short ID from the long one, please validate the candidate ID format before calling this API.
func getCandidateSequentialID(candidateLongID string) string {
	if len(candidateLongID) != len(Candidate_Long_ID_Sample) {
		return ""
	}
	return candidateLongID[len(Candidate_Long_ID_Sample)-len(Candidate_Short_ID_Sample):]
}

//Validate candidate long ID format.
//From the previous test result, using regex to do validation cause performance issue and fail the integration test.
//So here, do simple validation on the length and prefix.
func checkCandidateIDFormat(candidateLongID string) bool {
	if len(candidateLongID) != len(Candidate_Long_ID_Sample) {
		return false
	}

	i := strings.LastIndex(candidateLongID, Sequential_ID_Prefix)
	if i != (len(Candidate_Long_ID_Sample) - len(Candidate_Short_ID_Sample)) {
		return false
	}

	i = strings.Index(candidateLongID, Candidate_Prefix)
	if i != 0 {
		return false
	}
	return true
}

