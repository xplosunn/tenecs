// ##################################################################
// # The signatures of this file are generated via code-generation. #
// # Check gen.go                                                   #
// ##################################################################
package standard_library

func tenecs_compare_eq() Function {
	return function(
		params("first", "second"),
		body(`
function areDeeplyEqual(obj1, obj2) {
  if (typeof obj1 === 'function') return false;
  if (obj1 === obj2) return true;

  if (Array.isArray(obj1) && Array.isArray(obj2)) {

    if(obj1.length !== obj2.length) return false;
    
    return obj1.every((elem, index) => {
      return areDeeplyEqual(elem, obj2[index]);
    })


  }

  if(typeof obj1 === "object" && typeof obj2 === "object" && obj1 !== null && obj2 !== null) {
    if(Array.isArray(obj1) || Array.isArray(obj2)) return false;
    
    const keys1 = Object.keys(obj1)
    const keys2 = Object.keys(obj2)

    if(keys1.length !== keys2.length || !keys1.every(key => keys2.includes(key))) return false;
      
    for(let key in obj1) {
       let isEqual = areDeeplyEqual(obj1[key], obj2[key])
       if (!isEqual) { return false; }
    }

    return true;
    
  }

  return false;
}
return areDeeplyEqual(first, second);
`),
	)
}
