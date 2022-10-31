//分析代码仓库内的redis的kp
//结果为kp，redisClientName，location
import go
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

class KpSource extends DataFlow::Node {
    KpSource() {
        exists(DataFlow::StringLit stringLit|
            stringLit.getValue().toString().regexpMatch("(.*%v.*)|(.*%s.*)")
            and
            this.getStringValue() = stringLit.getStringValue())
    }

    string getStr() {
        result = this.getStringValue()
    }

    DataFlow::Node getNode() {
        result = this
    }
}

class RedisSink extends DataFlow::Node {
    RedisSink() {
      exists(DataFlow::MethodCallNode methodCallNode |
        methodCallNode.getAnArgument() = this
        and
        methodCallNode.getTarget().getQualifiedName().regexpMatch("(.*redis.*Set.*)|(.*redis.*Del.*)")
      )
    }

    override string toString() {
      // 拿到clientname
      result = this.(DataFlow::ArgumentNode).getCall().(DataFlow::MethodCallNode).getReceiver().asExpr().(SelectorExpr).getSelector().toString()
    }
}

class Configuration extends TaintTracking::Configuration {
    Configuration() { this = "Configuration" }
    override predicate isSource(DataFlow::Node source) { source instanceof KpSource }
    override predicate isSink(DataFlow::Node sink) {sink instanceof RedisSink }

    // 字段赋值也认为是一种traint
    override predicate isAdditionalTaintStep(DataFlow::Node fromNode, DataFlow::Node toNode) {
      any(DataFlow::Write w).writesComponent(toNode.(DataFlow::PostUpdateNode).getPreUpdateNode(),  fromNode)
    }
}

from
	Configuration configuration,
	DataFlow::PathNode sourceNode,
	DataFlow::PathNode sinkNode
where
	configuration.hasFlowPath(sourceNode, sinkNode)
select
	sourceNode.getNode().getStringValue().toString(),
	sinkNode.getNode().toString(),
	sinkNode.getNode().asExpr().getLocation().toString()