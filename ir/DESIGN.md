# IR Package Design Document

## Overview

The `ir` (Intermediate Representation) package transforms the typed AST into a language-agnostic representation for code generation. It handles complex language features that don't translate directly to target languages. The result should allow direct translation into a language such as JavaScript, without any additional logic. As such, Javascript is used for the examples below.

## Key Problems and Solutions

### 1. Top-Level Variables

**Problem**: Top-level variables in Tenecs can reference other variables regardless of declaration order, but target languages may require variables to be declared before use. This can be problematic when variables contain direct function invocations that reference other variables.

**Solution**: Convert all top-level variables to functions without arguments
- Top-level variables become functions that return the computed value
- This allows functions to call each other regardless of declaration order
- Functions can reference any other function without ordering constraints

**JavaScript Example**:
```tenecs
// Original Tenecs code
package main

str := "hello"
sameAs := (arg: String): String => { arg }
str2 := sameAs(str)  // references str
str3 := sameAs(str2) // references str2
```

```javascript
// Generated JavaScript
function main__str() {
  return { $type: "String", value: "hello" };
}

function main__sameAs() {
  return function(generics, arg) {
    return arg;
  }
}

function main__str2() {
  return main__sameAs()([], main__str());  // can call main__str() regardless of declaration order
}

function main__str3() {
  return main__sameAs()([], main__str2());  // can call main__str2() regardless of declaration order
}
```

### 2. Variable Name Conflicts

**Problem**: Variables can have names that aren't ideal in the target language, such as:
- Two variables from different packages having the same name
- A variable having the same name as a keyword in the target language

**Solution**: Systematic naming with prefixes and separators
- Local variables: `_name` (underscore prefix)
- Package variables: `package__name` (double underscore separator) - can't collide because tenecs variable names can't have double underscore in their name
- Generic variables: `genericsN` (numbered, no underscore prefix) - can't collide because package names in tenecs don't allow numbers

**JavaScript Example**:
```tenecs
// Original Tenecs code
package main

myFunction := () => {
  myVar := "hello"
  myVar
}
```

```javascript
// Generated JavaScript
var main__myFunction = function() {
  return function(generics) {
    var _myVar = { $type: "String", value: "hello" };
    return _myVar;
  }
};
```

### 3. Expression-Statement Duality

**Problem**: Everything in Tenecs is an expression, but target languages distinguish between expressions and statements. Some expressions need to be used in statement contexts (like blocks), while others need to be used as values.

**Solution**: Convert expressions to statements when needed
- Statements that need to be expressions become function calls with explicit returns
- Blocks maintain explicit return statements for the last expression

**Not solved**: Expressions in statement position
- Expressions in statement contexts are not addressed in this module and left to codegen to handle

**JavaScript Examples**:

```tenecs
// Statement that needs to be an expression
package main

getValue := (condition: Boolean): String => {
  if condition { "true" } else { "false" } // if expression used as value
}
```

```javascript
// Generated JavaScript
function main__getValue() {
  return function(generics, condition) {
    return (function() {
      if (condition["$type"] === "Boolean") {
        return { $type: "String", value: "true" };
      } else {
        return { $type: "String", value: "false" };
      }
    })();
  };
}
```

```tenecs
// Block with multiple expressions
package main

process := (value: String): String => {
  temp := value + " world" // expression as statement
  temp                     // last expression becomes return
}
```

```javascript
// Generated JavaScript
function main__process() {
  return function(generics, value) {
    return (function() {
      var _temp = { $type: "String", value: value.value + " world" };
      return { $type: "String", value: _temp.value };
    })();
  };
}
```


### 4. Type Information and Semantic Preservation

**Problem**: Target languages don't have the same type system or semantic constructs as Tenecs. We need to preserve type information, generic types, pattern matching, and expression semantics for runtime operations.

**Solution**: Runtime type system with object wrappers and explicit transformations
- All values become typed objects with `$type` and `value` fields
- Generic types are passed as arrays and accessed via indexing
- Pattern matching becomes chained if-else statements with type checking

**JavaScript Examples**:

```tenecs
// Generic types
package main

identity := <T>(value: T): T => {
  value
}

result := identity<String>("hello")
```

```javascript
// Generated JavaScript
function main__identity(generics, value) {
  return (function(generics, value) {
    return value;
  })
}

var main__result = main__identity()([["String"]], { $type: "String", value: "hello" });
```

```tenecs
// Pattern matching
package main

describe := (value: String | Int): String => {
  when value {
    is String => { "is string" }
    is Int => { "is number" }
  }
}

result := describe("hello")
```

```javascript
// Generated JavaScript
function main__describe() {
  return (function(generics, value) {
    var __over = value;
    if (Array.isArray("String") ? "String".includes(__over.$type) : "String" == __over.$type) {
      return { $type: "String", value: "is string" };
    } else {
      if (Array.isArray("Int") ? "Int".includes(__over.$type) : "Int" == __over.$type) {
        return { $type: "String", value: "is number" };
      } else {
        return null;
      }
    }
  })();
}

var main__result = main__describe()([], { $type: "String", value: "hello" });
```

```tenecs
// Generics and pattern matching together
package main

isLeftOrRight := <L, R>(value: L | R): String => {
  when value {
    is L => { "is left" }
    is R => { "is right" }
  }
}

result1 := isLeftOrRight<Int, String>(42)
result2 := isLeftOrRight<Int | Boolean, String>(42)
```

```javascript
// Generated JavaScript
function main__isLeftOrRight() {
  return (function(generics, value) {
    var __over = value;
    if (Array.isArray(generics[0]) ? generics[0].includes(__over.$type) : generics[0] == __over.$type) {
      return { $type: "String", value: "is left" };
    } else {
      if (Array.isArray(generics[1]) ? generics[1].includes(__over.$type) : generics[1] == __over.$type) {
        return { $type: "String", value: "is right" };
      } else {
        return null;
      }
    }
  })();
}

var main__result1 = main__isLeftOrRight([["Int"], ["String"]], { $type: "Int", value: 42 });
var main__result2 = main__isLeftOrRight([["Int", "Boolean"], ["String"]], { $type: "Int", value: 42 });
```