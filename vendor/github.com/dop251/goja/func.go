package goja

import (
	"reflect"

	"github.com/dop251/goja/unistring"
)

type baseFuncObject struct ***REMOVED***
	baseObject

	nameProp, lenProp valueProperty
***REMOVED***

type funcObject struct ***REMOVED***
	baseFuncObject

	stash *stash
	prg   *Program
	src   string
***REMOVED***

type nativeFuncObject struct ***REMOVED***
	baseFuncObject

	f         func(FunctionCall) Value
	construct func(args []Value, newTarget *Object) *Object
***REMOVED***

type boundFuncObject struct ***REMOVED***
	nativeFuncObject
	wrapped *Object
***REMOVED***

func (f *nativeFuncObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return f.f
***REMOVED***

func (f *nativeFuncObject) exportType() reflect.Type ***REMOVED***
	return reflect.TypeOf(f.f)
***REMOVED***

func (f *funcObject) _addProto(n unistring.String) Value ***REMOVED***
	if n == "prototype" ***REMOVED***
		if _, exists := f.values[n]; !exists ***REMOVED***
			return f.addPrototype()
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (f *funcObject) getStr(p unistring.String, receiver Value) Value ***REMOVED***
	return f.getStrWithOwnProp(f.getOwnPropStr(p), p, receiver)
***REMOVED***

func (f *funcObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if v := f._addProto(name); v != nil ***REMOVED***
		return v
	***REMOVED***

	return f.baseObject.getOwnPropStr(name)
***REMOVED***

func (f *funcObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	f._addProto(name)
	return f.baseObject.setOwnStr(name, val, throw)
***REMOVED***

func (f *funcObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return f._setForeignStr(name, f.getOwnPropStr(name), val, receiver, throw)
***REMOVED***

func (f *funcObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	f._addProto(name)
	return f.baseObject.deleteStr(name, throw)
***REMOVED***

func (f *funcObject) addPrototype() Value ***REMOVED***
	proto := f.val.runtime.NewObject()
	proto.self._putProp("constructor", f.val, true, false, true)
	return f._putProp("prototype", proto, true, false, false)
***REMOVED***

func (f *funcObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if r := f.baseObject.hasOwnPropertyStr(name); r ***REMOVED***
		return true
	***REMOVED***

	if name == "prototype" ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (f *funcObject) ownKeys(all bool, accum []Value) []Value ***REMOVED***
	if all ***REMOVED***
		if _, exists := f.values["prototype"]; !exists ***REMOVED***
			accum = append(accum, asciiString("prototype"))
		***REMOVED***
	***REMOVED***
	return f.baseFuncObject.ownKeys(all, accum)
***REMOVED***

func (f *funcObject) construct(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		newTarget = f.val
	***REMOVED***
	proto := newTarget.self.getStr("prototype", nil)
	var protoObj *Object
	if p, ok := proto.(*Object); ok ***REMOVED***
		protoObj = p
	***REMOVED*** else ***REMOVED***
		protoObj = f.val.runtime.global.ObjectPrototype
	***REMOVED***

	obj := f.val.runtime.newBaseObject(protoObj, classObject).val
	ret := f.call(FunctionCall***REMOVED***
		This:      obj,
		Arguments: args,
	***REMOVED***, newTarget)

	if ret, ok := ret.(*Object); ok ***REMOVED***
		return ret
	***REMOVED***
	return obj
***REMOVED***

func (f *funcObject) Call(call FunctionCall) Value ***REMOVED***
	return f.call(call, nil)
***REMOVED***

func (f *funcObject) call(call FunctionCall, newTarget Value) Value ***REMOVED***
	vm := f.val.runtime.vm
	pc := vm.pc

	vm.stack.expand(vm.sp + len(call.Arguments) + 1)
	vm.stack[vm.sp] = f.val
	vm.sp++
	if call.This != nil ***REMOVED***
		vm.stack[vm.sp] = call.This
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp] = _undefined
	***REMOVED***
	vm.sp++
	for _, arg := range call.Arguments ***REMOVED***
		if arg != nil ***REMOVED***
			vm.stack[vm.sp] = arg
		***REMOVED*** else ***REMOVED***
			vm.stack[vm.sp] = _undefined
		***REMOVED***
		vm.sp++
	***REMOVED***

	vm.pc = -1
	vm.pushCtx()
	vm.args = len(call.Arguments)
	vm.prg = f.prg
	vm.stash = f.stash
	vm.newTarget = newTarget
	vm.pc = 0
	vm.run()
	vm.pc = pc
	vm.halt = false
	return vm.pop()
***REMOVED***

func (f *funcObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return f.Call
***REMOVED***

func (f *funcObject) exportType() reflect.Type ***REMOVED***
	return reflect.TypeOf(f.Call)
***REMOVED***

func (f *funcObject) assertCallable() (func(FunctionCall) Value, bool) ***REMOVED***
	return f.Call, true
***REMOVED***

func (f *funcObject) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	return f.construct
***REMOVED***

func (f *baseFuncObject) init(name unistring.String, length int) ***REMOVED***
	f.baseObject.init()

	if name != "" ***REMOVED***
		f.nameProp.configurable = true
		f.nameProp.value = stringValueFromRaw(name)
		f._put("name", &f.nameProp)
	***REMOVED***

	f.lenProp.configurable = true
	f.lenProp.value = valueInt(length)
	f._put("length", &f.lenProp)
***REMOVED***

func (f *baseFuncObject) hasInstance(v Value) bool ***REMOVED***
	if v, ok := v.(*Object); ok ***REMOVED***
		o := f.val.self.getStr("prototype", nil)
		if o1, ok := o.(*Object); ok ***REMOVED***
			for ***REMOVED***
				v = v.self.proto()
				if v == nil ***REMOVED***
					return false
				***REMOVED***
				if o1 == v ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			f.val.runtime.typeErrorResult(true, "prototype is not an object")
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (f *nativeFuncObject) defaultConstruct(ccall func(ConstructorCall) *Object, args []Value) *Object ***REMOVED***
	proto := f.getStr("prototype", nil)
	var protoObj *Object
	if p, ok := proto.(*Object); ok ***REMOVED***
		protoObj = p
	***REMOVED*** else ***REMOVED***
		protoObj = f.val.runtime.global.ObjectPrototype
	***REMOVED***
	obj := f.val.runtime.newBaseObject(protoObj, classObject).val
	ret := ccall(ConstructorCall***REMOVED***
		This:      obj,
		Arguments: args,
	***REMOVED***)

	if ret != nil ***REMOVED***
		return ret
	***REMOVED***
	return obj
***REMOVED***

func (f *nativeFuncObject) assertCallable() (func(FunctionCall) Value, bool) ***REMOVED***
	if f.f != nil ***REMOVED***
		return f.f, true
	***REMOVED***
	return nil, false
***REMOVED***

func (f *nativeFuncObject) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	return f.construct
***REMOVED***

func (f *boundFuncObject) getStr(p unistring.String, receiver Value) Value ***REMOVED***
	return f.getStrWithOwnProp(f.getOwnPropStr(p), p, receiver)
***REMOVED***

func (f *boundFuncObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if name == "caller" || name == "arguments" ***REMOVED***
		return f.val.runtime.global.throwerProperty
	***REMOVED***

	return f.nativeFuncObject.getOwnPropStr(name)
***REMOVED***

func (f *boundFuncObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if name == "caller" || name == "arguments" ***REMOVED***
		return true
	***REMOVED***
	return f.nativeFuncObject.deleteStr(name, throw)
***REMOVED***

func (f *boundFuncObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if name == "caller" || name == "arguments" ***REMOVED***
		panic(f.val.runtime.NewTypeError("'caller' and 'arguments' are restricted function properties and cannot be accessed in this context."))
	***REMOVED***
	return f.nativeFuncObject.setOwnStr(name, val, throw)
***REMOVED***

func (f *boundFuncObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return f._setForeignStr(name, f.getOwnPropStr(name), val, receiver, throw)
***REMOVED***

func (f *boundFuncObject) hasInstance(v Value) bool ***REMOVED***
	return instanceOfOperator(v, f.wrapped)
***REMOVED***
