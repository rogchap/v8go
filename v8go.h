#ifndef V8GO_H
#define V8GO_H
#ifdef __cplusplus
extern "C" {
#endif

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

extern void Init();
extern IsolatePtr NewIsolate();
extern void IsolateDispose(IsolatePtr ptr);
extern void IsolateTerminateExecution(IsolatePtr ptr);

extern ContextPtr NewContext(IsolatePtr prt);
extern void ContextDispose(ContextPtr ptr);
extern RtnValue RunScript(ContextPtr ctx_ptr, const char* source, const char* origin);

extern void ValueDispose(ValuePtr ptr);
const char* ValueToString(ValuePtr ptr);

const char* Version();

#ifdef __cplusplus
}  // extern "C"
#endif
#endif  // V8GO_H
