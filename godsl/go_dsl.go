package godsl

//TODO
// - for
// - array

type GoDSL interface {
	sealedGoDSL()
}

func exhaustiveSwitch(goDSL GoDSL) (*Expression, *Statement, *TopLevelStatement) {
	caseExpression, ok := goDSL.(Expression)
	if ok {
		return &caseExpression, nil, nil
	}
	caseStatement, ok := goDSL.(Statement)
	if ok {
		return nil, &caseStatement, nil
	}
	caseTopLevelStatement, ok := goDSL.(TopLevelStatement)
	if ok {
		return nil, nil, &caseTopLevelStatement
	}
	return nil, nil, nil
}
