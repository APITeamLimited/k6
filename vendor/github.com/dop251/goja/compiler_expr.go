package goja

import (
	"fmt"
	"regexp"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/dop251/goja/unistring"
)

var (
	octalRegexp = regexp.MustCompile(`^0[0-7]`)
)

type compiledExpr interface ***REMOVED***
	emitGetter(putOnStack bool)
	emitSetter(valueExpr compiledExpr, putOnStack bool)
	emitUnary(prepare, body func(), postfix, putOnStack bool)
	deleteExpr() compiledExpr
	constant() bool
	addSrcMap()
***REMOVED***

type compiledExprOrRef interface ***REMOVED***
	compiledExpr
	emitGetterOrRef()
***REMOVED***

type compiledCallExpr struct ***REMOVED***
	baseCompiledExpr
	args   []compiledExpr
	callee compiledExpr
***REMOVED***

type compiledObjectLiteral struct ***REMOVED***
	baseCompiledExpr
	expr *ast.ObjectLiteral
***REMOVED***

type compiledArrayLiteral struct ***REMOVED***
	baseCompiledExpr
	expr *ast.ArrayLiteral
***REMOVED***

type compiledRegexpLiteral struct ***REMOVED***
	baseCompiledExpr
	expr *ast.RegExpLiteral
***REMOVED***

type compiledLiteral struct ***REMOVED***
	baseCompiledExpr
	val Value
***REMOVED***

type compiledAssignExpr struct ***REMOVED***
	baseCompiledExpr
	left, right compiledExpr
	operator    token.Token
***REMOVED***

type deleteGlobalExpr struct ***REMOVED***
	baseCompiledExpr
	name unistring.String
***REMOVED***

type deleteVarExpr struct ***REMOVED***
	baseCompiledExpr
	name unistring.String
***REMOVED***

type deletePropExpr struct ***REMOVED***
	baseCompiledExpr
	left compiledExpr
	name unistring.String
***REMOVED***

type deleteElemExpr struct ***REMOVED***
	baseCompiledExpr
	left, member compiledExpr
***REMOVED***

type constantExpr struct ***REMOVED***
	baseCompiledExpr
	val Value
***REMOVED***

type baseCompiledExpr struct ***REMOVED***
	c      *compiler
	offset int
***REMOVED***

type compiledIdentifierExpr struct ***REMOVED***
	baseCompiledExpr
	name unistring.String
***REMOVED***

type compiledFunctionLiteral struct ***REMOVED***
	baseCompiledExpr
	expr    *ast.FunctionLiteral
	lhsName unistring.String
	isExpr  bool
	strict  bool
***REMOVED***

type compiledBracketExpr struct ***REMOVED***
	baseCompiledExpr
	left, member compiledExpr
***REMOVED***

type compiledThisExpr struct ***REMOVED***
	baseCompiledExpr
***REMOVED***

type compiledNewExpr struct ***REMOVED***
	baseCompiledExpr
	callee compiledExpr
	args   []compiledExpr
***REMOVED***

type compiledNewTarget struct ***REMOVED***
	baseCompiledExpr
***REMOVED***

type compiledSequenceExpr struct ***REMOVED***
	baseCompiledExpr
	sequence []compiledExpr
***REMOVED***

type compiledUnaryExpr struct ***REMOVED***
	baseCompiledExpr
	operand  compiledExpr
	operator token.Token
	postfix  bool
***REMOVED***

type compiledConditionalExpr struct ***REMOVED***
	baseCompiledExpr
	test, consequent, alternate compiledExpr
***REMOVED***

type compiledLogicalOr struct ***REMOVED***
	baseCompiledExpr
	left, right compiledExpr
***REMOVED***

type compiledLogicalAnd struct ***REMOVED***
	baseCompiledExpr
	left, right compiledExpr
***REMOVED***

type compiledBinaryExpr struct ***REMOVED***
	baseCompiledExpr
	left, right compiledExpr
	operator    token.Token
***REMOVED***

type compiledVariableExpr struct ***REMOVED***
	baseCompiledExpr
	name        unistring.String
	initializer compiledExpr
***REMOVED***

type compiledEnumGetExpr struct ***REMOVED***
	baseCompiledExpr
***REMOVED***

type defaultDeleteExpr struct ***REMOVED***
	baseCompiledExpr
	expr compiledExpr
***REMOVED***

