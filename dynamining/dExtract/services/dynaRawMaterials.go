package dExtract

import "io/ioutil"

// 1. Raw data : 원시 자료
//   1) Raw data collection: 원시 자료 수집
//   2) Excel, PDF, ... etc Convert to learning materials (*.rmf CSV or text file )

func ConvertRawDataToLMaterials() {
	loadRawData()
	saveCandidate()
}

// load Raw data
// Excel file 을 읽어 들이는 작업
func loadRawData() {
	rData, err := ioutil.ReadFile("/Data/")

}

// Convert Raw data file to CSV file
// Excel file 에서 CSV file 로 변경
func convertRawDataToCSV() {

}

// Convert Raw data file to Text file
// Excel file 에서 text file 로 변경
func convertRawDateToText() {

}

func saveCandidate() {

}