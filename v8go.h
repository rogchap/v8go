#ifndef V8GO_H
#define V8GO_H
#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>
#include <stdint.h>

typedef void* IsolatePtr;
typedef void* ContextPtr;
typedef void* ValuePtr;

typedef struct {
  const char* msg;
  const char* location;
  const char* stack;
} RtnError;

typedef struct {
  ValuePtr value;
  RtnError error;
} RtnValue;

typedef struct {
  size_t total_heap_size;
  size_t total_heap_size_executable;
  size_t total_physical_size;
  size_t total_available_size;
  size_t used_heap_size;
  size_t heap_size_limit;
  size_t malloced_memory;
  size_t external_memory;
  size_t peak_malloced_memory;
  size_t number_of_native_contexts;
  size_t number_of_detached_contexts;
} IsolateHStatistics;

extern void Init();
extern IsolatePtr NewIsolate();
extern void IsolateDispose(IsolatePtr ptr);
extern void IsolateTerminateExecution(IsolatePtr ptr);
extern IsolateHStatistics IsolationGetHeapStatistics(IsolatePtr ptr);

extern ContextPtr NewContext(IsolatePtr prt);
extern void ContextDispose(ContextPtr ptr);
extern RtnValue RunScript(ContextPtr ctx_ptr,
                          const char* source,
                          const char* origin);

extern void ValueDispose(ValuePtr ptr);
const char* ValueToString(ValuePtr ptr);
const uint32_t* ValueToArrayIndex(ValuePtr ptr);
int ValueToBoolean(ValuePtr ptr);
int32_t ValueToInt32(ValuePtr ptr);
int ValueIsUndefined(ValuePtr ptr);
int ValueIsNull(ValuePtr ptr);
int ValueIsNullOrUndefined(ValuePtr ptr);
int ValueIsTrue(ValuePtr ptr);
int ValueIsFalse(ValuePtr ptr);
int ValueIsName(ValuePtr ptr);
int ValueIsString(ValuePtr ptr);
int ValueIsSymbol(ValuePtr ptr);
int ValueIsFunction(ValuePtr ptr);
int ValueIsObject(ValuePtr ptr);
int ValueIsBigInt(ValuePtr ptr);
int ValueIsBoolean(ValuePtr ptr);
int ValueIsNumber(ValuePtr ptr);
int ValueIsExternal(ValuePtr ptr);
int ValueIsInt32(ValuePtr ptr);
int ValueIsUint32(ValuePtr ptr);
int ValueIsDate(ValuePtr ptr);
int ValueIsArgumentsObject(ValuePtr ptr);
int ValueIsBigIntObject(ValuePtr ptr);
int ValueIsNumberObject(ValuePtr ptr);
int ValueIsStringObject(ValuePtr ptr);
int ValueIsSymbolObject(ValuePtr ptr);
int ValueIsNativeError(ValuePtr ptr);
int ValueIsRegExp(ValuePtr ptr);
int ValueIsAsyncFunction(ValuePtr ptr);
int ValueIsGeneratorFunction(ValuePtr ptr);
int ValueIsGeneratorObject(ValuePtr ptr);
int ValueIsPromise(ValuePtr ptr);
int ValueIsMap(ValuePtr ptr);
int ValueIsSet(ValuePtr ptr);
int ValueIsMapIterator(ValuePtr ptr);
int ValueIsSetIterator(ValuePtr ptr);
int ValueIsWeakMap(ValuePtr ptr);
int ValueIsWeakSet(ValuePtr ptr);
int ValueIsArray(ValuePtr ptr);
int ValueIsArrayBuffer(ValuePtr ptr);
int ValueIsArrayBufferView(ValuePtr ptr);
int ValueIsTypedArray(ValuePtr ptr);
int ValueIsUint8Array(ValuePtr ptr);
int ValueIsUint8ClampedArray(ValuePtr ptr);
int ValueIsInt8Array(ValuePtr ptr);
int ValueIsUint16Array(ValuePtr ptr);
int ValueIsInt16Array(ValuePtr ptr);
int ValueIsUint32Array(ValuePtr ptr);
int ValueIsInt32Array(ValuePtr ptr);
int ValueIsFloat32Array(ValuePtr ptr);
int ValueIsFloat64Array(ValuePtr ptr);
int ValueIsBigInt64Array(ValuePtr ptr);
int ValueIsBigUint64Array(ValuePtr ptr);
int ValueIsDataView(ValuePtr ptr);
int ValueIsSharedArrayBuffer(ValuePtr ptr);
int ValueIsProxy(ValuePtr ptr);
int ValueIsWasmModuleObject(ValuePtr ptr);
int ValueIsModuleNamespaceObject(ValuePtr ptr);

const char* Version();

#ifdef __cplusplus
}  // extern "C"
#endif
#endif  // V8GO_H