func (e *defaultDeleteExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.expr.emitGetter(false)
	if putOnStack ***REMOVED***
		e.c.emit(loadVal(e.c.p.defineLiteralValue(valueTrue)))
	***REMOVED***
***REMOVED***

func (c *compiler) compileExpression(v ast.Expression) compiledExpr ***REMOVED***
	// log.Printf("compileExpression: %T", v)
	switch v := v.(type) ***REMOVED***
	case nil:
		return nil
	case *ast.AssignExpression:
		return c.compileAssignExpression(v)
	case *ast.NumberLiteral:
		return c.compileNumberLiteral(v)
	case *ast.StringLiteral:
		return c.compileStringLiteral(v)
	case *ast.BooleanLiteral:
		return c.compileBooleanLiteral(v)
	case *ast.NullLiteral:
		r := &compiledLiteral***REMOVED***
			val: _null,
		***REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.Identifier:
		return c.compileIdentifierExpression(v)
	case *ast.CallExpression:
		return c.compileCallExpression(v)
	case *ast.ObjectLiteral:
		return c.compileObjectLiteral(v)
	case *ast.ArrayLiteral:
		return c.compileArrayLiteral(v)
	case *ast.RegExpLiteral:
		return c.compileRegexpLiteral(v)
	case *ast.VariableExpression:
		return c.compileVariableExpression(v)
	case *ast.BinaryExpression:
		return c.compileBinaryExpression(v)
	case *ast.UnaryExpression:
		return c.compileUnaryExpression(v)
	case *ast.ConditionalExpression:
		return c.compileConditionalExpression(v)
	case *ast.FunctionLiteral:
		return c.compileFunctionLiteral(v, true)
	case *ast.DotExpression:
		r := &compiledDotExpr***REMOVED***
			left: c.compileExpression(v.Left),
			name: v.Identifier.Name,
		***REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.BracketExpression:
		r := &compiledBracketExpr***REMOVED***
			left:   c.compileExpression(v.Left),
			member: c.compileExpression(v.Member),
		***REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.ThisExpression:
		r := &compiledThisExpr***REMOVED******REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.SequenceExpression:
		return c.compileSequenceExpression(v)
	case *ast.NewExpression:
		return c.compileNewExpression(v)
	case *ast.MetaProperty:
		return c.compileMetaProperty(v)
	default:
		panic(fmt.Errorf("Unknown expression type: %T", v))
	***REMOVED***
***REMOVED***

func (e *baseCompiledExpr) constant() bool ***REMOVED***
	return false
***REMOVED***

func (e *baseCompiledExpr) init(c *compiler, idx file.Idx) ***REMOVED***
	e.c = c
	e.offset = int(idx) - 1
***REMOVED***

func (e *baseCompiledExpr) emitSetter(compiledExpr, bool) ***REMOVED***
	e.c.throwSyntaxError(e.offset, "Not a valid left-value expression")
***REMOVED***

func (e *baseCompiledExpr) deleteExpr() compiledExpr ***REMOVED***
	r := &constantExpr***REMOVED***
		val: valueTrue,
	***REMOVED***
	r.init(e.c, file.Idx(e.offset+1))
	return r
***REMOVED***

func (e *baseCompiledExpr) emitUnary(func(), func(), bool, bool) ***REMOVED***
	e.c.throwSyntaxError(e.offset, "Not a valid left-value expression")
***REMOVED***

func (e *baseCompiledExpr) addSrcMap() ***REMOVED***
	if e.offset > 0 ***REMOVED***
		e.c.p.srcMap = append(e.c.p.srcMap, srcMapItem***REMOVED***pc: len(e.c.p.code), srcPos: e.offset***REMOVED***)
	***REMOVED***
***REMOVED***

func (e *constantExpr) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.addSrcMap()
		e.c.emit(loadVal(e.c.p.defineLiteralValue(e.val)))
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.addSrcMap()
	if b, noDynamics := e.c.scope.lookupName(e.name); noDynamics ***REMOVED***
		if b != nil ***REMOVED***
			if putOnStack ***REMOVED***
				b.emitGet()
			***REMOVED*** else ***REMOVED***
				b.emitGetP()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			panic("No dynamics and not found")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b != nil ***REMOVED***
			b.emitGetVar(false)
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadDynamic(e.name))
		***REMOVED***
		if !putOnStack ***REMOVED***
			e.c.emit(pop)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) emitGetterOrRef() ***REMOVED***
	e.addSrcMap()
	if b, noDynamics := e.c.scope.lookupName(e.name); noDynamics ***REMOVED***
		if b != nil ***REMOVED***
			b.emitGet()
		***REMOVED*** else ***REMOVED***
			panic("No dynamics and not found")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b != nil ***REMOVED***
			b.emitGetVar(false)
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadDynamicRef(e.name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) emitGetterAndCallee() ***REMOVED***
	e.addSrcMap()
	if b, noDynamics := e.c.scope.lookupName(e.name); noDynamics ***REMOVED***
		if b != nil ***REMOVED***
			e.c.emit(loadUndef)
			b.emitGet()
		***REMOVED*** else ***REMOVED***
			panic("No dynamics and not found")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b != nil ***REMOVED***
			b.emitGetVar(true)
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadDynamicCallee(e.name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) emitVarSetter1(name unistring.String, offset int, putOnStack bool, emitRight func(isRef bool)) ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierLName(name, offset)
	***REMOVED***

	if b, noDynamics := c.scope.lookupName(name); noDynamics ***REMOVED***
		emitRight(false)
		if b != nil ***REMOVED***
			if putOnStack ***REMOVED***
				b.emitSet()
			***REMOVED*** else ***REMOVED***
				b.emitSetP()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if c.scope.strict ***REMOVED***
				c.emit(setGlobalStrict(name))
			***REMOVED*** else ***REMOVED***
				c.emit(setGlobal(name))
			***REMOVED***
			if !putOnStack ***REMOVED***
				c.emit(pop)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b != nil ***REMOVED***
			b.emitResolveVar(c.scope.strict)
		***REMOVED*** else ***REMOVED***
			if c.scope.strict ***REMOVED***
				c.emit(resolveVar1Strict(name))
			***REMOVED*** else ***REMOVED***
				c.emit(resolveVar1(name))
			***REMOVED***
		***REMOVED***
		emitRight(true)
		if putOnStack ***REMOVED***
			c.emit(putValue)
		***REMOVED*** else ***REMOVED***
			c.emit(putValueP)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) emitVarSetter(name unistring.String, offset int, valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	c.emitVarSetter1(name, offset, putOnStack, func(bool) ***REMOVED***
		c.emitExpr(valueExpr, true)
	***REMOVED***)
***REMOVED***

func (e *compiledVariableExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.c.emitVarSetter(e.name, e.offset, valueExpr, putOnStack)
***REMOVED***

func (e *compiledIdentifierExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.c.emitVarSetter(e.name, e.offset, valueExpr, putOnStack)
***REMOVED***

func (e *compiledIdentifierExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.c.emitVarSetter1(e.name, e.offset, true, func(isRef bool) ***REMOVED***
			e.c.emit(loadUndef)
			if isRef ***REMOVED***
				e.c.emit(getValue)
			***REMOVED*** else ***REMOVED***
				e.emitGetter(true)
			***REMOVED***
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			if !postfix ***REMOVED***
				body()
			***REMOVED***
			e.c.emit(rdupN(1))
			if postfix ***REMOVED***
				body()
			***REMOVED***
		***REMOVED***)
		e.c.emit(pop)
	***REMOVED*** else ***REMOVED***
		e.c.emitVarSetter1(e.name, e.offset, false, func(isRef bool) ***REMOVED***
			if isRef ***REMOVED***
				e.c.emit(getValue)
			***REMOVED*** else ***REMOVED***
				e.emitGetter(true)
			***REMOVED***
			body()
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) deleteExpr() compiledExpr ***REMOVED***
	if e.c.scope.strict ***REMOVED***
		e.c.throwSyntaxError(e.offset, "Delete of an unqualified identifier in strict mode")
		panic("Unreachable")
	***REMOVED***
	if b, noDynamics := e.c.scope.lookupName(e.name); noDynamics ***REMOVED***
		if b == nil ***REMOVED***
			r := &deleteGlobalExpr***REMOVED***
				name: e.name,
			***REMOVED***
			r.init(e.c, file.Idx(0))
			return r
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b == nil ***REMOVED***
			r := &deleteVarExpr***REMOVED***
				name: e.name,
			***REMOVED***
			r.init(e.c, file.Idx(e.offset+1))
			return r
		***REMOVED***
	***REMOVED***
	r := &compiledLiteral***REMOVED***
		val: valueFalse,
	***REMOVED***
	r.init(e.c, file.Idx(e.offset+1))
	return r
***REMOVED***

type compiledDotExpr struct ***REMOVED***
	baseCompiledExpr
	left compiledExpr
	name unistring.String
***REMOVED***

func (e *compiledDotExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.addSrcMap()
	e.c.emit(getProp(e.name))
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledDotExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	valueExpr.emitGetter(true)
	if e.c.scope.strict ***REMOVED***
		if putOnStack ***REMOVED***
			e.c.emit(setPropStrict(e.name))
		***REMOVED*** else ***REMOVED***
			e.c.emit(setPropStrictP(e.name))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if putOnStack ***REMOVED***
			e.c.emit(setProp(e.name))
		***REMOVED*** else ***REMOVED***
			e.c.emit(setPropP(e.name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledDotExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	if !putOnStack ***REMOVED***
		e.left.emitGetter(true)
		e.c.emit(dup)
		e.c.emit(getProp(e.name))
		body()
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setPropStrict(e.name), pop)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setProp(e.name), pop)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !postfix ***REMOVED***
			e.left.emitGetter(true)
			e.c.emit(dup)
			e.c.emit(getProp(e.name))
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			body()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setPropStrict(e.name))
			***REMOVED*** else ***REMOVED***
				e.c.emit(setProp(e.name))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadUndef)
			e.left.emitGetter(true)
			e.c.emit(dup)
			e.c.emit(getProp(e.name))
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			e.c.emit(rdupN(2))
			body()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setPropStrict(e.name))
			***REMOVED*** else ***REMOVED***
				e.c.emit(setProp(e.name))
			***REMOVED***
			e.c.emit(pop)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledDotExpr) deleteExpr() compiledExpr ***REMOVED***
	r := &deletePropExpr***REMOVED***
		left: e.left,
		name: e.name,
	***REMOVED***
	r.init(e.c, file.Idx(0))
	return r
***REMOVED***

func (e *compiledBracketExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.member.emitGetter(true)
	e.addSrcMap()
	e.c.emit(getElem)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledBracketExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.member.emitGetter(true)
	valueExpr.emitGetter(true)
	if e.c.scope.strict ***REMOVED***
		if putOnStack ***REMOVED***
			e.c.emit(setElemStrict)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElemStrictP)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if putOnStack ***REMOVED***
			e.c.emit(setElem)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElemP)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledBracketExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	if !putOnStack ***REMOVED***
		e.left.emitGetter(true)
		e.member.emitGetter(true)
		e.c.emit(dupN(1), dupN(1))
		e.c.emit(getElem)
		body()
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setElemStrict, pop)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElem, pop)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !postfix ***REMOVED***
			e.left.emitGetter(true)
			e.member.emitGetter(true)
			e.c.emit(dupN(1), dupN(1))
			e.c.emit(getElem)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			body()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setElemStrict)
			***REMOVED*** else ***REMOVED***
				e.c.emit(setElem)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadUndef)
			e.left.emitGetter(true)
			e.member.emitGetter(true)
			e.c.emit(dupN(1), dupN(1))
			e.c.emit(getElem)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			e.c.emit(rdupN(3))
			body()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setElemStrict, pop)
			***REMOVED*** else ***REMOVED***
				e.c.emit(setElem, pop)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledBracketExpr) deleteExpr() compiledExpr ***REMOVED***
	r := &deleteElemExpr***REMOVED***
		left:   e.left,
		member: e.member,
	***REMOVED***
	r.init(e.c, file.Idx(0))
	return r
