package util

var (
	Db    IDbUtil    = &dbUtil{}
	Array IArrayUtil = &arrayUtil{}
	Math  IMathUtil  = &mathUtil{}
	Test  ITestUtil  = &testUtil{}
)
