/**
 * @kind path-problem
 * @id 2
 */

 import go
 import semmle.go.dataflow.TaintTracking
 import DataFlow::PathGraph

 FuncDef getAReachableFunction(FuncDef rootFunc) {
     result = rootFunc
     or
     exists(DataFlow::CallNode callNode |
         callNode.getRoot() = rootFunc
         and
         result = getAReachableFunction(callNode.getACallee())
     )
 }

 class TaintedFieldReadSource extends DataFlow::FieldReadNode {
     TaintedFieldReadSource() {
         exists(DataFlow::FieldReadNode fieldReadNode |
             this = fieldReadNode
             and
             fieldReadNode.getField().getName().regexpMatch("(?i).*email.*")
             and
             fieldReadNode.getRoot() = getAReachableFunction(any(FuncDef fd | fd.getName() = "CreateEmail"))
         )
     }
 }

 class RequestFieldAssignSink extends DataFlow::Node {
     string fieldName;
     DataFlow::Node base;

     RequestFieldAssignSink() {
         exists(Field field |
             field.getName() = fieldName
             and
             field.getAWrite().writesField(base, field, this)
             and
             base.getTypeBound().getName().regexpMatch("(?i)(.*Req.*)")
         )
     }

     string getFieldName() { result = fieldName }

     DataFlow::Node getBase() { result = base }
 }


 class FieldAssignConfig extends TaintTracking2::Configuration {
     FieldAssignConfig() { this = "FieldAssignConfig" }

     override predicate isSource(DataFlow::Node node) {
         node instanceof TaintedFieldReadSource
     }

     override predicate isSink(DataFlow::Node node) {
         node instanceof RequestFieldAssignSink
     }

     override predicate isAdditionalTaintStep(DataFlow::Node fromNode, DataFlow::Node toNode) {
         any(DataFlow::Write w).writesComponent(toNode.(DataFlow::PostUpdateNode).getPreUpdateNode(),  fromNode)
     }
 }


 class RpcCallSource extends DataFlow::Node {
     TaintedFieldReadSource taintedFieldReadSource;
     RequestFieldAssignSink requestFieldAssignSink;

     RpcCallSource() {
         exists(FieldAssignConfig fieldAssignConfig |
             fieldAssignConfig.hasFlow(taintedFieldReadSource, requestFieldAssignSink)
             and
             exists(DataFlow::Node node | node = requestFieldAssignSink.getBase() |
                 this = node
                 or
                 exists(SsaDefinition ssadef |
                     ssadef.getSourceVariable().getARead() = node
                     and
                     this = ssadef.getSourceVariable().getARead()
                     and
                     node.asInstruction().getASuccessor*() = this.asInstruction()
                 )
             )
         )
     }

     RequestFieldAssignSink getRequestFieldAssignSink() {
         result = requestFieldAssignSink
     }
 }

 private class RpcParamSink extends DataFlow::Node {
     RpcParamSink() {
         exists(DataFlow::CallNode callNode |
             callNode.getTarget().getPackage().getPath().regexpMatch("(.*sdk.*)|(.*client.*)")
             and
             callNode.getAnArgument() = this
       )
     }

     override string toString() {
       result = this.getPackagePath() + ":" + this.getMethodName()
     }

     string getPackagePath() {
       result = this.(DataFlow::ArgumentNode).getCall().(DataFlow::CallNode).getTarget().getPackage().getPath()
     }

     string getMethodName() {
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

 from
     DataFlow::PathNode source,
     DataFlow::PathNode sink,
     RpcCallSource rpcCallSource,
     RpcParamSink rpcParamSink,
     RpcCallConfig config

 where
     config.hasFlowPath(source, sink)
     and
     source.getNode() = rpcCallSource
     and
     sink.getNode() = rpcParamSink

 select
     sink, source, sink, "Field flows to $@ field, which is used in call to $@.",
     rpcCallSource.getRequestFieldAssignSink(), rpcCallSource.getRequestFieldAssignSink().getFieldName(),
     rpcParamSink, rpcParamSink.toString()