***REMOVED***

func (e *deleteElemExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.member.emitGetter(true)
	e.addSrcMap()
	if e.c.scope.strict ***REMOVED***
		e.c.emit(deleteElemStrict)
	***REMOVED*** else ***REMOVED***
		e.c.emit(deleteElem)
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *deletePropExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.addSrcMap()
	if e.c.scope.strict ***REMOVED***
		e.c.emit(deletePropStrict(e.name))
	***REMOVED*** else ***REMOVED***
		e.c.emit(deleteProp(e.name))
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *deleteVarExpr) emitGetter(putOnStack bool) ***REMOVED***
	/*if e.c.scope.strict ***REMOVED***
		e.c.throwSyntaxError(e.offset, "Delete of an unqualified identifier in strict mode")
		return
	***REMOVED****/
	e.c.emit(deleteVar(e.name))
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *deleteGlobalExpr) emitGetter(putOnStack bool) ***REMOVED***
	/*if e.c.scope.strict ***REMOVED***
		e.c.throwSyntaxError(e.offset, "Delete of an unqualified identifier in strict mode")
		return
	***REMOVED****/

	e.c.emit(deleteGlobal(e.name))
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledAssignExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.addSrcMap()
	switch e.operator ***REMOVED***
	case token.ASSIGN:
		if fn, ok := e.right.(*compiledFunctionLiteral); ok ***REMOVED***
			if fn.expr.Name == nil ***REMOVED***
				if id, ok := e.left.(*compiledIdentifierExpr); ok ***REMOVED***
					fn.lhsName = id.name
				***REMOVED***
			***REMOVED***
		***REMOVED***
		e.left.emitSetter(e.right, putOnStack)
	case token.PLUS:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(add)
		***REMOVED***, false, putOnStack)
	case token.MINUS:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(sub)
		***REMOVED***, false, putOnStack)
	case token.MULTIPLY:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(mul)
		***REMOVED***, false, putOnStack)
	case token.SLASH:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(div)
		***REMOVED***, false, putOnStack)
	case token.REMAINDER:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(mod)
		***REMOVED***, false, putOnStack)
	case token.OR:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(or)
		***REMOVED***, false, putOnStack)
	case token.AND:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(and)
		***REMOVED***, false, putOnStack)
	case token.EXCLUSIVE_OR:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(xor)
		***REMOVED***, false, putOnStack)
	case token.SHIFT_LEFT:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(sal)
		***REMOVED***, false, putOnStack)
	case token.SHIFT_RIGHT:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(sar)
		***REMOVED***, false, putOnStack)
	case token.UNSIGNED_SHIFT_RIGHT:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(shr)
		***REMOVED***, false, putOnStack)
	default:
		panic(fmt.Errorf("Unknown assign operator: %s", e.operator.String()))
	***REMOVED***
