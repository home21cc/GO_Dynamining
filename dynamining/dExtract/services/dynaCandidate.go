package dExtract

// 2. Candidate data : 후보자료 (*.cdf)
//   1) Candidate date extraction : 후보 자료 추출, 원시 자료에서 Frequency 반영 자료 추출
//   2) Candidate data Update : Candidate data, Frequency, Use state update

func ConvertLMaterialsToCandidate() {
	loadLearningMaterials()
	loadCandidate()
	updateCandidate()
}

// load Learning Materials Data (CSV, text file)
func loadLearningMaterials() {
}

// UpdateCandidateData
// add context, update Frequency
// Case study
// 1. Header, Body 가 하나의 Sheet에  있을경우
//   1) Head , body 구분 로직 추가
//   2) 이후 는 2. 항과 동일하게 운영
// 2. Head, Body Sheet 분리 되어 있을 경우
//   1) Extract Header
//      -
//   2) Extract Body
//      - Column properties 구분하여 추출 해야 함
//        ex) number, Quantity, Delivery ... etc 등은 Column Information 그대로 읽기
//            Item, Material ... 등은 Column 에서 단어, 절 등을 구분해야 함
//   3) 단어, 절 등으로 구분되는 건에 대하여 개발 로직 정의
//      - 단어 : 접두 %context, 접미어 context%, 포함된 단어 %context%, 일부포함 ?
//      - 절 :
func updateCandidate() {

}

func loadCandidate() {

}

// Determine the use of candidate data
func determineCandidate() {
}