function typeNameCollapseRecursive(input) {
	const str = typeof input === 'string' ? input : null
	if (str) {
		return str
	}
	const array = Array.isArray(input) ? input : null
	if (!array || array.length === 0) {
		throw new Error("typeNameCollapseRecursive input is not an array")
	}
	let resultOr = array.map(typeNameCollapseRecursive)
	resultOr.sort()
	resultOr = resultOr.filter((value, index, self) => self.indexOf(value) === index)
	return resultOr.join(' | ')
}