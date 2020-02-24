package config

import (
	"os"
	"testing"
)

func TestCompleteConfig(t *testing.T) {

	var config = `
project:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - analyzer: test-analyzer
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
    - analyzer: test-analyzer-2
      name: "The Testalyzer 2"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
        testfile: "/the/test"
  reporting:
    - reporter: test-reporter
      name: "The test reporter"
      config:
        tester: "Endocode"
`
	masterconf, err := ReadConfigFromBytes([]byte(config))
	if err != nil {
		t.Logf("Broken config %v", err)
		t.FailNow()
	}
	err = ValidateConfig(masterconf)
	if err != nil {
		t.Logf("Config validation failed: %v", err)
		t.Fail()
	}
	for _, ana := range masterconf.Analysis {
		if ana.TrustLevel == 0 {
			t.Logf("Analyzer trust level is 0: %s", ana.Name)
			t.Fail()
		}
	}

	projNode := CreateProjectNode(masterconf)
	value := projNode.GetMetaData("Vendor", "")
	if value != "Endocode" {
		t.Logf("Failed to create project node; %v not Endocode", value)
		t.Fail()
	}
}

func TestMissingModuleInstanceName(t *testing.T) {

	var config = `
project:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - analyzer: test-analyzer
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
    - analyzer: test-analyzer-2
      name: "The Testalyzer 2"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
        testfile: "/the/test"
  reporting:
    - reporter: test-reporter
      config:
        tester: "Endocode"
`
	masterconfig, err := ReadConfigFromBytes([]byte(config))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = ValidateConfig(masterconfig)
	if err == nil || err.Error() != "1. reporter misconfigured Name invalid" {
		t.Log(err)
		t.Fail()
	}
}

func TestDuplicateModuleInstanceNames(t *testing.T) {

	var config = `
project:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - analyzer: test-analyzer
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
    - analyzer: test-analyzer-2
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
        testfile: "/the/test"
  reporting:
    - reporter: test-reporter
      name: "The reporter"
      config:
        tester: "Endocode"
`
	masterconfig, err := ReadConfigFromBytes([]byte(config))
	if err != nil {
		t.Log(err)
	}

	err = ValidateConfig(masterconfig)
	if err == nil || err.Error() != "2. analyzer misconfigured duplicate value of The Testalyzer in Name" {
		t.Log(err)
		t.Fail()
	}
}

func TestMissingAnalyzerExecutableName(t *testing.T) {

	var config = `
project:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
    - analyzer: test-analyzer-2
      name: "The Testalyzer 2"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
        testfile: "/the/test"
  reporting:
    - reporter: test-reporter
      name: "The reporter"
      config:
        tester: "Endocode"
`
	masterconfig, err := ReadConfigFromBytes([]byte(config))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = ValidateConfig(masterconfig)
	if err == nil || err.Error() != "1. analyzer misconfigured Analyzer invalid" {
		t.Log(err)
		t.Fail()
	}
}

func TestMissingRepoterExecutableName(t *testing.T) {

	var config = `
project:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - name: "The Testalyzer"
      analyzer: test-analyzer
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
    - analyzer: test-analyzer
      name: "The Testalyzer 2"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
        testfile: "/the/test"
  reporting:
    - config:
        tester: "Endocode"
`
	masterconfig, err := ReadConfigFromBytes([]byte(config))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = ValidateConfig(masterconfig)
	if err == nil || err.Error() != "1. reporter misconfigured Name invalid" {
		t.Log(err)
		t.Fail()
	}
}

func TestImplicitDuplicatePosixNameConfig(t *testing.T) {

	var config = `
project:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - analyzer: test-analyzer
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
    - analyzer: test-analyzer-2
      name: "The_Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
        testfile: "/the/test"
  reporting:
    - reporter: test-reporter
      name: "The test reporter"
      config:
        tester: "Endocode"
`
	masterconfig, err := ReadConfigFromBytes([]byte(config))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = ValidateConfig(masterconfig)
	if err == nil || err.Error() != "2. analyzer misconfigured duplicate value of The_Testalyzer in PosixName" {
		t.Log(err)
		t.Fail()
	}
}

func TestRPCServerAddress(t *testing.T) {

	var config = `
project:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: "12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - analyzer: test-analyzer
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
  reporting:
    - reporter: test-reporter
      name: "The test reporter"
      config:
        tester: "Endocode"
`
	masterconfig, err := ReadConfigFromBytes([]byte(config))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = ValidateConfig(masterconfig)
	if err == nil || err.Error() != "Invalid RPC address" {
		t.Log(err)
		t.Fail()
	}
}

func TestConfigEnvOverride(t *testing.T) {

	var config = `
project:
  name: "The Test"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":12345"
    dbaddress: "testhost:54321"
    dbworkers: 4
  analysis:
    - analyzer: test-analyzer
      name: "The Testalyzer"
      selector: sourcecode
      pathsub:
        - old: "/the/path"
          new: "/buildroot"
      config:
        workdir: "/buildroot"
  reporting:
    - reporter: test-reporter
      name: "The test reporter"
      config:
        tester: "Endocode"
`
	os.Setenv("SERVER_DBADDRESS", "override:12345")
	os.Setenv("SERVER_RPCADDRESS", ":54321")
	os.Setenv("SERVER_BUILDPATH", "/override")

	masterConfig, _ := ReadConfigFromBytes([]byte(config))
	ConfigEnvOverride(masterConfig)

	if masterConfig.Server.DBAddress != os.Getenv("SERVER_DBADDRESS") ||
		masterConfig.Server.RPCAddress != os.Getenv("SERVER_RPCADDRESS") ||
		masterConfig.Server.BuildPath != os.Getenv("SERVER_BUILDPATH") {
		t.Log("Configuration override failed.")
		t.Fail()
	}
}