***REMOVED***

func (e *compiledLiteral) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.addSrcMap()
		e.c.emit(loadVal(e.c.p.defineLiteralValue(e.val)))
	***REMOVED***
***REMOVED***

func (e *compiledLiteral) constant() bool ***REMOVED***
	return true
***REMOVED***

func (e *compiledFunctionLiteral) emitGetter(putOnStack bool) ***REMOVED***
	savedPrg := e.c.p
	e.c.p = &Program***REMOVED***
		src: e.c.p.src,
	***REMOVED***
	e.c.newScope()
	e.c.scope.function = true

	var name unistring.String
	if e.expr.Name != nil ***REMOVED***
		name = e.expr.Name.Name
	***REMOVED*** else ***REMOVED***
		name = e.lhsName
	***REMOVED***

	if name != "" ***REMOVED***
		e.c.p.funcName = name
	***REMOVED***
	savedBlock := e.c.block
	defer func() ***REMOVED***
		e.c.block = savedBlock
	***REMOVED***()

	e.c.block = &block***REMOVED***
		typ: blockScope,
	***REMOVED***

	if !e.c.scope.strict ***REMOVED***
		e.c.scope.strict = e.strict
	***REMOVED***

	if e.c.scope.strict ***REMOVED***
		for _, item := range e.expr.ParameterList.List ***REMOVED***
			e.c.checkIdentifierName(item.Name, int(item.Idx)-1)
			e.c.checkIdentifierLName(item.Name, int(item.Idx)-1)
		***REMOVED***
	***REMOVED***

	length := len(e.expr.ParameterList.List)

	for _, item := range e.expr.ParameterList.List ***REMOVED***
		b, unique := e.c.scope.bindNameShadow(item.Name)
		if !unique && e.c.scope.strict ***REMOVED***
			e.c.throwSyntaxError(int(item.Idx)-1, "Strict mode function may not have duplicate parameter names (%s)", item.Name)
			return
		***REMOVED***
		b.isArg = true
		b.isVar = true
	***REMOVED***
	paramsCount := len(e.c.scope.bindings)
	e.c.scope.numArgs = paramsCount
	e.c.compileDeclList(e.expr.DeclarationList, true)
	body := e.expr.Body.List
	funcs := e.c.extractFunctions(body)
	e.c.createFunctionBindings(funcs)
	s := e.c.scope
	e.c.compileLexicalDeclarations(body, true)
	var calleeBinding *binding
	if e.isExpr && e.expr.Name != nil ***REMOVED***
		if b, created := s.bindName(e.expr.Name.Name); created ***REMOVED***
			calleeBinding = b
		***REMOVED***
	***REMOVED***
	preambleLen := 4 // enter, boxThis, createArgs, set
	e.c.p.code = make([]instruction, preambleLen, 8)

	if calleeBinding != nil ***REMOVED***
		e.c.emit(loadCallee)
		calleeBinding.emitSetP()
	***REMOVED***

	e.c.compileFunctions(funcs)
	e.c.compileStatements(body, false)

	var last ast.Statement
	if l := len(body); l > 0 ***REMOVED***
		last = body[l-1]
	***REMOVED***
	if _, ok := last.(*ast.ReturnStatement); !ok ***REMOVED***
		e.c.emit(loadUndef, ret)
	***REMOVED***

	delta := 0
	code := e.c.p.code

	if calleeBinding != nil && !s.isDynamic() && calleeBinding.useCount() == 1 ***REMOVED***
		s.deleteBinding(calleeBinding)
		preambleLen += 2
	***REMOVED***

	if (s.argsNeeded || s.isDynamic()) && !s.argsInStash ***REMOVED***
		s.moveArgsToStash()
	***REMOVED***

	if s.argsNeeded ***REMOVED***
		pos := preambleLen - 2
		delta += 2
		if s.strict ***REMOVED***
			code[pos] = createArgsStrict(length)
		***REMOVED*** else ***REMOVED***
			code[pos] = createArgs(length)
		***REMOVED***
		pos++
		b, _ := s.bindName("arguments")
		e.c.p.code = code[:pos]
		b.emitSetP()
		e.c.p.code = code
	***REMOVED***

	stashSize, stackSize := s.finaliseVarAlloc(0)

	if !s.strict && s.thisNeeded ***REMOVED***
		delta++
		code[preambleLen-delta] = boxThis
	***REMOVED***
	delta++
	delta = preambleLen - delta
	var enter instruction
	if stashSize > 0 || s.argsInStash ***REMOVED***
		enter1 := enterFunc***REMOVED***
			numArgs:     uint32(paramsCount),
			argsToStash: s.argsInStash,
			stashSize:   uint32(stashSize),
			stackSize:   uint32(stackSize),
			extensible:  s.dynamic,
		***REMOVED***
		if s.isDynamic() ***REMOVED***
			enter1.names = s.makeNamesMap()
		***REMOVED***
		enter = &enter1
	***REMOVED*** else ***REMOVED***
		enter = &enterFuncStashless***REMOVED***
			stackSize: uint32(stackSize),
			args:      uint32(paramsCount),
		***REMOVED***
	***REMOVED***
	code[delta] = enter
	if delta != 0 ***REMOVED***
		e.c.p.code = code[delta:]
		for i := range e.c.p.srcMap ***REMOVED***
			e.c.p.srcMap[i].pc -= delta
		***REMOVED***
		s.adjustBase(-delta)
	***REMOVED***

	strict := s.strict
	p := e.c.p
	// e.c.p.dumpCode()
	e.c.popScope()
	e.c.p = savedPrg
	e.c.emit(&newFunc***REMOVED***prg: p, length: uint32(length), name: name, srcStart: uint32(e.expr.Idx0() - 1), srcEnd: uint32(e.expr.Idx1() - 1), strict: strict***REMOVED***)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileFunctionLiteral(v *ast.FunctionLiteral, isExpr bool) *compiledFunctionLiteral ***REMOVED***
	strict := c.scope.strict || c.isStrictStatement(v.Body)
	if v.Name != nil && strict ***REMOVED***
		c.checkIdentifierLName(v.Name.Name, int(v.Name.Idx)-1)
	***REMOVED***
	r := &compiledFunctionLiteral***REMOVED***
		expr:   v,
		isExpr: isExpr,
		strict: strict,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledThisExpr) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.addSrcMap()
		scope := e.c.scope
		for ; scope != nil && !scope.function && !scope.eval; scope = scope.outer ***REMOVED***
		***REMOVED***

		if scope != nil ***REMOVED***
			scope.thisNeeded = true
			e.c.emit(loadStack(0))
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadGlobalObject)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledNewExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.callee.emitGetter(true)
	for _, expr := range e.args ***REMOVED***
		expr.emitGetter(true)
	***REMOVED***
	e.addSrcMap()
	e.c.emit(_new(len(e.args)))
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileNewExpression(v *ast.NewExpression) compiledExpr ***REMOVED***
	args := make([]compiledExpr, len(v.ArgumentList))
	for i, expr := range v.ArgumentList ***REMOVED***
		args[i] = c.compileExpression(expr)
	***REMOVED***
	r := &compiledNewExpr***REMOVED***
		callee: c.compileExpression(v.Callee),
		args:   args,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledNewTarget) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.addSrcMap()
		e.c.emit(loadNewTarget)
	***REMOVED***
