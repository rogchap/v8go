#ifndef V8GO_H
#define V8GO_H
#ifdef __cplusplus
extern "C" {
#endif

typedef void* IsolatePtr;
typedef void* ContextPtr;

extern void Init();
extern IsolatePtr NewIsolate();
extern void IsolateDispose(IsolatePtr ptr);
extern void TerminateExecution(IsolatePtr ptr);

extern ContextPtr NewContext(IsolatePtr prt);
extern void RunScript(ContextPtr ctx_ptr, const char* source, const char* origin);

const char* Version();

#ifdef __cplusplus
}  // extern "C"
#endif
#endif  // V8GO_H
