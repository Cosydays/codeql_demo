//分析字段流转
//结果为下游Field，pkgMethod，location
import go
import semmle.go.dataflow.TaintTracking3
import semmle.go.dataflow.DataFlow3
import DataFlow::PathGraph

module KitexCommon {
    //对所有的Func都进行污点追踪
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

// 将 MethodA 的第二个参数作为source
class TaintedParam extends DataFlow::Node {
    TaintedParam() {
        any(Parameter parameter |
            parameter.getIndex() = 1
            and
            parameter.getFunction().getName() = "CreateEmail")
        =
        this.asParameter()
    }
}

// 读取字段的Node
class TaintedFieldRead extends DataFlow::FieldReadNode {
    TaintedFieldRead() {
        this.getField().getName().regexpMatch("Email")
    }
}

// Map的Node
class TaintedMapAccess extends IndexExpr {
    TaintedMapAccess() {
        this.getIndex().(StringLit).getValue() = "Email"
    }
}

//从MethodA的第一个参数到字段读取/Map读取
class TaintedStructParamFlowConfig extends TaintTracking3::Configuration {
    TaintedStructParamFlowConfig() {
        this = "TaintedStructParamFlowConfig"
    }

    override predicate isSource(DataFlow::Node node) {
        node instanceof TaintedParam
    }

    override predicate isSink(DataFlow::Node node) {
        node instanceof TaintedFieldRead
        or
        node.asExpr() instanceof TaintedMapAccess
    }

    override predicate isAdditionalTaintStep(DataFlow::Node fromNode, DataFlow::Node toNode) {
        any(DataFlow::Write w).writesComponent(toNode.(DataFlow::PostUpdateNode).getPreUpdateNode(),  fromNode)
      }
}

//第一次污点追踪：Method的param到struct的field
class TaintedStructMemberSource extends DataFlow::Node {
    TaintedParam paramSource;

    TaintedStructMemberSource() {
        exists(TaintedStructParamFlowConfig config, DataFlow::Node sink |
            config.hasFlow(paramSource, sink)
            and
            this.asExpr() = sink.asExpr())
    }

    TaintedParam getParamSource() { result = paramSource }
}

// 字段赋值的终点
class FieldAssignSink extends DataFlow::Node {
    string fieldName;
    DataFlow::Node base;

    FieldAssignSink() {
        exists(Field field |
            field.getName() = fieldName
            and
            field.getAWrite().writesField(base, field, this)
            and
            // base.getTypeBound().getName().regexpMatch("(.*Request.*)|(.*Req.*)"))
            field.getQualifiedName().regexpMatch("(.*sdk.*)"))
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
        node instanceof TaintedStructMemberSource
    }

    override predicate isSink(DataFlow::Node node) {
        node instanceof FieldAssignSink
    }

    override predicate isAdditionalTaintStep(DataFlow::Node fromNode, DataFlow::Node toNode) {
        any(DataFlow::Write w).writesComponent(toNode.(DataFlow::PostUpdateNode).getPreUpdateNode(),  fromNode)
    }
}

// 第二次污点追踪
class RpcCallSource extends DataFlow::Node {
    TaintedStructMemberSource taintedStructMemberSource;
    FieldAssignSink fieldAssignSink;

    RpcCallSource() {
        exists(FieldAssignConfig fieldAssignConfig |
            fieldAssignConfig.hasFlow(taintedStructMemberSource, fieldAssignSink)
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

    TaintedStructMemberSource getTaintedStructMemberSource() {
        result = taintedStructMemberSource
    }

    FieldAssignSink getFieldAssignSink() {
        result = fieldAssignSink
    }
}

// RpcClient的参数为终点
class RpcParamSink extends DataFlow::Node {
    RpcParamSink() {
        exists(DataFlow::CallNode callNode |
            callNode.getTarget().getPackage().getPath().regexpMatch("(.*sdk.*)|(.*client.*)")
            // and
            // callNode.getArgument(1).getType().getPackage().getPath().regexpMatch("(.*thrift_gen.*)|(.*kitex_gen.*)")
            and
            callNode.getArgument(1) = this)
    }

    string getPkgFunc() {
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
        node instanceof RpcParamSink
    }

    override predicate isAdditionalTaintStep(DataFlow::Node fromNode, DataFlow::Node toNode) {
        any(DataFlow::Write w).writesComponent(toNode.(DataFlow::PostUpdateNode).getPreUpdateNode(),  fromNode)
    }
}

// 第三次污点追踪
from
    DataFlow::PathNode source,
    RpcCallSource rpcCallSource,
    DataFlow::PathNode sink,
    RpcParamSink rpcParamSink,
    RpcCallConfig rpcCallConfig

where
    rpcCallConfig.hasFlowPath(source, sink)
    and
    source.getNode() = rpcCallSource
    and
    sink.getNode() = rpcParamSink

select
    rpcCallSource.getTaintedStructMemberSource().getParamSource(),
    rpcCallSource.getFieldAssignSink(),
    rpcCallSource.getFieldAssignSink().getFieldName(),
    rpcCallSource.getTaintedStructMemberSource(),
    rpcParamSink.getPkgFunc(),
    rpcParamSink.asExpr().getLocation().toString()