***REMOVED***

func (c *compiler) compileMetaProperty(v *ast.MetaProperty) compiledExpr ***REMOVED***
	if v.Meta.Name == "new" || v.Property.Name != "target" ***REMOVED***
		r := &compiledNewTarget***REMOVED******REMOVED***
		r.init(c, v.Idx0())
		return r
	***REMOVED***
	c.throwSyntaxError(int(v.Idx)-1, "Unsupported meta property: %s.%s", v.Meta.Name, v.Property.Name)
	return nil
***REMOVED***

func (e *compiledSequenceExpr) emitGetter(putOnStack bool) ***REMOVED***
	if len(e.sequence) > 0 ***REMOVED***
		for i := 0; i < len(e.sequence)-1; i++ ***REMOVED***
			e.sequence[i].emitGetter(false)
		***REMOVED***
		e.sequence[len(e.sequence)-1].emitGetter(putOnStack)
	***REMOVED***
***REMOVED***

func (c *compiler) compileSequenceExpression(v *ast.SequenceExpression) compiledExpr ***REMOVED***
	s := make([]compiledExpr, len(v.Sequence))
	for i, expr := range v.Sequence ***REMOVED***
		s[i] = c.compileExpression(expr)
	***REMOVED***
	r := &compiledSequenceExpr***REMOVED***
		sequence: s,
	***REMOVED***
	var idx file.Idx
	if len(v.Sequence) > 0 ***REMOVED***
		idx = v.Idx0()
	***REMOVED***
	r.init(c, idx)
	return r
***REMOVED***

func (c *compiler) emitThrow(v Value) ***REMOVED***
	if o, ok := v.(*Object); ok ***REMOVED***
		t := nilSafe(o.self.getStr("name", nil)).toString().String()
		switch t ***REMOVED***
		case "TypeError":
			c.emit(loadDynamic(t))
			msg := o.self.getStr("message", nil)
			if msg != nil ***REMOVED***
				c.emit(loadVal(c.p.defineLiteralValue(msg)))
				c.emit(_new(1))
			***REMOVED*** else ***REMOVED***
				c.emit(_new(0))
			***REMOVED***
			c.emit(throw)
			return
		***REMOVED***
	***REMOVED***
	panic(fmt.Errorf("unknown exception type thrown while evaliating constant expression: %s", v.String()))
***REMOVED***

func (c *compiler) emitConst(expr compiledExpr, putOnStack bool) ***REMOVED***
	v, ex := c.evalConst(expr)
	if ex == nil ***REMOVED***
		if putOnStack ***REMOVED***
			c.emit(loadVal(c.p.defineLiteralValue(v)))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.emitThrow(ex.val)
	***REMOVED***
