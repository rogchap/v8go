// Copyright 2021 the V8 project authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef INCLUDE_V8_EMBEDDER_HEAP_H_
#define INCLUDE_V8_EMBEDDER_HEAP_H_

#include <stddef.h>
#include <stdint.h>

#include <utility>
#include <vector>

#include "cppgc/common.h"
#include "v8-local-handle.h"   // NOLINT(build/include_directory)
#include "v8-traced-handle.h"  // NOLINT(build/include_directory)
#include "v8config.h"          // NOLINT(build/include_directory)

namespace v8 {

class Data;
class Isolate;
class Value;

namespace internal {
class LocalEmbedderHeapTracer;
}  // namespace internal

/**
 * Handler for embedder roots on non-unified heap garbage collections.
 */
class V8_EXPORT EmbedderRootsHandler {
 public:
  virtual ~EmbedderRootsHandler() = default;

  /**
   * Returns true if the |TracedReference| handle should be considered as root
   * for the currently running non-tracing garbage collection and false
   * otherwise. The default implementation will keep all |TracedReference|
   * references as roots.
   *
   * If this returns false, then V8 may decide that the object referred to by
   * such a handle is reclaimed. In that case, V8 calls |ResetRoot()| for the
   * |TracedReference|.
   *
   * Note that the `handle` is different from the handle that the embedder holds
   * for retaining the object. The embedder may use |WrapperClassId()| to
   * distinguish cases where it wants handles to be treated as roots from not
   * being treated as roots.
   */
  virtual bool IsRoot(const v8::TracedReference<v8::Value>& handle) = 0;

  /**
   * Used in combination with |IsRoot|. Called by V8 when an
   * object that is backed by a handle is reclaimed by a non-tracing garbage
   * collection. It is up to the embedder to reset the original handle.
   *
   * Note that the |handle| is different from the handle that the embedder holds
   * for retaining the object. It is up to the embedder to find the original
   * handle via the object or class id.
   */
  virtual void ResetRoot(const v8::TracedReference<v8::Value>& handle) = 0;
};

/**
 * Interface for tracing through the embedder heap. During a V8 garbage
 * collection, V8 collects hidden fields of all potential wrappers, and at the
 * end of its marking phase iterates the collection and asks the embedder to
 * trace through its heap and use reporter to report each JavaScript object
 * reachable from any of the given wrappers.
 */
class V8_EXPORT
// GCC doesn't like combining __attribute__(()) with [[deprecated]].
#ifdef __clang__
V8_DEPRECATED("Use CppHeap when working with v8::TracedReference.")
#endif  // __clang__
    EmbedderHeapTracer {
 public:
  using EmbedderStackState = cppgc::EmbedderStackState;

  enum TraceFlags : uint64_t {
    kNoFlags = 0,
    kReduceMemory = 1 << 0,
    kForced = 1 << 2,
  };

  /**
   * Interface for iterating through |TracedReference| handles.
   */
  class V8_EXPORT TracedGlobalHandleVisitor {
   public:
    virtual ~TracedGlobalHandleVisitor() = default;
    virtual void VisitTracedReference(const TracedReference<Value>& handle) {}
  };

  /**
   * Summary of a garbage collection cycle. See |TraceEpilogue| on how the
   * summary is reported.
   */
  struct TraceSummary {
    /**
     * Time spent managing the retained memory in milliseconds. This can e.g.
     * include the time tracing through objects in the embedder.
     */
    double time = 0.0;

    /**
     * Memory retained by the embedder through the |EmbedderHeapTracer|
     * mechanism in bytes.
     */
    size_t allocated_size = 0;
  };

  virtual ~EmbedderHeapTracer() = default;

  /**
   * Iterates all |TracedReference| handles created for the |v8::Isolate| the
   * tracer is attached to.
   */
  void IterateTracedGlobalHandles(TracedGlobalHandleVisitor* visitor);

  /**
   * Called by the embedder to set the start of the stack which is e.g. used by
   * V8 to determine whether handles are used from stack or heap.
   */
  void SetStackStart(void* stack_start);

  /**
   * Called by v8 to register internal fields of found wrappers.
   *
   * The embedder is expected to store them somewhere and trace reachable
   * wrappers from them when called through |AdvanceTracing|.
   */
  virtual void RegisterV8References(
      const std::vector<std::pair<void*, void*>>& embedder_fields) = 0;

  void RegisterEmbedderReference(const BasicTracedReference<v8::Data>& ref);

  /**
   * Called at the beginning of a GC cycle.
   */
  virtual void TracePrologue(TraceFlags flags) {}

  /**
   * Called to advance tracing in the embedder.
   *
   * The embedder is expected to trace its heap starting from wrappers reported
   * by RegisterV8References method, and report back all reachable wrappers.
   * Furthermore, the embedder is expected to stop tracing by the given
   * deadline. A deadline of infinity means that tracing should be finished.
   *
   * Returns |true| if tracing is done, and false otherwise.
   */
  virtual bool AdvanceTracing(double deadline_in_ms) = 0;

  /*
   * Returns true if there no more tracing work to be done (see AdvanceTracing)
   * and false otherwise.
   */
  virtual bool IsTracingDone() = 0;

  /**
   * Called at the end of a GC cycle.
   *
   * Note that allocation is *not* allowed within |TraceEpilogue|. Can be
   * overriden to fill a |TraceSummary| that is used by V8 to schedule future
   * garbage collections.
   */
  virtual void TraceEpilogue(TraceSummary* trace_summary) {}

  /**
   * Called upon entering the final marking pause. No more incremental marking
   * steps will follow this call.
   */
  virtual void EnterFinalPause(EmbedderStackState stack_state) = 0;

  /*
   * Called by the embedder to request immediate finalization of the currently
   * running tracing phase that has been started with TracePrologue and not
   * yet finished with TraceEpilogue.
   *
   * Will be a noop when currently not in tracing.
   *
   * This is an experimental feature.
   */
  void FinalizeTracing();

  /**
   * See documentation on EmbedderRootsHandler.
   */
  virtual bool IsRootForNonTracingGC(
      const v8::TracedReference<v8::Value>& handle);

  /**
   * See documentation on EmbedderRootsHandler.
   */
  virtual void ResetHandleInNonTracingGC(
      const v8::TracedReference<v8::Value>& handle);

  /*
   * Called by the embedder to signal newly allocated or freed memory. Not bound
   * to tracing phases. Embedders should trade off when increments are reported
   * as V8 may consult global heuristics on whether to trigger garbage
   * collection on this change.
   */
  void IncreaseAllocatedSize(size_t bytes);
  void DecreaseAllocatedSize(size_t bytes);

  /*
   * Returns the v8::Isolate this tracer is attached too and |nullptr| if it
   * is not attached to any v8::Isolate.
   */
  v8::Isolate* isolate() const { return v8_isolate_; }

 protected:
  v8::Isolate* v8_isolate_ = nullptr;

  friend class internal::LocalEmbedderHeapTracer;
};

}  // namespace v8

#endif  // INCLUDE_V8_EMBEDDER_HEAP_H_
