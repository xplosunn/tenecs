package godsl

type GoDSL interface {
	sealedGoDSL()
}

func exhaustiveSwitch(goDSL GoDSL) (*Expression, *Statement) {
	caseExpression, ok := goDSL.(Expression)
	if ok {
		return &caseExpression, nil
	}
	caseStatement, ok := goDSL.(Statement)
	if ok {
		return nil, &caseStatement
	}
	return nil, nil
}
