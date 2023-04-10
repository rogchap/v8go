// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef V8GO_H
#define V8GO_H
#ifdef __cplusplus

#include "libplatform/libplatform.h"
#include "v8-profiler.h"
#include "v8.h"

typedef v8::Isolate* IsolatePtr;
typedef v8::CpuProfiler* CpuProfilerPtr;
typedef v8::CpuProfile* CpuProfilePtr;
typedef const v8::CpuProfileNode* CpuProfileNodePtr;
typedef v8::ScriptCompiler::CachedData* ScriptCompilerCachedDataPtr;

extern "C" {
#else
// Opaque to cgo, but useful to treat it as a pointer to a distinct type
typedef struct v8Isolate v8Isolate;
typedef v8Isolate* IsolatePtr;

typedef struct v8CpuProfiler v8CpuProfiler;
typedef v8CpuProfiler* CpuProfilerPtr;

typedef struct v8CpuProfile v8CpuProfile;
typedef v8CpuProfile* CpuProfilePtr;

typedef struct v8CpuProfileNode v8CpuProfileNode;
typedef const v8CpuProfileNode* CpuProfileNodePtr;

typedef struct v8ScriptCompilerCachedData v8ScriptCompilerCachedData;
typedef const v8ScriptCompilerCachedData* ScriptCompilerCachedDataPtr;
#endif

// Opaque to both C and C++
typedef struct v8BackingStore v8BackingStore;
typedef v8BackingStore* BackingStorePtr;

#include <stddef.h>
#include <stdint.h>

// ScriptCompiler::CompileOptions values
extern const int ScriptCompilerNoCompileOptions;
extern const int ScriptCompilerConsumeCodeCache;
extern const int ScriptCompilerEagerCompile;

typedef struct m_ctx m_ctx;
typedef struct m_value m_value;
typedef struct m_template m_template;
typedef struct m_unboundScript m_unboundScript;

typedef m_ctx* ContextPtr;
typedef m_value* ValuePtr;
typedef m_template* TemplatePtr;
typedef m_unboundScript* UnboundScriptPtr;

typedef struct {
  const char* msg;
  const char* location;
  const char* stack;
} RtnError;

typedef struct {
  UnboundScriptPtr ptr;
  int cachedDataRejected;
  RtnError error;
} RtnUnboundScript;

typedef struct {
  ScriptCompilerCachedDataPtr ptr;
  const uint8_t* data;
  int length;
  int rejected;
} ScriptCompilerCachedData;

typedef struct {
  ScriptCompilerCachedData cachedData;
  int compileOption;
} CompileOptions;

typedef struct {
  CpuProfilerPtr ptr;
  IsolatePtr iso;
} CPUProfiler;

typedef struct CPUProfileNode {
  CpuProfileNodePtr ptr;
  unsigned nodeId;
  int scriptId;
  const char* scriptResourceName;
  const char* functionName;
  int lineNumber;
  int columnNumber;
  unsigned hitCount;
  const char* bailoutReason;
  int childrenCount;
  struct CPUProfileNode** children;
} CPUProfileNode;

typedef struct {
  CpuProfilePtr ptr;
  const char* title;
  CPUProfileNode* root;
  int64_t startTime;
  int64_t endTime;
} CPUProfile;

typedef struct {
  ValuePtr value;
  RtnError error;
} RtnValue;

typedef struct {
  const char* data;
  int length;
  RtnError error;
} RtnString;

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

typedef struct {
  const uint64_t* word_array;
  int word_count;
  int sign_bit;
} ValueBigInt;

extern void Init();
extern IsolatePtr NewIsolate();
extern void IsolatePerformMicrotaskCheckpoint(IsolatePtr ptr);
extern void IsolateDispose(IsolatePtr ptr);
extern void IsolateTerminateExecution(IsolatePtr ptr);
extern int IsolateIsExecutionTerminating(IsolatePtr ptr);
extern IsolateHStatistics IsolationGetHeapStatistics(IsolatePtr ptr);

extern ValuePtr IsolateThrowException(IsolatePtr iso, ValuePtr value);

extern RtnUnboundScript IsolateCompileUnboundScript(IsolatePtr iso_ptr,
                                                    const char* source,
                                                    const char* origin,
                                                    CompileOptions options);
extern ScriptCompilerCachedData* UnboundScriptCreateCodeCache(
    IsolatePtr iso_ptr,
    UnboundScriptPtr us_ptr);
extern void ScriptCompilerCachedDataDelete(
    ScriptCompilerCachedData* cached_data);
extern RtnValue UnboundScriptRun(ContextPtr ctx_ptr, UnboundScriptPtr us_ptr);

extern CPUProfiler* NewCPUProfiler(IsolatePtr iso_ptr);
extern void CPUProfilerDispose(CPUProfiler* ptr);
extern void CPUProfilerStartProfiling(CPUProfiler* ptr, const char* title);
extern CPUProfile* CPUProfilerStopProfiling(CPUProfiler* ptr,
                                            const char* title);
extern void CPUProfileDelete(CPUProfile* ptr);

extern ContextPtr NewContext(IsolatePtr iso_ptr,
                             TemplatePtr global_template_ptr,
                             int ref);
extern int ContextRetainedValueCount(ContextPtr ctx);
extern void ContextFree(ContextPtr ptr);
extern RtnValue RunScript(ContextPtr ctx_ptr,
                          const char* source,
                          const char* origin);
extern RtnValue JSONParse(ContextPtr ctx_ptr, const char* str);
const char* JSONStringify(ContextPtr ctx_ptr, ValuePtr val_ptr);
extern ValuePtr ContextGlobal(ContextPtr ctx_ptr);

extern void TemplateFreeWrapper(TemplatePtr ptr);
extern void TemplateSetValue(TemplatePtr ptr,
                             const char* name,
                             ValuePtr val_ptr,
                             int attributes);
extern void TemplateSetTemplate(TemplatePtr ptr,
                                const char* name,
                                TemplatePtr obj_ptr,
                                int attributes);

extern TemplatePtr NewObjectTemplate(IsolatePtr iso_ptr);
extern RtnValue ObjectTemplateNewInstance(TemplatePtr ptr, ContextPtr ctx_ptr);
extern void ObjectTemplateSetInternalFieldCount(TemplatePtr ptr,
                                                int field_count);
extern int ObjectTemplateInternalFieldCount(TemplatePtr ptr);

extern TemplatePtr NewFunctionTemplate(IsolatePtr iso_ptr, int callback_ref);
extern RtnValue FunctionTemplateGetFunction(TemplatePtr ptr,
                                            ContextPtr ctx_ptr);

extern ValuePtr NewValueNull(IsolatePtr iso_ptr);
extern ValuePtr NewValueUndefined(IsolatePtr iso_ptr);
extern ValuePtr NewValueInteger(IsolatePtr iso_ptr, int32_t v);
extern ValuePtr NewValueIntegerFromUnsigned(IsolatePtr iso_ptr, uint32_t v);
extern RtnValue NewValueString(IsolatePtr iso_ptr, const char* v, int v_length);
extern ValuePtr NewValueBoolean(IsolatePtr iso_ptr, int v);
extern ValuePtr NewValueNumber(IsolatePtr iso_ptr, double v);
extern ValuePtr NewValueBigInt(IsolatePtr iso_ptr, int64_t v);
extern ValuePtr NewValueBigIntFromUnsigned(IsolatePtr iso_ptr, uint64_t v);
extern RtnValue NewValueBigIntFromWords(IsolatePtr iso_ptr,
                                        int sign_bit,
                                        int word_count,
                                        const uint64_t* words);
void ValueRelease(ValuePtr ptr);
extern RtnString ValueToString(ValuePtr ptr);
const uint32_t* ValueToArrayIndex(ValuePtr ptr);
int ValueToBoolean(ValuePtr ptr);
int32_t ValueToInt32(ValuePtr ptr);
int64_t ValueToInteger(ValuePtr ptr);
double ValueToNumber(ValuePtr ptr);
RtnString ValueToDetailString(ValuePtr ptr);
uint32_t ValueToUint32(ValuePtr ptr);
extern ValueBigInt ValueToBigInt(ValuePtr ptr);
extern RtnValue ValueToObject(ValuePtr ptr);
int ValueSameValue(ValuePtr ptr, ValuePtr otherPtr);
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

extern void ObjectSet(ValuePtr ptr, const char* key, ValuePtr val_ptr);
extern void ObjectSetIdx(ValuePtr ptr, uint32_t idx, ValuePtr val_ptr);
extern int ObjectSetInternalField(ValuePtr ptr, int idx, ValuePtr val_ptr);
extern int ObjectInternalFieldCount(ValuePtr ptr);
extern RtnValue ObjectGet(ValuePtr ptr, const char* key);
extern RtnValue ObjectGetIdx(ValuePtr ptr, uint32_t idx);
extern ValuePtr ObjectGetInternalField(ValuePtr ptr, int idx);
int ObjectHas(ValuePtr ptr, const char* key);
int ObjectHasIdx(ValuePtr ptr, uint32_t idx);
int ObjectDelete(ValuePtr ptr, const char* key);
int ObjectDeleteIdx(ValuePtr ptr, uint32_t idx);

extern RtnValue NewPromiseResolver(ContextPtr ctx_ptr);
extern ValuePtr PromiseResolverGetPromise(ValuePtr ptr);
int PromiseResolverResolve(ValuePtr ptr, ValuePtr val_ptr);
int PromiseResolverReject(ValuePtr ptr, ValuePtr val_ptr);
int PromiseState(ValuePtr ptr);
RtnValue PromiseThen(ValuePtr ptr, int callback_ref);
RtnValue PromiseThen2(ValuePtr ptr, int on_fulfilled_ref, int on_rejected_ref);
RtnValue PromiseCatch(ValuePtr ptr, int callback_ref);
extern ValuePtr PromiseResult(ValuePtr ptr);

extern RtnValue FunctionCall(ValuePtr ptr,
                             ValuePtr recv,
                             int argc,
                             ValuePtr argv[]);
RtnValue FunctionNewInstance(ValuePtr ptr, int argc, ValuePtr args[]);
ValuePtr FunctionSourceMapUrl(ValuePtr ptr);

const char* Version();
extern void SetFlags(const char* flags);

extern BackingStorePtr SharedArrayBufferGetBackingStore(ValuePtr ptr);
extern void BackingStoreRelease(BackingStorePtr ptr);
extern void* BackingStoreData(BackingStorePtr ptr);
extern size_t BackingStoreByteLength(BackingStorePtr ptr);

#ifdef __cplusplus
}  // extern "C"
#endif
#endif  // V8GO_H
