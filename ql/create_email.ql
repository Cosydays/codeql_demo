// Known: the field `Email` of method `CreateEmail` is tainted
// Analysis goal: the `Email` passed to rpc_sdk's RpcCreateEmail's `NewEmail`

import go
import semmle.go.dataflow.TaintTracking3
import semmle.go.dataflow.DataFlow3
import DataFlow::PathGraph

module CommonModule {
    class FunctionModels extends TaintTracking::FunctionModel {
        FunctionInput functionInput;
        FunctionOutput functionOutput;

        FunctionModels() {
            functionInput.isParameter([0 .. 7]) and
            functionOutput.isResult([0 .. 7])
        }

        override predicate hasTaintFlow(FunctionInput input, FunctionOutput output) {
          input = functionInput and output = functionOutput
        }
    }
}

class TaintedFieldReadSource extends DataFlow::FieldReadNode {
    TaintedFieldReadSource() {
        exists(DataFlow::FieldReadNode fieldReadNode |
            this = fieldReadNode
            and
            fieldReadNode.getField().getName().regexpMatch("(?i).*email.*")
            and
            fieldReadNode.getRoot() = getAReachableFunction(any(FuncDef fd | fd.getName() = "CreateEmail")))
    }
}

FuncDef getAReachableFunction(FuncDef rootFunc) {
    result = rootFunc
    or
    exists(DataFlow::CallNode callNode |
        callNode.getRoot() = rootFunc
        and
        result = getAReachableFunction(callNode.getACallee())
    )
}


class FieldAssignSink extends DataFlow::Node {
    string fieldName;
    DataFlow::Node base;

    FieldAssignSink() {
        exists(Field field |
            field.getName() = fieldName
            and
            field.getAWrite().writesField(base, field, this)
            and
            base.getTypeBound().getName().regexpMatch("(.*Request.*)|(.*Req.*)"))
    }

    string getFieldName() { result = fieldName }

    string getFieldType() {result = this.getBase().getTypeBound().getName()}

    DataFlow::Node getBase() { result = base }
}

class FieldAssignConfig extends TaintTracking2::Configuration {
    FieldAssignConfig() {
        this = "FieldAssignConfig"
    }

    override predicate isSource(DataFlow::Node node) {
        node instanceof TaintedFieldReadSource
    }

    override predicate isSink(DataFlow::Node node) {
        node instanceof FieldAssignSink
    }

    override predicate isAdditionalTaintStep(DataFlow::Node fromNode, DataFlow::Node toNode) {
        any(DataFlow::Write w).writesComponent(toNode.(DataFlow::PostUpdateNode).getPreUpdateNode(),  fromNode)
    }
}

class RpcCallSource extends DataFlow::Node {
    TaintedFieldReadSource taintedFieldReadSource;
    FieldAssignSink fieldAssignSink;

    RpcCallSource() {
        exists(FieldAssignConfig fieldAssignConfig |
            fieldAssignConfig.hasFlow(taintedFieldReadSource, fieldAssignSink)
            and
            exists(DataFlow::Node node | node = fieldAssignSink.getBase() |
                this = node
                or
                exists(SsaDefinition ssadef |
                    ssadef.getSourceVariable().getARead() = node
                    and
                    this = ssadef.getSourceVariable().getARead()
                    and
                    node.asInstruction().getASuccessor*() = this.asInstruction()
                  )))
    }

    FieldAssignSink getFieldAssignSink() {
        result = fieldAssignSink
    }
}

private class ThriftSink extends DataFlow::Node {
    ThriftSink() {
      exists(DataFlow::CallNode callNode |
        callNode.getTarget().getPackage().getPath().regexpMatch("(.*sdk.*)|(.*client.*)")
        and callNode.getAnArgument() = this
      )
    }

    override string toString() {
      result = this.getPackagePath() + ":" + this.getFunc()
    }

    string getPackagePath() {
      result = this.(DataFlow::ArgumentNode).getCall().(DataFlow::CallNode).getTarget().getPackage().getPath()
    }

    string getFunc() {
      result = this.(DataFlow::ArgumentNode).getCall().(DataFlow::CallNode).getTarget().getName()
    }
}

class RpcCallConfig extends TaintTracking::Configuration {
    RpcCallConfig() {
        this = "RpcCallConfig"
    }

    override predicate isSource(DataFlow::Node node) {
        node instanceof RpcCallSource
    }

    override predicate isSink(DataFlow::Node node) {
        node instanceof ThriftSink
    }

    override predicate isAdditionalTaintStep(DataFlow::Node fromNode, DataFlow::Node toNode) {
        any(DataFlow::Write w).writesComponent(toNode.(DataFlow::PostUpdateNode).getPreUpdateNode(),  fromNode)
    }
}

from
    RpcCallConfig rpcCallConfig,
    DataFlow::PathNode source,
    RpcCallSource rpcCallSource,
    DataFlow::PathNode sink,
    ThriftSink thriftSink

where
    rpcCallConfig.hasFlowPath(source, sink)
    and
    source.getNode() = rpcCallSource
    and
    sink.getNode() = thriftSink

select
    rpcCallSource.getFieldAssignSink().getFieldName(),
    thriftSink.toString(),
    thriftSink.asExpr().getLocation().toString()
