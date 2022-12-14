/**
 * @name Untrusted CreateEmail result passed to RpcCreateEmail
 * @kind path-problem
 * @id go/unsafe-create-email-result-used-in-rpc
 */

import go
import semmle.go.dataflow.TaintTracking
import DataFlow::PathGraph

private FuncDef getAReachableFunction(FuncDef rootFunc) {
  result = rootFunc
  or
  exists(DataFlow::CallNode callNode |
    callNode.getRoot() = rootFunc and
    result = getAReachableFunction(callNode.getACallee())
  )
}

class GetEmailFunctionModel extends TaintTracking::FunctionModel {
  GetEmailFunctionModel() {
    // TODO: correct qualified name
    this.hasQualifiedName(_, "GetEmail")
  }

  override predicate hasTaintFlow(FunctionInput input, FunctionOutput output) {
    input.isParameter(0) and output.isResult()
  }
}

class CreateEmailRequestEmailFieldReadSource extends DataFlow::FieldReadNode {
  CreateEmailRequestEmailFieldReadSource() {
    exists(DataFlow::FieldReadNode fieldReadNode |
      this = fieldReadNode and
      fieldReadNode.getField().getName().regexpMatch("(?i).*email.*") and
      fieldReadNode.getBase().getType().getName() = "CreateEmailRequest" and
      fieldReadNode.getRoot() =
        getAReachableFunction(any(FuncDef fd | fd.getName() = "CreateEmail"))
    )
  }
}

class CreateEmailToRpcConfiguration extends TaintTracking::Configuration {
  CreateEmailToRpcConfiguration() { this = "CreateEmailToRpcConfiguration" }

  override predicate isSource(DataFlow::Node source) {
    source instanceof CreateEmailRequestEmailFieldReadSource
  }

  override predicate isSink(DataFlow::Node sink) {
    exists(DataFlow::CallNode callNode |
      callNode.getTarget().getPackage().getPath().regexpMatch("(.*sdk.*)|(.*client.*)") and
      callNode.getAnArgument() = sink
    )
  }

  override predicate isSanitizer(DataFlow::Node node) {
    exists(DataFlow::CallNode cn |
      cn.getACallee().getName() = ["NormalizeEmail", "FormatInt"] and
      cn.getAnArgument() = node
    )
  }

  override predicate isAdditionalTaintStep(DataFlow::Node node1, DataFlow::Node node2) {
    any(DataFlow::Write w)
        .writesComponent(node2.(DataFlow::PostUpdateNode).getPreUpdateNode(), node1)
  }
}

from
  DataFlow::PathNode source, DataFlow::PathNode sink, DataFlow::FieldReadNode sourceFieldReadNode,
  DataFlow::CallNode sinkCallNode
where
  any(CreateEmailToRpcConfiguration config).hasFlowPath(source, sink) and
  sourceFieldReadNode = source.getNode() and
  sinkCallNode.getAnArgument() = sink.getNode()
select sink, source, sink, "Unsanitized $@ field flows to $@ call.", sourceFieldReadNode,
  sourceFieldReadNode.getField().getName(), sinkCallNode.getACallee(),
  sinkCallNode.getACallee().getName()