***REMOVED***

func (c *compiler) emitExpr(expr compiledExpr, putOnStack bool) ***REMOVED***
	if expr.constant() ***REMOVED***
		c.emitConst(expr, putOnStack)
	***REMOVED*** else ***REMOVED***
		expr.emitGetter(putOnStack)
	***REMOVED***
***REMOVED***

func (c *compiler) evalConst(expr compiledExpr) (Value, *Exception) ***REMOVED***
	if expr, ok := expr.(*compiledLiteral); ok ***REMOVED***
		return expr.val, nil
	***REMOVED***
	if c.evalVM == nil ***REMOVED***
		c.evalVM = New().vm
	***REMOVED***
	var savedPrg *Program
	createdPrg := false
	if c.evalVM.prg == nil ***REMOVED***
		c.evalVM.prg = &Program***REMOVED******REMOVED***
		savedPrg = c.p
		c.p = c.evalVM.prg
		createdPrg = true
	***REMOVED***
	savedPc := len(c.p.code)
	expr.emitGetter(true)
	c.emit(halt)
	c.evalVM.pc = savedPc
	ex := c.evalVM.runTry()
	if createdPrg ***REMOVED***
		c.evalVM.prg = nil
		c.evalVM.pc = 0
		c.p = savedPrg
	***REMOVED*** else ***REMOVED***
		c.evalVM.prg.code = c.evalVM.prg.code[:savedPc]
		c.p.code = c.evalVM.prg.code
	***REMOVED***
	if ex == nil ***REMOVED***
		return c.evalVM.pop(), nil
	***REMOVED***
	return nil, ex
***REMOVED***

func (e *compiledUnaryExpr) constant() bool ***REMOVED***
	return e.operand.constant()
***REMOVED***

func (e *compiledUnaryExpr) emitGetter(putOnStack bool) ***REMOVED***
	var prepare, body func()

	toNumber := func() ***REMOVED***
		e.c.emit(toNumber)
	***REMOVED***

	switch e.operator ***REMOVED***
	case token.NOT:
		e.operand.emitGetter(true)
		e.c.emit(not)
		goto end
	case token.BITWISE_NOT:
		e.operand.emitGetter(true)
		e.c.emit(bnot)
		goto end
	case token.TYPEOF:
		if o, ok := e.operand.(compiledExprOrRef); ok ***REMOVED***
			o.emitGetterOrRef()
		***REMOVED*** else ***REMOVED***
			e.operand.emitGetter(true)
		***REMOVED***
		e.c.emit(typeof)
		goto end
	case token.DELETE:
		e.operand.deleteExpr().emitGetter(putOnStack)
		return
	case token.MINUS:
		e.c.emitExpr(e.operand, true)
		e.c.emit(neg)
		goto end
	case token.PLUS:
		e.c.emitExpr(e.operand, true)
		e.c.emit(plus)
		goto end
	case token.INCREMENT:
		prepare = toNumber
		body = func() ***REMOVED***
			e.c.emit(inc)
		***REMOVED***
	case token.DECREMENT:
		prepare = toNumber
		body = func() ***REMOVED***
			e.c.emit(dec)
		***REMOVED***
	case token.VOID:
		e.c.emitExpr(e.operand, false)
		if putOnStack ***REMOVED***
			e.c.emit(loadUndef)
		***REMOVED***
		return
	default:
		panic(fmt.Errorf("Unknown unary operator: %s", e.operator.String()))
	***REMOVED***

	e.operand.emitUnary(prepare, body, e.postfix, putOnStack)
	return

end:
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileUnaryExpression(v *ast.UnaryExpression) compiledExpr ***REMOVED***
	r := &compiledUnaryExpr***REMOVED***
		operand:  c.compileExpression(v.Operand),
		operator: v.Operator,
		postfix:  v.Postfix,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledConditionalExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.test.emitGetter(true)
	j := len(e.c.p.code)
	e.c.emit(nil)
	e.consequent.emitGetter(putOnStack)
	j1 := len(e.c.p.code)
	e.c.emit(nil)
	e.c.p.code[j] = jne(len(e.c.p.code) - j)
	e.alternate.emitGetter(putOnStack)
	e.c.p.code[j1] = jump(len(e.c.p.code) - j1)
***REMOVED***

func (c *compiler) compileConditionalExpression(v *ast.ConditionalExpression) compiledExpr ***REMOVED***
	r := &compiledConditionalExpr***REMOVED***
		test:       c.compileExpression(v.Test),
		consequent: c.compileExpression(v.Consequent),
		alternate:  c.compileExpression(v.Alternate),
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledLogicalOr) constant() bool ***REMOVED***
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if v.ToBoolean() ***REMOVED***
				return true
			***REMOVED***
			return e.right.constant()
		***REMOVED*** else ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (e *compiledLogicalOr) emitGetter(putOnStack bool) ***REMOVED***
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if !v.ToBoolean() ***REMOVED***
				e.c.emitExpr(e.right, putOnStack)
			***REMOVED*** else ***REMOVED***
				if putOnStack ***REMOVED***
					e.c.emit(loadVal(e.c.p.defineLiteralValue(v)))
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emitThrow(ex.val)
		***REMOVED***
		return
	***REMOVED***
	e.c.emitExpr(e.left, true)
	j := len(e.c.p.code)
	e.addSrcMap()
	e.c.emit(nil)
	e.c.emit(pop)
	e.c.emitExpr(e.right, true)
	e.c.p.code[j] = jeq1(len(e.c.p.code) - j)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledLogicalAnd) constant() bool ***REMOVED***
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if !v.ToBoolean() ***REMOVED***
				return true
			***REMOVED*** else ***REMOVED***
				return e.right.constant()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (e *compiledLogicalAnd) emitGetter(putOnStack bool) ***REMOVED***
	var j int
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if !v.ToBoolean() ***REMOVED***
				e.c.emit(loadVal(e.c.p.defineLiteralValue(v)))
			***REMOVED*** else ***REMOVED***
				e.c.emitExpr(e.right, putOnStack)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emitThrow(ex.val)
		***REMOVED***
		return
	***REMOVED***
	e.left.emitGetter(true)
	j = len(e.c.p.code)
	e.addSrcMap()
	e.c.emit(nil)
	e.c.emit(pop)
	e.c.emitExpr(e.right, true)
	e.c.p.code[j] = jneq1(len(e.c.p.code) - j)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledBinaryExpr) constant() bool ***REMOVED***
	return e.left.constant() && e.right.constant()
