/**
 * @name Untrusted CreateEmail result passed to RpcCreateEmail
 * @kind path-problem
 * @id go/unsafe-create-email-result-used-in-rpc
 */

import go
import semmle.go.dataflow.TaintTracking
import DataFlow::PathGraph

abstract class RPCFieldAssignDataSource extends DataFlow::Node {
  abstract string getDescription();
}

/**
 * A `FieldReadNode` with the field's base type's name containing `Req`,
 * and the field name containing (case-insensitive) `Email`.
 */
class EmailRequestFieldReadSource extends RPCFieldAssignDataSource, DataFlow::FieldReadNode {
  EmailRequestFieldReadSource() {
    this.getField().getName().regexpMatch("(?i).*email.*") and
    this.getBase().getType().getName().regexpMatch(".*Req.*")
  }

  override string getDescription() { result = this.getFieldName() + " field" }
}

/**
 * The result of a `CallNode` to `GetRedisValue` with the argument `"id"`
 */
class RedisIdFieldReadSource extends RPCFieldAssignDataSource, DataFlow::Node {
  RedisIdFieldReadSource() {
    exists(DataFlow::CallNode cn |
      this = cn.getResult() and
      cn.getTarget().getName() = "GetRedisValue" and
      // improve this with local data-flow to also find calls where the arg is a var
      cn.getArgument(0).asExpr().getStringValue() = "id"
    )
  }

  override string getDescription() { result = "Redis 'id' field" }
}

/**
 * The result of a `CallNode` to `GetHttpData`
 */
class GetHttpDataFieldReadSource extends RPCFieldAssignDataSource, DataFlow::Node {
  GetHttpDataFieldReadSource() {
    exists(DataFlow::CallNode cn |
      this = cn.getResult() and
      cn.getTarget().getName() = "GetHttpData"
    )
  }

  override string getDescription() { result = "GetHttpData call result" }
}

/**
 * A data-flow node that is the RHS of an assignment to a field
 */
class RequestFieldAssignNode extends DataFlow::Node {
  string fieldName;
  DataFlow::Node base;

  RequestFieldAssignNode() {
    exists(Field field |
      fieldName = field.getName() and
      field.getAWrite().writesField(base, field, this)
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
    any(RPCFieldAssignToRpcConfiguration config)
        .hasFlow(this.getBase(), subsequentCall.getAnArgument())
  }

  DataFlow::CallNode getSubsequentCallNode() { result = subsequentCall }
}

class RPCFieldAssignToRpcConfiguration extends TaintTracking2::Configuration {
  RPCFieldAssignToRpcConfiguration() { this = "RPCFieldAssignToRpcConfiguration" }

  override predicate isSource(DataFlow::Node source) {
    source = any(RequestFieldAssignNode node).getBase()
  }

  override predicate isSink(DataFlow::Node sink) {
    sink =
      any(DataFlow::CallNode node |
        node.getTarget().getPackage().getPath().regexpMatch("(.*sdk.*)|(.*client.*)")
      ).getAnArgument()
  }

  override predicate isAdditionalTaintStep(DataFlow::Node node1, DataFlow::Node node2) {
    any(DataFlow::Write w)
        .writesComponent(node2.(DataFlow::PostUpdateNode).getPreUpdateNode(), node1)
  }
}

/**
 * A taint-tracking configuration representing flow from a `RPCFieldAssignDataSource`
 * to a `RPCFieldAssignSink` with additional taint steps from field writes to their base.
 */
class CreateEmailToRpcFieldAssignSinkConfiguration extends TaintTracking::Configuration {
  CreateEmailToRpcFieldAssignSinkConfiguration() { this = "CreateEmailToRpcConfiguration" }

  override predicate isSource(DataFlow::Node source) { source instanceof RPCFieldAssignDataSource }

  override predicate isSink(DataFlow::Node sink) { sink instanceof RPCFieldAssignSink }

  override predicate isAdditionalTaintStep(DataFlow::Node node1, DataFlow::Node node2) {
    any(DataFlow::Write w)
        .writesComponent(node2.(DataFlow::PostUpdateNode).getPreUpdateNode(), node1)
  }

  override predicate isSanitizerOut(DataFlow::Node node) { isSink(node) }
}

from
  DataFlow::PathNode source, DataFlow::PathNode sink, RPCFieldAssignDataSource sourceFieldReadNode,
  RPCFieldAssignSink sinkFieldAssignNode
where
  any(CreateEmailToRpcFieldAssignSinkConfiguration config).hasFlowPath(source, sink) and
  sourceFieldReadNode = source.getNode() and
  sinkFieldAssignNode = sink.getNode()
select sink, source, sink, "Unsanitized $@ flows to $@ field, which is used in call to $@.",
  sourceFieldReadNode, sourceFieldReadNode.getDescription(), sinkFieldAssignNode,
  sinkFieldAssignNode.getFieldName(), sinkFieldAssignNode.getSubsequentCallNode(),
  sinkFieldAssignNode.getSubsequentCallNode().getACallee().getName()
