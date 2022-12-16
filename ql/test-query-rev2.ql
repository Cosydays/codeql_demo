/**
 * @name Untrusted CreateEmail result passed to RpcCreateEmail
 * @kind path-problem
 * @id go/unsafe-create-email-result-used-in-rpc
 */

import go
import semmle.go.dataflow.TaintTracking
import DataFlow::PathGraph

/**
 * A `FieldReadNode` in a function that could be transitively called
 * by `CreateEmail`, with the field's base type's name containing `Req`,
 * and the field name containing (case-insensitive) `Email`.
 */
class CreateEmailRequestEmailFieldReadSource extends DataFlow::FieldReadNode {
  CreateEmailRequestEmailFieldReadSource() {
    this.getField().getName().regexpMatch("(?i).*email.*") and
    this.getBase().getType().getName().regexpMatch(".*Req.*")
  }
}

/**
 * A data-flow node that is the RHS of an assignment to a field with
 * a base type name containing the substring "Req" (e.g. "Request")
 */
class RequestFieldAssignNode extends DataFlow::Node {
  string fieldName;
  DataFlow::Node base;

  RequestFieldAssignNode() {
    exists(Field field |
      fieldName = field.getName() and
      field.getAWrite().writesField(base, field, this) and
      base.getTypeBound().getName().regexpMatch(".*Req.*")
    )
  }

  string getFieldName() { result = fieldName }

  DataFlow::Node getBase() { result = base }
}

/**
 * A `RequestFieldAssignNode` that flows to a subsequent RPC call
 */
class RPCFieldAssignSink extends RequestFieldAssignNode {
  DataFlow::CallNode subsequentCall;

  RPCFieldAssignSink() {
    TaintTracking::localTaint(this.getBase(), subsequentCall.getAnArgument()) and
    subsequentCall.getTarget().getPackage().getPath().regexpMatch("(.*sdk.*)|(.*client.*)")
  }

  DataFlow::CallNode getSubsequentCallNode() { result = subsequentCall }
}

/**
 * A taint-tracking configuration representing flow from a `CreateEmailRequestEmailFieldReadSource`
 * to a `RPCFieldAssignSink` with additional taint steps from field writes to their base.
 */
class CreateEmailToRpcConfiguration extends TaintTracking::Configuration {
  CreateEmailToRpcConfiguration() { this = "CreateEmailToRpcConfiguration" }

  override predicate isSource(DataFlow::Node source) {
    source instanceof CreateEmailRequestEmailFieldReadSource
  }

  override predicate isSink(DataFlow::Node sink) { sink instanceof RPCFieldAssignSink }

  override predicate isAdditionalTaintStep(DataFlow::Node node1, DataFlow::Node node2) {
    any(DataFlow::Write w)
        .writesComponent(node2.(DataFlow::PostUpdateNode).getPreUpdateNode(), node1)
  }
}

from
  DataFlow::PathNode source, DataFlow::PathNode sink, DataFlow::FieldReadNode sourceFieldReadNode,
  RPCFieldAssignSink sinkFieldAssignNode
where
  any(CreateEmailToRpcConfiguration config).hasFlowPath(source, sink) and
  sourceFieldReadNode = source.getNode() and
  sinkFieldAssignNode = sink.getNode()
select sink, source, sink, "Unsanitized $@ field flows to $@ field, which is used in call to $@.",
  sourceFieldReadNode, sourceFieldReadNode.getField().getName(), sinkFieldAssignNode,
  sinkFieldAssignNode.getFieldName(), sinkFieldAssignNode.getSubsequentCallNode(),
  sinkFieldAssignNode.getSubsequentCallNode().getACallee().getName()