***REMOVED***

func (e *compiledBinaryExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.c.emitExpr(e.left, true)
	e.c.emitExpr(e.right, true)
	e.addSrcMap()

	switch e.operator ***REMOVED***
	case token.LESS:
		e.c.emit(op_lt)
	case token.GREATER:
		e.c.emit(op_gt)
	case token.LESS_OR_EQUAL:
		e.c.emit(op_lte)
	case token.GREATER_OR_EQUAL:
		e.c.emit(op_gte)
	case token.EQUAL:
		e.c.emit(op_eq)
	case token.NOT_EQUAL:
		e.c.emit(op_neq)
	case token.STRICT_EQUAL:
		e.c.emit(op_strict_eq)
	case token.STRICT_NOT_EQUAL:
		e.c.emit(op_strict_neq)
	case token.PLUS:
		e.c.emit(add)
	case token.MINUS:
		e.c.emit(sub)
	case token.MULTIPLY:
		e.c.emit(mul)
	case token.SLASH:
		e.c.emit(div)
	case token.REMAINDER:
		e.c.emit(mod)
	case token.AND:
		e.c.emit(and)
	case token.OR:
		e.c.emit(or)
	case token.EXCLUSIVE_OR:
		e.c.emit(xor)
	case token.INSTANCEOF:
		e.c.emit(op_instanceof)
	case token.IN:
		e.c.emit(op_in)
	case token.SHIFT_LEFT:
		e.c.emit(sal)
	case token.SHIFT_RIGHT:
		e.c.emit(sar)
	case token.UNSIGNED_SHIFT_RIGHT:
		e.c.emit(shr)
	default:
		panic(fmt.Errorf("Unknown operator: %s", e.operator.String()))
	***REMOVED***

	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileBinaryExpression(v *ast.BinaryExpression) compiledExpr ***REMOVED***

	switch v.Operator ***REMOVED***
	case token.LOGICAL_OR:
		return c.compileLogicalOr(v.Left, v.Right, v.Idx0())
	case token.LOGICAL_AND:
		return c.compileLogicalAnd(v.Left, v.Right, v.Idx0())
	***REMOVED***

	r := &compiledBinaryExpr***REMOVED***
		left:     c.compileExpression(v.Left),
		right:    c.compileExpression(v.Right),
		operator: v.Operator,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileLogicalOr(left, right ast.Expression, idx file.Idx) compiledExpr ***REMOVED***
	r := &compiledLogicalOr***REMOVED***
		left:  c.compileExpression(left),
		right: c.compileExpression(right),
	***REMOVED***
	r.init(c, idx)
	return r
***REMOVED***

func (c *compiler) compileLogicalAnd(left, right ast.Expression, idx file.Idx) compiledExpr ***REMOVED***
	r := &compiledLogicalAnd***REMOVED***
		left:  c.compileExpression(left),
		right: c.compileExpression(right),
	***REMOVED***
	r.init(c, idx)
	return r
***REMOVED***

func (e *compiledVariableExpr) emitGetter(putOnStack bool) ***REMOVED***
	if e.initializer != nil ***REMOVED***
		idExpr := &compiledIdentifierExpr***REMOVED***
			name: e.name,
		***REMOVED***
		idExpr.init(e.c, file.Idx(0))
		idExpr.emitSetter(e.initializer, putOnStack)
	***REMOVED*** else ***REMOVED***
		if putOnStack ***REMOVED***
			e.c.emit(loadUndef)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileVariableExpression(v *ast.VariableExpression) compiledExpr ***REMOVED***
	r := &compiledVariableExpr***REMOVED***
		name:        v.Name,
		initializer: c.compileExpression(v.Initializer),
	***REMOVED***
	if fn, ok := r.initializer.(*compiledFunctionLiteral); ok ***REMOVED***
		fn.lhsName = v.Name
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledObjectLiteral) emitGetter(putOnStack bool) ***REMOVED***
	e.addSrcMap()
	e.c.emit(newObject)
	for _, prop := range e.expr.Value ***REMOVED***
		keyExpr := e.c.compileExpression(prop.Key)
		cl, ok := keyExpr.(*compiledLiteral)
		if !ok ***REMOVED***
			e.c.throwSyntaxError(e.offset, "non-literal properties in object literal are not supported yet")
		***REMOVED***
		key := cl.val.string()
		valueExpr := e.c.compileExpression(prop.Value)
		if fn, ok := valueExpr.(*compiledFunctionLiteral); ok ***REMOVED***
			if fn.expr.Name == nil ***REMOVED***
				fn.lhsName = key
			***REMOVED***
		***REMOVED***
		valueExpr.emitGetter(true)
		switch prop.Kind ***REMOVED***
		case "value":
			if key == __proto__ ***REMOVED***
				e.c.emit(setProto)
			***REMOVED*** else ***REMOVED***
				e.c.emit(setProp1(key))
			***REMOVED***
		case "method":
			e.c.emit(setProp1(key))
		case "get":
			e.c.emit(setPropGetter(key))
		case "set":
			e.c.emit(setPropSetter(key))
		default:
			panic(fmt.Errorf("unknown property kind: %s", prop.Kind))
		***REMOVED***
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileObjectLiteral(v *ast.ObjectLiteral) compiledExpr ***REMOVED***
	r := &compiledObjectLiteral***REMOVED***
		expr: v,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledArrayLiteral) emitGetter(putOnStack bool) ***REMOVED***
	e.addSrcMap()
	objCount := 0
	for _, v := range e.expr.Value ***REMOVED***
		if v != nil ***REMOVED***
			e.c.compileExpression(v).emitGetter(true)
			objCount++
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadNil)
		***REMOVED***
	***REMOVED***
	if objCount == len(e.expr.Value) ***REMOVED***
		e.c.emit(newArray(objCount))
	***REMOVED*** else ***REMOVED***
		e.c.emit(&newArraySparse***REMOVED***
			l:        len(e.expr.Value),
			objCount: objCount,
		***REMOVED***)
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileArrayLiteral(v *ast.ArrayLiteral) compiledExpr ***REMOVED***
	r := &compiledArrayLiteral***REMOVED***
		expr: v,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledRegexpLiteral) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		pattern, err := compileRegexp(e.expr.Pattern, e.expr.Flags)
		if err != nil ***REMOVED***
			e.c.throwSyntaxError(e.offset, err.Error())
		***REMOVED***

		e.c.emit(&newRegexp***REMOVED***pattern: pattern, src: newStringValue(e.expr.Pattern)***REMOVED***)
	***REMOVED***
