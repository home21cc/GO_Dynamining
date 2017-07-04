package dExtract

// 4. Real data (*.xlsx -> *.arf), (*.arf -> *.xlsx), ( *.xlsx -> database )
//   1) Copy Real data to RealData folder
//   2) Load Semantic dictionary
//   2) Real data read
//   3) Real data classification
//   4) Extract data
//   5) compare extracted data to Semantic dictionary
//   6) Crate conform data
//   7) Update Semantic Data
//   8) Apply other system : ERP, Estimation system  or etc

func RealExtraction() {
	readRealMaterials()
	createRealCandidate()
	loadDictionary()
	miningRealData()
	reinforcementDictionary()
	sendOtherSystem()
}

// load real data
func readRealMaterials() {

}

//
func createRealCandidate() {

}


// Extracting Real Data
func miningRealData() {

}

// memory to Real file  ( Excel )
func saveRealFile() {

}

func sendOtherSystem() {

}

