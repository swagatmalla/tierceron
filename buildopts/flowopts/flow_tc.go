//go:build tc
// +build tc

package flowopts

import (
	tccutil "VaultConfig.TenantConfig/util/controller"
	flowcore "github.com/trimble-oss/tierceron/trcflow/core"
	trcf "github.com/trimble-oss/tierceron/trcflow/core/flowcorehelper"
)

func GetAdditionalFlows() []flowcore.FlowNameType {
	return tccutil.GetAdditionalFlows()
}

func GetAdditionalTestFlows() []flowcore.FlowNameType {
	return []flowcore.FlowNameType{} // Noop
}

func GetAdditionalFlowsByState(teststate string) []flowcore.FlowNameType {
	return []flowcore.FlowNameType{} // Noop
}

func ProcessFlowController(tfmContext *flowcore.TrcFlowMachineContext, trcFlowContext *flowcore.TrcFlowContext) error {
	return tccutil.ProcessFlowController(tfmContext, trcFlowContext)
}

func ProcessTestFlowController(tfmContext *flowcore.TrcFlowMachineContext, trcFlowContext *flowcore.TrcFlowContext) error {
	return tccutil.ProcessFlowController(tfmContext, trcFlowContext)
}

func GetFlowDatabaseName() string {
	return trcf.GetFlowDBName()
}