***REMOVED***

func (c *compiler) compileRegexpLiteral(v *ast.RegExpLiteral) compiledExpr ***REMOVED***
	r := &compiledRegexpLiteral***REMOVED***
		expr: v,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledCallExpr) emitGetter(putOnStack bool) ***REMOVED***
	var calleeName unistring.String
	switch callee := e.callee.(type) ***REMOVED***
	case *compiledDotExpr:
		callee.left.emitGetter(true)
		e.c.emit(dup)
		e.c.emit(getPropCallee(callee.name))
	case *compiledBracketExpr:
		callee.left.emitGetter(true)
		e.c.emit(dup)
		callee.member.emitGetter(true)
		e.c.emit(getElemCallee)
	case *compiledIdentifierExpr:
		calleeName = callee.name
		callee.emitGetterAndCallee()
	default:
		e.c.emit(loadUndef)
		callee.emitGetter(true)
	***REMOVED***

	for _, expr := range e.args ***REMOVED***
		expr.emitGetter(true)
	***REMOVED***

	e.addSrcMap()
	if calleeName == "eval" ***REMOVED***
		foundFunc := false
		for sc := e.c.scope; sc != nil; sc = sc.outer ***REMOVED***
			if !foundFunc && sc.function ***REMOVED***
				foundFunc = true
				sc.thisNeeded, sc.argsNeeded = true, true
				if !sc.strict ***REMOVED***
					sc.dynamic = true
				***REMOVED***
			***REMOVED***
			sc.dynLookup = true
		***REMOVED***

		if e.c.scope.strict ***REMOVED***
			e.c.emit(callEvalStrict(len(e.args)))
		***REMOVED*** else ***REMOVED***
			e.c.emit(callEval(len(e.args)))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		e.c.emit(call(len(e.args)))
	***REMOVED***

	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledCallExpr) deleteExpr() compiledExpr ***REMOVED***
	r := &defaultDeleteExpr***REMOVED***
		expr: e,
	***REMOVED***
	r.init(e.c, file.Idx(e.offset+1))
	return r
***REMOVED***

func (c *compiler) compileCallExpression(v *ast.CallExpression) compiledExpr ***REMOVED***

	args := make([]compiledExpr, len(v.ArgumentList))
	for i, argExpr := range v.ArgumentList ***REMOVED***
		args[i] = c.compileExpression(argExpr)
	***REMOVED***

	r := &compiledCallExpr***REMOVED***
		args:   args,
		callee: c.compileExpression(v.Callee),
	***REMOVED***
	r.init(c, v.LeftParenthesis)
	return r
***REMOVED***

func (c *compiler) compileIdentifierExpression(v *ast.Identifier) compiledExpr ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierName(v.Name, int(v.Idx)-1)
	***REMOVED***

	r := &compiledIdentifierExpr***REMOVED***
		name: v.Name,
	***REMOVED***
	r.offset = int(v.Idx) - 1
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileNumberLiteral(v *ast.NumberLiteral) compiledExpr ***REMOVED***
	if c.scope.strict && octalRegexp.MatchString(v.Literal) ***REMOVED***
		c.throwSyntaxError(int(v.Idx)-1, "Octal literals are not allowed in strict mode")
		panic("Unreachable")
	***REMOVED***
	var val Value
	switch num := v.Value.(type) ***REMOVED***
	case int64:
		val = intToValue(num)
	case float64:
		val = floatToValue(num)
	default:
		panic(fmt.Errorf("Unsupported number literal type: %T", v.Value))
	***REMOVED***
	r := &compiledLiteral***REMOVED***
		val: val,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileStringLiteral(v *ast.StringLiteral) compiledExpr ***REMOVED***
	r := &compiledLiteral***REMOVED***
		val: stringValueFromRaw(v.Value),
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileBooleanLiteral(v *ast.BooleanLiteral) compiledExpr ***REMOVED***
	var val Value
	if v.Value ***REMOVED***
		val = valueTrue
	***REMOVED*** else ***REMOVED***
		val = valueFalse
	***REMOVED***

	r := &compiledLiteral***REMOVED***
		val: val,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileAssignExpression(v *ast.AssignExpression) compiledExpr ***REMOVED***
	// log.Printf("compileAssignExpression(): %+v", v)

	r := &compiledAssignExpr***REMOVED***
		left:     c.compileExpression(v.Left),
		right:    c.compileExpression(v.Right),
		operator: v.Operator,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledEnumGetExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.c.emit(enumGet)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***
